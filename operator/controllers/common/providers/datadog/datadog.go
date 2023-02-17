package datadog

import (
	"context"
	"github.com/go-logr/logr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"net/http"
)

type KeptnDataDogProvider struct {
	Log        logr.Logger
	HttpClient http.Client
}

func (k KeptnDataDogProvider) EvaluateQuery(ctx context.Context, objective klcv1alpha2.Objective, provider klcv1alpha2.KeptnEvaluationProvider) (string, []byte, error) {
	//TODO implement me
	panic("implement me")
}
