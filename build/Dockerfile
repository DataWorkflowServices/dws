# Copyright 2019 Cray Inc.  All Rights Reserved
FROM golang:latest AS build_base

COPY cmd $GOPATH/src/stash.us.cray.com/dpm/dws-operator/cmd/
COPY pkg $GOPATH/src/stash.us.cray.com/dpm/dws-operator/pkg/
COPY vendor $GOPATH/src/stash.us.cray.com/dpm/dws-operator/vendor/
COPY build $GOPATH/src/stash.us.cray.com/dpm/dws-operator/build/

WORKDIR $GOPATH/src/stash.us.cray.com/dpm/dws-operator/

RUN go build -v -i -o build/bin/dws-operator cmd/manager/main.go
#RUN go build -v -i -o /usr/local/bin/dws-operator cmd/manager/main.go

# The base testing container
FROM build_base AS base_testing

COPY test/ /test/
COPY test/ $GOPATH/src/stash.us.cray.com/dpm/dws-operator/test/

# The testing artifact container
FROM base_testing AS testing
WORKDIR $GOPATH/src/stash.us.cray.com/dpm/dws-operator/
ENTRYPOINT ["sh", "test/docker_test_entry.sh"]

# The final application release product container
#FROM build_base as application
FROM registry.access.redhat.com/ubi7/ubi-minimal:latest as application

WORKDIR /root

ENV OPERATOR=/usr/local/bin/dws-operator \
    USER_UID=1001 \
    USER_NAME=dws-operator

# install operator binary
ENV GOPATH /go
COPY --from=build_base $GOPATH/src/stash.us.cray.com/dpm/dws-operator/build/bin/dws-operator ${OPERATOR}
COPY --from=build_base $GOPATH/src/stash.us.cray.com/dpm/dws-operator/build/bin /usr/local/bin

RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
