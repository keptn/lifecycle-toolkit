package converter

import (
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConvertSLI(t *testing.T) {
	converter := NewSLIConverter()

	// map of slis is nil
	res := converter.Convert(nil, "provider", "default")
	require.Equal(t, 0, len(res))

	// map of slis is empty
	res = converter.Convert(map[string]string{}, "provider", "default")
	require.Equal(t, 0, len(res))

	// valid input
	in := map[string]string{
		"key1": "val1",
		"key2": "val2",
	}
	out := []*metricsapi.AnalysisValueTemplate{
		{
			TypeMeta: v1.TypeMeta{
				Kind:       "AnalysisValueTemplate",
				APIVersion: "metrics.keptn.sh/v1alpha3",
			},
			ObjectMeta: v1.ObjectMeta{
				Name: "key1",
			},
			Spec: metricsapi.AnalysisValueTemplateSpec{
				Query: "val1",
				Provider: metricsapi.ObjectReference{
					Name:      "provider",
					Namespace: "default",
				},
			},
		},
		{
			TypeMeta: v1.TypeMeta{
				Kind:       "AnalysisValueTemplate",
				APIVersion: "metrics.keptn.sh/v1alpha3",
			},
			ObjectMeta: v1.ObjectMeta{
				Name: "key2",
			},
			Spec: metricsapi.AnalysisValueTemplateSpec{
				Query: "val2",
				Provider: metricsapi.ObjectReference{
					Name:      "provider",
					Namespace: "default",
				},
			},
		},
	}
	res = converter.Convert(in, "provider", "default")
	require.Equal(t, 2, len(res))
	require.Equal(t, out, res)
}
