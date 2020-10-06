ifeq ($(REGISTRY),)
	REGISTRY = docker.io/enix
endif

ifeq ($(VERSION),)
	VERSION = latest
endif

ifeq ($(BIN),)
	BIN = dothill
endif

IMAGE = $(REGISTRY)/dothill-provisioner:$(VERSION)

all:		image
.PHONY: all

controller:
	go build -v -o $(BIN)-controller ./cmd/controller
.PHONY: bin

node:
	go build -v -o $(BIN)-node ./cmd/node
.PHONY: bin

image:
	docker build -t $(IMAGE) --build-arg version=$(VERSION) .
.PHONY: image

push:		image
	docker push $(IMAGE)
.PHONY: push

clean:
	rm -f $(BIN)
.PHONY: clean
