package dynatrace

import (
	"context"
	"github.com/go-logr/logr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KeptnDynatraceDQLProvider struct {
	Log        logr.Logger
	httpClient http.Client
	k8sClient  client.Client
}

// EvaluateQuery fetches the SLI values from dynatrace provider
func (d *KeptnDynatraceDQLProvider) EvaluateQuery(ctx context.Context, objective klcv1alpha2.Objective, provider klcv1alpha2.KeptnEvaluationProvider) (string, []byte, error) {
	// auth
	// submit DQL
	// attend result
	// parse result
	return "", []byte{}, nil
}
