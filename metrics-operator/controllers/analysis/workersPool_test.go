package analysis

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func TestRetrieveProvider(t *testing.T) {
	// Create a fake client
	client := fake.NewClient()

	// Create a fake WorkersPool instance for testing
	fakePool := WorkersPool{
		Analysis: &metricsapi.Analysis{},
		Objectives: map[int][]metricsapi.Objective{
			1: {{AnalysisValueTemplateRef: metricsapi.ObjectReference{Name: "template-1"}}},
			2: {{AnalysisValueTemplateRef: metricsapi.ObjectReference{Name: "template-2"}}},
		},
		Client: client,
		Log:    logr.Discard(),
	}

	// Simulate context
	ctx := context.TODO()

	// Test with an existing provider
	existingProvider := &metricsapi.KeptnMetricsProvider{
		ObjectMeta: metav1.ObjectMeta{Name: "provider-1"},
	}
	require.NoError(t, client.Create(ctx, existingProvider))

	fakePool.RetrieveProvider(ctx, 1)
	fakePool.RetrieveProvider(ctx, 2)

	// require that the error handling works as expected
}

func TestEvaluate(t *testing.T) { //TODO abstract provider
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
		Client:  fake.NewClient(),
		Log:     logr.Discard(),
		results: make(chan metricstypes.ProviderResult, 1),
	}

	// Simulate context
	ctx := context.TODO()

	// Test with a valid provider type
	providerType := "fakeProvider"
	obj := make(chan metricstypes.ProviderRequest, 1)
	go fakePool.Evaluate(ctx, providerType, obj)

	obj <- metricstypes.ProviderRequest{
		Objective: &metricsapi.Objective{},
		Query:     "query_fake_metric",
		Provider:  &metricsapi.KeptnMetricsProvider{Spec: metricsapi.KeptnMetricsProviderSpec{Type: providerType}},
	}
	close(obj)

	for res := range fakePool.results {
		require.Equal(t, "provider fakeProvider not supported", res.Err)
	}

	close(fakePool.results)

}
