package v1alpha2

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
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

	workloadVersionName := workload.GetWorkloadVersionName()
	require.Equal(t, "workload-version", workloadVersionName)

	workloadVersion := workload.GenerateWorkloadVersion("prev", map[string]string{})
	require.Equal(t, KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
			Name:        "workload-version",
			Namespace:   "namespace",
		},
		Spec: KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: KeptnWorkloadSpec{
				Version: "version",
				AppName: "app",
			},
			WorkloadName:    "workload",
			PreviousVersion: "prev",
		},
	}, workloadVersion)

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
