package common

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func IsExpired(resource client.Object, config *Config) (bool, time.Time, error) {
	duration, parseDurationErr := time.ParseDuration(resource.GetAnnotations()[config.ExpirationLabel])
	if parseDurationErr != nil {
		return false, time.Time{}, parseDurationErr
	}

	creationTime := resource.GetCreationTimestamp()
	expirationDate := creationTime.Add(duration)

	if expirationDate.Before(time.Now()) {
		return true, expirationDate, nil
	}

	return false, expirationDate, nil
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
