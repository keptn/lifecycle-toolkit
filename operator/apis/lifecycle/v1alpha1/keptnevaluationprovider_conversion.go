package v1alpha1

import (
	"fmt"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha1/common"
	corev1 "k8s.io/api/core/v1"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the src v1alpha1.KeptnEvaluationProvider to the hub version (v1alpha2.KeptnEvaluationProvider)
func (src *KeptnEvaluationProvider) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1alpha2.KeptnEvaluationProvider)

	if !ok {
		return fmt.Errorf("type %T %w", dstRaw, common.CannotCastKeptnEvaluationProviderErr)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.TargetServer = src.Spec.TargetServer

	// Set sensible defaults for new fields
	dst.Spec.SecretKeyRef = corev1.SecretKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{
			Name: src.Spec.SecretName,
		},
		Key: "apiToken",
	}

	return nil
}

// ConvertFrom converts from the hub version (v1alpha2.KeptnEvaluationProvider) to this version (v1alpha1.KeptnEvaluationProvider)
func (dst *KeptnEvaluationProvider) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1alpha2.KeptnEvaluationProvider)

	if !ok {
		return fmt.Errorf("type %T %w", srcRaw, common.CannotCastKeptnEvaluationProviderErr)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.TargetServer = src.Spec.TargetServer
	dst.Spec.SecretName = src.Spec.SecretKeyRef.Name

	return nil
}
