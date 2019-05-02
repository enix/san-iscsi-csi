FROM instrumentisto/dep:0.5-alpine AS build

WORKDIR /go/src/enix.io/dothil-provisioner

RUN apk add --update make

COPY . .

RUN dep ensure

RUN BIN="/go/bin/dothill-provisioner" make bin

###########################################

FROM alpine:3.7

COPY --from=build /go/bin/dothill-provisioner /usr/local/bin/dothill-provisioner

RUN chmod +x /usr/local/bin/dothill-provisioner

ENTRYPOINT [ "/usr/local/bin/dothill-provisioner" ]
