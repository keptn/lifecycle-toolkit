package dummy

import (
	"context"
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

func TestEvaluateQuery_Error(t *testing.T) {
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
			TargetServer: "",
		},
	}

	_, _, err := dummyProvider.EvaluateQuery(context.TODO(), metric, provider)

	require.Error(t, err)
}

func TestFetchAnalysisValue_HappyPath(t *testing.T) {
	dummyProvider := &KeptnDummyProvider{
		Log:        ctrl.Log.WithName("testytest"),
		HttpClient: http.Client{},
	}

	query := "random"
	analysis := metricsapi.Analysis{
		Spec: metricsapi.AnalysisSpec{
			Timeframe: metricsapi.Timeframe{
				From: metav1.Time{Time: time.Now()},
				To:   metav1.Time{Time: time.Now()},
			},
		},
	}
	provider := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: "http://www.randomnumberapi.com/api/v1.0/",
		},
	}

	value, err := dummyProvider.FetchAnalysisValue(context.TODO(), query, analysis, &provider)

	require.NoError(t, err)
	require.Equal(t, "dummy provider FetchAnalysisValue was called with query random", value)
}

func TestFetchAnalysisValue_Error(t *testing.T) {
	dummyProvider := &KeptnDummyProvider{
		Log:        ctrl.Log.WithName("testytest"),
		HttpClient: http.Client{},
	}

	query := "random"

	analysis := metricsapi.Analysis{
		Spec: metricsapi.AnalysisSpec{
			Timeframe: metricsapi.Timeframe{
				From: metav1.Time{Time: time.Now()},
				To:   metav1.Time{Time: time.Now()},
			},
		},
	}
	provider := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: "",
		},
	}

	_, err := dummyProvider.FetchAnalysisValue(context.TODO(), query, analysis, &provider)

	require.Error(t, err)
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
				Interval: "5m",
			},
		},
	}
	provider := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: "http://www.randomnumberapi.com/api/v1.0/",
		},
	}

	values, _, err := dummyProvider.EvaluateQueryForStep(context.TODO(), metric, provider)

	require.NoError(t, err)
	require.Equal(t, "dummy provider EvaluateQueryForStep was called with query random", values[0])
}

func TestEvaluateQueryForStep_Error(t *testing.T) {
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
			TargetServer: "",
		},
	}

	values, _, err := dummyProvider.EvaluateQueryForStep(context.TODO(), metric, provider)

	require.Error(t, err)
	require.Len(t, values[0], 0)
}
