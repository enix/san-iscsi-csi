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

package main

import (
	"testing"

	"github.com/enix/dothill-csi/pkg/controller"
	"github.com/enix/dothill-csi/pkg/node"
	"github.com/kubernetes-csi/csi-test/pkg/sanity"
)

// Test starts the drivers in background and runs k8s sanity checks
func Test(t *testing.T) {
	controllerSocketPath := "unix:///tmp/controller.sock"
	nodeSocketPath := "unix:///tmp/node.sock"

	ctrl := controller.New()
	node := node.New("/var/lib/kubelet")

	go ctrl.Start(controllerSocketPath)
	defer ctrl.Stop()

	go node.Start(nodeSocketPath)
	defer node.Stop()

	sanity.Test(t, &sanity.Config{
		Address:                  nodeSocketPath,
		ControllerAddress:        controllerSocketPath,
		SecretsFile:              "./secrets.yml",
		TestVolumeParametersFile: "./config.yml",
	})
}
