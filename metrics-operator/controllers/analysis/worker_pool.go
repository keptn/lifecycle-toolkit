package analysis

import (
	"time"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/analysispool_mock.go . IAnalysisPool
type IAnalysisPool interface {
	DispatchAndCollect(ctx context.Context) (map[string]metricsapi.ProviderResult, error)
}

type NewWorkersPoolFactory func(ctx context.Context, analysis *metricsapi.Analysis, objectives []metricsapi.Objective, numWorkers int, c client.Client, log logr.Logger, namespace string) (context.Context, IAnalysisPool)

func NewWorkersPool(ctx context.Context, analysis *metricsapi.Analysis, objectives []metricsapi.Objective, numWorkers int, c client.Client, log logr.Logger, namespace string) (context.Context, IAnalysisPool) {
	numJobs := len(objectives)
	if numJobs <= numWorkers { // do not start useless go routines
		numWorkers = numJobs
	}
	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	providerChans := make(map[string]chan metricstypes.ProviderRequest, len(providers.SupportedProviders))

	assigner := TaskAssigner{tasks: objectives, numWorkers: numWorkers}
	results := make(chan metricsapi.ProviderResult, numJobs)
	evaluator := ObjectivesEvaluator{
		ProviderFactory: providers.NewProvider,
		log:             log,
		Client:          c,
		Analysis:        analysis,
		results:         results,
		cancel:          cancel,
	}
	retriever := ProvidersPool{
		Client:               c,
		log:                  log,
		Analysis:             analysis,
		results:              results,
		Namespace:            namespace,
		Objectives:           assigner.assignTasks(),
		IObjectivesEvaluator: evaluator,
		providers:            providerChans,
		cancel:               cancel,
	}
	return childCtx, WorkersPool{
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

func (aw WorkersPool) DispatchAndCollect(ctx context.Context) (map[string]metricsapi.ProviderResult, error) {
	aw.StartProviders(ctx, aw.numJobs)
	for w := 1; w <= aw.numWorkers; w++ {
		go aw.DispatchToProviders(ctx, w)
	}
	return aw.CollectAnalysisResults(ctx)
}

func (aw WorkersPool) CollectAnalysisResults(ctx context.Context) (map[string]metricsapi.ProviderResult, error) {
	var err error
	results := make(map[string]metricsapi.ProviderResult, aw.numJobs)
	for a := 1; a <= aw.numJobs; a++ {
		select {
		case <-ctx.Done():
			err = errors.New("Collection terminated")
			break
		default:
			res, err2 := aw.GetResult(ctx)
			if err2 != nil {
				err = err2
			} else {
				results[analysis.ComputeKey(res.Objective)] = *res
				if res.ErrMsg != "" {
					err = errors.New(res.ErrMsg)
				}
			}
		}
	}
	aw.StopProviders()
	return results, err
}
