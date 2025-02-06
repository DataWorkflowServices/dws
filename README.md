# dws - Data Workflow Services

Data Workflow Services (DWS) contains the schema for the Data Workflow Services API. HPC batch software integrates with the DWS API on HPE HPC systems and HPC storage systems to provide intelligent data movement and ephemeral storage resources to user workloads.

## Contributing

Before opening an issue or pull request, please read the [Contributing] guide.

[contributing]: CONTRIBUTING.md

## Build

To build locally you need docker and Make packages installed.
For all commands, you must be at the top of the source tree.

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