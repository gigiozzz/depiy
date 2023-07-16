/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EntandoGatewayV2Spec defines the desired state of EntandoGatewayV2
type EntandoGatewayV2Spec struct {
	IngressName    string `json:"ingressName,omitempty"`
	IngressHost    string `json:"ingressHost,omitempty"`
	IngressPath    string `json:"ingressPath,omitempty"`
	IngressPort    string `json:"ingressPort,omitempty"`
	IngressService string `json:"ingressService,omitempty"`
}

// EntandoGatewayV2Status defines the observed state of EntandoGatewayV2
type EntandoGatewayV2Status struct {
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EntandoGatewayV2 is the Schema for the EntandoGatewayV2s API
//+kubebuilder:printcolumn:name="Ready",type="string",JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="state of Gateway"
type EntandoGatewayV2 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EntandoGatewayV2Spec   `json:"spec,omitempty"`
	Status EntandoGatewayV2Status `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EntandoGatewayV2List contains a list of EntandoGatewayV2
type EntandoGatewayV2List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EntandoGatewayV2 `json:"items"`
}

func (p *EntandoGatewayV2) GetConditions() []metav1.Condition {
	return p.Status.Conditions
}

func (p *EntandoGatewayV2) SetConditions(conditions []metav1.Condition) {
	p.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&EntandoGatewayV2{}, &EntandoGatewayV2List{})
}
