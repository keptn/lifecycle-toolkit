package analysis

import (
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis"
	"strings"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/analysispool_mock.go . IAnalysisPool:MyAnalysisPoolMock
type IAnalysisPool interface {
	DispatchObjectives(ctx context.Context)
	CollectAnalysisResults() map[string]string
}

type NewWorkersPoolFactory func(analysis *metricsapi.Analysis, definition *metricsapi.AnalysisDefinition, numWorkers int, c client.Client, log logr.Logger) IAnalysisPool

func NewWorkersPool(analysis *metricsapi.Analysis, definition *metricsapi.AnalysisDefinition, numWorkers int, c client.Client, log logr.Logger) IAnalysisPool {
	numJobs := len(definition.Spec.Objectives)
	if numJobs <= numWorkers { // do not start useless go routines
		numWorkers = numJobs
	}
	providerChans := make(map[string]chan metricstypes.ProviderRequest, len(providers.SupportedProviders))

	return WorkersPool{
		Analysis:   analysis,
		Objectives: assignTasks(definition.Spec.Objectives, numWorkers),
		Client:     c,
		Log:        log,
		numWorkers: numWorkers,
		numJobs:    numJobs,
		providers:  providerChans,
		results:    make(chan metricstypes.ProviderResult, numJobs),
	}
}

type WorkersPool struct {
	*metricsapi.Analysis
	Objectives map[int][]metricsapi.Objective
	client.Client
	Log        logr.Logger
	numWorkers int
	numJobs    int
	providers  map[string]chan metricstypes.ProviderRequest
	results    chan metricstypes.ProviderResult
}

func assignTasks(tasks []metricsapi.Objective, numWorkers int) map[int][]metricsapi.Objective {
	taskMap := make(map[int][]metricsapi.Objective, numWorkers)
	for i, task := range tasks {
		workerID := (i % numWorkers) + 1
		taskMap[workerID] = append(taskMap[workerID], task)
	}
	return taskMap
}

func (aw WorkersPool) DispatchObjectives(ctx context.Context) {

	for _, provider := range providers.SupportedProviders {
		channel := make(chan metricstypes.ProviderRequest, aw.numJobs)
		aw.providers[provider] = channel
		go aw.Evaluate(ctx, provider, channel)
	}

	for w := 1; w <= aw.numWorkers; w++ {
		go aw.RetrieveProvider(ctx, w)
	}
}

func (aw WorkersPool) CollectAnalysisResults() map[string]string {

	results := make(map[string]string, aw.numJobs)
	for a := 1; a <= aw.numJobs; a++ {
		res := <-aw.results
		aw.Log.Info("collected result")
		// Making sure error gets propagated
		if res.Err != "" {
			results[analysis.ComputeKey(res.Objective)] = res.Err
		} else {
			results[analysis.ComputeKey(res.Objective)] = res.Value
		}
	}
	aw.stopProviders()
	close(aw.results)
	return results
}

func (aw WorkersPool) stopProviders() {
	for _, ch := range aw.providers {
		close(ch)
	}
}

func generateQuery(query string, selectors map[string]string) string {
	for key, value := range selectors {
		query = strings.Replace(query, "$"+strings.ToUpper(key), value, -1)
	}
	return query
}

func (aw WorkersPool) RetrieveProvider(ctx context.Context, id int) {

	for _, j := range aw.Objectives[id] {

		aw.Log.Info("worker", "id:", id, "started  job:", j.AnalysisValueTemplateRef.Name)

		template := &metricsapi.AnalysisValueTemplate{}
		err := aw.Client.Get(ctx,
			types.NamespacedName{
				Name:      j.AnalysisValueTemplateRef.Name,
				Namespace: j.AnalysisValueTemplateRef.Namespace},
			template,
		)

		if err != nil {
			aw.Log.Error(err, "Failed to get the correct Provider")
			aw.results <- metricstypes.ProviderResult{Objective: j.AnalysisValueTemplateRef, Err: err.Error()}
			continue
		}

		providerRef := &metricsapi.KeptnMetricsProvider{}
		err = aw.Client.Get(ctx,
			types.NamespacedName{
				Name:      template.Spec.Provider.Name,
				Namespace: template.Spec.Provider.Namespace},
			providerRef,
		)

		if err != nil {
			aw.Log.Error(err, "Failed to get Provider")
			aw.results <- metricstypes.ProviderResult{Objective: j.AnalysisValueTemplateRef, Err: err.Error()}
			continue
		}

		templatedQuery := generateQuery(template.Spec.Query, aw.Analysis.Spec.Args)

		//send job to provider solver
		aw.providers[providerRef.Spec.Type] <- metricstypes.ProviderRequest{
			Objective: &j,
			Query:     templatedQuery,
			Provider:  providerRef,
		}
	}
}

// TODO add timeout and spec save of unfinished analysis
func (aw WorkersPool) Evaluate(ctx context.Context, providerType string, obj chan metricstypes.ProviderRequest) {
	provider, err := providers.NewProvider(providerType, aw.Log, aw.Client)
	if err != nil {
		aw.Log.Error(err, "Failed to get the correct Provider")
	}
	for o := range obj {
		value := ""
		raw := []byte{}
		if err == nil {
			value, raw, err = provider.RunAnalysis(ctx, o.Query, aw.Analysis.Spec, o.Provider)
		}
		result := metricstypes.ProviderResult{
			Objective: o.Objective.AnalysisValueTemplateRef,
			Value:     value,
			Raw:       raw,
			Err:       err.Error(),
		}
		aw.Log.Info("provider", "id:", providerType, "finished job:", o.Objective.AnalysisValueTemplateRef.Name, "result:", result)
		aw.results <- result
	}
}
