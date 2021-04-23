package main

import (
	"flag"
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"stash.us.cray.com/dpm/dws-operator/pkg/apis"
	"stash.us.cray.com/dpm/dws-operator/pkg/apis/dws/v1alpha1"
)

var kubeconfig string
var resource string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.StringVar(&resource, "resource", "workflows", "The resource kind to list.")
	flag.Parse()
}

func main() {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		fmt.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		fmt.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		panic(err)
	}

	apis.AddToScheme(scheme.Scheme)

	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &v1alpha1.SchemeGroupVersion
	crdConfig.APIPath = "/apis"
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	dwsClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		panic(err)
	}

	switch resource {
	case "workflows":
		var list = v1alpha1.WorkflowList{}
		err = dwsClient.Get().Resource(resource).Do().Into(&list)
		if err != nil {
			fmt.Printf("Error listing all workflows: %v", err)
			os.Exit(1)
		}
		for _, r := range list.Items {
			fmt.Printf("workflow %s\n%v\nStatus: %v\n", r.Name, r.Spec, r.Status)
		}
	case "storagepools":
		var list = v1alpha1.StoragePoolList{}
		err = dwsClient.Get().Resource(resource).Do().Into(&list)
		if err != nil {
			fmt.Printf("Error listing all StoragePools: %v", err)
			os.Exit(1)
		}
		for _, r := range list.Items {
			fmt.Printf("StoragePool %s granularity %s quantity %d free %d\n", r.Name, r.Spec.Granularity, r.Spec.Quantity, r.Spec.Free)
		}
	case "dwddirectiverules":
		var list = v1alpha1.DWDirectiveRuleList{}
		err = dwsClient.Get().Resource(resource).Do().Into(&list)
		if err != nil {
			fmt.Printf("Error listing all DWDirectiveRules: %v", err)
			os.Exit(1)
		}
		for _, r := range list.Items {
			fmt.Printf("parser %s rules %+v\n", r.Name, r.Spec)
		}
	default:
		fmt.Printf("Unknown resource kind %s.\n", resource)
		os.Exit(1)

	}
}
