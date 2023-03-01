package v1alpha2

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v2 "sigs.k8s.io/controller-runtime/pkg/webhook/conversion/testdata/api/v2"
)

func TestKeptnEvalProvider_ConvertFrom(t *testing.T) {
	tests := []struct {
		name    string
		srcObj  *v1alpha3.KeptnEvaluationProvider
		wantErr bool
		wantObj *KeptnEvaluationProvider
	}{
		{
			name: "Test that conversion from v1alpha3 to v1alpha2 works",
			srcObj: &v1alpha3.KeptnEvaluationProvider{
				TypeMeta: v1.TypeMeta{
					Kind:       "KeptnEvaluationProvider",
					APIVersion: "lifecycle.keptn.sh/v1alpha3",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      "some-keptn-app-name",
					Namespace: "",
					Labels: map[string]string{
						"some-label": "some-label-value",
					},
					Annotations: map[string]string{
						"some-annotation": "some-annotation-value",
					},
				},
				Spec: v1alpha3.KeptnEvaluationProviderSpec{
					TargetServer: "my-server",
					SecretKeyRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "my-secret-name",
						},
						Key: "my-secret-key",
					},
				},
				Status: v1alpha3.KeptnEvaluationProviderStatus{},
			},
			wantErr: false,
			wantObj: &KeptnEvaluationProvider{
				ObjectMeta: v1.ObjectMeta{
					Name:      "some-keptn-app-name",
					Namespace: "",
					Labels: map[string]string{
						"some-label": "some-label-value",
					},
					Annotations: map[string]string{
						"some-annotation": "some-annotation-value",
					},
				},
				Spec: KeptnEvaluationProviderSpec{
					TargetServer: "my-server",
					SecretKeyRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "my-secret-name",
						},
						Key: "my-secret-key",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &KeptnEvaluationProvider{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec:       KeptnEvaluationProviderSpec{},
				Status:     KeptnEvaluationProviderStatus{},
			}
			if err := dst.ConvertFrom(tt.srcObj); (err != nil) != tt.wantErr {
				t.Errorf("ConvertFrom() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantObj != nil {
				require.Equal(t, tt.wantObj, dst, "Object was not converted correctly")
			}
		})
	}
}

func TestKeptnEvalProvider_ConvertTo(t *testing.T) {
	tests := []struct {
		name    string
		src     *KeptnEvaluationProvider
		wantErr bool
		wantObj *v1alpha3.KeptnEvaluationProvider
	}{
		{
			name: "Test that conversion from v1alpha2 to v1alpha3 works",
			src: &KeptnEvaluationProvider{
				TypeMeta: v1.TypeMeta{
					Kind:       "KeptnEvaluationProvider",
					APIVersion: "lifecycle.keptn.sh/v1alpha2",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      "some-keptn-app-name",
					Namespace: "",
					Labels: map[string]string{
						"some-label": "some-label-value",
					},
					Annotations: map[string]string{
						"some-annotation": "some-annotation-value",
					},
				},
				Spec: KeptnEvaluationProviderSpec{
					TargetServer: "my-server",
					SecretKeyRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "my-secret-name",
						},
						Key: "apiToken",
					},
				},
				Status: KeptnEvaluationProviderStatus{},
			},
			wantErr: false,
			wantObj: &v1alpha3.KeptnEvaluationProvider{
				ObjectMeta: v1.ObjectMeta{
					Name:      "some-keptn-app-name",
					Namespace: "",
					Labels: map[string]string{
						"some-label": "some-label-value",
					},
					Annotations: map[string]string{
						"some-annotation": "some-annotation-value",
					},
				},
				Spec: v1alpha3.KeptnEvaluationProviderSpec{
					TargetServer: "my-server",
					SecretKeyRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "my-secret-name",
						},
						Key: "apiToken",
					},
				},
				Status: v1alpha3.KeptnEvaluationProviderStatus{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := v1alpha3.KeptnEvaluationProvider{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec:       v1alpha3.KeptnEvaluationProviderSpec{},
				Status:     v1alpha3.KeptnEvaluationProviderStatus{},
			}
			if err := tt.src.ConvertTo(&dst); (err != nil) != tt.wantErr {
				t.Errorf("ConvertTo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantObj != nil {
				require.Equal(t, tt.wantObj, &dst, "Object was not converted correctly")
			}
		})
	}
}

func TestKeptnEvalProvider_ConvertFrom_Errorcase(t *testing.T) {
	// A random different object is used here to simulate a different API version
	testObj := v2.ExternalJob{}

	dst := &KeptnEvaluationProvider{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{},
		Spec:       KeptnEvaluationProviderSpec{},
		Status:     KeptnEvaluationProviderStatus{},
	}

	if err := dst.ConvertFrom(&testObj); err == nil {
		t.Errorf("ConvertFrom() error = %v", err)
	} else {
		require.ErrorIs(t, err, common.ErrCannotCastKeptnEvaluationProvider)
	}
}

func TestKeptnEvalProvider_ConvertTo_Errorcase(t *testing.T) {
	testObj := KeptnEvaluationProvider{}

	// A random different object is used here to simulate a different API version
	dst := v2.ExternalJob{}

	if err := testObj.ConvertTo(&dst); err == nil {
		t.Errorf("ConvertTo() error = %v", err)
	} else {
		require.ErrorIs(t, err, common.ErrCannotCastKeptnEvaluationProvider)
	}
}
