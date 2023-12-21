package v1alpha1

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
	ConditionCreated   Condition = "Created"
	ConditionScheduled Condition = "Scheduled"
	ConditionFailed    Condition = "Failed"
)

var ErrObjectIsNotValid = errors.New("object is not valid")

func init() {
	SchemeBuilder.Register(&ScheduledResource{}, &ScheduledResourceList{})
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"

type ScheduledResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScheduledSpec   `json:"spec,omitempty"`
	Status ScheduledStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

type ScheduledResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScheduledResource `json:"items"`
}

type ScheduledSpec struct {
	In      string `json:"in"`
	Content string `json:"content"`
}

type ScheduledStatus struct {
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
