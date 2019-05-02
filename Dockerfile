FROM instrumentisto/dep:0.5-alpine AS build

WORKDIR /go/src/enix.io/msa

RUN apk add --update make

COPY . .

RUN dep ensure

RUN BIN="/go/bin/msa-provisioner" make bin

###########################################

FROM alpine:3.7

COPY --from=build /go/bin/msa-provisioner /usr/local/bin/msa-provisioner

RUN chmod +x /usr/local/bin/msa-provisioner

ENTRYPOINT [ "/usr/local/bin/msa-provisioner" ]
