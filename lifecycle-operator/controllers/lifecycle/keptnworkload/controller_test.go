package keptnworkload

import (
	"testing"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	klcv1alpha4 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha4"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnWorkload(t *testing.T) {
	workload := &klcv1alpha3.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workload",
			Namespace: "namespace",
		},
		Spec: klcv1alpha3.KeptnWorkloadSpec{
			Version: "version",
			AppName: "app",
		},
	}

	workloadVersion := generateWorkloadVersion("prev", map[string]string{}, workload)
	require.Equal(t, klcv1alpha4.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
			Name:        "workload-version",
			Namespace:   "namespace",
		},
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
				Version: "version",
				AppName: "app",
			},
			WorkloadName:    "workload",
			PreviousVersion: "prev",
		},
	}, workloadVersion)
}
