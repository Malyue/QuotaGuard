package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type QuotaPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec QuotaPolicySpec `json:"spec"`
}

type QuotaPolicySpec struct {
	Rule []QuotaRule `json:"rules"`
}

type QuotaRule struct {
	Target QuotaTarget `json:"target"`
	Limit  QuotaLimit  `json:"limit"`
}

type QuotaTarget struct {
	Kind string `json:"kind"`
	Key  string `json:"key"`
}

type QuotaLimit struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type QuotaPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QuotaPolicy `json:"items"`
}
