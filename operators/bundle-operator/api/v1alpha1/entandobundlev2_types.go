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

// SignatureType identifies the type of key to use to verify signature
// +enum
type SignatureType string

const (
	// SignatureKeyPair means that the signature will use a priv/pub key pair
	SignatureKeyPair SignatureType = "KEY_PAIR"
	// SignatureKeyPair means that the signature will use the Fulcio OIDC flow
	SignatureKeyLess SignatureType = "KEY_LESS"
)

type SignatureInfo struct {
	Type SignatureType `json:"type,omitempty"`
	//Image        string        `json:"image,omitempty"`
	PubKey       string `json:"pubKey,omitempty"`
	PubKeySecret string `json:"pubKeySecret,omitempty"`
}

type EntandoBundleTag struct {
	Tag           string          `json:"tag,omitempty"`
	Digest        string          `json:"digest,omitempty"`
	SignatureInfo []SignatureInfo `json:"signatureInfo,omitempty"`
}

// EntandoBundleV2Spec defines the desired state of EntandoBundleV2
type EntandoBundleV2Spec struct {
	Title         string             `json:"title,omitempty"`
	Icon          string             `json:"icon,omitempty"`
	SignatureInfo string             `json:"signatureInfo,omitempty"`
	Repository    string             `json:"repository,omitempty"`
	TagList       []EntandoBundleTag `json:"tagList,omitempty"`
}

// EntandoBundleV2Status defines the observed state of EntandoBundleV2
type EntandoBundleV2Status struct {
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EntandoBundleV2 is the Schema for the entandobundlev2s API
//+kubebuilder:printcolumn:name="Ready",type="string",JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="state of Plugin"
type EntandoBundleV2 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EntandoBundleV2Spec   `json:"spec,omitempty"`
	Status EntandoBundleV2Status `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EntandoBundleV2List contains a list of EntandoBundleV2
type EntandoBundleV2List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EntandoBundleV2 `json:"items"`
}

func (p *EntandoBundleV2) GetConditions() []metav1.Condition {
	return p.Status.Conditions
}

func (p *EntandoBundleV2) SetConditions(conditions []metav1.Condition) {
	p.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&EntandoBundleV2{}, &EntandoBundleV2List{})
}
