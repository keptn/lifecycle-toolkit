package analysis

import (
	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers"
	"golang.org/x/net/context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/analysispool_mock.go . IAnalysisPool:MyAnalysisPoolMock
type IAnalysisPool interface {
	DispatchObjectives(ctx context.Context)
	CollectAnalysisResults() map[string]metricstypes.ProviderResult
}

type NewWorkersPoolFactory func(analysis *metricsapi.Analysis, definition *metricsapi.AnalysisDefinition, numWorkers int, c client.Client, log logr.Logger, namespace string) IAnalysisPool

func NewWorkersPool(analysis *metricsapi.Analysis, definition *metricsapi.AnalysisDefinition, numWorkers int, c client.Client, log logr.Logger, namespace string) IAnalysisPool {
	numJobs := len(definition.Spec.Objectives)
	if numJobs <= numWorkers { // do not start useless go routines
		numWorkers = numJobs
	}
	providerChans := make(map[string]chan metricstypes.ProviderRequest, len(providers.SupportedProviders))

	assigner := TaskAssigner{tasks: definition.Spec.Objectives, numWorkers: numWorkers}
	results := make(chan metricstypes.ProviderResult, numJobs)
	evaluator := ObjectivesEvaluator{
		ProviderFactory: providers.NewProvider,
		Log:             log,
		Client:          c,
		Analysis:        analysis,
		results:         results,
	}
	retriever := ProvidersPool{
		Namespace:            namespace,
		Objectives:           assigner.AssignTasks(),
		IObjectivesEvaluator: evaluator,
		providers:            providerChans,
	}
	return WorkersPool{
		numWorkers:     numWorkers,
		numJobs:        numJobs,
		IProvidersPool: retriever,
	}
}

type WorkersPool struct {
	IProvidersPool
	numWorkers int
	numJobs    int
}

func (aw WorkersPool) DispatchObjectives(ctx context.Context) {
	aw.StartProviders(ctx, aw.numJobs)
	for w := 1; w <= aw.numWorkers; w++ {
		go aw.DispatchToProviders(ctx, w)
	}
}

func (aw WorkersPool) CollectAnalysisResults() map[string]metricstypes.ProviderResult {

	results := make(map[string]metricstypes.ProviderResult, aw.numJobs)
	for a := 1; a <= aw.numJobs; a++ {
		res := aw.GetResult()
		results[analysis.ComputeKey(res.Objective)] = res

	}
	aw.StopProviders()
	return results
}
