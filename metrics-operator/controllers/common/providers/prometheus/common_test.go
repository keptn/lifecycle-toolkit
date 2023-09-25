package prometheus

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/prometheus/common/config"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const prometheusPayload = "test"

func TestGetSecret_NoKeyDefined(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(prometheusPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()
	fakeClient := fake.NewClient()

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: svr.URL,
		},
	}
	r1, e := getPrometheusSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.ErrorIs(t, ErrSecretKeyRefNotDefined, e)
	require.Empty(t, r1)

}

func TestGetSecret_NoSecretDefined(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(prometheusPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	secretName := "testSecret"

	fakeClient := fake.NewClient()

	b := true
	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				Key: "apiKey",
				LocalObjectReference: v1.LocalObjectReference{
					Name: secretName,
				},
				Optional: &b,
			},
			TargetServer: svr.URL,
		},
	}
	r1, e := getPrometheusSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	t.Log(e.Error())
	require.True(t, k8serrors.IsNotFound(e))
	require.Empty(t, r1)

}

func TestGetSecret_HappyPath(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(prometheusPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	secretName := "mySecret"
	apiToken := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "default",
		},
		Data: map[string][]byte{
			"user":     []byte("myuser"),
			"password": []byte("mytoken"),
		},
	}
	fakeClient := fake.NewClient(apiToken)

	p := metricsapi.KeptnMetricsProvider{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default"},
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				Key: "login",
				LocalObjectReference: v1.LocalObjectReference{
					Name: secretName,
				},
			},
			TargetServer: svr.URL,
		},
	}
	r1, e := getPrometheusSecret(context.TODO(), p, fakeClient)
	require.Nil(t, e)
	require.Equal(t, "myuser", r1.User)
	require.Equal(t, config.Secret("mytoken"), r1.Password)

}
