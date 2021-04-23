# Build
To build locally (outside of Jenkins), you will need docker and Make packages
installed nad have access to DTR.  For all commands, you must be at the top of the source tree.

To create a deployable image with the dws-operator installed (update version as appropriate):
```
$ make
```
Or
```
$ make image
```
Or
```
$ docker build -f build/Dockerfile --label arti.dev.cray.com/cray/dws-operator:0.0.1 \
	-t arti.dev.cray.com/cray/dws-operator:0.0.1 .
```
To rebuild the operator-sdk auto-generated source after updating api/controller definitions:
```
$ make code-generation
```

To re-format source to meet go fmt conventions:
```
$ make fmt
```

To clean/remove all images:
```
$ make clean
```

# Test
[Testing using Kind](test/k8s/README.md)
