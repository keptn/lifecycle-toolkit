package v1alpha1

import (
	"fmt"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the src v1alpha1.KeptnApp to the hub version (v1alpha2.KeptnApp)
func (src *KeptnApp) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1alpha2.KeptnApp)

	if !ok {
		return fmt.Errorf("cannot cast KeptnApp to v1alpha2. Got type %T", dstRaw)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Version = src.Spec.Version
	for _, srcWl := range src.Spec.Workloads {
		dst.Spec.Workloads = append(dst.Spec.Workloads, v1alpha2.KeptnWorkloadRef{
			Name:    srcWl.Name,
			Version: srcWl.Version,
		})
	}
	dst.Spec.PreDeploymentTasks = src.Spec.PreDeploymentTasks
	dst.Spec.PostDeploymentTasks = src.Spec.PostDeploymentTasks
	dst.Spec.PreDeploymentEvaluations = src.Spec.PreDeploymentEvaluations
	dst.Spec.PostDeploymentEvaluations = src.Spec.PostDeploymentEvaluations

	dst.Status.CurrentVersion = src.Status.CurrentVersion

	// Set sensible defaults for new fields
	dst.Spec.Revision = "1"

	return nil
}

// ConvertFrom converts from the hub version (v1alpha2.KeptnApp) to this version (v1alpha1.KeptnApp)
func (dst *KeptnApp) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1alpha2.KeptnApp)

	if !ok {
		return fmt.Errorf("cannot cast KeptnApp to v1alpha1. Got type %T", srcRaw)
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
	dst.Spec.PreDeploymentTasks = src.Spec.PreDeploymentTasks
	dst.Spec.PostDeploymentTasks = src.Spec.PostDeploymentTasks
	dst.Spec.PreDeploymentEvaluations = src.Spec.PreDeploymentEvaluations
	dst.Spec.PostDeploymentEvaluations = src.Spec.PostDeploymentEvaluations

	dst.Status.CurrentVersion = src.Status.CurrentVersion

	return nil
}
