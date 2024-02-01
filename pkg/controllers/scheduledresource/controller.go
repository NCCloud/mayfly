package scheduledresource

import (
	"context"
	errors2 "errors"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/NCCloud/mayfly/pkg/apis/v1alpha1"
	"github.com/NCCloud/mayfly/pkg/common"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Controller struct {
	config    *common.Config
	client    client.Client
	scheduler common.Scheduler
}

func NewController(config *common.Config, client client.Client,
	scheduler common.Scheduler,
) *Controller {
	return &Controller{
		config:    config,
		client:    client,
		scheduler: scheduler,
	}
}

func (r *Controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		logger            = log.FromContext(ctx)
		scheduledResource = &v1alpha1.ScheduledResource{}
		tag               = fmt.Sprintf("v1alpha1/ScheduledResource/%s/%s/create", req.Name, req.Namespace)
	)

	logger.Info("Reconciliation started.")
	defer logger.Info("Reconciliation finished.")

	if getErr := r.client.Get(ctx, req.NamespacedName, scheduledResource); getErr != nil {
		if errors.IsNotFound(getErr) {
			_ = r.scheduler.DeleteTask(tag)
		}

		return ctrl.Result{}, client.IgnoreNotFound(getErr)
	}

	if scheduledResource.Status.Condition == v1alpha1.ConditionCreated {
		return ctrl.Result{}, nil
	}

	schedule, scheduleErr := common.ResolveSchedule(scheduledResource.CreationTimestamp, scheduledResource.Spec.In)
	if scheduleErr != nil {
		return ctrl.Result{}, scheduleErr
	}

	if createOrUpdateTaskErr := r.scheduler.CreateOrUpdateTask(tag, schedule, func() error {
		content, contentErr := scheduledResource.GetContent()
		if contentErr != nil {
			logger.Error(contentErr, "Error while parsing content.")

			return contentErr
		}

		if getErr := r.client.Get(context.Background(), client.
			ObjectKeyFromObject(scheduledResource), scheduledResource); client.IgnoreNotFound(getErr) != nil {
			logger.Error(contentErr, "Error while getting resource.")

			return getErr
		}

		if createErr := r.client.Create(context.Background(),
			content); client.IgnoreAlreadyExists(createErr) != nil {
			logger.Error(contentErr, "An error occurred while creating resource.")

			scheduledResource.Status.Condition = v1alpha1.ConditionFailed

			return errors2.Join(createErr, r.client.Status().Update(context.Background(), scheduledResource))
		}

		logger.Info(fmt.Sprintf("%s created.", tag))

		_ = r.scheduler.DeleteTask(tag)

		scheduledResource.Status.Condition = v1alpha1.ConditionCreated

		return r.client.Status().Update(context.Background(), scheduledResource)
	}); createOrUpdateTaskErr != nil {
		logger.Error(createOrUpdateTaskErr, "Error while creating or updating task.")

		return ctrl.Result{}, createOrUpdateTaskErr
	}

	scheduledResource.Status.Condition = v1alpha1.ConditionScheduled

	return ctrl.Result{}, r.client.Status().Update(ctx, scheduledResource)
}

func (r *Controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ScheduledResource{}).
		Complete(r)
}
