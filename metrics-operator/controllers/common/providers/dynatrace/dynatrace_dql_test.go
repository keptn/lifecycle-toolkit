// nolint: dupl
package dynatrace

import (
	"context"
	"encoding/json"
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

const dqlPayload = `{
  "state": "SUCCEEDED",
  "progress": 100,
  "result": {
    "records": [
      {
        "timeframe": {
          "start": "2023-08-23T11:00:00.000000000Z",
          "end": "2023-08-23T14:00:00.000000000Z"
        },
        "interval": "3600000000000",
        "avg(dt.host.cpu.usage)": [
          20.44886593058413,
          20.43724084597563,
          20.417480020188446
        ]
      }
    ],
    "types": [
      {
        "indexRange": [
          0,
          0
        ],
        "mappings": {
          "timeframe": {
            "type": "timeframe"
          },
          "interval": {
            "type": "duration"
          },
          "avg(dt.host.cpu.usage)": {
            "type": "array",
            "types": [
              {
                "indexRange": [
                  0,
                  2
                ],
                "mappings": {
                  "element": {
                    "type": "double"
                  }
                }
              }
            ]
          }
        }
      }
    ],
    "metadata": {
      "grail": {
        "canonicalQuery": "timeseries interval:1h, avg(dt.host.cpu.usage)",
        "timezone": "Z",
        "query": "timeseries avg(dt.host.cpu.usage), interval:1h",
        "scannedRecords": 0,
        "dqlVersion": "V1_0",
        "scannedBytes": 0,
        "analysisTimeframe": {
          "start": "2023-08-23T11:00:00.000000000Z",
          "end": "2023-08-23T14:00:00.000000000Z"
        },
        "locale": "en_US",
        "executionTimeMilliseconds": 72,
        "notifications": [],
        "queryId": "76f215bf-7ef0-4374-8fd6-4677aaf4d816",
        "sampled": false
      },
      "metrics": [
        {
          "metric.key": "dt.host.cpu.usage",
          "displayName": "CPU usage %",
          "description": "Percentage of CPU time when CPU was utilized. A value close to 100% means most host processing resources are in use, and host CPUs can’t handle additional work",
          "unit": "%",
          "fieldName": "avg(dt.host.cpu.usage)"
        }
      ]
    }
  }
}`

const dqlPayloadNotFinished = `{
  "state": "",
  "progress": 100,
  "result": {
    "records": [],
    "types": [
      {
        "indexRange": [
          0,
          0
        ],
        "mappings": {
          "timeframe": {
            "type": "timeframe"
          },
          "interval": {
            "type": "duration"
          },
          "avg(dt.host.cpu.usage)": {
            "type": "array",
            "types": [
              {
                "indexRange": [
                  0,
                  2
                ],
                "mappings": {
                  "element": {
                    "type": "double"
                  }
                }
              }
            ]
          }
        }
      }
    ],
    "metadata": {
      "grail": {
        "canonicalQuery": "timeseries interval:1h, avg(dt.host.cpu.usage)",
        "timezone": "Z",
        "query": "timeseries avg(dt.host.cpu.usage), interval:1h",
        "scannedRecords": 0,
        "dqlVersion": "V1_0",
        "scannedBytes": 0,
        "analysisTimeframe": {
          "start": "2023-08-23T11:00:00.000000000Z",
          "end": "2023-08-23T14:00:00.000000000Z"
        },
        "locale": "en_US",
        "executionTimeMilliseconds": 72,
        "notifications": [],
        "queryId": "76f215bf-7ef0-4374-8fd6-4677aaf4d816",
        "sampled": false
      },
      "metrics": [
        {
          "metric.key": "dt.host.cpu.usage",
          "displayName": "CPU usage %",
          "description": "Percentage of CPU time when CPU was utilized. A value close to 100% means most host processing resources are in use, and host CPUs can’t handle additional work",
          "unit": "%",
          "fieldName": "avg(dt.host.cpu.usage)"
        }
      ]
    }
  }
}`

const dqlPayloadEmpty = `{
  "state": "SUCCEEDED",
  "progress": 100,
  "result": {
    "records": [],
    "types": [
      {
        "indexRange": [
          0,
          0
        ],
        "mappings": {
          "timeframe": {
            "type": "timeframe"
          },
          "interval": {
            "type": "duration"
          },
          "avg(dt.host.cpu.usage)": {
            "type": "array",
            "types": [
              {
                "indexRange": [
                  0,
                  2
                ],
                "mappings": {
                  "element": {
                    "type": "double"
                  }
                }
              }
            ]
          }
        }
      }
    ],
    "metadata": {
      "grail": {
        "canonicalQuery": "timeseries interval:1h, avg(dt.host.cpu.usage)",
        "timezone": "Z",
        "query": "timeseries avg(dt.host.cpu.usage), interval:1h",
        "scannedRecords": 0,
        "dqlVersion": "V1_0",
        "scannedBytes": 0,
        "analysisTimeframe": {
          "start": "2023-08-23T11:00:00.000000000Z",
          "end": "2023-08-23T14:00:00.000000000Z"
        },
        "locale": "en_US",
        "executionTimeMilliseconds": 72,
        "notifications": [],
        "queryId": "76f215bf-7ef0-4374-8fd6-4677aaf4d816",
        "sampled": false
      },
      "metrics": [
        {
          "metric.key": "dt.host.cpu.usage",
          "displayName": "CPU usage %",
          "description": "Percentage of CPU time when CPU was utilized. A value close to 100% means most host processing resources are in use, and host CPUs can’t handle additional work",
          "unit": "%",
          "fieldName": "avg(dt.host.cpu.usage)"
        }
      ]
    }
  }
}`

const dqlPayloadMultipleRecords = `{
  "state": "SUCCEEDED",
  "progress": 100,
  "result": {
    "records": [
      {
        "timeframe": {
          "start": "2023-08-23T11:00:00.000000000Z",
          "end": "2023-08-23T14:00:00.000000000Z"
        },
        "interval": "3600000000000",
        "avg(dt.host.cpu.usage)": [
          20.44886593058413,
          20.43724084597563,
          20.417480020188446
        ]
      },
	{
        "timeframe": {
          "start": "2023-08-23T11:00:00.000000000Z",
          "end": "2023-08-23T14:00:00.000000000Z"
        },
        "interval": "3600000000000",
        "avg(some-other-metric)": [
          30.44886593058413,
          30.43724084597563,
          30.417480020188446
        ]
      }
    ],
    "types": [
      {
        "indexRange": [
          0,
          0
        ],
        "mappings": {
          "timeframe": {
            "type": "timeframe"
          },
          "interval": {
            "type": "duration"
          },
          "avg(dt.host.cpu.usage)": {
            "type": "array",
            "types": [
              {
                "indexRange": [
                  0,
                  2
                ],
                "mappings": {
                  "element": {
                    "type": "double"
                  }
                }
              }
            ]
          }
        }
      }
    ],
    "metadata": {
      "grail": {
        "canonicalQuery": "timeseries interval:1h, avg(dt.host.cpu.usage)",
        "timezone": "Z",
        "query": "timeseries avg(dt.host.cpu.usage), interval:1h",
        "scannedRecords": 0,
        "dqlVersion": "V1_0",
        "scannedBytes": 0,
        "analysisTimeframe": {
          "start": "2023-08-23T11:00:00.000000000Z",
          "end": "2023-08-23T14:00:00.000000000Z"
        },
        "locale": "en_US",
        "executionTimeMilliseconds": 72,
        "notifications": [],
        "queryId": "76f215bf-7ef0-4374-8fd6-4677aaf4d816",
        "sampled": false
      },
      "metrics": [
        {
          "metric.key": "dt.host.cpu.usage",
          "displayName": "CPU usage %",
          "description": "Percentage of CPU time when CPU was utilized. A value close to 100% means most host processing resources are in use, and host CPUs can’t handle additional work",
          "unit": "%",
          "fieldName": "avg(dt.host.cpu.usage)"
        }
      ]
    }
  }
}`

// const dqlPayload = "{\"state\":\"SUCCEEDED\",\"result\":{\"records\":[{\"value\":{\"count\":1,\"sum\":36.50,\"min\":36.50,\"avg\":36.50,\"max\":36.50},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:11:00.000Z\",\"end\":\"2023-01-31T09:12:00.`00Z\"},\"Container\":\"frontend\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"hipstershop\",\"k8s.pod.uid\":\"632df64d-474c-4410-968d-666f639ad358\"}],\"types\":[{\"mappings\":{\"value\":{\"type\":\"summary_stats\"},\"metric.key\":{\"type\":\"string\"},\"timeframe\":{\"type\":\"timeframe\"},\"Container\":{\"type\":\"string\"},\"host.name\":{\"type\":\"string\"},\"k8s.namespace.name\":{\"type\":\"string\"},\"k8s.pod.uid\":{\"type\":\"string\"}},\"indexRange\":[0,1]}]}}"
// const dqlPayloadEmpty = "{\"state\":\"SUCCEEDED\",\"result\":{\"records\":[],\"types\":[{\"mappings\":{\"value\":{\"type\":\"summary_stats\"},\"metric.key\":{\"type\":\"string\"},\"timeframe\":{\"type\":\"timeframe\"},\"Container\":{\"type\":\"string\"},\"host.name\":{\"type\":\"string\"},\"k8s.namespace.name\":{\"type\":\"string\"},\"k8s.pod.uid\":{\"type\":\"string\"}},\"indexRange\":[0,1]}]}}"
// const dqlPayloadNotFinished = "{\"state\":\"\",\"result\":{\"records\":[{\"value\":{\"count\":1,\"sum\":36.50,\"min\":36.78336878333334,\"avg\":36.50,\"max\":36.50},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:11:00.000Z\",\"end\":\"2023-01-31T09:12:00.`00Z\"},\"Container\":\"frontend\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"hipstershop\",\"k8s.pod.uid\":\"632df64d-474c-4410-968d-666f639ad358\"}],\"types\":[{\"mappings\":{\"value\":{\"type\":\"summary_stats\"},\"metric.key\":{\"type\":\"string\"},\"timeframe\":{\"type\":\"timeframe\"},\"Container\":{\"type\":\"string\"},\"host.name\":{\"type\":\"string\"},\"k8s.namespace.name\":{\"type\":\"string\"},\"k8s.pod.uid\":{\"type\":\"string\"}},\"indexRange\":[0,1]}]}}"
const dqlPayloadError = "{\"error\":{\"code\":403,\"message\":\"Token is missing required scope\"}}"

var ErrUnexpected = errors.New("unexpected path")

//nolint:dupl
func TestGetDQL_EvaluateQueryResultAvailableImmediately(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlPayload), 200, nil
		}
		return nil, 0, ErrUnexpected
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
	require.Equal(t, "20.417480", result)

	require.Len(t, mockClient.DoCalls(), 1)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
}

func TestGetDQL_EvaluateQueryResultAvailableViaRequestToken(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), 202, nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayload), 202, nil
		}
		return nil, 0, ErrUnexpected
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
	require.Equal(t, "20.417480", result)

	require.Len(t, mockClient.DoCalls(), 2)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
	require.Contains(t, mockClient.DoCalls()[1].Path, "query:poll")
}

func TestGetDQL_EvaluateQueryWithRange(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			reqPayload := &DQLRequest{}
			err := json.Unmarshal(payload, reqPayload)

			require.Nil(t, err)

			parsedStartTime, err := time.Parse(time.RFC3339, reqPayload.DefaultTimeframeStart)
			require.Nil(t, err)
			parsedEndTime, err := time.Parse(time.RFC3339, reqPayload.DefaultTimeframeEnd)
			require.Nil(t, err)

			require.WithinDuration(t, time.Now().UTC(), parsedEndTime, 5*time.Second)
			require.Equal(t, 5*time.Minute, parsedEndTime.Sub(parsedStartTime))
			return []byte(dqlPayload), 200, nil
		}
		return nil, 0, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, raw, err := dqlProvider.EvaluateQuery(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "",
				Range: &metricsapi.RangeSpec{
					Interval: "5m",
				},
			},
		},
		metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.Nil(t, err)
	require.NotEmpty(t, raw)
	require.Equal(t, "20.417480", result)

	require.Len(t, mockClient.DoCalls(), 1)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
}

func TestGetDQL_EvaluateQueryWithWrongRange(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlPayload), 200, nil
		}
		// the second if can be left out as in this case the dql provider will return the result without needing to call query:poll
		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayload), 202, nil
		}
		return nil, 0, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, raw, err := dqlProvider.EvaluateQuery(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "",
				Range: &metricsapi.RangeSpec{
					Interval: "5mins",
				},
			},
		},
		metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "time: unknown unit \"mins\" in duration \"5mins\"")
	require.Empty(t, raw)
	require.Empty(t, result)
}

//nolint:dupl
func TestGetDQLMultipleRecords_EvaluateQuery(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), 202, nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadMultipleRecords), 202, nil
		}

		return nil, 0, ErrUnexpected
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
	require.Equal(t, "20.417480", result)

	require.Len(t, mockClient.DoCalls(), 2)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
	require.Contains(t, mockClient.DoCalls()[1].Path, "query:poll")
}

func TestGetDQLAPIError_EvaluateQuery(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), 202, nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadError), 202, nil
		}

		return nil, 0, ErrUnexpected
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

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), 202, nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadNotFinished), 202, nil
		}

		return nil, 0, ErrUnexpected
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

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return nil, 0, errors.New("oops")
		}

		return nil, 0, ErrUnexpected
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

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), 202, nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadEmpty), 202, nil
		}
		return nil, 0, ErrUnexpected
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

func TestGetDQL_EvaluateQueryForStep200(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlPayload), 200, nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayload), 202, nil
		}
		return nil, 0, ErrUnexpected
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
	require.Equal(t, []string{"20.448866", "20.437241", "20.417480"}, result)

	require.Len(t, mockClient.DoCalls(), 1)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
}

func TestGetDQL_EvaluateQueryForStep202(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), 202, nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayload), 202, nil
		}
		return nil, 0, ErrUnexpected
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
	require.Equal(t, []string{"20.448866", "20.437241", "20.417480"}, result)

	require.Len(t, mockClient.DoCalls(), 2)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
	require.Contains(t, mockClient.DoCalls()[1].Path, "query:poll")
}

func TestGetDQL_EvaluateQueryForStepWithRange(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), 202, nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayload), 202, nil
		}
		return nil, 0, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, raw, err := dqlProvider.EvaluateQueryForStep(context.TODO(),
		metricsapi.KeptnMetric{
			Spec: metricsapi.KeptnMetricSpec{
				Query: "",
				Range: &metricsapi.RangeSpec{
					Interval: "5m",
				},
			},
		},
		metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.Nil(t, err)
	require.NotEmpty(t, raw)
	require.Equal(t, []string{"20.448866", "20.437241", "20.417480"}, result)

	require.Len(t, mockClient.DoCalls(), 2)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
	require.Contains(t, mockClient.DoCalls()[1].Path, "query:poll")
}

//nolint:dupl
func TestGetDQLMultipleRecords_EvaluateQueryForStep(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), 202, nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadMultipleRecords), 202, nil
		}

		return nil, 0, ErrUnexpected
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
	require.Equal(t, []string{"20.448866", "20.437241", "20.417480"}, result)

	require.Len(t, mockClient.DoCalls(), 2)
	require.Contains(t, mockClient.DoCalls()[0].Path, "query:execute")
	require.Contains(t, mockClient.DoCalls()[1].Path, "query:poll")
}

func TestGetDQLAPIError_EvaluateQueryForStep(t *testing.T) {

	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), 202, nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadError), 202, nil
		}

		return nil, 0, ErrUnexpected
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

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), 202, nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadNotFinished), 202, nil
		}

		return nil, 0, ErrUnexpected
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

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return nil, 0, errors.New("oops")
		}

		return nil, 0, ErrUnexpected
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

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		if strings.Contains(path, "query:execute") {
			return []byte(dqlRequestHandler), 202, nil
		}

		if strings.Contains(path, "query:poll") {
			return []byte(dqlPayloadEmpty), 202, nil
		}
		return nil, 0, ErrUnexpected
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

func Test_keptnDynatraceDQLProvider_FetchAnalysisValue(t *testing.T) {
	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		require.Contains(t, path, "query:execute")
		return []byte(dqlPayload), 200, nil
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, err := dqlProvider.FetchAnalysisValue(context.TODO(), "timeseries (dt.host.cpu)",
		metricsapi.Analysis{
			Status: metricsapi.AnalysisStatus{},
		},
		&metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.Nil(t, err)
	require.Equal(t, "20.417480", result)
}

func Test_keptnDynatraceDQLProvider_FetchAnalysisValue_ReceiveUnexpectedError(t *testing.T) {
	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		return nil, 0, ErrUnexpected
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, err := dqlProvider.FetchAnalysisValue(context.TODO(), "timeseries (dt.host.cpu)",
		metricsapi.Analysis{
			Status: metricsapi.AnalysisStatus{},
		},
		&metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.NotNil(t, err)
	require.Empty(t, result)
}

func Test_keptnDynatraceDQLProvider_FetchAnalysisValue_ReceiveMultipleRecords(t *testing.T) {
	mockClient := &fake.DTAPIClientMock{}

	mockClient.DoFunc = func(ctx context.Context, path string, method string, payload []byte) ([]byte, int, error) {
		return []byte(dqlPayloadMultipleRecords), 200, nil
	}

	dqlProvider := NewKeptnDynatraceDQLProvider(
		nil,
		WithDTAPIClient(mockClient),
		WithLogger(logr.New(klog.NewKlogr().GetSink())),
	)

	result, err := dqlProvider.FetchAnalysisValue(context.TODO(), "timeseries (dt.host.cpu)",
		metricsapi.Analysis{
			Status: metricsapi.AnalysisStatus{},
		},
		&metricsapi.KeptnMetricsProvider{
			Spec: metricsapi.KeptnMetricsProviderSpec{},
		},
	)

	require.Nil(t, err)
	require.Equal(t, "20.417480", result)
}

func TestExtractValuesFromRecord(t *testing.T) {
	type args struct {
		records map[string]any
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "values available",
			args: args{
				records: map[string]any{
					"avg(cpu-usage)": []any{
						25,
						13,
					},
				},
			},
			want: []string{"25.000000", "13.000000"},
		},
		{
			name: "no values",
			args: args{
				records: map[string]any{},
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := extractValuesFromRecord(tt.args.records)
			require.Equal(t, tt.want, res)
		})
	}
}

func Test_toFloatArray(t *testing.T) {
	type args struct {
		obj any
	}
	tests := []struct {
		name    string
		args    args
		wantRes []float64
		wantOk  bool
	}{
		{
			name: "array of float64 values",
			args: args{
				obj: []any{
					13.0,
					12.0,
				},
			},
			wantRes: []float64{
				13.0,
				12.0,
			},
			wantOk: true,
		},
		{
			name: "array of int values",
			args: args{
				obj: []any{
					13,
					12,
				},
			},
			wantRes: []float64{
				13.0,
				12.0,
			},
			wantOk: true,
		},
		{
			name: "array of string values",
			args: args{
				obj: []any{
					"foo",
					"bar",
				},
			},
			wantRes: nil,
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := toFloatArray(tt.args.obj)
			require.Equal(t, tt.wantRes, res)
			require.Equal(t, tt.wantOk, ok)
		})
	}
}

func Test_newMetricRequestFromMetric(t *testing.T) {
	type args struct {
		metric metricsapi.KeptnMetric
	}
	tests := []struct {
		name    string
		args    args
		want    *metricRequest
		wantErr bool
	}{
		{
			name: "metric with interval",
			args: args{
				metric: metricsapi.KeptnMetric{
					Spec: metricsapi.KeptnMetricSpec{
						Query: "queryviderci",
						Range: &metricsapi.RangeSpec{
							Interval: "5m",
						},
					},
				},
			},
			want: &metricRequest{
				query: "queryviderci",
				timeframe: &timeframe{
					from: time.Now().UTC().Add(-5 * time.Minute),
					to:   time.Now().UTC(),
				},
			},
			wantErr: false,
		},
		{
			name: "metric with invalid interval",
			args: args{
				metric: metricsapi.KeptnMetric{
					Spec: metricsapi.KeptnMetricSpec{
						Query: "queryviderci",
						Range: &metricsapi.RangeSpec{
							Interval: "5socks",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "metric without range",
			args: args{
				metric: metricsapi.KeptnMetric{
					Spec: metricsapi.KeptnMetricSpec{
						Query: "queryviderci",
					},
				},
			},
			want: &metricRequest{
				query:     "queryviderci",
				timeframe: nil,
			},
			wantErr: false,
		},
		{
			name: "metric with empty interval",
			args: args{
				metric: metricsapi.KeptnMetric{
					Spec: metricsapi.KeptnMetricSpec{
						Query: "queryviderci",
						Range: &metricsapi.RangeSpec{
							Interval: "",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newMetricRequestFromMetric(tt.args.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("newMetricRequestFromMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want == nil {
				require.Nil(t, got)
				return
			}
			require.Equal(t, tt.want.query, got.query)

			if tt.want.timeframe == nil {
				require.Nil(t, got.timeframe)
				return
			}
			require.NotNil(t, got.timeframe)
			require.WithinDuration(t, tt.want.timeframe.from, got.timeframe.from, 10*time.Second)
			require.WithinDuration(t, tt.want.timeframe.to, got.timeframe.to, 10*time.Second)
		})
	}
}

func Test_newMetricRequestFromAnalysis(t *testing.T) {
	type args struct {
		query    string
		analysis metricsapi.Analysis
	}
	tests := []struct {
		name    string
		args    args
		want    *metricRequest
		wantErr bool
	}{
		{
			name: "from analysis",
			args: args{
				query: "my-query",
				analysis: metricsapi.Analysis{
					Status: metricsapi.AnalysisStatus{
						Timeframe: metricsapi.Timeframe{
							From: metav1.Time{
								Time: time.Now().Add(-5 * time.Minute),
							},
							To: metav1.Time{
								Time: time.Now(),
							},
						},
					},
				},
			},
			want: &metricRequest{
				query: "my-query",
				timeframe: &timeframe{
					from: time.Now().Add(-5 * time.Minute),
					to:   time.Now(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newMetricRequestFromAnalysis(tt.args.query, tt.args.analysis)
			if (err != nil) != tt.wantErr {
				t.Errorf("newMetricRequestFromMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want == nil {
				require.Nil(t, got)
				return
			}
			require.Equal(t, tt.want.query, got.query)

			if tt.want.timeframe == nil {
				require.Nil(t, got.timeframe)
				return
			}
			require.NotNil(t, got.timeframe)
			require.WithinDuration(t, tt.want.timeframe.from, got.timeframe.from, 10*time.Second)
			require.WithinDuration(t, tt.want.timeframe.to, got.timeframe.to, 10*time.Second)
		})
	}
}
