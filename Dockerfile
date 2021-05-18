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
      org.opencontainers.image.licenses="MIT"

COPY host-chrooted.sh /usr/local/bin/

ENV CHROOTED_BINARIES='scsi_id,iscsiadm,multipath,multipathd,lsblk,blockdev,findmnt,mount,umount,mountpoint,resize2fs,e2fsck,blkid,mkfs.ext4'

RUN chmod +x /usr/local/bin/host-chrooted.sh \
 && bash -c 'cd /usr/local/bin; for bin in ${CHROOTED_BINARIES//,/ }; do ln -s ./host-chrooted.sh ./$bin; done'

COPY --from=build /dothill-* /usr/local/bin/

ENV PATH="${PATH}:/lib/udev"

CMD [ "/usr/local/bin/dothill-controller" ]

ARG version
ARG vcs_ref
ARG build_date
LABEL org.opencontainers.image.version="$version" \
      org.opencontainers.image.revision="$vcs_ref" \
      org.opencontainers.image.created="$build_date"
