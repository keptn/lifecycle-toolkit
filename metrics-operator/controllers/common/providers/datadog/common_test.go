package datadog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetSecret_NoKeyDefined(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()
	fakeClient := fake.NewClient()

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: svr.URL,
		},
	}
	r1, r2, e := getDDSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.ErrorIs(t, e, ErrSecretKeyRefNotDefined)
	require.Empty(t, r1)
	require.Empty(t, r2)

}

func TestGetSecret_NoSecretDefined(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	secretName := "datadogSecret"
	apiKey, apiKeyValue := "DD_CLIENT_API_KEY", "fake-api-key"
	appKey, appKeyValue := "DD_CLIENT_APP_KEY", "fake-app-key"
	apiToken := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "garbage",
			Namespace: "",
		},
		Data: map[string][]byte{
			apiKey: []byte(apiKeyValue),
			appKey: []byte(appKeyValue),
		},
	}
	fakeClient := fake.NewClient(apiToken)

	b := true
	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: secretName,
				},
				Optional: &b,
			},
			TargetServer: svr.URL,
		},
	}
	r1, r2, e := getDDSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.True(t, strings.Contains(e.Error(), "secrets \""+secretName+"\" not found"))
	require.Empty(t, r1)
	require.Empty(t, r2)

}

func TestGetSecret_HappyPath(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	secretName := "datadogSecret"
	apiKey, apiKeyValue := "DD_CLIENT_API_KEY", "fake-api-key"
	appKey, appKeyValue := "DD_CLIENT_APP_KEY", "fake-app-key"
	apiToken := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "",
		},
		Data: map[string][]byte{
			apiKey: []byte(apiKeyValue),
			appKey: []byte(appKeyValue),
		},
	}
	fakeClient := fake.NewClient(apiToken)

	b := true
	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: secretName,
				},
				Optional: &b,
			},
			TargetServer: svr.URL,
		},
	}
	r1, r2, e := getDDSecret(context.TODO(), p, fakeClient)
	require.Nil(t, e)
	require.Equal(t, apiKeyValue, r1)
	require.Equal(t, appKeyValue, r2)

}
