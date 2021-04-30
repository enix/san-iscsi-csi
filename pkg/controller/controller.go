/*
 * Copyright (c) 2021 Enix, SAS
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 *
 * Authors:
 * Paul Laffitte <paul.laffitte@enix.fr>
 * Arthur Chaloin <arthur.chaloin@enix.fr>
 * Alexandre Buisine <alexandre.buisine@enix.fr>
 */

package controller

import (
	"context"
	"errors"
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/dothill-api-go"
	"github.com/enix/dothill-csi/pkg/common"
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
	volumeHasSnapshot             = -10183
	snapshotNotFoundErrorCode     = -10050
	snapshotAlreadyExists         = -10186
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
	"/csi.v1.Controller/ControllerExpandVolume":    {},
}

var nonAuthenticatedMethods = []string{
	"/csi.v1.Controller/ControllerGetCapabilities",
	"/csi.v1.Controller/ListVolumes",
	"/csi.v1.Controller/GetCapacity",
	"/csi.v1.Controller/ControllerGetVolume",
	"/csi.v1.Identity/Probe",
	"/csi.v1.Identity/GetPluginInfo",
	"/csi.v1.Identity/GetPluginCapabilities",
}

// Controller is the implementation of csi.ControllerServer
type Controller struct {
	*common.Driver

	dothillClient *dothill.Client
}

// DriverCtx contains data common to most calls
type DriverCtx struct {
	Credentials map[string]string
	Parameters  *map[string]string
	VolumeCaps  *[]*csi.VolumeCapability
}

// New is a convenience fn for creating a controller driver
func New() *Controller {
	dothillClient := dothill.NewClient()
	controller := &Controller{
		Driver:        common.NewDriver(dothillClient.Collector),
		dothillClient: dothillClient,
	}

	controller.InitServer(
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			if mutex, exists := csiMutexes[info.FullMethod]; exists {
				mutex.Lock()
				defer mutex.Unlock()
			}
			return handler(ctx, req)
		},
		common.NewLogRoutineServerInterceptor(func(string) bool {
			return true
		}),
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

			err := controller.beginRoutine(&driverContext, info.FullMethod)
			defer controller.endRoutine()
			if err != nil {
				return nil, err
			}

			return handler(ctx, req)
		},
	)

	csi.RegisterIdentityServer(controller.Server, controller)
	csi.RegisterControllerServer(controller.Server, controller)

	return controller
}

// ControllerGetCapabilities returns the capabilities of the controller service.
func (controller *Controller) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	var csc []*csi.ControllerServiceCapability
	cl := []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
		csi.ControllerServiceCapability_RPC_LIST_SNAPSHOTS,
		csi.ControllerServiceCapability_RPC_CLONE_VOLUME,
		csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
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
func (controller *Controller) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	volumeID := req.GetVolumeId()
	if len(volumeID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot validate volume with empty ID")
	}
	if len(req.GetVolumeCapabilities()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot validate volume without capabilities")
	}
	_, _, err := controller.dothillClient.ShowVolumes(volumeID)
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
func (controller *Controller) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListVolumes is unimplemented and should not be called")
}

// GetCapacity returns the capacity of the storage pool
func (controller *Controller) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetCapacity is unimplemented and should not be called")
}

// ControllerGetVolume fetch current information about a volume
func (controller *Controller) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ControllerGetVolume is unimplemented and should not be called")
}

// Probe returns the health and readiness of the plugin
func (controller *Controller) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return &csi.ProbeResponse{}, nil
}

func (controller *Controller) beginRoutine(ctx *DriverCtx, methodName string) error {
	if err := runPreflightChecks(ctx.Parameters, ctx.VolumeCaps); err != nil {
		return err
	}

	needsAuthentication := true
	for _, name := range nonAuthenticatedMethods {
		if methodName == name {
			needsAuthentication = false
			break
		}
	}

	if !needsAuthentication {
		return nil
	}

	if ctx.Credentials == nil {
		return errors.New("missing API credentials")
	}

	return controller.configureClient(ctx.Credentials)
}

func (controller *Controller) endRoutine() {
	controller.dothillClient.HTTPClient.CloseIdleConnections()
}

func (controller *Controller) configureClient(credentials map[string]string) error {
	username := string(credentials[common.UsernameSecretKey])
	password := string(credentials[common.PasswordSecretKey])
	apiAddr := string(credentials[common.APIAddressConfigKey])

	if len(apiAddr) == 0 || len(username) == 0 || len(password) == 0 {
		return status.Error(codes.InvalidArgument, "at least one field is missing in credentials secret")
	}

	klog.Infof("using dothill API at address %s", apiAddr)
	if controller.dothillClient.Addr == apiAddr && controller.dothillClient.Username == username {
		klog.Info("dothill client is already configured for this API, skipping login")
		return nil
	}

	controller.dothillClient.Username = username
	controller.dothillClient.Password = password
	controller.dothillClient.Addr = apiAddr
	klog.Infof("login into %q as user %q", controller.dothillClient.Addr, controller.dothillClient.Username)
	err := controller.dothillClient.Login()
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
