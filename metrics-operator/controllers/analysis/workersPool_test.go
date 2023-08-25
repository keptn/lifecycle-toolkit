package analysis

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers"
	fake2 "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/fake"
	"github.com/stretchr/testify/require"
)

func TestNewWorkerPool(t *testing.T) {

	analysis := metricsapi.Analysis{
		Spec: metricsapi.AnalysisSpec{Args: map[string]string{"hi": "there"}},
	}
	def := metricsapi.AnalysisDefinition{
		Spec: metricsapi.AnalysisDefinitionSpec{
			Objectives: []metricsapi.Objective{{
				Weight:       10,
				KeyObjective: true,
			},
			},
		},
	}

	log := testr.New(t)
	want := WorkersPool{
		Objectives: map[int][]metricsapi.Objective{1: def.Spec.Objectives},
		Analysis:   &analysis,
		numWorkers: 4,
		numJobs:    1,
	}

	got := NewWorkersPool(&analysis, &def, 4, nil, log, "default")
	require.Equal(t, want.Objectives, got.(WorkersPool).Objectives)
	require.Equal(t, want.Analysis, got.(WorkersPool).Analysis)
	//make sure never to create more workers than needed
	require.Equal(t, want.numJobs, got.(WorkersPool).numWorkers)
	//make sure all objectives are processed
	require.Equal(t, want.numJobs, got.(WorkersPool).numJobs)
}

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

func TestWorkersPool_CollectAnalysisResults(t *testing.T) {
	// Create a fake WorkersPool instance for testing
	fakePool := WorkersPool{
		Analysis: &metricsapi.Analysis{},
		results:  make(chan metricstypes.ProviderResult, 2),
		numJobs:  2,
		Log:      logr.Discard(),
	}

	res1 := metricstypes.ProviderResult{
		Objective: metricsapi.ObjectReference{Name: "t1"},
		Value:     "result1",
		Err:       nil,
	}

	res2 := metricstypes.ProviderResult{
		Objective: metricsapi.ObjectReference{Name: "t2"},
		Value:     "result2",
		Err:       nil,
	}

	// Create and send mock results to the results channel
	go func() {
		time.Sleep(time.Second)
		fakePool.results <- res1
		fakePool.results <- res2
	}()

	// Collect the results
	results := fakePool.CollectAnalysisResults()

	// Check the collected results
	require.Equal(t, res1, results["t1"])
	require.Equal(t, res2, results["t2"])
}

func TestAssignTasks(t *testing.T) {
	tests := []struct {
		name       string
		tasks      []metricsapi.Objective
		numWorkers int
		expected   map[int][]metricsapi.Objective
	}{
		{
			name:       "Equal Workers and Tasks",
			tasks:      []metricsapi.Objective{{}, {}, {}},
			numWorkers: 3,
			expected: map[int][]metricsapi.Objective{
				1: {{}},
				2: {{}},
				3: {{}},
			},
		},
		{
			name:       "More Workers than Tasks",
			tasks:      []metricsapi.Objective{{}, {}},
			numWorkers: 3,
			expected: map[int][]metricsapi.Objective{
				1: {{}},
				2: {{}},
			},
		},
		{
			name:       "More Tasks",
			tasks:      []metricsapi.Objective{{}, {}, {}},
			numWorkers: 2,
			expected: map[int][]metricsapi.Objective{
				1: {{}, {}},
				2: {{}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := assignTasks(test.tasks, test.numWorkers)
			require.Equal(t, test.expected, result)
		})
	}
}

func TestEvaluate(t *testing.T) {
	fakeProvider := fake2.KeptnSLIProviderMock{
		FetchAnalysisValueFunc: func(ctx context.Context, query string, spec metricsapi.AnalysisSpec, provider *metricsapi.KeptnMetricsProvider) (string, []byte, error) {
			return "10", []byte{}, nil
		},
	}

	fakePool := WorkersPool{
		Analysis: &metricsapi.Analysis{
			Spec: metricsapi.AnalysisSpec{
				Args: map[string]string{"a": "b"},
				AnalysisDefinition: metricsapi.ObjectReference{
					Name:      "a",
					Namespace: "d",
				},
			},
		},
		ProviderFactory: func(providerType string, log logr.Logger, k8sClient client.Client) (providers.KeptnSLIProvider, error) {
			return &fakeProvider, nil
		},
		Client:  fake.NewClient(),
		Log:     logr.Discard(),
		results: make(chan metricstypes.ProviderResult, 1),
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

	fakePool.Evaluate(ctx, providerType, obj)
	close(fakePool.results)
	for res := range fakePool.results {
		require.Equal(t, "mytemp", res.Objective.Name)
		require.Contains(t, "10", res.Value)
		require.Nil(t, res.Err)
	}

}
