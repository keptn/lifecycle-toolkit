package analysis

import (
	"context"
	"testing"
	"time"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	analysistypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestResultsReporter(t *testing.T) {
	res := make(chan analysistypes.AnalysisCompletion)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reporter := GetResultsReporter(ctx, res)

	require.NotNil(t, reporter)

	go func() {
		res <- analysistypes.AnalysisCompletion{
			Result: analysistypes.AnalysisResult{
				ObjectiveResults: []analysistypes.ObjectiveResult{
					{
						Result: analysistypes.TargetResult{},
						Objective: metricsapi.Objective{
							AnalysisValueTemplateRef: metricsapi.ObjectReference{
								Name:      "my-av",
								Namespace: "my-namespace",
							},
							Weight:       2,
							KeyObjective: true,
						},
						Value: 10,
						Score: 0,
						Error: nil,
					},
				},
				TotalScore:   2,
				MaximumScore: 2,
				Pass:         true,
				Warning:      false,
			},
			Analysis: metricsapi.Analysis{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-analysis",
					Namespace: "my-namespace",
				},
				Spec: metricsapi.AnalysisSpec{
					Timeframe: metricsapi.Timeframe{
						From: metav1.Now(),
						To:   metav1.Now(),
					},
				},
			},
		}
	}()

	var analysisMetric *io_prometheus_client.MetricFamily
	var objectiveMetric *io_prometheus_client.MetricFamily

	require.Eventually(t, func() bool {
		gather, err := prometheus.DefaultGatherer.Gather()
		if err != nil {
			return false
		}
		for i, metric := range gather {
			if metric.Name == nil {
				continue
			}
			if *metric.Name == analysisResultMetricName {
				analysisMetric = gather[i]
			}
			if *metric.Name == objectiveResultMetricName {
				objectiveMetric = gather[i]
			}
		}
		return analysisMetric != nil && objectiveMetric != nil
	}, 10*time.Second, 100*time.Millisecond)

	require.NotNil(t, analysisMetric)
	require.Len(t, analysisMetric.Metric[0].Label, 4)
	require.Len(t, analysisMetric.Metric, 1)
	require.EqualValues(t, 100.0, *analysisMetric.Metric[0].Gauge.Value)

	require.NotNil(t, objectiveMetric)
	require.Len(t, objectiveMetric.Metric[0].Label, 8)
	require.Len(t, objectiveMetric.Metric, 1)
	require.Equal(t, 10.0, *objectiveMetric.Metric[0].Gauge.Value)
}
