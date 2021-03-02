ARG GO_VER
FROM golang:${GO_VER}

ARG PROTOC_VER
ARG PROTOC_PLATFORM
ARG GOGO_PROTO_TAG
ARG UID
ARG GID

ENV TARBALL protoc-${PROTOC_VER}-${PROTOC_PLATFORM}.zip
ENV GOGOPROTO_ROOT /go/src/github.com/gogo/protobuf
ENV PROTOC_URL https://github.com/google/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-${PROTOC_PLATFORM}.zip

RUN set -ex && apt-get -q -y update --fix-missing && apt-get -q -y install unzip dumb-init libc6

RUN set -ex && \
    getent group  $GID || groupadd builder --gid=$GID -o; \
    getent passwd $UID || useradd builder --uid=$UID --gid=$GID --create-home --shell=/bin/sh;

RUN set -ex && (mkdir -p /opt/protoc && \
     mkdir -p /.cache && \
     chown -R $UID:$GID /go && \
     chown -R $UID:$GID /opt/protoc && \
     chmod 777 /.cache && \
     chmod 777 /tmp)

USER $UID:$GID

ENV LANGUAGE="en_US.UTF-8" \
     LANG="en_US.UTF-8" \
     LC_ALL="en_US.UTF-8" \
     LC_CTYPE="en_US.UTF-8" \
     PATH="$PATH:/opt/protoc/bin:/go/bin"

RUN set -ex && \
     wget --quiet -O /tmp/${TARBALL} ${PROTOC_URL} && \
     unzip -d /opt/protoc /tmp/${TARBALL} && \
     mkdir -p /go/src/github.com/gogo/ && \
     git clone https://github.com/gogo/protobuf --branch ${GOGO_PROTO_TAG} /go/src/github.com/gogo/protobuf && \
     cd /go/src/github.com/gogo/protobuf && make install

ENV PROTO_INCLUDE "/usr/local/include":"/go/src"

