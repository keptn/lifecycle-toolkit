package providers

const DynatraceProviderType = "dynatrace"
const DynatraceDQLProviderType = "dql"
const PrometheusProviderType = "prometheus"
const DataDogProviderType = "datadog"
const KeptnPlaceholderProviderType = "placeholder"

var SupportedProviders = []string{
	DynatraceProviderType,
	DynatraceDQLProviderType,
	PrometheusProviderType,
	DataDogProviderType,
	KeptnPlaceholderProviderType,
}
