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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretType identifies the type of use for a secret inside a deployment
// +enum
type SecretType string

const (
	// SecretTypeEnv means that the secret will be used via env
	SecretTypeEnv SecretType = "ENV"
	// SecretTypeFile means that the secret will be used via mounted file
	SecretTypeFile SecretType = "FILE"
)

type EntandoPluginV2Secret struct {
	SecretType SecretType `json:"secretType,omitempty"`
	Name       string     `json:"name,omitempty"`
	Prefix     string     `json:"prefix,omitempty"`
	MountPath  string     `json:"mountPath,omitempty"`
}

type EntandoPluginV2Volume struct {
	StorageClass string `json:"storagClass,omitempty"`
	Size         string `json:"size,omitempty"`
	MountPath    string `json:"mountPath,omitempty"`
}

// EntandoPluginV2Spec defines the desired state of EntandoPluginV2
type EntandoPluginV2Spec struct {
	// +kubebuilder:default:="none"
	Database             string                  `json:"database,omitempty"`
	EnvironmentVariables []corev1.EnvVar         `json:"environmentVariables,omitempty"`
	Secrets              []EntandoPluginV2Secret `json:"secrets,omitempty"`
	Volumes              []EntandoPluginV2Volume `json:"volumes,omitempty"`
	HealthCheckPath      string                  `json:"healthCheckPath,omitempty"`
	IngressName          string                  `json:"ingressName,omitempty"`
	IngressHost          string                  `json:"ingressHost,omitempty"`
	IngressPath          string                  `json:"ingressPath,omitempty"`
	Image                string                  `json:"image,omitempty"`
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`
	// +kubebuilder:default:=8080
	Port int32 `json:"port,omitempty"`
}

// EntandoPluginV2Status defines the observed state of EntandoPluginV2
type EntandoPluginV2Status struct {
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EntandoPluginV2 is the Schema for the entandopluginv2s API
//+kubebuilder:printcolumn:name="Ready",type="string",JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="state of Plugin"
type EntandoPluginV2 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EntandoPluginV2Spec   `json:"spec,omitempty"`
	Status EntandoPluginV2Status `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EntandoPluginV2List contains a list of EntandoPluginV2
type EntandoPluginV2List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EntandoPluginV2 `json:"items"`
}

func (p *EntandoPluginV2) GetConditions() []metav1.Condition {
	return p.Status.Conditions
}

func (p *EntandoPluginV2) SetConditions(conditions []metav1.Condition) {
	p.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&EntandoPluginV2{}, &EntandoPluginV2List{})
}
