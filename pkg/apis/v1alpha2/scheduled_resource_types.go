package v1alpha2

import (
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

type (
	FailureReason string
	Condition     string
)

const (
	ConditionFinished  Condition = "Finished"
	ConditionScheduled Condition = "Scheduled"
	ConditionFailed    Condition = "Failed"
)

var ErrObjectIsNotValid = errors.New("object is not valid")

func init() {
	SchemeBuilder.Register(&ScheduledResource{}, &ScheduledResourceList{})
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Schedule",type=string,JSONPath=".spec.schedule"
// +kubebuilder:printcolumn:name="Next Run",type=string,JSONPath=".status.nextRun"
// +kubebuilder:printcolumn:name="Last Run",type=string,JSONPath=".status.lastRun"
// +kubebuilder:printcolumn:name="Condition",type=string,JSONPath=".status.condition"
// +kubebuilder:printcolumn:name="Completions",type=string,JSONPath=".status.completions"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"

type ScheduledResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScheduledResourceSpec   `json:"spec,omitempty"`
	Status ScheduledResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

type ScheduledResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScheduledResource `json:"items"`
}

type ScheduledResourceSpec struct {
	Schedule string `json:"schedule"`
	// +kubebuilder:validation:Minimum=1
	Completions int    `json:"completions,omitempty"`
	Content     string `json:"content"`
}

type ScheduledResourceStatus struct {
	NextRun     string    `json:"nextRun,omitempty"`
	LastRun     string    `json:"lastRun,omitempty"`
	Condition   Condition `json:"condition,omitempty"`
	Completions int       `json:"completions,omitempty"`
}

func (in *ScheduledResource) IsBeingDeleted() bool {
	return in.DeletionTimestamp != nil
}

func (in *ScheduledResource) GetContent() (*unstructured.Unstructured, error) {
	object, _, decodeErr := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).
		Decode([]byte(in.Spec.Content), nil, nil)
	if decodeErr != nil {
		return nil, decodeErr
	}

	unstructuredObj, isUnstructuredObj := object.(*unstructured.Unstructured)
	if !isUnstructuredObj {
		return nil, ErrObjectIsNotValid
	}

	return unstructuredObj, nil
}

func (in *ScheduledResource) IsCompletionsLimitReached(isOneTimeSchedule bool) bool {
	return isOneTimeSchedule && in.Status.Completions >= 1 ||
		in.Spec.Completions > 0 && in.Status.Completions >= in.Spec.Completions
}
