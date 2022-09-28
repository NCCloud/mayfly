package main

import (
	"github.com/NCCloud/mayfly/pkg/common"
	"github.com/NCCloud/mayfly/pkg/controllers/mayfly"
	"github.com/NCCloud/mayfly/pkg/controllers/mayfly/resource"
	"github.com/NCCloud/mayfly/pkg/scheduler"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	logger := zap.New()
	logger.Info("Starting the Mayfly manager")

	config, configErr := resource.NewConfig("./config.yaml")
	if configErr != nil {
		panic(configErr)
	}

	scheme := runtime.NewScheme()

	operatorConfig := common.NewOperatorConfig()

	manager, managerErr := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Port:                   9443,
		Scheme:                 scheme,
		Logger:                 logger,
		MetricsBindAddress:     ":8082",
		HealthProbeBindAddress: ":8083",
		LeaderElection:         operatorConfig.EnableLeaderElection,
		LeaderElectionID:       "mayfly-operator-leader.spaceship.com",
	})
	if managerErr != nil {
		panic(managerErr)
	}

	client := manager.GetClient()

	mayflyScheduler := scheduler.NewScheduler(operatorConfig, client)
	mayflyScheduler.StartMonitor()

	for _, resource := range config.Resources {
		if newControllerErr := mayfly.NewController(operatorConfig, client, resource, mayflyScheduler).SetupWithManager(manager); newControllerErr != nil {
			panic(newControllerErr)
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
