package config

import (
	"sync"
	"time"
)

const defaultKeptnAppCreationRequestTimeout = 30 * time.Second

//go:generate moq -pkg fake -skip-ensure -out ./fake/config_mock.go . IConfig:MockConfig
type IConfig interface {
	SetCreationRequestTimeout(value time.Duration)
	GetCreationRequestTimeout() time.Duration
	SetCloudEventsEndpoint(endpoint string)
	GetCloudEventsEndpoint() string
}

type ControllerConfig struct {
	keptnAppCreationRequestTimeout time.Duration
	cloudEventsEndpoint            string
}

var instance *ControllerConfig
var once = sync.Once{}

func Instance() *ControllerConfig {
	once.Do(func() {
		instance = &ControllerConfig{keptnAppCreationRequestTimeout: defaultKeptnAppCreationRequestTimeout}
	})
	return instance
}

func (o *ControllerConfig) SetCreationRequestTimeout(value time.Duration) {
	o.keptnAppCreationRequestTimeout = value
}

func (o *ControllerConfig) GetCreationRequestTimeout() time.Duration {
	return o.keptnAppCreationRequestTimeout
}

func (o *ControllerConfig) SetCloudEventsEndpoint(endpoint string) {
	o.cloudEventsEndpoint = endpoint
}

func (o *ControllerConfig) GetCloudEventsEndpoint() string {
	return o.cloudEventsEndpoint
}
