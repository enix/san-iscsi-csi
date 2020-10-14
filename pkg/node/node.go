package node

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"regexp"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/dothill-storage-controller/pkg/common"
	"github.com/kubernetes-csi/csi-lib-iscsi/iscsi"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

// Driver is the implementation of csi.NodeServer
type Driver struct {
	mutex sync.Mutex
	kubeletPath string
}

// NewDriver is a convenience function for creating a node driver
func NewDriver(kubeletPath string) *Driver {
	if klog.V(8) {
		iscsi.EnableDebugLogging(os.Stderr)
	}

	return &Driver{kubeletPath: kubeletPath}
}

// NodeGetInfo returns info about the node
func (driver *Driver) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	initiatorName, err := readInitiatorName()
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return &csi.NodeGetInfoResponse{
		NodeId:            initiatorName,
		MaxVolumesPerNode: 255,
	}, nil
}

// NodeGetCapabilities returns the supported capabilities of the node server
func (driver *Driver) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	var csc []*csi.NodeServiceCapability
	cl := []csi.NodeServiceCapability_RPC_Type{
		// csi.NodeServiceCapability_RPC_EXPAND_VOLUME,
	}

	for _, cap := range cl {
		klog.V(4).Infof("enabled node service capability: %v", cap.String())
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
func (driver *Driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot publish volume with empty id")
	}
	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot publish volume at an empty path")
	}
	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, "cannot publish volume without capabilities")
	}

	driver.beginRoutine(&common.DriverCtx{Req: req})
	defer driver.endRoutine()
	klog.Infof("publishing volume %s", req.GetVolumeId())

	portals := strings.Split(req.GetVolumeContext()[common.PortalsConfigKey], ",")
	klog.Infof("ISCSI portals: %s", portals)

	lun, _ := strconv.ParseInt(req.GetPublishContext()["lun"], 10, 32)
	klog.Infof("LUN: %d", lun)

	klog.Info("initiating ISCSI connection...")
	targets := make([]iscsi.TargetInfo, 0)
	for _, portal := range portals {
		targets = append(targets, iscsi.TargetInfo{
			Iqn:    req.GetVolumeContext()[common.TargetIQNConfigKey],
			Portal: portal,
		})
	}
	connector := &iscsi.Connector{
		Targets: targets,
		Lun:     int32(lun),
	}
	path, err := iscsi.Connect(*connector)
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	klog.Infof("attached device at %s", path)

	connector.DevicePath = path[4:]
	if connector.DevicePath[1:4] == "dm-" {
		klog.Info("device is using multipath")
		connector.Multipath = true
	} else {
		klog.Warning("device is NOT using multipath")
	}

	fsType := req.GetVolumeContext()[common.FsTypeConfigKey]
	err = ensureFsType(fsType, path)
	if err != nil {
		return nil, err
	}

	klog.Infof("mounting volume at %s", req.GetTargetPath())
	os.Mkdir(req.GetTargetPath(), 00755)
	out, err := exec.Command("mount", "-t", fsType, path, req.GetTargetPath()).CombinedOutput()
	if err != nil {
		klog.Error(string(out))
		return nil, errors.New(string(out))
	}

	iscsiInfoPath := fmt.Sprintf("%s/plugins/%s/iscsi-%s.json", driver.kubeletPath, common.PluginName, req.GetVolumeId())
	klog.Infof("saving ISCSI connection info in %s", iscsiInfoPath)
	err = iscsi.PersistConnector(connector, iscsiInfoPath)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	klog.Infof("succesfully mounted volume at %s", req.GetTargetPath())
	return &csi.NodePublishVolumeResponse{}, nil
}

// NodeUnpublishVolume unmounts the volume from the target path
func (driver *Driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot unpublish volume with empty id")
	}
	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot publish volume at an empty path")
	}

	driver.beginRoutine(&common.DriverCtx{Req: req})
	defer driver.endRoutine()
	klog.Infof("unpublishing volume %s", req.GetVolumeId())

	_, err := os.Stat(req.GetTargetPath())
	if err == nil {
		klog.Infof("unmounting volume at %s", req.GetTargetPath())
		out, err := exec.Command("umount", req.GetTargetPath()).CombinedOutput()
		if err != nil && !os.IsNotExist(err) {
			klog.Error(errors.New(string(out)))
			return nil, errors.New(string(out))
		}
		os.Remove(req.GetTargetPath())
	}

	iscsiInfoPath := fmt.Sprintf("%s/plugins/%s/iscsi-%s.json", driver.kubeletPath, common.PluginName, req.GetVolumeId())
	klog.Infof("loading ISCSI connection info from %s", iscsiInfoPath)
	connector, err := iscsi.GetConnectorFromFile(iscsiInfoPath)
	if err != nil {
		klog.Error(errors.Wrap(err, "assuming ISCSI connection is already closed"))
		return &csi.NodeUnpublishVolumeResponse{}, nil
	}

	klog.Info("detaching ISCSI device")
	err = iscsi.DisconnectVolume(*connector)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	klog.Info("rescaning ISCSI sessions")
	out, err := exec.Command("iscsiadm", "-m", "session", "--rescan").CombinedOutput()
	if err != nil {
		klog.Error(errors.New(string(out)))
		return nil, errors.New(string(out))
	}

	klog.Infof("deleting ISCSI connection info file %s", iscsiInfoPath)
	os.Remove(iscsiInfoPath)

	klog.Info("successfully detached ISCSI device")
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeExpandVolume finalizes volume expansion on the node
func (driver *Driver) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	fmt.Println("NodeExpandVolume call")
	return nil, status.Error(codes.Unimplemented, "NodeExpandVolume unimplemented yet")
}

// NodeGetVolumeStats return info about a given volume
// Will not be called as the plugin does not have the GET_VOLUME_STATS capability
func (driver *Driver) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "NodeGetVolumeStats is unimplemented and should not be called")
}

// NodeStageVolume mounts the volume to a staging path on the node. This is
// called by the CO before NodePublishVolume and is used to temporary mount the
// volume to a staging path. Once mounted, NodePublishVolume will make sure to
// mount it to the appropriate path
// Will not be called as the plugin does not have the STAGE_UNSTAGE_VOLUME capability
func (driver *Driver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "NodeStageVolume is unimplemented and should not be called")
}

// NodeUnstageVolume unstages the volume from the staging path
// Will not be called as the plugin does not have the STAGE_UNSTAGE_VOLUME capability
func (driver *Driver) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "NodeUnstageVolume is unimplemented and should not be called")
}

func (driver *Driver) beginRoutine(ctx *common.DriverCtx) {
	driver.mutex.Lock()
	ctx.BeginRoutine()
}

func (driver *Driver) endRoutine() {
	klog.Infof("=== [ROUTINE END] ===\n\n")
	driver.mutex.Unlock()
}

// see https://github.com/kubernetes-csi/driver-registrar/blob/795af1899f3c94dd0c6dda2a25ed301123479bb9/vendor/k8s.io/kubernetes/pkg/util/mount/mount_linux.go#L543
func getDiskFormat(disk string) (string, error) {
	args := []string{"-p", "-s", "TYPE", "-s", "PTTYPE", "-o", "export", disk}
	klog.V(4).Infof("Attempting to determine if disk %q is formatted using blkid with args: (%v)", disk, args)
	output, err := exec.Command("blkid", args...).CombinedOutput()
	klog.V(4).Infof("Output: %q, err: %v", output, err)

	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			if exit.ExitCode() == 2 {
				// Disk device is unformatted.
				// For `blkid`, if the specified token (TYPE/PTTYPE, etc) was
				// not found, or no (specified) devices could be identified, an
				// exit code of 2 is returned.
				return "", nil
			}
		}
		klog.Errorf("Could not determine if disk %q is formatted (%v)", disk, err)
		return "", err
	}

	var fsType, ptType string

	re := regexp.MustCompile(`([A-Z]+)="?([^"\n]+)"?`) // Handles alpine and debian outputs
	matches := re.FindAllSubmatch(output, -1)
	for _, match := range matches {
		if len(match) != 3 {
			return "", fmt.Errorf("blkid returns invalid output: %s", output)
		}
		// TYPE is filesystem type, and PTTYPE is partition table type, according
		// to https://www.kernel.org/pub/linux/utils/util-linux/v2.21/libblkid-docs/.
		if string(match[1]) == "TYPE" {
			fsType = string(match[2])
		} else if string(match[1]) == "PTTYPE" {
			ptType = string(match[2])
		}
	}

	if len(ptType) > 0 {
		klog.V(4).Infof("Disk %s detected partition table type: %s", ptType)
		// Returns a special non-empty string as filesystem type, then kubelet
		// will not format it.
		return "unknown data, probably partitions", nil
	}

	return fsType, nil
}


func ensureFsType(fsType string, disk string) (error) {
	currentFsType, err := getDiskFormat(disk)

	if err != nil {
		return err
	}

	klog.V(1).Infof("detected filesystem: %q", currentFsType)
	if currentFsType != "ext4" {
		klog.Infof("creating %s filesystem on device %s", fsType, disk)
		out, err := exec.Command(fmt.Sprintf("mkfs.%s", fsType), disk).CombinedOutput()
		if err != nil {
			klog.Error(string(out))
			return errors.New(string(out))
		}
	}

	return nil
}

func readInitiatorName() (string, error) {
	initiatorNameFilePath := "/etc/iscsi/initiatorname.iscsi"
	file, err := os.Open(initiatorNameFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if strings.TrimSpace(line[:equal]) == "InitiatorName" {
				return strings.TrimSpace(line[equal+1:]), nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("InitiatorName key is missing from %s", initiatorNameFilePath)
}