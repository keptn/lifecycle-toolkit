package v1alpha1

import (
	"testing"

	v1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha1/common"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v2 "sigs.k8s.io/controller-runtime/pkg/webhook/conversion/testdata/api/v2"
)

func TestKeptnApp_ConvertFrom(t *testing.T) {
	tests := []struct {
		name    string
		srcObj  *v1.KeptnApp
		wantErr bool
		wantObj *KeptnApp
	}{
		{
			name: "Test that conversion from v1 to v1alpha1 works",
			srcObj: &v1.KeptnApp{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeptnApp",
					APIVersion: "lifecycle.keptn.sh/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "some-keptn-app-name",
					Namespace: "",
					Labels: map[string]string{
						"some-label": "some-label-value",
					},
					Annotations: map[string]string{
						"some-annotation": "some-annotation-value",
					},
				},
				Spec: v1.KeptnAppSpec{
					Version:  "1.2.3",
					Revision: 1,
					Workloads: []v1.KeptnWorkloadRef{
						{
							Name:    "workload-1",
							Version: "1.2.3",
						},
						{
							Name:    "workload-2",
							Version: "4.5.6",
						},
					},
				},
				Status: v1.KeptnAppStatus{
					CurrentVersion: "1.2.3",
				},
			},
			wantErr: false,
			wantObj: &KeptnApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "some-keptn-app-name",
					Namespace: "",
					Labels: map[string]string{
						"some-label": "some-label-value",
					},
					Annotations: map[string]string{
						"some-annotation": "some-annotation-value",
					},
				},
				Spec: KeptnAppSpec{
					Version: "1.2.3",
					Workloads: []KeptnWorkloadRef{
						{
							Name:    "workload-1",
							Version: "1.2.3",
						},
						{
							Name:    "workload-2",
							Version: "4.5.6",
						},
					},
				},
				Status: KeptnAppStatus{
					CurrentVersion: "1.2.3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &KeptnApp{
				TypeMeta:   metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{},
				Spec:       KeptnAppSpec{},
				Status:     KeptnAppStatus{},
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

func TestKeptnApp_ConvertTo(t *testing.T) {
	tests := []struct {
		name    string
		src     *KeptnApp
		wantErr bool
		wantObj *v1.KeptnApp
	}{
		{
			name: "Test that conversion from v1alpha1 to v1 works",
			src: &KeptnApp{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeptnApp",
					APIVersion: "lifecycle.keptn.sh/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "some-keptn-app-name",
					Namespace: "",
					Labels: map[string]string{
						"some-label": "some-label-value",
					},
					Annotations: map[string]string{
						"some-annotation": "some-annotation-value",
					},
				},
				Spec: KeptnAppSpec{
					Version: "1.2.3",
					Workloads: []KeptnWorkloadRef{
						{
							Name:    "workload-1",
							Version: "1.2.3",
						},
						{
							Name:    "workload-2",
							Version: "4.5.6",
						},
					},
				},
				Status: KeptnAppStatus{
					CurrentVersion: "1.2.3",
				},
			},
			wantErr: false,
			wantObj: &v1.KeptnApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "some-keptn-app-name",
					Namespace: "",
					Labels: map[string]string{
						"some-label": "some-label-value",
					},
					Annotations: map[string]string{
						"some-annotation": "some-annotation-value",
					},
				},
				Spec: v1.KeptnAppSpec{
					Version:  "1.2.3",
					Revision: 1,
					Workloads: []v1.KeptnWorkloadRef{
						{
							Name:    "workload-1",
							Version: "1.2.3",
						},
						{
							Name:    "workload-2",
							Version: "4.5.6",
						},
					},
				},
				Status: v1.KeptnAppStatus{
					CurrentVersion: "1.2.3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := v1.KeptnApp{
				TypeMeta:   metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{},
				Spec:       v1.KeptnAppSpec{},
				Status:     v1.KeptnAppStatus{},
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

func TestKeptnApp_ConvertFrom_Errorcase(t *testing.T) {
	// A random different object is used here to simulate a different API version
	testObj := v2.ExternalJob{}

	dst := &KeptnApp{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       KeptnAppSpec{},
		Status:     KeptnAppStatus{},
	}

	if err := dst.ConvertFrom(&testObj); err == nil {
		t.Errorf("ConvertFrom() error = %v", err)
	} else {
		require.ErrorIs(t, err, common.ErrCannotCastKeptnApp)
	}
}

func TestKeptnApp_ConvertTo_Errorcase(t *testing.T) {
	testObj := KeptnApp{}

	// A random different object is used here to simulate a different API version
	dst := v2.ExternalJob{}

	if err := testObj.ConvertTo(&dst); err == nil {
		t.Errorf("ConvertTo() error = %v", err)
	} else {
		require.ErrorIs(t, err, common.ErrCannotCastKeptnApp)
	}
}
