package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigProvider(t *testing.T) {
	provider := NewKubeConfigProvider()

	assert.NotNil(t, provider)
}
