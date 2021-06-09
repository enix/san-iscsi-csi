# Copyright (c) 2021 Enix, SAS
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
# or implied. See the License for the specific language governing
# permissions and limitations under the License.
#
# Authors:
# Paul Laffitte <paul.laffitte@enix.fr>
# Arthur Chaloin <arthur.chaloin@enix.fr>
# Alexandre Buisine <alexandre.buisine@enix.fr>

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
