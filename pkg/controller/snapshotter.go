package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/dothill-api-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

// CreateSnapshot creates a snapshot of the given volume
func (controller *Controller) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	name := strings.Replace(req.Name[9:], "-", "", -1)

	_, respStatus, err := controller.dothillClient.CreateSnapshot(req.SourceVolumeId, name)
	if err != nil && respStatus.ReturnCode != snapshotAlreadyExists {
		return nil, err
	}

	response, _, err := controller.dothillClient.ShowSnapshots(name)
	if err != nil {
		return nil, err
	}

	var snapshot *csi.Snapshot
	for _, object := range response.Objects {
		if object.Typ != "snapshots" {
			continue
		}

		snapshot, err = newSnapshotFromResponse(&object)
		if err != nil {
			return nil, err
		}
	}

	if snapshot == nil {
		return nil, errors.New("snapshot not found")
	}

	if snapshot.SourceVolumeId != req.SourceVolumeId {
		return nil, status.Error(codes.AlreadyExists, "cannot validate volume with empty ID")
	}

	return &csi.CreateSnapshotResponse{Snapshot: snapshot}, nil
}

// DeleteSnapshot deletes a snapshot of the given volume
func (controller *Controller) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	_, status, err := controller.dothillClient.DeleteSnapshot(req.SnapshotId)
	if err != nil {
		if status != nil && status.ReturnCode == snapshotNotFoundErrorCode {
			klog.Infof("snapshot %s does not exist, assuming it has already been deleted", req.SnapshotId)
			return &csi.DeleteSnapshotResponse{}, nil
		}
		return nil, err
	}
	return &csi.DeleteSnapshotResponse{}, nil
}

// ListSnapshots list existing snapshots
func (controller *Controller) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	response, _, err := controller.dothillClient.ShowSnapshots()
	if err != nil {
		return nil, err
	}

	snapshots := []*csi.ListSnapshotsResponse_Entry{}
	for _, object := range response.Objects {
		if object.Typ != "snapshots" {
			continue
		}

		snapshot, err := newSnapshotFromResponse(&object)
		if err != nil {
			return nil, err
		}

		snapshots = append(snapshots, &csi.ListSnapshotsResponse_Entry{
			Snapshot: snapshot,
		})
	}

	return &csi.ListSnapshotsResponse{
		Entries: snapshots,
	}, nil
}

func newSnapshotFromResponse(object *dothill.Object) (*csi.Snapshot, error) {
	properties, err := object.GetProperties("total-size-numeric", "name", "master-volume-name", "creation-date-time-numeric")
	if err != nil {
		return nil, fmt.Errorf("could not read snapshot %v", err)
	}

	sizeBytes, err := strconv.ParseInt(properties[0].Data, 10, 64)
	snapshotId := properties[1].Data
	sourceVolumeId := properties[2].Data
	creationTime, err := creationTimeFromString(properties[3].Data)

	return &csi.Snapshot{
		SizeBytes:      sizeBytes,
		SnapshotId:     snapshotId,
		SourceVolumeId: sourceVolumeId,
		CreationTime:   creationTime,
		ReadyToUse:     true,
	}, nil
}

func creationTimeFromString(creationTime string) (*timestamp.Timestamp, error) {
	creationTimestamp, err := strconv.ParseInt(creationTime, 10, 64)
	if err != nil {
		return nil, err
	}

	return ptypes.TimestampProto(time.Unix(creationTimestamp, 0))
}
