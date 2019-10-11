package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DriverSpec defines the desired state of Driver
// +k8s:openapi-gen=true
type DriverSpec struct {
    DriverId             string `json:"driverId"`
    WatchStates          string `json:"watchStates"`
}

// DriverStatus defines the observed state of Driver
// +k8s:openapi-gen=true
type DriverStatusSpec struct {
    DriverId      string `json:"driverId"`
    WatchStates   string `json:"watchStates"`
    Status        []string `json:"status"`
}

type DriverStatus struct {
    Registrations        []DriverStatusSpec `json:"registrations"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Driver is the Schema for the drivers API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Driver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DriverSpec   `json:"spec,omitempty"`
	Status DriverStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DriverList contains a list of Driver
type DriverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Driver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Driver{}, &DriverList{})
}
