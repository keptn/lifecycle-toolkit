package v1alpha3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalysisValueTemplate_GetProviderNamespace(t *testing.T) {
	a := AnalysisValueTemplate{}

	require.Equal(t, "default", a.GetProviderNamespace("default"))

	a.Spec.Provider.Namespace = "ns"

	require.Equal(t, "ns", a.GetProviderNamespace("default"))
}
