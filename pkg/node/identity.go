package node

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/dothill-storage-controller/pkg/common"
	"os/exec"
	"k8s.io/klog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPluginInfo returns metadata of the plugin
func (d *Driver) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	return &csi.GetPluginInfoResponse{
		Name:          common.PluginName,
		VendorVersion: common.Version,
	}, nil
}

// GetPluginCapabilities returns available capabilities of the plugin
func (d *Driver) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{},
	}, nil
}

// Probe returns the health and readiness of the plugin
func (d *Driver) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	if !isKernelModLoaded("iscsi_tcp") {
		return nil, status.Error(codes.FailedPrecondition, "kernel mod iscsi_tcp is not loaded")
	}
	if !isKernelModLoaded("dm_multipath") {
		return nil, status.Error(codes.FailedPrecondition, "kernel mod dm_multipath is not loaded")
	}

	return &csi.ProbeResponse{}, nil
}

func isKernelModLoaded(modName string) bool {
	klog.V(5).Infof("verifiying that %q kernel mod is loaded", modName)
	err := exec.Command("grep", "^" + modName, "/proc/modules", "-q").Run()
	
	if err != nil {
		klog.Errorf("required kernel mod %q is not loaded", modName)
		return false
	}

	klog.V(5).Infof("kernel mod %q is loaded", modName)

	return true
}
