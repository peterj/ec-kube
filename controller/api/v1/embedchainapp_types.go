/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EmbedchainAppSpec defines the desired state of EmbedchainApp
type EmbedchainAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ConfigRef is a reference to a ConfigMap that contains the configuration
	//+kubebuilder:validation:Required
	ConfigRef *ConfigRefSpec `json:"configRef,omitempty"`

	// SecretRef is a reference to a Secret that contains the configuration
	//+kubebuilder:validation:Required
	SecretRef *SecretRefSpec `json:"secretRef,omitempty"`
}

// +kubebuilder:printcolumn:name="ConfigMap",type="string",JSONPath=".spec.configRef.name"
type ConfigRefSpec struct {
	// ConfigMapName is the name of the ConfigMap that contains the configuration
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// +kubebuilder:printcolumn:name="Secret",type="string",JSONPath=".spec.secretRef.name"
type SecretRefSpec struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// EmbedchainAppStatus defines the observed state of EmbedchainApp
type EmbedchainAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EmbedchainApp is the Schema for the embedchainapps API
type EmbedchainApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmbedchainAppSpec   `json:"spec,omitempty"`
	Status EmbedchainAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EmbedchainAppList contains a list of EmbedchainApp
type EmbedchainAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EmbedchainApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EmbedchainApp{}, &EmbedchainAppList{})
}
