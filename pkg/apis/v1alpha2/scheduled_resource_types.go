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

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Schedule",type=string,JSONPath=".spec.schedule"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="Condition",type=string,JSONPath=".status.condition"

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
	// +kubebuilder:validation:Immutable=true
	Schedule string `json:"schedule"`
	Content  string `json:"content"`
}

type ScheduledResourceStatus struct {
	Condition Condition `json:"condition"`
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
