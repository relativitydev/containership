/*
Copyright 2021.

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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RegistriesConfigSpec defines the desired state of RegistriesConfig
type RegistriesConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Regsitries is a list of container registries authentication credentials
	// kubebuilder:validation:UniqueItems
	Registries []Registry `json:"registries,omitempty"`
}

type Registry struct {
	// Name is a unique name for the registry. This name is the key to match when listed as an image destination
	Name string `json:"name"`

	// Uri describes the URI for the registry (example: docker.io)
	URI string `json:"uri"`

	// The name of the secret containing the authorization credentials for the registry. It must exist in the same namespace as the operator. Secret type must be kubernetes.io/dockerconfigjson
	SecretName string `json:"secretName,omitempty"`
}

// RegistriesConfigStatus defines the observed state of RegistriesConfig
type RegistriesConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RegistriesConfig is the Schema for the registriesconfigs API
type RegistriesConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegistriesConfigSpec   `json:"spec,omitempty"`
	Status RegistriesConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RegistriesConfigList contains a list of RegistriesConfig
type RegistriesConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RegistriesConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RegistriesConfig{}, &RegistriesConfigList{})
}
