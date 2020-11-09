module github.com/enix/dothill-storage-controller

go 1.12

require (
	github.com/container-storage-interface/spec v1.2.0
	github.com/enix/dothill-api-go v1.5.0
	github.com/kubernetes-csi/csi-lib-iscsi v0.0.0-20200118015005-959f12c91ca8
	github.com/kubernetes-csi/csi-test v0.0.0-20191016154743-6931aedb3df0
	github.com/pkg/errors v0.9.1
	google.golang.org/grpc v1.26.0
	k8s.io/klog v1.0.0
)

replace (
	github.com/kubernetes-csi/csi-lib-iscsi => github.com/enix/csi-lib-iscsi 9fff3f45a09f1c4904d599172f2498015e85d27a
)
