/*
Copyright 2021-2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"
	"runtime"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	kruntime "k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	dwsv1alpha1 "github.com/DataWorkflowServices/dws/api/v1alpha1"
	dwsv1alpha2 "github.com/DataWorkflowServices/dws/api/v1alpha2"
	dwsv1alpha3 "github.com/DataWorkflowServices/dws/api/v1alpha3"
	dwsv1alpha4 "github.com/DataWorkflowServices/dws/api/v1alpha4"
	dwsv1alpha5 "github.com/DataWorkflowServices/dws/api/v1alpha5"
	"github.com/DataWorkflowServices/dws/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = kruntime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(dwsv1alpha1.AddToScheme(scheme))
	utilruntime.Must(dwsv1alpha2.AddToScheme(scheme))
	utilruntime.Must(dwsv1alpha3.AddToScheme(scheme))
	utilruntime.Must(dwsv1alpha4.AddToScheme(scheme))
	utilruntime.Must(dwsv1alpha5.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var mode string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&mode, "mode", "controller", "What mode to run in (controller, webhook)")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	setupLog.Info("GOMAXPROCS", "value", runtime.GOMAXPROCS(0))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: metricsAddr},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "a08857a2.dataworkflowservices.github.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	switch mode {
	case "controller":
		if err = (&controllers.WorkflowReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("Workflow"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Workflow")
			os.Exit(1)
		}

		if err = (&controllers.SystemConfigurationReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("SystemConfiguration"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "SystemConfiguration")
			os.Exit(1)
		}

		if os.Getenv("ENVIRONMENT") == "kind" {
			if err = (&controllers.ClientMountReconciler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("controllers").WithName("ClientMount"),
				Scheme: mgr.GetScheme(),
			}).SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "Workflow")
				os.Exit(1)
			}
		}
	case "webhook":
		if err = (&dwsv1alpha5.ClientMount{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "ClientMount")
			os.Exit(1)
		}
		if err = (&dwsv1alpha5.Computes{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Computes")
			os.Exit(1)
		}
		if err = (&dwsv1alpha5.DWDirectiveRule{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "DWDirectiveRule")
			os.Exit(1)
		}
		if err = (&dwsv1alpha5.DirectiveBreakdown{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "DirectiveBreakdown")
			os.Exit(1)
		}
		if err = (&dwsv1alpha5.PersistentStorageInstance{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "PersistentStorageInstance")
			os.Exit(1)
		}
		if err = (&dwsv1alpha5.Servers{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Servers")
			os.Exit(1)
		}
		if err = (&dwsv1alpha5.Storage{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Storage")
			os.Exit(1)
		}
		if err = (&dwsv1alpha5.SystemConfiguration{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "SystemConfiguration")
			os.Exit(1)
		}
		if err = (&dwsv1alpha5.Workflow{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Workflow")
			os.Exit(1)
		}
		if err = (&dwsv1alpha5.SystemStatus{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "SystemStatus")
			os.Exit(1)
		}
	default:
		setupLog.Info("unsupported mode", "mode", mode)
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
