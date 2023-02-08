package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

const (
	defaultAuthURL                  = "https://sso-dev.dynatracelabs.com/sso/oauth2/token"
	oAuthGrantType                  = "grant_type"
	oAuthGrantTypeClientCredentials = "client_credentials"
	oAuthScope                      = "scope"
	oAuthClientID                   = "client_id"
	oAuthClientSecret               = "client_secret"
)

// OAuthScope represents a scope provided for the registered OAuth client interacting with the DT API
type OAuthScope string

// These constants define the scopes that we currently need for the DQL metric functionality. This list might extend as new features will be added.
// For now, we keep this at the minimum set of scopes required, as these are currently likely to change
const (
	OAuthScopeStorageMetricsRead    = "storage:metrics:read"
	OAuthScopeEnvironmentRoleViewer = "environment:roles:viewer"
)

type OAuthResponse struct {
	AccessToken string `json:"access_token"`
}

type apiConfig struct {
	serverURL        string
	authURL          string
	oAuthCredentials oAuthCredentials
}

type APIConfigOption func(config *apiConfig)

func WithAuthURL(authURL string) APIConfigOption {
	return func(config *apiConfig) {
		config.authURL = authURL
	}
}

// WithScopes passes the given scopes to the client config
func WithScopes(scopes []OAuthScope) APIConfigOption {
	return func(config *apiConfig) {
		config.oAuthCredentials.scopes = scopes
	}
}

// NewAPIConfig returns a new apiConfig that can be used for initializing a DTAPIClient with the NewAPIClient function
func NewAPIConfig(serverURL string, secret string, opts ...APIConfigOption) (*apiConfig, error) {
	if err := validateOAuthSecret(secret); err != nil {
		return nil, err
	}

	secretParts := strings.Split(secret, ".")
	clientId := fmt.Sprintf("%s.%s", secretParts[0], secretParts[1])
	clientSecret := fmt.Sprintf("%s.%s", clientId, secretParts[2])

	cfg := &apiConfig{
		serverURL: serverURL,
		authURL:   defaultAuthURL,
		oAuthCredentials: oAuthCredentials{
			clientID:     clientId,
			clientSecret: clientSecret,
			scopes:       []OAuthScope{OAuthScopeStorageMetricsRead, OAuthScopeEnvironmentRoleViewer},
		},
	}

	for _, o := range opts {
		o(cfg)
	}

	return cfg, nil
}

//go:generate moq -pkg fake --skip-ensure -out ./fake/dt_client_mock.go . DTAPIClient
type DTAPIClient interface {
	Do(ctx context.Context, path, method string, payload []byte) ([]byte, error)
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
func (client *apiClient) Do(ctx context.Context, path, method string, payload []byte) ([]byte, error) {
	if err := client.auth(ctx); err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s%s", client.config.serverURL, path)
	req, err := http.NewRequestWithContext(ctx, method, api, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.config.oAuthCredentials.accessToken))

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			client.Log.Error(err, "Could not close request body")
		}
	}()
	if isErrorStatus(res.StatusCode) {
		return nil, ErrRequestFailed
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
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

type oAuthCredentials struct {
	clientID     string
	clientSecret string
	scopes       []OAuthScope
	accessToken  string
}

func (oac oAuthCredentials) urlValues() url.Values {
	values := url.Values{}
	values.Add(oAuthGrantType, oAuthGrantTypeClientCredentials)
	values.Add(oAuthScope, oac.getScopesAsString())
	values.Add(oAuthClientID, oac.clientID)
	values.Add(oAuthClientSecret, oac.clientSecret)

	return values
}

func (oac oAuthCredentials) getScopesAsString() string {
	scopeStr := ""

	for i := 0; i < len(oac.scopes); i++ {
		if i == 0 {
			scopeStr += string(oac.scopes[i])
		} else {
			scopeStr += " " + string(oac.scopes[i])
		}
	}
	return scopeStr
}

func isErrorStatus(statusCode int) bool {
	return statusCode < 200 || statusCode >= 300
}
