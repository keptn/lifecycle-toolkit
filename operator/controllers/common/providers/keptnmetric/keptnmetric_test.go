package keptnmetric

import (
	"context"
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_keptnmetric(t *testing.T) {
	tests := []struct {
		name      string
		metric    *metricsapi.KeptnMetric
		out       string
		outraw    []byte
		wantError bool
	}{
		{
			name:      "no KeptnMetric",
			metric:    &metricsapi.KeptnMetric{},
			out:       "",
			outraw:    []byte(nil),
			wantError: true,
		},
		{
			name: "KeptnMetric without results",
			metric: &metricsapi.KeptnMetric{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "metric",
					Namespace: "default",
				},
			},
			out:       "",
			outraw:    []byte(nil),
			wantError: true,
		},
		{
			name: "KeptnMetric with results",
			metric: &metricsapi.KeptnMetric{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "metric",
					Namespace: "default",
				},
				Status: metricsapi.KeptnMetricStatus{
					Value:    "1",
					RawValue: []byte("1"),
				},
			},
			out:       "1",
			outraw:    []byte("1"),
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := metricsapi.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			client := fake.NewClientBuilder().WithObjects(tt.metric).Build()

			kmp := KeptnMetricProvider{
				Log:       ctrl.Log.WithName("testytest"),
				K8sClient: client,
			}

			obj := klcv1alpha2.Objective{
				KeptnMetricRef: klcv1alpha2.KeptnMetricRef{
					Name:      "metric",
					Namespace: "default",
				},
			}

			r, raw, e := kmp.FetchData(context.TODO(), obj)
			require.Equal(t, tt.out, r)
			require.Equal(t, tt.outraw, raw)
			if tt.wantError != (e != nil) {
				t.Errorf("want error: %t, got: %v", tt.wantError, e)
			}

		})

	}
}
