package scheduled_resource

import (
	"context"
	errors2 "errors"
	"fmt"
	common2 "github.com/NCCloud/mayfly/mocks/github.com/NCCloud/mayfly/pkg/common"
	cache2 "github.com/NCCloud/mayfly/mocks/sigs.k8s.io/controller-runtime/pkg/cache"
	client2 "github.com/NCCloud/mayfly/mocks/sigs.k8s.io/controller-runtime/pkg/client"
	manager2 "github.com/NCCloud/mayfly/mocks/sigs.k8s.io/controller-runtime/pkg/manager"
	"github.com/NCCloud/mayfly/pkg/apis/v1alpha1"
	"github.com/NCCloud/mayfly/pkg/common"
	"github.com/araddon/dateparse"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/config"
	controller2 "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"strings"
	"testing"
	"time"
)

var testVars = struct {
	config         *common.Config
	k8sClient      client.Client
	apiVersionKind string
}{
	config:         common.NewConfig(),
	apiVersionKind: "v1;Secret",
}

func init() {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))

	kubeConfig, testEnvStartErr := (&envtest.Environment{
		ControlPlane: envtest.ControlPlane{
			APIServer: &envtest.APIServer{
				StartTimeout: 5 * time.Minute,
				StopTimeout:  5 * time.Minute,
			},
			Etcd: &envtest.Etcd{
				StartTimeout: 5 * time.Minute,
				StopTimeout:  5 * time.Minute,
			},
		},
		ErrorIfCRDPathMissing: true,
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "deploy", "crds"),
			filepath.Join("..", "..", "..", ".envtest", "crds"),
		},
		BinaryAssetsDirectory:    "../../../.envtest/bins",
		ControlPlaneStartTimeout: 5 * time.Minute,
		ControlPlaneStopTimeout:  5 * time.Minute,
	}).Start()
	if testEnvStartErr != nil {
		panic(testEnvStartErr)
	}

	manager, managerErr := ctrl.NewManager(kubeConfig, ctrl.Options{
		Scheme: scheme,
		Metrics: server.Options{
			BindAddress: ":0",
		},
		Logger: zap.New(),
	})
	if managerErr != nil {
		panic(managerErr)
	}

	if setupErr := (&Controller{
		config:    common.NewConfig(),
		client:    manager.GetClient(),
		scheduler: common.NewScheduler(testVars.config),
	}).SetupWithManager(manager); setupErr != nil {
		panic(setupErr)
	}

	testVars.k8sClient = manager.GetClient()

	go func() {
		log.SetLogger(logr.New(log.NullLogSink{}))

		if managerStartErr := manager.Start(context.Background()); managerStartErr != nil {
			panic(managerStartErr)
		}
	}()
}

func TestController_New(t *testing.T) {
	// given
	var (
		config        = common.NewConfig()
		mockClient    = new(client2.MockClient)
		mockScheduler = new(common2.MockScheduler)
	)

	// when
	controller := NewController(config, mockClient, mockScheduler)

	// then
	assert.NotNil(t, controller)
	assert.IsType(t, controller, &Controller{})
}

func TestController_Reconcile(t *testing.T) {
	// given
	var (
		ctx               = context.Background()
		mockClient        = new(client2.MockClient)
		mockStatusClient  = new(client2.MockSubResourceClient)
		mockScheduler     = new(common2.MockScheduler)
		controller        = NewController(common.NewConfig(), mockClient, mockScheduler)
		scheduledResource = &v1alpha1.ScheduledResource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      gofakeit.Name(),
				Namespace: gofakeit.Name(),
			},
			Spec: v1alpha1.ScheduledResourceSpec{
				In: gofakeit.FutureDate().String(),
			},
		}
	)

	mockScheduler.EXPECT().CreateOrUpdateTask(mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockStatusClient.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	mockClient.EXPECT().Status().Return(mockStatusClient)
	mockClient.EXPECT().Get(mock.Anything, client.ObjectKeyFromObject(scheduledResource),
		mock.AnythingOfType("*v1alpha1.ScheduledResource")).RunAndReturn(
		func(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
			scheduledResource.DeepCopyInto(obj.(*v1alpha1.ScheduledResource))
			return nil
		})

	// when
	result, reconcileErr := controller.Reconcile(ctx, ctrl.Request{
		NamespacedName: client.ObjectKeyFromObject(scheduledResource),
	})

	// then
	date, _ := dateparse.ParseAny(scheduledResource.Spec.In)
	mockScheduler.AssertCalled(t, "CreateOrUpdateTask", fmt.Sprintf("v1alpha1/ScheduledResource/%s/%s/create",
		scheduledResource.Name, scheduledResource.Namespace), date, mock.Anything)
	mockStatusClient.AssertCalled(t, "Update", mock.Anything, mock.MatchedBy(func(obj client.Object) bool {
		return obj.(*v1alpha1.ScheduledResource).Status.Condition == v1alpha1.ConditionScheduled
	}))
	assert.Nil(t, reconcileErr)
	assert.False(t, result.Requeue)
}

func TestController_ReconcileIntegration(t *testing.T) {
	// given
	var (
		ctx               = context.Background()
		scheduledResource = &v1alpha1.ScheduledResource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      strings.ToLower(strings.ReplaceAll(gofakeit.Name(), " ", "")),
				Namespace: "default",
			},
			Spec: v1alpha1.ScheduledResourceSpec{
				In: "5s",
				Content: `apiVersion: v1
kind: Secret
metadata:
  name: my-resource
  namespace: default`,
			},
		}
	)

	// when
	content, contentErr := scheduledResource.GetContent()
	createErr := testVars.k8sClient.Create(ctx, scheduledResource)

	// then
	assert.Nil(t, contentErr)
	assert.Nil(t, createErr)
	assert.Eventually(t, func() bool {
		return errors.IsNotFound(testVars.k8sClient.Get(ctx, client.ObjectKeyFromObject(content), content))
	}, 60*time.Second, 100*time.Millisecond)
	assert.Eventually(t, func() bool {
		return testVars.k8sClient.Get(ctx, client.ObjectKeyFromObject(content), content) == nil
	}, 60*time.Second, 100*time.Millisecond)
}

func TestController_Reconcile_ShouldDeleteTaskWhenNotFound(t *testing.T) {
	// given
	var (
		ctx           = context.Background()
		mockClient    = new(client2.MockClient)
		mockScheduler = new(common2.MockScheduler)
		controller    = NewController(common.NewConfig(), mockClient, mockScheduler)
	)

	mockScheduler.EXPECT().DeleteTask(mock.Anything).Return(nil)
	mockClient.EXPECT().Get(mock.Anything, mock.Anything,
		mock.AnythingOfType("*v1alpha1.ScheduledResource")).Return(errors.NewNotFound(schema.GroupResource{}, ""))

	// when
	result, reconcileErr := controller.Reconcile(ctx, ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "my-secret",
			Namespace: "my-namespace",
		},
	})

	// then
	mockScheduler.AssertCalled(t, "DeleteTask", fmt.Sprintf("v1alpha1/ScheduledResource/%s/%s/create",
		"my-secret", "my-namespace"))
	assert.Nil(t, reconcileErr)
	assert.False(t, result.Requeue)
}

func TestController_Reconcile_ShouldReturnErrWhenInFieldDoesNotMakeSense(t *testing.T) {
	// given
	var (
		ctx               = context.Background()
		mockClient        = new(client2.MockClient)
		mockScheduler     = new(common2.MockScheduler)
		controller        = NewController(common.NewConfig(), mockClient, mockScheduler)
		scheduledResource = &v1alpha1.ScheduledResource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      gofakeit.Name(),
				Namespace: gofakeit.Name(),
			},
			Spec: v1alpha1.ScheduledResourceSpec{
				In: gofakeit.BeerName(),
			},
		}
	)

	mockClient.EXPECT().Get(mock.Anything, client.ObjectKeyFromObject(scheduledResource),
		mock.AnythingOfType("*v1alpha1.ScheduledResource")).RunAndReturn(
		func(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
			scheduledResource.DeepCopyInto(obj.(*v1alpha1.ScheduledResource))
			return nil
		})

	// when
	result, reconcileErr := controller.Reconcile(ctx, ctrl.Request{
		NamespacedName: client.ObjectKeyFromObject(scheduledResource),
	})

	// then
	assert.NotNil(t, reconcileErr)
	assert.False(t, result.Requeue)
}

func TestController_Reconcile_TaskShouldCreateItem(t *testing.T) {
	// given
	var (
		ctx               = context.Background()
		mockClient        = new(client2.MockClient)
		mockStatusClient  = new(client2.MockSubResourceClient)
		mockScheduler     = new(common2.MockScheduler)
		controller        = NewController(common.NewConfig(), mockClient, mockScheduler)
		scheduledResource = &v1alpha1.ScheduledResource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      gofakeit.Name(),
				Namespace: gofakeit.Name(),
			},
			Spec: v1alpha1.ScheduledResourceSpec{
				In: gofakeit.FutureDate().String(),
				Content: `apiVersion: v1
kind: Secret
metadata:
  name: my-resource
  namespace: default`,
			},
		}
	)

	mockScheduler.EXPECT().DeleteTask(mock.Anything).Return(nil)
	mockScheduler.EXPECT().CreateOrUpdateTask(mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockStatusClient.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	mockClient.EXPECT().Status().Return(mockStatusClient)
	mockClient.EXPECT().Get(mock.Anything, client.ObjectKeyFromObject(scheduledResource),
		mock.AnythingOfType("*v1alpha1.ScheduledResource")).RunAndReturn(
		func(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
			scheduledResource.DeepCopyInto(obj.(*v1alpha1.ScheduledResource))
			return nil
		})
	mockClient.EXPECT().Create(mock.Anything, mock.Anything).Return(nil)
	_, reconcileErr := controller.Reconcile(ctx, ctrl.Request{
		NamespacedName: client.ObjectKeyFromObject(scheduledResource),
	})

	// when
	taskErr := mockScheduler.Calls[0].Arguments[2].(func() error)()

	// then
	mockScheduler.AssertCalled(t, "DeleteTask", fmt.Sprintf("v1alpha1/ScheduledResource/%s/%s/create",
		scheduledResource.Name, scheduledResource.Namespace))
	mockClient.AssertCalled(t, "Create", mock.Anything, mock.Anything)
	assert.Nil(t, reconcileErr)
	assert.Nil(t, taskErr)
}

func TestController_Reconcile_ShouldFailedIfAnyErrorHappens(t *testing.T) {
	// given
	var (
		ctx               = context.Background()
		mockClient        = new(client2.MockClient)
		mockStatusClient  = new(client2.MockSubResourceClient)
		mockScheduler     = new(common2.MockScheduler)
		controller        = NewController(common.NewConfig(), mockClient, mockScheduler)
		scheduledResource = &v1alpha1.ScheduledResource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      gofakeit.Name(),
				Namespace: gofakeit.Name(),
			},
			Spec: v1alpha1.ScheduledResourceSpec{
				In: gofakeit.FutureDate().String(),
				Content: `apiVersion: v1
kind: Secret
metadata:
  name: my-resource
  namespace: default`,
			},
		}
	)

	mockScheduler.EXPECT().DeleteTask(mock.Anything).Return(nil)
	mockScheduler.EXPECT().CreateOrUpdateTask(mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockStatusClient.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	mockClient.EXPECT().Status().Return(mockStatusClient)
	mockClient.EXPECT().Get(mock.Anything, client.ObjectKeyFromObject(scheduledResource),
		mock.AnythingOfType("*v1alpha1.ScheduledResource")).RunAndReturn(
		func(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
			scheduledResource.DeepCopyInto(obj.(*v1alpha1.ScheduledResource))
			return nil
		})
	mockClient.EXPECT().Create(mock.Anything, mock.Anything).Return(errors2.New("an error"))
	_, reconcileErr := controller.Reconcile(ctx, ctrl.Request{
		NamespacedName: client.ObjectKeyFromObject(scheduledResource),
	})

	// when
	taskErr := mockScheduler.Calls[0].Arguments[2].(func() error)()

	// then
	mockStatusClient.AssertCalled(t, "Update", mock.Anything, mock.MatchedBy(func(obj client.Object) bool {
		return obj.(*v1alpha1.ScheduledResource).Status.Condition == v1alpha1.ConditionFailed
	}))
	assert.Nil(t, reconcileErr)
	assert.NotNil(t, taskErr)
}

func TestController_SetupWithManager(t *testing.T) {
	// given
	var (
		mockClient    = new(client2.MockClient)
		mockManager   = new(manager2.MockManager)
		mockCache     = new(cache2.MockCache)
		mockScheduler = new(common2.MockScheduler)
		controller    = NewController(common.NewConfig(), mockClient, mockScheduler)
		scheme        = runtime.NewScheme()
	)

	addToSchemeErr := v1alpha1.AddToScheme(scheme)
	mockManager.EXPECT().GetControllerOptions().Return(config.Controller{})
	mockManager.EXPECT().GetScheme().Return(scheme)
	mockManager.EXPECT().GetCache().Return(mockCache)
	mockManager.EXPECT().GetRESTMapper().Return(meta.MultiRESTMapper{})
	mockManager.EXPECT().GetLogger().Return(zap.New())
	mockManager.EXPECT().GetFieldIndexer().Return(mockCache)
	mockManager.EXPECT().Add(mock.MatchedBy(func(ct controller2.Controller) bool {
		return ct != nil
	})).Return(nil)

	// when
	setupErr := controller.SetupWithManager(mockManager)

	// then
	assert.Nil(t, addToSchemeErr)
	assert.Nil(t, setupErr)
}
