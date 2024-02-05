package dummy

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1beta1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestEvaluateQuery_HappyPath(t *testing.T) {
	dummyProvider := &KeptnDummyProvider{
		Log:        ctrl.Log.WithName("testytest"),
		HttpClient: http.Client{},
	}

	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "random",
			Range: &metricsapi.RangeSpec{
				Interval: "5m",
			},
		},
	}
	provider := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: "http://www.randomnumberapi.com/api/v1.0/",
		},
	}

	value, _, err := dummyProvider.EvaluateQuery(context.TODO(), metric, provider)

	require.NoError(t, err)
	require.Equal(t, "dummy provider EvaluateQuery was called with query random", value)
}

func TestFetchAnalysisValue_HappyPath(t *testing.T) {
	dummyProvider := &KeptnDummyProvider{
		Log:        ctrl.Log.WithName("testytest"),
		HttpClient: http.Client{},
	}

	query := "random"
	currentTime := metav1.Time{Time: time.Now()}
	analysis := metricsapi.Analysis{
		Spec: metricsapi.AnalysisSpec{
			Timeframe: metricsapi.Timeframe{
				From: currentTime,
				To:   currentTime,
			},
		},
	}
	provider := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: "http://www.randomnumberapi.com/api/v1.0/",
		},
	}

	value, err := dummyProvider.FetchAnalysisValue(context.TODO(), query, analysis, &provider)

	expected := fmt.Sprintf("dummy provider EvaluateQueryForStep was called with query random from %q to %q", currentTime, currentTime)
	require.NoError(t, err)
	require.Equal(t, expected, value)
}

func TestEvaluateQueryForStep_HappyPath(t *testing.T) {
	dummyProvider := &KeptnDummyProvider{
		Log:        ctrl.Log.WithName("testytest"),
		HttpClient: http.Client{},
	}

	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "random",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}
	provider := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: "http://www.randomnumberapi.com/api/v1.0/",
		},
	}

	intervalDuration, _ := time.ParseDuration("5m")

	stepDuration, _ := time.ParseDuration("1m")
	fromTime := time.Now().Add(-intervalDuration).Unix()
	toTime := time.Now().Unix()
	stepInterval := stepDuration.Milliseconds()

	values, _, err := dummyProvider.EvaluateQueryForStep(context.TODO(), metric, provider)

	expected := fmt.Sprintf("dummy provider EvaluateQueryForStep was called with query random from %q to %q at an interval %q", fromTime, toTime, stepInterval)
	require.NoError(t, err)

	require.Equal(t, expected, values[0])
}
