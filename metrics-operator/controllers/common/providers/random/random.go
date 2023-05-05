package random

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
)

type KeptnRandomProvider struct {
	Log logr.Logger
}

// Every metric provider must implement this EvaluateQuery method.
func (r *KeptnRandomProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {

	// create a context for cancelling the request if it takes too long.
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// since our endpoint only returns a value when hitting /random route
	query := "random"

	result, err := r.Query(ctx, query)
	if err != nil {
		return "", nil, err
	}
	return "", result, nil
}

func (r *KeptnRandomProvider) Query(ctx context.Context, metric string) ([]byte, error) {

	// create a new request with context
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://www.randomnumberapi.com/api/v1.0/random", nil)
	if err != nil {
		r.Log.Error(err, "Error in creating the request")
		return nil, err
	}

	// make an http call using the default client.
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		r.Log.Error(err, "Error in making the request")
		return nil, err
	}

	// parse the response data
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		r.Log.Error(err, "Error in reading the response")
	}

	// return the metric
	return responseData, nil
}
