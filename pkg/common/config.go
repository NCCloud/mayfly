package common

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type OperatorConfig struct {
	EnableLeaderElection          bool `env:"ENABLE_LEADER_ELECTION" envDefault:"false"`
	ResourceConfiguration         ResourceConfig
	GarbageCollectorConfiguration GarbageCollectorConfig
}

type ResourceConfig struct {
	MayflyExpireLabel string `env:"MAYFLY_EXPIRE_LABEL" envDefault:"mayfly.cloud.spaceship.com/expire"`
}

type GarbageCollectorConfig struct {
	GarbageCollectionPeriod time.Duration `env:"GARBAGE_COLLECTION_PERIOD" envDefault:"30m"`
}

func NewOperatorConfig() *OperatorConfig {
	operatorConfig := &OperatorConfig{}
	if err := env.Parse(operatorConfig); err != nil {
		panic(err)
	}

	return operatorConfig
}
