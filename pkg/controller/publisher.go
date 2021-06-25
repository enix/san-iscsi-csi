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
	"sort"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/dothill-api-go/v2"
	"github.com/enix/san-iscsi-csi/pkg/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	klog "k8s.io/klog/v2"
)

func getVolumeMapsHostNames(client *dothill.Client, name string) ([]string, *dothill.ResponseStatus, error) {
	if name != "" {
		name = fmt.Sprintf("\"%s\"", name)
	}
	res, status, err := client.Request(fmt.Sprintf("/show/volume-maps/%s", name))
	if err != nil {
		return []string{}, status, err
	}

	hostNames := []string{}
	for _, rootObj := range res.Objects {
		if rootObj.Name != "volume-view" {
			continue
		}

		for _, object := range rootObj.Objects {
			hostName := object.PropertiesMap["identifier"].Data
			if object.Name == "host-view" && hostName != "all other hosts" {
				hostNames = append(hostNames, hostName)
			}
		}
	}

	return hostNames, status, err
}

// ControllerPublishVolume attaches the given volume to the node
func (driver *Controller) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot publish volume with empty ID")
	}
	if len(req.GetNodeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot publish volume to a node with empty ID")
	}
	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, "cannot publish volume without capabilities")
	}

	initiatorName := req.GetNodeId()

	common.AddLogTag(ctx, "volumeId", req.GetVolumeId())
	common.AddLogTag(ctx, "nodeId", req.GetNodeId())
	common.AddLogTag(ctx, "initiatorName", initiatorName)

	common.LogInfoS(ctx, "attach request")
	hostNames, _, err := getVolumeMapsHostNames(driver.dothillClient, req.GetVolumeId())
	if err != nil {
		return nil, err
	}
	for _, hostName := range hostNames {
		if hostName != initiatorName {
			return nil, status.Errorf(codes.FailedPrecondition, "volume %s is already attached to another node", req.GetVolumeId())
		}
	}

	lun, err := driver.chooseLUN(ctx, initiatorName)
	if err != nil {
		return nil, err
	}
	common.LogInfoS(ctx, "LUN choosed", "lun", lun)

	if err = driver.mapVolume(ctx, req.GetVolumeId(), initiatorName, lun); err != nil {
		return nil, err
	}

	common.LogInfoS(ctx, "successfully mapped volume", "lun", lun)
	return &csi.ControllerPublishVolumeResponse{
		PublishContext: map[string]string{"lun": strconv.Itoa(lun)},
	}, nil
}

// ControllerUnpublishVolume deattaches the given volume from the node
func (driver *Controller) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "cannot unpublish volume with empty ID")
	}

	common.AddLogTag(ctx, "volumeId", req.GetVolumeId())
	common.AddLogTag(ctx, "nodeId", req.GetNodeId())

	common.LogInfoS(ctx, "unmapping volume")
	_, status, err := driver.dothillClient.UnmapVolume(req.GetVolumeId(), req.GetNodeId())
	if err != nil {
		if status != nil && status.ReturnCode == unmapFailedErrorCode {
			common.LogInfoS(ctx, "unmap failed, assuming volume is already unmapped")
			return &csi.ControllerUnpublishVolumeResponse{}, nil
		}

		return nil, err
	}

	common.LogInfoS(ctx, "successfully unmapped volume from all initiators")
	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

func (driver *Controller) chooseLUN(ctx context.Context, initiatorName string) (int, error) {
	common.LogInfoS(ctx, "listing all LUN mappings")
	volumes, responseStatus, err := driver.dothillClient.ShowHostMaps(initiatorName)
	if err != nil && responseStatus == nil {
		return -1, err
	}
	if responseStatus.ReturnCode == hostMapDoesNotExistsErrorCode {
		common.LogInfoS(ctx, "initiator does not exist, assuming there is no LUN mappings yet and using LUN 1")
		return 1, nil
	}
	if err != nil {
		return -1, err
	}

	sort.Sort(Volumes(volumes))

	klog.V(5).InfoS("checking if LUN 1 is not already in use", common.GetLogKeyAndValues(ctx)...)
	if len(volumes) == 0 || volumes[0].LUN > 1 {
		return 1, nil
	}

	klog.V(5).InfoS("searching for an available LUN between LUNs in use", common.GetLogKeyAndValues(ctx)...)
	for index := 1; index < len(volumes); index++ {
		if volumes[index].LUN-volumes[index-1].LUN > 1 {
			return volumes[index-1].LUN + 1, nil
		}
	}

	klog.V(5).InfoS("checking if next LUN is not above maximum LUNs limit", common.GetLogKeyAndValues(ctx)...)
	if volumes[len(volumes)-1].LUN+1 < common.MaximumLUN {
		return volumes[len(volumes)-1].LUN + 1, nil
	}

	return -1, status.Error(codes.ResourceExhausted, "no more available LUNs")
}

func (driver *Controller) mapVolume(ctx context.Context, volumeName, initiatorName string, lun int) error {
	common.LogInfoS(ctx, "trying to map volume", "volumeName", volumeName, "initiatorName", initiatorName, "lun", lun)
	_, metadata, err := driver.dothillClient.MapVolume(volumeName, initiatorName, "rw", lun)
	if err != nil && metadata == nil {
		return err
	}
	if metadata.ReturnCode == hostDoesNotExistsErrorCode {
		nodeIDParts := strings.Split(initiatorName, ":")
		if len(nodeIDParts) < 2 {
			return status.Error(codes.InvalidArgument, "specified node ID is not a valid IQN")
		}

		nodeName := strings.Join(nodeIDParts[1:], ":")
		common.LogInfoS(ctx, "initiator does not exist, creating it", "nodeName", nodeName, "initiatorName", initiatorName)
		_, _, err = driver.dothillClient.CreateHost(nodeName, initiatorName)
		if err != nil {
			return err
		}
		common.LogInfoS(ctx, "retrying to map volume")
		_, _, err = driver.dothillClient.MapVolume(volumeName, initiatorName, "rw", lun)
		if err != nil {
			return err
		}
	} else if metadata.ReturnCode == volumeNotFoundErrorCode {
		return status.Errorf(codes.NotFound, "volume %s not found", volumeName)
	} else if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
