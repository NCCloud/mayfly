package pkg

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	EnableLeaderElection bool           `env:"ENABLE_LEADER_ELECTION" envDefault:"false"`
	SyncPeriod           *time.Duration `env:"SYNC_PERIOD" envDefault:"30m"`
	ExpirationLabel      string         `env:"EXPIRATION_LABEL" envDefault:"mayfly.cloud.namecheap.com/expire"`
	Resources            []string       `env:"RESOURCES" envSeparator:"," envDefault:"v1;Secret"`
}

func NewConfig() *Config {
	operatorConfig := &Config{}
	if err := env.Parse(operatorConfig); err != nil {
		panic(err)
	}

	return operatorConfig
}
