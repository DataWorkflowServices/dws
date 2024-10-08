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

	dwsv1alpha2 "github.com/DataWorkflowServices/dws/api/v1alpha2"
)

var (
	oldWorkflowGVK = schema.GroupVersionKind{
		Group:   dwsv1alpha2.GroupVersion.Group,
		Version: "v1old",
		Kind:    "Workflow",
	}

	// +crdbumper:scaffold:gvk
)

func TestMarshalData(t *testing.T) {
	g := NewWithT(t)

	t.Run("Workflow should write source object to destination", func(*testing.T) {
		src := &dwsv1alpha2.Workflow{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
				Labels: map[string]string{
					"label1": "",
				},
			},
			Spec: dwsv1alpha2.WorkflowSpec{
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
		src := &dwsv1alpha2.Workflow{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-1",
			},
		}
		dst := &unstructured.Unstructured{}
		dst.SetGroupVersionKind(dwsv1alpha2.GroupVersion.WithKind("Workflow"))
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
		src := &dwsv1alpha2.Workflow{
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

		dst := &dwsv1alpha2.Workflow{
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

		dst := &dwsv1alpha2.Workflow{
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
