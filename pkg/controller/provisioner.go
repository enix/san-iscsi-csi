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
	"fmt"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/san-iscsi-csi/pkg/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (controller *Controller) checkVolumeExists(volumeID string, size int64) (bool, error) {
	data, responseStatus, err := controller.dothillClient.ShowVolumes(volumeID)
	if err != nil && responseStatus.ReturnCode != -10058 {
		return false, err
	}

	for _, object := range data.Objects {
		if object.Name == "volume" && object.PropertiesMap["volume-name"].Data == volumeID {
			blocks, _ := strconv.ParseInt(object.PropertiesMap["blocks"].Data, 10, 64)
			blocksize, _ := strconv.ParseInt(object.PropertiesMap["blocksize"].Data, 10, 64)

			if blocks*blocksize == size {
				return true, nil
			}
			return true, status.Error(codes.AlreadyExists, "cannot create volume with same name but different capacity than the existing one")
		}
	}

	return false, nil
}

// CreateVolume creates a new volume from the given request. The function is idempotent.
func (controller *Controller) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "cannot create volume with empty name")
	}

	size := req.GetCapacityRange().GetRequiredBytes()
	sizeStr := getSizeStr(size)
	parameters := req.GetParameters()

	volumeID := req.GetName()
	if len(volumeID) > common.VolumeNameMaxLength {
		volumeID = volumeID[4:]
		volumeID = strings.ReplaceAll(volumeID, "-", "")
		volumeID = volumeID[:common.VolumeNameMaxLength]
	}

	common.AddLogTag(ctx, "volumeId", volumeID)
	common.AddLogTag(ctx, "size", sizeStr)

	common.LogInfoS(ctx, "creating volume", "pool", parameters[common.PoolConfigKey])
	volumeExists, err := controller.checkVolumeExists(volumeID, size)
	if err != nil {
		return nil, err
	}

	if !volumeExists {
		var sourceID string

		if volume := req.VolumeContentSource.GetVolume(); volume != nil {
			sourceID = volume.VolumeId
		}

		if snapshot := req.VolumeContentSource.GetSnapshot(); sourceID == "" && snapshot != nil {
			sourceID = snapshot.SnapshotId
		}

		if sourceID != "" {
			_, _, err = controller.dothillClient.CopyVolume(sourceID, volumeID, parameters[common.PoolConfigKey])
		} else {
			_, _, err = controller.dothillClient.CreateVolume(volumeID, sizeStr, parameters[common.PoolConfigKey])
		}
		if err != nil {
			return nil, err
		}
	}

	volume := &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      volumeID,
			VolumeContext: parameters,
			CapacityBytes: req.GetCapacityRange().GetRequiredBytes(),
			ContentSource: req.GetVolumeContentSource(),
		},
	}

	common.LogInfoS(ctx, "volume created")
	return volume, nil
}

// DeleteVolume deletes the given volume. The function is idempotent.
func (controller *Controller) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot delete volume with empty ID")
	}

	common.AddLogTag(ctx, "volumeId", req.GetVolumeId())

	common.LogInfoS(ctx, "deleting volume")
	_, respStatus, err := controller.dothillClient.DeleteVolume(req.GetVolumeId())
	if err != nil {
		if respStatus != nil {
			if respStatus.ReturnCode == volumeNotFoundErrorCode {
				common.LogInfoS(ctx, "volume does not exist, assuming it has already been deleted")
				return &csi.DeleteVolumeResponse{}, nil
			} else if respStatus.ReturnCode == volumeHasSnapshot {
				return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("volume %s cannot be deleted since it has snapshots", req.GetVolumeId()))
			}
		}
		return nil, err
	}

	common.LogInfoS(ctx, "successfully deleted volume")
	return &csi.DeleteVolumeResponse{}, nil
}

func getSizeStr(size int64) string {
	if size == 0 {
		size = 4096
	}

	return fmt.Sprintf("%dB", size)
}
