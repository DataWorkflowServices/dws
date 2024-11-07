/*
Copyright 2019 The Kubernetes Authors.

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

package conversion

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"

	dwsv1alpha3 "github.com/DataWorkflowServices/dws/api/v1alpha3"
)

var (
	oldWorkflowGVK = schema.GroupVersionKind{
		Group:   dwsv1alpha3.GroupVersion.Group,
		Version: "v1old",
		Kind:    "Workflow",
	}

	oldClientMountGVK = schema.GroupVersionKind{
		Group:   dwsv1alpha3.GroupVersion.Group,
		Version: "v1old",
		Kind:    "ClientMount",
	}

	oldComputesGVK = schema.GroupVersionKind{
		Group:   dwsv1alpha3.GroupVersion.Group,
		Version: "v1old",
		Kind:    "Computes",
	}

	oldDWDirectiveRuleGVK = schema.GroupVersionKind{
		Group:   dwsv1alpha3.GroupVersion.Group,
		Version: "v1old",
		Kind:    "DWDirectiveRule",
	}

	oldDirectiveBreakdownGVK = schema.GroupVersionKind{
		Group:   dwsv1alpha3.GroupVersion.Group,
		Version: "v1old",
		Kind:    "DirectiveBreakdown",
	}

	oldPersistentStorageInstanceGVK = schema.GroupVersionKind{
		Group:   dwsv1alpha3.GroupVersion.Group,
		Version: "v1old",
		Kind:    "PersistentStorageInstance",
	}

	oldServersGVK = schema.GroupVersionKind{
		Group:   dwsv1alpha3.GroupVersion.Group,
		Version: "v1old",
		Kind:    "Servers",
	}

	oldStorageGVK = schema.GroupVersionKind{
		Group:   dwsv1alpha3.GroupVersion.Group,
		Version: "v1old",
		Kind:    "Storage",
	}

	oldSystemConfigurationGVK = schema.GroupVersionKind{
		Group:   dwsv1alpha3.GroupVersion.Group,
		Version: "v1old",
		Kind:    "SystemConfiguration",
	}

	// +crdbumper:scaffold:gvk
)

func TestMarshalData(t *testing.T) {
	g := NewWithT(t)

	t.Run("Workflow should write source object to destination", func(*testing.T) {
		src := &dwsv1alpha3.Workflow{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
				Labels: map[string]string{
					"label1": "",
				},
			},
			Spec: dwsv1alpha3.WorkflowSpec{
				DesiredState: "Proposal",
				WLMID:        "special-id",
				JobID:        intstr.FromString("my wlm job 8128"),
				UserID:       9129,
				GroupID:      7127,
			},
		}

		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldWorkflowGVK)
		dst.SetName("test-1")

		g.Expect(MarshalData(src, dst)).To(Succeed())
		// ensure the src object is not modified
		g.Expect(src.GetLabels()).ToNot(BeEmpty())

		g.Expect(dst.GetAnnotations()[DataAnnotation]).ToNot(BeEmpty())
		g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("Proposal"))
		g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("special-id"))
		g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("8128"))
		g.Expect(dst.GetAnnotations()[DataAnnotation]).ToNot(ContainSubstring("metadata"))
		g.Expect(dst.GetAnnotations()[DataAnnotation]).ToNot(ContainSubstring("label1"))
	})

	t.Run("Workflow should append the annotation", func(*testing.T) {
		src := &dwsv1alpha3.Workflow{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(dwsv1alpha3.GroupVersion.WithKind("Workflow"))
		dst.SetName("test-1")
		dst.SetAnnotations(map[string]string{
			"annotation": "1",
		})

		g.Expect(MarshalData(src, dst)).To(Succeed())
		g.Expect(dst.GetAnnotations()).To(HaveLen(2))
	})

	t.Run("ClientMount should write source object to destination", func(*testing.T) {
		src := &dwsv1alpha3.ClientMount{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
				Labels: map[string]string{
					"label1": "",
				},
			},
			//Spec: dwsv1alpha3.ClientMountSpec{
			//	// ACTION: Fill in a few valid fields so
			//	// they can be tested in the annotation checks
			//	// below.
			//},
		}

		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldClientMountGVK)
		dst.SetName("test-1")

		g.Expect(MarshalData(src, dst)).To(Succeed())
		// ensure the src object is not modified
		g.Expect(src.GetLabels()).ToNot(BeEmpty())

		g.Expect(dst.GetAnnotations()[DataAnnotation]).ToNot(BeEmpty())

		// ACTION: Fill in a few valid fields above in the Spec so
		// they can be tested here in the annotation checks.

		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mgsNids"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("rabbit-03@tcp"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mountRoot"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("/lus/w0"))
	})

	t.Run("ClientMount should append the annotation", func(*testing.T) {
		src := &dwsv1alpha3.ClientMount{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(dwsv1alpha3.GroupVersion.WithKind("ClientMount"))
		dst.SetName("test-1")
		dst.SetAnnotations(map[string]string{
			"annotation": "1",
		})

		g.Expect(MarshalData(src, dst)).To(Succeed())
		g.Expect(dst.GetAnnotations()).To(HaveLen(2))
	})

	t.Run("Computes should write source object to destination", func(*testing.T) {
		src := &dwsv1alpha3.Computes{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
				Labels: map[string]string{
					"label1": "",
				},
			},
			//Spec: dwsv1alpha3.ComputesSpec{
			//	// ACTION: Fill in a few valid fields so
			//	// they can be tested in the annotation checks
			//	// below.
			//},
		}

		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldComputesGVK)
		dst.SetName("test-1")

		g.Expect(MarshalData(src, dst)).To(Succeed())
		// ensure the src object is not modified
		g.Expect(src.GetLabels()).ToNot(BeEmpty())

		g.Expect(dst.GetAnnotations()[DataAnnotation]).ToNot(BeEmpty())

		// ACTION: Fill in a few valid fields above in the Spec so
		// they can be tested here in the annotation checks.

		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mgsNids"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("rabbit-03@tcp"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mountRoot"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("/lus/w0"))
	})

	t.Run("Computes should append the annotation", func(*testing.T) {
		src := &dwsv1alpha3.Computes{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(dwsv1alpha3.GroupVersion.WithKind("Computes"))
		dst.SetName("test-1")
		dst.SetAnnotations(map[string]string{
			"annotation": "1",
		})

		g.Expect(MarshalData(src, dst)).To(Succeed())
		g.Expect(dst.GetAnnotations()).To(HaveLen(2))
	})

	t.Run("DWDirectiveRule should write source object to destination", func(*testing.T) {
		src := &dwsv1alpha3.DWDirectiveRule{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
				Labels: map[string]string{
					"label1": "",
				},
			},
			//Spec: dwsv1alpha3.DWDirectiveRuleSpec{
			//	// ACTION: Fill in a few valid fields so
			//	// they can be tested in the annotation checks
			//	// below.
			//},
		}

		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldDWDirectiveRuleGVK)
		dst.SetName("test-1")

		g.Expect(MarshalData(src, dst)).To(Succeed())
		// ensure the src object is not modified
		g.Expect(src.GetLabels()).ToNot(BeEmpty())

		g.Expect(dst.GetAnnotations()[DataAnnotation]).ToNot(BeEmpty())

		// ACTION: Fill in a few valid fields above in the Spec so
		// they can be tested here in the annotation checks.

		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mgsNids"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("rabbit-03@tcp"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mountRoot"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("/lus/w0"))
	})

	t.Run("DWDirectiveRule should append the annotation", func(*testing.T) {
		src := &dwsv1alpha3.DWDirectiveRule{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(dwsv1alpha3.GroupVersion.WithKind("DWDirectiveRule"))
		dst.SetName("test-1")
		dst.SetAnnotations(map[string]string{
			"annotation": "1",
		})

		g.Expect(MarshalData(src, dst)).To(Succeed())
		g.Expect(dst.GetAnnotations()).To(HaveLen(2))
	})

	t.Run("DirectiveBreakdown should write source object to destination", func(*testing.T) {
		src := &dwsv1alpha3.DirectiveBreakdown{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
				Labels: map[string]string{
					"label1": "",
				},
			},
			//Spec: dwsv1alpha3.DirectiveBreakdownSpec{
			//	// ACTION: Fill in a few valid fields so
			//	// they can be tested in the annotation checks
			//	// below.
			//},
		}

		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldDirectiveBreakdownGVK)
		dst.SetName("test-1")

		g.Expect(MarshalData(src, dst)).To(Succeed())
		// ensure the src object is not modified
		g.Expect(src.GetLabels()).ToNot(BeEmpty())

		g.Expect(dst.GetAnnotations()[DataAnnotation]).ToNot(BeEmpty())

		// ACTION: Fill in a few valid fields above in the Spec so
		// they can be tested here in the annotation checks.

		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mgsNids"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("rabbit-03@tcp"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mountRoot"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("/lus/w0"))
	})

	t.Run("DirectiveBreakdown should append the annotation", func(*testing.T) {
		src := &dwsv1alpha3.DirectiveBreakdown{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(dwsv1alpha3.GroupVersion.WithKind("DirectiveBreakdown"))
		dst.SetName("test-1")
		dst.SetAnnotations(map[string]string{
			"annotation": "1",
		})

		g.Expect(MarshalData(src, dst)).To(Succeed())
		g.Expect(dst.GetAnnotations()).To(HaveLen(2))
	})

	t.Run("PersistentStorageInstance should write source object to destination", func(*testing.T) {
		src := &dwsv1alpha3.PersistentStorageInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
				Labels: map[string]string{
					"label1": "",
				},
			},
			//Spec: dwsv1alpha3.PersistentStorageInstanceSpec{
			//	// ACTION: Fill in a few valid fields so
			//	// they can be tested in the annotation checks
			//	// below.
			//},
		}

		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldPersistentStorageInstanceGVK)
		dst.SetName("test-1")

		g.Expect(MarshalData(src, dst)).To(Succeed())
		// ensure the src object is not modified
		g.Expect(src.GetLabels()).ToNot(BeEmpty())

		g.Expect(dst.GetAnnotations()[DataAnnotation]).ToNot(BeEmpty())

		// ACTION: Fill in a few valid fields above in the Spec so
		// they can be tested here in the annotation checks.

		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mgsNids"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("rabbit-03@tcp"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mountRoot"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("/lus/w0"))
	})

	t.Run("PersistentStorageInstance should append the annotation", func(*testing.T) {
		src := &dwsv1alpha3.PersistentStorageInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(dwsv1alpha3.GroupVersion.WithKind("PersistentStorageInstance"))
		dst.SetName("test-1")
		dst.SetAnnotations(map[string]string{
			"annotation": "1",
		})

		g.Expect(MarshalData(src, dst)).To(Succeed())
		g.Expect(dst.GetAnnotations()).To(HaveLen(2))
	})

	t.Run("Servers should write source object to destination", func(*testing.T) {
		src := &dwsv1alpha3.Servers{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
				Labels: map[string]string{
					"label1": "",
				},
			},
			//Spec: dwsv1alpha3.ServersSpec{
			//	// ACTION: Fill in a few valid fields so
			//	// they can be tested in the annotation checks
			//	// below.
			//},
		}

		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldServersGVK)
		dst.SetName("test-1")

		g.Expect(MarshalData(src, dst)).To(Succeed())
		// ensure the src object is not modified
		g.Expect(src.GetLabels()).ToNot(BeEmpty())

		g.Expect(dst.GetAnnotations()[DataAnnotation]).ToNot(BeEmpty())

		// ACTION: Fill in a few valid fields above in the Spec so
		// they can be tested here in the annotation checks.

		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mgsNids"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("rabbit-03@tcp"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mountRoot"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("/lus/w0"))
	})

	t.Run("Servers should append the annotation", func(*testing.T) {
		src := &dwsv1alpha3.Servers{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(dwsv1alpha3.GroupVersion.WithKind("Servers"))
		dst.SetName("test-1")
		dst.SetAnnotations(map[string]string{
			"annotation": "1",
		})

		g.Expect(MarshalData(src, dst)).To(Succeed())
		g.Expect(dst.GetAnnotations()).To(HaveLen(2))
	})

	t.Run("Storage should write source object to destination", func(*testing.T) {
		src := &dwsv1alpha3.Storage{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
				Labels: map[string]string{
					"label1": "",
				},
			},
			//Spec: dwsv1alpha3.StorageSpec{
			//	// ACTION: Fill in a few valid fields so
			//	// they can be tested in the annotation checks
			//	// below.
			//},
		}

		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldStorageGVK)
		dst.SetName("test-1")

		g.Expect(MarshalData(src, dst)).To(Succeed())
		// ensure the src object is not modified
		g.Expect(src.GetLabels()).ToNot(BeEmpty())

		g.Expect(dst.GetAnnotations()[DataAnnotation]).ToNot(BeEmpty())

		// ACTION: Fill in a few valid fields above in the Spec so
		// they can be tested here in the annotation checks.

		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mgsNids"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("rabbit-03@tcp"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mountRoot"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("/lus/w0"))
	})

	t.Run("Storage should append the annotation", func(*testing.T) {
		src := &dwsv1alpha3.Storage{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(dwsv1alpha3.GroupVersion.WithKind("Storage"))
		dst.SetName("test-1")
		dst.SetAnnotations(map[string]string{
			"annotation": "1",
		})

		g.Expect(MarshalData(src, dst)).To(Succeed())
		g.Expect(dst.GetAnnotations()).To(HaveLen(2))
	})

	t.Run("SystemConfiguration should write source object to destination", func(*testing.T) {
		src := &dwsv1alpha3.SystemConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
				Labels: map[string]string{
					"label1": "",
				},
			},
			//Spec: dwsv1alpha3.SystemConfigurationSpec{
			//	// ACTION: Fill in a few valid fields so
			//	// they can be tested in the annotation checks
			//	// below.
			//},
		}

		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldSystemConfigurationGVK)
		dst.SetName("test-1")

		g.Expect(MarshalData(src, dst)).To(Succeed())
		// ensure the src object is not modified
		g.Expect(src.GetLabels()).ToNot(BeEmpty())

		g.Expect(dst.GetAnnotations()[DataAnnotation]).ToNot(BeEmpty())

		// ACTION: Fill in a few valid fields above in the Spec so
		// they can be tested here in the annotation checks.

		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mgsNids"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("rabbit-03@tcp"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("mountRoot"))
		//g.Expect(dst.GetAnnotations()[DataAnnotation]).To(ContainSubstring("/lus/w0"))
	})

	t.Run("SystemConfiguration should append the annotation", func(*testing.T) {
		src := &dwsv1alpha3.SystemConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(dwsv1alpha3.GroupVersion.WithKind("SystemConfiguration"))
		dst.SetName("test-1")
		dst.SetAnnotations(map[string]string{
			"annotation": "1",
		})

		g.Expect(MarshalData(src, dst)).To(Succeed())
		g.Expect(dst.GetAnnotations()).To(HaveLen(2))
	})

	// +crdbumper:scaffold:marshaldata
}

func TestUnmarshalData(t *testing.T) {
	g := NewWithT(t)

	t.Run("Workflow should return false without errors if annotation doesn't exist", func(*testing.T) {
		src := &dwsv1alpha3.Workflow{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldWorkflowGVK)
		dst.SetName("test-1")

		ok, err := UnmarshalData(src, dst)
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Workflow should return true when a valid annotation with data exists", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldWorkflowGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.Workflow{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(dst.GetLabels()).To(HaveLen(1))
		g.Expect(dst.GetName()).To(Equal("test-1"))
		g.Expect(dst.GetLabels()).To(HaveKeyWithValue("label1", ""))
		g.Expect(dst.GetAnnotations()).To(BeEmpty())
	})

	t.Run("Workflow should clean the annotation on successful unmarshal", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldWorkflowGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			"annotation-1": "",
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.Workflow{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(src.GetAnnotations()).ToNot(HaveKey(DataAnnotation))
		g.Expect(src.GetAnnotations()).To(HaveLen(1))
	})

	t.Run("ClientMount should return false without errors if annotation doesn't exist", func(*testing.T) {
		src := &dwsv1alpha3.ClientMount{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldClientMountGVK)
		dst.SetName("test-1")

		ok, err := UnmarshalData(src, dst)
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("ClientMount should return true when a valid annotation with data exists", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldClientMountGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.ClientMount{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(dst.GetLabels()).To(HaveLen(1))
		g.Expect(dst.GetName()).To(Equal("test-1"))
		g.Expect(dst.GetLabels()).To(HaveKeyWithValue("label1", ""))
		g.Expect(dst.GetAnnotations()).To(BeEmpty())
	})

	t.Run("ClientMount should clean the annotation on successful unmarshal", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldClientMountGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			"annotation-1": "",
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.ClientMount{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(src.GetAnnotations()).ToNot(HaveKey(DataAnnotation))
		g.Expect(src.GetAnnotations()).To(HaveLen(1))
	})

	t.Run("Computes should return false without errors if annotation doesn't exist", func(*testing.T) {
		src := &dwsv1alpha3.Computes{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldComputesGVK)
		dst.SetName("test-1")

		ok, err := UnmarshalData(src, dst)
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Computes should return true when a valid annotation with data exists", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldComputesGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.Computes{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(dst.GetLabels()).To(HaveLen(1))
		g.Expect(dst.GetName()).To(Equal("test-1"))
		g.Expect(dst.GetLabels()).To(HaveKeyWithValue("label1", ""))
		g.Expect(dst.GetAnnotations()).To(BeEmpty())
	})

	t.Run("Computes should clean the annotation on successful unmarshal", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldComputesGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			"annotation-1": "",
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.Computes{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(src.GetAnnotations()).ToNot(HaveKey(DataAnnotation))
		g.Expect(src.GetAnnotations()).To(HaveLen(1))
	})

	t.Run("DWDirectiveRule should return false without errors if annotation doesn't exist", func(*testing.T) {
		src := &dwsv1alpha3.DWDirectiveRule{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldDWDirectiveRuleGVK)
		dst.SetName("test-1")

		ok, err := UnmarshalData(src, dst)
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("DWDirectiveRule should return true when a valid annotation with data exists", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldDWDirectiveRuleGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.DWDirectiveRule{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(dst.GetLabels()).To(HaveLen(1))
		g.Expect(dst.GetName()).To(Equal("test-1"))
		g.Expect(dst.GetLabels()).To(HaveKeyWithValue("label1", ""))
		g.Expect(dst.GetAnnotations()).To(BeEmpty())
	})

	t.Run("DWDirectiveRule should clean the annotation on successful unmarshal", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldDWDirectiveRuleGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			"annotation-1": "",
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.DWDirectiveRule{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(src.GetAnnotations()).ToNot(HaveKey(DataAnnotation))
		g.Expect(src.GetAnnotations()).To(HaveLen(1))
	})

	t.Run("DirectiveBreakdown should return false without errors if annotation doesn't exist", func(*testing.T) {
		src := &dwsv1alpha3.DirectiveBreakdown{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldDirectiveBreakdownGVK)
		dst.SetName("test-1")

		ok, err := UnmarshalData(src, dst)
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("DirectiveBreakdown should return true when a valid annotation with data exists", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldDirectiveBreakdownGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.DirectiveBreakdown{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(dst.GetLabels()).To(HaveLen(1))
		g.Expect(dst.GetName()).To(Equal("test-1"))
		g.Expect(dst.GetLabels()).To(HaveKeyWithValue("label1", ""))
		g.Expect(dst.GetAnnotations()).To(BeEmpty())
	})

	t.Run("DirectiveBreakdown should clean the annotation on successful unmarshal", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldDirectiveBreakdownGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			"annotation-1": "",
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.DirectiveBreakdown{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(src.GetAnnotations()).ToNot(HaveKey(DataAnnotation))
		g.Expect(src.GetAnnotations()).To(HaveLen(1))
	})

	t.Run("PersistentStorageInstance should return false without errors if annotation doesn't exist", func(*testing.T) {
		src := &dwsv1alpha3.PersistentStorageInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldPersistentStorageInstanceGVK)
		dst.SetName("test-1")

		ok, err := UnmarshalData(src, dst)
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("PersistentStorageInstance should return true when a valid annotation with data exists", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldPersistentStorageInstanceGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.PersistentStorageInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(dst.GetLabels()).To(HaveLen(1))
		g.Expect(dst.GetName()).To(Equal("test-1"))
		g.Expect(dst.GetLabels()).To(HaveKeyWithValue("label1", ""))
		g.Expect(dst.GetAnnotations()).To(BeEmpty())
	})

	t.Run("PersistentStorageInstance should clean the annotation on successful unmarshal", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldPersistentStorageInstanceGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			"annotation-1": "",
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.PersistentStorageInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(src.GetAnnotations()).ToNot(HaveKey(DataAnnotation))
		g.Expect(src.GetAnnotations()).To(HaveLen(1))
	})

	t.Run("Servers should return false without errors if annotation doesn't exist", func(*testing.T) {
		src := &dwsv1alpha3.Servers{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldServersGVK)
		dst.SetName("test-1")

		ok, err := UnmarshalData(src, dst)
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Servers should return true when a valid annotation with data exists", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldServersGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.Servers{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(dst.GetLabels()).To(HaveLen(1))
		g.Expect(dst.GetName()).To(Equal("test-1"))
		g.Expect(dst.GetLabels()).To(HaveKeyWithValue("label1", ""))
		g.Expect(dst.GetAnnotations()).To(BeEmpty())
	})

	t.Run("Servers should clean the annotation on successful unmarshal", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldServersGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			"annotation-1": "",
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.Servers{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(src.GetAnnotations()).ToNot(HaveKey(DataAnnotation))
		g.Expect(src.GetAnnotations()).To(HaveLen(1))
	})

	t.Run("Storage should return false without errors if annotation doesn't exist", func(*testing.T) {
		src := &dwsv1alpha3.Storage{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldStorageGVK)
		dst.SetName("test-1")

		ok, err := UnmarshalData(src, dst)
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("Storage should return true when a valid annotation with data exists", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldStorageGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.Storage{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(dst.GetLabels()).To(HaveLen(1))
		g.Expect(dst.GetName()).To(Equal("test-1"))
		g.Expect(dst.GetLabels()).To(HaveKeyWithValue("label1", ""))
		g.Expect(dst.GetAnnotations()).To(BeEmpty())
	})

	t.Run("Storage should clean the annotation on successful unmarshal", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldStorageGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			"annotation-1": "",
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.Storage{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(src.GetAnnotations()).ToNot(HaveKey(DataAnnotation))
		g.Expect(src.GetAnnotations()).To(HaveLen(1))
	})

	t.Run("SystemConfiguration should return false without errors if annotation doesn't exist", func(*testing.T) {
		src := &dwsv1alpha3.SystemConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(oldSystemConfigurationGVK)
		dst.SetName("test-1")

		ok, err := UnmarshalData(src, dst)
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("SystemConfiguration should return true when a valid annotation with data exists", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldSystemConfigurationGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.SystemConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(dst.GetLabels()).To(HaveLen(1))
		g.Expect(dst.GetName()).To(Equal("test-1"))
		g.Expect(dst.GetLabels()).To(HaveKeyWithValue("label1", ""))
		g.Expect(dst.GetAnnotations()).To(BeEmpty())
	})

	t.Run("SystemConfiguration should clean the annotation on successful unmarshal", func(*testing.T) {
		src := &unstructured.Unstructured{}
		src.SetGroupVersionKind(oldSystemConfigurationGVK)
		src.SetName("test-1")
		src.SetAnnotations(map[string]string{
			"annotation-1": "",
			DataAnnotation: "{\"metadata\":{\"name\":\"test-1\",\"creationTimestamp\":null,\"labels\":{\"label1\":\"\"}},\"spec\":{},\"status\":{}}",
		})

		dst := &dwsv1alpha3.SystemConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}

		ok, err := UnmarshalData(src, dst)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(ok).To(BeTrue())

		g.Expect(src.GetAnnotations()).ToNot(HaveKey(DataAnnotation))
		g.Expect(src.GetAnnotations()).To(HaveLen(1))
	})

	// +crdbumper:scaffold:unmarshaldata
}

// Just touch ginkgo, so it's here to interpret any ginkgo args from
// "make test", so that doesn't fail on this test file.
var _ = BeforeSuite(func() {})
