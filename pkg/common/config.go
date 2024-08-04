package common

import (
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	EnableLeaderElection bool          `env:"ENABLE_LEADER_ELECTION" envDefault:"false"`
	SyncPeriod           time.Duration `env:"SYNC_PERIOD" envDefault:"60m"`
	MonitoringInterval   time.Duration `env:"MONITORING_INTERVAL" envDefault:"5s"`
	ExpirationLabel      string        `env:"EXPIRATION_LABEL" envDefault:"mayfly.cloud.namecheap.com/expire"`
	Resources            []string      `env:"RESOURCES" envSeparator:"," envDefault:"v1;Secret,cloud.namecheap.com/v1alpha1;ScheduledResource"`
}

func NewConfig() *Config {
	operatorConfig := &Config{}
	if err := env.Parse(operatorConfig); err != nil {
		panic(err)
	}

	operatorConfig.GroupAdjacentGroupVersionKinds()

	return operatorConfig
}

func (c *Config) GroupAdjacentGroupVersionKinds() {
	var resources []string
	for _, r := range c.Resources {
		apiVersionKind := strings.Split(r, ";")
		apiVersion := apiVersionKind[0]

		for _, kind := range apiVersionKind[1:] {
			resources = append(resources, apiVersion+";"+kind)
		}
	}

	c.Resources = resources
}
