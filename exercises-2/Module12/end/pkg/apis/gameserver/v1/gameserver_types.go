package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GameserverSpec defines the desired state of Gameserver
type GameserverSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	GameID      string `json:"gameid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ServerPort  int32  `json:"port"`
}

// GameserverStatus defines the observed state of Gameserver
type GameserverStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Gameserver is the Schema for the gameservers API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=gameservers,scope=Namespaced
type Gameserver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GameserverSpec   `json:"spec,omitempty"`
	Status GameserverStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GameserverList contains a list of Gameserver
type GameserverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gameserver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Gameserver{}, &GameserverList{})
}
