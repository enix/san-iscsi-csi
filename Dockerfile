FROM golang:1.12-buster AS build

ARG version

RUN apt update \
 && apt install -y make git \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./go.* ./
COPY ./pkg/controller/go.* ./pkg/controller/
COPY ./pkg/node/go.* ./pkg/node/
COPY ./pkg/common/go.* ./pkg/common/

RUN go mod download

COPY cmd cmd
COPY pkg pkg
COPY Makefile ./

RUN echo "package common\nconst Version = \"${version}\"" > ./pkg/common/version.go

RUN BIN="/dothill" make controller

RUN BIN="/dothill" make node

###########################################

FROM debian:buster

RUN apt update \
 && apt install -y dosfstools e2fsprogs xfsprogs jfsutils libisns0 \
 && rm -rf /var/lib/apt/lists/*

COPY --from=build /dothill-* /usr/local/bin/

CMD [ "/usr/local/bin/dothill-controller" ]
