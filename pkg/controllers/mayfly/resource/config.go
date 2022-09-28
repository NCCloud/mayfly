package resource

import (
	"os"

	"k8s.io/apimachinery/pkg/util/yaml"
)

func NewConfig(path string) (*Resources, error) {
	file, readFileErr := os.ReadFile(path)
	if readFileErr != nil {
		return nil, readFileErr
	}

	var cfg Resources
	configUnmarshallErr := yaml.Unmarshal(file, &cfg)

	return &cfg, configUnmarshallErr
}
