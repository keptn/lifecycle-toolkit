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
	objs := []metricsapi.Objective{{
		Weight:       10,
		KeyObjective: true,
	},
	}

	log := testr.New(t)
	want := WorkersPool{
		numWorkers: 4,
		numJobs:    1,
	}

	// no objectives to evaluate
	_, got := NewWorkersPool(context.TODO(), &analysis, []metricsapi.Objective{}, 4, nil, log, "default")
	require.Equal(t, 0, got.(WorkersPool).numWorkers)
	require.Equal(t, 0, got.(WorkersPool).numJobs)

	_, got = NewWorkersPool(context.TODO(), &analysis, objs, 4, nil, log, "default")
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
		ErrMsg:    "",
	}

	res2 := metricsapi.ProviderResult{
		Objective: metricsapi.ObjectReference{Name: "t2"},
		Value:     "result2",
		ErrMsg:    "",
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
	require.Equal(t, res1, results["t1"])
	require.Equal(t, res2, results["t2"])
}

func TestWorkersPool_CollectAnalysisResultsWithError(t *testing.T) {
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
		ErrMsg:    "",
	}

	res2 := metricsapi.ProviderResult{
		Objective: metricsapi.ObjectReference{Name: "t2"},
		Value:     "result2",
		ErrMsg:    "unexpected error",
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
	require.NotNil(t, err)
	require.Equal(t, res1, results["t1"])
	require.Equal(t, res2, results["t2"])
}

func TestWorkersPool_CollectAnalysisResultsTimeout(t *testing.T) {
	// Create a fake WorkersPool instance for testing
	resChan := make(chan metricsapi.ProviderResult, 2)
	fakePool := WorkersPool{
		IProvidersPool: ProvidersPool{
			results: resChan,
		},
		numJobs: 2,
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)

	fakePool.cancel = cancel

	// Collect the results
	results, err := fakePool.CollectAnalysisResults(ctx)

	// Check the collected results
	require.NotNil(t, err)
	require.Empty(t, results)
}

func TestWorkersPool_CollectAnalysisResultsNoJob(t *testing.T) {
	// Create a fake WorkersPool instance for testing
	resChan := make(chan metricsapi.ProviderResult, 1)
	fakePool := WorkersPool{
		IProvidersPool: ProvidersPool{
			results: resChan,
		},
		numJobs: 0,
	}
	// Collect that func returns
	results, err := fakePool.CollectAnalysisResults(context.TODO())
	_, open := <-resChan
	// Check the collected results
	require.Nil(t, err)
	require.Equal(t, results, map[string]metricsapi.ProviderResult{})
	require.False(t, open)

}
