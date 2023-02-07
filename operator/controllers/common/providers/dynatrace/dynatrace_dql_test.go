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
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/providers/dynatrace/client/fake"
	"github.com/stretchr/testify/require"
	"k8s.io/klog/v2"
)

const dqlRequestHandler = `{"requestToken": "my-token"}`

const dqlPayload = "{\"state\":\"SUCCEEDED\",\"result\":{\"records\":[{\"value\":{\"count\":1,\"sum\":36.50,\"min\":36.78336878333334,\"avg\":36.78336878333334,\"max\":36.78336878333334},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:11:00.000Z\",\"end\":\"2023-01-31T09:12:00.`00Z\"},\"Container\":\"frontend\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"hipstershop\",\"k8s.pod.uid\":\"632df64d-474c-4410-968d-666f639ad358\"}],\"types\":[{\"mappings\":{\"value\":{\"type\":\"summary_stats\"},\"metric.key\":{\"type\":\"string\"},\"timeframe\":{\"type\":\"timeframe\"},\"Container\":{\"type\":\"string\"},\"host.name\":{\"type\":\"string\"},\"k8s.namespace.name\":{\"type\":\"string\"},\"k8s.pod.uid\":{\"type\":\"string\"}},\"indexRange\":[0,1]}]}}"
const dqlPayloadNotFinished = "{\"state\":\"\",\"result\":{\"records\":[{\"value\":{\"count\":1,\"sum\":36.50,\"min\":36.78336878333334,\"avg\":36.50,\"max\":36.50},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:11:00.000Z\",\"end\":\"2023-01-31T09:12:00.`00Z\"},\"Container\":\"frontend\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"hipstershop\",\"k8s.pod.uid\":\"632df64d-474c-4410-968d-666f639ad358\"}],\"types\":[{\"mappings\":{\"value\":{\"type\":\"summary_stats\"},\"metric.key\":{\"type\":\"string\"},\"timeframe\":{\"type\":\"timeframe\"},\"Container\":{\"type\":\"string\"},\"host.name\":{\"type\":\"string\"},\"k8s.namespace.name\":{\"type\":\"string\"},\"k8s.pod.uid\":{\"type\":\"string\"}},\"indexRange\":[0,1]}]}}"
const dqlPayloadTooManyItems = "{\"records\":[{\"value\":{\"count\":1,\"sum\":6.293549483333334,\"min\":6.293549483333334,\"avg\":6.293549483333334,\"max\":6.293549483333334},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:07:00.000Z\",\"end\":\"2023-01-31T09:08:00.000Z\"},\"Container\":\"loginservice\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"easytrade\",\"k8s.pod.uid\":\"fc084e57-11a0-4a95-b8a0-76191c31d839\"},{\"value\":{\"count\":1,\"sum\":1.0421756,\"min\":1.0421756,\"avg\":1.0421756,\"max\":1.0421756},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:07:00.000Z\",\"end\":\"2023-01-31T09:08:00.000Z\"},\"Container\":\"frontendreverseproxy\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"easytrade\",\"k8s.pod.uid\":\"41b5d6e0-98fc-4dce-a1b4-bb269a03d72b\"},{\"value\":{\"count\":1,\"sum\":6.3881383000000005,\"min\":6.3881383000000005,\"avg\":6.3881383000000005,\"max\":6.3881383000000005},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:07:00.000Z\",\"end\":\"2023-01-31T09:08:00.000Z\"},\"Container\":\"shippingservice\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"hipstershop\",\"k8s.pod.uid\":\"96fcf9d7-748a-47f7-b1b3-ca6427e20edd\"}],\"types\":[{\"mappings\":{\"value\":{\"type\":\"summary_stats\"},\"metric.key\":{\"type\":\"string\"},\"timeframe\":{\"type\":\"timeframe\"},\"Container\":{\"type\":\"string\"},\"host.name\":{\"type\":\"string\"},\"k8s.namespace.name\":{\"type\":\"string\"},\"k8s.pod.uid\":{\"type\":\"string\"}},\"indexRange\":[0,3]}]}"

const dqlOAuthCorrect = "{\"scope\":\"scope_1\",\"token_type\":\"Bearer\",\"expires_in\":300,\"access_token\":\"token\"}"
const dqlOAuthWrong = "{ \"token\":\"<token>\"}"

func TestGetDQL(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayload), nil
		}

		return nil, errors.New("unexpected path")
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, raw, err := dqlProvider.EvaluateQuery(context.TODO(), klcv1alpha2.Objective{
		Name:             "",
		Query:            "",
		EvaluationTarget: "",
	}, klcv1alpha2.KeptnEvaluationProvider{
		Spec: klcv1alpha2.KeptnEvaluationProviderSpec{},
	})

	require.Nil(t, err)
	require.NotEmpty(t, raw)
	require.Equal(t, "36.50", result)
}

func TestGetDQLTimeout(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadNotFinished), nil
		}

		return nil, errors.New("unexpected path")
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
		result, raw, err := dqlProvider.EvaluateQuery(context.TODO(), klcv1alpha2.Objective{
			Name:             "",
			Query:            "",
			EvaluationTarget: "",
		}, klcv1alpha2.KeptnEvaluationProvider{
			Spec: klcv1alpha2.KeptnEvaluationProviderSpec{},
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
}

func TestDQL_TooManyItems(t *testing.T) {

}
