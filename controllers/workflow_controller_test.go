/*
 * Copyright 2022-2023 Hewlett Packard Enterprise Development LP
 * Other additional copyright holders may be indicated within.
 *
 * The entirety of this work is licensed under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 *
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controllers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dwsv1alpha2 "github.com/HewlettPackard/dws/api/v1alpha2"
)

var _ = Describe("Workflow Controller Test", func() {

	var (
		wf *dwsv1alpha2.Workflow
	)

	BeforeEach(func() {
		wfid := uuid.NewString()[0:8]
		wf = &dwsv1alpha2.Workflow{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s", wfid),
				Namespace: corev1.NamespaceDefault,
			},
			Spec: dwsv1alpha2.WorkflowSpec{
				DesiredState: dwsv1alpha2.StateProposal,
				WLMID:        "test",
				JobID:        intstr.FromString("wlm job 442"),
				UserID:       0,
				GroupID:      0,
				DWDirectives: []string{},
			},
		}
	})

	AfterEach(func() {
		if wf != nil {
			Expect(k8sClient.Delete(context.TODO(), wf)).To(Succeed())

			wfExpected := &dwsv1alpha2.Workflow{}
			Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present.
				return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(wf), wfExpected)
			}).ShouldNot(Succeed())
		}
	})

	It("Creates workflow", func() {
		Expect(k8sClient.Create(context.TODO(), wf)).To(Succeed())

		Eventually(func(g Gomega) string {
			g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(wf), wf)).To(Succeed())
			return wf.Status.Status
		}).Should(Equal(dwsv1alpha2.StatusCompleted))

	})

	It("Fails to create workflow with hurry flag set", func() {
		wf.Spec.Hurry = true
		Expect(k8sClient.Create(context.TODO(), wf)).ToNot(Succeed())
		wf = nil
	})

	It("Creates workflow, then fails to add hurry flag", func() {
		Expect(k8sClient.Create(context.TODO(), wf)).To(Succeed())

		Eventually(func(g Gomega) string {
			g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(wf), wf)).To(Succeed())
			return wf.Status.Status
		}).Should(Equal(dwsv1alpha2.StatusCompleted))

		wf.Spec.Hurry = true
		Expect(k8sClient.Update(context.TODO(), wf)).ToNot(Succeed())
	})

	It("Creates workflow, goes to teardown with hurry flag", func() {
		Expect(k8sClient.Create(context.TODO(), wf)).To(Succeed())

		Eventually(func(g Gomega) string {
			g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(wf), wf)).To(Succeed())
			return wf.Status.Status
		}).Should(Equal(dwsv1alpha2.StatusCompleted))

		wf.Spec.DesiredState = dwsv1alpha2.StateTeardown
		wf.Spec.Hurry = true
		Expect(k8sClient.Update(context.TODO(), wf)).To(Succeed())
	})
})
