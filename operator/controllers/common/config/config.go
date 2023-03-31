package config

import (
	"sync"
	"time"
)

const defaultKeptnAppCreationRequestTimeout = 30 * time.Second

type Config struct {
	keptnAppCreationRequestTimeout time.Duration
}

var instance *Config
var once = sync.Once{}

func Instance() *Config {
	once.Do(func() {
		instance = &Config{keptnAppCreationRequestTimeout: defaultKeptnAppCreationRequestTimeout}
	})
	return instance
}

func (o *Config) SetCreationRequestTimeout(value time.Duration) {
	o.keptnAppCreationRequestTimeout = value
}

func (o *Config) GetCreationRequestTimeout() time.Duration {
	return o.keptnAppCreationRequestTimeout
}
