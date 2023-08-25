package analysis

import (
	"context"
	"fmt"
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
	// Define test cases
	testCases := []struct {
		name            string
		providerType    string
		mockProvider    providers.KeptnSLIProvider
		providerRequest metricstypes.ProviderRequest
		expectedResult  metricstypes.ProviderResult
		expectedError   error
	}{
		{
			name:         "SuccessfulEvaluation",
			providerType: "mockProvider",
			mockProvider: &fake2.KeptnSLIProviderMock{
				FetchAnalysisValueFunc: func(ctx context.Context, query string, spec metricsapi.AnalysisSpec, provider *metricsapi.KeptnMetricsProvider) (string, []byte, error) {
					return "10", []byte{}, nil
				},
			},
			providerRequest: metricstypes.ProviderRequest{
				Objective: &metricsapi.Objective{
					AnalysisValueTemplateRef: metricsapi.ObjectReference{
						Name:      "mytemp",
						Namespace: "default",
					},
				},
				Query:    "query_fake_metric",
				Provider: &metricsapi.KeptnMetricsProvider{Spec: metricsapi.KeptnMetricsProviderSpec{Type: "prometheus"}},
			},
			expectedResult: metricstypes.ProviderResult{
				Objective: metricsapi.ObjectReference{
					Name:      "mytemp",
					Namespace: "default",
				},
				Value: "10",
				Err:   nil,
			},
			expectedError: nil,
		},
		{
			name:         "FailedEvaluation",
			providerType: "mockProvider",
			mockProvider: &fake2.KeptnSLIProviderMock{
				FetchAnalysisValueFunc: func(ctx context.Context, query string, spec metricsapi.AnalysisSpec, provider *metricsapi.KeptnMetricsProvider) (string, []byte, error) {
					return "", []byte{}, fmt.Errorf("something bad")
				},
			},
			providerRequest: metricstypes.ProviderRequest{
				Objective: &metricsapi.Objective{
					AnalysisValueTemplateRef: metricsapi.ObjectReference{
						Name:      "mytemp",
						Namespace: "default",
					},
				},
				Query:    "query_fake_metric",
				Provider: &metricsapi.KeptnMetricsProvider{Spec: metricsapi.KeptnMetricsProviderSpec{Type: "prometheus"}},
			},
			expectedResult: metricstypes.ProviderResult{
				Objective: metricsapi.ObjectReference{
					Name:      "mytemp",
					Namespace: "default",
				},
				Value: "",
				Err:   fmt.Errorf("something bad"),
			},
			expectedError: fmt.Errorf("something bad"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockProviderFactory := func(providerType string, log logr.Logger, client client.Client) (providers.KeptnSLIProvider, error) {
				// Define your mock provider implementation
				return tc.mockProvider, nil
			}
			objectivesEvaluator := ObjectivesEvaluator{
				// Initialize ObjectivesEvaluator fields here
				Client:          fake.NewClient(),
				Log:             logr.Discard(),
				ProviderFactory: mockProviderFactory,
				Analysis: &metricsapi.Analysis{
					Spec: metricsapi.AnalysisSpec{
						AnalysisDefinition: metricsapi.ObjectReference{
							Name:      "a",
							Namespace: "d",
						},
					},
				},
				results: make(chan metricstypes.ProviderResult, 1),
			}

			ctx := context.TODO()
			objChan := make(chan metricstypes.ProviderRequest, 1)
			go func() {
				objChan <- tc.providerRequest
				close(objChan)
			}()
			objectivesEvaluator.Evaluate(ctx, tc.providerType, objChan)
			close(objectivesEvaluator.results)
			result := <-objectivesEvaluator.results

			require.Equal(t, tc.expectedResult, result)
			require.Equal(t, tc.expectedError, result.Err)
		})
	}
}
