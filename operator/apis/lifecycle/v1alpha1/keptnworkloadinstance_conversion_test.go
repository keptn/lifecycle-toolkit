//nolint:dupl
package v1alpha1

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha1/common"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	v1alpha3common "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/propagation"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	v2 "sigs.k8s.io/controller-runtime/pkg/webhook/conversion/testdata/api/v2"
)

func TestKeptnWorkloadInstance_ConvertFrom(t *testing.T) {
	tests := []struct {
		name    string
		srcObj  *v1alpha3.KeptnWorkloadInstance
		wantErr bool
		wantObj *KeptnWorkloadInstance
	}{
		{
			name: "Test that conversion from v1alpha3 to v1alpha1 works",
			srcObj: &v1alpha3.KeptnWorkloadInstance{
				TypeMeta: v1.TypeMeta{
					Kind:       "KeptnWorkloadInstance",
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
				Spec: v1alpha3.KeptnWorkloadInstanceSpec{
					KeptnWorkloadSpec: v1alpha3.KeptnWorkloadSpec{
						Version: "1.2.3",
						ResourceReference: v1alpha3.ResourceReference{
							UID:  types.UID("1"),
							Kind: "Pod",
							Name: "pod",
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
						AppName: "app",
					},
					WorkloadName:    "workload",
					PreviousVersion: "1.0",
					TraceId: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
				Status: v1alpha3.KeptnWorkloadInstanceStatus{
					PreDeploymentStatus:            v1alpha3common.StateFailed,
					PostDeploymentStatus:           v1alpha3common.StateFailed,
					PreDeploymentEvaluationStatus:  v1alpha3common.StateFailed,
					PostDeploymentEvaluationStatus: v1alpha3common.StateFailed,
					DeploymentStatus:               v1alpha3common.StateFailed,
					CurrentPhase:                   "phase",
					PreDeploymentTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "def1",
							Name:           "name1",
							Status:         v1alpha3common.StateFailed,
						},
						{
							DefinitionName: "def12",
							Name:           "name12",
							Status:         v1alpha3common.StateFailed,
						},
					},
					PostDeploymentTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "def2",
							Name:           "name2",
							Status:         v1alpha3common.StateFailed,
						},
						{
							DefinitionName: "def22",
							Name:           "name22",
							Status:         v1alpha3common.StateFailed,
						},
					},
					PreDeploymentEvaluationTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "def3",
							Name:           "name3",
							Status:         v1alpha3common.StateFailed,
						},
						{
							DefinitionName: "def32",
							Name:           "name32",
							Status:         v1alpha3common.StateFailed,
						},
					},
					PostDeploymentEvaluationTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "def4",
							Name:           "name4",
							Status:         v1alpha3common.StateFailed,
						},
						{
							DefinitionName: "def42",
							Name:           "name42",
							Status:         v1alpha3common.StateFailed,
						},
					},
					PhaseTraceIDs: v1alpha3common.PhaseTraceID{
						"key": propagation.MapCarrier{
							"key1": "value1",
							"key2": "value2",
						},
						"key22": propagation.MapCarrier{
							"key122": "value122",
							"key222": "value222",
						},
					},
					Status: v1alpha3common.StateFailed,
				},
			},
			wantErr: false,
			wantObj: &KeptnWorkloadInstance{
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
				Spec: KeptnWorkloadInstanceSpec{
					KeptnWorkloadSpec: KeptnWorkloadSpec{
						Version: "1.2.3",
						ResourceReference: ResourceReference{
							UID:  types.UID("1"),
							Kind: "Pod",
							Name: "pod",
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
						AppName: "app",
					},
					WorkloadName:    "workload",
					PreviousVersion: "1.0",
					TraceId: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
				Status: KeptnWorkloadInstanceStatus{
					PreDeploymentStatus:            common.StateFailed,
					PostDeploymentStatus:           common.StateFailed,
					PreDeploymentEvaluationStatus:  common.StateFailed,
					PostDeploymentEvaluationStatus: common.StateFailed,
					DeploymentStatus:               common.StateFailed,
					CurrentPhase:                   "phase",
					PreDeploymentTaskStatus: []TaskStatus{
						{
							TaskDefinitionName: "def1",
							TaskName:           "name1",
							Status:             common.StateFailed,
						},
						{
							TaskDefinitionName: "def12",
							TaskName:           "name12",
							Status:             common.StateFailed,
						},
					},
					PostDeploymentTaskStatus: []TaskStatus{
						{
							TaskDefinitionName: "def2",
							TaskName:           "name2",
							Status:             common.StateFailed,
						},
						{
							TaskDefinitionName: "def22",
							TaskName:           "name22",
							Status:             common.StateFailed,
						},
					},
					PreDeploymentEvaluationTaskStatus: []EvaluationStatus{
						{
							EvaluationDefinitionName: "def3",
							EvaluationName:           "name3",
							Status:                   common.StateFailed,
						},
						{
							EvaluationDefinitionName: "def32",
							EvaluationName:           "name32",
							Status:                   common.StateFailed,
						},
					},
					PostDeploymentEvaluationTaskStatus: []EvaluationStatus{
						{
							EvaluationDefinitionName: "def4",
							EvaluationName:           "name4",
							Status:                   common.StateFailed,
						},
						{
							EvaluationDefinitionName: "def42",
							EvaluationName:           "name42",
							Status:                   common.StateFailed,
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
			dst := &KeptnWorkloadInstance{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec:       KeptnWorkloadInstanceSpec{},
				Status:     KeptnWorkloadInstanceStatus{},
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

func TestKeptnWorkloadInstance_ConvertTo(t *testing.T) {
	tests := []struct {
		name    string
		src     *KeptnWorkloadInstance
		wantErr bool
		wantObj *v1alpha3.KeptnWorkloadInstance
	}{
		{
			name: "Test that conversion from v1alpha1 to v1alpha3 works",
			src: &KeptnWorkloadInstance{
				TypeMeta: v1.TypeMeta{
					Kind:       "KeptnWorkloadInstance",
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
				Spec: KeptnWorkloadInstanceSpec{
					KeptnWorkloadSpec: KeptnWorkloadSpec{
						Version: "1.2.3",
						ResourceReference: ResourceReference{
							UID:  types.UID("1"),
							Kind: "Pod",
							Name: "pod",
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
						AppName: "app",
					},
					WorkloadName:    "workload",
					PreviousVersion: "1.0",
					TraceId: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
				Status: KeptnWorkloadInstanceStatus{
					PreDeploymentStatus:            common.StateFailed,
					PostDeploymentStatus:           common.StateFailed,
					PreDeploymentEvaluationStatus:  common.StateFailed,
					PostDeploymentEvaluationStatus: common.StateFailed,
					DeploymentStatus:               common.StateFailed,
					CurrentPhase:                   "phase",
					PreDeploymentTaskStatus: []TaskStatus{
						{
							TaskDefinitionName: "def1",
							TaskName:           "name1",
							Status:             common.StateFailed,
						},
						{
							TaskDefinitionName: "def12",
							TaskName:           "name12",
							Status:             common.StateFailed,
						},
					},
					PostDeploymentTaskStatus: []TaskStatus{
						{
							TaskDefinitionName: "def2",
							TaskName:           "name2",
							Status:             common.StateFailed,
						},
						{
							TaskDefinitionName: "def22",
							TaskName:           "name22",
							Status:             common.StateFailed,
						},
					},
					PreDeploymentEvaluationTaskStatus: []EvaluationStatus{
						{
							EvaluationDefinitionName: "def3",
							EvaluationName:           "name3",
							Status:                   common.StateFailed,
						},
						{
							EvaluationDefinitionName: "def32",
							EvaluationName:           "name32",
							Status:                   common.StateFailed,
						},
					},
					PostDeploymentEvaluationTaskStatus: []EvaluationStatus{
						{
							EvaluationDefinitionName: "def4",
							EvaluationName:           "name4",
							Status:                   common.StateFailed,
						},
						{
							EvaluationDefinitionName: "def42",
							EvaluationName:           "name42",
							Status:                   common.StateFailed,
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
			wantObj: &v1alpha3.KeptnWorkloadInstance{
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
				Spec: v1alpha3.KeptnWorkloadInstanceSpec{
					KeptnWorkloadSpec: v1alpha3.KeptnWorkloadSpec{
						Version: "1.2.3",
						ResourceReference: v1alpha3.ResourceReference{
							UID:  types.UID("1"),
							Kind: "Pod",
							Name: "pod",
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
						AppName: "app",
					},
					WorkloadName:    "workload",
					PreviousVersion: "1.0",
					TraceId: map[string]string{
						"key1": "value1",
						"key2": "value2",
					},
				},
				Status: v1alpha3.KeptnWorkloadInstanceStatus{
					PreDeploymentStatus:            v1alpha3common.StateFailed,
					PostDeploymentStatus:           v1alpha3common.StateFailed,
					PreDeploymentEvaluationStatus:  v1alpha3common.StateFailed,
					PostDeploymentEvaluationStatus: v1alpha3common.StateFailed,
					DeploymentStatus:               v1alpha3common.StateFailed,
					CurrentPhase:                   "phase",
					PreDeploymentTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "def1",
							Name:           "name1",
							Status:         v1alpha3common.StateFailed,
						},
						{
							DefinitionName: "def12",
							Name:           "name12",
							Status:         v1alpha3common.StateFailed,
						},
					},
					PostDeploymentTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "def2",
							Name:           "name2",
							Status:         v1alpha3common.StateFailed,
						},
						{
							DefinitionName: "def22",
							Name:           "name22",
							Status:         v1alpha3common.StateFailed,
						},
					},
					PreDeploymentEvaluationTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "def3",
							Name:           "name3",
							Status:         v1alpha3common.StateFailed,
						},
						{
							DefinitionName: "def32",
							Name:           "name32",
							Status:         v1alpha3common.StateFailed,
						},
					},
					PostDeploymentEvaluationTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "def4",
							Name:           "name4",
							Status:         v1alpha3common.StateFailed,
						},
						{
							DefinitionName: "def42",
							Name:           "name42",
							Status:         v1alpha3common.StateFailed,
						},
					},
					PhaseTraceIDs: v1alpha3common.PhaseTraceID{
						"key": propagation.MapCarrier{
							"key1": "value1",
							"key2": "value2",
						},
						"key22": propagation.MapCarrier{
							"key122": "value122",
							"key222": "value222",
						},
					},
					Status: v1alpha3common.StateFailed,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := v1alpha3.KeptnWorkloadInstance{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec:       v1alpha3.KeptnWorkloadInstanceSpec{},
				Status:     v1alpha3.KeptnWorkloadInstanceStatus{},
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

func TestKeptnWorkloadInstance_ConvertFrom_Errorcase(t *testing.T) {
	// A random different object is used here to simulate a different API version
	testObj := v2.ExternalJob{}

	dst := &KeptnWorkloadInstance{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{},
		Spec:       KeptnWorkloadInstanceSpec{},
		Status:     KeptnWorkloadInstanceStatus{},
	}

	if err := dst.ConvertFrom(&testObj); err == nil {
		t.Errorf("ConvertFrom() error = %v", err)
	} else {
		require.ErrorIs(t, err, common.ErrCannotCastKeptnWorkloadInstance)
	}
}

func TestKeptnWorkloadInstance_ConvertTo_Errorcase(t *testing.T) {
	testObj := KeptnWorkloadInstance{}

	// A random different object is used here to simulate a different API version
	dst := v2.ExternalJob{}

	if err := testObj.ConvertTo(&dst); err == nil {
		t.Errorf("ConvertTo() error = %v", err)
	} else {
		require.ErrorIs(t, err, common.ErrCannotCastKeptnWorkloadInstance)
	}
}
