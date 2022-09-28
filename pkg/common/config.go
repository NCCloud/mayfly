package common

import "github.com/caarlos0/env/v6"

type OperatorConfig struct {
	EnableLeaderElection  bool `env:"ENABLE_LEADER_ELECTION" envDefault:"false"`
	ResourceConfiguration ResourceConfig
}

type ResourceConfig struct {
	MayflyExpireAnnotation string `env:"MAYFLY_EXPIRE_ANNOTATION" envDefault:"mayfly.cloud.spaceship.com/expire"`
}

func NewOperatorConfig() *OperatorConfig {
	operatorConfig := &OperatorConfig{}
	if err := env.Parse(operatorConfig); err != nil {
		panic(err)
	}

	return operatorConfig
}
