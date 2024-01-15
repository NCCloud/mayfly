package controllers

import (
	"context"
	"fmt"

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
		logger           = log.FromContext(ctx)
		resource         = common.NewResourceInstance(r.apiVersionKind)
		apiVersion, kind = resource.GroupVersionKind().ToAPIVersionAndKind()
		tag              = fmt.Sprintf("%s/%s/%s/%s/delete", apiVersion, kind, req.Name, req.Namespace)
	)

	logger.Info("Reconciliation started.")
	defer logger.Info("Reconciliation finished.")

	if getErr := r.client.Get(ctx, req.NamespacedName, resource); getErr != nil {
		if errors.IsNotFound(getErr) {
			_ = r.scheduler.DeleteTask(tag)
		}

		return ctrl.Result{}, client.IgnoreNotFound(getErr)
	}

	hasExpired, expirationDate, hasExpiredErr := common.IsExpired(resource, r.config)
	if hasExpiredErr != nil {
		logger.Error(hasExpiredErr, "Error while checking if resource has expired.")

		return ctrl.Result{}, hasExpiredErr
	}

	if hasExpired {
		logger.Info("Resource already expired will be removed.")

		_ = r.scheduler.DeleteTask(tag)

		return ctrl.Result{}, client.IgnoreNotFound(r.client.Delete(ctx, resource))
	}

	if createOrUpdateTaskErr := r.scheduler.CreateOrUpdateTask(tag, expirationDate, func() error {
		return r.client.Delete(context.Background(), resource)
	}); createOrUpdateTaskErr != nil {
		logger.Error(createOrUpdateTaskErr, "Error while creating or updating task.")

		return ctrl.Result{}, createOrUpdateTaskErr
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
