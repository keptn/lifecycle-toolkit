package client

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	"k8s.io/klog/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewConfigInvalidSecretFormat(t *testing.T) {

	config, err := NewAPIConfig("", "my-secret")

	require.ErrorIs(t, err, ErrClientSecretInvalid)
	require.Nil(t, config)
}

func TestAPIClient(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/auth" {
			_, _ = writer.Write([]byte(`{"accessToken": "my-token"}`))
			return
		}
		_, _ = writer.Write([]byte("success"))
	}))

	defer server.Close()

	mockSecret := "dts08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	config, err := NewAPIConfig(
		server.URL,
		mockSecret,
		WithScopes("my-scopes"),
		WithAuthURL(server.URL+"/auth"),
	)

	require.Nil(t, err)
	require.NotNil(t, config)

	expectedApiConfig := apiConfig{
		serverURL: server.URL,
		authURL:   server.URL + "/auth",
		oAuthCredentials: OAuthCredentials{
			clientID:     "dts08.XX",
			clientSecret: mockSecret,
			scopes:       "my-scopes",
		},
	}
	require.Equal(t, expectedApiConfig, *config)

	apiClient := NewAPIClient(
		*config,
		WithHTTPClient(http.Client{}),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	require.NotNil(t, apiClient)

	resp, err := apiClient.Do(context.TODO(), "/query", http.MethodPost, nil)

	require.Nil(t, err)
	require.Equal(t, "success", string(resp))

	require.Equal(t, "my-token", apiClient.config.oAuthCredentials.accessToken)
}
