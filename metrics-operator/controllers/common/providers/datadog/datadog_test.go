package datadog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const ddErrorPayload = "{\"error\":\"Token is missing required scope\"}"
const ddPayload = "{\"from_date\":1677736306000,\"group_by\":[],\"message\":\"\",\"query\":\"system.cpu.idle{*}\",\"res_type\":\"time_series\",\"series\":[{\"aggr\":null,\"display_name\":\"system.cpu.idle\",\"end\":1677821999000,\"expression\":\"system.cpu.idle{*}\",\"interval\":300,\"length\":7,\"metric\":\"system.cpu.idle\",\"pointlist\":[[1677781200000,92.37997436523438],[1677781500000,91.46615447998047],[1677781800000,92.05865631103515],[1677782100000,97.49858474731445],[1677782400000,95.95263163248698],[1677821400000,69.67094268798829],[1677821700000,84.78184509277344]],\"query_index\":0,\"scope\":\"*\",\"start\":1677781200000,\"tag_set\":[],\"unit\":[{\"family\":\"percentage\",\"name\":\"percent\",\"plural\":\"percent\",\"scale_factor\":1,\"short_name\":\"%\"},{}]}],\"status\":\"ok\",\"to_date\":1677822706000}"
const ddEmptyPayload = "{\"from_date\":1677736306000,\"group_by\":[],\"message\":\"\",\"query\":\"system.cpu.idle{*}\",\"res_type\":\"time_series\",\"series\":[],\"status\":\"ok\",\"to_date\":1677822706000}"

func TestEvaluateQuery_APIError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddErrorPayload))
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
	kdd := setupTest(apiToken)
	metrics := []metricsapi.KeptnMetric{
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
			},
		},
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
				Range: &metricsapi.RangeSpec{Interval: "5m"},
			},
		},
	}
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
	for _, metric := range metrics {
		r, raw, e := kdd.EvaluateQuery(context.TODO(), metric, p)
		require.Error(t, e)
		require.Contains(t, e.Error(), "Token is missing required scope")
		require.Equal(t, []byte(ddErrorPayload), raw)
		require.Empty(t, r)
	}
}

func TestEvaluateQuery_HappyPath(t *testing.T) {
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
	kdd := setupTest(apiToken)
	metrics := []metricsapi.KeptnMetric{
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
			},
		},
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
				Range: &metricsapi.RangeSpec{Interval: "5m"},
			},
		},
	}
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
	for _, metric := range metrics {
		r, raw, e := kdd.EvaluateQuery(context.TODO(), metric, p)
		require.Nil(t, e)
		require.Equal(t, []byte(ddPayload), raw)
		require.Equal(t, fmt.Sprintf("%.3f", 89.116), r)
	}
}
func TestEvaluateQuery_WrongPayloadHandling(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("garbage"))
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
	kdd := setupTest(apiToken)
	metrics := []metricsapi.KeptnMetric{
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
			},
		},
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
				Range: &metricsapi.RangeSpec{Interval: "5m"},
			},
		},
	}
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
	for _, metric := range metrics {
		r, raw, e := kdd.EvaluateQuery(context.TODO(), metric, p)
		require.Equal(t, "", r)
		require.Equal(t, []byte("garbage"), raw)
		require.NotNil(t, e)
	}
}
func TestEvaluateQuery_MissingSecret(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	kdd := setupTest()
	metrics := []metricsapi.KeptnMetric{
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
			},
		},
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
				Range: &metricsapi.RangeSpec{Interval: "5m"},
			},
		},
	}
	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: svr.URL,
		},
	}
	for _, metric := range metrics {
		_, _, e := kdd.EvaluateQuery(context.TODO(), metric, p)
		require.NotNil(t, e)
		require.ErrorIs(t, e, ErrSecretKeyRefNotDefined)
	}
}
func TestEvaluateQuery_SecretNotFound(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	secretName := "datadogSecret"

	kdd := setupTest()
	metrics := []metricsapi.KeptnMetric{
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
			},
		},
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
				Range: &metricsapi.RangeSpec{Interval: "5m"},
			},
		},
	}
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
	for _, metric := range metrics {
		_, _, e := kdd.EvaluateQuery(context.TODO(), metric, p)
		require.NotNil(t, e)
		require.True(t, errors.IsNotFound(e))
	}
}
func TestEvaluateQuery_RefNonExistingKey(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	secretName := "datadogSecret"
	apiKey, apiKeyValue := "I_AM_NOT_DD_CLIENT_API_KEY", "value"
	apiToken := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "",
		},
		Data: map[string][]byte{
			apiKey: []byte(apiKeyValue),
		},
	}
	kdd := setupTest(apiToken)
	metrics := []metricsapi.KeptnMetric{
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
			},
		},
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
				Range: &metricsapi.RangeSpec{Interval: "5m"},
			},
		},
	}
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
	for _, metric := range metrics {
		_, _, e := kdd.EvaluateQuery(context.TODO(), metric, p)
		require.NotNil(t, e)
		require.True(t, strings.Contains(e.Error(), "secret does not contain DD_CLIENT_API_KEY or DD_CLIENT_APP_KEY"))
	}
}
func TestEvaluateQuery_EmptyPayload(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddEmptyPayload))
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
	kdd := setupTest(apiToken)
	metrics := []metricsapi.KeptnMetric{
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
			},
		},
		{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "system.cpu.idle{*}",
				Range: &metricsapi.RangeSpec{Interval: "5m"},
			},
		},
	}
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
	for _, metric := range metrics {
		r, raw, e := kdd.EvaluateQuery(context.TODO(), metric, p)
		t.Log(string(raw))
		require.Nil(t, raw)
		require.Equal(t, "", r)
		require.True(t, strings.Contains(e.Error(), "no values in query result"))
	}
}

func TestEvaluateQuery_WrongInterval(t *testing.T) {
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
	kdd := setupTest(apiToken)
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "system.cpu.idle{*}",
			Range: &metricsapi.RangeSpec{Interval: "5mins"},
		},
	}
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
	r, raw, e := kdd.EvaluateQuery(context.TODO(), metric, p)
	require.Error(t, e)
	require.Contains(t, e.Error(), "unknown unit \"mins\" in duration \"5mins\"")
	require.Equal(t, []byte(nil), raw)
	require.Empty(t, r)
}

func TestEvaluateQueryForStep_APIError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddErrorPayload))
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
	kdd := setupTest(apiToken)
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "system.cpu.idle{*}",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}
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
	r, raw, e := kdd.EvaluateQueryForStep(context.TODO(), metric, p)
	require.Error(t, e)
	require.Contains(t, e.Error(), "Token is missing required scope")
	require.Equal(t, []byte(ddErrorPayload), raw)
	require.Empty(t, r)
}

func TestEvaluateQueryForStep_HappyPath(t *testing.T) {
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
	kdd := setupTest(apiToken)
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "system.cpu.idle{*}",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}
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
	r, raw, e := kdd.EvaluateQueryForStep(context.TODO(), metric, p)
	require.Nil(t, e)
	require.Equal(t, []byte(ddPayload), raw)
	require.Equal(t, []string{"92.379974", "91.466154", "92.058656", "97.498585", "95.952632", "69.670943", "84.781845"}, r)
}
func TestEvaluateQueryForStep_WrongPayloadHandling(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("garbage"))
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
	kdd := setupTest(apiToken)
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "system.cpu.idle{*}",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}
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
	r, raw, e := kdd.EvaluateQueryForStep(context.TODO(), metric, p)
	require.Equal(t, []string(nil), r)
	require.Equal(t, []byte("garbage"), raw)
	require.NotNil(t, e)
}
func TestEvaluateQueryForStep_MissingSecret(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	kdd := setupTest()
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "system.cpu.idle{*}",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}
	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: svr.URL,
		},
	}
	_, _, e := kdd.EvaluateQueryForStep(context.TODO(), metric, p)
	require.NotNil(t, e)
	require.ErrorIs(t, e, ErrSecretKeyRefNotDefined)
}
func TestEvaluateQueryForStep_SecretNotFound(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	secretName := "datadogSecret"

	kdd := setupTest()
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "system.cpu.idle{*}",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}
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
	_, _, e := kdd.EvaluateQueryForStep(context.TODO(), metric, p)
	require.NotNil(t, e)
	require.True(t, errors.IsNotFound(e))
}
func TestEvaluateQueryForStep_RefNonExistingKey(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	secretName := "datadogSecret"
	apiKey, apiKeyValue := "I_AM_NOT_DD_CLIENT_API_KEY", "value"
	apiToken := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "",
		},
		Data: map[string][]byte{
			apiKey: []byte(apiKeyValue),
		},
	}
	kdd := setupTest(apiToken)
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "system.cpu.idle{*}",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}
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
	_, _, e := kdd.EvaluateQueryForStep(context.TODO(), metric, p)
	require.NotNil(t, e)
	require.True(t, strings.Contains(e.Error(), "secret does not contain DD_CLIENT_API_KEY or DD_CLIENT_APP_KEY"))

}
func TestEvaluateQueryForStep_EmptyPayload(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(ddEmptyPayload))
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
	kdd := setupTest(apiToken)
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "system.cpu.idle{*}",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}
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
	r, raw, e := kdd.EvaluateQueryForStep(context.TODO(), metric, p)
	t.Log(string(raw))
	require.Nil(t, raw)
	require.Equal(t, []string(nil), r)
	require.True(t, strings.Contains(e.Error(), "no values in query result"))
}

func TestEvaluateQueryForStep_WrongInterval(t *testing.T) {
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
	kdd := setupTest(apiToken)
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "system.cpu.idle{*}",
			Range: &metricsapi.RangeSpec{
				Interval:    "5mins",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}
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
	r, raw, e := kdd.EvaluateQueryForStep(context.TODO(), metric, p)
	require.Error(t, e)
	require.Contains(t, e.Error(), "unknown unit \"mins\" in duration \"5mins\"")
	require.Equal(t, []byte(nil), raw)
	require.Empty(t, r)
}

func TestEvaluateQueryForStep_WrongStep(t *testing.T) {
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
	kdd := setupTest(apiToken)
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "system.cpu.idle{*}",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1mins",
				Aggregation: "max",
			},
		},
	}
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
	r, raw, e := kdd.EvaluateQueryForStep(context.TODO(), metric, p)
	require.Error(t, e)
	require.Contains(t, e.Error(), "unknown unit \"mins\" in duration \"1mins\"")
	require.Equal(t, []byte(nil), raw)
	require.Empty(t, r)
}

func TestGetSingleValue_EmptyPoints(t *testing.T) {
	kdd := setupTest()
	var points [][]*float64
	value := kdd.getSingleValue(points)

	require.Zero(t, value)
}
func TestGetSingleValue_HappyPath(t *testing.T) {

	kdd := setupTest()
	result := datadogV1.MetricsQueryResponse{}
	_ = json.Unmarshal([]byte(ddPayload), &result)
	points := (result.Series)[0].Pointlist
	value := kdd.getSingleValue(points)

	require.NotZero(t, value)
	require.Equal(t, 89.11554133097331, value)
}

func TestGetResultSlice_EmptyPoints(t *testing.T) {
	kdd := setupTest()
	var points [][]*float64
	value := kdd.getResultSlice(points)

	require.Equal(t, []string{}, value)
}

func TestGetResultSlice_HappyPath(t *testing.T) {

	kdd := setupTest()
	result := datadogV1.MetricsQueryResponse{}
	_ = json.Unmarshal([]byte(ddPayload), &result)
	points := (result.Series)[0].Pointlist
	resultSlice := kdd.getResultSlice(points)

	require.NotZero(t, resultSlice)
	require.Equal(t, []string{"92.379974", "91.466154", "92.058656", "97.498585", "95.952632", "69.670943", "84.781845"}, resultSlice)
}

func setupTest(objs ...client.Object) KeptnDataDogProvider {

	fakeClient := fake.NewClient(objs...)

	kdd := KeptnDataDogProvider{
		HttpClient: http.Client{},
		Log:        ctrl.Log.WithName("testytest"),
		K8sClient:  fakeClient,
	}
	return kdd
}
