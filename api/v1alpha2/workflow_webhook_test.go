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

package v1alpha2

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
				DesiredState: StateProposal,
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
		func(desiredState WorkflowState, expectSuccess bool) {
			workflow.Spec.DesiredState = desiredState
			if expectSuccess {
				Expect(k8sClient.Create(context.TODO(), workflow)).Should(Succeed())
			} else {
				Expect(k8sClient.Create(context.TODO(), workflow)).ShouldNot(Succeed())
			}

			workflow = nil
		},
		Entry("When Spec.DesiredState Proposal", StateProposal, true),
		Entry("When Spec.DesiredState Setup", StateSetup, false),
		Entry("When Spec.DesiredState DataIn", StateDataIn, false),
		Entry("When Spec.DesiredState PreRun", StatePreRun, false),
		Entry("When Spec.DesiredState PostRun", StatePostRun, false),
		Entry("When Spec.DesiredState DataOut", StateDataOut, false),
		Entry("When Spec.DesiredState Teardown", StateTeardown, false),
	)

	DescribeTable("Fails to create workflow with Status.State set",
		func(statusState WorkflowState) {
			workflow.Status.State = statusState
			Expect(k8sClient.Create(context.TODO(), workflow)).ShouldNot(Succeed())
			workflow = nil
		},
		Entry("When Status.State Proposal", StateProposal),
		Entry("When Status.State Setup", StateSetup),
		Entry("When Status.State DataIn", StateDataIn),
		Entry("When Status.State PreRun", StatePreRun),
		Entry("When Status.State PostRun", StatePostRun),
		Entry("When Status.State DataOut", StateDataOut),
		Entry("When Status.State Teardown", StateTeardown),
	)

	Describe("Invalid transitions after create", Ordered, func() {
		BeforeEach(func() {
			Expect(k8sClient.Create(context.TODO(), workflow)).Should(Succeed())
		})

		DescribeTable("Fails to transition out of proposal", func(desiredState WorkflowState) {
			workflow.Spec.DesiredState = desiredState
			Expect(k8sClient.Update(context.TODO(), workflow)).ShouldNot(Succeed())
		},
			Entry("When Spec.DesiredState Setup", StateSetup), // This is a valid state transition, but status is not ready
			Entry("When Spec.DesiredState DataIn", StateDataIn),
			Entry("When Spec.DesiredState PreRun", StatePreRun),
			Entry("When Spec.DesiredState PostRun", StatePostRun),
			Entry("When Spec.DesiredState DataOut", StateDataOut),
			//Entry("When Spec.DesiredState Teardown", StateTeardown), // Transition to Teardown is always permitted
		)

		DescribeTable("Fails to transition out of teardown", func(desiredState WorkflowState) {
			workflow.Spec.DesiredState = StateTeardown
			Expect(k8sClient.Update(context.TODO(), workflow)).Should(Succeed())

			workflow.Spec.DesiredState = desiredState
			Expect(k8sClient.Update(context.TODO(), workflow)).ShouldNot(Succeed())
		},
			Entry("When Spec.DesiredState Setup", StateSetup),
			Entry("When Spec.DesiredState DataIn", StateDataIn),
			Entry("When Spec.DesiredState PreRun", StatePreRun),
			Entry("When Spec.DesiredState PostRun", StatePostRun),
			Entry("When Spec.DesiredState DataOut", StateDataOut),
		)
	})
})
