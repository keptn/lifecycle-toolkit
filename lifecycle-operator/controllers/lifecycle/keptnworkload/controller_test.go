package keptnworkload

import (
	"testing"

	klcv1beta1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnWorkload(t *testing.T) {
	workload := &klcv1beta1.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workload",
			Namespace: "namespace",
		},
		Spec: klcv1beta1.KeptnWorkloadSpec{
			Version: "version",
			AppName: "app",
			Metadata: map[string]string{
				"foo": "bar",
			},
		},
	}

	workloadVersion := generateWorkloadVersion("prev", map[string]string{}, workload)
	require.Equal(t, klcv1beta1.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
			Name:        "workload-version",
			Namespace:   "namespace",
		},
		Spec: klcv1beta1.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1beta1.KeptnWorkloadSpec{
				Version: "version",
				AppName: "app",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			WorkloadName:    "workload",
			PreviousVersion: "prev",
		},
	}, workloadVersion)
}
