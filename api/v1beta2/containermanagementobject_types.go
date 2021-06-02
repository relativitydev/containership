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

// ContainerManagementObjectSpec defines the desired state of ContainerManagementObject
type ContainerManagementObjectSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Images is a list of images that containership will manage
	// kubebuilder:validation:Required
	Images []Image `json:"images"`
}

type Image struct {
	// SourceRepository is where the image will be pulled from. It is the source of truth.
	SourceRepository string `json:"sourceRepository"`

	/*
	 TargetRepository is an optional field that allows the image repository to be renamed.
	 If sourceRepository is "docker.io/library/busybox", setting targetRepository to "hello-world/busybox" will
	 rename the image "hello-world/busybox"
	*/
	TargetRepository string `json:"targetRepository,omitempty"`

	// SupportedTags are the image tags that will pulled from the source repository. Any extra tags found in the target image destinations will be deleted.
	SupportedTags []string `json:"supportedTags,omitempty"`
}

// ContainerManagementObjectStatus defines the observed state of ContainerManagementObject
type ContainerManagementObjectStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ContainerManagementObject is the Schema for the containermanagementobjects API
type ContainerManagementObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContainerManagementObjectSpec   `json:"spec,omitempty"`
	Status ContainerManagementObjectStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ContainerManagementObjectList contains a list of ContainerManagementObject
type ContainerManagementObjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContainerManagementObject `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContainerManagementObject{}, &ContainerManagementObjectList{})
}
