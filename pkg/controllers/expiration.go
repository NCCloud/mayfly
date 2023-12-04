package controllers

import (
	"context"

	"github.com/NCCloud/mayfly/pkg/common"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type ExpirationController struct {
	config         *common.Config
	client         client.Client
	scheduler      *common.Scheduler
	apiVersionKind string
}

func NewExpirationController(config *common.Config, client client.Client,
	apiVersionKind string, scheduler *common.Scheduler,
) *ExpirationController {
	return &ExpirationController{
		config:         config,
		client:         client,
		scheduler:      scheduler,
		apiVersionKind: apiVersionKind,
	}
}

func (r *ExpirationController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		logger   = log.FromContext(ctx)
		resource = common.NewResourceInstance(r.apiVersionKind)
	)

	logger.Info("Reconciliation started.")

	if getErr := r.client.Get(ctx, req.NamespacedName, resource); getErr != nil {
		if errors.IsNotFound(getErr) {
			_ = r.scheduler.DeleteDeletionJob(resource)
		}

		return ctrl.Result{}, client.IgnoreNotFound(getErr)
	}

	hasExpired, expirationDate, hasExpiredErr := common.IsExpired(resource, r.config)
	if hasExpiredErr != nil {
		logger.Error(hasExpiredErr, "Error while checking if resource has expired.")

		return ctrl.Result{}, hasExpiredErr
	}

	if hasExpired {
		logger.Info("Resource already expired. Removing")

		_ = r.client.Delete(ctx, resource)
		_ = r.scheduler.DeleteDeletionJob(resource)

		return ctrl.Result{}, nil
	}

	if createOrUpdateErr := r.scheduler.CreateOrUpdateDeletionJob(expirationDate,
		func(resource client.Object) error {
			return r.client.Delete(context.Background(), resource)
		}, resource); createOrUpdateErr != nil {
		logger.Error(createOrUpdateErr, "Error while starting job.")

		return ctrl.Result{}, createOrUpdateErr
	}

	return ctrl.Result{}, nil
}

func (r *ExpirationController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(common.NewResourceInstance(r.apiVersionKind)).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(createEvent event.CreateEvent) bool {
				return len(createEvent.Object.GetAnnotations()[r.config.ExpirationLabel]) != 0
			},
			DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
				return len(deleteEvent.Object.GetAnnotations()[r.config.ExpirationLabel]) != 0
			},
			UpdateFunc: func(updateEvent event.UpdateEvent) bool {
				oldMayFlyAnnotation := updateEvent.ObjectOld.GetAnnotations()[r.config.ExpirationLabel]
				newMayFlyAnnotation := updateEvent.ObjectNew.GetAnnotations()[r.config.ExpirationLabel]

				return len(newMayFlyAnnotation) != 0 && oldMayFlyAnnotation != newMayFlyAnnotation
			},
		}).Complete(r)
}
