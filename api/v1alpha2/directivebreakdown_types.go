/*
 * Copyright 2023 Hewlett Packard Enterprise Development LP
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DirectiveBreakdownSpec defines the desired state of DirectiveBreakdown
type DirectiveBreakdownSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of DirectiveBreakdown. Edit directivebreakdown_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// DirectiveBreakdownStatus defines the observed state of DirectiveBreakdown
type DirectiveBreakdownStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DirectiveBreakdown is the Schema for the directivebreakdowns API
type DirectiveBreakdown struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DirectiveBreakdownSpec   `json:"spec,omitempty"`
	Status DirectiveBreakdownStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DirectiveBreakdownList contains a list of DirectiveBreakdown
type DirectiveBreakdownList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DirectiveBreakdown `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DirectiveBreakdown{}, &DirectiveBreakdownList{})
}
