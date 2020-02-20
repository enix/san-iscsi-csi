package controller

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/dothill-storage-controller/pkg/common"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
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

	err = driver.configureClient(req.GetSecrets(), parameters[common.APIAddressConfigKey])
	if err != nil {
		return nil, err
	}

	volumeID := uuid.NewUUID().String()[:common.VolumeNameMaxLength]
	klog.Infof("creating volume %s (size %s) in pool %s", volumeID, sizeStr, parameters[common.PoolConfigKey])
	_, _, err = driver.dothillClient.CreateVolume(volumeID, sizeStr, parameters[common.PoolConfigKey])
	if err != nil {
		return nil, err
	}

	volume := &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      volumeID,
			VolumeContext: req.GetParameters(),
			CapacityBytes: req.GetCapacityRange().GetRequiredBytes(),
			ContentSource: req.GetVolumeContentSource(),
		},
	}

	klog.Infof("created volume %s (%s)", volumeID, sizeStr)
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

func runPreflightChecks(parameters map[string]string, capabilities []*csi.VolumeCapability) error {
	checkIfKeyExistsInConfig := func(key string) error {
		klog.V(2).Infof("checking for %s in storage class parameters", key)
		_, ok := parameters[key]
		if !ok {
			return status.Errorf(codes.FailedPrecondition, "'%s' is missing from configuration", key)
		}
		return nil
	}

	if err := checkIfKeyExistsInConfig(common.FsTypeConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(common.PoolConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(common.TargetIQNConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(common.PortalsConfigKey); err != nil {
		return err
	}
	if err := checkIfKeyExistsInConfig(common.InitiatorNameConfigKey); err != nil {
		if err2 := checkIfKeyExistsInConfig(common.UniqueInitiatorNameByPvcConfigKey); err2 != nil {
			return errors.Wrap(err, err2.Error())
		}
	}
	if err := checkIfKeyExistsInConfig(common.APIAddressConfigKey); err != nil {
		return err
	}

	for _, capability := range capabilities {
		if capability.GetAccessMode().GetMode() != csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER {
			return status.Error(codes.FailedPrecondition, "dothill storage only supports ReadWriteOnce access mode")
		}
	}

	return nil
}
