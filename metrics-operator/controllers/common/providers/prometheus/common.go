package prometheus

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"sync"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
	promapi "github.com/prometheus/client_golang/api"
	"github.com/prometheus/common/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	secretKeyUserName = "user"
	secretKeyPassword = "password"
	secretKeyCAFile   = "caFile"
	secretKeyCertFile = "certFile"
	secretKeyKeyFile  = "keyFile"
)

var (
	ErrSecretKeyRefNotDefined = errors.New("the SecretKeyRef property with the Prometheus API Key is missing")
	ErrInvalidSecretFormat    = errors.New("secret key does not contain required fields")
)

type SecretData struct {
	User     string        `json:"user"`
	Password config.Secret `json:"password"`
	CAFile   string        `json:"caFile"`
	CertFile string        `json:"certFile"`
	KeyFile  string        `json:"keyFile"`
}

type IRoundTripper interface {
	GetRoundTripper(context.Context, metricsapi.KeptnMetricsProvider, client.Client) (http.RoundTripper, error)
}

type RoundTripperRetriever struct{}

type TLSRoundTripperSettings struct {
	CAFile             string
	CertFile           string
	KeyFile            string
	ServerName         string
	InsecureSkipVerify bool
}

type tlsRoundTripper struct {
	settings  TLSRoundTripperSettings
	newRT     func(*tls.Config) (http.RoundTripper, error)
	mtx       sync.RWMutex
	rt        http.RoundTripper
	tlsConfig *tls.Config
}

func NewTLSRoundTripper(settings TLSRoundTripperSettings) (*tlsRoundTripper, error) {
	return &tlsRoundTripper{
		settings: settings,
		newRT: func(config *tls.Config) (http.RoundTripper, error) {
			return &http.Transport{
				TLSClientConfig: config,
			}, nil
		},
	}, nil
}

func (t *tlsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mtx.RLock()
	rt := t.rt
	t.mtx.RUnlock()

	if rt == nil {
		t.mtx.Lock()
		defer t.mtx.Unlock()
		if t.rt == nil {
			var err error
			t.tlsConfig, err = newTLSConfig(&t.settings)
			if err != nil {
				return nil, err
			}
			rt, err = t.newRT(t.tlsConfig)
			if err != nil {
				return nil, err
			}
			t.rt = rt
		}
		rt = t.rt
	}

	return rt.RoundTrip(req)
}

func newTLSConfig(settings *TLSRoundTripperSettings) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: settings.InsecureSkipVerify,
		ServerName:         settings.ServerName,
	}

	if settings.CAFile != "" {
		caConfig := &config.TLSConfig{
			CAFile: settings.CAFile,
		}
		ca, err := config.NewTLSConfig(caConfig)
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs = ca.RootCAs
	}

	if settings.CertFile != "" && settings.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(settings.CertFile, settings.KeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}

func (r RoundTripperRetriever) GetRoundTripper(ctx context.Context, provider metricsapi.KeptnMetricsProvider, k8sClient client.Client) (http.RoundTripper, error) {
	secret, err := getPrometheusSecret(ctx, provider, k8sClient)
	if err != nil {
		if errors.Is(err, ErrSecretKeyRefNotDefined) {
			return promapi.DefaultRoundTripper, nil
		}
		return nil, err
	}

	settings := TLSRoundTripperSettings{
		CAFile:     secret.CAFile,
		CertFile:   secret.CertFile,
		KeyFile:    secret.KeyFile,
		ServerName: provider.Spec.TargetServer,
	}

	tlsRT, err := NewTLSRoundTripper(settings)
	if err != nil {
		return nil, err
	}

	return config.NewBasicAuthRoundTripper(secret.User, secret.Password, "", "", tlsRT), nil
}

func getPrometheusSecret(ctx context.Context, provider metricsapi.KeptnMetricsProvider, k8sClient client.Client) (*SecretData, error) {
	if !provider.HasSecretDefined() {
		return nil, ErrSecretKeyRefNotDefined
	}

	secret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: provider.Spec.SecretKeyRef.Name, Namespace: provider.Namespace}, secret); err != nil {
		return nil, err
	}

	var secretData SecretData
	user, ok := secret.Data[secretKeyUserName]
	pw, pwOk := secret.Data[secretKeyPassword]
	if !ok || !pwOk {
		return nil, ErrInvalidSecretFormat
	}

	secretData.User = string(user)
	secretData.Password = config.Secret(pw)
	secretData.CAFile = string(secret.Data[secretKeyCAFile])
	secretData.CertFile = string(secret.Data[secretKeyCertFile])
	secretData.KeyFile = string(secret.Data[secretKeyKeyFile])

	return &secretData, nil
}
