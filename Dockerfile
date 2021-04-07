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

RUN echo "package common\nconst Version = \"${version}\"" > ./pkg/common/version.go

RUN BIN="/dothill" make controller

RUN BIN="/dothill" make node

###########################################

FROM ubuntu:18.04

LABEL maintainer="Enix <no-reply@enix.fr>" \
      org.opencontainers.image.title="Dothill CSI" \
      org.opencontainers.image.description="A dynamic persistent volume provisioner for Dothill AssuredSAN based storage systems." \
      org.opencontainers.image.url="https://github.com/enix/dothill-csi" \
      org.opencontainers.image.source="https://github.com/enix/dothill-csi/blob/master/Dockerfile" \
      org.opencontainers.image.documentation="https://github.com/enix/dothill-csi/blob/master/README.md" \
      org.opencontainers.image.authors="Enix <no-reply@enix.fr>" \
      org.opencontainers.image.licenses="MIT"

RUN apt update \
 && apt dist-upgrade -y \
 && apt install -y dosfstools e2fsprogs xfsprogs jfsutils libisns0 open-iscsi kmod multipath-tools \
 && rm -rf /var/lib/apt/lists/*

COPY --from=build /dothill-* /usr/local/bin/

CMD [ "/usr/local/bin/dothill-controller" ]

ARG version
ARG vcs_ref
ARG build_date
LABEL org.opencontainers.image.version="$version" \
      org.opencontainers.image.revision="$vcs_ref" \
      org.opencontainers.image.created="$build_date"
