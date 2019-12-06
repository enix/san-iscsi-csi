package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

// Provision : Called when a PVC is created
func (p *DothillController) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	defer p.dothillClient.HTTPClient.CloseIdleConnections()

	klog.V(2).Infof("Provision() called with: %+v", options)
	size := options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]
	sizeStr := fmt.Sprintf("%sB", size.String())
	klog.Infof("received %s volume request\n", sizeStr)

	err := runPreflightChecks(options.Parameters, options.PVC.Spec.AccessModes)
	if err != nil {
		return nil, err
	}

	err = p.configureClient(options.Parameters)
	if err != nil {
		return nil, err
	}

	lun, err := p.chooseLUN()
	if err != nil {
		return nil, err
	}
	klog.V(1).Infof("using LUN %d", lun)

	initiatorName := options.Parameters[initiatorNameConfigKey]
	overrideInitiatorName, overrideExists := options.PVC.Annotations[initiatorNameConfigKey]
	if overrideExists {
		initiatorName = overrideInitiatorName
		klog.V(1).Infof("custom initiator name was specified in PVC annotation: %s", initiatorName)
	} else if options.Parameters[uniqueInitiatorNameByPvcConfigKey] == "true" {
		year, month, _ := time.Now().Date()
		uniquePart := fmt.Sprintf("%d", rand.Int())[:8]
		initiatorName = fmt.Sprintf("iqn.%d-%02d.local.cluster:%s", year, int(month), uniquePart)
		klog.V(1).Infof("generated initiator name: %s", initiatorName)
	}

	iqnUniquePart := strings.Split(initiatorName, ":")[1]
	volumeName := fmt.Sprintf("%s.lun%d", iqnUniquePart, lun)
	klog.V(1).Infof("creating volume %s (size %s) in pool %s", volumeName, sizeStr, options.Parameters[poolConfigKey])
	_, _, err = p.dothillClient.CreateVolume(volumeName, sizeStr, options.Parameters[poolConfigKey])
	if err != nil {
		return nil, err
	}

	err = p.mapVolume(volumeName, initiatorName, lun)
	if err != nil {
		klog.Infof("volume %s couldn't be mapped, deleting it", volumeName)
		p.dothillClient.DeleteVolume(volumeName)
		return nil, err
	}

	klog.Infof("created volume %s (%s) for initiator %s (mapped on LUN %d)", volumeName, sizeStr, initiatorName, lun)
	pv := generatePersistentVolume(volumeName, initiatorName, lun, options)
	klog.V(2).Infof("created persitent volume %+v", pv)
	return pv, nil
}

// Delete : Called when a PVC is deleted
func (p *DothillController) Delete(volume *v1.PersistentVolume) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	defer p.dothillClient.HTTPClient.CloseIdleConnections()

	klog.V(2).Infof("Delete() called with: %+v", volume)
	klog.Infof("received delete request for volume %s", volume.ObjectMeta.Name)
	initiatorName := volume.ObjectMeta.Annotations[initiatorNameConfigKey]
	storageClassName := volume.ObjectMeta.Annotations[storageClassAnnotationKey]
	klog.V(1).Infof("fetching storage class %s", storageClassName)
	storageClass, err := p.kubeClient.StorageV1().StorageClasses().Get(storageClassName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	klog.V(2).Info(storageClass)

	if err = runPreflightChecks(storageClass.Parameters, nil); err != nil {
		return err
	}

	err = p.configureClient(storageClass.Parameters)
	if err != nil {
		return err
	}

	klog.V(1).Infof("unmapping volume %s from initiator %s", volume.ObjectMeta.Name, initiatorName)
	_, _, err = p.dothillClient.UnmapVolume(volume.ObjectMeta.Name, initiatorName)
	if err != nil {
		return err
	}

	klog.V(1).Infof("deleting volume %s", volume.ObjectMeta.Name)
	_, _, err = p.dothillClient.DeleteVolume(volume.ObjectMeta.Name)
	if err != nil {
		return err
	}

	klog.V(1).Infof("listing LUN mappings for %s", initiatorName)
	volumes, _, err := p.dothillClient.ShowHostMaps(initiatorName)
	if err != nil {
		return err
	}

	if len(volumes) == 0 {
		klog.V(1).Infof("no more mappings, deleting host %s", initiatorName)
		_, _, err := p.dothillClient.DeleteHost(initiatorName)
		if err != nil {
			return err
		}
		klog.V(1).Info("delete was successful")
	} else {
		klog.V(1).Infof("not deleting host %s as other mappings exists", initiatorName)
	}

	return nil
}

func (p *DothillController) configureClient(parameters map[string]string) error {
	klog.V(1).Infof("fetching dothill credentials from secret %s in namespace %s", parameters[credentialsSecretNameConfigKey], p.namespace)
	credentials, err := p.kubeClient.CoreV1().Secrets(p.namespace).Get(parameters[credentialsSecretNameConfigKey], metav1.GetOptions{})
	if err != nil {
		return err
	}

	username := string(credentials.Data[usernameSecretKey])
	password := string(credentials.Data[passwordSecretKey])
	klog.V(1).Infof("using dothill API at address %s", parameters[apiAddressConfigKey])
	if p.dothillClient.Addr == parameters[apiAddressConfigKey] && p.dothillClient.Username == username {
		klog.V(1).Info("dothill client is already configured for this API, skipping login")
		return nil
	}

	p.dothillClient.Username = username
	p.dothillClient.Password = password
	p.dothillClient.Addr = parameters[apiAddressConfigKey]

	klog.V(1).Infof("login into %s as user %s", p.dothillClient.Addr, p.dothillClient.Username)
	err = p.dothillClient.Login()
	if err != nil {
		return err
	}

	klog.V(1).Info("login was successful")
	return nil
}

func (p *DothillController) chooseLUN() (int, error) {
	klog.V(1).Infof("listing all LUN mappings")
	volumes, status, err := p.dothillClient.ShowHostMaps("")
	if err != nil && status == nil {
		return -1, err
	}
	if status.ReturnCode == hostMapDoesNotExistsErrorCode {
		klog.V(1).Info("initiator does not exists, assuming there is no LUN mappings yet and using LUN 1")
		return 1, nil
	}
	if err != nil {
		return -1, err
	}

	sort.Sort(Volumes(volumes))
	index := 0
	for ; index < len(volumes); index++ {
		if volumes[index].LUN != index+1 {
			return index + 1, nil
		}
	}

	if index+1 < maximumLUN {
		return index + 1, nil
	}

	return -1, errors.New("no more available LUNs")
}

func (p *DothillController) mapVolume(volumeName, initiatorName string, lun int) error {
	klog.V(1).Infof("trying to map volume %s for initiator %s on LUN %d", volumeName, initiatorName, lun)
	_, status, err := p.dothillClient.MapVolume(volumeName, initiatorName, "rw", lun)
	if err != nil && status == nil {
		return err
	}
	if status.ReturnCode == hostDoesNotExistsErrorCode {
		nodeName := strings.Split(initiatorName, ":")[1]
		klog.V(1).Infof("initiator does not exist, creating it with nickname %s", nodeName)
		_, _, err = p.dothillClient.CreateHost(nodeName, initiatorName)
		if err != nil {
			return err
		}
		klog.V(1).Info("retrying to map volume")
		_, _, err = p.dothillClient.MapVolume(volumeName, initiatorName, "rw", lun)
		if err != nil {
			return err
		}
	}

	klog.V(1).Info("mapping was successful")
	return nil
}

func generatePersistentVolume(name, initiatorName string, lun int, options controller.VolumeOptions) *v1.PersistentVolume {
	portals := strings.Split(options.Parameters[portalsConfigKey], ",")
	klog.V(1).Infof("generating persistent volume spec, ISCSI portals: %s", portals)

	return &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				initiatorNameConfigKey:    initiatorName,
				storageClassAnnotationKey: *options.PVC.Spec.StorageClassName,
			},
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			VolumeMode:                    options.PVC.Spec.VolumeMode,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				ISCSI: &v1.ISCSIPersistentVolumeSource{
					InitiatorName: &initiatorName,
					TargetPortal:  portals[0],
					Portals:       portals,
					IQN:           options.Parameters[targetIQNConfigKey],
					Lun:           int32(lun),
					FSType:        options.Parameters[fsTypeConfigKey],
					ReadOnly:      false,
				},
			},
		},
	}
}

func runPreflightChecks(parameters map[string]string, accessModes []v1.PersistentVolumeAccessMode) error {
	checkIfKeyExistsInConfig := func(key string) error {
		klog.V(2).Infof("checking for %s in storage class parameters", key)
		_, ok := parameters[key]
		if !ok {
			return fmt.Errorf("'%s' is missing from configuration", key)
		}
		return nil
	}

	if err := checkIfKeyExistsInConfig(fsTypeConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(poolConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(targetIQNConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(portalsConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(initiatorNameConfigKey); err != nil {
		if err2 := checkIfKeyExistsInConfig(uniqueInitiatorNameByPvcConfigKey); err2 != nil {
			return errors.Wrap(err, err2.Error())
		}
	}
	if err := checkIfKeyExistsInConfig(apiAddressConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(credentialsSecretNameConfigKey); err != nil {
		return err
	}

	for _, mode := range accessModes {
		if mode != v1.ReadWriteOnce {
			return errors.New("dothill storage only supports ReadWriteOnce access mode")
		}
	}

	return nil
}
