package providers

const DynatraceProviderType = "dynatrace"
const DynatraceDQLProviderType = "dql"
const PrometheusProviderType = "prometheus"
const ThanosProviderType = "thanos"
const CortexProviderType = "cortex"
const DataDogProviderType = "datadog"

var SupportedProviders = []string{
	DynatraceProviderType,
	DynatraceDQLProviderType,
	PrometheusProviderType,
	DataDogProviderType,
	CortexProviderType,
	ThanosProviderType,
}
