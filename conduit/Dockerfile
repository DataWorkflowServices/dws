FROM ghcr.io/lunarmodules/busted:master as base

WORKDIR /

COPY burst_buffer.lua burst_buffer.lua

FROM base as testing-base

RUN apk add --update --no-cache python3 py3-pip

RUN pip3 install openapi-schema-validator pyyaml

FROM testing-base as testing

ARG DWS_CRD_DIR

WORKDIR /crds

COPY crd/bases/* /crd/bases/

WORKDIR /

COPY validate.py validate.py

COPY test.lua test.lua

RUN busted /test.lua

From base as production
