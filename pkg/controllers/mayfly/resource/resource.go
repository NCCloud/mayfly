package resource

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r Resource) NewResourceInstance() client.Object {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": r.APIVersion,
			"kind":       r.Kind,
		},
	}
}

func (r Resource) NewResourceInstanceList() *unstructured.UnstructuredList {
	return &unstructured.UnstructuredList{
		Object: map[string]interface{}{
			"apiVersion": r.APIVersion,
			"kind":       r.Kind + "List",
		},
		Items: []unstructured.Unstructured{},
	}
}
