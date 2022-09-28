package utils

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DeleteResource(ctx context.Context, client client.Client, resource client.Object) error {
	deleteErr := client.Delete(ctx, resource)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}
