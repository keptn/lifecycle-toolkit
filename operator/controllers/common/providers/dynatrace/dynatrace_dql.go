package dynatrace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"io"
	"net/http"
	"net/url"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

type KeptnDynatraceDQLProvider struct {
	Log        logr.Logger
	httpClient http.Client
	k8sClient  client.Client
}

type DynatraceOAuthResponse struct {
	accessToken string `json:"access_token"`
}

type DynatraceDQLHandler struct {
	requestToken string `json:"requestToken"`
}

type DynatraceDQLResult struct {
	state  string    `json:"state"`
	result DQLResult `json:"result,omitempty"`
}

type DQLResult struct {
	records []DQLRecord `json:"records"`
}

type DQLRecord struct {
	value DQLMetric `json:"value"`
}

type DQLMetric struct {
	count int64 `json:"count"`
	sum   int64 `json:"sum"`
	min   int64 `json:"min"`
	avg   int64 `json:"avg"`
	max   int64 `json:"max"`
}

// EvaluateQuery fetches the SLI values from dynatrace provider
func (d *KeptnDynatraceDQLProvider) EvaluateQuery(ctx context.Context, objective klcv1alpha2.Objective, provider klcv1alpha2.KeptnEvaluationProvider) (string, []byte, error) {
	// auth
	jwt, err := d.doOAuth(ctx, provider)
	if err != nil {
		return "", []byte{}, nil
	}
	// submit DQL
	dqlHandler, err := d.postDQL(ctx, jwt, objective, provider)
	if err != nil {
		d.Log.Error(err, "", dqlHandler)
		return "", nil, err
	}
	// attend result
	// parse result
	return "", []byte{}, nil
}

func (d *KeptnDynatraceDQLProvider) getScopes() string {
	return "storage:events:read environment:roles:viewer"
}

func (d *KeptnDynatraceDQLProvider) doOAuth(ctx context.Context, provider klcv1alpha2.KeptnEvaluationProvider) (jwt DynatraceOAuthResponse, err error) {
	d.Log.V(10).Info("OAuth login")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	secret, err := getDTSecret(ctx, provider, d.k8sClient)
	if err != nil {
		d.Log.Error(err, "Error while fetching the Dynatrace secret")
		return jwt, err
	}
	if err2 := validateOAuthSecret(secret); err2 != nil {
		d.Log.Error(err, "Error while parsing the Dynatrace secret")
		return jwt, err2
	}
	secretParts := strings.Split(secret, ".")
	clientId := fmt.Sprintf("%s.%s", secretParts[0], secretParts[1])
	clientSecret := fmt.Sprintf("%s.%s", clientId, secretParts[2])
	values := url.Values{}
	values.Add("grant_type", "client_credentials")
	values.Add("scope", d.getScopes())
	values.Add("client_id", clientId)
	values.Add("client_secret", clientSecret)
	body := []byte(values.Encode())
	authURL := "https://sso-dev.dynatracelabs.com/sso/oauth2/token"
	req, err := http.NewRequestWithContext(ctx, "POST", authURL, bytes.NewBuffer(body))
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return jwt, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := d.httpClient.Do(req)
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return jwt, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			d.Log.Error(err, "Could not close request body")
		}
	}()
	// we ignore the error here because we fail later while unmarshalling
	b, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(b, &jwt)
	if err != nil {
		d.Log.Error(err, "Error while parsing response")
		return jwt, err
	}
	defer res.Body.Close()
	return jwt, nil
}

func (d *KeptnDynatraceDQLProvider) postDQL(ctx context.Context, jwt DynatraceOAuthResponse, o klcv1alpha2.Objective, p klcv1alpha2.KeptnEvaluationProvider) (dqlHandler DynatraceDQLHandler, err error) {
	d.Log.V(10).Info("posting DQL")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	values := url.Values{}
	values.Add("query", o.Query)

	api := fmt.Sprintf("%s/platform/storage/query/v0.7/query:execute?%s", p.Spec.TargetServer, values.Encode())
	req, err := http.NewRequestWithContext(ctx, "POST", api, bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return dqlHandler, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt.accessToken))
	res, err := d.httpClient.Do(req)
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return dqlHandler, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			d.Log.Error(err, "Could not close request body")
		}
	}()
	// we ignore the error here because we fail later while unmarshalling
	b, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(b, &dqlHandler)
	if err != nil {
		d.Log.Error(err, "Error while parsing response")
		return dqlHandler, err
	}
	defer res.Body.Close()
	return dqlHandler, nil
}
