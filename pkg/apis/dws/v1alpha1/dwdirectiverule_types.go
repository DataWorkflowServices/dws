package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DWDirectiveRuleDef defines the DWDirective parser rules
type DWDirectiveRuleDef struct {
	Key				string		`json:"key"`
	Type			string		`json:"type"`
	Pattern			string		`json:"pattern,omitempty"`
	Min				int			`json:"min,omitempty"`
	Max				int			`json:"max,omitempty"`
	IsRequired		bool		`json:"isRrequired,omitempty"`
	IsValueRequired	bool		`json:"isValueRequired,omitempty"`
}

// DWDirectiveRuleSpec defines the desired state of DWDirective
type DWDirectiveRuleSpec struct {
	Command		string					`json:"command"`
	RuleDefs	[]DWDirectiveRuleDef	`json:"ruleDefs"`
}

// +genclient
// +genclient:noStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DWDirectiveRule is the Schema for the DWDirective API
type DWDirectiveRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   []DWDirectiveRuleSpec   `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DWDirectiveRuleList contains a list of DWDirective
type DWDirectiveRuleList struct {
	metav1.TypeMeta	`json:",inline"`
	metav1.ListMeta	`json:"metadata,omitempty"`
	Items	[]DWDirectiveRule	`json:"items"`
}

func init() {
	SchemeBuilder.Register(&DWDirectiveRule{}, &DWDirectiveRuleList{})
}
