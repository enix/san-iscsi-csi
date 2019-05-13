package main

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"k8s.io/client-go/kubernetes"

	dothill "enix.io/dothill-api-go"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/controller"
)

type dothillProvisioner struct {
	dothillClient *dothill.Client
	kubeClient    *kubernetes.Clientset
}

// NewDothillProvisioner : Creates the provisionner instance that implements
// the controller.Provisioner interface
func NewDothillProvisioner(kubeClient *kubernetes.Clientset) controller.Provisioner {
	return &dothillProvisioner{
		dothillClient: &dothill.Client{},
		kubeClient:    kubeClient,
	}
}

func (p *dothillProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	size := options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]
	sizeStr := fmt.Sprintf("%sB", size.String())
	log.Printf("creating %s volume for host %s\n", sizeStr, options.Parameters[initiatorNameConfigKey])

	err := runPreflightChecks(options.Parameters, options.PVC.Spec.AccessModes)
	if err != nil {
		return nil, err
	}

	err = p.configureClient(options.Parameters)
	if err != nil {
		return nil, err
	}

	lun, err := p.chooseLUN(options.Parameters[initiatorNameConfigKey])
	if err != nil {
		return nil, err
	}

	dnsFormattedIQN := strings.ReplaceAll(options.Parameters[initiatorNameConfigKey], ":", ".")
	volumeName := fmt.Sprintf("%s.lun%d", dnsFormattedIQN, lun)
	_, _, err = p.dothillClient.CreateVolume(volumeName, sizeStr, options.Parameters[poolConfigKey])
	if err != nil {
		return nil, err
	}

	err = p.mapVolume(volumeName, options.Parameters[initiatorNameConfigKey], lun)
	if err != nil {
		p.dothillClient.DeleteVolume(volumeName)
		return nil, err
	}

	log.Printf("created volume %s (%s) for host %s\n", volumeName, sizeStr, options.Parameters[initiatorNameConfigKey])
	return generatePersistentVolume(volumeName, options.Parameters[initiatorNameConfigKey], lun, options), nil
}

func (p *dothillProvisioner) Delete(volume *v1.PersistentVolume) error {
	log.Printf("deleting volume %s\n", volume.ObjectMeta.Name)
	initiatorName := volume.ObjectMeta.Annotations[initiatorNameConfigKey]
	storageClassName := volume.ObjectMeta.Annotations[storageClassAnnotationKey]
	storageClass, err := p.kubeClient.StorageV1().StorageClasses().Get(storageClassName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	err = p.configureClient(storageClass.Parameters)
	if err != nil {
		return err
	}

	_, _, err = p.dothillClient.UnmapVolume(volume.ObjectMeta.Name, initiatorName)
	if err != nil {
		return err
	}

	_, _, err = p.dothillClient.DeleteVolume(volume.ObjectMeta.Name)
	if err != nil {
		return err
	}

	log.Printf("deleted volume %s\n", volume.ObjectMeta.Name)
	return nil
}

func (p *dothillProvisioner) configureClient(parameters map[string]string) error {
	credentials, err := p.kubeClient.CoreV1().Secrets("default").Get(parameters[credentialsSecretNameConfigKey], metav1.GetOptions{})
	if err != nil {
		return err
	}

	username := string(credentials.Data[usernameSecretKey])
	password := string(credentials.Data[passwordSecretKey])
	if p.dothillClient.Addr == parameters[apiAddressConfigKey] && p.dothillClient.Username == username {
		return nil
	}

	p.dothillClient.Username = username
	p.dothillClient.Password = password
	p.dothillClient.Addr = parameters[apiAddressConfigKey]

	err = p.dothillClient.Login()
	if err != nil {
		return err
	}

	return nil
}

func (p *dothillProvisioner) chooseLUN(initiatorName string) (int, error) {
	volumes, status, err := p.dothillClient.ShowHostMaps(initiatorName)
	if err != nil && status == nil {
		return -1, err
	}
	if status.ReturnCode == hostMapDoesNotExistsErrorCode {
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

func (p *dothillProvisioner) mapVolume(volumeName, initiatorName string, lun int) error {
	_, status, err := p.dothillClient.MapVolume(volumeName, initiatorName, "rw", lun)
	if err != nil && status == nil {
		return err
	}
	if status.ReturnCode == hostDoesNotExistsErrorCode {
		nodeName := strings.Split(initiatorName, ":")[1]
		_, _, err = p.dothillClient.CreateHost(nodeName, initiatorName)
		if err != nil {
			return err
		}
		_, _, err = p.dothillClient.MapVolume(volumeName, initiatorName, "rw", lun)
	}
	return err
}

func generatePersistentVolume(name, initiatorName string, lun int, options controller.VolumeOptions) *v1.PersistentVolume {
	portals := strings.Split(options.Parameters[portalsConfigKey], ",")
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
		return err
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
