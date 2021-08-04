# dws operator

## Build

To build locally (outside of Jenkins), you will need docker and Make packages
installed.  For all commands, you must be at the top of the source tree.

To create a deployable image with the dws-operator installed and rebuild:

```bash
make docker-build
```

To push the docker image to a kind environment:

```bash
make kind-push
```

To deploy dws-operator to a kind environment

```bash
make deploy
```

To remove dws-operator from a kind environment

```bash
make undeploy
```

To rebuild the operator-sdk auto-generated source after updating api/controller definitions:

```bash
make generate
```

To rebuild the operator-sdk auto-generated yaml files after updating api/controller definitions:

```bash
make manifests
```
