package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WorkflowSpec defines the desired state of Workflow
// +k8s:openapi-gen=true
// +genclient
// +genclient:noStatus
type WorkflowSpec struct {
	DesiredState	string				`json:"desiredState"`
	WLMID			string				`json:"wlmID"`
	JobID			int					`json:"jobID"`
	UserID			int					`json:"userID"`
	DWDirectives	[]string			`json:"dwDirectives"`
}

type WorkflowDriverStatus struct {
	DriverID	string				`json:"driverID"`
    DWDIndex	int					`json:"dwdIndex"`
	WatchState	string				`json:"watchState"`
    LastHB		int64				`json:"lastHB"`
    Completed	bool				`json:"completed"`
	// User readable reason.
	// For the CDS driver, this could be the state of the underlying
	// data movement request:  Pending, Queued, Running, Completed or Error
	Reason		string				`json:"reason,omitempty"`
}

// WorkflowStatus defines the observed state of the Workflow
// +k8s:openapi-gen=true
type WorkflowStatus struct {
	// The state the resource is currently transitioning to.
	// Updated by the controller once started.
    State		string					`json:"state"`

	// Ready can be 'True', 'False'
	// Indicates whether State has been reached.
    Ready       bool					`json:"ready"`

	// User readable reason and status message
	Reason		string					 `json:"reason,omitempty"`
	Message		string					 `json:"message,omitempty"`

	// Set of DW environment variable settings for WLM to apply to the job.
	//		- DW_JOB_STRIPED
	//		- DW_JOB_PRIVATE
	//		- DW_JOB_STRIPED_CACHE
	//		- DW_JOB_LDBAL_CACHE
	//		- DW_PERSISTENT_STRIPED_{resname}
    Env         []string				`json:"env,omitempty"`

	// List of registered drivers and related status.  Updated by drivers.
    Drivers     []WorkflowDriverStatus	`json:"drivers,omitempty"`
}

// Workflow is the Schema for the workflows API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
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
