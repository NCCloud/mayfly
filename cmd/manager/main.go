package main

import (
	"fmt"

	"github.com/NCCloud/mayfly/pkg"
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
	config := pkg.NewConfig()

	logger.Info("Configuration", "config", config)

	manager, managerErr := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Logger:                 logger,
		MetricsBindAddress:     fmt.Sprintf(":%d", metricPort),
		HealthProbeBindAddress: fmt.Sprintf(":%d", healthPort),
		LeaderElection:         config.EnableLeaderElection,
		LeaderElectionID:       "mayfly-leader.cloud.namecheap.com",
		Cache: cache.Options{
			SyncPeriod: config.SyncPeriod,
		},
	})
	if managerErr != nil {
		panic(managerErr)
	}

	client := manager.GetClient()
	scheduler := pkg.NewScheduler(config, client)

	for _, resource := range config.Resources {
		if err := pkg.NewController(config, client, resource, scheduler).SetupWithManager(manager); err != nil {
			panic(err)
		}
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
