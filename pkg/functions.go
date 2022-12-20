package pkg

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func IsExpired(resource client.Object, config *Config) (bool, time.Time, error) {
	var expirationDate time.Time

	hasAnnotation, annotationLabel, annotationValue := HasMayFlyAnnotation(resource, config)

	if hasAnnotation {
		if annotationLabel == config.ExpirationLabel {
			duration, parseDurationErr := time.ParseDuration(annotationValue)
			if parseDurationErr != nil {
				return false, time.Time{}, parseDurationErr
			}
			creationTime := resource.GetCreationTimestamp()
			expirationDate = creationTime.Add(duration)
		} else if annotationLabel == config.ExpirationDateLabel {
			var parseTimeErr error
			expirationDate, parseTimeErr = time.Parse(time.RFC3339, annotationValue)
			if parseTimeErr != nil {
				return false, time.Time{}, parseTimeErr
			}
		}
	}

	return expirationDate.Before(time.Now()), expirationDate, nil
}

func NewResourceInstance(apiVersionKind string) *unstructured.Unstructured {
	apiVersionKindArr := strings.Split(apiVersionKind, ";")

	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersionKindArr[0],
			"kind":       apiVersionKindArr[1],
		},
	}
}

func NewResourceInstanceList(apiVersionKind string) *unstructured.UnstructuredList {
	resourceInstance := NewResourceInstance(apiVersionKind)

	return &unstructured.UnstructuredList{
		Object: map[string]interface{}{
			"apiVersion": resourceInstance.GetAPIVersion(),
			"kind":       fmt.Sprintf("%sList", resourceInstance.GetKind()),
		},
	}
}

func HasMayFlyAnnotation(resource client.Object, config *Config) (bool, string, string) {
	var annotation string
	var value string
	if resource.GetAnnotations()[config.ExpirationLabel] != "" {
		annotation = config.ExpirationLabel
		value = resource.GetAnnotations()[config.ExpirationLabel]
	} else if resource.GetAnnotations()[config.ExpirationDateLabel] != "" {
		annotation = config.ExpirationDateLabel
		value = resource.GetAnnotations()[config.ExpirationDateLabel]
	}
	return value != "", annotation, value
}
