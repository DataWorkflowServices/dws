/*
 * Copyright 2023-2025 Hewlett Packard Enterprise Development LP
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

package controller

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dwsv1alpha6 "github.com/DataWorkflowServices/dws/api/v1alpha6"
	dwsv1alpha7 "github.com/DataWorkflowServices/dws/api/v1alpha7"
	utilconversion "github.com/DataWorkflowServices/dws/github/cluster-api/util/conversion"
	"github.com/DataWorkflowServices/dws/utils/dwdparse"
)

var _ = Describe("Conversion Webhook Test", func() {

	// Don't get deep into verifying the conversion.
	// We have api/<spoke_ver>/conversion_test.go that is digging deep.
	// We're just verifying that the conversion webhook is hooked up.

	// Note: if a resource is accessed by its spoke API, then it should
	// have the utilconversion.DataAnnotation annotation.  It will not
	// have that annotation when it is accessed by its hub API.

	Context("ClientMount", func() {
		var resHub *dwsv1alpha7.ClientMount

		BeforeEach(func() {
			id := uuid.NewString()[0:8]
			resHub = &dwsv1alpha7.ClientMount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      id,
					Namespace: corev1.NamespaceDefault,
				},
				Spec: dwsv1alpha7.ClientMountSpec{
					Node:         "client-01",
					DesiredState: "unmounted",
					Mounts: []dwsv1alpha7.ClientMountInfo{
						{
							MountPath:      "",
							SetPermissions: false,
							Options:        "",
							Device: dwsv1alpha7.ClientMountDevice{
								Type: "reference",
							},
							Type:       "none",
							TargetType: "directory",
						},
					},
				},
			}

			Expect(k8sClient.Create(context.TODO(), resHub)).To(Succeed())
		})

		AfterEach(func() {
			if resHub != nil {
				Expect(k8sClient.Delete(context.TODO(), resHub)).To(Succeed())
				expected := &dwsv1alpha7.ClientMount{}
				Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present.
					return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), expected)
				}).ShouldNot(Succeed())
			}
		})

		It("reads ClientMount resource via hub and via spoke v1alpha6", func() {
			// Spoke should have annotation.
			resSpoke := &dwsv1alpha6.ClientMount{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resSpoke)).To(Succeed())
				anno := resSpoke.GetAnnotations()
				g.Expect(anno).To(HaveLen(1))
				g.Expect(anno).Should(HaveKey(utilconversion.DataAnnotation))
			}).Should(Succeed())

			// Hub should not have annotation.
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resHub)).To(Succeed())
				anno := resHub.GetAnnotations()
				g.Expect(anno).To(HaveLen(0))
			}).Should(Succeed())
		})

		// +crdbumper:scaffold:spoketest="dataworkflowservices.ClientMount"
	})

	Context("DWDirectiveRule", func() {
		var resHub *dwsv1alpha7.DWDirectiveRule

		BeforeEach(func() {
			id := uuid.NewString()[0:8]
			resHub = &dwsv1alpha7.DWDirectiveRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      id,
					Namespace: corev1.NamespaceDefault,
				},
				Spec: []dwdparse.DWDirectiveRuleSpec{},
			}

			Expect(k8sClient.Create(context.TODO(), resHub)).To(Succeed())
		})

		AfterEach(func() {
			if resHub != nil {
				Expect(k8sClient.Delete(context.TODO(), resHub)).To(Succeed())
				expected := &dwsv1alpha7.DWDirectiveRule{}
				Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present.
					return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), expected)
				}).ShouldNot(Succeed())
			}
		})

		It("reads DWDirectiveRule resource via hub and via spoke v1alpha6", func() {
			// Spoke should have annotation.
			resSpoke := &dwsv1alpha6.DWDirectiveRule{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resSpoke)).To(Succeed())
				anno := resSpoke.GetAnnotations()
				g.Expect(anno).To(HaveLen(1))
				g.Expect(anno).Should(HaveKey(utilconversion.DataAnnotation))
			}).Should(Succeed())

			// Hub should not have annotation.
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resHub)).To(Succeed())
				anno := resHub.GetAnnotations()
				g.Expect(anno).To(HaveLen(0))
			}).Should(Succeed())
		})

		// +crdbumper:scaffold:spoketest="dataworkflowservices.DWDirectiveRule"
	})

	Context("DirectiveBreakdown", func() {
		var resHub *dwsv1alpha7.DirectiveBreakdown

		BeforeEach(func() {
			id := uuid.NewString()[0:8]
			resHub = &dwsv1alpha7.DirectiveBreakdown{
				ObjectMeta: metav1.ObjectMeta{
					Name:      id,
					Namespace: corev1.NamespaceDefault,
				},
				Spec: dwsv1alpha7.DirectiveBreakdownSpec{},
			}

			Expect(k8sClient.Create(context.TODO(), resHub)).To(Succeed())
		})

		AfterEach(func() {
			if resHub != nil {
				Expect(k8sClient.Delete(context.TODO(), resHub)).To(Succeed())
				expected := &dwsv1alpha7.DirectiveBreakdown{}
				Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present.
					return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), expected)
				}).ShouldNot(Succeed())
			}
		})

		It("reads DirectiveBreakdown resource via hub and via spoke v1alpha6", func() {
			// Spoke should have annotation.
			resSpoke := &dwsv1alpha6.DirectiveBreakdown{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resSpoke)).To(Succeed())
				anno := resSpoke.GetAnnotations()
				g.Expect(anno).To(HaveLen(1))
				g.Expect(anno).Should(HaveKey(utilconversion.DataAnnotation))
			}).Should(Succeed())

			// Hub should not have annotation.
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resHub)).To(Succeed())
				anno := resHub.GetAnnotations()
				g.Expect(anno).To(HaveLen(0))
			}).Should(Succeed())
		})

		// +crdbumper:scaffold:spoketest="dataworkflowservices.DirectiveBreakdown"
	})

	Context("PersistentStorageInstance", func() {
		var resHub *dwsv1alpha7.PersistentStorageInstance

		BeforeEach(func() {
			id := uuid.NewString()[0:8]
			resHub = &dwsv1alpha7.PersistentStorageInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      id,
					Namespace: corev1.NamespaceDefault,
				},
				Spec: dwsv1alpha7.PersistentStorageInstanceSpec{
					Name:        "",
					FsType:      "raw",
					DWDirective: "",
					UserID:      0,
					State:       "Active",
				},
			}

			Expect(k8sClient.Create(context.TODO(), resHub)).To(Succeed())
		})

		AfterEach(func() {
			if resHub != nil {
				Expect(k8sClient.Delete(context.TODO(), resHub)).To(Succeed())
				expected := &dwsv1alpha7.PersistentStorageInstance{}
				Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present.
					return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), expected)
				}).ShouldNot(Succeed())
			}
		})

		It("reads PersistentStorageInstance resource via hub and via spoke v1alpha6", func() {
			// Spoke should have annotation.
			resSpoke := &dwsv1alpha6.PersistentStorageInstance{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resSpoke)).To(Succeed())
				anno := resSpoke.GetAnnotations()
				g.Expect(anno).To(HaveLen(1))
				g.Expect(anno).Should(HaveKey(utilconversion.DataAnnotation))
			}).Should(Succeed())

			// Hub should not have annotation.
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resHub)).To(Succeed())
				anno := resHub.GetAnnotations()
				g.Expect(anno).To(HaveLen(0))
			}).Should(Succeed())
		})

		// +crdbumper:scaffold:spoketest="dataworkflowservices.PersistentStorageInstance"
	})

	Context("Servers", func() {
		var resHub *dwsv1alpha7.Servers

		BeforeEach(func() {
			id := uuid.NewString()[0:8]
			resHub = &dwsv1alpha7.Servers{
				ObjectMeta: metav1.ObjectMeta{
					Name:      id,
					Namespace: corev1.NamespaceDefault,
				},
				Spec: dwsv1alpha7.ServersSpec{},
			}

			Expect(k8sClient.Create(context.TODO(), resHub)).To(Succeed())
		})

		AfterEach(func() {
			if resHub != nil {
				Expect(k8sClient.Delete(context.TODO(), resHub)).To(Succeed())
				expected := &dwsv1alpha7.Servers{}
				Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present.
					return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), expected)
				}).ShouldNot(Succeed())
			}
		})

		It("reads Servers resource via hub and via spoke v1alpha6", func() {
			// Spoke should have annotation.
			resSpoke := &dwsv1alpha6.Servers{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resSpoke)).To(Succeed())
				anno := resSpoke.GetAnnotations()
				g.Expect(anno).To(HaveLen(1))
				g.Expect(anno).Should(HaveKey(utilconversion.DataAnnotation))
			}).Should(Succeed())

			// Hub should not have annotation.
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resHub)).To(Succeed())
				anno := resHub.GetAnnotations()
				g.Expect(anno).To(HaveLen(0))
			}).Should(Succeed())
		})

		// +crdbumper:scaffold:spoketest="dataworkflowservices.Servers"
	})

	Context("Storage", func() {
		var resHub *dwsv1alpha7.Storage

		BeforeEach(func() {
			id := uuid.NewString()[0:8]
			resHub = &dwsv1alpha7.Storage{
				ObjectMeta: metav1.ObjectMeta{
					Name:      id,
					Namespace: corev1.NamespaceDefault,
				},
				Spec: dwsv1alpha7.StorageSpec{},
			}

			Expect(k8sClient.Create(context.TODO(), resHub)).To(Succeed())
		})

		AfterEach(func() {
			if resHub != nil {
				Expect(k8sClient.Delete(context.TODO(), resHub)).To(Succeed())
				expected := &dwsv1alpha7.Storage{}
				Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present.
					return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), expected)
				}).ShouldNot(Succeed())
			}
		})

		It("reads Storage resource via hub and via spoke v1alpha6", func() {
			// Spoke should have annotation.
			resSpoke := &dwsv1alpha6.Storage{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resSpoke)).To(Succeed())
				anno := resSpoke.GetAnnotations()
				g.Expect(anno).To(HaveLen(1))
				g.Expect(anno).Should(HaveKey(utilconversion.DataAnnotation))
			}).Should(Succeed())

			// Hub should not have annotation.
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resHub)).To(Succeed())
				anno := resHub.GetAnnotations()
				g.Expect(anno).To(HaveLen(0))
			}).Should(Succeed())
		})

		// +crdbumper:scaffold:spoketest="dataworkflowservices.Storage"
	})

	Context("SystemConfiguration", func() {
		var resHub *dwsv1alpha7.SystemConfiguration

		BeforeEach(func() {
			id := uuid.NewString()[0:8]
			resHub = &dwsv1alpha7.SystemConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name:      id,
					Namespace: corev1.NamespaceDefault,
				},
				Spec: dwsv1alpha7.SystemConfigurationSpec{},
			}

			Expect(k8sClient.Create(context.TODO(), resHub)).To(Succeed())
		})

		AfterEach(func() {
			if resHub != nil {
				Expect(k8sClient.Delete(context.TODO(), resHub)).To(Succeed())
				expected := &dwsv1alpha7.SystemConfiguration{}
				Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present.
					return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), expected)
				}).ShouldNot(Succeed())
			}
		})

		It("reads SystemConfiguration resource via hub and via spoke v1alpha6", func() {
			// Spoke should have annotation.
			resSpoke := &dwsv1alpha6.SystemConfiguration{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resSpoke)).To(Succeed())
				anno := resSpoke.GetAnnotations()
				g.Expect(anno).To(HaveLen(1))
				g.Expect(anno).Should(HaveKey(utilconversion.DataAnnotation))
			}).Should(Succeed())

			// Hub should not have annotation.
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resHub)).To(Succeed())
				anno := resHub.GetAnnotations()
				g.Expect(anno).To(HaveLen(0))
			}).Should(Succeed())
		})

		// +crdbumper:scaffold:spoketest="dataworkflowservices.SystemConfiguration"
	})

	Context("Workflow", func() {
		var resHub *dwsv1alpha7.Workflow

		BeforeEach(func() {
			id := uuid.NewString()[0:8]
			resHub = &dwsv1alpha7.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Name:      id,
					Namespace: corev1.NamespaceDefault,
				},
				Spec: dwsv1alpha7.WorkflowSpec{
					DesiredState: dwsv1alpha7.StateProposal,
					WLMID:        "test",
					JobID:        intstr.FromString("a job id 42"),
					UserID:       0,
					GroupID:      0,
					DWDirectives: []string{},
				},
			}

			// The workflow_controller's Reconcile() will also
			// create a Computes resource to go with this.
			Expect(k8sClient.Create(context.TODO(), resHub)).To(Succeed())
		})

		AfterEach(func() {
			if resHub != nil {
				Expect(k8sClient.Delete(context.TODO(), resHub)).To(Succeed())
				expected := &dwsv1alpha7.Workflow{}
				Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present.
					return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), expected)
				}).ShouldNot(Succeed())
			}
		})

		It("reads Workflow resource via hub and via spoke v1alpha6", func() {
			// Spoke should have annotation.
			resSpoke := &dwsv1alpha6.Workflow{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resSpoke)).To(Succeed())
				anno := resSpoke.GetAnnotations()
				g.Expect(anno).To(HaveLen(1))
				g.Expect(anno).Should(HaveKey(utilconversion.DataAnnotation))
			}).Should(Succeed())

			// Hub should not have annotation.
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resHub)).To(Succeed())
				anno := resHub.GetAnnotations()
				g.Expect(anno).To(HaveLen(0))
			}).Should(Succeed())
		})

		// +crdbumper:scaffold:spoketest="dataworkflowservices.Workflow"
	})

	Context("Computes", func() {
		var resHub *dwsv1alpha7.Computes

		BeforeEach(func() {
			id := uuid.NewString()[0:8]
			resHub = &dwsv1alpha7.Computes{
				ObjectMeta: metav1.ObjectMeta{
					Name:      id,
					Namespace: corev1.NamespaceDefault,
				},
				//Spec: dwsv1alpha7.ComputesSpec{},
			}

			Expect(k8sClient.Create(context.TODO(), resHub)).To(Succeed())
		})

		AfterEach(func() {
			if resHub != nil {
				Expect(k8sClient.Delete(context.TODO(), resHub)).To(Succeed())
				expected := &dwsv1alpha7.Computes{}
				Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present.
					return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), expected)
				}).ShouldNot(Succeed())
			}
		})

		It("reads Computes resource via hub and via spoke v1alpha6", func() {
			// Spoke should have annotation.
			resSpoke := &dwsv1alpha6.Computes{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resSpoke)).To(Succeed())
				anno := resSpoke.GetAnnotations()
				g.Expect(anno).To(HaveLen(1))
				g.Expect(anno).Should(HaveKey(utilconversion.DataAnnotation))
			}).Should(Succeed())

			// Hub should not have annotation.
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resHub)).To(Succeed())
				anno := resHub.GetAnnotations()
				g.Expect(anno).To(HaveLen(0))
			}).Should(Succeed())
		})

		// +crdbumper:scaffold:spoketest="dataworkflowservices.Computes"
	})

	Context("SystemStatus", func() {
		var resHub *dwsv1alpha7.SystemStatus

		BeforeEach(func() {
			id := uuid.NewString()[0:8]
			resHub = &dwsv1alpha7.SystemStatus{
				ObjectMeta: metav1.ObjectMeta{
					Name:      id,
					Namespace: corev1.NamespaceDefault,
				},
				//Spec: dwsv1alpha7.SystemStatusSpec{},
			}

			Expect(k8sClient.Create(context.TODO(), resHub)).To(Succeed())
		})

		AfterEach(func() {
			if resHub != nil {
				Expect(k8sClient.Delete(context.TODO(), resHub)).To(Succeed())
				expected := &dwsv1alpha7.SystemStatus{}
				Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present.
					return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), expected)
				}).ShouldNot(Succeed())
			}
		})

		It("reads SystemStatus resource via hub and via spoke v1alpha6", func() {
			// Spoke should have annotation.
			resSpoke := &dwsv1alpha6.SystemStatus{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resSpoke)).To(Succeed())
				anno := resSpoke.GetAnnotations()
				g.Expect(anno).To(HaveLen(1))
				g.Expect(anno).Should(HaveKey(utilconversion.DataAnnotation))
			}).Should(Succeed())

			// Hub should not have annotation.
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resHub)).To(Succeed())
				anno := resHub.GetAnnotations()
				g.Expect(anno).To(HaveLen(0))
			}).Should(Succeed())
		})

		// +crdbumper:scaffold:spoketest="dataworkflowservices.SystemStatus"
	})

	// +crdbumper:scaffold:webhooksuitetest
})
