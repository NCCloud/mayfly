package mayfly

import (
	"context"
	"fmt"

	"github.com/NCCloud/mayfly/pkg/common"
	"github.com/NCCloud/mayfly/pkg/controllers/mayfly/resource"
	"github.com/NCCloud/mayfly/pkg/controllers/mayfly/utils"
	"github.com/NCCloud/mayfly/pkg/scheduler"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type Controller struct {
	Resource  resource.Resource
	Config    *common.OperatorConfig
	Client    client.Client
	Scheduler *scheduler.Scheduler
}

func NewController(config *common.OperatorConfig, client client.Client, resource resource.Resource, scheduler *scheduler.Scheduler) *Controller {
	return &Controller{
		Resource:  resource,
		Config:    config,
		Client:    client,
		Scheduler: scheduler,
	}
}

func (r *Controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconcilliation started.")

	resource := r.Resource.NewResourceInstance()

	err := r.Client.Get(ctx, req.NamespacedName, resource)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Scheduler.RemoveJob(fmt.Sprintf("%v", resource.GetUID()))
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	hasExpired, expirationDate, hasExpiredErr := utils.HasExpired(resource, r.Config)
	if hasExpiredErr != nil {
		logger.Error(hasExpiredErr, "Error while checking if resource has expired.")
		return ctrl.Result{}, hasExpiredErr
	}

	if hasExpired {
		logger.Info("Resource already expired. Removing")
		_ = utils.DeleteResource(ctx, r.Client, resource)
	}

	startJobErr := r.Scheduler.StartOrUpdateJob(expirationDate, utils.DeleteResource, ctx, r.Client, resource)
	if startJobErr != nil {
		logger.Error(startJobErr, "Error while starting job.")
		return ctrl.Result{}, startJobErr
	}
	return ctrl.Result{}, nil
}

func (r *Controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(r.Resource.NewResourceInstance()).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(createEvent event.CreateEvent) bool {
				mayFlylabel := createEvent.Object.GetLabels()[r.Config.ResourceConfiguration.MayflyExpireLabel]
				return mayFlylabel != ""
			},
			DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
				mayFlylabel := deleteEvent.Object.GetLabels()[r.Config.ResourceConfiguration.MayflyExpireLabel]
				return mayFlylabel != ""
			},
			UpdateFunc: func(updateEvent event.UpdateEvent) bool {
				oldMayFlylabel := updateEvent.ObjectOld.GetLabels()[r.Config.ResourceConfiguration.MayflyExpireLabel]
				newMayFlylabel := updateEvent.ObjectNew.GetLabels()[r.Config.ResourceConfiguration.MayflyExpireLabel]
				if newMayFlylabel != "" && oldMayFlylabel != newMayFlylabel {
					return true
				}
				return false
			},
		}).Complete(r)
}
