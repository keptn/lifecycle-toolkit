package providers

const DynatraceProviderType = "dynatrace"
const DynatraceDQLProviderType = "dql"
const PrometheusProviderType = "prometheus"
const DataDogProviderType = "datadog"
const KeptnDummyProviderType = "dummy"

var SupportedProviders = []string{
	DynatraceProviderType,
	DynatraceDQLProviderType,
	PrometheusProviderType,
	DataDogProviderType,
	KeptnDummyProviderType,
}
