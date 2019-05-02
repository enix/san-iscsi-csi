ifeq ($(REGISTRY),)
	REGISTRY = docker-registry.enix.io/enix
endif

ifeq ($(VERSION),)
	VERSION = latest
endif

ifeq ($(BIN),)
	BIN = msa-provisioner
endif

IMAGE = $(REGISTRY)/$(BIN):$(VERSION)

SRC		=	src/main.go src/provisioner.go

all:		image
.PHONY: all

bin:
	go build -v -o $(BIN)	$(SRC)
.PHONY: bin

image:
	docker build -t $(IMAGE) .
.PHONY: image

push:		image
	docker push $(IMAGE)
.PHONY: push

clean:
	rm -f $(BIN)
.PHONY: clean
