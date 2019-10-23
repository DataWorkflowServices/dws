package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WorkflowSpec defines the desired state of Workflow
// +k8s:openapi-gen=true
type WorkflowSpec struct {
	DesiredState	  string	`json:"desiredState"`
    WLMID             string	`json:"wlmID"`
    JobID             string	`json:"jobID"`
    DWDirectives      []string	`json:"dwDirectives"`
    UserID            int		`json:"userID"`
    Env               []string	`json:"env,omitempty"`
    DefaultPool       string	`json:"defaultPool,omitempty"`
}

// WorkflowStatus defines the observed state of Workflow
// +k8s:openapi-gen=true
type WorkflowStatus struct {
	State			string				`json:"state"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Workflow is the Schema for the workflows API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Workflow struct {
	metav1.TypeMeta			`json:",inline"`
	metav1.ObjectMeta		`json:"metadata,omitempty"`

	Spec   WorkflowSpec		`json:"spec,omitempty"`
	Status WorkflowStatus	`json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WorkflowList contains a list of Workflow
type WorkflowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workflow `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Workflow{}, &WorkflowList{})
}
