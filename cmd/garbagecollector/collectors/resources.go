package collectors

import (
	"context"

	"github.com/NCCloud/mayfly/pkg/common"
	"github.com/NCCloud/mayfly/pkg/controllers/mayfly/resource"
	"github.com/NCCloud/mayfly/pkg/controllers/mayfly/utils"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ResourcesCollector struct {
	logger         logr.Logger
	mgrClient      client.Client
	resource       resource.Resource
	operatorConfig *common.OperatorConfig
}

func NewResourcesCollector(logger logr.Logger, mgrClient client.Client, resource resource.Resource, operatorConfig *common.OperatorConfig) Collector {
	return &ResourcesCollector{
		logger:         logger,
		mgrClient:      mgrClient,
		resource:       resource,
		operatorConfig: operatorConfig,
	}
}

func (e ResourcesCollector) Collect() error {
	resources := e.resource.NewResourceInstanceList()

	listOpts := []client.ListOption{}
	listOpts = append(listOpts, client.HasLabels{e.operatorConfig.ResourceConfiguration.MayflyExpireLabel})
	resourcesListErr := e.mgrClient.List(context.Background(), resources, listOpts...)

	if resourcesListErr != nil {
		e.logger.Error(resourcesListErr, "Failed to list resources")
		return resourcesListErr
	}

	for _, resource := range resources.Items {
		hasExpired, _, hasExpiredErr := utils.HasExpired(&resource, e.operatorConfig)
		if hasExpiredErr != nil {
			e.logger.Error(hasExpiredErr, "Failed to check if resource has expired", "kind", resource.GetKind(), "name", resource.GetName())
			continue
		}
		if hasExpired {
			e.logger.Info("Deleting resource", "kind", resource.GetKind(), "name", resource.GetName())
			_ = utils.DeleteResource(context.Background(), e.mgrClient, &resource)
		}
	}
	return nil
}
