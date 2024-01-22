package common

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func ResolveSchedule(creationTimestamp metav1.Time, dateTimeDuration string) (time.Time, error) {
	duration, parseDurationErr := time.ParseDuration(dateTimeDuration)
	if parseDurationErr == nil {
		return creationTimestamp.Add(duration), nil
	}

	date, parseDateErr := dateparse.ParseAny(dateTimeDuration)
	if parseDateErr == nil {
		return date, nil
	}

	return time.Time{}, errors.Join(parseDurationErr, parseDateErr)
}

func NewResourceInstance(apiVersionKind string) *unstructured.Unstructured {
	apiVersionKindArr := strings.Split(apiVersionKind, ";")

	return &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": apiVersionKindArr[0],
			"kind":       apiVersionKindArr[1],
		},
	}
}

func NewResourceInstanceList(apiVersionKind string) *unstructured.UnstructuredList {
	resourceInstance := NewResourceInstance(apiVersionKind)

	return &unstructured.UnstructuredList{
		Object: map[string]any{
			"apiVersion": resourceInstance.GetAPIVersion(),
			"kind":       fmt.Sprintf("%sList", resourceInstance.GetKind()),
		},
	}
}
