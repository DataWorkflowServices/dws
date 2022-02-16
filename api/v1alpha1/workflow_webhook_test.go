/*
Copyright 2021 Hewlett Packard Enterprise Development LP
*/

package v1alpha1

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// These tests are written in BDD-style using Ginkgo framework. Refer to
// http://onsi.github.io/ginkgo to learn more.

var _ = Describe("Workflow Webhook", func() {
	var (
		workflow *Workflow
	)

	BeforeEach(func() {
		workflow = &Workflow{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: metav1.NamespaceDefault,
			},
			Spec: WorkflowSpec{
				DesiredState: StateProposal.String(),
				DWDirectives: []string{},
			},
		}

		Expect(k8sClient.Create(context.TODO(), workflow)).To(Succeed())
	})

	AfterEach(func() {
		Expect(k8sClient.Delete(context.TODO(), workflow)).To(Succeed())
	})

	It("should have workflow environmental variables set successfully", func() {
		Expect(workflow.Status.Env).To(HaveKeyWithValue("DW_WORKFLOW_NAME", workflow.Name))
		Expect(workflow.Status.Env).To(HaveKeyWithValue("DW_WORKFLOW_NAMESPACE", workflow.Namespace))
	})

})
