package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

//go:generate moq -pkg fake --skip-ensure -out ./fake/dt_client_mock.go . DTAPIClient
type DTAPIClient interface {
	Do(ctx context.Context, path, method string, payload []byte) ([]byte, int, error)
}

type apiClient struct {
	Log        logr.Logger
	httpClient http.Client
	config     apiConfig
}

type APIClientOption func(client *apiClient)

// WithLogger injects the given logger into an APIClient
func WithLogger(logger logr.Logger) APIClientOption {
	return func(client *apiClient) {
		client.Log = logger
	}
}

// WithHTTPClient injects the given HTTP client into an APIClient
func WithHTTPClient(httpClient http.Client) APIClientOption {
	return func(client *apiClient) {
		client.httpClient = httpClient
	}
}

// NewAPIClient creates and returns a new APIClient
func NewAPIClient(config apiConfig, options ...APIClientOption) *apiClient {
	client := &apiClient{
		Log:        logr.New(klog.NewKlogr().GetSink()),
		httpClient: http.Client{},
		config:     config,
	}

	for _, o := range options {
		o(client)
	}

	return client
}

// Do sends and API request to the Dynatrace API and returns its result as a string containing the raw response payload
func (client *apiClient) Do(ctx context.Context, path, method string, payload []byte) ([]byte, int, error) {
	if err := client.auth(ctx); err != nil {
		return nil, http.StatusUnauthorized, err
	}
	api := fmt.Sprintf("%s%s", client.config.serverURL, path)
	req, err := http.NewRequestWithContext(ctx, method, api, bytes.NewBuffer(payload))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.config.oAuthCredentials.accessToken))

	client.Log.Info(fmt.Sprintf(" sending request payload: %s", string(payload)))
	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			client.Log.Error(err, "Could not close request body")
		}
	}()
	if isErrorStatus(res.StatusCode) {
		return nil, res.StatusCode, ErrRequestFailed
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return b, res.StatusCode, nil
}

func (client *apiClient) auth(ctx context.Context) error {
	// return if we already have a token
	if client.config.oAuthCredentials.accessToken != "" {
		return nil
	}
	client.Log.V(10).Info("OAuth login")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	values := client.config.oAuthCredentials.urlValues()
	body := []byte(values.Encode())

	req, err := http.NewRequestWithContext(ctx, "POST", client.config.authURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			client.Log.Error(err, "Could not close request body")
		}
	}()
	if isErrorStatus(res.StatusCode) {
		return ErrRequestFailed
	}
	// we ignore the error here because we fail later while unmarshalling
	oauthResponse := OAuthResponse{}
	b, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(b, &oauthResponse)
	if err != nil {
		return err
	}

	if oauthResponse.AccessToken == "" {
		return ErrAuthenticationFailed
	}

	client.config.oAuthCredentials.accessToken = oauthResponse.AccessToken
	return nil
}
