package main

import (
	"fmt"
	"time"

	"github.com/NCCloud/mayfly/cmd/garbagecollector/collectors"
	"github.com/NCCloud/mayfly/pkg/common"
	"github.com/NCCloud/mayfly/pkg/controllers/mayfly/resource"
	"github.com/go-co-op/gocron"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	logger            logr.Logger
	operatorConfig    *common.OperatorConfig
	scheduler         *gocron.Scheduler
	garbageCollectors []collectors.Collector
)

func init() {
	scheme := runtime.NewScheme()
	logger = zap.New().WithName("GarbageCollector")
	operatorConfig = common.NewOperatorConfig()
	scheduler = gocron.NewScheduler(time.UTC)

	config, configErr := resource.NewConfig("./config.yaml")
	if configErr != nil {
		panic(configErr)
	}

	mgrClient, mgrClientErr := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme})
	if mgrClientErr != nil {
		panic(mgrClientErr)
	}
	garbageCollectors = []collectors.Collector{}
	for _, resource := range config.Resources {

		logger.Info(fmt.Sprintf("Adding garbage collector for %s", resource.Kind))
		garbageCollectors = append(garbageCollectors, collectors.NewResourcesCollector(logger, mgrClient, resource, operatorConfig))
	}
}

func main() {
	logger.Info(fmt.Sprintf("Starting garbage collectors with %s interval",
		operatorConfig.GarbageCollectorConfiguration.GarbageCollectionPeriod.String()))

	for _, garbageCollector := range garbageCollectors {
		if _, doErr := scheduler.Every(operatorConfig.GarbageCollectorConfiguration.GarbageCollectionPeriod).Do(func(gCollector collectors.Collector) {
			if collectErr := gCollector.Collect(); collectErr != nil {
				logger.Error(collectErr, "failed to collect garbage")
			}
		}, garbageCollector); doErr != nil {
			panic(doErr)
		}
	}

	scheduler.StartBlocking()
}
