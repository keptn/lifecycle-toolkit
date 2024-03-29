//nolint:dupl
package v1alpha3

import (
	"testing"

	v1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	v1common "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/propagation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v2 "sigs.k8s.io/controller-runtime/pkg/webhook/conversion/testdata/api/v2"
)

func TestKeptnAppVersion_ConvertFrom(t *testing.T) {
	tests := []struct {
		name    string
		srcObj  *v1.KeptnAppVersion
		wantErr bool
		wantObj *KeptnAppVersion
	}{
		{
			name: "Test that conversion from v1 to v1alpha3 works",
			srcObj: &v1.KeptnAppVersion{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeptnAppVersion",
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
				Spec: v1.KeptnAppVersionSpec{
					KeptnAppSpec: v1.KeptnAppSpec{
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
					KeptnAppContextSpec: v1.KeptnAppContextSpec{
						DeploymentTaskSpec: v1.DeploymentTaskSpec{
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
					},
					AppName:         "app",
					PreviousVersion: "1.0",
					TraceId: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
				Status: v1.KeptnAppVersionStatus{
					PreDeploymentStatus:            v1common.StateFailed,
					PostDeploymentStatus:           v1common.StateFailed,
					PreDeploymentEvaluationStatus:  v1common.StateFailed,
					PostDeploymentEvaluationStatus: v1common.StateFailed,
					WorkloadOverallStatus:          v1common.StateFailed,
					WorkloadStatus: []v1.WorkloadStatus{
						{
							Workload: v1.KeptnWorkloadRef{
								Name:    "name1",
								Version: "1",
							},
							Status: v1common.StateFailed,
						},
						{
							Workload: v1.KeptnWorkloadRef{
								Name:    "name2",
								Version: "2",
							},
							Status: v1common.StateFailed,
						},
					},
					CurrentPhase: "phase",
					PreDeploymentTaskStatus: []v1.ItemStatus{
						{
							DefinitionName: "def1",
							Name:           "name1",
							Status:         v1common.StateFailed,
						},
						{
							DefinitionName: "def12",
							Name:           "name12",
							Status:         v1common.StateFailed,
						},
					},
					PostDeploymentTaskStatus: []v1.ItemStatus{
						{
							DefinitionName: "def2",
							Name:           "name2",
							Status:         v1common.StateFailed,
						},
						{
							DefinitionName: "def22",
							Name:           "name22",
							Status:         v1common.StateFailed,
						},
					},
					PreDeploymentEvaluationTaskStatus: []v1.ItemStatus{
						{
							DefinitionName: "def3",
							Name:           "name3",
							Status:         v1common.StateFailed,
						},
						{
							DefinitionName: "def32",
							Name:           "name32",
							Status:         v1common.StateFailed,
						},
					},
					PostDeploymentEvaluationTaskStatus: []v1.ItemStatus{
						{
							DefinitionName: "def4",
							Name:           "name4",
							Status:         v1common.StateFailed,
						},
						{
							DefinitionName: "def42",
							Name:           "name42",
							Status:         v1common.StateFailed,
						},
					},
					PhaseTraceIDs: v1common.PhaseTraceID{
						"key": propagation.MapCarrier{
							"key1": "value1",
							"key2": "value2",
						},
						"key22": propagation.MapCarrier{
							"key122": "value122",
							"key222": "value222",
						},
					},
					Status: v1common.StateFailed,
				},
			},
			wantErr: false,
			wantObj: &KeptnAppVersion{
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
				Spec: KeptnAppVersionSpec{
					KeptnAppSpec: KeptnAppSpec{
						Version:  "1.2.3",
						Revision: 1,
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
					AppName:         "app",
					PreviousVersion: "1.0",
					TraceId: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
				Status: KeptnAppVersionStatus{
					PreDeploymentStatus:            common.StateFailed,
					PostDeploymentStatus:           common.StateFailed,
					PreDeploymentEvaluationStatus:  common.StateFailed,
					PostDeploymentEvaluationStatus: common.StateFailed,
					WorkloadOverallStatus:          common.StateFailed,
					WorkloadStatus: []WorkloadStatus{
						{
							Workload: KeptnWorkloadRef{
								Name:    "name1",
								Version: "1",
							},
							Status: common.StateFailed,
						},
						{
							Workload: KeptnWorkloadRef{
								Name:    "name2",
								Version: "2",
							},
							Status: common.StateFailed,
						},
					},
					CurrentPhase: "phase",
					PreDeploymentTaskStatus: []ItemStatus{
						{
							DefinitionName: "def1",
							Name:           "name1",
							Status:         common.StateFailed,
						},
						{
							DefinitionName: "def12",
							Name:           "name12",
							Status:         common.StateFailed,
						},
					},
					PostDeploymentTaskStatus: []ItemStatus{
						{
							DefinitionName: "def2",
							Name:           "name2",
							Status:         common.StateFailed,
						},
						{
							DefinitionName: "def22",
							Name:           "name22",
							Status:         common.StateFailed,
						},
					},
					PreDeploymentEvaluationTaskStatus: []ItemStatus{
						{
							DefinitionName: "def3",
							Name:           "name3",
							Status:         common.StateFailed,
						},
						{
							DefinitionName: "def32",
							Name:           "name32",
							Status:         common.StateFailed,
						},
					},
					PostDeploymentEvaluationTaskStatus: []ItemStatus{
						{
							DefinitionName: "def4",
							Name:           "name4",
							Status:         common.StateFailed,
						},
						{
							DefinitionName: "def42",
							Name:           "name42",
							Status:         common.StateFailed,
						},
					},
					PhaseTraceIDs: common.PhaseTraceID{
						"key": propagation.MapCarrier{
							"key1": "value1",
							"key2": "value2",
						},
						"key22": propagation.MapCarrier{
							"key122": "value122",
							"key222": "value222",
						},
					},
					Status: common.StateFailed,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &KeptnAppVersion{
				TypeMeta:   metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{},
				Spec:       KeptnAppVersionSpec{},
				Status:     KeptnAppVersionStatus{},
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

func TestKeptnAppVersion_ConvertTo(t *testing.T) {
	tests := []struct {
		name    string
		src     *KeptnAppVersion
		wantErr bool
		wantObj *v1.KeptnAppVersion
	}{
		{
			name: "Test that conversion from v1 to v1alpha3 works",
			src: &KeptnAppVersion{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeptnAppVersion",
					APIVersion: "lifecycle.keptn.sh/v1alpha3",
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
				Spec: KeptnAppVersionSpec{
					KeptnAppSpec: KeptnAppSpec{
						Version:  "1.2.3",
						Revision: 1,
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
					AppName:         "app",
					PreviousVersion: "1.0",
					TraceId: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
				Status: KeptnAppVersionStatus{
					PreDeploymentStatus:            common.StateFailed,
					PostDeploymentStatus:           common.StateFailed,
					PreDeploymentEvaluationStatus:  common.StateFailed,
					PostDeploymentEvaluationStatus: common.StateFailed,
					WorkloadOverallStatus:          common.StateFailed,
					WorkloadStatus: []WorkloadStatus{
						{
							Workload: KeptnWorkloadRef{
								Name:    "name1",
								Version: "1",
							},
							Status: common.StateFailed,
						},
						{
							Workload: KeptnWorkloadRef{
								Name:    "name2",
								Version: "2",
							},
							Status: common.StateFailed,
						},
					},
					CurrentPhase: "phase",
					PreDeploymentTaskStatus: []ItemStatus{
						{
							DefinitionName: "def1",
							Name:           "name1",
							Status:         common.StateFailed,
						},
						{
							DefinitionName: "def12",
							Name:           "name12",
							Status:         common.StateFailed,
						},
					},
					PostDeploymentTaskStatus: []ItemStatus{
						{
							DefinitionName: "def2",
							Name:           "name2",
							Status:         common.StateFailed,
						},
						{
							DefinitionName: "def22",
							Name:           "name22",
							Status:         common.StateFailed,
						},
					},
					PreDeploymentEvaluationTaskStatus: []ItemStatus{
						{
							DefinitionName: "def3",
							Name:           "name3",
							Status:         common.StateFailed,
						},
						{
							DefinitionName: "def32",
							Name:           "name32",
							Status:         common.StateFailed,
						},
					},
					PostDeploymentEvaluationTaskStatus: []ItemStatus{
						{
							DefinitionName: "def4",
							Name:           "name4",
							Status:         common.StateFailed,
						},
						{
							DefinitionName: "def42",
							Name:           "name42",
							Status:         common.StateFailed,
						},
					},
					PhaseTraceIDs: common.PhaseTraceID{
						"key": propagation.MapCarrier{
							"key1": "value1",
							"key2": "value2",
						},
						"key22": propagation.MapCarrier{
							"key122": "value122",
							"key222": "value222",
						},
					},
					Status: common.StateFailed,
				},
			},
			wantErr: false,
			wantObj: &v1.KeptnAppVersion{
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
				Spec: v1.KeptnAppVersionSpec{
					KeptnAppSpec: v1.KeptnAppSpec{
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
					KeptnAppContextSpec: v1.KeptnAppContextSpec{
						DeploymentTaskSpec: v1.DeploymentTaskSpec{
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
					},
					AppName:         "app",
					PreviousVersion: "1.0",
					TraceId: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
				Status: v1.KeptnAppVersionStatus{
					PreDeploymentStatus:            v1common.StateFailed,
					PostDeploymentStatus:           v1common.StateFailed,
					PreDeploymentEvaluationStatus:  v1common.StateFailed,
					PostDeploymentEvaluationStatus: v1common.StateFailed,
					WorkloadOverallStatus:          v1common.StateFailed,
					WorkloadStatus: []v1.WorkloadStatus{
						{
							Workload: v1.KeptnWorkloadRef{
								Name:    "name1",
								Version: "1",
							},
							Status: v1common.StateFailed,
						},
						{
							Workload: v1.KeptnWorkloadRef{
								Name:    "name2",
								Version: "2",
							},
							Status: v1common.StateFailed,
						},
					},
					CurrentPhase: "phase",
					PreDeploymentTaskStatus: []v1.ItemStatus{
						{
							DefinitionName: "def1",
							Name:           "name1",
							Status:         v1common.StateFailed,
						},
						{
							DefinitionName: "def12",
							Name:           "name12",
							Status:         v1common.StateFailed,
						},
					},
					PostDeploymentTaskStatus: []v1.ItemStatus{
						{
							DefinitionName: "def2",
							Name:           "name2",
							Status:         v1common.StateFailed,
						},
						{
							DefinitionName: "def22",
							Name:           "name22",
							Status:         v1common.StateFailed,
						},
					},
					PreDeploymentEvaluationTaskStatus: []v1.ItemStatus{
						{
							DefinitionName: "def3",
							Name:           "name3",
							Status:         v1common.StateFailed,
						},
						{
							DefinitionName: "def32",
							Name:           "name32",
							Status:         v1common.StateFailed,
						},
					},
					PostDeploymentEvaluationTaskStatus: []v1.ItemStatus{
						{
							DefinitionName: "def4",
							Name:           "name4",
							Status:         v1common.StateFailed,
						},
						{
							DefinitionName: "def42",
							Name:           "name42",
							Status:         v1common.StateFailed,
						},
					},
					PhaseTraceIDs: v1common.PhaseTraceID{
						"key": propagation.MapCarrier{
							"key1": "value1",
							"key2": "value2",
						},
						"key22": propagation.MapCarrier{
							"key122": "value122",
							"key222": "value222",
						},
					},
					Status: v1common.StateFailed,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := v1.KeptnAppVersion{
				TypeMeta:   metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{},
				Spec:       v1.KeptnAppVersionSpec{},
				Status:     v1.KeptnAppVersionStatus{},
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

func TestKeptnAppVersion_ConvertFrom_Errorcase(t *testing.T) {
	// A random different object is used here to simulate a different API version
	testObj := v2.ExternalJob{}

	dst := &KeptnAppVersion{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       KeptnAppVersionSpec{},
		Status:     KeptnAppVersionStatus{},
	}

	if err := dst.ConvertFrom(&testObj); err == nil {
		t.Errorf("ConvertFrom() error = %v", err)
	} else {
		require.ErrorIs(t, err, common.ErrCannotCastKeptnAppVersion)
	}
}

func TestKeptnAppVersion_ConvertTo_Errorcase(t *testing.T) {
	testObj := KeptnAppVersion{}

	// A random different object is used here to simulate a different API version
	dst := v2.ExternalJob{}

	if err := testObj.ConvertTo(&dst); err == nil {
		t.Errorf("ConvertTo() error = %v", err)
	} else {
		require.ErrorIs(t, err, common.ErrCannotCastKeptnAppVersion)
	}
}
