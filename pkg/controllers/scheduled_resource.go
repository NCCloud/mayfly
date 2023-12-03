package controllers

import (
	"context"
	errors2 "errors"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/NCCloud/mayfly/pkg/apis/v1alpha1"
	"github.com/NCCloud/mayfly/pkg/common"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ScheduledResourceController struct {
	config    *common.Config
	client    client.Client
	scheduler *common.Scheduler
}

func NewScheduledResourceController(config *common.Config, client client.Client,
	scheduler *common.Scheduler,
) *ScheduledResourceController {
	return &ScheduledResourceController{
		config:    config,
		client:    client,
		scheduler: scheduler,
	}
}

func (r *ScheduledResourceController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		logger            = log.FromContext(ctx)
		scheduledResource = &v1alpha1.ScheduledResource{}
	)

	if getErr := r.client.Get(ctx, req.NamespacedName, scheduledResource); getErr != nil {
		if errors.IsNotFound(getErr) {
			_ = r.scheduler.RemoveCreationJob(scheduledResource)
		}

		return ctrl.Result{}, client.IgnoreNotFound(getErr)
	}

	logger.Info("Reconciliation started.")

	if scheduledResource.Status.Condition == v1alpha1.ConditionCreated {
		return ctrl.Result{}, nil
	}

	duration, parseDurationErr := time.ParseDuration(scheduledResource.Spec.In)
	if parseDurationErr != nil {
		return ctrl.Result{}, parseDurationErr
	}

	if createOrUpdateErr := r.scheduler.CreateOrUpdateCreationJob(scheduledResource.CreationTimestamp.Add(duration),
		func(resource v1alpha1.ScheduledResource) error {
			unstructured, toUnstructuredErr := resource.ToUnstructured()
			if toUnstructuredErr != nil {
				return toUnstructuredErr
			}

			if getErr := r.client.Get(context.Background(),
				client.ObjectKeyFromObject(&resource), &resource); getErr != nil {
				return getErr
			}

			if createErr := r.client.Create(context.Background(),
				unstructured); client.IgnoreAlreadyExists(createErr) != nil {
				resource.Status.Condition = v1alpha1.ConditionFailed

				return errors2.Join(createErr, r.client.Status().Update(context.Background(), &resource))
			}

			logger.Info(fmt.Sprintf("%s/%s created", resource.Name, resource.Namespace))

			if removeErr := r.scheduler.RemoveCreationJob(&resource); removeErr != nil {
				return removeErr
			}

			resource.Status.Condition = v1alpha1.ConditionCreated

			return r.client.Status().Update(context.Background(), &resource)
		}, *scheduledResource); createOrUpdateErr != nil {
		return ctrl.Result{}, createOrUpdateErr
	}

	logger.Info("Reconciliation finished.")

	scheduledResource.Status.Condition = v1alpha1.ConditionScheduled

	return ctrl.Result{}, r.client.Status().Update(ctx, scheduledResource)
}

func (r *ScheduledResourceController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ScheduledResource{}).
		Complete(r)
}
