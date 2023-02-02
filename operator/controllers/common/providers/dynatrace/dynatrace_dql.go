package dynatrace

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	count int64   `json:"count"`
	sum   float64 `json:"sum"`
	min   float64 `json:"min"`
	avg   float64 `json:"avg"`
	max   float64 `json:"max"`
}

// EvaluateQuery fetches the SLI values from dynatrace provider
func (d *KeptnDynatraceDQLProvider) EvaluateQuery(ctx context.Context, objective klcv1alpha2.Objective, provider klcv1alpha2.KeptnEvaluationProvider) (string, []byte, error) {
	// auth
	jwt, err := d.doOAuth(ctx, provider)
	if err != nil {
		return "", []byte{}, nil
	}
	// submit DQL
	dqlHandler, err := d.postDQL(ctx, jwt, provider.Spec.TargetServer, objective.Query)
	if err != nil {
		d.Log.Error(err, "Error while posting the DQL query: %s", objective.Query)
		return "", nil, err
	}
	// attend result
	results, err := d.getDQL(ctx, jwt, provider.Spec.TargetServer, dqlHandler)
	if err != nil {
		d.Log.Error(err, "Error while waiting for DQL query: %s", dqlHandler)
		return "", nil, err
	}
	// parse result
	if len(results.records) > 1 {
		d.Log.Error(err, "More than a single result, the first one will be used")
	}
	if len(results.records) == 0 {
		return "", nil, ErrInvalidResult
	}
	r := fmt.Sprintf("%f", results.records[0].value.avg)
	b, err := json.Marshal(r)
	if err != nil {
		d.Log.Error(err, "Error marshaling DQL results")
	}
	return r, b, nil
}

func (d *KeptnDynatraceDQLProvider) getScopes() string {
	return "storage:metrics:read environment:roles:viewer"
}

func (d *KeptnDynatraceDQLProvider) doOAuth(ctx context.Context, provider klcv1alpha2.KeptnEvaluationProvider) (jwt DynatraceOAuthResponse, err error) {
	d.Log.V(10).Info("OAuth login")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	secret, err := getDTSecret(ctx, provider, d.k8sClient)
	if err != nil {
		return jwt, err
	}
	if err2 := validateOAuthSecret(secret); err2 != nil {
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
		return jwt, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := d.httpClient.Do(req)
	if err != nil {
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
		return jwt, err
	}
	return jwt, nil
}

func (d *KeptnDynatraceDQLProvider) postDQL(ctx context.Context, jwt DynatraceOAuthResponse, server, query string) (dqlHandler DynatraceDQLHandler, err error) {
	d.Log.V(10).Info("posting DQL")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	values := url.Values{}
	values.Add("query", query)

	api := fmt.Sprintf("%s/platform/storage/query/v0.7/query:execute?%s", server, values.Encode())
	req, err := http.NewRequestWithContext(ctx, "POST", api, bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		return dqlHandler, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt.accessToken))
	res, err := d.httpClient.Do(req)
	if err != nil {
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
		return dqlHandler, err
	}
	return dqlHandler, nil
}

func (d *KeptnDynatraceDQLProvider) getDQL(ctx context.Context, jwt DynatraceOAuthResponse, server string, handler DynatraceDQLHandler) (*DQLResult, error) {
	d.Log.V(10).Info("posting DQL")
	for true {
		r, err := d.retrieveDQLResults(ctx, jwt, server, handler)
		if err != nil {
			return &DQLResult{}, err
		}
		if r.state == "SUCCEEDED" {
			return &r.result, nil
		}
		d.Log.V(10).Info("DQL not finished, got: %s", r.state)
		time.Sleep(5 * time.Second)
	}
	return nil, errors.New("something went wrong while waiting for the DQL to be finished")
}

func (d *KeptnDynatraceDQLProvider) retrieveDQLResults(ctx context.Context, jwt DynatraceOAuthResponse, server string, handler DynatraceDQLHandler) (result DynatraceDQLResult, err error) {
	d.Log.V(10).Info("Getting DQL")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	values := url.Values{}
	values.Add("request-token", handler.requestToken)

	api := fmt.Sprintf("%s/platform/storage/query/v0.7/query:poll?%s", server, values.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", api, nil)
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return result, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt.accessToken))
	res, err := d.httpClient.Do(req)
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return result, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			d.Log.Error(err, "Could not close request body")
		}
	}()
	// we ignore the error here because we fail later while unmarshalling
	b, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(b, &result)
	if err != nil {
		d.Log.Error(err, "Error while parsing response")
		return result, err
	}
	defer res.Body.Close()
	return result, nil
}
