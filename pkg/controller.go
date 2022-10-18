package pkg

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type Controller struct {
	APIVersionKind string
	Config         *Config
	MgrClient      client.Client
	Scheduler      *Scheduler
}

func NewController(config *Config, client client.Client, apiVersionKind string, scheduler *Scheduler) *Controller {
	return &Controller{
		APIVersionKind: apiVersionKind,
		Config:         config,
		MgrClient:      client,
		Scheduler:      scheduler,
	}
}

func (r *Controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciliation started.")

	resource := NewResourceInstance(r.APIVersionKind)

	if getErr := r.MgrClient.Get(ctx, req.NamespacedName, resource); getErr != nil {
		if errors.IsNotFound(getErr) {
			r.Scheduler.RemoveJob(fmt.Sprintf("%v", resource.GetUID()))

			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, getErr
	}

	hasExpired, expirationDate, hasExpiredErr := IsExpired(resource, r.Config)
	if hasExpiredErr != nil {
		logger.Error(hasExpiredErr, "Error while checking if resource has expired.")

		return ctrl.Result{}, hasExpiredErr
	}

	if hasExpired {
		logger.Info("Resource already expired. Removing")

		_ = r.MgrClient.Delete(ctx, resource)
	}

	startJobErr := r.Scheduler.StartOrUpdateJob(ctx, expirationDate, func(ctx context.Context, client client.Client,
		resource client.Object,
	) error {
		return r.MgrClient.Delete(ctx, resource)
	}, r.MgrClient, resource)
	if startJobErr != nil {
		logger.Error(startJobErr, "Error while starting job.")

		return ctrl.Result{}, startJobErr
	}

	return ctrl.Result{}, nil
}

func (r *Controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(NewResourceInstance(r.APIVersionKind)).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(createEvent event.CreateEvent) bool {
				mayFlyAnnotation := createEvent.Object.GetAnnotations()[r.Config.ExpirationLabel]

				return mayFlyAnnotation != ""
			},
			DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
				mayFlyAnnotation := deleteEvent.Object.GetAnnotations()[r.Config.ExpirationLabel]

				return mayFlyAnnotation != ""
			},
			UpdateFunc: func(updateEvent event.UpdateEvent) bool {
				oldMayFlyAnnotation := updateEvent.ObjectOld.GetAnnotations()[r.Config.ExpirationLabel]
				newMayFlyAnnotation := updateEvent.ObjectNew.GetAnnotations()[r.Config.ExpirationLabel]
				if newMayFlyAnnotation != "" && oldMayFlyAnnotation != newMayFlyAnnotation {
					return true
				}

				return false
			},
		}).Complete(r)
}
