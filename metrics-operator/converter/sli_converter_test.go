package converter

import (
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const sliContent = `
spec_version: "1.0"
indicators:
  throughput: "builtin:service.requestCount.total:merge(0):count?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"
  response_time_p95: "builtin:service.response.time:merge(0):percentile(95)?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"`

const expectedOutput1 = `---
apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisValueTemplate
metadata:
  creationTimestamp: null
  name: response_time_p95
spec:
  provider:
    name: dynatrace
    namespace: keptn
  query: builtin:service.response.time:merge(0):percentile(95)?scope=tag(keptn_project:{{.project}}),tag(keptn_stage:{{.stage}}),tag(keptn_service:{{.service}}),tag(keptn_deployment:{{.deployment}})
`

const expectedOutput2 = `---
apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisValueTemplate
metadata:
  creationTimestamp: null
  name: throughput
spec:
  provider:
    name: dynatrace
    namespace: keptn
  query: builtin:service.requestCount.total:merge(0):count?scope=tag(keptn_project:{{.project}}),tag(keptn_stage:{{.stage}}),tag(keptn_service:{{.service}}),tag(keptn_deployment:{{.deployment}})
`

func TestConvertMapToAnalysisValueTemplate(t *testing.T) {
	converter := NewSLIConverter()

	// map of slis is nil
	res := converter.convertMapToAnalysisValueTemplate(nil, "provider", "default")
	require.Equal(t, 0, len(res))

	// map of slis is empty
	res = converter.convertMapToAnalysisValueTemplate(map[string]string{}, "provider", "default")
	require.Equal(t, 0, len(res))

	// valid input
	in := map[string]string{
		"key1": "val1",
		"key2": "val2",
	}
	out1 := &metricsapi.AnalysisValueTemplate{
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
	}
	out2 := &metricsapi.AnalysisValueTemplate{
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
	}
	res = converter.convertMapToAnalysisValueTemplate(in, "provider", "default")
	require.Equal(t, 2, len(res))
	// need to check the result with contains, as there is no guarantee
	// of order of the AnalysisValueTemplates in resulting manifest
	require.Contains(t, res, out1)
	require.Contains(t, res, out2)
}

func TestConvertSLI(t *testing.T) {
	c := NewSLIConverter()
	// no provider nor namespace
	res, err := c.Convert([]byte(sliContent), "", "")
	require.NotNil(t, err)
	require.Equal(t, "", res)

	// invalid file content
	res, err = c.Convert([]byte("invalid"), "dynatrace", "keptn")
	require.NotNil(t, err)
	require.Equal(t, "", res)

	// happy path
	res, err = c.Convert([]byte(sliContent), "dynatrace", "keptn")
	require.Nil(t, err)
	// need to check the result with contains, as there is no guarantee
	// of order of the AnalysisValueTemplates in resulting manifest
	require.Contains(t, res, expectedOutput1)
	require.Contains(t, res, expectedOutput2)
}

func TestConvertQuary(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "no substitutions",
			in:   "hello",
			out:  "hello",
		},
		{
			name: "no substitutions - uppercase",
			in:   "HELLO",
			out:  "HELLO",
		},
		{
			name: "no substitutions - dollar with lowercase",
			in:   "$hello",
			out:  "$hello",
		},
		{
			name: "substitution - dollar with uppercase single",
			in:   "$HELLO",
			out:  "{{.hello}}",
		},
		{
			name: "substitution - dollar with uppercase and number",
			in:   "$HELLO2",
			out:  "{{.hello2}}",
		},
		{
			name: "substitution - dollar with uppercase",
			in:   "hello:$HELLO,hi:$HI",
			out:  "hello:{{.hello}},hi:{{.hi}}",
		},
		{
			name: "substitution - real query",
			in:   "builtin:service.response.time:merge(0):percentile(95)?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)",
			out:  "builtin:service.response.time:merge(0):percentile(95)?scope=tag(keptn_project:{{.project}}),tag(keptn_stage:{{.stage}}),tag(keptn_service:{{.service}}),tag(keptn_deployment:{{.deployment}})",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.out, convertQuery(tt.in))
		})

	}
}
