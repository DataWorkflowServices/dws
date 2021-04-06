To build locally (outside of Jenkins), you will need docker and Make packages
installed nad have access to DTR.  For all commands, you must be at the top of the source tree.

You must be logged into DTR to access the containers for any of these build options to work:

```
$ docker login dtr.dev.cray.com
```

To rebuild the operator-sdk auto-generated source after updating api/controller definitions:
```
$ make code-generation
```

To re-format source to meet go fmt conventions:
```
$ make fmt
```

To create a deployable image with the dws-operator installed (update version as appropriate):
```
$ make image
```
Or
```
$ docker build -f build/Dockerfile --label dtr.dev.cray.com/${USER}/dws-operator:0.0.1 \
	-t dtr.dev.cray.com/${USER}/dws-operator:0.0.1 .
```

To clean/remove all images:
```
$ make clean
```


