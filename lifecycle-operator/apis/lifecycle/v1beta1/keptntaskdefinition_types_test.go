package v1beta1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTaskDefinition_GetServiceAccountNoName(t *testing.T) {
	d := &KeptnTaskDefinition{
		Spec: KeptnTaskDefinitionSpec{},
	}
	svcAccname := d.GetServiceAccount()
	require.Equal(t, svcAccname, "")
}

func TestTaskDefinition_GetServiceAccountName(t *testing.T) {
	sAName := "sva"
	d := &KeptnTaskDefinition{
		Spec: KeptnTaskDefinitionSpec{
			ServiceAccount: &ServiceAccountSpec{
				Name: sAName,
			},
		},
	}
	svcAccname := d.GetServiceAccount()
	require.Equal(t, svcAccname, sAName)
}

func TestTaskDefinition_GetAutomountServiceAccountToken(t *testing.T) {
	token := true
	d := &KeptnTaskDefinition{
		Spec: KeptnTaskDefinitionSpec{
			AutomountServiceAccountToken: &AutomountServiceAccountTokenSpec{
				Type: &token,
			},
		},
	}
	require.True(t, *d.GetAutomountServiceAccountToken())
}
