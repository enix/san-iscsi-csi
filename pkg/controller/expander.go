package controller

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ControllerExpandVolume expands a volume to the given new size
func (driver *Driver) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	fmt.Println("ControllerExpandVolume call")
	return nil, status.Error(codes.Unimplemented, "ControllerExpandVolume unimplemented yet")
}
