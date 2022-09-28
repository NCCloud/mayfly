package mayfly

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"mayfly/pkg/common"
	"mayfly/pkg/controllers/mayfly/resource"
	"mayfly/pkg/controllers/mayfly/utils"
	"mayfly/pkg/scheduler"
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

	duration, parseDurationErr := time.ParseDuration(resource.GetAnnotations()[r.Config.ResourceConfiguration.MayflyExpireAnnotation])
	if parseDurationErr != nil {
		logger.Error(parseDurationErr, "Failed to parse duration.")
		return ctrl.Result{}, parseDurationErr
	}

	creationTime := resource.GetCreationTimestamp()
	expirationDate := creationTime.Add(duration)

	if expirationDate.Before(time.Now()) {
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
				annotation := createEvent.Object.GetAnnotations()[r.Config.ResourceConfiguration.MayflyExpireAnnotation]
				return annotation != ""
			},
			DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
				annotation := deleteEvent.Object.GetAnnotations()[r.Config.ResourceConfiguration.MayflyExpireAnnotation]
				return annotation != ""
			},
			UpdateFunc: func(updateEvent event.UpdateEvent) bool {
				oldAnnotation := updateEvent.ObjectOld.GetAnnotations()[r.Config.ResourceConfiguration.MayflyExpireAnnotation]
				newAnnotation := updateEvent.ObjectNew.GetAnnotations()[r.Config.ResourceConfiguration.MayflyExpireAnnotation]
				if newAnnotation != "" && oldAnnotation != newAnnotation {
					return true
				}
				return false
			},
		}).Complete(r)
}
