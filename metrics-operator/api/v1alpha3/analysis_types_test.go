package v1alpha3

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAnalysis_GetAnalysisDefinitionNamespace(t *testing.T) {
	a := Analysis{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "ns",
		},
	}

	require.Equal(t, "ns", a.GetAnalysisDefinitionNamespace())

	a.Spec.AnalysisDefinition.Namespace = "ns2"

	require.Equal(t, "ns2", a.GetAnalysisDefinitionNamespace())
}
