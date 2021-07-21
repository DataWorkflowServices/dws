# dws operator

## Build

To build locally (outside of Jenkins), you will need docker and Make packages
installed nad have access to DTR.  For all commands, you must be at the top of the source tree.

To create a deployable image with the dws-operator installed (update version as appropriate):

```bash
make
```

Or

```bash
make image
```

Or

```bash
$ docker build -f build/Dockerfile --label arti.dev.cray.com/cray/dws-operator:0.0.1 \
     -t arti.dev.cray.com/cray/dws-operator:0.0.1 .
```

To rebuild the operator-sdk auto-generated source after updating api/controller definitions:

```bash
make code-generation
```

To run go test:

```bash
make test
```

To re-format source to meet go fmt conventions:

```bash
make fmt
```

To clean/remove all images:

```bash
make clean
```

## Test

[Testing using Kind](test/k8s/README.md)
