package v1alpha1

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v2 "sigs.k8s.io/controller-runtime/pkg/webhook/conversion/testdata/api/v2"
)

func TestKeptnApp_ConvertFrom(t *testing.T) {
	tests := []struct {
		name    string
		srcObj  *v1alpha3.KeptnApp
		wantErr bool
		wantObj *KeptnApp
	}{
		{
			name: "Test that conversion from v1alpha3 to v1alpha1 works",
			srcObj: &v1alpha3.KeptnApp{
				TypeMeta: v1.TypeMeta{
					Kind:       "KeptnApp",
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
				Spec: v1alpha3.KeptnAppSpec{
					Version:  "1.2.3",
					Revision: 1,
					Workloads: []v1alpha3.KeptnWorkloadRef{
						{
							Name:    "workload-1",
							Version: "1.2.3",
						},
						{
							Name:    "workload-2",
							Version: "4.5.6",
						},
					},
					PreDeploymentTasks: []string{
						"some-pre-deployment-task1",
					},
					PostDeploymentTasks: []string{
						"some-post-deployment-task2",
					},
					PreDeploymentEvaluations: []string{
						"some-pre-evaluation-task1",
					},
					PostDeploymentEvaluations: []string{
						"some-pre-evaluation-task2",
					},
				},
				Status: v1alpha3.KeptnAppStatus{
					CurrentVersion: "1.2.3",
				},
			},
			wantErr: false,
			wantObj: &KeptnApp{
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
					PreDeploymentTasks: []string{
						"some-pre-deployment-task1",
					},
					PostDeploymentTasks: []string{
						"some-post-deployment-task2",
					},
					PreDeploymentEvaluations: []string{
						"some-pre-evaluation-task1",
					},
					PostDeploymentEvaluations: []string{
						"some-pre-evaluation-task2",
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
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
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
		wantObj *v1alpha3.KeptnApp
	}{
		{
			name: "Test that conversion from v1alpha1 to v1alpha3 works",
			src: &KeptnApp{
				TypeMeta: v1.TypeMeta{
					Kind:       "KeptnApp",
					APIVersion: "lifecycle.keptn.sh/v1alpha1",
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
					PreDeploymentTasks: []string{
						"some-pre-deployment-task1",
					},
					PostDeploymentTasks: []string{
						"some-post-deployment-task2",
					},
					PreDeploymentEvaluations: []string{
						"some-pre-evaluation-task1",
					},
					PostDeploymentEvaluations: []string{
						"some-pre-evaluation-task2",
					},
				},
				Status: KeptnAppStatus{
					CurrentVersion: "1.2.3",
				},
			},
			wantErr: false,
			wantObj: &v1alpha3.KeptnApp{
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
				Spec: v1alpha3.KeptnAppSpec{
					Version:  "1.2.3",
					Revision: 1,
					Workloads: []v1alpha3.KeptnWorkloadRef{
						{
							Name:    "workload-1",
							Version: "1.2.3",
						},
						{
							Name:    "workload-2",
							Version: "4.5.6",
						},
					},
					PreDeploymentTasks: []string{
						"some-pre-deployment-task1",
					},
					PostDeploymentTasks: []string{
						"some-post-deployment-task2",
					},
					PreDeploymentEvaluations: []string{
						"some-pre-evaluation-task1",
					},
					PostDeploymentEvaluations: []string{
						"some-pre-evaluation-task2",
					},
				},
				Status: v1alpha3.KeptnAppStatus{
					CurrentVersion: "1.2.3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := v1alpha3.KeptnApp{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec:       v1alpha3.KeptnAppSpec{},
				Status:     v1alpha3.KeptnAppStatus{},
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
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{},
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
