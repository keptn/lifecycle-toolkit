package v1alpha1

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the src v1alpha1.KeptnApp to the hub version (v1alpha2.KeptnApp)
func (src *KeptnEvaluationProvider) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1alpha2.KeptnEvaluationProvider)

	if !ok {
		return fmt.Errorf("cannot cast KeptnEvaluationProvider to v1alpha2. Got type %T", dstRaw)
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

// ConvertFrom converts from the hub version (v1alpha2.KeptnApp) to this version (v1alpha1.KeptnApp)
func (dst *KeptnEvaluationProvider) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1alpha2.KeptnEvaluationProvider)

	if !ok {
		return fmt.Errorf("cannot cast KeptnEvaluationProvider to v1alpha1. Got type %T", srcRaw)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.TargetServer = src.Spec.TargetServer
	dst.Spec.SecretName = src.Spec.SecretKeyRef.Name

	return nil
}
