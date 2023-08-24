package converter

import (
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SLIConverter struct {
}

type SLI struct {
	Indicators map[string]string `yaml:"indicators"`
}

func NewSLIConverter() *SLIConverter {
	return &SLIConverter{}
}

func (c *SLIConverter) Convert(slis map[string]string, provider string, namespace string) []*metricsapi.AnalysisValueTemplate {
	result := make([]*metricsapi.AnalysisValueTemplate, len(slis))
	i := 0
	for key, query := range slis {
		template := &metricsapi.AnalysisValueTemplate{
			TypeMeta: v1.TypeMeta{
				Kind:       "AnalysisValueTemplate",
				APIVersion: "metrics.keptn.sh/v1alpha3",
			},
			ObjectMeta: v1.ObjectMeta{
				Name: key,
			},
			Spec: metricsapi.AnalysisValueTemplateSpec{
				Query: query,
				Provider: metricsapi.ObjectReference{
					Name:      provider,
					Namespace: namespace,
				},
			},
		}
		result[i] = template
		i++
	}
	return result
}
