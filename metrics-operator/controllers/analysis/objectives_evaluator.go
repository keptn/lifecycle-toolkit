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
	results chan metricstypes.ProviderResult
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
		if err == nil {
			value, _, err = provider.FetchAnalysisValue(ctx, o.Query, oe.Analysis.Spec, o.Provider)
		}
		result := metricstypes.ProviderResult{
			Objective: o.Objective.AnalysisValueTemplateRef,
			Value:     value,
			Err:       err,
		}
		oe.Log.Info("provider", "id:", providerType, "finished job:", o.Objective.AnalysisValueTemplateRef.Name, "result:", result)
		oe.results <- result
	}
}
