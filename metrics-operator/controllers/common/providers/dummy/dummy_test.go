package dummy

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1beta1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestEvaluateQuery_HappyPath(t *testing.T) {
	// Create a dummy HTTP server that responds with a predefined payload
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("42")) // Respond with a dummy value
	}))
	defer svr.Close()

	// Create a new instance of KeptnDummyProvider
	dummyProvider := &KeptnDummyProvider{
		Log:        ctrl.Log.WithName("testytest"),
		HttpClient: http.Client{},
	}

	// Create a sample metric and provider
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "random",
		},
	}
	provider := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: svr.URL,
		},
	}

	// Call the EvaluateQuery method
	value, _, err := dummyProvider.EvaluateQuery(context.TODO(), metric, provider)

	// Check if the result is as expected
	require.NoError(t, err)
	require.Equal(t, "42", value)
}

func TestEvaluateQuery_Error(t *testing.T) {
	// Create a dummy HTTP server that always returns an error
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "mock error", http.StatusInternalServerError)
	}))
	defer svr.Close()

	// Create a new instance of KeptnDummyProvider
	dummyProvider := &KeptnDummyProvider{
		Log:        ctrl.Log.WithName("testytest"),
		HttpClient: http.Client{},
	}

	// Create a sample metric and provider
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "random",
		},
	}
	provider := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: svr.URL,
		},
	}

	// Call the EvaluateQuery method
	_, _, err := dummyProvider.EvaluateQuery(context.TODO(), metric, provider)

	// Check if an error occurred
	require.Error(t, err)
	require.Contains(t, err.Error(), "mock error")
}

func TestFetchAnalysisValue_HappyPath(t *testing.T) {
	// Create a new instance of KeptnDummyProvider
	dummyProvider := &KeptnDummyProvider{
		Log:        ctrl.Log.WithName("testytest"),
		HttpClient: http.Client{},
	}

	// Create a sample query and analysis
	query := "random"
	analysis := metricsapi.Analysis{
		Spec: metricsapi.AnalysisSpec{
			Timeframe: metricsapi.Timeframe{
				From: metav1.Time{Time: time.Now()},
				To:   metav1.Time{Time: time.Now()},
			},
		},
	}
	// Create a sample provider
	provider := &metricsapi.KeptnMetricsProvider{}

	// Call the FetchAnalysisValue method
	value, err := dummyProvider.FetchAnalysisValue(context.TODO(), query, analysis, provider)

	// Check if the result is as expected
	require.NoError(t, err)
	require.NotEmpty(t, value)
}

func TestFetchAnalysisValue_Error(t *testing.T) {
	// Create a new instance of KeptnDummyProvider
	dummyProvider := &KeptnDummyProvider{
		Log:        ctrl.Log.WithName("testytest"),
		HttpClient: http.Client{},
	}

	// Create a sample query and analysis
	query := "random"

	analysis := metricsapi.Analysis{
		Spec: metricsapi.AnalysisSpec{
			Timeframe: metricsapi.Timeframe{
				From: metav1.Time{Time: time.Now()},
				To:   metav1.Time{Time: time.Now()},
			},
		},
	}
	// Create a sample provider that will return an error
	provider := &metricsapi.KeptnMetricsProvider{}

	// Call the FetchAnalysisValue method
	_, err := dummyProvider.FetchAnalysisValue(context.TODO(), query, analysis, provider)

	// Check if an error occurred
	require.Error(t, err)
	///require.True(t, errors.Is(err, context.DeadlineExceeded))
}
