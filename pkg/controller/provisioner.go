package controller

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/dothill-storage-controller/pkg/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

// CreateVolume creates a new volume from the given request. The function is
// idempotent.
func (driver *Driver) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot create volume with empty name")
	}

	parameters := req.GetParameters()
	caps := req.GetVolumeCapabilities()
	err := driver.beginRoutine(&common.DriverCtx{
		Req:         req,
		Credentials: req.GetSecrets(),
		Parameters:  &parameters,
		VolumeCaps:  &caps,
	})
	defer driver.endRoutine()
	if err != nil {
		return nil, err
	}

	size := req.GetCapacityRange().GetRequiredBytes()
	if size == 0 {
		size = 4096
	}

	sizeStr := fmt.Sprintf("%diB", size)
	klog.Infof("received %s volume request\n", sizeStr)

	volumeID := req.GetName()
	if len(volumeID) > common.VolumeNameMaxLength {
		volumeID = volumeID[:common.VolumeNameMaxLength]
	}

	klog.Infof("creating volume %s (size %s) in pool %s", volumeID, sizeStr, parameters[common.PoolConfigKey])
	_, _, err = driver.dothillClient.CreateVolume(volumeID, sizeStr, parameters[common.PoolConfigKey])
	if err != nil {
		return nil, err
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
func (driver *Driver) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot delete volume with empty ID")
	}

	err := driver.beginRoutine(&common.DriverCtx{
		Req:         req,
		Credentials: req.GetSecrets(),
	})
	defer driver.endRoutine()
	if err != nil {
		return nil, err
	}

	klog.Infof("deleting volume %s", req.GetVolumeId())
	_, status, err := driver.dothillClient.DeleteVolume(req.GetVolumeId())
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
