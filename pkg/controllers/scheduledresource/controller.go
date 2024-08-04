package scheduledresource

import (
	"context"
	errors2 "errors"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/NCCloud/mayfly/pkg/apis/v1alpha2"
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
		scheduledResource = &v1alpha2.ScheduledResource{}
		tag               = fmt.Sprintf("v1alpha2/ScheduledResource/%s/%s/create", req.Name, req.Namespace)
	)

	logger.Info("Reconciliation started.")
	defer logger.Info("Reconciliation finished.")

	if getErr := r.client.Get(ctx, req.NamespacedName, scheduledResource); getErr != nil {
		if errors.IsNotFound(getErr) {
			_ = r.scheduler.DeleteTask(tag)
		}

		return ctrl.Result{}, client.IgnoreNotFound(getErr)
	}

	oneTimeSchedule, oneTimeScheduleErr := common.ResolveOneTimeSchedule(
		scheduledResource.CreationTimestamp, scheduledResource.Spec.Schedule)
	isOneTimeSchedule := oneTimeScheduleErr == nil

	task := func() error {
		if scheduledResource.Status.Condition == v1alpha2.ConditionFinished {
			return nil
		}

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

			scheduledResource.Status.Condition = v1alpha2.ConditionFailed

			return errors2.Join(createErr, r.client.Status().Update(context.Background(), scheduledResource))
		}

		logger.Info(fmt.Sprintf("%s created.", tag))
		if isOneTimeSchedule {
			scheduledResource.Status.Condition = v1alpha2.ConditionFinished
		}

		return r.client.Status().Update(context.Background(), scheduledResource)
	}

	if isOneTimeSchedule {

		if createOrUpdateTaskErr := r.scheduler.CreateOrUpdateOneTimeTask(
			tag, oneTimeSchedule, task); createOrUpdateTaskErr != nil {
			return ctrl.Result{}, createOrUpdateTaskErr
		}
	} else {
		if createOrUpdateTaskErr := r.scheduler.CreateOrUpdateRecurringTask(
			tag, scheduledResource.Spec.Schedule, task); createOrUpdateTaskErr != nil {
			return ctrl.Result{}, createOrUpdateTaskErr
		}
	}

	scheduledResource.Status.Condition = v1alpha2.ConditionScheduled

	return ctrl.Result{}, r.client.Status().Update(ctx, scheduledResource)
}

func (r *Controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha2.ScheduledResource{}).
		Complete(r)
}
