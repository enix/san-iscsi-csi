module github.com/enix/dothill-csi

go 1.12

require (
	github.com/container-storage-interface/spec v1.3.0
	github.com/enix/dothill-api-go v1.5.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/kubernetes-csi/csi-lib-iscsi v0.0.0-20200118015005-959f12c91ca8
	github.com/kubernetes-csi/csi-test v0.0.0-20191016154743-6931aedb3df0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	google.golang.org/grpc v1.29.1
	k8s.io/klog v1.0.0
)

// replace github.com/kubernetes-csi/csi-lib-iscsi => ./pkg/csi-lib-iscsi
replace github.com/kubernetes-csi/csi-lib-iscsi => github.com/enix/csi-lib-iscsi v0.0.0-dothill-3-0-3
