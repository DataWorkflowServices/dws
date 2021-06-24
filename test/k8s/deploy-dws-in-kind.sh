#!/bin/bash

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

workdir=${workdir:-'./work-deploy-in-kind'}
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
include_nnf=${include_nnf:-'false'}
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
    echo "    include_nnf         \"$include_nnf\""
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
    echo "    $action_reapply_charts        Removes charts and re-applies the charts in $workdir repos"
    echo ""
    echo "  Options:"
    echo "    -s/--show_parameters     Display parameter variables"
    echo "    -c/--clean_docker_cache  Find images in manifest files and remove them from the docker cache."
    echo "                             Used in conjunction with \"$main_manifest\"."
    echo "    --include_cds            Include CDS componenets.  Default: \"$include_cds\""
    echo "    --include_nnf            Include NNF componenets.  Default: \"$include_nnf\""
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
    echo "      $0 --include_nnf $action_create_deployment"
    echo "      $0 --include_nnf $action_deploy"
    echo "    ..or.."
    echo "      $0 --include_cds $action_create_deployment"
    echo "      $0 --include_cds $action_deploy"
    echo "    ..or.."
    echo "      $0 --include_nnf --include_cds $action_create_deployment"
    echo "      $0 --include_nnf --include_cds $action_deploy"
    echo ""
    echo "    These are NOT allowed:"
    echo "      $0 --include_nnf $action_create_deployment"
    echo "      $0 --include_cds $action_deploy"
    echo "    ..or.."
    echo "      $0 --include_cds $action_create_deployment"
    echo "      $0 --include_nnf $action_deploy"
    echo "    ..or.."
    echo "      $0 --include_nnf --include_cds $action_create_deployment"
    echo "      $0 --include_nnf  $action_deploy"
    echo ""
    echo "  Examples"
    echo "    To deploy DWS in Kind using the latest images from artifactory, run:"
    echo "      $0 $action_create_deployment"
    echo ""
    echo "    To deploy with a different artifactory image.  Update the image version in $workdir/$main_manifest"
    echo "    and run:"
    echo "      $0 <--include_cds/--include_nnf> $action_deploy"
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
    echo "    To persist your include_cds or include_nnf setting in your environment so you don't have to specify it with each $action_create_deployment or $action_deploy:"
    echo "      export include_<cds,nnf>=true"
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
    if kind get clusters 2>&1 >/dev/null | grep -qv 'No kind clusters found'; then
        echo "Existing cluster found. Run using '$action_recreate_deployment' or '$action_delete_kind_cluster' followed by '$action_create_deployment'"
        exit 1
    fi
}
delete_kind_cluster ()
{
    kind delete cluster || exit
}

wait_for_helm_tiller ()
{
    # Acquire tiller's Ready setting.
    # Looking for the 2 numbers in the READY column to match
    # afloeder@Anthonys-MacBook-Pro ~/dev1/dws-operator/test/k8s (issue/RABSW-xxx)$ kubectl get pods --namespace kube-system
    # NAME                                         READY   STATUS    RESTARTS   AGE
    # coredns-74ff55c5b-jzhvb                      1/1     Running   0          22m
    # coredns-74ff55c5b-scpc9                      1/1     Running   0          22m
    # etcd-kind-control-plane                      1/1     Running   0          22m
    # kindnet-95kc7                                1/1     Running   0          22m
    # kindnet-lsjtx                                1/1     Running   0          22m
    # kindnet-qw6tg                                1/1     Running   0          22m
    # kube-apiserver-kind-control-plane            1/1     Running   0          22m
    # kube-controller-manager-kind-control-plane   1/1     Running   0          22m
    # kube-proxy-7g2gk                             1/1     Running   0          22m
    # kube-proxy-lnkmp                             1/1     Running   0          22m
    # kube-proxy-nsg9f                             1/1     Running   0          22m
    # kube-scheduler-kind-control-plane            1/1     Running   0          22m
    # tiller-deploy-7b56c8dfb7-7mqcx            -->1/1<--  Running   0          22m
    helmReady=$(kubectl get pods --namespace kube-system | grep tiller | awk '{split($2,a,"/"); if(a[1] == a[2]) {print "Ready"} else {print "Waiting"};}')

    echo "Helm-Tiller: " "$helmReady"
    while [ "$helmReady" == "Waiting" ];
    do
        echo "Waiting 5 seconds for helm's tiller to be Ready"
        sleep 5

        helmReady=$(kubectl get pods --namespace kube-system | grep tiller | awk '{split($2,a,"/"); if(a[1] == a[2]) {print "Ready"} else {print "Waiting"};}')
    done
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

    kind create cluster --config="$kind_config"		# see embedded config file above
    helm init --service-account tiller --upgrade
    kubectl create serviceaccount -n kube-system tiller
    kubectl create clusterrolebinding tiller-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
    kubectl --namespace kube-system patch deploy tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'
    kubectl label node kind-worker  cray.dpm.dg.data-workflow-services=true
    kubectl label node kind-worker2 cray.dpm.dg.data-workflow-services=true

    while kubectl get nodes | grep -q NotReady
    do
        echo "Waiting 5 seconds for nodes to become Ready"
        sleep 5
    done

    # Before we can submit helm charts, we need to make sure Tiller is Ready
    wait_for_helm_tiller

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

    ordered_rollback_list=('dws-operator'
        'dws-validation-webhook'
        'dws-nnf-webhook'
        'dws-cds-webhook'
        'dws-cds-integration')

    for chart in "${ordered_rollback_list[@]}";
    do
        chartToDelete=$(helm list | grep "$chart" | awk '{print $1}')
        if [ "$chartToDelete" != "" ]; then
            helm delete --purge "$chartToDelete"
        fi
    done

    resource_kinds=$(kubectl get crds | grep -E 'dws.cray.hpe.com' | awk '{print $1}' | sed -e 's/.dws.cray.hpe.com//')
    for res_type in $resource_kinds
    do
        if ! kubectl get "$res_type" -n default > /dev/null 2>&1
        then
            continue
        fi
        if cnt=$(kubectl get "$res_type" -n default -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' | wc -l)
        then
            if (( cnt > 0 ))
            then
                echo "Force-delete $res_type"
                if ! kubectl get "$res_type" --no-headers | awk '{print $1}' | xargs kubectl patch "$res_type" --type merge -p '{"metadata":{"finalizers":null}}'
                then
                    echo "Failed to delete $res_type"
                fi
            fi
        else
            echo "Unable to count $res_type" > /dev/stderr
        fi
    done

    if [ "$clean_docker_cache" = "true" ]; then
        clean_docker_cache "${required_images[@]}"

        if [ "$include_cds" = "true" ]; then
            clean_docker_cache "${cds_images[@]}"
        fi

        if [ "$include_nnf" = "true" ]; then
            clean_docker_cache "${nnf_images[@]}"
        fi
    fi
}

deploy_images ()
{
    # see https://stackoverflow.com/questions/16461656/how-to-pass-array-as-an-argument-to-a-function-in-bash
    # for info on passing arrays into functions
    local image_array_name="$1[@]"
    local image_array=("${!image_array_name}")
    local deployment_nodes_name="$2[@]"
    local deployment_nodes=("${!deployment_nodes_name}")

    for image_base_name in "${image_array[@]}";
    do
        basic_name=$image_base_name

        # Get the name and version from the manifest
        docker_image_name=$(grep -h "$basic_name" "$main_manifest" | awk '{print $1}' | sed 's|.$||') || exit
        ver=$(grep -h "$basic_name" "$main_manifest" | awk '{print $2}') || exit

        if [ -z "$docker_image_name" ] || [ -z "$ver" ]; then
            echo "Failed to retrieve info for $basic_name ($docker_image_name:$ver)"
            exit
        fi

        # Check to see if it's already been loaded into the docker cache
        if [ -z "$(docker images -q "$docker_image_name:$ver" 2> /dev/null)" ]; then
            # Check to see if the image has been downloaded
            artifactory_name="cray-$image_base_name-$ver-dockerimage.tar"
            if [ ! -f "$artifactory_name" ]; then
                # Try to fetch the image from artifactory where the DWS images are kept
                if ! wget "http://car.dev.cray.com/artifactory/kj/DPM/noos/noarch/dev/master/ssm-team/$artifactory_name"; then
                    echo "Unable to find image: $artifactory_name"
                    exit
                fi
            fi
            # Load the image into the docker cache
            echo "Docker load $artifactory_name"
            docker load -i "$artifactory_name" || exit
        fi

        # Tag the image according to what the helm charts will look for, and load it onto the Kind nodes
        echo "Docker tag $docker_image_name:$ver $docker_image_name_prefix$image_base_name:$deploy_tag"
        docker tag "$docker_image_name:$ver" "$docker_image_name_prefix$image_base_name:$deploy_tag" || exit
        echo "Kind load --nodes ${deployment_nodes[*]} $docker_image_name_prefix$image_base_name:$deploy_tag"
        kind load docker-image --nodes "${deployment_nodes[*]}" "$docker_image_name_prefix$image_base_name:$deploy_tag" || exit
    done

}

deploy ()
{
    verify_k8s_is_running

    if [ ! -f "$main_manifest" ]; then
        # Fetch the manifest from artifactory and update each image in the file to reflect its name (i.e. "Repository") after loaded into docker
        wget https://arti.dev.cray.com/artifactory/kj-misc-master-local/manifest/manifest.txt -O "tmp-$main_manifest" || exit
        awk '{print "sms.local:5000/cray/" $0}' "tmp-$main_manifest" > "$main_manifest"
        rm "tmp-$main_manifest"
    fi

    DWS_NODES=$(kubectl get node -l cray.dpm.dg.data-workflow-services=true -o json | jq -r '.items[].metadata.name' | paste -s -d, -)
    echo "DWS_NODES: $DWS_NODES"

    deploy_images required_images DWS_NODES

    if [ "$include_cds" = "true" ]; then
        echo "deploying cds images"
        deploy_images cds_images DWS_NODES
    fi

    if [ "$include_nnf" = "true" ]; then
        echo "deploying nnf images"
        deploy_images nnf_images DWS_NODES
    fi

    if ! helm repo list | grep -q cray-internal; then
        helm repo add cray-internal http://helmrepo.dev.cray.com:8080 || exit
    fi
}

apply_charts ()
{
set -x

    # Get the helm charts and install them
    for sdp_manifest in "${sdp_manifests[@]}";
    do
        # shellcheck disable=SC2207
        charts_metadata=($(curl "$sdp_manifest" | jq | grep 'metadata.json' | grep -vE "$chart_filter" | sed 's|,||' | sed 's|"||g' | tr -d ' '))

        for helm_metadata in "${charts_metadata[@]}";
        do
            # Get the repo that contains the helm chart
            repo=$(curl "$helm_metadata" | jq -r '.git_clone_url')

            # We clone only the repos we need.
            # If a repo is independent of either cds or nnf, we clone it.
            # Otherwise, we check whether $include_cds and/or $include_nnf is set
            # to determine whether we clone the repo
            if [[ "$repo" =~ "cds" ]] && [[ "$include_cds" == false ]]; then
                echo "CDS repo, but skipping because include_cds: $include_cds"
                continue
            fi
            if [[ "$repo" =~ "nnf" ]] && [[ "$include_nnf" == false ]]; then
                echo "NNF repo, but skipping because include_nnf: $include_nnf"
                continue
            fi

            # Clone the repo if it doesn't exist locally
            repo_dir=$(basename "$repo" | sed 's|.git||')
            if [ ! -d "$repo_dir" ]; then
                git clone "$repo" || exit
            fi


            # We install only the helm charts we need
            # If a chart_dir is independent of either cds or nnf, we install it.
            # Otherwise, we check whether $include_cds and/or $include_nnf is set
            # to determine whether we install the helm chart
            chart_dir="$repo_dir"/kubernetes/$(curl "$helm_metadata" | jq -r '.name')

            if [[ "$chart_dir" =~ "cds" ]] && [[ "$include_cds" == false ]]; then
                echo "CDS chart_dir, but skipping because include_cds: $include_cds"
                continue
            fi
            if [[ "$chart_dir" =~ "nnf" ]] && [[ "$include_nnf" == false ]]; then
                echo "NNF chart_dir, but skipping because include_nnf: $include_nnf"
                continue
            fi

            needs_cray_service=$(awk '!/^#/ && /cray-service/ {print}' "$chart_dir"/requirements.yaml)
            if [ "$needs_cray_service" ] && [ ! -d "$chart_dir"/charts ]; then
                mkdir -p "$chart_dir"/charts
                if [ ! -f cray-service-2.0.0.tgz ]; then
                    helm dependency update "$chart_dir" || exit
                    cp "$chart_dir"/charts/cray-service-2.0.0.tgz ./ || exit
                    rm "$chart_dir"/requirements.lock || exit
                else
                    cp 'cray-service-2.0.0.tgz' "$chart_dir"/charts/ || exit
                fi
            fi

            helm install --set cray-service.imagesHost=arti.dev.cray.com "$chart_dir" || exit
        done
    done

}

ready_to_run_initial_checks ()
{
    # Check for the required version of helm
    if ! type helm >/dev/null 2>&1 || ! helm version 2>/dev/null | grep SemVer | grep -q "v2"; then
        echo "This tool requires helm v2"
        exit
    fi

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
        include_nnf )	         include_nnf="true" ;;
        include_cds )            include_cds="true" ;;
        ??* )                    die "Illegal option --$OPT" ;;  # bad long option
        ? )                      exit 2 ;;  # bad short option (error reported via getopts)
    esac
done
shift $((OPTIND-1)) # remove parsed options and args from $@ list


if [ "$display_parameters" = "true" ]; then
    show_parameters
fi

# Verify we are ready to proceed
ready_to_run_initial_checks

action="$1"
case "$action" in
    "$action_create_deployment")
        create_kind_cluster
        deploy
        apply_charts
        ;;
    "$action_recreate_deployment")
        delete_kind_cluster
        create_kind_cluster
        deploy
        apply_charts
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
    "$action_reapply_charts")
        rollback
        apply_charts
        ;;
    *)
        usage
esac