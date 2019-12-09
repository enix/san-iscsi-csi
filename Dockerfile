FROM golang:1.12-alpine3.9 AS build

ARG version

RUN apk add --update make

COPY . .

RUN echo -e "package main\nconst version = \"${version}\"" > src/version.go

RUN BIN="/dothill-provisioner" make bin

###########################################

FROM alpine:3.7

COPY --from=build /dothill-provisioner /usr/local/bin/dothill-provisioner

RUN chmod +x /usr/local/bin/dothill-provisioner

ENTRYPOINT [ "/usr/local/bin/dothill-provisioner" ]
