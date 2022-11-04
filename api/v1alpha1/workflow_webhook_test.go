/*
 * Copyright 2021, 2022 Hewlett Packard Enterprise Development LP
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

package v1alpha1

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
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
		wfid := uuid.NewString()[0:8]
		workflow = &Workflow{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("w%s", wfid),
				Namespace: metav1.NamespaceDefault,
			},
			Spec: WorkflowSpec{
				DesiredState: StateProposal.String(),
				DWDirectives: []string{},
			},
		}
	})

	AfterEach(func() {
		if workflow != nil {
			Expect(k8sClient.Delete(context.TODO(), workflow)).To(Succeed())
		}
	})

	It("should have workflow environmental variables set successfully", func() {
		Expect(k8sClient.Create(context.TODO(), workflow)).To(Succeed())
		Expect(workflow.Status.Env).To(HaveKeyWithValue("DW_WORKFLOW_NAME", workflow.Name))
		Expect(workflow.Status.Env).To(HaveKeyWithValue("DW_WORKFLOW_NAMESPACE", workflow.Namespace))
	})

	It("Fails to create workflow with hurry flag set", func() {
		workflow.Spec.Hurry = true
		Expect(k8sClient.Create(context.TODO(), workflow)).ShouldNot(Succeed())
		workflow = nil
	})

	DescribeTable("Workflow created only when Spec.DesiredState is Proposal",
		func(statusState string, expectSuccess bool) {
			workflow.Spec.DesiredState = statusState
			if expectSuccess {
				Expect(k8sClient.Create(context.TODO(), workflow)).Should(Succeed())
			} else {
				Expect(k8sClient.Create(context.TODO(), workflow)).ShouldNot(Succeed())
			}

			workflow = nil
		},
		Entry("When Spec.DesiredState Proposal", StateProposal.String(), true),
		Entry("When Spec.DesiredState Setup", StateSetup.String(), false),
		Entry("When Spec.DesiredState DataIn", StateDataIn.String(), false),
		Entry("When Spec.DesiredState PreRun", StatePreRun.String(), false),
		Entry("When Spec.DesiredState PostRun", StatePostRun.String(), false),
		Entry("When Spec.DesiredState DataOut", StateDataOut.String(), false),
		Entry("When Spec.DesiredState Teardown", StateTeardown.String(), false),
	)

	DescribeTable("Fails to create workflow with Status.State set",
		func(statusState string) {
			workflow.Status.State = statusState
			Expect(k8sClient.Create(context.TODO(), workflow)).ShouldNot(Succeed())
			workflow = nil
		},
		Entry("When Status.State Proposal", StateProposal.String()),
		Entry("When Status.State Setup", StateSetup.String()),
		Entry("When Status.State DataIn", StateDataIn.String()),
		Entry("When Status.State PreRun", StatePreRun.String()),
		Entry("When Status.State PostRun", StatePostRun.String()),
		Entry("When Status.State DataOut", StateDataOut.String()),
		Entry("When Status.State Teardown", StateTeardown.String()),
	)
})
