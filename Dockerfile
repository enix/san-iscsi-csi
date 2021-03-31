FROM golang:1.12-buster AS build

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

RUN echo "package common\nconst Version = \"${version}\"" > ./pkg/common/version.go

RUN BIN="/dothill" make controller

RUN BIN="/dothill" make node

###########################################

FROM ubuntu:18.04

RUN apt update \
 && apt dist-upgrade -y \
 && apt install -y dosfstools e2fsprogs xfsprogs jfsutils libisns0 open-iscsi kmod multipath-tools \
 && rm -rf /var/lib/apt/lists/*

COPY --from=build /dothill-* /usr/local/bin/

CMD [ "/usr/local/bin/dothill-controller" ]
