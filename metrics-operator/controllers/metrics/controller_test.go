package metrics

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

func TestKeptnMetricReconciler_fetchProvider(t *testing.T) {
	provider := metricsapi.KeptnMetricsProvider{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KeptnMetricsProvider",
			APIVersion: "metrics.keptn.sh/v1alpha3"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myprovider",
			Namespace: "default",
		},
		Spec: metricsapi.KeptnMetricsProviderSpec{
			Type: "prometheus",
		},
	}

	client := fake.NewClient(&provider)
	r := &KeptnMetricReconciler{
		Client: client,
		Scheme: client.Scheme(),
		Log:    testr.New(t),
	}

	// fetch existing provider based on source
	namespacedProvider := types.NamespacedName{Namespace: "default", Name: "myprovider"}
	got, err := r.fetchProvider(context.TODO(), namespacedProvider)
	require.Nil(t, err)
	require.Equal(t, provider, *got)

	// fetch unexisting provider

	namespacedProvider2 := types.NamespacedName{Namespace: "default", Name: "myunexistingprovider"}
	got, err = r.fetchProvider(context.TODO(), namespacedProvider2)
	require.Error(t, err)
	require.True(t, errors.IsNotFound(err))
	require.Nil(t, got)
}

func TestKeptnMetricReconciler_Reconcile(t *testing.T) {

	metric := &metricsapi.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mymetric",
			Namespace: "default",
		},
		Spec: metricsapi.KeptnMetricSpec{
			Provider:             metricsapi.ProviderRef{},
			Query:                "",
			FetchIntervalSeconds: 1,
		},
		Status: metricsapi.KeptnMetricStatus{
			Value:       "12",
			RawValue:    nil,
			LastUpdated: metav1.Time{Time: time.Now().Add(-1 * time.Minute)},
		},
	}
	metric2 := &metricsapi.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mymetric2",
			Namespace: "default",
		},
		Spec: metricsapi.KeptnMetricSpec{
			Provider:             metricsapi.ProviderRef{},
			Query:                "",
			FetchIntervalSeconds: 1,
		},
		Status: metricsapi.KeptnMetricStatus{
			Value:       "12",
			RawValue:    nil,
			LastUpdated: metav1.Time{Time: time.Now().Add(-1 * time.Minute)},
		},
	}

	metric3 := &metricsapi.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mymetric3",
			Namespace: "default",
		},
		Spec: metricsapi.KeptnMetricSpec{
			Provider: metricsapi.ProviderRef{
				Name: "myprov",
			},
			Query:                "",
			FetchIntervalSeconds: 10,
		},
		Status: metricsapi.KeptnMetricStatus{
			Value:       "12",
			RawValue:    nil,
			LastUpdated: metav1.Time{Time: time.Now().Add(-1 * time.Minute)},
		},
	}

	metric4 := &metricsapi.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mymetric4",
			Namespace: "default",
		},
		Spec: metricsapi.KeptnMetricSpec{
			Provider: metricsapi.ProviderRef{
				Name: "provider-name",
			},
			Query:                "",
			FetchIntervalSeconds: 10,
		},
		Status: metricsapi.KeptnMetricStatus{
			Value:       "12",
			RawValue:    nil,
			LastUpdated: metav1.Time{Time: time.Now().Add(-1 * time.Minute)},
		},
	}

	metric5 := &metricsapi.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mymetric5",
			Namespace: "default",
		},
		Spec: metricsapi.KeptnMetricSpec{
			Provider: metricsapi.ProviderRef{
				Name: "prometheus",
			},
			Query:                "",
			FetchIntervalSeconds: 10,
		},
		Status: metricsapi.KeptnMetricStatus{
			Value:       "12",
			RawValue:    nil,
			LastUpdated: metav1.Time{Time: time.Now().Add(-1 * time.Minute)},
		},
	}

	metric6 := &metricsapi.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mymetric6",
			Namespace: "default",
		},
		Spec: metricsapi.KeptnMetricSpec{
			Provider: metricsapi.ProviderRef{
				Name: "prometheus",
			},
			Query:                "",
			FetchIntervalSeconds: 10,
			Range: &metricsapi.RangeSpec{
				Interval: "5m",
			},
		},
		Status: metricsapi.KeptnMetricStatus{
			Value:       "12",
			RawValue:    nil,
			LastUpdated: metav1.Time{Time: time.Now().Add(-1 * time.Minute)},
		},
	}

	metric7 := &metricsapi.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mymetric7",
			Namespace: "default",
		},
		Spec: metricsapi.KeptnMetricSpec{
			Provider: metricsapi.ProviderRef{
				Name: "prometheus",
			},
			Query:                "",
			FetchIntervalSeconds: 10,
			Range: &metricsapi.RangeSpec{
				Interval:    "5m",
				Step:        "1m",
				Aggregation: "max",
			},
		},
		Status: metricsapi.KeptnMetricStatus{
			Value:       "12",
			RawValue:    nil,
			LastUpdated: metav1.Time{Time: time.Now().Add(-1 * time.Minute)},
		},
	}

	unsupportedProvider := &metricsapi.KeptnMetricsProvider{
		ObjectMeta: metav1.ObjectMeta{Name: "myprov", Namespace: "default"},
		Spec: metricsapi.KeptnMetricsProviderSpec{
			Type: "unsupported-type",
		},
	}

	supportedProvider := &metricsapi.KeptnMetricsProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "provider-name",
			Namespace: "default",
		},
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: "http://keptn.sh",
			Type:         "prometheus",
		},
	}

	oldSupportedProvider := &metricsapi.KeptnMetricsProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prometheus",
			Namespace: "default",
		},
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: "http://keptn.sh",
		},
	}

	client := fake.NewClient(metric, metric2, metric3, metric4, metric5, metric6, metric7, unsupportedProvider, supportedProvider, oldSupportedProvider)

	r := &KeptnMetricReconciler{
		Client: client,
		Scheme: client.Scheme(),
		Log:    testr.New(t),
	}

	tests := []struct {
		name       string
		ctx        context.Context
		req        controllerruntime.Request
		want       controllerruntime.Result
		wantMetric *metricsapi.KeptnMetric
		wantErr    error
	}{
		{
			name: "metric not found, ignoring",
			ctx:  context.TODO(),
			req: controllerruntime.Request{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "myunexistingmetric"},
			},
			want:       controllerruntime.Result{},
			wantMetric: nil,
		},

		{
			name: "metric exists, not time to fetch",
			ctx:  context.TODO(),
			req: controllerruntime.Request{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "mymetric"},
			},
			want:       controllerruntime.Result{Requeue: true, RequeueAfter: 10 * time.Second},
			wantMetric: nil,
		},

		{
			name: "metric exists, needs to fetch, provider not found ignoring",
			ctx:  context.TODO(),
			req: controllerruntime.Request{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "mymetric2"},
			},
			want:       controllerruntime.Result{Requeue: true, RequeueAfter: 10 * time.Second},
			wantMetric: nil,
		},

		{
			name: "metric exists, needs to fetch, provider unsupported",
			ctx:  context.TODO(),
			req: controllerruntime.Request{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "mymetric3"},
			},
			want:       controllerruntime.Result{Requeue: false, RequeueAfter: 0},
			wantErr:    fmt.Errorf("provider unsupported-type not supported"),
			wantMetric: nil,
		},

		{
			name: "metric exists, needs to fetch, prometheus supported, bad query",
			ctx:  context.TODO(),
			req: controllerruntime.Request{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "mymetric4"},
			},
			want:    controllerruntime.Result{Requeue: false, RequeueAfter: 0},
			wantErr: fmt.Errorf("client_error: client error: 404"),
			wantMetric: &metricsapi.KeptnMetric{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mymetric4",
					Namespace: "default",
				},
				Status: metricsapi.KeptnMetricStatus{
					ErrMsg:   "client_error: client error: 404",
					Value:    "",
					RawValue: []byte{},
				},
			},
		},

		{
			name: "metric exists, needs to fetch, using old provider API, bad query",
			ctx:  context.TODO(),
			req: controllerruntime.Request{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "mymetric5"},
			},
			want:    controllerruntime.Result{Requeue: false, RequeueAfter: 0},
			wantErr: fmt.Errorf("client_error: client error: 404"),
			wantMetric: &metricsapi.KeptnMetric{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mymetric5",
					Namespace: "default",
				},
				Status: metricsapi.KeptnMetricStatus{
					ErrMsg:   "client_error: client error: 404",
					Value:    "",
					RawValue: []byte{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			got, err := r.Reconcile(tt.ctx, tt.req)
			if tt.wantErr != nil {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				require.Nil(t, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reconcile() got = %v, want %v", got, tt.want)
			}

			if tt.wantMetric != nil {
				metric := &metricsapi.KeptnMetric{}
				err := client.Get(context.TODO(), types.NamespacedName{Namespace: tt.wantMetric.Namespace, Name: tt.wantMetric.Name}, metric)
				require.Nil(t, err)
				require.Equal(t, tt.wantMetric.Status.ErrMsg, metric.Status.ErrMsg)
				require.Equal(t, tt.wantMetric.Status.Value, metric.Status.Value)
				require.Equal(t, tt.wantMetric.Status.RawValue, metric.Status.RawValue)
			}
		})
	}
}

func Test_cupSize(t *testing.T) {
	myVeryBigSlice := make([]byte, MB+1)
	mySmallSlice := []byte("I am small")
	myAtLimitSlice := make([]byte, MB)

	res1 := cupSize(myVeryBigSlice)
	res2 := cupSize(mySmallSlice)
	res3 := cupSize(myAtLimitSlice)

	require.Equal(t, len(res1), MB)
	require.Equal(t, len(res2), len(mySmallSlice))
	require.Equal(t, len(res3), MB)

}

func Test_AggregateValues(t *testing.T) {
	tests := []struct {
		name        string
		aggFunc     string
		stringSlice []string
		want        string
	}{
		{
			name:        "test-max-for-even-length",
			aggFunc:     "max",
			stringSlice: []string{"1", "2", "3", "4"},
			want:        "4",
		},
		{
			name:        "test-max-for-odd-length",
			aggFunc:     "max",
			stringSlice: []string{"1", "2", "3", "4","5"},
			want:        "5",
		},
		{
			name:        "test-min-for-even-length",
			aggFunc:     "min",
			stringSlice: []string{"1", "2", "3", "4"},
			want:        "1",
		},
		{
			name:        "test-min-for-odd-length",
			aggFunc:     "min",
			stringSlice: []string{"1", "2", "3", "4","5"},
			want:        "1",
		},
		{
			name:        "test-median-for-even-length",
			aggFunc:     "median",
			stringSlice: []string{"1", "2", "3", "4"},
			want:        "2.5",
		},
		{
			name:        "test-median-for-odd-length",
			aggFunc:     "median",
			stringSlice: []string{"1", "2", "3", "4", "5"},
			want:        "3",
		},
		{
			name:        "test-avg-for-even-length",
			aggFunc:     "avg",
			stringSlice: []string{"1", "2", "3", "4"},
			want:        "2.5",
		},
		{
			name:        "test-avg-for-odd-length",
			aggFunc:     "avg",
			stringSlice: []string{"1", "2", "3", "4","5"},
			want:        "3",
		},
		{
			name:        "test-p90-for-even-length",
			aggFunc:     "p90",
			stringSlice: []string{"1", "2", "3", "4"},
			want:        "4",
		},
		{
			name:        "test-p90-for-odd-length",
			aggFunc:     "p90",
			stringSlice: []string{"1", "2", "3", "4","5"},
			want:        "5",
		},
		{
			name:        "test-p95-for-even-length",
			aggFunc:     "p95",
			stringSlice: []string{"1", "2", "3", "4"},
			want:        "4",
		},
		{
			name:        "test-p95-for-odd-length",
			aggFunc:     "p95",
			stringSlice: []string{"1", "2", "3", "4","5"},
			want:        "5",
		},
		{
			name:        "test-p99-for-even-length",
			aggFunc:     "p99",
			stringSlice: []string{"1", "2", "3", "4"},
			want:        "4",
		},
		{
			name:        "test-p99-for-odd-length",
			aggFunc:     "p99",
			stringSlice: []string{"1", "2", "3", "4","5"},
			want:        "5",
		},
		{
			name:        "test-max-empty-string",
			aggFunc:     "max",
			stringSlice: []string(nil),
			want:        "0",
		},
		{
			name:        "test-min-empty-string",
			aggFunc:     "min",
			stringSlice: []string(nil),
			want:        "0",
		},
		{
			name:        "test-median-empty-string",
			aggFunc:     "median",
			stringSlice: []string(nil),
			want:        "0",
		},
		{
			name:        "test-avg-empty-string",
			aggFunc:     "avg",
			stringSlice: []string(nil),
			want:        "0",
		},
		{
			name:        "test-p90-empty-string",
			aggFunc:     "p90",
			stringSlice: []string(nil),
			want:        "0",
		},
		{
			name:        "test-p95-empty-string",
			aggFunc:     "p95",
			stringSlice: []string(nil),
			want:        "0",
		},
		{
			name:        "test-p99-empty-string",
			aggFunc:     "p99",
			stringSlice: []string(nil),
			want:        "0",
		},
		{
			name:        "wrong-aggFunc",
			aggFunc:     "p50",
			stringSlice: []string{"1", "2", "3", "4"},
			want:        "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			res, err := aggregateValues(tt.stringSlice, tt.aggFunc)
			require.Equal(t, tt.want, res)
			require.Nil(t, err)
		})
	}
}
