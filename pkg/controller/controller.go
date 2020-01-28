package controller

import (
	"context"
	"fmt"
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/dothill-api-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

// Driver is the implementation of csi.ControllerServer
type Driver struct {
	dothillClient *dothill.Client
	lock          sync.Mutex
}

// NewDriver is a convenience fn for creating a controller driver
func NewDriver() *Driver {
	return &Driver{
		dothillClient: dothill.NewClient(),
	}
}

// ControllerGetCapabilities returns the capabilities of the controller service.
func (driver *Driver) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	var csc []*csi.ControllerServiceCapability
	cl := []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
		// csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
		// csi.ControllerServiceCapability_RPC_LIST_SNAPSHOTS,
		// csi.ControllerServiceCapability_RPC_CLONE_VOLUME,
		// csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
	}

	for _, cap := range cl {
		klog.Infof("enabling controller service capability: %v", cap.String())
		csc = append(csc, &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: cap,
				},
			},
		})
	}

	return &csi.ControllerGetCapabilitiesResponse{Capabilities: csc}, nil
}

// DeleteVolume deletes the given volume. The function is idempotent.
func (driver *Driver) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "DeleteVolume unimplemented yet")
}

// ControllerPublishVolume attaches the given volume to the node
func (driver *Driver) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	fmt.Println("go it")
	return nil, status.Error(codes.Unimplemented, "ControllerPublishVolume unimplemented yet")
}

// ControllerUnpublishVolume deattaches the given volume from the node
func (driver *Driver) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ControllerUnpublishVolume unimplemented yet")
}

// ValidateVolumeCapabilities checks whether the volume capabilities requested
// are supported.
func (driver *Driver) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ValidateVolumeCapabilities unimplemented yet")
}

// ListVolumes returns a list of all requested volumes
func (driver *Driver) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListVolumes unimplemented yet")
}

// GetCapacity returns the capacity of the storage pool
func (driver *Driver) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetCapacity unimplemented yet")
}

// ControllerExpandVolume expands a volume to the given new size
func (driver *Driver) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ControllerExpandVolume unimplemented yet")
}

// CreateSnapshot creates a snapshot of the given volume
func (driver *Driver) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CreateSnapshot unimplemented yet")
}

// DeleteSnapshot deletes a snapshot of the given volume
func (driver *Driver) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "DeleteSnapshot unimplemented yet")
}

// ListSnapshots list existing snapshots
func (driver *Driver) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListSnapshots unimplemented yet")
}

func (driver *Driver) configureClient(credentials map[string]string, apiAddr string) error {
	username := string(credentials[usernameSecretKey])
	password := string(credentials[passwordSecretKey])
	klog.Infof("using dothill API at address %s", apiAddr)
	if driver.dothillClient.Addr == apiAddr && driver.dothillClient.Username == username {
		klog.Info("dothill client is already configured for this API, skipping login")
		return nil
	}

	driver.dothillClient.Username = username
	driver.dothillClient.Password = password
	driver.dothillClient.Addr = apiAddr

	klog.Infof("login into %s as user %s", driver.dothillClient.Addr, driver.dothillClient.Username)
	err := driver.dothillClient.Login()
	if err != nil {
		return err
	}

	klog.Info("login was successful")
	return nil
}
