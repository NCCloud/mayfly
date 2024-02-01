package main

import (
	"fmt"

	"github.com/NCCloud/mayfly/pkg/controllers/expiration"
	"github.com/NCCloud/mayfly/pkg/controllers/scheduledresource"

	"github.com/NCCloud/mayfly/pkg/apis/v1alpha1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/NCCloud/mayfly/pkg/common"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const (
	metricPort = 8082
	healthPort = 8083
)

func main() {
	logger := zap.New()
	scheme := runtime.NewScheme()
	config := common.NewConfig()

	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	ctrl.SetLogger(logger)

	logger.Info("Configuration", "config", config)

	manager, managerErr := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Logger: logger,
		Metrics: server.Options{
			BindAddress: fmt.Sprintf(":%d", metricPort),
		},
		HealthProbeBindAddress: fmt.Sprintf(":%d", healthPort),
		LeaderElection:         config.EnableLeaderElection,
		LeaderElectionID:       "mayfly-leader.cloud.namecheap.com",
		Cache: cache.Options{
			SyncPeriod: &config.SyncPeriod,
		},
	})
	if managerErr != nil {
		panic(managerErr)
	}

	client := manager.GetClient()
	scheduler := common.NewScheduler(config)

	for _, resource := range config.Resources {
		if expirationControllerErr := expiration.NewController(config, client, resource, scheduler).
			SetupWithManager(manager); expirationControllerErr != nil {
			panic(expirationControllerErr)
		}
	}

	if scheduledResourceControllerErr := scheduledresource.NewController(config, client, scheduler).
		SetupWithManager(manager); scheduledResourceControllerErr != nil {
		panic(scheduledResourceControllerErr)
	}

	if addHealthCheckErr := manager.AddHealthzCheck("healthz", healthz.Ping); addHealthCheckErr != nil {
		panic(addHealthCheckErr)
	}

	if addReadyCheckErr := manager.AddReadyzCheck("readyz", healthz.Ping); addReadyCheckErr != nil {
		panic(addReadyCheckErr)
	}

	if startManagerErr := manager.Start(ctrl.SetupSignalHandler()); startManagerErr != nil {
		panic(startManagerErr)
	}
}
