package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ContainerManagementObjectSpec defines the desired state of ContainerManagementObject
type ContainerManagementObjectSpec struct {
	// Important: Run "make" to regenerate code after modifying this file
	// SlackWebhookEndpoint is the API endpoint for the Slack channel to which you want alerts sent
	SlackWebhookEndpoint string `json:"slackWebhookEndpoint,omitempty"`

	// Data regarding the images you want pushed with Containership
	Images []Image `json:"images"`
}

// Images contains the metadata for the images to be managed by Containership
type Image struct {
	// SourceImage is the name of the external image to pull. Should not include tags. ex: docker.io/library/busybox
	// +kubebuilder:validation:Required
	SourceImage string `json:"sourceImage"`

	// TargetRepository is the name of the image repository that will be pushed to the destination ACR. Only use if you want to rename the image. Do not include tags.
	// +kubebuilder:validation:Optional
	TargetRepository string `json:"targetRepository,omitempty"`

	// SupportedTags are the tags that will be pulled externally and pushed into the destinations
	// +kubebuilder:validation:Required
	SupportedTags []string `json:"supportedTags"`

	// Destinations are the intended destinations for all imported tags
	// +kubebuilder:validation:Required
	Destinations Destinations `json:"destinations"`
}

// Destinations are the intended destinations for all imported tags
type Destinations struct {
	// +kubebuilder:validation:Optional
	AzureContainerRegistries []AzureContainerRegistries `json:"azurecontainerregistries,omitempty"`
	// +kubebuilder:validation:Optional
	GoogleContainerRegistry []GoogleContainerRegistries `json:"googlecontainerregistries,omitempty"`
}

// AzureContainerRegistries is required information to interact with Azure container registries
type AzureContainerRegistries struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// SubscriptionID is optional. Set if you want to override the global subscriptionID (envvar AZURE_SUBSCRIPTION_ID)
	// +kubebuilder:validation:Optional
	SubscriptionID string `json:"subscriptionId,omitempty"`

	// ResourceGroup is optional. Set if you want to override the global resource group (envvar AZURE_RESOURCE_GROUP)
	// +kubebuilder:validation:Optional
	ResourceGroup string `json:"resourceGroup,omitempty"`

	// Ring is an optional field that indicates an image tag should import from a destination with a subsequent ring value rather than source. Destinations with higher ring values import tags from destinations with lower ring values.
	Ring int `json:"ring,omitempty"`
}

// GoogleContainerRegistries is required information to interact with Google container registries
type GoogleContainerRegistries struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// ContainerManagementObjectStatus defines the observed state of ContainerManagementObject
type ContainerManagementObjectStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// ContainerManagementObject is the Schema for the containermanagementobjects API
// +kubebuilder:resource:path=containermanagementobjects,scope=Namespaced,shortName=cmo
type ContainerManagementObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContainerManagementObjectSpec   `json:"spec,omitempty"`
	Status ContainerManagementObjectStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ContainerManagementObjectList contains a list of ContainerManagementObject
type ContainerManagementObjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContainerManagementObject `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContainerManagementObject{}, &ContainerManagementObjectList{})
}
