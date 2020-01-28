module github.com/enix/dothill-storage-controller

go 1.12

require (
	github.com/enix/dothill-storage-controller/pkg/common v0.0.0-00010101000000-000000000000
	github.com/enix/dothill-storage-controller/pkg/controller v0.0.0-00010101000000-000000000000
	github.com/enix/dothill-storage-controller/pkg/node v0.0.0-00010101000000-000000000000
	github.com/kubernetes-csi/csi-test v2.2.0+incompatible
	honnef.co/go/tools v0.0.1-2019.2.3 // indirect
	k8s.io/klog v1.0.0
)

replace (
	github.com/enix/dothill-storage-controller/pkg/common => ./pkg/common
	github.com/enix/dothill-storage-controller/pkg/controller => ./pkg/controller
	github.com/enix/dothill-storage-controller/pkg/node => ./pkg/node
)
