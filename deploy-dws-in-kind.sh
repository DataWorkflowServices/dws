#!/bin/bash

# This script is the original deployment script for use with DWS and CDS.
# To deploy dws-operator with nnf, use the nnf-sos/playground.sh script.

#set -vx
# Immediately exit in wake of failure
set -e

action_create_deployment="create_deployment"
action_recreate_deployment="recreate_deployment"
action_create_kind_cluster="create_kind_cluster"
action_delete_kind_cluster="delete_kind_cluster"
action_deploy="deploy"
action_rollback="rollback"
action_reapply_charts="reapply_charts"

docker_image_name_prefix="arti.dev.cray.com/cray/cray-"

 # Build things in the current branch
 # Older implementations created this directory, 'git cloned' the branch into it,
 # built it from there and then deployed. Setting workdir to ./ builds the current branch
 # without the clone operations.
workdir=${workdir:-'./'}
helm_dir=${helm_dir:-$workdir}
kind_config=${kind_config:-kind-1-2.yaml}

# Manifests of images.
# the 'kj' product.
main_manifest=${main_manifest:-'main-manifest.txt'}

required_images=('dws-operator'
    'dws-admission-webhook')

cds_images=('dws-cds-driver')
nnf_images=()

# Version/tag of images that helm charts expects
deploy_tag=${deploy_tag:-'latest'}

# Switches to include various components
include_cds=${include_cds:-'false'}
clean_docker_cache=${clean_docker_cache:-'false'}

# Used to retrieve charts
default_sdp_manifests=('http://car.dev.cray.com/artifactory/kj/DPM/noos/noarch/dev/master/ssm-team/sdp_manifest/data-workflow-services_1.0.0_metadata.json'
                       'http://car.dev.cray.com/artifactory/kj/DPM/noos/noarch/dev/master/ssm-team/sdp_manifest/dws-cds-integration_1.0.0_metadata.json')

sdp_manifests=("${sdp_manifests:-${default_sdp_manifests[@]}}")

# Don't deploy these charts
chart_filter=${chart_filter:-'.*-deployment-group-manager'}

show_parameters ()
{
    echo
    echo "parameters for:         $0"
    echo "    workdir             \"$workdir\""
    echo "    main_manifest       \"$workdir/$main_manifest\""
    echo "    deploy_tag          \"$deploy_tag\""
    echo "    include_cds         \"$include_cds\""
    echo "    clean_docker_cache  \"$clean_docker_cache\""
    echo
}

usage ()
{
    echo "usage: $0 <OPTIONS> action"
    echo ""
    echo "  Actions"
    echo "    $action_create_deployment     Creates a new deployment using the latest from DWS"
    echo "    $action_recreate_deployment   Deletes existing deployment then creates a new deployment using prior configuration"
    echo "    $action_create_kind_cluster   Creates a Kind cluster without DWS deployments"
    echo "    $action_delete_kind_cluster   Deletes the cluster"
    echo "    $action_deploy                Deploys DWS"
    echo "    $action_rollback              Removes DWS deployments"
    echo ""
    echo "  Options:"
    echo "    -s/--show_parameters     Display parameter variables"
    echo "    -c/--clean_docker_cache  Find images in manifest files and remove them from the docker cache."
    echo "                             Used in conjunction with \"$main_manifest\"."
    echo "    --include_cds            Include CDS componenets.  Default: \"$include_cds\""
    echo ""
    echo "  Parameters"
    echo "    main_manifest     Location of the main manifest file.  Default: \"$main_manifest\""
    echo "    sdp_manifests     Contains helm metadata"
    echo "    workdir           Working directory to hold things like manifests, repos, and image artifacts.  Default: \"$workdir\""
    echo "    deploy_tag        Version/tag of images that helm charts expects.  Default: \"$deploy_tag\""
    echo ""
    echo "  NOTE"
    echo "    Once you've created a deployment for cds, nnf, or both, rollback/deploy must include the same components."
    echo ""
    echo "    These are allowed:"
    echo "      $0 --include_cds $action_create_deployment"
    echo "      $0 --include_cds $action_deploy"
    echo ""
    echo "  Examples"
    echo "    To deploy DWS in Kind using the latest images from artifactory, run:"
    echo "      $0 $action_create_deployment"
    echo ""
    echo "    To deploy with a different artifactory image.  Update the image version in $workdir/$main_manifest"
    echo "    and run:"
    echo "      $0 <--include_cds> $action_deploy"
    echo ""
    echo "      Note:  To deploy a custom image, load the image into your local docker cache then update $workdir/$main_manifest"
    echo "             with the docker name and tag of the custom image.  Finally, run the deploy command."
    echo ""
    echo "     To use a development branch, <branchToUse>, other than 'master' in order to pick up helm charts:"
    echo "     NOTE: helm charts are installed from the repos in $workdir ..."
    echo ""
    echo "      From your working repo where you have committed your changes to <branchToUse>:"
    echo "      git push <branchToUse>                      #to save your changes to stash/github"
    echo "      cd $workdir/<repo>"
    echo "      git checkout <branchToUse>"
    echo ""
    echo "    To start from scratch, remove the working directory and appropriate images from docker cache.  Then create a new deployment:"
    echo "      $0 $action_rollback"
    echo "      rm -rf $workdir"
    echo "      $0 <--include_cds/--include-nnf> $action_create_deployment"
    echo ""
    echo "    To set a different workdir for a deployment:"
    echo "      export workdir=nonDefaultWorkDirectory; $0 $action_create_deployment"
    echo ""
    echo "    To persist your include_cds setting in your environment so you don't have to specify it with each $action_create_deployment or $action_deploy:"
    echo "      export include_cds=true"
    echo ""
    echo ""
}

verify_k8s_is_running ()
{
    if kind get clusters 2>&1 >/dev/null | grep -q 'No kind clusters found'; then
        echo "No kind clusters found.  If you would like to create a cluster.  Run this tool using '$action_create_kind_cluster' or '$action_create_deployment'"
        exit 1
    fi
}
verify_k8s_is_not_running ()
{
    if kind get clusters 2>&1  | grep -qv 'No kind clusters found'; then
        echo "Existing cluster found. Run using '$action_recreate_deployment' or '$action_delete_kind_cluster' followed by '$action_create_deployment'"
        exit 1
    fi
}
delete_kind_cluster ()
{
    kind delete cluster || exit
}

create_kind_cluster ()
{
    verify_k8s_is_not_running
    # Create a config file used to create the Kind cluster,
    # if it doesn't already exist
    if [ ! -f "$kind_config" ]; then
        cat << EOF > "$kind_config"
# three node (two workers) cluster config
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerAddress: "127.0.0.1"
nodes:
- role: control-plane
- role: worker
- role: worker
EOF
    fi
    kind create cluster --image=kindest/node:v1.20.0 --config="$kind_config"		# see embedded config file above

    while kubectl get nodes | grep -q NotReady
    do
        echo "Waiting 5 seconds for nodes to become Ready"
        sleep 5
    done

    kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml

    while kubectl get nodes | grep -q NotReady
    do
        echo "Waiting 5 seconds for nodes to become Ready"
        sleep 5
    done

    echo "Kind cluster is ready"
}

# The method takes an array of strings that represents the docker images to clean from the cache
clean_docker_cache ()
{
    echo "Cleaning docker cache"
    local images_to_clean=$1

    for image_base_name in "${images_to_clean[@]}";
    do
        basic_name=$image_base_name

        # Get the name and version from the manifest
        docker_image_name=$(grep -h "$basic_name" "$main_manifest" | awk '{print $1}' | sed 's|.$||') || exit
        ver=$(grep -h "$basic_name" "$main_manifest" | awk '{print $2}') || exit
        if [ -z "$docker_image_name" ] || [ -z "$ver" ]; then
            echo "Failed to retrieve info for $basic_name ($docker_image_name:$ver)"
            continue
        fi

        docker rmi "$docker_image_name:$ver"
        docker rmi "$docker_image_name_prefix$image_base_name:$deploy_tag"
    done
}

rollback ()
{
    verify_k8s_is_running
    make undeploy
}

deploy_images ()
{
    make docker-build
    make kind-push
}

deploy ()
{
    verify_k8s_is_running

    deploy_images
    sleep 10
    make deploy
}


ready_to_run_initial_checks ()
{
    # Check other required dependencies
    type wget >/dev/null 2>&1 || { echo "wget is required but not found.  Exiting" ; exit; }
    type jq   >/dev/null 2>&1 || { echo "jq is required but not found.  Exiting" ; exit; }

    mkdir -p "$workdir" || exit
    cd "$workdir" || exit
}


# Command line handling
die() { echo "$*" >&2; exit 2; }  # complain to STDERR and exit with error
needs_arg() { if [ -z "$OPTARG" ]; then die "No arg for --$OPT option"; fi; }
while getopts "cs-:" OPT; do
    # support long options: https://stackoverflow.com/a/28466267/519360
    if [ "$OPT" = "-" ]; then   # long option: reformulate OPT and OPTARG
        OPT="${OPTARG%%=*}"       # extract long option name
        OPTARG="${OPTARG#$OPT}"   # extract long option argument (may be empty)
        OPTARG="${OPTARG#=}"      # if long option argument, remove assigning `=`
    fi

    # 'clean_docker_cache' not defined in getopts options list - disable shellcheck
    # shellcheck disable=SC2214
    case "$OPT" in
        c | clean_docker_cache ) clean_docker_cache="true" ;;
        s | show_parameters )    display_parameters="true" ;;
        include_cds )            include_cds="true" ;;
        ??* )                    die "Illegal option --$OPT" ;;  # bad long option
        ? )                      exit 2 ;;  # bad short option (error reported via getopts)
    esac
done
shift $((OPTIND-1)) # remove parsed options and args from $@ list


# Verify we are ready to proceed
ready_to_run_initial_checks

action="$1"
case "$action" in
    "$action_create_deployment")
        create_kind_cluster
        deploy
        ;;
    "$action_recreate_deployment")
        delete_kind_cluster
        create_kind_cluster
        deploy
        ;;
    "$action_create_kind_cluster")
        create_kind_cluster
        ;;
    "$action_delete_kind_cluster")
        delete_kind_cluster
        ;;
    "$action_deploy")
        rollback
        deploy
        apply_charts
        ;;
    "$action_rollback")
        rollback
        ;;
    *)
        usage
esac