package v1alpha2

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha2/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnWorkload(t *testing.T) {
	workload := &KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workload",
			Namespace: "namespace",
		},
		Spec: KeptnWorkloadSpec{
			Version: "version",
			AppName: "app",
		},
	}

	workloadInstanceName := workload.GetWorkloadInstanceName()
	require.Equal(t, "workload-version", workloadInstanceName)

	workloadInstance := workload.GenerateWorkloadInstance("prev", map[string]string{})
	require.Equal(t, KeptnWorkloadInstance{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
			Name:        "workload-version",
			Namespace:   "namespace",
		},
		Spec: KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: KeptnWorkloadSpec{
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

	require.Equal(t, map[string]string{
		"appName":         "app",
		"workloadName":    "workload",
		"workloadVersion": "version",
	}, workload.GetEventAnnotations())
}
