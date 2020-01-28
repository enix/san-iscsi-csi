module github.com/enix/dothill-storage-controller/pkg/controller

go 1.12

require (
	github.com/container-storage-interface/spec v1.2.0
	github.com/enix/dothill-api-go v1.4.1
	github.com/enix/dothill-storage-controller/pkg/common v0.0.0-00010101000000-000000000000
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	google.golang.org/grpc v1.26.0
	k8s.io/klog v1.0.0
)

replace github.com/enix/dothill-storage-controller/pkg/common => ../common
