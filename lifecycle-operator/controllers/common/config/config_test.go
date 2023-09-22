package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
