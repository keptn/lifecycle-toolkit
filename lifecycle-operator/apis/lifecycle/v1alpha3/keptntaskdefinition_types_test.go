package v1alpha3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTaskDefinition_GetServiceAccountNoName(t *testing.T) {
	d := &KeptnTaskDefinition{
		Spec: KeptnTaskDefinitionSpec{},
	}
	svcAccname := d.GetServiceAccount()
	require.Len(t, svcAccname, 0)
}

func TestTaskDefinition_GetServiceAccountName(t *testing.T) {
	svcAccName := "sva"
	d := &KeptnTaskDefinition{
		Spec: KeptnTaskDefinitionSpec{
			ServiceAccount: &ServiceAccountSpec{
				Name: svcAccName,
			},
		},
	}
	svcAccname := d.GetServiceAccount()
	require.Len(t, svcAccname, 3)
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
	automountSaToken := d.GetAutomountServiceAccountToken()
	require.True(t, *automountSaToken)
}
