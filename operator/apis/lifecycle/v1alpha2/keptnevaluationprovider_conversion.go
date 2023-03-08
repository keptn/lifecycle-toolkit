package v1alpha2

import (
	"fmt"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the src v1alpha2.KeptnEvaluationProvider to the hub version (v1alpha3.KeptnEvaluationProvider)
func (src *KeptnEvaluationProvider) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1alpha3.KeptnEvaluationProvider)

	if !ok {
		return fmt.Errorf("type %T %w", dstRaw, common.ErrCannotCastKeptnEvaluationProvider)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.TargetServer = src.Spec.TargetServer

	// Set sensible defaults for new fields
	dst.Spec.SecretKeyRef = corev1.SecretKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{
			Name: src.Spec.SecretKeyRef.Name,
		},
		Key: src.Spec.SecretKeyRef.Key,
	}

	return nil
}

// ConvertFrom converts from the hub version (v1alpha3.KeptnEvaluationProvider) to this version (v1alpha2.KeptnEvaluationProvider)
func (dst *KeptnEvaluationProvider) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1alpha3.KeptnEvaluationProvider)

	if !ok {
		return fmt.Errorf("type %T %w", srcRaw, common.ErrCannotCastKeptnEvaluationProvider)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.TargetServer = src.Spec.TargetServer
	// Set sensible defaults for new fields
	dst.Spec.SecretKeyRef = corev1.SecretKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{
			Name: src.Spec.SecretKeyRef.Name,
		},
		Key: src.Spec.SecretKeyRef.Key,
	}

	return nil
}
