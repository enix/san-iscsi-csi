FROM golang:1.12-alpine3.9 AS build

ARG version

RUN apk add --update make git

WORKDIR /app

COPY . .

RUN echo -e "package common\nconst Version = \"${version}\"" > pkg/common/version.go

RUN BIN="/dothill" make controller

RUN BIN="/dothill" make node

###########################################

FROM alpine:3.7

COPY --from=build /dothill-controller /usr/local/bin/dothill-controller

RUN chmod +x /usr/local/bin/dothill-controller

ENTRYPOINT [ "/usr/local/bin/dothill-controller" ]
