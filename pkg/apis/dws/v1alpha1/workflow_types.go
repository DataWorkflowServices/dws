package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Capacity struct {
	Units string `json:"units"`
	Size  int    `json:"size"`
}

type DWRecord struct {
	// Array index of the #DW directive in original WFR
	DWDirectiveIndex int `json:"dwDirectiveIndex"`
	// Copy of the #DW for this breakdown
	DWDirective string `json:"dwDirective"`
}

// AllocationSetComponents define the details of the allocation
type AllocationSetComponents struct {
	DW DWRecord `json:"dwRecord"`
	// The allowed set of AllocationStategy's:
	// AllocatePerCompute
	// DivideAcrossRabbits
	// SingleRabbit
	AllocationStrategy string   `json:"allocationStrategy"`
	MinimumCapacity    Capacity `json:"minimumCapacity"`
	Labels             []string `json:"labels"`
	ComputeBindings    []string `json:"computeBindings"`
}

type StorageResourceDescriptor struct {
	DW        DWRecord                `json:"dwRecord"`
	Name      string                  `json:"name"`
	Reference *corev1.ObjectReference `json:"storageResourceRef"`
}

// DWDirectiveBreakdowns define the storage information WLM needs to select NNF Nodes and request storage from the selected nodes
type DWDirectiveBreakdownAllocationSet struct {
	AllocationSet []AllocationSetComponents `json:"allocationSet"`
}

// WorkflowSpec defines the desired state of Workflow
type WorkflowSpec struct {
	DesiredState string   `json:"desiredState"`
	WLMID        string   `json:"wlmID"`
	JobID        int      `json:"jobID"`
	UserID       int      `json:"userID"`
	DWDirectives []string `json:"dwDirectives"`
}

// WorkflowDriverStatus defines the status information provided by integration drivers.
type WorkflowDriverStatus struct {
	DriverID   string `json:"driverID"`
	TaskID     string `json:"taskID"`
	DWDIndex   int    `json:"dwdIndex"`
	WatchState string `json:"watchState"`
	LastHB     int64  `json:"lastHB"`
	Completed  bool   `json:"completed"`
	// User readable reason.
	// For the CDS driver, this could be the state of the underlying
	// data movement request:  Pending, Queued, Running, Completed or Error
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// WorkflowStatus defines the observed state of the Workflow
type WorkflowStatus struct {
	// The state the resource is currently transitioning to.
	// Updated by the controller once started.
	State string `json:"state"`

	// Ready can be 'True', 'False'
	// Indicates whether State has been reached.
	Ready bool `json:"ready"`

	// User readable reason and status message
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`

	// #DW directive breakdowns indicating to WLM what to allocate on what Rabbit
	DWDirectiveBreakdowns []DWDirectiveBreakdownAllocationSet `json:"dwDirectiveBreakdowns,omitempty"`

	// A StorageResource is created for each #DW to express to the NNF Driver how to create storage
	StorageResource []StorageResourceDescriptor `json:"storageResourceDesc,omitempty"`

	// Set of DW environment variable settings for WLM to apply to the job.
	//		- DW_JOB_STRIPED
	//		- DW_JOB_PRIVATE
	//		- DW_JOB_STRIPED_CACHE
	//		- DW_JOB_LDBAL_CACHE
	//		- DW_PERSISTENT_STRIPED_{resname}
	Env []string `json:"env,omitempty"`

	// List of registered drivers and related status.  Updated by drivers.
	Drivers []WorkflowDriverStatus `json:"drivers,omitempty"`
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Workflow is the Schema for the workflows API
type Workflow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkflowSpec   `json:"spec,omitempty"`
	Status WorkflowStatus `json:"status,omitempty"`
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
