#!/bin/bash

# Generate code coverage testing report.  This step is informative only
# at this point and will only fail if there is a problem when executing
# the go command.

echo "Code coverage for Go"
go test -cover ./pkg/... ./cmd/...
if [ $? -ne 0 ] ; then
	exit 1
fi

exit 0
