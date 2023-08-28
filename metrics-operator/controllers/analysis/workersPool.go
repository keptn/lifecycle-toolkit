package analysis

import (
	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/analysispool_mock.go . IAnalysisPool:MyAnalysisPoolMock
type IAnalysisPool interface {
	DispatchAndCollect(ctx context.Context) (map[string]string, error)
}

type NewWorkersPoolFactory func(ctx context.Context, analysis *metricsapi.Analysis, definition *metricsapi.AnalysisDefinition, numWorkers int, c client.Client, log logr.Logger, namespace string) (context.Context, IAnalysisPool)

func NewWorkersPool(ctx context.Context, analysis *metricsapi.Analysis, definition *metricsapi.AnalysisDefinition, numWorkers int, c client.Client, log logr.Logger, namespace string) (context.Context, IAnalysisPool) {
	numJobs := len(definition.Spec.Objectives)
	if numJobs <= numWorkers { // do not start useless go routines
		numWorkers = numJobs
	}
	ctx, cancel := context.WithCancel(context.Background())

	providerChans := make(map[string]chan metricstypes.ProviderRequest, len(providers.SupportedProviders))

	assigner := TaskAssigner{tasks: definition.Spec.Objectives, numWorkers: numWorkers}
	results := make(chan metricstypes.ProviderResult, numJobs)
	evaluator := ObjectivesEvaluator{
		ProviderFactory: providers.NewProvider,
		Log:             log,
		Client:          c,
		Analysis:        analysis,
		results:         results,
		cancel:          cancel,
	}
	retriever := ProvidersPool{
		Client:               c,
		Log:                  log,
		Analysis:             analysis,
		results:              results,
		Namespace:            namespace,
		Objectives:           assigner.AssignTasks(),
		IObjectivesEvaluator: evaluator,
		providers:            providerChans,
		cancel:               cancel,
	}
	return ctx, WorkersPool{
		numWorkers:     numWorkers,
		numJobs:        numJobs,
		cancel:         cancel,
		IProvidersPool: retriever,
	}
}

type WorkersPool struct {
	IProvidersPool
	numWorkers int
	numJobs    int
	cancel     context.CancelFunc
}

func (aw WorkersPool) DispatchAndCollect(ctx context.Context) (map[string]string, error) {
	aw.StartProviders(ctx, aw.numJobs)
	for w := 1; w <= aw.numWorkers; w++ {
		go aw.DispatchToProviders(ctx, w)
	}
	return aw.CollectAnalysisResults(ctx)
}

func (aw WorkersPool) CollectAnalysisResults(ctx context.Context) (map[string]string, error) {
	var err error
	results := make(map[string]string, aw.numJobs)
loop:
	for a := 1; a <= aw.numJobs; a++ {
		select {
		case <-ctx.Done():
			err = errors.New("Collection terminated")
			break loop
		default:
			res := aw.GetResult()
			results[analysis.ComputeKey(res.Objective)] = res.Value
			if res.Err != nil {
				err = res.Err
				aw.cancel()
				break loop
			}
		}
	}
	aw.StopProviders()
	return results, err
}
