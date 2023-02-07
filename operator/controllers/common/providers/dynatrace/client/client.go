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

const defaultAuthURL = "https://sso-dev.dynatracelabs.com/sso/oauth2/token"
const defaultScopes = "storage:metrics:read environment:roles:viewer"

type OAuthResponse struct {
	AccessToken string `json:"accessToken"`
}

type OAuthCredentials struct {
	clientID     string
	clientSecret string
	scopes       string
	accessToken  string
}

func (oac OAuthCredentials) URLValues() url.Values {
	values := url.Values{}
	values.Add("grant_type", "client_credentials")
	values.Add("scope", oac.scopes)
	values.Add("client_id", oac.clientID)
	values.Add("client_secret", oac.clientSecret)

	return values
}

type apiConfig struct {
	serverURL        string
	authURL          string
	oAuthCredentials OAuthCredentials
}

type APIConfigOption func(config *apiConfig)

func WithAuthURL(authURL string) APIConfigOption {
	return func(config *apiConfig) {
		config.authURL = authURL
	}
}

func WithScopes(scopes string) APIConfigOption {
	return func(config *apiConfig) {
		config.oAuthCredentials.scopes = scopes
	}
}

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
		oAuthCredentials: OAuthCredentials{
			clientID:     clientId,
			clientSecret: clientSecret,
			scopes:       defaultScopes,
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

func WithLogger(logger logr.Logger) APIClientOption {
	return func(client *apiClient) {
		client.Log = logger
	}
}

func WithHTTPClient(httpClient http.Client) APIClientOption {
	return func(client *apiClient) {
		client.httpClient = httpClient
	}
}

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

func (client *apiClient) Do(ctx context.Context, path, method string, payload []byte) ([]byte, error) {
	if err := client.Auth(ctx); err != nil {
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
	// we ignore the error here because we fail later while unmarshalling
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (client *apiClient) Auth(ctx context.Context) error {
	// return if we already have a token
	if client.config.oAuthCredentials.accessToken != "" {
		return nil
	}
	client.Log.V(10).Info("OAuth login")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	values := client.config.oAuthCredentials.URLValues()
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
	// we ignore the error here because we fail later while unmarshalling
	oauthResponse := OAuthResponse{}
	b, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(b, &oauthResponse)
	if err != nil {
		return err
	}

	client.config.oAuthCredentials.accessToken = oauthResponse.AccessToken
	return nil
}
