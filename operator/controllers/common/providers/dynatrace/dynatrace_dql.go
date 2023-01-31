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

// EvaluateQuery fetches the SLI values from dynatrace provider
func (d *KeptnDynatraceDQLProvider) EvaluateQuery(ctx context.Context, objective klcv1alpha2.Objective, provider klcv1alpha2.KeptnEvaluationProvider) (string, []byte, error) {
	// auth
	jwt, err := d.doOAuth(ctx, provider)
	if err != nil {
		return "", []byte{}, nil
	}

	// submit DQL
	// attend result
	// parse result
	return "", []byte{}, nil
}

func (d *KeptnDynatraceDQLProvider) getScopes() string {
	return "storage:events:read storage:events:write environment:roles:viewer"
}

func (d *KeptnDynatraceDQLProvider) doOAuth(ctx context.Context, provider klcv1alpha2.KeptnEvaluationProvider) (*DynatraceOAuthResponse, error) {
	d.Log.V(10).Info("OAuth login")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	secret, err := getDTSecret(ctx, provider, d.k8sClient)
	if err != nil {
		d.Log.Error(err, "Error while fetching the Dynatrace secret")
		return nil, err
	}
	if err2 := validateOAuthSecret(secret); err2 != nil {
		d.Log.Error(err, "Error while parsing the Dynatrace secret")
		return nil, err2
	}
	secretParts := strings.Split(secret, ".")
	values := url.Values{}
	values.Add("grant_type", "client_credentials")
	values.Add("scope", d.getScopes())
	values.Add("client_secret", fmt.Sprintf("%s.%s.%s", secretParts[0], secretParts[1], secretParts[2]))
	values.Add("client_id", fmt.Sprintf("%s.%s", secretParts[0], secretParts[1]))
	body := []byte(values.Encode())
	authURL := "https://sso-dev.dynatracelabs.com/sso/oauth2/token"
	req, err := http.NewRequestWithContext(ctx, "POST", authURL, bytes.NewBuffer(body))
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := d.httpClient.Do(req)
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return nil, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			d.Log.Error(err, "Could not close request body")
		}
	}()
	// we ignore the error here because we fail later while unmarshalling
	b, _ := io.ReadAll(res.Body)
	jwt := DynatraceOAuthResponse{}
	err = json.Unmarshal(b, &jwt)
	if err != nil {
		d.Log.Error(err, "Error while parsing response")
		return nil, err
	}
	defer res.Body.Close()
	return &jwt, nil
}
