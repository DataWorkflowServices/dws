#!/bin/bash

export CGO_ENABLED=0

apk add git
if [ $? -ne 0 ] ; then
	echo "pre-run setup error"
	exit 1
fi

go get -u golang.org/x/lint/golint
if [ $? -ne 0 ] ; then
	echo "pre-run setup error, unable to fetch golint"
	exit 1
fi

GOLINTBIN=$(dirname "$(go list -f "{{.Target}}" golang.org/x/lint/golint)")
export PATH=$PATH:$GOLINTBIN

# Check non-vendor packages
mypkgs=$(go list -f "{{.Dir}}" ./...)
if [ -z "${mypkgs}" ] ; then
	echo "Exiting: No packages to check"
	exit 0
fi

# Check formatting wiithin those packages
echo "Checking non-vendor packages - code format"
foundCount=0
for p in $mypkgs; do
	# If gofmt returns the file name, there is a formatting issue with it.
	formatResult=$(gofmt -l "$p")
	if [ -n "$formatResult" ] ; then
		echo "go fmt needed:  $formatResult"
		foundCount=$((foundCount + 1))
	fi
done

if [ $foundCount -ne 0 ] ; then
	echo "go fmt failed"
	exit 1
fi

# Check non-vendor package files ignoring all generated files
echo "Running golint"
for pkg in $mypkgs; do
	mypkgfiles=$(find "${pkg##*dws-operator/}" -maxdepth 1 -type f \( ! -iname "*zz_generated*" \))
	printf "checking $mypkgfiles\n"
	golint -min_confidence 0.8 -set_exit_status "$mypkgfiles"
	if [ $? -ne 0 ] ; then
		echo "go lint failed"
		exit 1
	fi
done

echo "success"
exit 0
