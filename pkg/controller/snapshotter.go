package controller

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateSnapshot creates a snapshot of the given volume
func (controller *Controller) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	fmt.Println("CreateSnapshot call")
	return nil, status.Error(codes.Unimplemented, "CreateSnapshot unimplemented yet")
}

// DeleteSnapshot deletes a snapshot of the given volume
func (controller *Controller) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	fmt.Println("DeleteSnapshot call")
	return nil, status.Error(codes.Unimplemented, "DeleteSnapshot unimplemented yet")
}

// ListSnapshots list existing snapshots
func (controller *Controller) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	fmt.Println("ListSnapshots call")
	return nil, status.Error(codes.Unimplemented, "ListSnapshots unimplemented yet")
}
