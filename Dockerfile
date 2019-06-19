FROM enix/go-dep:0.5 AS build

ARG version

WORKDIR /go/src/enix.io/dothil-provisioner

RUN apk add --update make

COPY . .

RUN dep ensure

RUN echo -e "package main\nconst version = \"${version}\"" > src/version.go

RUN BIN="/go/bin/dothill-provisioner" make bin

###########################################

FROM alpine:3.7

COPY --from=build /go/bin/dothill-provisioner /usr/local/bin/dothill-provisioner

RUN chmod +x /usr/local/bin/dothill-provisioner

ENTRYPOINT [ "/usr/local/bin/dothill-provisioner" ]
