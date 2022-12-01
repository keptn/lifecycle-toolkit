package api_test

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnWorkload(t *testing.T) {
	workload := &v1alpha2.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workload",
			Namespace: "namespace",
		},
		Spec: v1alpha2.KeptnWorkloadSpec{
			Version: "version",
			AppName: "app",
		},
	}

	workloadInstanceName := workload.GetWorkloadInstanceName()
	require.Equal(t, "workload-version", workloadInstanceName)

	workloadInstance := workload.GenerateWorkloadInstance("prev", map[string]string{})
	require.Equal(t, v1alpha2.KeptnWorkloadInstance{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
			Name:        "workload-version",
			Namespace:   "namespace",
		},
		Spec: v1alpha2.KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: v1alpha2.KeptnWorkloadSpec{
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
