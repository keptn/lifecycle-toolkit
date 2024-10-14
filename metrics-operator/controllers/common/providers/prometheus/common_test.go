package prometheus

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"

	"github.com/prometheus/common/config"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const prometheusPayload = "test"

func TestGetSecret_NoKeyDefined(t *testing.T) {

	fakeClient := fake.NewClient()

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: "svr.URL",
		},
	}
	r1, e := getPrometheusSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.ErrorIs(t, ErrSecretKeyRefNotDefined, e)
	require.Empty(t, r1)

}

func TestGetSecret_NoSecretDefined(t *testing.T) {

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
			TargetServer: "svr",
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

func Test_GetRoundtripper(t *testing.T) {
	goodsecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"user":     []byte("myuser"),
			"password": []byte("mytoken"),
		},
	}
	tests := []struct {
		name      string
		provider  metricsapi.KeptnMetricsProvider
		k8sClient client.Client
		wantUser  string
		wantPass  string
		wantErr   bool
		errorStr  string
	}{
		{
			name: "TestSuccess",
			provider: metricsapi.KeptnMetricsProvider{
				ObjectMeta: metav1.ObjectMeta{Namespace: "default"},
				Spec: metricsapi.KeptnMetricsProviderSpec{
					Type:         "",
					TargetServer: "",
					SecretKeyRef: v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "test",
						},
						Key:      "",
						Optional: nil,
					},
				},
			},
			k8sClient: fake.NewClient(goodsecret),
			wantUser:  "myuser",
			wantPass:  "mytoken",
			wantErr:   false,
		},
		{
			name:      "TestSecretNotDefined",
			provider:  metricsapi.KeptnMetricsProvider{},
			k8sClient: fake.NewClient(),
			wantUser:  "",
			wantPass:  "",
			wantErr:   false,
		},
		{
			name: "TestErrorFromGetPrometheusSecretNotExists",
			provider: metricsapi.KeptnMetricsProvider{
				ObjectMeta: metav1.ObjectMeta{Namespace: "default"},
				Spec: metricsapi.KeptnMetricsProviderSpec{
					Type:         "",
					TargetServer: "",
					SecretKeyRef: v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "test",
						},
						Key:      "",
						Optional: nil,
					},
				},
			},
			k8sClient: fake.NewClient(),
			wantUser:  "",
			wantPass:  "",
			wantErr:   true,
			errorStr:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RoundTripperRetriever{}.GetRoundTripper(context.TODO(), tt.provider, tt.k8sClient)
			t.Log(err)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRoundtripper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errorStr != "" && !strings.Contains(err.Error(), tt.errorStr) {
				t.Errorf("getRoundtripper() error = %s, wantErr %s", err.Error(), tt.errorStr)
				return
			}

			if tt.wantErr {
				return
			}

			if _, ok := got.(*http.Transport); ok {
				if tt.wantUser != "" || tt.wantPass != "" {
					t.Errorf("getRoundtripper() got default RoundTripper, want BasicAuth")
				}
				return
			}

			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					t.Errorf("Authorization header not set")
					return
				}

				auth := strings.SplitN(authHeader, " ", 2)
				if len(auth) != 2 || auth[0] != "Basic" {
					t.Errorf("Invalid Authorization header format")
					return
				}
				payload, _ := base64.StdEncoding.DecodeString(auth[1])
				pair := strings.SplitN(string(payload), ":", 2)

				if pair[0] != tt.wantUser {
					t.Errorf("got user = %v, want %v", pair[0], tt.wantUser)
				}
				if pair[1] != tt.wantPass {
					t.Errorf("got password = %v, want %v", pair[1], tt.wantPass)
				}
			}))
			defer testServer.Close()

			req, _ := http.NewRequest("GET", testServer.URL, nil)
			_, err = got.RoundTrip(req)
			if err != nil {
				t.Errorf("RoundTrip error: %v", err)
			}
		})
	}
}
