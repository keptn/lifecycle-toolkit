package v1alpha3

import (
	"fmt"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the src v1beta1.KeptnAppVersion to the hub version (v1beta1.KeptnAppVersion)
//
//nolint:gocyclo
func (src *KeptnApp) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1beta1.KeptnApp)

	if !ok {
		return fmt.Errorf("type %T %w", dstRaw, "Cannot cast KeptnAppVersion")
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Version = src.Spec.Version
	for _, srcWl := range src.Spec.Workloads {
		dst.Spec.Workloads = append(dst.Spec.Workloads, v1beta1.KeptnWorkloadRef{
			Name:    srcWl.Name,
			Version: srcWl.Version,
		})
	}

	// Set sensible defaults for new fields
	dst.Spec.Revision = src.Spec.Revision

	return nil
}

// ConvertFrom converts from the hub version (v1beta1.KeptnAppVersion) to this version (v1alpha3.KeptnAppVersion)
//
//nolint:gocyclo
func (dst *KeptnApp) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1beta1.KeptnApp)

	if !ok {
		return fmt.Errorf("type %T %w", srcRaw, "cannot cast KeptnAppVersion")
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Version = src.Spec.Version
	for _, srcWl := range src.Spec.Workloads {
		dst.Spec.Workloads = append(dst.Spec.Workloads, KeptnWorkloadRef{
			Name:    srcWl.Name,
			Version: srcWl.Version,
		})
	}

	dst.Spec.Revision = src.Spec.Revision

	return nil
}
