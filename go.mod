module github.com/enix/dothill-storage-controller

go 1.12

require (
	github.com/enix/dothill-storage-controller/pkg/common v0.0.0-00010101000000-000000000000
	github.com/enix/dothill-storage-controller/pkg/controller v0.0.0-00010101000000-000000000000
	github.com/enix/dothill-storage-controller/pkg/node v0.0.0-00010101000000-000000000000
	github.com/kubernetes-csi/csi-test v2.2.0+incompatible
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	k8s.io/klog v1.0.0
)

replace (
	// github.com/enix/dothill-api-go => ../dothill-api-go
	github.com/enix/dothill-storage-controller/pkg/common => ./pkg/common
	github.com/enix/dothill-storage-controller/pkg/controller => ./pkg/controller
	github.com/enix/dothill-storage-controller/pkg/node => ./pkg/node

	github.com/kubernetes-csi/csi-lib-iscsi => github.com/27149chen/csi-lib-iscsi v0.0.0-20200113115836-da1b94e79a4c
)
