# Overview
DWS can be run in Kind,  or "kubernetes in docker".  Kind will give you a set of nodes with k8s running on them.
The following instructions are for use on a Macbook.

# Install Software
Install kind: 
```
$ brew install kind
```
We'll use helm charts for DWS, and shasta-land requires helm version 2:
```
$ brew install helm@2"
```
You must also have docker installed and running.

# The auto-deploy tool
## Preparing for the auto-deploy tool
Add the following to your ~/.gitconfig file, to allow the deploy tool to clone DWS repositories:
```
[url "ssh://git@stash.us.cray.com:7999/"]
    insteadOf = https://stash.us.cray.com/scm/
```

## Running the auto-deploy tool
Create your kind environment with DWS.  This will configure and start kind; pull and load containers for DWS; pull helm charts for DWS; deploy all helm charts; and leave you with a running system.
```
$ test/k8s/deploy-dws-in-kind.sh create_deployment
helm is /usr/local/opt/helm@2/bin/helm
wget is /usr/local/bin/wget
jq is /usr/local/bin/jq
Deleting cluster "kind" ...
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.19.1) 🖼
 ✓ Preparing nodes 📦 📦 📦 
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
 ✓ Joining worker nodes 🚜
[...]
```

## The Captured Environment
The deploy-dws-in-kind tool created a directory named deploy-in-kind-work where it has saved all of the containers and helm charts that it used.  If you preserve this directory then it will be used the next time you deploy an environment, which means a quick start-up.  It also means you can save a combination of containers and charts that you know works, should you need that ability.

## About That KIND Version
Notice in the output of the deploy tool that it pulled the kind container.  In this case it's pulling version v1.19.1.  It's going to pull every time you start kind, which is a pain on VPN, so let's fix that.  Check your local docker cache and you'll see it doesn't include that tag:
```
$ docker images | grep kind
kindest/node                                           <none>                          37ddbc9063d2   6 months ago    1.33GB
```
So let's tag it:
```
$ docker tag 37ddbc9063d2  kindest/node:v1.19.1
$ docker images | grep kind
kindest/node                                           v1.19.1                         37ddbc9063d2   6 months ago    1.33GB
```
Now when you start kind it'll use the cached container, and won't pull a new container unless there's an update.

## Clean The DWS Environment
Clean up the DWS resources and containers used in an active kind environment so you can redeploy a new DWS environment without restarting a new kind cluster.
test/k8s/deploy-dws-in-kind.sh rollback

## Deploy DWS Again
If you have done a rollback then you still have a working kind cluster and you can re-deploy the DWS containers and charts.
test/k8s/deploy-dws-in-kind.sh deploy

## Destroy the DWS Environment and KIND cluster
Destroy the entire KIND cluster and all DWS resources.  This preserves the deploy-in-kind-work directory so it can be re-used for the next deployment.
test/k8s/deploy-dws-in-kind.sh delete_kind_cluster

# References
