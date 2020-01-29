/*

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

type GeneratedSecretData struct {
	Key string `json:"key"`

	// Specifies the wanted length of the generated secret.
	// +optional
	Length *int `json:"length,omitempty"`

	// Specifies if letters should be used in the generated secret.
	// +optional
	Letters *bool `json:"letters,omitempty"`

	// Specifies if numbers should be used in the generated secret.
	// +optional
	Numbers *bool `json:"numbers,omitempty"`

	// Specifies if special characters should be used in the generated secret.
	// +optional
	Special *bool `json:"special,omitempty"`
}

// +kubebuilder:subresource:data

// GeneratedSecretSpec defines the desired state of GeneratedSecret
type GeneratedSecretSpec struct {
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
