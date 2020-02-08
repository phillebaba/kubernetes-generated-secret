package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:validation:Enum=Uppercase;Lowercase;Numbers;Symbols
type ValueOption string

const (
	Uppercase = "Uppercase"
	Lowercase = "Lowercase"
	Numbers   = "Numbers"
	Symbols   = "Symbols"
)

func (c ValueOption) Regex() string {
	switch c {
	case Uppercase:
		return "[A-Z]"
	case Lowercase:
		return "[a-z]"
	case Numbers:
		return "[0-9]"
	case Symbols:
		return "[!#$%&()*+,-./:;<=>?@[\\]^_`{|}~]"
	default:
		return ""
	}
}

type GeneratedSecretData struct {
	// Key of the secret
	Key string `json:"key"`

	// Wanted length of the secret value
	// +kubebuilder:validation:Minimum=1
	Length int `json:"length"`

	// Options to apply to the generated secret value
	ValueOptions []ValueOption `json:"options,uniqueItems"`
}

// +kubebuilder:subresource:data

// GeneratedSecretSpec defines the desired state of GeneratedSecret
type GeneratedSecretSpec struct {
	// Addtional metadata to add to the generated secret.
	// +optional
	SecretMeta metav1.ObjectMeta `json:"secretMetadata,omitempty"`

	// Data configuration of the secret
	DataList []GeneratedSecretData `json:"data"`
}

// GeneratedSecretStatus defines the observed state of GeneratedSecret
type GeneratedSecretStatus struct {
}

// +kubebuilder:object:root=true

// GeneratedSecret is the Schema for the generatedsecrets API
type GeneratedSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GeneratedSecretSpec   `json:"spec,omitempty"`
	Status GeneratedSecretStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GeneratedSecretList contains a list of GeneratedSecret
type GeneratedSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GeneratedSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GeneratedSecret{}, &GeneratedSecretList{})
}
