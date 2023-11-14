package providers

const DynatraceProviderType = "dynatrace"
const DynatraceDQLProviderType = "dql"
const PrometheusProviderType = "prometheus"
const DataDogProviderType = "datadog"

var SupportedProviders = []string{
	DynatraceProviderType,
	DynatraceDQLProviderType,
	PrometheusProviderType,
	DataDogProviderType,
}

func IsProviderSupported(providerName string) bool {
	for _, supportedProvider := range SupportedProviders {
		if providerName == supportedProvider {
			return true
		}
	}
	return false
}
