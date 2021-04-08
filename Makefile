ifndef DOCKER_HUB_REPOSITORY
	DOCKER_HUB_REPOSITORY = docker.io/enix
endif

ifndef VERSION
	VERSION = latest
else
	VERSION_FLAG = -X github.com/enix/dothill-csi/pkg/common.Version=$(VERSION)
endif

ifndef BIN
	BIN = dothill
endif

IMAGE = $(DOCKER_HUB_REPOSITORY)/dothill-csi:$(VERSION)

all:		bin image
.PHONY: all

bin: controller node
.PHONY: bin

controller:
	go build -v -ldflags "$(VERSION_FLAG)" -o $(BIN)-controller ./cmd/controller
.PHONY: controller

node:
	echo "$(VERSION_FLAG)"
	go build -v -ldflags "$(VERSION_FLAG)" -o $(BIN)-node ./cmd/node
.PHONY: node

test:
	./test/sanity
.PHONY: test

image:
	docker build -t $(IMAGE) --build-arg version="$(VERSION)" --build-arg vcs_ref="$(shell git rev-parse HEAD)" --build-arg build_date="$(shell date --rfc-3339=seconds)" .
.PHONY: image

push:		image
	docker push $(IMAGE)
.PHONY: push

clean:
	rm -vf $(BIN)-controller $(BIN)-node
.PHONY: clean
