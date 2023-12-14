package common

import (
	"time"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	EnableLeaderElection bool          `env:"ENABLE_LEADER_ELECTION" envDefault:"false"`
	SyncPeriod           time.Duration `env:"SYNC_PERIOD" envDefault:"30m"`
	MonitoringInterval   time.Duration `env:"MONITORING_INTERVAL" envDefault:"5s"`
	ExpirationLabel      string        `env:"EXPIRATION_LABEL" envDefault:"mayfly.cloud.namecheap.com/expire"`
	Resources            []string      `env:"RESOURCES" envSeparator:"," envDefault:"v1;Secret,cloud.namecheap.com/v1alpha1;ScheduledResource"`
}

func NewConfig() *Config {
	operatorConfig := &Config{}
	if err := env.Parse(operatorConfig); err != nil {
		panic(err)
	}

	return operatorConfig
}
