package v1alpha1_test

import (
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestKeptnWorkload(t *testing.T) {
	workload := &v1alpha1.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workload",
			Namespace: "namespace",
		},
		Spec: v1alpha1.KeptnWorkloadSpec{
			Version: "version",
			AppName: "app",
		},
	}

	workloadInstanceName := workload.GetWorkloadInstanceName()
	require.Equal(t, "workload-version", workloadInstanceName)

	workloadInstance := workload.GenerateWorkloadInstance("prev", map[string]string{})
	require.Equal(t, v1alpha1.KeptnWorkloadInstance{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
			Name:        "workload-version",
			Namespace:   "namespace",
		},
		Spec: v1alpha1.KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: v1alpha1.KeptnWorkloadSpec{
				Version: "version",
				AppName: "app",
			},
			WorkloadName:    "workload",
			PreviousVersion: "prev",
		},
	}, workloadInstance)

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("app"),
		common.WorkloadName.String("workload"),
		common.WorkloadVersion.String("version"),
	}, workload.GetSpanAttributes())
}
