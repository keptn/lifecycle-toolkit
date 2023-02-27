package datadog

import (
	"context"
	"fmt"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"net/http"
	"os"
	"testing"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func Test_datadog(t *testing.T) {

	t.Run("datadog-test", func(t *testing.T) {

		kpp := KeptnDataDogProvider{
			HttpClient: http.Client{},
			Log:        ctrl.Log.WithName("testytest"),
		}
		obj := klcv1alpha2.Objective{
			Query: "garbage",
		}
		p := klcv1alpha2.KeptnEvaluationProvider{
			Spec: klcv1alpha2.KeptnEvaluationProviderSpec{
				SecretKeyRef: v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "myapitoken",
					},
					Key: "mykey",
				},
				TargetServer: "http://mockserver", // TODO: pass mock server url
			},
		}
		ctx := context.WithValue(
			context.Background(),
			datadog.ContextAPIKeys,
			map[string]datadog.APIKey{
				"apiKeyAuth": {
					Key: "DD_CLIENT_API_KEY",
				},
				"appKeyAuth": {
					Key: "DD_CLIENT_APP_KEY",
				},
			},
		)
		r, _, _ := kpp.EvaluateQuery(ctx, obj, p)
		fmt.Fprintf(os.Stdout, r)
		require.Equal(t, "0", "1")

	})

}
