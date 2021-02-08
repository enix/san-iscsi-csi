ifeq ($(DOCKER_HUB_REPOSITORY),)
	DOCKER_HUB_REPOSITORY = docker.io/enix
endif

ifeq ($(VERSION),)
	VERSION = latest
endif

ifeq ($(BIN),)
	BIN = dothill
endif

IMAGE = $(DOCKER_HUB_REPOSITORY)/dothill-csi:$(VERSION)

all:		bin image
.PHONY: all

bin: controller node
.PHONY: bin

controller:
	go build -v -o $(BIN)-controller ./cmd/controller
.PHONY: controller

node:
	go build -v -o $(BIN)-node ./cmd/node
.PHONY: node

test:
	./test/sanity
.PHONY: test

image:
	docker build -t $(IMAGE) --build-arg version=$(VERSION) .
.PHONY: image

push:		image
	docker push $(IMAGE)
.PHONY: push

clean:
	rm -vf $(BIN)-controller $(BIN)-node
.PHONY: clean
