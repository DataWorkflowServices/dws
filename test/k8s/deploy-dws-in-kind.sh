#!/bin/bash

#set -vx

action_create_deployment="create_deployment"
action_create_kind_cluster="create_kind_cluster"
action_delete_kind_cluster="delete_kind_cluster"
action_deploy="deploy"
action_rollback="rollback"

dws_tools_dir=$(pwd)/$(dirname "${0}")
docker_image_name_prefix="arti.dev.cray.com/cray/cray-"

workdir=${workdir:-'./deploy-in-kind-work'}
helm_dir=${helm_dir:-${workdir}}
kind_config=${kind_config:-kind-1-2.yaml}

# Manifests of images.
# the 'kj' product.
main_manifest=${main_manifest:-'main-manifest.txt'}

required_images=('dws-operator'
	'dws-admission-webhook'
	'dws-cds-driver')

# Version/tag of images that helm charts expects
deploy_tag=${deploy_tag:-'latest'}

# Used to retrieve charts
default_sdp_manifests=('http://car.dev.cray.com/artifactory/kj/DPM/noos/noarch/dev/master/ssm-team/sdp_manifest/data-workflow-services_1.0.0_metadata.json'
	                   'http://car.dev.cray.com/artifactory/kj/DPM/noos/noarch/dev/master/ssm-team/sdp_manifest/dws-cds-integration_1.0.0_metadata.json')

sdp_manifests=("${sdp_manifests:-${default_sdp_manifests[@]}}")

# Don't deploy these charts
chart_filter=${chart_filter:-'.*-deployment-group-manager'}

usage ()
{
	echo "usage: ${0} action"
	echo ""
	echo "  Actions"
	echo "    ${action_create_deployment}     Creates a new deployment using the latest from DWS"
	echo "    ${action_create_kind_cluster}   Creates a Kind cluster without DWS deployments"
	echo "    ${action_delete_kind_cluster}   Deletes the cluster"
	echo "    ${action_deploy}                Deploys DWS"
	echo "    ${action_rollback}              Removes DWS deployments"
	echo ""
	echo "  Options:"
	echo "    -c  Find images in manifest files and remove them from the docker cache."
	echo "        Used in conjunction with \"${main_manifest}\"."
	echo ""
	echo "  Parameters"
	echo "    main_manifest     Location of the main manifest file.  Defaults to: \"${main_manifest}\""
	echo "    sdp_manifests     Contains helm metadata"
	echo "    workdir           Working directory to hold things like manifests, repos, and image artifacts.  Defaults to: \"${workdir}\""
	echo "    deploy_tag        Version/tag of images that helm charts expects.  Defaults to: \"${deploy_tag}\""
	echo ""
	echo "  Examples"
	echo "    To deploy DWS in Kind using the latest images from artifactory, run:"
	echo "      ${0} ${action_create_deployment}"
	echo ""
	echo "    To deploy with a different artifactory image.  Update the image version in ${main_manifest}"
	echo "    or ${main_manifest}, and run:"
	echo "      ${0} ${action_deploy}"
	echo ""
	echo "      Note:  To deploy a custom image, load the image into your local docker cache then update the manifest"
	echo "             with the docker name and tag of the custom image.  Finally, run the above command."
	echo ""
	echo "    To start from scratch, remove the working directory and appropriate images from docker cache.  Then create a new deployment:"
	echo "      ${0} ${action_rollback}"
	echo "      rm -rf ${workdir}"
	echo "      ${0} ${action_create_deployment}"
	echo ""
	echo ""
}

verify_k8s_is_running ()
{
	if kind get clusters 2>&1 >/dev/null | grep -q 'No kind clusters found'; then
		echo "No kind clusters found.  If you would like to create a cluster.  Run this tool using '${action_create_kind_cluster}' or '${action_create_deployment}'"
		exit
	fi
}
delete_kind_cluster ()
{
	kind delete cluster || exit
}

create_kind_cluster ()
{
	# Create a config file used to create the Kind cluster,
	# if it doesn't already exist
	if [ ! -f "${kind_config}" ]; then
		cat << EOF > "${kind_config}"
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

	delete_kind_cluster

	kind create cluster --config="${kind_config}"		# see embedded config file above
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
	
	echo "Kind cluster is ready"
}

rollback ()
{
	verify_k8s_is_running

	ordered_rollback_list=('dws-operator'
		'dws-validation-webhook'
		'dws-cds-webhook'
		'dws-cds-integration')

	for chart in "${ordered_rollback_list[@]}";
	do
		helm list | grep "${chart}" | awk '{print $1}' | xargs helm delete || exit
	done

	resource_kinds=$(kubectl get crds | grep -E 'dws.cray.hpe.com' | awk '{print $1}' | sed -e 's/.dws.cray.hpe.com//')
	for res_type in ${resource_kinds}
	do
		if ! kubectl get "${res_type}" -n default > /dev/null 2>&1
		then
			continue
		fi
		if cnt=$(kubectl get "${res_type}" -n default -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' | wc -l)
		then
			if (( cnt > 0 ))
			then
				echo "Force-delete ${res_type}"
				if ! kubectl get "${res_type}" --no-headers | awk '{print $1}' | xargs kubectl patch "${res_type}" --type merge -p '{"metadata":{"finalizers":null}}'
				then
					echo "Failed to delete ${res_type}"
				fi
			fi
		else
			echo "Unable to count ${res_type}" > /dev/stderr
		fi
	done

	if [ "${clean_docker_cache}" = "true" ]; then
		for image_base_name in "${required_images[@]}";
		do
			basic_name=${image_base_name}

			# Get the name and version from the manifest
			docker_image_name=$(grep -h "${basic_name}" "${main_manifest}" | awk '{print $1}' | sed 's|.$||') || exit
			ver=$(grep -h "${basic_name}" "${main_manifest}" | awk '{print $2}') || exit
			if [ -z "${docker_image_name}" ] || [ -z "${ver}" ]; then
				echo "Failed to retrieve info for ${basic_name} (${docker_image_name}:${ver})"
				continue
			fi

			docker rmi "${docker_image_name}:${ver}"
			docker rmi "${docker_image_name_prefix}${image_base_name}:${deploy_tag}"
		done
	fi
}

deploy ()
{
	verify_k8s_is_running

	if [ ! -f "${main_manifest}" ]; then
		# Fetch the manifest from artifactory and update each image in the file to reflect its name (i.e. "Repository") after loaded into docker
		wget https://arti.dev.cray.com/artifactory/kj-misc-master-local/manifest/manifest.txt -O "tmp-${main_manifest}" || exit
		awk '{print "sms.local:5000/cray/" $0}' "tmp-${main_manifest}" > "${main_manifest}"
		rm "tmp-${main_manifest}"
	fi

	DWS_NODES=$(kubectl get node -l cray.dpm.dg.data-workflow-services=true -o json | jq -r '.items[].metadata.name' | paste -s -d, -)
	echo "DWS_NODES: ${DWS_NODES}"

	for image_base_name in "${required_images[@]}";
	do
		basic_name=${image_base_name}

		# Get the name and version from the manifest
		docker_image_name=$(grep -h "${basic_name}" "${main_manifest}" | awk '{print $1}' | sed 's|.$||') || exit
		ver=$(grep -h "${basic_name}" "${main_manifest}" | awk '{print $2}') || exit

		if [ -z "${docker_image_name}" ] || [ -z "${ver}" ]; then
			echo "Failed to retrieve info for ${basic_name} (${docker_image_name}:${ver})"
			exit
		fi

		# Check to see if it's already been loaded into the docker cache
		if [ -z "$(docker images -q "${docker_image_name}:${ver}" 2> /dev/null)" ]; then
			# Check to see if the image has been downloaded
			artifactory_name="cray-${image_base_name}-${ver}-dockerimage.tar"
			if [ ! -f "${artifactory_name}" ]; then
				# Try to fetch the image from artifactory where the DWS images are kept
				if ! wget "http://car.dev.cray.com/artifactory/kj/DPM/noos/noarch/dev/master/ssm-team/${artifactory_name}"; then
					echo "Unable to find image: ${artifactory_name}"
					exit
				fi
			fi
			# Load the image into the docker cache
			echo "Docker load $artifactory_name"
			docker load -i "${artifactory_name}" || exit
		fi

		deployment_nodes="${DWS_NODES}"

		# Tag the image according to what the helm charts will look for, and load it onto the Kind nodes
		echo "Docker tag ${docker_image_name}:${ver} ${docker_image_name_prefix}${image_base_name}:${deploy_tag}"
		docker tag "${docker_image_name}:${ver}" "${docker_image_name_prefix}${image_base_name}:${deploy_tag}" || exit
		echo "Kind load --nodes ${deployment_nodes} ${docker_image_name_prefix}${image_base_name}:${deploy_tag}"
		kind load docker-image --nodes "${deployment_nodes}" "${docker_image_name_prefix}${image_base_name}:${deploy_tag}" || exit
	done

set -x

	if ! helm repo list | grep -q cray-internal; then
		helm repo add cray-internal http://helmrepo.dev.cray.com:8080 || exit
	fi

	# Get the helm charts and install them
	for sdp_manifest in "${sdp_manifests[@]}";
	do
		# shellcheck disable=SC2207
		charts_metadata=($(curl "${sdp_manifest}" | jq | grep 'metadata.json' | grep -vE "${chart_filter}" | sed 's|,||' | sed 's|"||g' | tr -d ' '))

		for helm_metadata in "${charts_metadata[@]}";
		do
			# Get the repo that contains the helm chart
			repo=$(curl "${helm_metadata}" | jq -r '.git_clone_url')

			# Clone the repo if it doesn't exist locally
			repo_dir=$(basename "${repo}" | sed 's|.git||')
			if [ ! -d "${repo_dir}" ]; then
				git clone --depth 1 "${repo}" || exit
			fi

			# Get/Update charts if needed and install
			chart_dir="${repo_dir}"/kubernetes/$(curl "${helm_metadata}" | jq -r '.name')
			needs_cray_service=$(cat ${chart_dir}/requirements.yaml | awk '!/^#/ && /cray-service/ {print}')
			if [ "$needs_cray_service" -a ! -d "${chart_dir}"/charts ]; then
				mkdir -p "${chart_dir}"/charts
				if [ ! -f cray-service-2.0.0.tgz ]; then
					helm dependency update "${chart_dir}" || exit
					cp "${chart_dir}"/charts/cray-service-2.0.0.tgz ./ || exit
					rm "${chart_dir}"/requirements.lock || exit
				else
					cp 'cray-service-2.0.0.tgz' "${chart_dir}"/charts/ || exit
				fi
			fi

			helm install --set cray-service.imagesHost=arti.dev.cray.com "${chart_dir}" || exit
		done
	done
}


# Check for the required version of helm
if ! type helm 2>/dev/null || ! helm version 2> /dev/null | grep SemVer | grep -q "v2"; then
	echo "This tool requires helm v2"
	exit
fi

# Check other required dependencies
type wget 2>/dev/null || { echo "wget is required but not found.  Exiting" ; exit; }
type jq   2>/dev/null || { echo "jq is required but not found.  Exiting" ; exit; }

mkdir -p "${workdir}" || exit
cd "${workdir}" || exit

while getopts "c" opt; do
	case "${opt}" in
		c)
			clean_docker_cache="true"
			;;
		*)
			usage
	esac
done
shift $((OPTIND - 1))

action="${1}"
case "${action}" in
	"${action_create_deployment}")
		create_kind_cluster
		deploy
		;;
	"${action_create_kind_cluster}")
		create_kind_cluster
		;;
	"${action_delete_kind_cluster}")
		delete_kind_cluster
		;;
	"${action_deploy}")
		rollback
		deploy
		;;
	"${action_rollback}")
		rollback
		;;
	*)
		usage
esac
