package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// APISpec Contains the Spec of the API object
type APISpec struct {
	Expose   bool   `json:"expose"` //TODO: Make expose readonly after creation
	Hostname string `json:"hostname"`
	// +patchMergeKey=name
	// +patchStrategy=merge,retainKeys
	Endpoints []Endpoint `json:"endpoints" patchStrategy:"merge" patchMergeKey:"name"`
	// +patchMergeKey=name
	// +patchStrategy=merge,retainKeys
	RateLimits []RateLimit `json:"rate_limits,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

type APIConditionType string

// APIStatus Contains the Status of the API object
type APIStatus struct { //TODO: Make this struct not user editable
	Deployed bool `json:"deployed"`

	// ObservedGeneration reflects the generation of the most recently observed ReplicaSet.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Represents the latest available observations of a replica set's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []APICondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

type APICondition struct {
	// Type of replica set condition.
	Type APIConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// Endpoint is a struct used to define the different upstream services
type Endpoint struct {
	Name       string      `json:"name"` // Not really needed?
	Host       string      `json:"host"`
	Path       string      `json:"path"`
	RateLimits []RateLimit `json:"rate_limits,omitempty"`
}

// RateLimit is a struct used to define different types of rate limiting rules
type RateLimit struct {
	Burst      *int       `json:"burst"`
	Conn       *int       `json:"conn"`
	Delay      *int       `json:"delay"`
	Limit      string     `json:"limit"`
	Name       string     `json:"name"`   //TODO - This needs to reference and endpoint name currently but this relationship will reverse.
	Source     string     `json:"source"` // Source will allow user to limit based on jwt, source ip etc
	Type       string     `json:"type"`
	Conditions *Condition `json:"conditions, omitempty"`
}

// Condition wraps a generic rate limit condition
type Condition struct {
	Operator   string               `json:"operator,omitempty"`
	Operations []RateLimitCondition `json:"operations"`
}

// RateLimitCondition is an interface for a type which should marshal to apicast config
type RateLimitCondition interface {
	MarshalJSON() ([]byte, error)
	DeepCopyRateLimitCondition() RateLimitCondition
}

// +k8s:deepcopy-gen:interfaces=github.com/3scale/ostia/ostia-operator/pkg/apis/ostia/v1alpha1.RateLimitCondition
type HeaderBasedCondition struct {
	Header    string `json:"header"`
	Operation string `json:"op, omitempty"`
	Value     string `json:"value"`
}

// +k8s:deepcopy-gen:interfaces=github.com/3scale/ostia/ostia-operator/pkg/apis/ostia/v1alpha1.RateLimitCondition
type MethodBasedCondition struct {
	Method    string `json:"http_method"`
	Operation string `json:"op, omitempty"`
}

// +k8s:deepcopy-gen:interfaces=github.com/3scale/ostia/ostia-operator/pkg/apis/ostia/v1alpha1.RateLimitCondition
type PathBasedCondition struct {
	Path      string `json:"request_path,omitempty"`
	Operation string `json:"op, omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// API is the Schema for the apis API
// +k8s:openapi-gen=true
type API struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APISpec   `json:"spec,omitempty"`
	Status APIStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIList contains a list of API
type APIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []API `json:"items"`
}

func init() {
	SchemeBuilder.Register(&API{}, &APIList{})
}
