package v1alpha3

import (
	"fmt"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the src v1alpha3.KeptnApp to the hub version (v1alpha2.KeptnApp)
func (src *KeptnApp) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1alpha2.KeptnApp)

	if !ok {
		return fmt.Errorf("type %T %w", dstRaw, common.ErrCannotCastKeptnApp)
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

	dst.Spec.Revision = src.Spec.Revision

	return nil
}

// ConvertFrom converts from the hub version (v1alpha2.KeptnApp) to this version (v1alpha3.KeptnApp)
func (dst *KeptnApp) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1alpha2.KeptnApp)

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
	dst.Spec.PreDeploymentTasks = src.Spec.PreDeploymentTasks
	dst.Spec.PostDeploymentTasks = src.Spec.PostDeploymentTasks
	dst.Spec.PreDeploymentEvaluations = src.Spec.PreDeploymentEvaluations
	dst.Spec.PostDeploymentEvaluations = src.Spec.PostDeploymentEvaluations

	dst.Status.CurrentVersion = src.Status.CurrentVersion

	dst.Spec.Revision = src.Spec.Revision

	return nil
}
