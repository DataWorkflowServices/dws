#!/bin/bash

# Copyright 2021 Hewlett Packard Enterprise Development LP

set -e
trap 'catch $? $LINENO' EXIT

catch() {
    if [ "$1" != "0" ]; then
        echo "Error $1 occurred on line $2"
    fi
}

# We need to run the unit test within the docker container.
make test