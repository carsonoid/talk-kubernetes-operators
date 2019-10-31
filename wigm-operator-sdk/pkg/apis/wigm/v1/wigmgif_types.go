package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WigmGifSpec defines the desired state of WigmGif
type WigmGifSpec struct {
	Gif GifProperties `json:"gif"`

	Service *ServiceProperties `json:"service,omitempty"`

	Ingress *IngressProperties `json:"ingress,omitempty"`
}

type GifProperties struct {
	Title string `json:"title,omitempty"`

	Link string `json:"link"`
}

type ServiceProperties struct {
	CreateCloudLB bool `json:"create_cloud_lb"`
}

type IngressProperties struct {
	// Enabled is a pointer to a bool so that it is possible
	// to tell the difference between the value not being set
	// and the value being explicitly set to false
	Enabled *bool `json:"enabled"`
}

// WigmGifStatus defines the observed state of WigmGif
type WigmGifStatus struct {
	Deployment DeploymentStatus `json:"deployment"`

	Service ServiceStatus `json:"service"`

	Ingress IngressStatus `json:"ingress"`
}

type DeploymentStatus struct {
	Created bool `json:"created"`
}

type ServiceStatus struct {
	Created bool               `json:"created"`
	Type    corev1.ServiceType `json:"type"`
}

type IngressStatus struct {
	Created bool `json:"created"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WigmGif is the Schema for the wigmgifs API
// +k8s:openapi-gen=true
type WigmGif struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WigmGifSpec   `json:"spec,omitempty"`
	Status WigmGifStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WigmGifList contains a list of WigmGif
type WigmGifList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WigmGif `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WigmGif{}, &WigmGifList{})
}
