module github.com/enix/dothill-storage-controller/pkg/node

go 1.12

require (
	github.com/container-storage-interface/spec v1.2.0
	github.com/enix/dothill-storage-controller/pkg/common v0.0.0-00010101000000-000000000000
	github.com/kubernetes-csi/csi-lib-iscsi v0.0.0-20200118015005-959f12c91ca8
	github.com/pkg/errors v0.9.1
	google.golang.org/grpc v1.26.0
	k8s.io/klog v1.0.0
)

replace (
	github.com/enix/dothill-storage-controller/pkg/common => ../common

	github.com/kubernetes-csi/csi-lib-iscsi => github.com/27149chen/csi-lib-iscsi v0.0.0-20200113115836-da1b94e79a4c
)
