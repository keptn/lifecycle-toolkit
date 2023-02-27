package datadog

import (
	"context"
	"encoding/json"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/go-logr/logr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"net/http"
	"time"
)

type KeptnDataDogProvider struct {
	Log        logr.Logger
	HttpClient http.Client
}

func (k KeptnDataDogProvider) EvaluateQuery(ctx context.Context, objective klcv1alpha2.Objective, provider klcv1alpha2.KeptnEvaluationProvider) (string, []byte, error) {

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	k.Log.Info("Running query: " + provider.Spec.TargetServer + "/api/v1/slo")

	cfg := datadog.NewConfiguration()
	client := datadog.NewAPIClient(cfg)
	api := datadogV1.NewServiceLevelObjectivesApi(client)
	resp, _, err := api.ListSLOs(ctx, *datadogV1.NewListSLOsOptionalParameters())

	if err != nil {
		return "", nil, err
	}

	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	return string(responseContent), responseContent, nil
}
