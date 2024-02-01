package expiration

import (
	"context"
	common2 "github.com/NCCloud/mayfly/mocks/github.com/NCCloud/mayfly/pkg/common"
	cache2 "github.com/NCCloud/mayfly/mocks/sigs.k8s.io/controller-runtime/pkg/cache"
	client2 "github.com/NCCloud/mayfly/mocks/sigs.k8s.io/controller-runtime/pkg/client"
	manager2 "github.com/NCCloud/mayfly/mocks/sigs.k8s.io/controller-runtime/pkg/manager"
	"github.com/NCCloud/mayfly/pkg/apis/v1alpha1"
	"github.com/araddon/dateparse"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

	"github.com/NCCloud/mayfly/pkg/common"
	"github.com/stretchr/testify/assert"
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
		config:         common.NewConfig(),
		client:         manager.GetClient(),
		scheduler:      common.NewScheduler(testVars.config),
		apiVersionKind: testVars.apiVersionKind,
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
		config         = common.NewConfig()
		mockClient     = new(client2.MockClient)
		mockScheduler  = new(common2.MockScheduler)
		apiVersionKind = "v1;Secret"
	)

	// when
	controller := NewController(config, mockClient, apiVersionKind, mockScheduler)

	// then
	assert.NotNil(t, controller)
	assert.IsType(t, controller, &Controller{})
}

func TestController_Reconcile(t *testing.T) {
	// given
	var (
		ctx            = context.Background()
		config         = common.NewConfig()
		mockClient     = new(client2.MockClient)
		mockScheduler  = new(common2.MockScheduler)
		apiVersionKind = "v1;Secret"
		controller     = NewController(config, mockClient, apiVersionKind, mockScheduler)
		secret         = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion":     "v1",
				"apiVersionKind": "Secret",
				"type":           "Opaque",
				"metadata": map[string]interface{}{
					"name":      "my-secret",
					"namespace": "my-namespace",
					"annotations": map[string]interface{}{
						config.ExpirationLabel: gofakeit.FutureDate().String(),
					},
				},
				"data": map[string]interface{}{
					"my-key": "my-value",
				},
			},
		}
	)

	mockScheduler.EXPECT().CreateOrUpdateTask(mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockClient.EXPECT().Get(mock.Anything, client.ObjectKeyFromObject(secret),
		mock.AnythingOfType("*unstructured.Unstructured")).RunAndReturn(
		func(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
			secret.DeepCopyInto(obj.(*unstructured.Unstructured))
			return nil
		})

	// when
	result, reconcileErr := controller.Reconcile(ctx, ctrl.Request{
		NamespacedName: client.ObjectKeyFromObject(secret),
	})

	// then
	date, _ := dateparse.ParseAny(secret.
		Object["metadata"].(map[string]interface{})["annotations"].(map[string]interface{})[config.
		ExpirationLabel].(string))
	mockScheduler.AssertCalled(t, "CreateOrUpdateTask", "v1/Secret/my-secret/my-namespace/delete",
		date, mock.Anything)
	assert.Nil(t, reconcileErr)
	assert.False(t, result.Requeue)
}

func TestController_ReconcileIntegration(t *testing.T) {
	// given
	var (
		ctx    = context.Background()
		secret = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Secret",
				"metadata": map[string]interface{}{
					"name":      strings.ToLower(strings.ReplaceAll(gofakeit.Name(), " ", "")),
					"namespace": "default",
					"annotations": map[string]interface{}{
						testVars.config.ExpirationLabel: "5s",
					},
				},
			},
		}
	)

	// when
	createErr := testVars.k8sClient.Create(ctx, secret)

	// then
	assert.Nil(t, createErr)
	assert.Eventually(t, func() bool {
		return testVars.k8sClient.Get(ctx, client.ObjectKeyFromObject(secret), secret) == nil
	}, 60*time.Second, 100*time.Millisecond)
	assert.Eventually(t, func() bool {
		return errors.IsNotFound(testVars.k8sClient.Get(ctx, client.ObjectKeyFromObject(secret), secret))
	}, 60*time.Second, 100*time.Millisecond)
}

func TestController_Reconcile_ShouldDeleteTaskWhenNotFound(t *testing.T) {
	// given
	var (
		ctx            = context.Background()
		config         = common.NewConfig()
		mockClient     = new(client2.MockClient)
		mockScheduler  = new(common2.MockScheduler)
		apiVersionKind = "v1;Secret"
		controller     = NewController(config, mockClient, apiVersionKind, mockScheduler)
	)

	mockScheduler.EXPECT().DeleteTask(mock.Anything).Return(nil)
	mockClient.EXPECT().Get(mock.Anything, mock.Anything,
		mock.AnythingOfType("*unstructured.Unstructured")).Return(errors.NewNotFound(schema.GroupResource{}, ""))

	// when
	result, reconcileErr := controller.Reconcile(ctx, ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "my-secret",
			Namespace: "my-namespace",
		},
	})

	// then
	mockScheduler.AssertCalled(t, "DeleteTask", "v1/Secret/my-secret/my-namespace/delete")
	assert.Nil(t, reconcileErr)
	assert.False(t, result.Requeue)
}

func TestController_Reconcile_ShouldDeleteTaskWhenAnnotationNotFound(t *testing.T) {
	// given
	var (
		ctx            = context.Background()
		config         = common.NewConfig()
		mockClient     = new(client2.MockClient)
		mockScheduler  = new(common2.MockScheduler)
		apiVersionKind = "v1;Secret"
		controller     = NewController(config, mockClient, apiVersionKind, mockScheduler)
		secret         = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion":     "v1",
				"apiVersionKind": "Secret",
				"type":           "Opaque",
				"metadata": map[string]interface{}{
					"name":      "my-secret",
					"namespace": "my-namespace",
					"annotations": map[string]interface{}{
						"not-related": gofakeit.FutureDate().String(),
					},
				},
				"data": map[string]interface{}{
					"my-key": "my-value",
				},
			},
		}
	)

	mockScheduler.EXPECT().DeleteTask(mock.Anything).Return(nil)
	mockClient.EXPECT().Get(mock.Anything, client.ObjectKeyFromObject(secret),
		mock.AnythingOfType("*unstructured.Unstructured")).RunAndReturn(
		func(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
			secret.DeepCopyInto(obj.(*unstructured.Unstructured))
			return nil
		})

	// when
	result, reconcileErr := controller.Reconcile(ctx, ctrl.Request{
		NamespacedName: client.ObjectKeyFromObject(secret),
	})

	// then
	mockScheduler.AssertCalled(t, "DeleteTask", "v1/Secret/my-secret/my-namespace/delete")
	assert.Nil(t, reconcileErr)
	assert.False(t, result.Requeue)
}

func TestController_Reconcile_ShouldReturnErrWhenAnnotationValueDoesNotMakeSense(t *testing.T) {
	// given
	var (
		ctx            = context.Background()
		config         = common.NewConfig()
		mockClient     = new(client2.MockClient)
		mockScheduler  = new(common2.MockScheduler)
		apiVersionKind = "v1;Secret"
		controller     = NewController(config, mockClient, apiVersionKind, mockScheduler)
		secret         = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion":     "v1",
				"apiVersionKind": "Secret",
				"type":           "Opaque",
				"metadata": map[string]interface{}{
					"name":      "my-secret",
					"namespace": "my-namespace",
					"annotations": map[string]interface{}{
						config.ExpirationLabel: gofakeit.BeerName(),
					},
				},
				"data": map[string]interface{}{
					"my-key": "my-value",
				},
			},
		}
	)

	mockClient.EXPECT().Get(mock.Anything, client.ObjectKeyFromObject(secret),
		mock.AnythingOfType("*unstructured.Unstructured")).RunAndReturn(
		func(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
			secret.DeepCopyInto(obj.(*unstructured.Unstructured))
			return nil
		})

	// when
	result, reconcileErr := controller.Reconcile(ctx, ctrl.Request{
		NamespacedName: client.ObjectKeyFromObject(secret),
	})

	// then
	assert.NotNil(t, reconcileErr)
	assert.False(t, result.Requeue)
}

func TestController_Reconcile_ShouldDeleteTaskWhenAnnotationValueIsPast(t *testing.T) {
	// given
	var (
		ctx            = context.Background()
		config         = common.NewConfig()
		mockClient     = new(client2.MockClient)
		mockScheduler  = new(common2.MockScheduler)
		apiVersionKind = "v1;Secret"
		controller     = NewController(config, mockClient, apiVersionKind, mockScheduler)
		secret         = &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion":     "v1",
				"apiVersionKind": "Secret",
				"type":           "Opaque",
				"metadata": map[string]interface{}{
					"name":      "my-secret",
					"namespace": "my-namespace",
					"annotations": map[string]interface{}{
						config.ExpirationLabel: gofakeit.PastDate().String(),
					},
				},
				"data": map[string]interface{}{
					"my-key": "my-value",
				},
			},
		}
	)
	mockScheduler.EXPECT().CreateOrUpdateTask(mock.Anything, mock.Anything, mock.Anything).
		Return(gocron.ErrOneTimeJobStartDateTimePast)
	mockScheduler.EXPECT().DeleteTask(mock.Anything).Return(nil)
	mockClient.EXPECT().Delete(mock.Anything, mock.Anything).Return(nil)
	mockClient.EXPECT().Get(mock.Anything, client.ObjectKeyFromObject(secret),
		mock.AnythingOfType("*unstructured.Unstructured")).RunAndReturn(
		func(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
			secret.DeepCopyInto(obj.(*unstructured.Unstructured))
			return nil
		})

	// when
	result, reconcileErr := controller.Reconcile(ctx, ctrl.Request{
		NamespacedName: client.ObjectKeyFromObject(secret),
	})

	// then
	mockClient.AssertCalled(t, "Delete", mock.Anything, secret)
	assert.Nil(t, reconcileErr)
	assert.False(t, result.Requeue)
}

func TestController_SetupWithManager(t *testing.T) {
	// given
	var (
		mockClient     = new(client2.MockClient)
		mockManager    = new(manager2.MockManager)
		mockCache      = new(cache2.MockCache)
		mockScheduler  = new(common2.MockScheduler)
		apiVersionKind = "v1;Secret"

		controller = NewController(common.NewConfig(), mockClient, apiVersionKind, mockScheduler)
	)

	mockManager.EXPECT().GetControllerOptions().Return(config.Controller{})
	mockManager.EXPECT().GetScheme().Return(runtime.NewScheme())
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
	assert.Nil(t, setupErr)
}
