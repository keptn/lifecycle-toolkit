package dynatrace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KeptnDynatraceProvider struct {
	Log        logr.Logger
	httpClient http.Client
	k8sClient  client.Client
}

type DynatraceResponse struct {
	TotalCount int               `json:"totalCount"`
	Resolution string            `json:"resolution"`
	Result     []DynatraceResult `json:"result"`
}

type DynatraceResult struct {
	MetricId string          `json:"metricId"`
	Data     []DynatraceData `json:"data"`
}

type DynatraceData struct {
	Timestamps []int64    `json:"timestamps"`
	Values     []*float64 `json:"values"`
}

// EvaluateQuery fetches the SLI values from dynatrace provider
func (d *KeptnDynatraceProvider) EvaluateQuery(ctx context.Context, objective klcv1alpha2.Objective, provider klcv1alpha2.KeptnEvaluationProvider) (string, []byte, error) {
	qURL := provider.Spec.TargetServer + "/api/v2/metrics/query?metricSelector=" + objective.Query

	d.Log.Info("Running query: " + qURL)
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", qURL, nil)
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return "", nil, err
	}

	token, err := getDTApiToken(ctx, provider, d.k8sClient)
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("Authorization", "Api-Token "+token)
	res, err := d.httpClient.Do(req)
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return "", nil, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			d.Log.Error(err, "Could not close request body")
		}
	}()

	// we ignore the error here because we fail later while unmarshalling
	b, _ := io.ReadAll(res.Body)
	result := DynatraceResponse{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		d.Log.Error(err, "Error while parsing response")
		return "", nil, err
	}

	r := fmt.Sprintf("%f", d.getSingleValue(result))
	return r, b, nil
}

func (d *KeptnDynatraceProvider) getSingleValue(result DynatraceResponse) float64 {
	var sum float64 = 0
	var count uint64 = 0
	for _, r := range result.Result {
		for _, points := range r.Data {
			for _, v := range points.Values {
				if v != nil {
					sum += *v
					count++
				}
			}
		}
	}
	if count < 1 {
		// cannot dive by zero
		return 0
	}
	return sum / float64(count)
}

func getDTApiToken(ctx context.Context, provider klcv1alpha2.KeptnEvaluationProvider, k8sClient client.Client) (string, error) {
	if !provider.HasSecretDefined() {
		return "", errors.New("the SecretKeyRef property with the DT API token is missing")
	}
	dtCredsSecret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: provider.Spec.SecretKeyRef.Name, Namespace: provider.Namespace}, dtCredsSecret); err != nil {
		return "", err
	}

	apiToken := dtCredsSecret.Data[provider.Spec.SecretKeyRef.Key]
	if len(apiToken) == 0 {
		return "", fmt.Errorf("secret contains invalid key %s", provider.Spec.SecretKeyRef.Key)
	}
	return string(apiToken), nil
}
