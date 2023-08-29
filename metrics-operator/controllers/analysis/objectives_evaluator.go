package analysis

import (
	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers"
	"golang.org/x/net/context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/evaluator_mock.go . IObjectivesEvaluator
type IObjectivesEvaluator interface {
	Evaluate(ctx context.Context, providerType string, obj chan metricstypes.ProviderRequest)
}

type ObjectivesEvaluator struct {
	*metricsapi.Analysis
	providers.ProviderFactory
	client.Client
	Log     logr.Logger
	results chan metricsapi.ProviderResult
	cancel  context.CancelFunc
}

func (oe ObjectivesEvaluator) Evaluate(ctx context.Context, providerType string, obj chan metricstypes.ProviderRequest) {
	provider, err := oe.ProviderFactory(providerType, oe.Log, oe.Client)
	if err != nil {
		oe.Log.Error(err, "Failed to get the correct Provider")
		oe.cancel()
		return
	}
	for o := range obj {
		value := ""
		var strErr string
		if err == nil {
			value, err = provider.FetchAnalysisValue(ctx, o.Query, oe.Analysis.Spec, o.Provider)
		}
		if err != nil {
			strErr = err.Error()
		}
		result := metricsapi.ProviderResult{
			Objective: o.Objective.AnalysisValueTemplateRef,
			Value:     value,
			Err:       strErr,
		}
		oe.Log.Info("provider", "id:", providerType, "finished job:", o.Objective.AnalysisValueTemplateRef.Name, "result:", result)
		oe.results <- result
	}
}
