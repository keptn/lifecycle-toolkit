package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConfig_GetDefaultCreationRequestTimeout(t *testing.T) {
	i := Instance()

	timeout := i.GetCreationRequestTimeout()

	require.Equal(t, defaultKeptnAppCreationRequestTimeout, timeout)
}

func TestConfig_SetDefaultCreationRequestTimeout(t *testing.T) {
	i := Instance()

	i.SetCreationRequestTimeout(5 * time.Second)

	timeout := i.GetCreationRequestTimeout()
	require.Equal(t, 5*time.Second, timeout)
}

func TestGetOptionsInstance(t *testing.T) {
	o := Instance()
	require.NotNil(t, o)

	o.SetCreationRequestTimeout(5 * time.Second)

	// verify that all sets/gets operator on the same instance
	o2 := Instance()

	timeout := o2.GetCreationRequestTimeout()
	require.Equal(t, 5*time.Second, timeout)
}

func TestConfig_SetAndGetDefaultNamespace(t *testing.T) {
	i := Instance()

	ns := i.GetDefaultNamespace()

	require.Empty(t, ns)
	i.SetDefaultNamespace("test")
	require.Equal(t, "test", i.GetDefaultNamespace())
}

func TestConfig_SetAndGetCloudEventEndpoint(t *testing.T) {
	i := Instance()

	ns := i.GetCloudEventsEndpoint()
	require.Empty(t, ns)
	i.SetCloudEventsEndpoint("mytestendpoint")
	require.Equal(t, "mytestendpoint", i.GetCloudEventsEndpoint())
}

func TestConfig_SetAndGetBlockDeployment(t *testing.T) {
	i := Instance()

	blocked := i.GetBlockDeployment()
	require.True(t, blocked)
	i.SetBlockDeployment(false)
	require.False(t, i.GetBlockDeployment())
}

func TestConfig_SetAndGetObservabilityTimeout(t *testing.T) {
	i := Instance()

	require.Equal(t, metav1.Duration{
		Duration: time.Duration(5 * time.Minute),
	}, i.GetObservabilityTimeout())

	i.SetObservabilityTimeout(metav1.Duration{
		Duration: time.Duration(10 * time.Minute),
	})

	require.Equal(t, metav1.Duration{
		Duration: time.Duration(10 * time.Minute),
	}, i.GetObservabilityTimeout())
}
