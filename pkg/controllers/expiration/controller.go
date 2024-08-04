package expiration

import (
	"context"
	"errors"
	"fmt"
	"github.com/NCCloud/mayfly/pkg/common"
	"github.com/go-co-op/gocron/v2"
	errors2 "k8s.io/apimachinery/pkg/api/errors"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Controller struct {
	config         *common.Config
	client         client.Client
	scheduler      common.Scheduler
	apiVersionKind string
}

func NewController(config *common.Config, client client.Client,
	apiVersionKind string, scheduler common.Scheduler,
) *Controller {
	return &Controller{
		config:         config,
		client:         client,
		scheduler:      scheduler,
		apiVersionKind: apiVersionKind,
	}
}

func (r *Controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		logger           = log.FromContext(ctx)
		resource         = common.NewResourceInstance(r.apiVersionKind)
		apiVersion, kind = resource.GroupVersionKind().ToAPIVersionAndKind()
		tag              = fmt.Sprintf("%s/%s/%s/%s/delete", apiVersion, kind, req.Name, req.Namespace)
	)

	if getErr := r.client.Get(ctx, req.NamespacedName, resource); getErr != nil {
		if errors2.IsNotFound(getErr) {
			return ctrl.Result{}, r.scheduler.DeleteTask(tag)
		}

		return ctrl.Result{}, client.IgnoreNotFound(getErr)
	}

	expiration, hasExpiration := resource.GetAnnotations()[r.config.ExpirationLabel]
	if !hasExpiration {
		return ctrl.Result{}, r.scheduler.DeleteTask(tag)
	}

	date, dateErr := common.ResolveOneTimeSchedule(resource.GetCreationTimestamp(), expiration)
	if dateErr != nil {
		return ctrl.Result{}, dateErr
	}

	createOrUpdateTaskErr := r.scheduler.CreateOrUpdateOneTimeTask(tag, date, func() error {
		logger.Info("Deleted")

		return client.IgnoreNotFound(r.client.Delete(ctx, resource))
	})

	if errors.Is(createOrUpdateTaskErr, gocron.ErrOneTimeJobStartDateTimePast) {
		logger.Info("Deleted")

		return ctrl.Result{}, client.IgnoreNotFound(r.client.Delete(ctx, resource))
	}

	logger.Info("Scheduled")

	return ctrl.Result{}, createOrUpdateTaskErr
}

func (r *Controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(common.NewResourceInstance(r.apiVersionKind)).Complete(r)
}
