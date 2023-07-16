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

// EntandoBundleInstanceV2Spec defines the desired state of EntandoBundleInstanceV2
type EntandoBundleInstanceV2Spec struct {
	Tag        string `json:"tag,omitempty"`
	Digest     string `json:"digest,omitempty"`
	Repository string `json:"repository,omitempty"`
	// FIXME vanno inserite in annotations Dependencies  []string `json:"dependencies,omitempty"`
	// FIXME vanno inserite in annotations Components    []string `json:"components,omitempty"`
	DesiredStatus string `json:"desiredStatus,omitempty"`
	Configuration string `json:"configuration,omitempty"`
}

// EntandoBundleInstanceV2Status defines the observed state of EntandoBundleInstanceV2
type EntandoBundleInstanceV2Status struct {
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EntandoBundleInstanceV2 is the Schema for the entandobundleinstancev2s API
type EntandoBundleInstanceV2 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EntandoBundleInstanceV2Spec   `json:"spec,omitempty"`
	Status EntandoBundleInstanceV2Status `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EntandoBundleInstanceV2List contains a list of EntandoBundleInstanceV2
type EntandoBundleInstanceV2List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EntandoBundleInstanceV2 `json:"items"`
}

func (p *EntandoBundleInstanceV2) GetConditions() []metav1.Condition {
	return p.Status.Conditions
}

func (p *EntandoBundleInstanceV2) SetConditions(conditions []metav1.Condition) {
	p.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&EntandoBundleInstanceV2{}, &EntandoBundleInstanceV2List{})
}
