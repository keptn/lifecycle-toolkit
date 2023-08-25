package analysis

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"text/template"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"k8s.io/apimachinery/pkg/types"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/providers_pool_mock.go . IProvidersPool
type IProvidersPool interface {
	StartProviders(ctx context.Context, numJobs int)
	DispatchToProviders(ctx context.Context, id int)
	GetResult() metricstypes.ProviderResult
	StopProviders()
}

type ProvidersPool struct {
	IObjectivesEvaluator
	client.Client
	Log        logr.Logger
	Namespace  string
	Objectives map[int][]metricsapi.Objective
	*metricsapi.Analysis
	results   chan metricstypes.ProviderResult
	providers map[string]chan metricstypes.ProviderRequest
}

func (ps ProvidersPool) StartProviders(ctx context.Context, numJobs int) {
	for _, provider := range providers.SupportedProviders {
		channel := make(chan metricstypes.ProviderRequest, numJobs)
		ps.providers[provider] = channel
		go ps.Evaluate(ctx, provider, channel)
	}

}

func (ps ProvidersPool) DispatchToProviders(ctx context.Context, id int) {

	for _, j := range ps.Objectives[id] {

		ps.Log.Info("worker", "id:", id, "started  job:", j.AnalysisValueTemplateRef.Name)

		templ := &metricsapi.AnalysisValueTemplate{}
		if j.AnalysisValueTemplateRef.Namespace == "" {
			j.AnalysisValueTemplateRef.Namespace = ps.Namespace
		}
		err := ps.Client.Get(ctx,
			types.NamespacedName{
				Name:      j.AnalysisValueTemplateRef.Name,
				Namespace: j.AnalysisValueTemplateRef.Namespace},
			templ,
		)

		if err != nil {
			ps.Log.Error(err, "Failed to get the correct Provider")
			ps.results <- metricstypes.ProviderResult{Objective: j.AnalysisValueTemplateRef, Err: err}
			continue
		}

		providerRef := &metricsapi.KeptnMetricsProvider{}
		if templ.Spec.Provider.Namespace == "" {
			templ.Spec.Provider.Namespace = ps.Namespace
		}
		err = ps.Client.Get(ctx,
			types.NamespacedName{
				Name:      templ.Spec.Provider.Name,
				Namespace: templ.Spec.Provider.Namespace},
			providerRef,
		)

		if err != nil {
			ps.Log.Error(err, "Failed to get Provider")
			ps.results <- metricstypes.ProviderResult{Objective: j.AnalysisValueTemplateRef, Err: err}
			continue
		}

		templatedQuery, err := generateQuery(templ.Spec.Query, ps.Analysis.Spec.Args)
		if err != nil {
			ps.Log.Error(err, "Failed to substitute args in templ")
			ps.results <- metricstypes.ProviderResult{Objective: j.AnalysisValueTemplateRef, Err: err}
			continue
		}
		//send job to provider solver
		ps.providers[providerRef.Spec.Type] <- metricstypes.ProviderRequest{
			Objective: &j,
			Query:     templatedQuery,
			Provider:  providerRef,
		}
	}
}

func (ps ProvidersPool) StopProviders() {
	for _, ch := range ps.providers {
		close(ch)
	}
	close(ps.results)
}

func (ps ProvidersPool) GetResult() metricstypes.ProviderResult {
	res := <-ps.results
	return res
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
