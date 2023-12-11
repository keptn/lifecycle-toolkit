package config

import (
	"sync"
	"time"
)

const (
	defaultKeptnAppCreationRequestTimeout = 30 * time.Second
	defaultBlockDeployment                = true
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/config_mock.go . IConfig:MockConfig
type IConfig interface {
	SetCreationRequestTimeout(value time.Duration)
	GetCreationRequestTimeout() time.Duration
	SetCloudEventsEndpoint(endpoint string)
	GetCloudEventsEndpoint() string
	SetDefaultNamespace(namespace string)
	GetDefaultNamespace() string
	SetBlockDeployment(blockDeployment bool)
	GetBlockDeployment() bool
}

type ControllerConfig struct {
	keptnAppCreationRequestTimeout time.Duration
	cloudEventsEndpoint            string
	defaultNamespace               string
	blockDeployment                bool
}

var instance *ControllerConfig
var once = sync.Once{}

func Instance() *ControllerConfig {
	once.Do(func() {
		instance = &ControllerConfig{
			keptnAppCreationRequestTimeout: defaultKeptnAppCreationRequestTimeout,
			blockDeployment:                defaultBlockDeployment,
		}
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

func (o *ControllerConfig) SetDefaultNamespace(ns string) {
	o.defaultNamespace = ns
}

func (o *ControllerConfig) GetDefaultNamespace() string {
	return o.defaultNamespace
}

func (o *ControllerConfig) SetBlockDeployment(blockDeployment bool) {
	o.blockDeployment = blockDeployment
}

func (o *ControllerConfig) GetBlockDeployment() bool {
	return o.blockDeployment
}
