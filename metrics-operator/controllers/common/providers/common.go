package providers

const DynatraceProviderName = "dynatrace"
const DynatraceDQLProviderName = "dql"
const PrometheusProviderName = "prometheus"
const KeptnMetricProviderName = "keptn-metric"

type MetricQuery struct {
	Query string
}
