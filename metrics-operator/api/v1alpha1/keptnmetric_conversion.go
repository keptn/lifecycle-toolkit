package v1alpha1

import (
	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this KeptnMetric to the Hub version (v1alpha3)
func (src *KeptnMetric) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha3.KeptnMetric)

	// Spec
	dst.Spec = v1alpha3.KeptnMetricSpec{
		Provider: v1alpha3.ProviderRef{
			Name: src.Spec.Provider.Name,
		},
		Query:                src.Spec.Query,
		FetchIntervalSeconds: src.Spec.FetchIntervalSeconds,
	}

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Status
	dst.Status = v1alpha3.KeptnMetricStatus{
		Value:       src.Status.Value,
		RawValue:    src.Status.RawValue,
		LastUpdated: src.Status.LastUpdated,
	}

	return nil
}

// ConvertFrom converts from the Hub version (v1alpha3) to this version.
func (dst *KeptnMetric) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha3.KeptnMetric)

	// Spec
	dst.Spec = KeptnMetricSpec{
		Provider: ProviderRef{
			Name: src.Spec.Provider.Name,
		},
		Query:                src.Spec.Query,
		FetchIntervalSeconds: src.Spec.FetchIntervalSeconds,
	}

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Status
	dst.Status = KeptnMetricStatus{
		Value:       src.Status.Value,
		RawValue:    src.Status.RawValue,
		LastUpdated: src.Status.LastUpdated,
	}

	return nil
}
