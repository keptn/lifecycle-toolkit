package api_test

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnWorkload(t *testing.T) {
	app := &v1alpha1.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workload",
			Namespace: "namespace",
		},
		Spec: v1alpha1.KeptnWorkloadSpec{
			Version: "version",
		},
	}

	workloadInstanceName := app.GetWorkloadInstanceName()
	require.Equal(t, "workload-version", workloadInstanceName)

	workloadInstance := app.GenerateWorkloadInstance("prev", map[string]string{})
	require.Equal(t, v1alpha1.KeptnWorkloadInstance{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
			Name:        "workload-version",
			Namespace:   "namespace",
		},
		Spec: v1alpha1.KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: v1alpha1.KeptnWorkloadSpec{
				Version: "version",
			},
			WorkloadName:    "workload",
			PreviousVersion: "prev",
		},
	}, workloadInstance)
}
