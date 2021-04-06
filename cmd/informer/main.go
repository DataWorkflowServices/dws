package main

import (
	"flag"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"stash.us.cray.com/dpm/dws-operator/pkg/apis/dws/v1alpha1"
	clientset "stash.us.cray.com/dpm/dws-operator/pkg/client/clientset/versioned"
	informers "stash.us.cray.com/dpm/dws-operator/pkg/client/informers/externalversions"
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

	dwsClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		fmt.Printf("Error building dws clientset: %v", err)
		os.Exit(1)
	}

	factory := informers.NewSharedInformerFactory(dwsClient, 0)
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	var informer cache.SharedIndexInformer

	switch *resource {
	case "workflows":
		informer = factory.Dws().V1alpha1().Workflows().Informer()
		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc:    onWorkflowAdd,
			UpdateFunc: onWorkflowUpdate,
		})
		//	case "storagepools":
		//		informer = factory.Core().V1alpha1().StoragPools()Informer()
		//		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		//			AddFunc: onStoragPoolAdd,
		//			UpdateFunc: onStoragPoolUpdate,
		//		})
		//	case "dwddirectiverules":
		//		informer = factory.Core().V1alpha1().DWDirectiveRules().Informer()
		//		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		//			AddFunc: onDWDirectiveRuleAdd,
		//			UpdateFunc: onDWDirectiveRuleUpdate,
		//		})
	default:
		fmt.Printf("Unknown resource kind %s\n", *resource)
		os.Exit(1)
	}

	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	<-stopper
}

func onWorkflowAdd(obj interface{}) {
	wf := obj.(*v1alpha1.Workflow)
	fmt.Printf("Workflow:\n%+v", wf)
}

func onWorkflowUpdate(old interface{}, new interface{}) {
	wfo := old.(*v1alpha1.Workflow)
	wfn := new.(*v1alpha1.Workflow)
	fmt.Printf("Workflow Old:\n%+v", wfo)
	fmt.Printf("Workflow New:\n%+v", wfn)
}

func onStoragePoolAdd(obj interface{}) {
	sp := obj.(*v1alpha1.StoragePool)
	fmt.Printf("StoragePool:\n%+v", sp)
}

func onDWDirectiveRuleAdd(obj interface{}) {
	dr := obj.(*v1alpha1.DWDirectiveRule)
	fmt.Printf("DWDirectiveRule:\n%+v", dr)
}
