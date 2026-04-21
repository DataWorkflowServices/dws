/*
 * Copyright 2025 Hewlett Packard Enterprise Development LP
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

package v1alpha7

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("PersistentStorageInstance Webhook", func() {

	var (
		psi *PersistentStorageInstance
	)

	BeforeEach(func() {
		psiid := uuid.NewString()[0:8]
		psi = &PersistentStorageInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("psi-%s", psiid),
				Namespace: metav1.NamespaceDefault,
			},
			Spec: PersistentStorageInstanceSpec{
				Name:        fmt.Sprintf("psi-%s", psiid),
				FsType:      "lustre",
				DWDirective: "#DW create_persistent name=psi type=lustre capacity=1GiB",
				UserID:      1000,
				State:       PSIStateActive,
			},
		}
	})

	AfterEach(func() {
		if psi != nil {
			Expect(k8sClient.Delete(context.TODO(), psi)).To(Succeed())
		}
	})

	It("should create a PersistentStorageInstance with valid spec", func() {
		Expect(k8sClient.Create(context.TODO(), psi)).To(Succeed())
	})

	It("should fail to create when spec.state is not Active", func() {
		psi.Spec.State = PSIStateDestroying
		Expect(k8sClient.Create(context.TODO(), psi)).ShouldNot(Succeed())
		psi = nil
	})

	Describe("Immutable fields after create", Ordered, func() {
		BeforeEach(func() {
			Expect(k8sClient.Create(context.TODO(), psi)).Should(Succeed())
		})

		It("should reject changes to FsType", func() {
			psi.Spec.FsType = "xfs"
			Expect(k8sClient.Update(context.TODO(), psi)).ShouldNot(Succeed())
		})

		It("should reject changes to DWDirective", func() {
			psi.Spec.DWDirective = "#DW create_persistent name=psi type=lustre capacity=10GiB"
			Expect(k8sClient.Update(context.TODO(), psi)).ShouldNot(Succeed())
		})

		It("should reject changes to Name", func() {
			psi.Spec.Name = "different-name"
			Expect(k8sClient.Update(context.TODO(), psi)).ShouldNot(Succeed())
		})

		It("should reject changes to UserID", func() {
			psi.Spec.UserID = 9999
			Expect(k8sClient.Update(context.TODO(), psi)).ShouldNot(Succeed())
		})

		It("should allow transitioning State to Destroying", func() {
			psi.Spec.State = PSIStateDestroying
			Expect(k8sClient.Update(context.TODO(), psi)).Should(Succeed())
		})
	})
})
