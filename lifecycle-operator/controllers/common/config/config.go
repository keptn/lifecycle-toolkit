package config

import (
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultKeptnAppCreationRequestTimeout = 30 * time.Second

//go:generate moq -pkg fake -skip-ensure -out ./fake/config_mock.go . IConfig:MockConfig
type IConfig interface {
	SetCreationRequestTimeout(value time.Duration)
	GetCreationRequestTimeout() time.Duration
	SetCloudEventsEndpoint(endpoint string)
	GetCloudEventsEndpoint() string
	SetDefaultNamespace(namespace string)
	GetDefaultNamespace() string
	SetBlockDeployment(value bool)
	GetBlockDeployment() bool
	SetObservabilityTimeout(timeout metav1.Duration)
	GetObservabilityTimeout() metav1.Duration
}

type ControllerConfig struct {
	keptnAppCreationRequestTimeout time.Duration
	cloudEventsEndpoint            string
	defaultNamespace               string
	blockDeployment                bool
	observabilityTimeout           metav1.Duration
}

var instance *ControllerConfig
var once = sync.Once{}

func Instance() *ControllerConfig {
	once.Do(func() {
		instance = &ControllerConfig{
			keptnAppCreationRequestTimeout: defaultKeptnAppCreationRequestTimeout,
			blockDeployment:                true,
			observabilityTimeout: metav1.Duration{
				Duration: time.Duration(5 * time.Minute),
			},
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

func (o *ControllerConfig) SetBlockDeployment(value bool) {
	o.blockDeployment = value
}

func (o *ControllerConfig) GetBlockDeployment() bool {
	return o.blockDeployment
}

func (o *ControllerConfig) SetObservabilityTimeout(timeout metav1.Duration) {
	o.observabilityTimeout = timeout
}

func (o *ControllerConfig) GetObservabilityTimeout() metav1.Duration {
	return o.observabilityTimeout
}
