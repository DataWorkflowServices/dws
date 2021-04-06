#!/bin/bash

# Copyright 2020 Hewlett Packard Enterprise Development LP

set -e
trap 'catch $? $LINENO' EXIT

catch() {
    if [ "$1" != "0" ]; then
        echo "Error $1 occurred on line $2"
    fi
}

(
    cd stash.us.cray.com/dpm/$1
    operator-sdk generate k8s
    operator-sdk generate crds

    cd deploy/crds
    CRDFILES=(*.yaml)
    cd ../..

    for fn in "${CRDFILES[@]}"
    do
        cp deploy/crds/${fn} kubernetes/$1/templates/${fn}
        cat << EOF | patch kubernetes/$1/templates/${fn}
4a5,11
>   labels:
>     app.kubernetes.io/name: {{ include "cray-service.name" . }}
>     {{- include "cray-service.common-labels" . | nindent 4 }}
>   annotations:
>     {{- include "cray-service.common-annotations" . | nindent 4 }}
>   annotations:
> {{ toYaml .Values.crdAnnotations | indent 4 }}
EOF
    done
)
