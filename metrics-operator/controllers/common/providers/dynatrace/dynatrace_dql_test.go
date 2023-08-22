// nolint: dupl
package dynatrace

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/dynatrace/client/fake"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	k8sfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const dqlRequestHandler = `{"requestToken": "my-token"}`

const dqlPayload = "{\"state\":\"SUCCEEDED\",\"result\":{\"records\":[{\"value\":{\"count\":1,\"sum\":36.50,\"min\":36.50,\"avg\":36.50,\"max\":36.50},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:11:00.000Z\",\"end\":\"2023-01-31T09:12:00.`00Z\"},\"Container\":\"frontend\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"hipstershop\",\"k8s.pod.uid\":\"632df64d-474c-4410-968d-666f639ad358\"}],\"types\":[{\"mappings\":{\"value\":{\"type\":\"summary_stats\"},\"metric.key\":{\"type\":\"string\"},\"timeframe\":{\"type\":\"timeframe\"},\"Container\":{\"type\":\"string\"},\"host.name\":{\"type\":\"string\"},\"k8s.namespace.name\":{\"type\":\"string\"},\"k8s.pod.uid\":{\"type\":\"string\"}},\"indexRange\":[0,1]}]}}"
const dqlPayloadEmpty = "{\"state\":\"SUCCEEDED\",\"result\":{\"records\":[],\"types\":[{\"mappings\":{\"value\":{\"type\":\"summary_stats\"},\"metric.key\":{\"type\":\"string\"},\"timeframe\":{\"type\":\"timeframe\"},\"Container\":{\"type\":\"string\"},\"host.name\":{\"type\":\"string\"},\"k8s.namespace.name\":{\"type\":\"string\"},\"k8s.pod.uid\":{\"type\":\"string\"}},\"indexRange\":[0,1]}]}}"
const dqlPayloadNotFinished = "{\"state\":\"\",\"result\":{\"records\":[{\"value\":{\"count\":1,\"sum\":36.50,\"min\":36.78336878333334,\"avg\":36.50,\"max\":36.50},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:11:00.000Z\",\"end\":\"2023-01-31T09:12:00.`00Z\"},\"Container\":\"frontend\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"hipstershop\",\"k8s.pod.uid\":\"632df64d-474c-4410-968d-666f639ad358\"}],\"types\":[{\"mappings\":{\"value\":{\"type\":\"summary_stats\"},\"metric.key\":{\"type\":\"string\"},\"timeframe\":{\"type\":\"timeframe\"},\"Container\":{\"type\":\"string\"},\"host.name\":{\"type\":\"string\"},\"k8s.namespace.name\":{\"type\":\"string\"},\"k8s.pod.uid\":{\"type\":\"string\"}},\"indexRange\":[0,1]}]}}"
const dqlPayloadError = "{\"error\":{\"code\":403,\"message\":\"Token is missing required scope\"}}"

const dqlPayloadTooManyItems = "{\"state\":\"SUCCEEDED\",\"result\":{\"records\":[{\"value\":{\"count\":1,\"sum\":6.293549483333334,\"min\":6.293549483333334,\"avg\":6.293549483333334,\"max\":6.293549483333334},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:07:00.000Z\",\"end\":\"2023-01-31T09:08:00.000Z\"},\"Container\":\"loginservice\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"easytrade\",\"k8s.pod.uid\":\"fc084e57-11a0-4a95-b8a0-76191c31d839\"},{\"value\":{\"count\":1,\"sum\":1.0421756,\"min\":1.0421756,\"avg\":1.0421756,\"max\":1.0421756},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:07:00.000Z\",\"end\":\"2023-01-31T09:08:00.000Z\"},\"Container\":\"frontendreverseproxy\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"easytrade\",\"k8s.pod.uid\":\"41b5d6e0-98fc-4dce-a1b4-bb269a03d72b\"},{\"value\":{\"count\":1,\"sum\":6.3881383000000005,\"min\":6.3881383000000005,\"avg\":6.3881383000000005,\"max\":6.3881383000000005},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:07:00.000Z\",\"end\":\"2023-01-31T09:08:00.000Z\"},\"Container\":\"shippingservice\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"hipstershop\",\"k8s.pod.uid\":\"96fcf9d7-748a-47f7-b1b3-ca6427e20edd\"}],\"types\":[{\"mappings\":{\"value\":{\"type\":\"summary_stats\"},\"metric.key\":{\"type\":\"string\"},\"timeframe\":{\"type\":\"timeframe\"},\"Container\":{\"type\":\"string\"},\"host.name\":{\"type\":\"string\"},\"k8s.namespace.name\":{\"type\":\"string\"},\"k8s.pod.uid\":{\"type\":\"string\"}},\"indexRange\":[0,3]}]}}"

var ErrUnexpected = errors.New("unexpected path")

//nolint:dupl
func TestGetDQL_EvaluateQuery(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayload), nil
		}
		return nil, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, raw, err := dqlProvider.EvaluateQuery(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{Query: ""},
		},
		metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.Nil(t, err)
	require.NotEmpty(t, raw)
	require.Equal(t, "36.500000", result)

	require.Len(t, mockClient.DoCalls(), 2)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
	require.Contains(t, mockClient.DoCalls()[1].Path, "query:poll")
}

//nolint:dupl
func TestGetDQLMultipleRecords_EvaluateQuery(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadTooManyItems), nil
		}

		return nil, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, raw, err := dqlProvider.EvaluateQuery(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{Query: ""},
		}, metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.Nil(t, err)
	require.NotEmpty(t, raw)
	require.Equal(t, "6.293549", result)

	require.Len(t, mockClient.DoCalls(), 2)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
	require.Contains(t, mockClient.DoCalls()[1].Path, "query:poll")
}

func TestGetDQLAPIError_EvaluateQuery(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadError), nil
		}

		return nil, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, raw, err := dqlProvider.EvaluateQuery(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{Query: ""},
		}, metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "Token is missing required scope")
	require.Empty(t, raw)
	require.Empty(t, result)

	require.Len(t, mockClient.DoCalls(), 2)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
	require.Contains(t, mockClient.DoCalls()[1].Path, "query:poll")
}

func TestGetDQLTimeout_EvaluateQuery(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadNotFinished), nil
		}

		return nil, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	mockClock := clock.NewMock()
	dqlProvider.clock = mockClock

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		result, raw, err := dqlProvider.EvaluateQuery(context.TODO(),
			metricsapi.KeptnMetric{
				Spec: metricsapi.KeptnMetricSpec{Query: ""},
			}, metricsapi.KeptnMetricsProvider{
				Spec: metricsapi.KeptnMetricsProviderSpec{},
			})

		require.ErrorIs(t, err, ErrDQLQueryTimeout)
		require.Empty(t, raw)
		require.Empty(t, result)
	}(&wg)

	// wait for the mockClient to be called at least one time before adding to the clock
	require.Eventually(t, func() bool {
		return len(mockClient.DoCalls()) > 0
	}, 5*time.Second, 100*time.Millisecond)

	mockClock.Add(retryFetchInterval * (maxRetries + 1))
	wg.Wait()
	require.Len(t, mockClient.DoCalls(), maxRetries+1)
}

func TestGetDQLCannotPostQuery_EvaluateQuery(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return nil, errors.New("oops")
		}

		return nil, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	mockClock := clock.NewMock()
	dqlProvider.clock = mockClock

	result, raw, err := dqlProvider.EvaluateQuery(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{Query: ""},
		},
		metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.NotNil(t, err, err)
	require.Empty(t, raw)
	require.Empty(t, result)

	require.Len(t, mockClient.DoCalls(), 1)
}

func TestDQLInitClientWithSecret_EvaluateQuery(t *testing.T) {

	namespace := "keptn-lifecycle-toolkit-system"

	mySecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"my-key": []byte("dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"),
		},
		Type: corev1.SecretTypeOpaque,
	}
	fakeClient := k8sfake.NewClientBuilder().WithScheme(clientgoscheme.Scheme).WithObjects(mySecret).Build()

	dqlProvider := NewKeptnDynatraceDQLProvider(
		fakeClient,
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	require.NotNil(t, dqlProvider)

	err := dqlProvider.ensureDTClientIsSetUp(context.TODO(), metricsapi.KeptnMetricsProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dql",
			Namespace: namespace,
		},
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "my-secret",
				},
				Key: "my-key",
			},
		},
	})

	require.Nil(t, err)
	require.NotNil(t, dqlProvider.dtClient)
}

func TestGetDQLEmptyPayload_EvaluateQuery(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadEmpty), nil
		}
		return nil, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, raw, err := dqlProvider.EvaluateQuery(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{Query: ""},
		},
		metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.Equal(t, ErrInvalidResult, err)
	require.Empty(t, raw)
	require.Equal(t, "", result)
}

func TestGetDQL_EvaluateQueryForStep(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayload), nil
		}
		return nil, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, raw, err := dqlProvider.EvaluateQueryForStep(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{Query: ""},
		},
		metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.Nil(t, err)
	require.NotEmpty(t, raw)
	require.Equal(t, []string{"36.500000"}, result)

	require.Len(t, mockClient.DoCalls(), 2)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
	require.Contains(t, mockClient.DoCalls()[1].Path, "query:poll")
}

//nolint:dupl
func TestGetDQLMultipleRecords_EvaluateQueryForStep(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadTooManyItems), nil
		}

		return nil, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, raw, err := dqlProvider.EvaluateQueryForStep(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{Query: ""},
		}, metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.Nil(t, err)
	require.NotEmpty(t, raw)
	require.Equal(t, []string{"6.293549", "1.042176", "6.388138"}, result)

	require.Len(t, mockClient.DoCalls(), 2)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
	require.Contains(t, mockClient.DoCalls()[1].Path, "query:poll")
}

func TestGetDQLAPIError_EvaluateQueryForStep(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadError), nil
		}

		return nil, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, raw, err := dqlProvider.EvaluateQueryForStep(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{Query: ""},
		}, metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "Token is missing required scope")
	require.Empty(t, raw)
	require.Empty(t, result)

	require.Len(t, mockClient.DoCalls(), 2)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
	require.Contains(t, mockClient.DoCalls()[1].Path, "query:poll")
}

func TestGetDQLTimeout_EvaluateQueryForStep(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadNotFinished), nil
		}

		return nil, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	mockClock := clock.NewMock()
	dqlProvider.clock = mockClock

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		result, raw, err := dqlProvider.EvaluateQueryForStep(context.TODO(),
			metricsapi.KeptnMetric{
				Spec: metricsapi.KeptnMetricSpec{Query: ""},
			}, metricsapi.KeptnMetricsProvider{
				Spec: metricsapi.KeptnMetricsProviderSpec{},
			})

		require.ErrorIs(t, err, ErrDQLQueryTimeout)
		require.Empty(t, raw)
		require.Empty(t, result)
	}(&wg)

	// wait for the mockClient to be called at least one time before adding to the clock
	require.Eventually(t, func() bool {
		return len(mockClient.DoCalls()) > 0
	}, 5*time.Second, 100*time.Millisecond)

	mockClock.Add(retryFetchInterval * (maxRetries + 1))
	wg.Wait()
	require.Len(t, mockClient.DoCalls(), maxRetries+1)
}

func TestGetDQLCannotPostQuery_EvaluateQueryForStep(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return nil, errors.New("oops")
		}

		return nil, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	mockClock := clock.NewMock()
	dqlProvider.clock = mockClock

	result, raw, err := dqlProvider.EvaluateQueryForStep(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{Query: ""},
		},
		metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.NotNil(t, err, err)
	require.Empty(t, raw)
	require.Empty(t, result)

	require.Len(t, mockClient.DoCalls(), 1)
}

func TestDQLInitClientWithSecret_EvaluateQueryForStep(t *testing.T) {

	namespace := "keptn-lifecycle-toolkit-system"

	mySecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"my-key": []byte("dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"),
		},
		Type: corev1.SecretTypeOpaque,
	}
	fakeClient := k8sfake.NewClientBuilder().WithScheme(clientgoscheme.Scheme).WithObjects(mySecret).Build()

	dqlProvider := NewKeptnDynatraceDQLProvider(
		fakeClient,
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	require.NotNil(t, dqlProvider)

	err := dqlProvider.ensureDTClientIsSetUp(context.TODO(), metricsapi.KeptnMetricsProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dql",
			Namespace: namespace,
		},
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "my-secret",
				},
				Key: "my-key",
			},
		},
	})

	require.Nil(t, err)
	require.NotNil(t, dqlProvider.dtClient)
}

func TestGetDQLEmptyPayload_EvaluateQueryForStep(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadEmpty), nil
		}
		return nil, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, raw, err := dqlProvider.EvaluateQueryForStep(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{Query: ""},
		},
		metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.Equal(t, ErrInvalidResult, err)
	require.Empty(t, raw)
	require.Equal(t, []string(nil), result)
}

func TestGetResultForSlice_HappyPath(t *testing.T) {

	namespace := "keptn-lifecycle-toolkit-system"

	mySecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"my-key": []byte("dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"),
		},
		Type: corev1.SecretTypeOpaque,
	}
	fakeClient := k8sfake.NewClientBuilder().WithScheme(clientgoscheme.Scheme).WithObjects(mySecret).Build()
	dqlProvider := NewKeptnDynatraceDQLProvider(
		fakeClient,
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)
	result := &DQLResult{
		Records: []DQLRecord{
			{
				Value: DQLMetric{
					Count: 1,
					Sum:   25.0,
					Min:   25.0,
					Avg:   25.0,
					Max:   25.0,
				},
			},
			{
				Value: DQLMetric{
					Count: 1,
					Sum:   13.0,
					Min:   13.0,
					Avg:   13.0,
					Max:   13.0,
				},
			},
		},
	}
	resultSlice := dqlProvider.getResultSlice(result)
	require.NotZero(t, resultSlice)
	require.Equal(t, []string{"25.000000", "13.000000"}, resultSlice)
}

func TestGetResultForSlice_Empty(t *testing.T) {

	namespace := "keptn-lifecycle-toolkit-system"

	mySecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"my-key": []byte("dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"),
		},
		Type: corev1.SecretTypeOpaque,
	}
	fakeClient := k8sfake.NewClientBuilder().WithScheme(clientgoscheme.Scheme).WithObjects(mySecret).Build()
	dqlProvider := NewKeptnDynatraceDQLProvider(
		fakeClient,
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)
	result := &DQLResult{
		Records: []DQLRecord{},
	}
	resultSlice := dqlProvider.getResultSlice(result)
	require.Equal(t, []string(nil), resultSlice)
}
