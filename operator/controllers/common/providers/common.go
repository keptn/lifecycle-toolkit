package providers

import (
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const DynatraceProviderName = "dynatrace"
const PrometheusProviderName = "prometheus"
const KeptnMetricProviderName = "keptn-metric"
const KLTNamespace = "keptn-lifecycle-toolkit-system"

var MetricDefaultProvider = &klcv1alpha2.KeptnEvaluationProvider{
	ObjectMeta: metav1.ObjectMeta{
		Name:      KeptnMetricProviderName,
		Namespace: KLTNamespace,
	},
}
