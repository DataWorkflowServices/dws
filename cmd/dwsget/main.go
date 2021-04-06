package main

import (
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	"stash.us.cray.com/dpm/dws-operator/pkg/client/clientset/versioned/typed/dws/v1alpha1"
)

var (
	kuberconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	master      = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	resource    = flag.String("resource", "workflows", "The resource kind to list.")
)

func main() {
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kuberconfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v", err)
		os.Exit(1)
	}

	dwsClient, err := v1alpha1.NewForConfig(cfg)
	if err != nil {
		fmt.Printf("Error building dws clientset: %v", err)
		os.Exit(1)
	}

	switch *resource {
	case "workflows":
		list, err := dwsClient.Workflows("default").List(metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing all workflows: %v", err)
			os.Exit(1)
		}
		for _, r := range list.Items {
			fmt.Printf("workflow %s:\n%v\n", r.Name, r.Spec)
		}
	case "storagepools":
		list, err := dwsClient.StoragePools("default").List(metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing all StoragePools: %v", err)
			os.Exit(1)
		}
		for _, r := range list.Items {
			fmt.Printf("StoragePools %s granularity %s quantity %d free %d\n", r.Name, r.Spec.Granularity, r.Spec.Quantity, r.Spec.Free)
		}
	case "dwddirectiverules":
		list, err := dwsClient.DWDirectiveRules("default").List(metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing all workflows: %v", err)
			os.Exit(1)
		}
		for _, r := range list.Items {
			fmt.Printf("parser %s rules %+v\n", r.Name, r.Spec)
		}
	default:
		fmt.Printf("Unknown resource kind %s\n", *resource)
		os.Exit(1)

	}
}
