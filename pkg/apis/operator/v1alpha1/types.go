package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FederationV2List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []FederationV2 `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FederationV2 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              FederationV2Spec   `json:"spec"`
	Status            FederationV2Status `json:"status,omitempty"`
}

type FederationV2Spec struct {
	PushReconciliation bool `json:"pushReconciliation,omitempty"`
	Scheduling         bool `json:"scheduling,omitempty"`
	ServiceDiscovery   bool `json:"serviceDiscovery,omitempty"`
}
type FederationV2Status struct {
	// Fill me
}
