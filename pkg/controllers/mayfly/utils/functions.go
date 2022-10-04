package utils

import (
	"context"
	"time"

	"github.com/NCCloud/mayfly/pkg/common"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DeleteResource(ctx context.Context, client client.Client, resource client.Object) error {
	deleteErr := client.Delete(ctx, resource)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

func HasExpired(resource client.Object, config *common.OperatorConfig) (bool, time.Time, error) {
	duration, parseDurationErr := time.ParseDuration(resource.GetAnnotations()[config.ResourceConfiguration.MayflyExpireLabel])
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
