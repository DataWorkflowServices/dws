# Copyright 2019 Cray Inc.  All Rights Reserved
FROM dtr.dev.cray.com/alpine/golang:alpine AS build_base

COPY cmd $GOPATH/src/stash.us.cray.com/dpm/dws-operator/cmd/
COPY pkg $GOPATH/src/stash.us.cray.com/dpm/dws-operator/pkg/
COPY vendor $GOPATH/src/stash.us.cray.com/dpm/dws-operator/vendor/
COPY build $GOPATH/src/stash.us.cray.com/dpm/dws-operator/build/

WORKDIR $GOPATH/src/stash.us.cray.com/dpm/dws-operator/

RUN CGO_ENABLED=0 go build -v -i -o build/bin/dws-operator ./cmd/manager/main.go
RUN CGO_ENABLED=0 go build -v -i -o build/bin/dwsget ./cmd/dwsget/main.go
RUN CGO_ENABLED=0 go build -v -i -o build/bin/dwsinformer ./cmd/informer/main.go

# The base testing container
FROM build_base AS base_testing

COPY test/ $GOPATH/src/stash.us.cray.com/dpm/dws-operator/test/
COPY static-analysis/ $GOPATH/src/stash.us.cray.com/dpm/dws-operator/static-analysis/

# The testing artifact container
FROM base_testing AS testing
WORKDIR $GOPATH/src/stash.us.cray.com/dpm/dws-operator/
ENTRYPOINT ["sh", "test/docker_test_entry.sh"]

# The code coverage analysis container
FROM testing as coverage
WORKDIR $GOPATH/src/stash.us.cray.com/dpm/dws-operator
ENTRYPOINT ["sh", "test/docker_coverage_entry.sh"]

# the static-analysis-codestyle-container
FROM testing as codestyle
WORKDIR $GOPATH/src/stash.us.cray.com/dpm/dws-operator
ENTRYPOINT ["sh", "static-analysis/docker_codestyle_entry.sh"]

# the static-analysis-lint-container
FROM testing as lint
WORKDIR $GOPATH/src/stash.us.cray.com/dpm/dws-operator
ENTRYPOINT ["sh", "static-analysis/docker_lint_entry.sh"]

# The final application release product container
FROM dtr.dev.cray.com/baseos/alpine:latest as application

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
