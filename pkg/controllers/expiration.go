package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/go-co-op/gocron/v2"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/NCCloud/mayfly/pkg/common"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
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

	if getErr := r.client.Get(ctx, req.NamespacedName, resource); getErr != nil {
		if errors2.IsNotFound(getErr) {
			return ctrl.Result{}, r.scheduler.DeleteTask(tag)
		}

		return ctrl.Result{}, client.IgnoreNotFound(getErr)
	}

	annotation, hasAnnotation := resource.GetAnnotations()[r.config.ExpirationLabel]
	if !hasAnnotation {
		return ctrl.Result{}, r.scheduler.DeleteTask(tag)
	}

	expiration, expirationErr := r.ResolveExpiration(resource.GetCreationTimestamp(), annotation)
	if expirationErr != nil {
		return ctrl.Result{}, expirationErr
	}

	createOrUpdateTaskErr := r.scheduler.CreateOrUpdateTask(tag, expiration, func() error {
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

func (r *ExpirationController) ResolveExpiration(creationTimestamp metav1.Time, expiration string) (time.Time, error) {
	duration, parseDurationErr := time.ParseDuration(expiration)
	if parseDurationErr == nil {
		return creationTimestamp.Add(duration), nil
	}

	date, parseDateErr := dateparse.ParseAny(expiration)
	if parseDateErr == nil {
		return date, nil
	}

	return time.Time{}, errors.Join(parseDurationErr, parseDateErr)
}

func (r *ExpirationController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(common.NewResourceInstance(r.apiVersionKind)).Complete(r)
}
