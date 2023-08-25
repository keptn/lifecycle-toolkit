package analysis

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/types"
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

	return WorkersPool{
		Analysis:   analysis,
		Objectives: assignTasks(definition.Spec.Objectives, numWorkers),
		Client:     c,
		Log:        log,
		Namespace:  namespace,
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
	Namespace  string
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

func (aw WorkersPool) CollectAnalysisResults() map[string]metricstypes.ProviderResult {

	results := make(map[string]metricstypes.ProviderResult, aw.numJobs)
	for a := 1; a <= aw.numJobs; a++ {
		res := <-aw.results
		aw.Log.Info("collected result")
		// Making sure error gets propagated
		results[analysis.ComputeKey(res.Objective)] = res

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

func generateQuery(query string, selectors map[string]string) (string, error) {
	tmpl, err := template.New("").Parse(query)
	if err != nil {
		return "", fmt.Errorf("could not create a template: %w", err)
	}

	var resultBuf bytes.Buffer
	err = tmpl.Execute(&resultBuf, selectors)
	if err != nil {
		return "", fmt.Errorf("could not template the args: %w", err)
	}

	return resultBuf.String(), nil
}

func (aw WorkersPool) RetrieveProvider(ctx context.Context, id int) {

	for _, j := range aw.Objectives[id] {

		aw.Log.Info("worker", "id:", id, "started  job:", j.AnalysisValueTemplateRef.Name)

		template := &metricsapi.AnalysisValueTemplate{}
		if j.AnalysisValueTemplateRef.Namespace == "" {
			j.AnalysisValueTemplateRef.Namespace = aw.Namespace
		}
		err := aw.Client.Get(ctx,
			types.NamespacedName{
				Name:      j.AnalysisValueTemplateRef.Name,
				Namespace: j.AnalysisValueTemplateRef.Namespace},
			template,
		)

		if err != nil {
			aw.Log.Error(err, "Failed to get the correct Provider")
			aw.results <- metricstypes.ProviderResult{Objective: j.AnalysisValueTemplateRef, Err: err}
			continue
		}

		providerRef := &metricsapi.KeptnMetricsProvider{}
		if template.Spec.Provider.Namespace == "" {
			template.Spec.Provider.Namespace = aw.Namespace
		}
		err = aw.Client.Get(ctx,
			types.NamespacedName{
				Name:      template.Spec.Provider.Name,
				Namespace: template.Spec.Provider.Namespace},
			providerRef,
		)

		if err != nil {
			aw.Log.Error(err, "Failed to get Provider")
			aw.results <- metricstypes.ProviderResult{Objective: j.AnalysisValueTemplateRef, Err: err}
			continue
		}

		templatedQuery, err := generateQuery(template.Spec.Query, aw.Analysis.Spec.Args)
		if err != nil {
			aw.Log.Error(err, "Failed to substitute args in template")
			aw.results <- metricstypes.ProviderResult{Objective: j.AnalysisValueTemplateRef, Err: err}
			continue
		}
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
		if err == nil {
			value, _, err = provider.FetchAnalysisValue(ctx, o.Query, aw.Analysis.Spec, o.Provider)
		}
		result := metricstypes.ProviderResult{
			Objective: o.Objective.AnalysisValueTemplateRef,
			Value:     value,
			Err:       err,
		}
		aw.Log.Info("provider", "id:", providerType, "finished job:", o.Objective.AnalysisValueTemplateRef.Name, "result:", result)
		aw.results <- result
	}
}
