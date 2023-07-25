package pkg

import (
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SLIConverter struct {
}

func NewSLIConverter() *SLIConverter {
	return &SLIConverter{}
}

func (c *SLIConverter) Convert(slis map[string]string, provider string) []*metricsapi.AnalysisTemplate {
	result := make([]*metricsapi.AnalysisTemplate, 0, len(slis))
	for key, query := range slis {
		template := &metricsapi.AnalysisTemplate{
			ObjectMeta: v1.ObjectMeta{
				Name: key,
			},
			Spec: metricsapi.AnalysisTemplateSpec{
				Query: query,
				ProviderRef: metricsapi.ProviderReference{
					Name: provider,
				},
			},
		}
		result = append(result, template)
	}
	return result
}
