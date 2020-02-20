package node

import (
	"context"
	"fmt"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-iscsi/iscsi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

// Driver is the implementation of csi.NodeServer
type Driver struct{}

// NewDriver is a convenience function for creating a node driver
func NewDriver() *Driver {
	return &Driver{}
}

// NodeGetInfo returns info about the node
func (d *Driver) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	hostname, _ := os.Hostname()
	return &csi.NodeGetInfoResponse{
		NodeId:            hostname,
		MaxVolumesPerNode: 255,
	}, nil
}

// NodeGetCapabilities returns the supported capabilities of the node server
func (d *Driver) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	var csc []*csi.NodeServiceCapability
	cl := []csi.NodeServiceCapability_RPC_Type{
		// csi.NodeServiceCapability_RPC_EXPAND_VOLUME,
	}

	for _, cap := range cl {
		klog.Infof("enabled node service capability: %v", cap.String())
		csc = append(csc, &csi.NodeServiceCapability{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: cap,
				},
			},
		})
	}

	return &csi.NodeGetCapabilitiesResponse{Capabilities: csc}, nil
}

// NodePublishVolume mounts the volume mounted to the staging path to the target path
func (d *Driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	path, err := iscsi.Connect(iscsi.Connector{
		Targets: []iscsi.TargetInfo{{
			Iqn:    "iqn.2015-11.com.hpe:storage.msa2050.18323cc9ed",
			Portal: "10.14.84.211",
			Port:   "3260",
		}},
		Lun:       3,
		Multipath: true,
	})

	if err != nil {
		return nil, err
	}

	fmt.Println(path)

	return &csi.NodePublishVolumeResponse{}, nil
}

// NodeUnpublishVolume unmounts the volume from the target path
func (d *Driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	fmt.Println("NodeUnpublishVolume call")
	return nil, status.Error(codes.Unimplemented, "NodeUnpublishVolume unimplemented yet")
}

// NodeExpandVolume finalizes volume expansion on the node
func (d *Driver) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	fmt.Println("NodeExpandVolume call")
	return nil, status.Error(codes.Unimplemented, "NodeExpandVolume unimplemented yet")
}

// NodeGetVolumeStats return info about a given volume
func (d *Driver) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	fmt.Println("NodeGetVolumeStats call")
	return nil, status.Error(codes.Unimplemented, "NodeGetVolumeStats unimplemented yet")
}

// NodeStageVolume mounts the volume to a staging path on the node. This is
// called by the CO before NodePublishVolume and is used to temporary mount the
// volume to a staging path. Once mounted, NodePublishVolume will make sure to
// mount it to the appropriate path
func (d *Driver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "NodeStageVolume is unimplemented and should not be called")
}

// NodeUnstageVolume unstages the volume from the staging path
func (d *Driver) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "NodeUnstageVolume is unimplemented and should not be called")
}
