package analysis

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/analysis/fake"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	fake2 "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_generateQuery(t *testing.T) {

	tests := []struct {
		name      string
		query     string
		selectors map[string]string
		want      string
		wanterror string
	}{
		{
			name:  "successful, all args exist",
			query: "this is a {{.good}} query{{.dot}}",
			selectors: map[string]string{
				"good": "good",
				"dot":  ".",
			},
			want: "this is a good query.",
		},
		{
			name:  "no substitution, all args missing",
			query: "this is a {{.good}} query{{.dot}}",
			selectors: map[string]string{
				"bad":    "good",
				"dotted": ".",
			},
			want: "this is a <no value> query<no value>",
		},
		{
			name:  "no substitution, bad template",
			query: "this is a {{.good} query{{.dot}}",
			selectors: map[string]string{
				"bad":    "good",
				"dotted": ".",
			},
			want:      "",
			wanterror: "could not create a template:",
		},
		{
			name:      "nothing to do",
			query:     "this is a query",
			selectors: map[string]string{},
			want:      "this is a query",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateQuery(tt.query, tt.selectors)
			if tt.wanterror != "" {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tt.wanterror)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestProvidersPool(t *testing.T) {
	// Define your test cases

	analysis, analysisDef, template, provider := getTestCRDs()

	provider2 := metricsapi.KeptnMetricsProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-provider",
			Namespace: "default2",
		},
		Spec: metricsapi.KeptnMetricsProviderSpec{
			Type:         "prometheus",
			TargetServer: "localhost:2000",
		},
	}

	template2 := metricsapi.AnalysisValueTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-template",
			Namespace: "default",
		},
		Spec: metricsapi.AnalysisValueTemplateSpec{
			Provider: metricsapi.ObjectReference{
				Name: "my-provider",
			},
			Query: "this is a {{.good}} query{{.dot}}",
		},
	}

	analysisDef2 := metricsapi.AnalysisDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-analysis-def",
			Namespace: "default",
		},
		Spec: metricsapi.AnalysisDefinitionSpec{
			Objectives: []metricsapi.Objective{
				{
					AnalysisValueTemplateRef: metricsapi.ObjectReference{
						Name: "my-template",
					},
					Weight:       1,
					KeyObjective: false,
				},
			},
			TotalScore: metricsapi.TotalScore{
				PassPercentage:    0,
				WarningPercentage: 0,
			},
		},
	}

	template3 := metricsapi.AnalysisValueTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-template",
			Namespace: "default2",
		},
		Spec: metricsapi.AnalysisValueTemplateSpec{
			Provider: metricsapi.ObjectReference{
				Name:      "my-provider",
				Namespace: "default",
			},
			Query: "this is a {{.good}} query{{.dot}}",
		},
	}

	provider.Spec.Type = "mock-provider"
	provider2.Spec.Type = "mock-provider"

	testCases := []struct {
		name           string
		expectedErr    string
		mockClient     client.Client
		analysisDef    metricsapi.AnalysisDefinition
		providerResult *metricstypes.ProviderRequest
	}{

		{
			name:        "MissingTemplate",
			expectedErr: "analysisvaluetemplates.metrics.keptn.sh \"my-template\" not found",
			analysisDef: analysisDef,
			mockClient:  fake2.NewClient(&analysis, &analysisDef),
		},
		{
			name:        "MissingProvider",
			expectedErr: "keptnmetricsproviders.metrics.keptn.sh \"my-provider\" not found",
			analysisDef: analysisDef,
			mockClient:  fake2.NewClient(&analysis, &analysisDef, &template),
		},
		{
			name:        "Success",
			mockClient:  fake2.NewClient(&analysis, &analysisDef, &template, &provider),
			analysisDef: analysisDef,
			providerResult: &metricstypes.ProviderRequest{
				Query: "this is a good query.",
			},
		},
		{
			name:        "Success - provider in same namespace",
			mockClient:  fake2.NewClient(&analysis, &analysisDef, &template2, &provider2),
			analysisDef: analysisDef,
			providerResult: &metricstypes.ProviderRequest{
				Query: "this is a good query.",
			},
		},
		{
			name:        "Success - analysisValueTemplate in same namespace",
			mockClient:  fake2.NewClient(&analysis, &analysisDef2, &template3, &provider),
			analysisDef: analysisDef2,
			providerResult: &metricstypes.ProviderRequest{
				Query: "this is a good query.",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock context for testing
			ctx, cancel := context.WithCancel(context.TODO())

			resultChan := make(chan metricsapi.ProviderResult, 1)

			// Create a mock IObjectivesEvaluator and Logger for testing
			mockEvaluator := &fake.IObjectivesEvaluatorMock{}
			mockLogger := logr.Discard()
			providerChan := make(chan metricstypes.ProviderRequest, 1)
			pool := ProvidersPool{
				IObjectivesEvaluator: mockEvaluator,
				Client:               tc.mockClient,
				log:                  mockLogger,
				Namespace:            "default2",
				Objectives: map[int][]metricsapi.Objective{
					1: tc.analysisDef.Spec.Objectives,
				},
				Analysis: &analysis,
				results:  resultChan,
				cancel:   cancel,
				providers: map[string]chan metricstypes.ProviderRequest{
					"mock-provider": providerChan,
				},
			}

			// Call DispatchToProviders with the test context and example ID
			pool.DispatchToProviders(ctx, 1)

			if tc.expectedErr == "" {
				res := <-providerChan
				require.Equal(t, tc.providerResult.Query, res.Query)
			} else {
				res := <-resultChan
				require.Contains(t, res.ErrMsg, tc.expectedErr)
			}
			pool.StopProviders()
		})
	}
}

func TestProvidersPool_StartProviders(t *testing.T) {

	numJobs := 4
	ctx, cancel := context.WithCancel(context.Background())
	resChan := make(chan metricsapi.ProviderResult)
	// Create a mock IObjectivesEvaluator, Client, and Logger for testing
	mockEvaluator := &fake.IObjectivesEvaluatorMock{
		EvaluateFunc: func(ctx context.Context, providerType string, obj chan metricstypes.ProviderRequest) {
		},
	}

	// Create a ProvidersPool instance with the mock objects
	pool := ProvidersPool{
		IObjectivesEvaluator: mockEvaluator,
		Namespace:            "test-namespace",
		Objectives:           make(map[int][]metricsapi.Objective),
		Analysis:             &metricsapi.Analysis{},
		results:              resChan,
		cancel:               cancel,
		providers:            make(map[string]chan metricstypes.ProviderRequest),
	}

	// Call StartProviders with the test context and example numJobs
	pool.StartProviders(ctx, numJobs)

	// Wait for a short time to allow the goroutines to start
	time.Sleep(time.Millisecond * 100)

	// Assert the expected number of workers (goroutines) were started
	require.Equal(t, 4, len(pool.providers))
	require.Equal(t, numJobs, cap(pool.providers["prometheus"]))
	// Stop the providers after testing
	pool.StopProviders()

}
