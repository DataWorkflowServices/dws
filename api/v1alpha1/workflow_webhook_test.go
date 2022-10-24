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
		Expect(k8sClient.Create(context.TODO(), workflow)).ToNot(Succeed())
		workflow = nil
	})
})
