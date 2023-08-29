package analysis

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
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

	_, got := NewWorkersPool(context.TODO(), &analysis, &def, 4, nil, log, "default")
	//make sure never to create more workers than needed
	require.Equal(t, want.numJobs, got.(WorkersPool).numWorkers)
	//make sure all objectives are processed
	require.Equal(t, want.numJobs, got.(WorkersPool).numJobs)
}

func TestWorkersPool_CollectAnalysisResults(t *testing.T) {
	// Create a fake WorkersPool instance for testing
	resChan := make(chan metricsapi.ProviderResult, 2)
	fakePool := WorkersPool{
		IProvidersPool: ProvidersPool{
			results: resChan,
		},
		numJobs: 2,
	}

	res1 := metricsapi.ProviderResult{
		Objective: metricsapi.ObjectReference{Name: "t1"},
		Value:     "result1",
		Err:       "",
	}

	res2 := metricsapi.ProviderResult{
		Objective: metricsapi.ObjectReference{Name: "t2"},
		Value:     "result2",
		Err:       "",
	}

	// Create and send mock results to the results channel
	go func() {
		time.Sleep(time.Second)
		resChan <- res1
		resChan <- res2
	}()

	// Collect the results
	results, err := fakePool.CollectAnalysisResults(context.TODO())

	// Check the collected results
	require.Nil(t, err)
	require.Equal(t, res1.Value, results["t1"])
	require.Equal(t, res2.Value, results["t2"])
}
