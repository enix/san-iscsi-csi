// Copyright (c) 2021 Enix, SAS
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing
// permissions and limitations under the License.
//
// Authors:
// Paul Laffitte <paul.laffitte@enix.fr>
// Arthur Chaloin <arthur.chaloin@enix.fr>
// Alexandre Buisine <alexandre.buisine@enix.fr>

module github.com/enix/dothill-csi

go 1.16

require (
	github.com/container-storage-interface/spec v1.4.0
	github.com/enix/dothill-api-go v1.7.1
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/kubernetes-csi/csi-lib-iscsi v0.0.0-20200118015005-959f12c91ca8
	github.com/kubernetes-csi/csi-test v0.0.0-20191016154743-6931aedb3df0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	google.golang.org/grpc v1.29.1
	k8s.io/klog v1.0.0
)

// replace github.com/enix/dothill-api-go => ./pkg/dothill-api-go

// replace github.com/kubernetes-csi/csi-lib-iscsi => ./pkg/csi-lib-iscsi
replace github.com/kubernetes-csi/csi-lib-iscsi => github.com/enix/csi-lib-iscsi v0.0.0-dothill-3-1-6-2
