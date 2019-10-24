#!/bin/bash

echo "Running go vet"
go vet -c=5 ./pkg/... ./cmd/...
if [ $? -ne 0 ] ; then
	echo "failed"
	exit 1
fi

echo "success"
exit 0
