package controller

import (
	"context"
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/dothill-api-go"
	"github.com/enix/dothill-storage-controller/pkg/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

const (
	hostDoesNotExistsErrorCode    = -10386
	hostMapDoesNotExistsErrorCode = -10074
	unmapFailedErrorCode          = -10509
	volumeNotFoundErrorCode       = -10075
)

var volumeCapabilities = []*csi.VolumeCapability{
	{
		AccessType: &csi.VolumeCapability_Mount{
			Mount: &csi.VolumeCapability_MountVolume{},
		},
		AccessMode: &csi.VolumeCapability_AccessMode{
			Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
		},
	},
}

var csiMutexes = map[string]*sync.Mutex{
	"/csi.v1.Controller/CreateVolume":              {},
	"/csi.v1.Controller/ControllerPublishVolume":   {},
	"/csi.v1.Controller/DeleteVolume":              {},
	"/csi.v1.Controller/ControllerUnpublishVolume": {},
}

// Driver is the implementation of csi.ControllerServer
type Driver struct {
	dothillClient *dothill.Client
}

// DriverCtx contains data common to most calls
type DriverCtx struct {
	Credentials map[string]string
	Parameters  *map[string]string
	VolumeCaps  *[]*csi.VolumeCapability
}

// NewDriver is a convenience fn for creating a controller driver
func NewDriver() *Driver {
	return &Driver{dothillClient: dothill.NewClient()}
}

// NewServerInterceptors implements DriverImpl.NewServerInterceptors
func (driver *Driver) NewServerInterceptors(logRoutineServerInterceptor grpc.UnaryServerInterceptor) *[]grpc.UnaryServerInterceptor {
	serverInterceptors := []grpc.UnaryServerInterceptor{
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			if mutex, exists := csiMutexes[info.FullMethod]; exists {
				mutex.Lock()
				defer mutex.Unlock()
			}
			return handler(ctx, req)
		},
		logRoutineServerInterceptor,
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			driverContext := DriverCtx{}
			if reqWithSecrets, ok := req.(common.WithSecrets); ok {
				driverContext.Credentials = reqWithSecrets.GetSecrets()
			}
			if reqWithParameters, ok := req.(common.WithParameters); ok {
				driverContext.Parameters = reqWithParameters.GetParameters()
			}
			if reqWithVolumeCaps, ok := req.(common.WithVolumeCaps); ok {
				driverContext.VolumeCaps = reqWithVolumeCaps.GetVolumeCapabilities()
			}

			err := driver.beginRoutine(&driverContext)
			defer driver.endRoutine()
			if err != nil {
				return nil, err
			}

			return handler(ctx, req)
		},
	}

	return &serverInterceptors
}

// ShouldLogRoutine implements DriverImpl.ShouldLogRoutine
func (driver *Driver) ShouldLogRoutine(fullMethod string) bool {
	return true
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
		klog.V(4).Infof("enabled controller service capability: %v", cap.String())
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

// ValidateVolumeCapabilities checks whether the volume capabilities requested
// are supported.
func (driver *Driver) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	volumeID := req.GetVolumeId()
	if len(volumeID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot validate volume with empty ID")
	}
	if len(req.GetVolumeCapabilities()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot validate volume without capabilities")
	}
	_, _, err := driver.dothillClient.ShowVolumes(volumeID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "cannot validate volume not found")
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeCapabilities: volumeCapabilities,
		},
	}, nil
}

// ListVolumes returns a list of all requested volumes
func (driver *Driver) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListVolumes is unimplemented and should not be called")
}

// GetCapacity returns the capacity of the storage pool
func (driver *Driver) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetCapacity is unimplemented and should not be called")
}

// ControllerGetVolume fetch current information about a volume
func (driver *Driver) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ControllerGetVolume is unimplemented and should not be called")
}

func (driver *Driver) beginRoutine(ctx *DriverCtx) error {
	if err := runPreflightChecks(ctx.Parameters, ctx.VolumeCaps); err != nil {
		return err
	}

	if ctx.Credentials == nil {
		klog.Info("skipping login as this RPC does not require any API call")
		return nil
	}

	return driver.configureClient(ctx.Credentials)
}

func (driver *Driver) endRoutine() {
	driver.dothillClient.HTTPClient.CloseIdleConnections()
}

func (driver *Driver) configureClient(credentials map[string]string) error {
	username := string(credentials[common.UsernameSecretKey])
	password := string(credentials[common.PasswordSecretKey])
	apiAddr := string(credentials[common.APIAddressConfigKey])

	if len(apiAddr) == 0 || len(username) == 0 || len(password) == 0 {
		return status.Error(codes.InvalidArgument, "at least one field is missing in credentials secret")
	}

	klog.Infof("using dothill API at address %s", apiAddr)
	if driver.dothillClient.Addr == apiAddr && driver.dothillClient.Username == username {
		klog.Info("dothill client is already configured for this API, skipping login")
		return nil
	}

	driver.dothillClient.Username = username
	driver.dothillClient.Password = password
	driver.dothillClient.Addr = apiAddr
	klog.Infof("login into %q as user %q", driver.dothillClient.Addr, driver.dothillClient.Username)
	err := driver.dothillClient.Login()
	if err != nil {
		return status.Error(codes.Unauthenticated, err.Error())
	}

	klog.Info("login was successful")
	return nil
}

func runPreflightChecks(parameters *map[string]string, capabilities *[]*csi.VolumeCapability) error {
	checkIfKeyExistsInConfig := func(key string) error {
		if parameters == nil {
			return nil
		}

		klog.V(2).Infof("checking for %s in storage class parameters", key)
		_, ok := (*parameters)[key]
		if !ok {
			return status.Errorf(codes.InvalidArgument, "'%s' is missing from configuration", key)
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

	if capabilities != nil {
		if len(*capabilities) == 0 {
			return status.Error(codes.InvalidArgument, "missing volume capabilities")
		}
		for _, capability := range *capabilities {
			if capability.GetAccessMode().GetMode() != csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER {
				return status.Error(codes.FailedPrecondition, "dothill storage only supports ReadWriteOnce access mode")
			}
		}
	}

	return nil
}
