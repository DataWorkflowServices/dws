This directory holds a set of tests for the k8s workflow REST api for the
create, get, update, patch, and delete verbs.  It assumes that DWS workflow
CRD has been created and the dws-operator has been deployed and is running.

They need to be run on node-1 of the k8s cluster with kubectl set up to run
in 'proxy' mode.  This requirement will go away once RBAC support in DWS hs
been added.

Remember that in the kind environment, the default namespace is 'default', so
for the `delete` script, you need:
$ env NAMESPACE=default ./delete dws-workflow-test-68240

To start the the kubectl proxy:

$ kubectl proxy --port 8080 &

NOTE: If you tear down you kind cluster, you likely need to kill kubectl proxy as well.:wq

Run the tests.
