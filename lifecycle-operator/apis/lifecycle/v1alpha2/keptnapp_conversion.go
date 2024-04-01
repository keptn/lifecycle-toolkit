package v1alpha2

import (
	"fmt"

	v1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha2/common"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the src v1alpha2.KeptnApp to the hub version (v1.KeptnApp)
func (src *KeptnApp) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1.KeptnApp)

	if !ok {
		return fmt.Errorf("type %T %w", dstRaw, common.ErrCannotCastKeptnApp)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Version = src.Spec.Version
	for _, srcWl := range src.Spec.Workloads {
		dst.Spec.Workloads = append(dst.Spec.Workloads, v1.KeptnWorkloadRef{
			Name:    srcWl.Name,
			Version: srcWl.Version,
		})
	}

	dst.Status.CurrentVersion = src.Status.CurrentVersion

	dst.Spec.Revision = src.Spec.Revision

	return nil
}

// ConvertFrom converts from the hub version (v1.KeptnApp) to this version (v1alpha2.KeptnApp)
func (dst *KeptnApp) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1.KeptnApp)

	if !ok {
		return fmt.Errorf("type %T %w", srcRaw, common.ErrCannotCastKeptnApp)
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

	dst.Status.CurrentVersion = src.Status.CurrentVersion

	dst.Spec.Revision = src.Spec.Revision

	return nil
}
