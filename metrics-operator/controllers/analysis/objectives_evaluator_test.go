package analysis

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers"
	fake2 "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/fake"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestEvaluate(t *testing.T) {
	fakeProvider := fake2.KeptnSLIProviderMock{
		FetchAnalysisValueFunc: func(ctx context.Context, query string, spec metricsapi.AnalysisSpec, provider *metricsapi.KeptnMetricsProvider) (string, []byte, error) {
			return "10", []byte{}, nil
		},
	}

	evaluator := ObjectivesEvaluator{
		Analysis: &metricsapi.Analysis{
			Spec: metricsapi.AnalysisSpec{
				Args: map[string]string{"a": "b"},
				AnalysisDefinition: metricsapi.ObjectReference{
					Name:      "a",
					Namespace: "d",
				},
			},
		},

		Client:  fake.NewClient(),
		Log:     logr.Discard(),
		results: make(chan metricstypes.ProviderResult, 1),
		ProviderFactory: func(providerType string, log logr.Logger, k8sClient client.Client) (providers.KeptnSLIProvider, error) {
			return &fakeProvider, nil
		},
	}

	// Simulate context
	ctx := context.TODO()

	// Test with a valid provider type
	providerType := "prometheus"
	obj := make(chan metricstypes.ProviderRequest, 1)
	go func() {
		obj <- metricstypes.ProviderRequest{
			Objective: &metricsapi.Objective{
				AnalysisValueTemplateRef: metricsapi.ObjectReference{
					Name:      "mytemp",
					Namespace: "default",
				},
			},
			Query:    "query_fake_metric",
			Provider: &metricsapi.KeptnMetricsProvider{Spec: metricsapi.KeptnMetricsProviderSpec{Type: providerType}},
		}
		close(obj)
	}()

	evaluator.Evaluate(ctx, providerType, obj)
	close(evaluator.results)
	for res := range evaluator.results {
		require.Equal(t, "mytemp", res.Objective.Name)
		require.Contains(t, "10", res.Value)
		require.Nil(t, res.Err)
	}

}
