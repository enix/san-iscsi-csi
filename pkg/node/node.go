package node

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/dothill-storage-controller/pkg/common"
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
	klog.V(8).Infof("NodePublishVolume %+v", req)
	klog.Infof("publishing volume %s", req.GetVolumeId())

	portals := strings.Split(req.GetVolumeContext()[common.PortalsConfigKey], ",")
	klog.Infof("ISCSI portals: %s", portals)

	lun, _ := strconv.ParseInt(req.GetPublishContext()["lun"], 10, 32)
	klog.Infof("LUN: %d", lun)

	klog.Info("initiating ISCSI connection...")
	path, err := iscsi.Connect(iscsi.Connector{
		Targets: []iscsi.TargetInfo{{
			Iqn:    req.GetVolumeContext()[common.TargetIQNConfigKey],
			Portal: portals[0],
			Port:   "3260",
		}},
		Lun:       int32(lun),
		Multipath: true,
	})

	if err != nil {
		return nil, err
	}
	klog.Infof("attached device at %s", path)

	fsType := req.GetVolumeContext()[common.FsTypeConfigKey]
	klog.Infof("creating %s filesystem on device", fsType)
	out, err := exec.Command(fmt.Sprintf("mkfs.%s", fsType), path).CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out))
	}

	klog.Infof("mounting volume at %s", req.GetTargetPath())
	os.Mkdir(req.GetTargetPath(), 00755)
	out, err = exec.Command("mount", "-t", fsType, path, req.GetTargetPath()).CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out))
	}

	klog.Infof("succesfully mounted volume at %s", req.GetTargetPath())
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
