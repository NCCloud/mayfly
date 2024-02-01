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
