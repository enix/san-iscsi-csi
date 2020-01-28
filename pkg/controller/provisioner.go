package controller

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

const (
	fsTypeConfigKey                   = "fsType"
	poolConfigKey                     = "pool"
	targetIQNConfigKey                = "iqn"
	portalsConfigKey                  = "portals"
	initiatorNameConfigKey            = "initiatorName"
	apiAddressConfigKey               = "apiAddress"
	uniqueInitiatorNameByPvcConfigKey = "uniqueInitiatorNameByPvc"
	usernameSecretKey                 = "username"
	passwordSecretKey                 = "password"
	storageClassAnnotationKey         = "storageClass"

	maximumLUN                    = 255
	hostDoesNotExistsErrorCode    = -10386
	hostMapDoesNotExistsErrorCode = -10074
)

// CreateVolume creates a new volume from the given request. The function is
// idempotent.
func (driver *Driver) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	driver.lock.Lock()
	defer driver.lock.Unlock()
	defer driver.dothillClient.HTTPClient.CloseIdleConnections()

	klog.V(9).Infof("CreateVolume() called with: %+v", req)
	parameters := req.GetParameters()
	size := req.GetCapacityRange().GetRequiredBytes()
	sizeStr := fmt.Sprintf("%diB", size)
	klog.Infof("received %s volume request\n", sizeStr)

	err := runPreflightChecks(parameters, req.GetVolumeCapabilities())
	if err != nil {
		return nil, err
	}

	err = driver.configureClient(req.GetSecrets(), parameters[apiAddressConfigKey])
	if err != nil {
		return nil, err
	}

	lun, err := driver.chooseLUN()
	if err != nil {
		return nil, err
	}
	klog.Infof("using LUN %d", lun)

	initiatorName := parameters[initiatorNameConfigKey]
	// overrideInitiatorName, overrideExists := options.PVC.Annotations[initiatorNameConfigKey]
	// if overrideExists {
	// 	initiatorName = overrideInitiatorName
	// 	klog.Infof("custom initiator name was specified in PVC annotation: %s", initiatorName)
	// } else if options.Parameters[uniqueInitiatorNameByPvcConfigKey] == "true" {
	// 	year, month, _ := time.Now().Date()
	// 	uniquePart := fmt.Sprintf("%d", rand.Int())[:8]
	// 	initiatorName = fmt.Sprintf("iqn.%d-%02d.local.cluster:%s", year, int(month), uniquePart)
	// 	klog.Infof("generated initiator name: %s", initiatorName)
	// }

	iqnUniquePart := strings.Split(initiatorName, ":")[1]
	volumeName := fmt.Sprintf("%s.lun%d", iqnUniquePart, lun)
	klog.Infof("creating volume %s (size %s) in pool %s", volumeName, sizeStr, parameters[poolConfigKey])
	_, _, err = driver.dothillClient.CreateVolume(volumeName, sizeStr, parameters[poolConfigKey])
	if err != nil {
		return nil, err
	}

	err = driver.mapVolume(volumeName, initiatorName, lun)
	if err != nil {
		// klog.Infof("volume %s couldn't be mapped, deleting it", volumeName)
		// driver.dothillClient.DeleteVolume(volumeName)
		return nil, err
	}

	klog.Infof("created volume %s (%s) for initiator %s (mapped on LUN %d)", volumeName, sizeStr, initiatorName, lun)
	portals := strings.Split(req.GetParameters()[portalsConfigKey], ",")
	klog.Infof("generating volume spec, ISCSI portals: %s", portals)
	volume := &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      uuid.NewUUID().String(),
			CapacityBytes: req.GetCapacityRange().GetRequiredBytes(),
			VolumeContext: req.GetParameters(),
			ContentSource: req.GetVolumeContentSource(),
		},
	}

	klog.V(8).Infof("created volume %+v", volume)
	return volume, nil
}

// Delete : Called when a PVC is deleted
// func (driver *Driver) Delete(volume *v1.PersistentVolume) error {
// p.lock.Lock()
// defer p.lock.Unlock()
// defer p.dothillClient.HTTPClient.CloseIdleConnections()

// klog.V(2).Infof("Delete() called with: %+v", volume)
// klog.Infof("received delete request for volume %s", volume.ObjectMeta.Name)
// initiatorName := volume.ObjectMeta.Annotations[initiatorNameConfigKey]
// storageClassName := volume.ObjectMeta.Annotations[storageClassAnnotationKey]
// klog.Infof("fetching storage class %s", storageClassName)
// storageClass, err := p.kubeClient.StorageV1().StorageClasses().Get(storageClassName, metav1.GetOptions{})
// if err != nil {
// 	return err
// }
// klog.V(2).Info(storageClass)

// if err = runPreflightChecks(storageClass.Parameters, nil); err != nil {
// 	return err
// }

// err = p.configureClient(storageClass.Parameters)
// if err != nil {
// 	return err
// }

// klog.Infof("unmapping volume %s from initiator %s", volume.ObjectMeta.Name, initiatorName)
// _, _, err = p.dothillClient.UnmapVolume(volume.ObjectMeta.Name, initiatorName)
// if err != nil {
// 	return err
// }

// klog.Infof("deleting volume %s", volume.ObjectMeta.Name)
// _, _, err = p.dothillClient.DeleteVolume(volume.ObjectMeta.Name)
// if err != nil {
// 	return err
// }

// klog.Infof("listing LUN mappings for %s", initiatorName)
// volumes, _, err := p.dothillClient.ShowHostMaps(initiatorName)
// if err != nil {
// 	return err
// }

// if len(volumes) == 0 {
// 	klog.Infof("no more mappings, deleting host %s", initiatorName)
// 	_, _, err := p.dothillClient.DeleteHost(initiatorName)
// 	if err != nil {
// 		klog.Error(errors.Wrap(err, "host deletion failed, skipping"))
// 		return nil
// 	}
// 	klog.Info("delete was successful")
// } else {
// 	klog.Infof("not deleting host %s as other mappings exists", initiatorName)
// }

// return nil
// }

func (driver *Driver) chooseLUN() (int, error) {
	klog.Infof("listing all LUN mappings")
	volumes, status, err := driver.dothillClient.ShowHostMaps("")
	if err != nil && status == nil {
		return -1, err
	}
	if status.ReturnCode == hostMapDoesNotExistsErrorCode {
		klog.Info("initiator does not exist, assuming there is no LUN mappings yet and using LUN 1")
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

	if volumes[len(volumes)-1].LUN+1 < maximumLUN {
		return volumes[len(volumes)-1].LUN + 1, nil
	}

	return -1, errors.New("no more available LUNs")
}

func (driver *Driver) mapVolume(volumeName, initiatorName string, lun int) error {
	klog.Infof("trying to map volume %s for initiator %s on LUN %d", volumeName, initiatorName, lun)
	_, status, err := driver.dothillClient.MapVolume(volumeName, initiatorName, "rw", lun)
	if err != nil && status == nil {
		return err
	}
	if status.ReturnCode == hostDoesNotExistsErrorCode {
		nodeName := strings.Split(initiatorName, ":")[1]
		klog.Infof("initiator does not exist, creating it with nickname %s", nodeName)
		_, _, err = driver.dothillClient.CreateHost(nodeName, initiatorName)
		if err != nil {
			return err
		}
		klog.Info("retrying to map volume")
		_, _, err = driver.dothillClient.MapVolume(volumeName, initiatorName, "rw", lun)
		if err != nil {
			return err
		}
	}

	klog.Info("mapping was successful")
	return nil
}

func runPreflightChecks(parameters map[string]string, capabilities []*csi.VolumeCapability) error {
	checkIfKeyExistsInConfig := func(key string) error {
		klog.V(2).Infof("checking for %s in storage class parameters", key)
		_, ok := parameters[key]
		if !ok {
			return status.Errorf(codes.FailedPrecondition, "'%s' is missing from configuration", key)
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

	for _, capability := range capabilities {
		if capability.GetAccessMode().GetMode() != csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER {
			return status.Error(codes.FailedPrecondition, "dothill storage only supports ReadWriteOnce access mode")
		}
	}

	return nil
}
