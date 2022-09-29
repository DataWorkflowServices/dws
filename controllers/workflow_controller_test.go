package controllers

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dwsv1alpha1 "github.com/HewlettPackard/dws/api/v1alpha1"
)

var _ = Describe("Workflow Controller Test", func() {
	It("Creates workflow", func() {
		wf := &dwsv1alpha1.Workflow{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: corev1.NamespaceDefault,
			},
			Spec: dwsv1alpha1.WorkflowSpec{
				DesiredState: "proposal",
				WLMID:        "test",
				JobID:        0,
				UserID:       0,
				GroupID:      0,
				DWDirectives: []string{
					"Sup Yo",
				},
			},
		}

		Expect(k8sClient.Create(context.TODO(), wf)).To(Succeed())

		Eventually(func(g Gomega) string {
			g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(wf), wf)).To(Succeed())
			return wf.Status.Status
		}).Should(Equal(dwsv1alpha1.StatusCompleted))

	})
})
