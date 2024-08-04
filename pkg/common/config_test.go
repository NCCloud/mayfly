package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// given

	// when
	config := NewConfig()

	// then
	assert.NotEqual(t, *config, Config{})
}

func TestGroupAdjacentGroupVersionKinds(t *testing.T) {
	t.Setenv("RESOURCES", "v1;Secret;ConfigMap")

	config := NewConfig()

	assert.Equal(t, config.Resources, []string{"v1;Secret", "v1;ConfigMap"})
}
