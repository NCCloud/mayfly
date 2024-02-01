package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	// given

	// when
	config := NewConfig()

	// then
	assert.NotEqual(t, *config, Config{})
}
