FROM golang:1.12-alpine3.9 AS build

ARG version

RUN apk add --update make git

WORKDIR /app

COPY ./go.* ./
COPY ./pkg/controller/go.* ./pkg/controller/
COPY ./pkg/node/go.* ./pkg/node/
COPY ./pkg/common/go.* ./pkg/common/

RUN go mod download

COPY cmd cmd
COPY pkg pkg
COPY Makefile ./

RUN echo -e "package common\nconst Version = \"${version}\"" > pkg/common/version.go

RUN BIN="/dothill" make controller

RUN BIN="/dothill" make node

###########################################

FROM alpine:3.7

RUN apk add --update dosfstools e2fsprogs xfsprogs jfsutils

RUN echo -e '#! /bin/sh\nchroot /host /usr/bin/env -i PATH="/bin:/sbin:/usr/bin" iscsiadm $@' > /usr/local/bin/iscsiadm

RUN echo -e '#! /bin/sh\nchroot /host /usr/bin/env -i PATH="/bin:/sbin:/usr/bin" multipath $@' > /usr/local/bin/multipath

COPY --from=build /dothill-* /usr/local/bin/

RUN chmod +x /usr/local/bin/iscsiadm /usr/local/bin/multipath /usr/local/bin/dothill-*

CMD [ "/usr/local/bin/dothill-controller" ]
