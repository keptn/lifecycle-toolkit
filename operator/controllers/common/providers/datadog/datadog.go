package datadog

import (
	"context"
	"fmt"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/go-logr/logr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"net/http"
	"strconv"
	"time"
)

type KeptnDataDogProvider struct {
	Log        logr.Logger
	HttpClient http.Client
}

func (k KeptnDataDogProvider) EvaluateQuery(ctx context.Context, objective klcv1alpha2.Objective, provider klcv1alpha2.KeptnEvaluationProvider) (string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// TODO: get DD_API_KEY and DD_APP_KEY from kubernetes secret
	// TODO: patch the context with the keys
	// Ref: https://github.com/DataDog/datadog-api-client-go#getting-started

	ctx = context.WithValue(
		ctx,
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: "DD_KEY",
			},
			"appKeyAuth": {
				Key: "DD_APP_KEY",
			},
		},
	)

	fromTime := time.Now().AddDate(0, 0, -1)
	queryTime := time.Now()

	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV1.NewMetricsApi(apiClient)

	resp, _, err := api.QueryMetrics(
		ctx,
		fromTime.Unix(),
		queryTime.Unix(),
		objective.Query,
	)
	if err != nil {
		k.Log.Error(err, "Error while creating request")
		return "", nil, err
	}
	if len(resp.Series) == 0 {
		k.Log.Info("No values in query result")
		return "", nil, fmt.Errorf("no values in query result")
	}
	points := (resp.Series)[0].Pointlist
	value := strconv.FormatFloat(*points[len(points)-1][1], 'g', 5, 64)
	return value, []byte(value), nil

}
