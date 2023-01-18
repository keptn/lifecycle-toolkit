package metrics

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	metricsv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/apis/metrics/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

func TestKeptnMetricReconciler_fetchProvider(t *testing.T) {
	provider := klcv1alpha2.KeptnEvaluationProvider{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KeptnEvaluationProvider",
			APIVersion: "lifecycle.keptn.sh/v1alpha2"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myprovider",
			Namespace: "default",
		},
		Spec:   klcv1alpha2.KeptnEvaluationProviderSpec{},
		Status: klcv1alpha2.KeptnEvaluationProviderStatus{},
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

	//fetch unexisting provider

	namespacedProvider2 := types.NamespacedName{Namespace: "default", Name: "myunexistingprovider"}
	got, err = r.fetchProvider(context.TODO(), namespacedProvider2)
	require.Error(t, err)
	require.True(t, errors.IsNotFound(err))
	require.Nil(t, got)
}

func TestKeptnMetricReconciler_Reconcile(t *testing.T) {

	metric := &metricsv1alpha1.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mymetric",
			Namespace: "default",
		},
		Spec: metricsv1alpha1.KeptnMetricSpec{
			Provider:             metricsv1alpha1.ProviderRef{},
			Query:                "",
			FetchIntervalSeconds: 10,
		},
		Status: metricsv1alpha1.KeptnMetricStatus{
			Value:       "12",
			RawValue:    nil,
			LastUpdated: metav1.Time{Time: time.Now().Add(-10 * time.Second)},
		},
	}
	metric2 := &metricsv1alpha1.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mymetric2",
			Namespace: "default",
		},
		Spec: metricsv1alpha1.KeptnMetricSpec{
			Provider:             metricsv1alpha1.ProviderRef{},
			Query:                "",
			FetchIntervalSeconds: 10,
		},
		Status: metricsv1alpha1.KeptnMetricStatus{
			Value:       "12",
			RawValue:    nil,
			LastUpdated: metav1.Time{Time: time.Now()},
		},
	}

	metric3 := &metricsv1alpha1.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mymetric3",
			Namespace: "default",
		},
		Spec: metricsv1alpha1.KeptnMetricSpec{
			Provider: metricsv1alpha1.ProviderRef{
				Name: "myprov",
			},
			Query:                "",
			FetchIntervalSeconds: 10,
		},
		Status: metricsv1alpha1.KeptnMetricStatus{
			Value:       "12",
			RawValue:    nil,
			LastUpdated: metav1.Time{Time: time.Now()},
		},
	}

	provider := &klcv1alpha2.KeptnEvaluationProvider{
		ObjectMeta: metav1.ObjectMeta{Name: "myprov", Namespace: "default"},
		Spec: klcv1alpha2.KeptnEvaluationProviderSpec{
			TargetServer: "localhost",
		},
		Status: klcv1alpha2.KeptnEvaluationProviderStatus{},
	}

	client := fake.NewClient(metric, metric2, metric3, provider)

	r := &KeptnMetricReconciler{
		Client: client,
		Scheme: client.Scheme(),
		Log:    testr.New(t),
	}

	tests := []struct {
		name    string
		ctx     context.Context
		req     controllerruntime.Request
		want    controllerruntime.Result
		wantErr bool
	}{
		{
			name: "metric not found",
			ctx:  context.TODO(),
			req: controllerruntime.Request{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "myunexistingmetric"},
			},
			want:    controllerruntime.Result{},
			wantErr: false,
		},

		{
			name: "metric exists, needs to fetch, provider not found",
			ctx:  context.TODO(),
			req: controllerruntime.Request{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "mymetric2"},
			},
			want:    controllerruntime.Result{Requeue: true, RequeueAfter: 10 * time.Second},
			wantErr: false,
		},
		{
			name: "metric exists, not time to fetch",
			ctx:  context.TODO(),
			req: controllerruntime.Request{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "mymetric3"},
			},
			want:    controllerruntime.Result{Requeue: true, RequeueAfter: 10 * time.Second},
			wantErr: false,
		},

		{
			name: "metric exists, needs to fetch, provider found",
			ctx:  context.TODO(),
			req: controllerruntime.Request{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "mymetric3"},
			},
			want:    controllerruntime.Result{Requeue: true, RequeueAfter: 10 * time.Second},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := r.Reconcile(tt.ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reconcile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
