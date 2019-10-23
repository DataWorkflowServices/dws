package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StoragePoolSpec defines the desired state of StoragePool
// +k8s:openapi-gen=true
type StoragePoolSpec struct {
    PoolID            string `json:"poolID"`
    Units             string `json:"units"`
    Granularity       int `json:"granularity"`
    Quantity          int `json:"quantity"`
    Free              int `json:"free"`
}

// StoragePoolStatus defines the observed state of StoragePool
// +k8s:openapi-gen=true
type StoragePoolStatus struct {
    State  string `json:"state"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StoragePool is the Schema for the storagepools API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type StoragePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StoragePoolSpec   `json:"spec,omitempty"`
	Status StoragePoolStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StoragePoolList contains a list of StoragePool
type StoragePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StoragePool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StoragePool{}, &StoragePoolList{})
}
