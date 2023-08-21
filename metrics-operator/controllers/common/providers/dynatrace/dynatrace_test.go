package dynatrace

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

const dtpayload = "{\"totalCount\":1,\"nextPageKey\":null,\"resolution\":\"1m\",\"result\":[{\"metricId\":\"dsfm:billing.hostunit.assigned:splitBy():sort(value(auto,descending)):avg\",\"dataPointCountRatio\":6.0E-6,\"dimensionCountRatio\":1.0E-5,\"data\":[{\"dimensions\":[],\"dimensionMap\":{},\"timestamps\":[1666090140000,1666090200000,1666090260000,1666090320000,1666090380000,1666090440000,1666090500000,1666090560000,1666090620000,1666090680000,1666090740000,1666090800000,1666090860000,1666090920000,1666090980000,1666091040000,1666091100000,1666091160000,1666091220000,1666091280000,1666091340000,1666091400000,1666091460000,1666091520000,1666091580000,1666091640000,1666091700000,1666091760000,1666091820000,1666091880000,1666091940000,1666092000000,1666092060000,1666092120000,1666092180000,1666092240000,1666092300000,1666092360000,1666092420000,1666092480000,1666092540000,1666092600000,1666092660000,1666092720000,1666092780000,1666092840000,1666092900000,1666092960000,1666093020000,1666093080000,1666093140000,1666093200000,1666093260000,1666093320000,1666093380000,1666093440000,1666093500000,1666093560000,1666093620000,1666093680000,1666093740000,1666093800000,1666093860000,1666093920000,1666093980000,1666094040000,1666094100000,1666094160000,1666094220000,1666094280000,1666094340000,1666094400000,1666094460000,1666094520000,1666094580000,1666094640000,1666094700000,1666094760000,1666094820000,1666094880000,1666094940000,1666095000000,1666095060000,1666095120000,1666095180000,1666095240000,1666095300000,1666095360000,1666095420000,1666095480000,1666095540000,1666095600000,1666095660000,1666095720000,1666095780000,1666095840000,1666095900000,1666095960000,1666096020000,1666096080000,1666096140000,1666096200000,1666096260000,1666096320000,1666096380000,1666096440000,1666096500000,1666096560000,1666096620000,1666096680000,1666096740000,1666096800000,1666096860000,1666096920000,1666096980000,1666097040000,1666097100000,1666097160000,1666097220000,1666097280000,1666097340000],\"values\":[null,null,null,null,null,null,50,null,null,null,null,null,null,null,null,null,null,null,null,null,null,50,null,null,null,null,null,null,null,null,null,null,null,null,null,null,50,null,null,null,null,null,null,null,null,null,null,null,null,null,null,50,null,null,null,null,null,null,null,null,null,null,null,null,null,null,50,null,null,null,null,null,null,null,null,null,null,null,null,null,null,50,null,null,null,null,null,null,null,null,null,null,null,null,null,null,50,null,null,null,null,null,null,null,null,null,null,null,null,null,null,50,null,null,null,null,null,null,null,null,null]}]}]}"

func TestGetSingleValue(t *testing.T) {
	v := 5.0
	tests := []struct {
		name   string
		input  *DynatraceResponse
		result float64
	}{
		{
			name: "happy path",
			input: &DynatraceResponse{
				Result: []DynatraceResult{
					{
						Data: []DynatraceData{
							{
								Values: []*float64{&v, &v, &v},
							},
						},
					},
				},
			},
			result: v,
		},
		{
			name: "empty path",
			input: &DynatraceResponse{
				Result: []DynatraceResult{},
			},
			result: 0.0,
		},
		{
			name: "no data",
			input: &DynatraceResponse{
				Result: []DynatraceResult{
					{
						Data: []DynatraceData{},
					},
				},
			},
			result: 0.0,
		},
		{
			name: "no values",
			input: &DynatraceResponse{
				Result: []DynatraceResult{
					{
						Data: []DynatraceData{
							{
								Values: []*float64{},
							},
						},
					},
				},
			},
			result: 0.0,
		},
		{
			name: "nil values",
			input: &DynatraceResponse{
				Result: []DynatraceResult{
					{
						Data: []DynatraceData{
							{
								Values: []*float64{nil, nil, nil},
							},
						},
					},
				},
			},
			result: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kdp := KeptnDynatraceProvider{}
			r := kdp.getSingleValue(tt.input)
			require.Equal(t, tt.result, r)
		})

	}
}

func TestGetResultSlice(t *testing.T) {
	v := 5.0
	tests := []struct {
		name   string
		input  *DynatraceResponse
		result []string
	}{
		{
			name: "happy path",
			input: &DynatraceResponse{
				Result: []DynatraceResult{
					{
						Data: []DynatraceData{
							{
								Values: []*float64{&v, &v, &v},
							},
						},
					},
				},
			},
			result: []string{"5.000000", "5.000000", "5.000000"},
		},
		{
			name: "empty path",
			input: &DynatraceResponse{
				Result: []DynatraceResult{},
			},
			result: []string{},
		},
		{
			name: "no data",
			input: &DynatraceResponse{
				Result: []DynatraceResult{
					{
						Data: []DynatraceData{},
					},
				},
			},
			result: []string{},
		},
		{
			name: "no values",
			input: &DynatraceResponse{
				Result: []DynatraceResult{
					{
						Data: []DynatraceData{
							{
								Values: []*float64{},
							},
						},
					},
				},
			},
			result: []string{},
		},
		{
			name: "nil values",
			input: &DynatraceResponse{
				Result: []DynatraceResult{
					{
						Data: []DynatraceData{
							{
								Values: []*float64{nil, nil, nil},
							},
						},
					},
				},
			},
			result: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kdp := KeptnDynatraceProvider{}
			r := kdp.getResultSlice(tt.input)
			require.Equal(t, tt.result, r)
		})

	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result string
	}{
		{
			name:   "happy path",
			input:  "http://mydttenant.com/api/",
			result: "http://mydttenant.com/api/",
		},
		{
			name:   "missing final /",
			input:  "http://mydttenant.com/api",
			result: "http://mydttenant.com/api/",
		},
		{
			name:   "missing final /api",
			input:  "http://mydttenant.com/",
			result: "http://mydttenant.com/api/",
		},
		{
			name:   "base url",
			input:  "http://mydttenant.com",
			result: "http://mydttenant.com/api/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kdp := KeptnDynatraceProvider{}
			r := kdp.normalizeAPIURL(tt.input)
			require.Equal(t, tt.result, r)
		})

	}
}

func TestEvaluateQuery_MetricWithoutRange(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
		require.Equal(t, "GET", r.Method)
		require.Equal(t, "/api/v2/metrics/query", r.URL.Path)
		require.Equal(t, "my-query", r.URL.Query().Get("metricSelector"))
		require.Equal(t, "now-2h", r.URL.Query().Get("from"))
		require.Equal(t, 1, len(r.Header["Authorization"]))
	}))
	defer svr.Close()

	apiTokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myapitoken",
			Namespace: "",
		},
		Data: map[string][]byte{
			"mykey": []byte("secret-value"),
		},
	}

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "myapitoken",
				},
				Key: "mykey",
			},
			TargetServer: svr.URL,
		},
	}

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query&from=now-2h",
		},
	}

	fakeClient := fake.NewClient(apiTokenSecret)

	kdp := KeptnDynatraceProvider{
		HttpClient: http.Client{},
		Log:        ctrl.Log.WithName("testytest"),
		K8sClient:  fakeClient,
	}

	r, raw, err := kdp.EvaluateQuery(context.TODO(), keptnMetric, p)
	require.Nil(t, err)
	require.Equal(t, []byte(dtpayload), raw)
	require.Equal(t, "50.000000", r)
}

func TestEvaluateQuery_MetricWithRange(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
		require.Equal(t, "GET", r.Method)
		require.Equal(t, "/api/v2/metrics/query", r.URL.Path)
		require.Equal(t, "my-query", r.URL.Query().Get("metricSelector"))
		require.Equal(t, "now-5m", r.URL.Query().Get("from"))
		require.Equal(t, 1, len(r.Header["Authorization"]))
	}))
	defer svr.Close()

	apiTokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myapitoken",
			Namespace: "",
		},
		Data: map[string][]byte{
			"mykey": []byte("secret-value"),
		},
	}

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "myapitoken",
				},
				Key: "mykey",
			},
			TargetServer: svr.URL,
		},
	}

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query",
			Range: &metricsapi.RangeSpec{
				Interval: "5m",
			},
		},
	}

	fakeClient := fake.NewClient(apiTokenSecret)

	kdp := KeptnDynatraceProvider{
		HttpClient: http.Client{},
		Log:        ctrl.Log.WithName("testytest"),
		K8sClient:  fakeClient,
	}

	r, raw, err := kdp.EvaluateQuery(context.TODO(), keptnMetric, p)
	require.Nil(t, err)
	require.Equal(t, []byte(dtpayload), raw)
	require.Equal(t, "50.000000", r)
}

func TestEvaluateQuery_APIError(t *testing.T) {
	errorResponse := []byte("{\"error\":{\"code\":403,\"message\":\"Token is missing required scope. Use one of: metrics.read (Read metrics)\"}}")
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(errorResponse)
		require.Nil(t, err)
	}))
	defer svr.Close()
	secretName, secretKey, secretValue := "secretName", "secretKey", "secretValue"
	apiToken := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "",
		},
		Data: map[string][]byte{
			secretKey: []byte(secretValue),
		},
	}
	kdp := setupTest(apiToken)

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query&from=now-2h",
		},
	}

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: secretName,
				},
				Key: secretKey,
			},
			TargetServer: svr.URL,
		},
	}
	r, raw, e := kdp.EvaluateQuery(context.TODO(), keptnMetric, p)
	require.NotNil(t, e)
	require.Contains(t, e.Error(), "Token is missing required scope.")
	require.Equal(t, "", r)
	t.Log(string(raw))
	require.Equal(t, errorResponse, raw) //we still return the raw answer to help user debug

}

func TestEvaluateQuery_WrongPayloadHandling(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("garbage"))
		require.Nil(t, err)
	}))
	defer svr.Close()
	secretName, secretKey, secretValue := "secretName", "secretKey", "secretValue"
	apiToken := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "",
		},
		Data: map[string][]byte{
			secretKey: []byte(secretValue),
		},
	}

	kdp := setupTest(apiToken)

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query&from=now-2h",
		},
	}

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: secretName,
				},
				Key: secretKey,
			},
			TargetServer: svr.URL,
		},
	}

	r, raw, e := kdp.EvaluateQuery(context.TODO(), keptnMetric, p)
	require.Equal(t, "", r)
	t.Log(string(raw), e)
	require.Equal(t, []byte("garbage"), raw) //we still return the raw answer to help user debug
	require.NotNil(t, e)
}

func TestEvaluateQuery_MissingSecret(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
	}))
	defer svr.Close()
	kdp := setupTest()

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query&from=now-2h",
		},
	}

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: svr.URL,
		},
	}

	_, _, e := kdp.EvaluateQuery(context.TODO(), keptnMetric, p)
	require.NotNil(t, e)
	require.ErrorIs(t, e, ErrSecretKeyRefNotDefined)

}

func TestEvaluateQuery_SecretNotFound(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
	}))
	defer svr.Close()
	kdp := setupTest()

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query&from=now-2h",
		},
	}

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "myapitoken",
				},
				Key: "mykey",
			},
			TargetServer: svr.URL,
		},
	}

	_, _, e := kdp.EvaluateQuery(context.TODO(), keptnMetric, p)
	require.NotNil(t, e)
	require.True(t, errors.IsNotFound(e))
}

func TestEvaluateQuery_RefNotExistingKey(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
	}))
	defer svr.Close()
	secretName, secretKey, secretValue := "secretName", "secretKey", "secretValue"
	apiToken := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "",
		},
		Data: map[string][]byte{
			secretKey: []byte(secretValue),
		},
	}
	kdp := setupTest(apiToken)

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query&from=now-2h",
		},
	}

	missingKey := "key_not_found"
	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: secretName,
				},
				Key: missingKey,
			},
			TargetServer: svr.URL,
		},
	}

	_, _, e := kdp.EvaluateQuery(context.TODO(), keptnMetric, p)
	require.NotNil(t, e)
	require.True(t, strings.Contains(e.Error(), "invalid key "+missingKey))
}

func TestEvaluateQueryForStep(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
		require.Equal(t, "GET", r.Method)
		require.Equal(t, "/api/v2/metrics/query", r.URL.Path)
		require.Equal(t, "GET", r.Method)
		require.Equal(t, "/api/v2/metrics/query", r.URL.Path)
		require.Equal(t, "my-query", r.URL.Query().Get("metricSelector"))
		require.Equal(t, "now-5m", r.URL.Query().Get("from"))
		require.Equal(t, "1m", r.URL.Query().Get("resolution"))
		require.Equal(t, 1, len(r.Header["Authorization"]))
	}))
	defer svr.Close()

	apiTokenSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myapitoken",
			Namespace: "",
		},
		Data: map[string][]byte{
			"mykey": []byte("secret-value"),
		},
	}

	kdp := setupTest(apiTokenSecret)

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}
	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "myapitoken",
				},
				Key: "mykey",
			},
			TargetServer: svr.URL,
		},
	}
	r, raw, err := kdp.EvaluateQueryForStep(context.TODO(), keptnMetric, p)
	require.Nil(t, err)
	require.Equal(t, []byte(dtpayload), raw)
	require.Equal(t, []string{"50.000000", "50.000000", "50.000000", "50.000000", "50.000000", "50.000000", "50.000000", "50.000000"}, r)
}

func TestEvaluateQueryForStep_APIError(t *testing.T) {
	errorResponse := []byte("{\"error\":{\"code\":403,\"message\":\"Token is missing required scope. Use one of: metrics.read (Read metrics)\"}}")
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(errorResponse)
		require.Nil(t, err)
	}))
	defer svr.Close()
	secretName, secretKey, secretValue := "secretName", "secretKey", "secretValue"
	apiToken := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "",
		},
		Data: map[string][]byte{
			secretKey: []byte(secretValue),
		},
	}
	kdp := setupTest(apiToken)

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: secretName,
				},
				Key: secretKey,
			},
			TargetServer: svr.URL,
		},
	}
	r, raw, e := kdp.EvaluateQueryForStep(context.TODO(), keptnMetric, p)
	require.Equal(t, []string(nil), r)
	t.Log(string(raw))
	require.Equal(t, errorResponse, raw) //we still return the raw answer to help user debug
	require.NotNil(t, e)
	require.Contains(t, e.Error(), "Token is missing required scope.")
}

func TestEvaluateQueryForStep_WrongPayloadHandling(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("garbage"))
		require.Nil(t, err)
	}))
	defer svr.Close()
	secretName, secretKey, secretValue := "secretName", "secretKey", "secretValue"
	apiToken := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "",
		},
		Data: map[string][]byte{
			secretKey: []byte(secretValue),
		},
	}

	kdp := setupTest(apiToken)

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}
	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: secretName,
				},
				Key: secretKey,
			},
			TargetServer: svr.URL,
		},
	}
	r, raw, e := kdp.EvaluateQueryForStep(context.TODO(), keptnMetric, p)
	require.Equal(t, []string(nil), r)
	t.Log(string(raw), e)
	require.Equal(t, []byte("garbage"), raw) //we still return the raw answer to help user debug
	require.NotNil(t, e)
}

func TestEvaluateQueryForStep_MissingSecret(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
	}))
	defer svr.Close()
	kdp := setupTest()

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query",
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
	_, _, e := kdp.EvaluateQueryForStep(context.TODO(), keptnMetric, p)
	require.NotNil(t, e)
	require.ErrorIs(t, e, ErrSecretKeyRefNotDefined)
}

func TestEvaluateQueryForStep_SecretNotFound(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
	}))
	defer svr.Close()
	kdp := setupTest()

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "myapitoken",
				},
				Key: "mykey",
			},
			TargetServer: svr.URL,
		},
	}
	_, _, e := kdp.EvaluateQueryForStep(context.TODO(), keptnMetric, p)
	require.NotNil(t, e)
	require.True(t, errors.IsNotFound(e))
}

func TestEvaluateQueryForStep_RefNotExistingKey(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
	}))
	defer svr.Close()
	secretName, secretKey, secretValue := "secretName", "secretKey", "secretValue"
	apiToken := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "",
		},
		Data: map[string][]byte{
			secretKey: []byte(secretValue),
		},
	}
	kdp := setupTest(apiToken)

	keptnMetric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "my-query",
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
	}

	missingKey := "key_not_found"
	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: secretName,
				},
				Key: missingKey,
			},
			TargetServer: svr.URL,
		},
	}
	_, _, e := kdp.EvaluateQueryForStep(context.TODO(), keptnMetric, p)
	require.NotNil(t, e)
	require.True(t, strings.Contains(e.Error(), "invalid key "+missingKey))
}

func setupTest(objs ...client.Object) KeptnDynatraceProvider {

	fakeClient := fake.NewClient(objs...)

	kdp := KeptnDynatraceProvider{
		HttpClient: http.Client{},
		Log:        ctrl.Log.WithName("testytest"),
		K8sClient:  fakeClient,
	}

	return kdp
}
