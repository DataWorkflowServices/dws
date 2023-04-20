/*
Copyright 2021.

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

package controllers

import (
	"context"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	zapcr "sigs.k8s.io/controller-runtime/pkg/log/zap"

	dwsv1alpha1 "github.com/HewlettPackard/dws/api/v1alpha1"
	dwsv1alpha2 "github.com/HewlettPackard/dws/api/v1alpha2"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var ctx context.Context
var cancel context.CancelFunc

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {

	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	zaplogger := zapcr.New(zapcr.WriteTo(GinkgoWriter), zapcr.Encoder(encoder), zapcr.UseDevMode(true))
	logf.SetLogger(zaplogger)

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")

	// See https://github.com/kubernetes-sigs/controller-runtime/issues/1882
	// about getting the conversion webhook to register properly.
	// Begin by relocating the code that builds the scheme, so it happens
	// before calling envtest.Start().
	// Then add the scheme to envtest.CRDInstallOptions.

	var err error
	err = dwsv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = dwsv1alpha2.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	testEnv = &envtest.Environment{
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: []string{filepath.Join("..", "config", "webhook")},
		},
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},

		CRDInstallOptions: envtest.CRDInstallOptions{
			// This adds the conversion webhook configuration to
			// the CRDs.
			Scheme: scheme.Scheme,
		},

		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	// start webhook server using Manager
	webhookInstallOptions := &testEnv.WebhookInstallOptions
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:  scheme.Scheme,
		Host:    webhookInstallOptions.LocalServingHost,
		Port:    webhookInstallOptions.LocalServingPort,
		CertDir: webhookInstallOptions.LocalServingCertDir,
	})
	Expect(err).NotTo(HaveOccurred())

	// start webhooks

	err = (&dwsv1alpha2.ClientMount{}).SetupWebhookWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())
	/*
		err = (&dwsv1alpha2.Computes{}).SetupWebhookWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())
	*/

	err = (&dwsv1alpha2.DWDirectiveRule{}).SetupWebhookWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&dwsv1alpha2.DirectiveBreakdown{}).SetupWebhookWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&dwsv1alpha2.PersistentStorageInstance{}).SetupWebhookWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&dwsv1alpha2.Servers{}).SetupWebhookWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&dwsv1alpha2.Storage{}).SetupWebhookWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&dwsv1alpha2.SystemConfiguration{}).SetupWebhookWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&dwsv1alpha2.Workflow{}).SetupWebhookWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	// start reconcilers

	err = (&WorkflowReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Workflow"),
		Scheme: testEnv.Scheme,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err := k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()

})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	cancel()
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
