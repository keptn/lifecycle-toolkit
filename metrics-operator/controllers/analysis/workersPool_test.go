package analysis

import (
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
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
		numWorkers: 4,
		numJobs:    1,
	}

	got := NewWorkersPool(&analysis, &def, 4, nil, log, "default")
	//make sure never to create more workers than needed
	require.Equal(t, want.numJobs, got.(WorkersPool).numWorkers)
	//make sure all objectives are processed
	require.Equal(t, want.numJobs, got.(WorkersPool).numJobs)
}

func TestWorkersPool_CollectAnalysisResults(t *testing.T) {
	// Create a fake WorkersPool instance for testing
	resChan := make(chan metricstypes.ProviderResult, 2)
	fakePool := WorkersPool{
		IProvidersPool: ProvidersPool{
			results: resChan,
		},
		numJobs: 2,
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
		resChan <- res1
		resChan <- res2
	}()

	// Collect the results
	results := fakePool.CollectAnalysisResults()

	// Check the collected results
	require.Equal(t, res1, results["t1"])
	require.Equal(t, res2, results["t2"])
}
