package controller

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/dothill-csi/pkg/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

func (controller *Controller) checkVolumeExists(volumeID string, size int64) (bool, error) {
	data, responseStatus, err := controller.dothillClient.ShowVolumes(volumeID)
	if err != nil && responseStatus.ReturnCode != -10058 {
		return false, err
	}

	for _, object := range data.Objects {
		if object.Name == "volume" && object.PropertiesMap["volume-name"].Data == volumeID {
			blocks, _ := strconv.ParseInt(object.PropertiesMap["blocks"].Data, 10, 64)
			blocksize, _ := strconv.ParseInt(object.PropertiesMap["blocksize"].Data, 10, 64)

			if blocks*blocksize == size {
				return true, nil
			}
			return true, status.Error(codes.AlreadyExists, "cannot create volume with same name but different capacity than the existing one")
		}
	}

	return false, nil
}

// CreateVolume creates a new volume from the given request. The function is idempotent.
func (controller *Controller) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "cannot create volume with empty name")
	}

	size := req.GetCapacityRange().GetRequiredBytes()
	sizeStr := getSizeStr(size)
	parameters := req.GetParameters()
	klog.Infof("received %s volume request\n", sizeStr)

	volumeID := req.GetName()
	if len(volumeID) > common.VolumeNameMaxLength {
		volumeID = volumeID[4:]
		volumeID = strings.ReplaceAll(volumeID, "-", "")
		volumeID = volumeID[:common.VolumeNameMaxLength]
	}

	klog.Infof("creating volume %s (size %s) in pool %s", volumeID, sizeStr, parameters[common.PoolConfigKey])

	volumeExists, err := controller.checkVolumeExists(volumeID, size)
	if err != nil {
		return nil, err
	}

	if !volumeExists {
		_, _, err = controller.dothillClient.CreateVolume(volumeID, sizeStr, parameters[common.PoolConfigKey])
		if err != nil {
			return nil, err
		}
	}

	volume := &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      volumeID,
			VolumeContext: parameters,
			CapacityBytes: req.GetCapacityRange().GetRequiredBytes(),
			ContentSource: req.GetVolumeContentSource(),
		},
	}

	klog.Infof("created volume %s (%s)", volumeID, sizeStr)
	klog.V(8).Infof("created volume %+v", volume)
	return volume, nil
}

// DeleteVolume deletes the given volume. The function is idempotent.
func (controller *Controller) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot delete volume with empty ID")
	}

	klog.Infof("deleting volume %s", req.GetVolumeId())
	_, status, err := controller.dothillClient.DeleteVolume(req.GetVolumeId())
	if err != nil {
		if status != nil && status.ReturnCode == volumeNotFoundErrorCode {
			klog.Infof("volume %s does not exist, assuming it has already been deleted", req.GetVolumeId())
			return &csi.DeleteVolumeResponse{}, nil
		}
		return nil, err
	}

	klog.Infof("successfully deleted volume %s", req.GetVolumeId())
	return &csi.DeleteVolumeResponse{}, nil
}

func getSizeStr(size int64) string {
	if size == 0 {
		size = 4096
	}

	return fmt.Sprintf("%dB", size)
}
