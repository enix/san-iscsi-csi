module github.com/enix/dothill-storage-controller

go 1.12

replace (
	github.com/enix/dothill-api-go => ../dothill-api-go
	github.com/enix/dothill-storage-controller/pkg/common => ./pkg/common
	github.com/enix/dothill-storage-controller/pkg/controller => ./pkg/controller
	github.com/enix/dothill-storage-controller/pkg/node => ./pkg/node
)

require (
	cloud.google.com/go v0.38.0 // indirect
	github.com/container-storage-interface/spec v1.2.0
	github.com/enix/dothill-storage-controller/pkg/common v0.0.0-00010101000000-000000000000
	github.com/enix/dothill-storage-controller/pkg/controller v0.0.0-00010101000000-000000000000
	github.com/enix/dothill-storage-controller/pkg/node v0.0.0-00010101000000-000000000000
	github.com/google/btree v1.0.0 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/kubernetes-csi/csi-test v2.2.0+incompatible
	github.com/kubernetes-csi/external-resizer v0.4.0 // indirect
	golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	google.golang.org/appengine v1.5.0 // indirect
	google.golang.org/grpc v1.26.0
)
