package pkg

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	EnableLeaderElection bool           `env:"ENABLE_LEADER_ELECTION" envDefault:"false"`
	SyncPeriod           *time.Duration `env:"SYNC_PERIOD" envDefault:"30m"`
	ExpirationLabel      string         `env:"EXPIRATION_LABEL" envDefault:"mayfly.cloud.spaceship.com/expire"`
	Resources            []string       `env:"RESOURCES" envSeparator:"," envDefault:"spaceship.com/v1alpha1;SSHip"`
}

func NewConfig() *Config {
	operatorConfig := &Config{}
	if err := env.Parse(operatorConfig); err != nil {
		panic(err)
	}

	return operatorConfig
}
