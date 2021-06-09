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

FROM golang:1.16-buster AS build

RUN apt update \
 && apt install -y make git \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./go.* ./

RUN go mod download

COPY cmd cmd
COPY pkg pkg
COPY Makefile ./

ARG version

RUN BIN="/dothill" VERSION="$version" make controller
RUN BIN="/dothill" VERSION="$version" make node

###########################################

FROM ubuntu:18.04

LABEL maintainer="Enix <no-reply@enix.fr>" \
      org.opencontainers.image.title="Dothill CSI" \
      org.opencontainers.image.description="A dynamic persistent volume provisioner for Dothill AssuredSAN based storage systems." \
      org.opencontainers.image.url="https://github.com/enix/dothill-csi" \
      org.opencontainers.image.source="https://github.com/enix/dothill-csi/blob/master/Dockerfile" \
      org.opencontainers.image.documentation="https://github.com/enix/dothill-csi/blob/master/README.md" \
      org.opencontainers.image.authors="Enix <no-reply@enix.fr>" \
      org.opencontainers.image.licenses="Apache 2.0"

RUN apt update \
 && apt dist-upgrade -y \
 && apt install -y dosfstools e2fsprogs xfsprogs jfsutils libisns0 open-iscsi kmod multipath-tools \
 && rm -rf /var/lib/apt/lists/*

COPY --from=build /dothill-* /usr/local/bin/

ENV PATH="${PATH}:/lib/udev"

CMD [ "/usr/local/bin/dothill-controller" ]

ARG version
ARG vcs_ref
ARG build_date
LABEL org.opencontainers.image.version="$version" \
      org.opencontainers.image.revision="$vcs_ref" \
      org.opencontainers.image.created="$build_date"
