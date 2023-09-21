package common

import (
	"context"
	"strings"
	"testing"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	kltfake "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/fake"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

//nolint:dupl
func TestTaskHandler(t *testing.T) {
	tests := []struct {
		name            string
		object          client.Object
		createAttr      CreateTaskAttributes
		wantStatus      []v1alpha3.ItemStatus
		wantSummary     apicommon.StatusSummary
		taskObj         v1alpha3.KeptnTask
		taskDef         *v1alpha3.KeptnTaskDefinition
		wantErr         error
		getSpanCalls    int
		unbindSpanCalls int
	}{
		{
			name:            "cannot unwrap object",
			object:          &v1alpha3.KeptnTask{},
			taskObj:         v1alpha3.KeptnTask{},
			createAttr:      CreateTaskAttributes{},
			wantStatus:      nil,
			wantSummary:     apicommon.StatusSummary{},
			wantErr:         controllererrors.ErrCannotWrapToPhaseItem,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name:    "no tasks",
			object:  &v1alpha3.KeptnAppVersion{},
			taskObj: v1alpha3.KeptnTask{},
			createAttr: CreateTaskAttributes{
				SpanName:   "",
				Definition: v1alpha3.KeptnTaskDefinition{},
				CheckType:  apicommon.PreDeploymentCheckType,
			},
			wantStatus:      []v1alpha3.ItemStatus(nil),
			wantSummary:     apicommon.StatusSummary{},
			wantErr:         nil,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name: "task not started - could not find taskDefinition",
			object: &v1alpha3.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentTasks: []string{"task-def"},
					},
				},
			},
			taskObj: v1alpha3.KeptnTask{},
			createAttr: CreateTaskAttributes{
				SpanName: "",
				Definition: v1alpha3.KeptnTaskDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "task-def",
					},
				},
				CheckType: apicommon.PreDeploymentCheckType,
			},
			wantStatus:      nil,
			wantSummary:     apicommon.StatusSummary{Total: 1, Pending: 0},
			wantErr:         nil,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name: "tasks not started - could not find taskDefinition of one task",
			object: &v1alpha3.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentTasks: []string{"task-def", "other-task-def"},
					},
				},
			},
			taskDef: &v1alpha3.KeptnTaskDefinition{
				ObjectMeta: v1.ObjectMeta{
					Namespace: KeptnNamespace,
					Name:      "task-def",
				},
			},
			taskObj: v1alpha3.KeptnTask{},
			createAttr: CreateTaskAttributes{
				SpanName: "",
				Definition: v1alpha3.KeptnTaskDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "task-def",
					},
				},
				CheckType: apicommon.PreDeploymentCheckType,
			},
			wantStatus: []v1alpha3.ItemStatus{
				{
					DefinitionName: "task-def",
					Status:         apicommon.StatePending,
					Name:           "pre-task-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 2, Pending: 1},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 0,
		},
		{
			name: "task not started - taskDefinition in default KLT namespace",
			object: &v1alpha3.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentTasks: []string{"task-def"},
					},
				},
			},
			taskDef: &v1alpha3.KeptnTaskDefinition{
				ObjectMeta: v1.ObjectMeta{
					Namespace: KeptnNamespace,
					Name:      "task-def",
				},
			},
			taskObj: v1alpha3.KeptnTask{},
			createAttr: CreateTaskAttributes{
				SpanName: "",
				Definition: v1alpha3.KeptnTaskDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "task-def",
					},
				},
				CheckType: apicommon.PreDeploymentCheckType,
			},
			wantStatus: []v1alpha3.ItemStatus{
				{
					DefinitionName: "task-def",
					Status:         apicommon.StatePending,
					Name:           "pre-task-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 1, Pending: 1},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 0,
		},
		{
			name: "task not started",
			object: &v1alpha3.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentTasks: []string{"task-def"},
					},
				},
			},
			taskDef: &v1alpha3.KeptnTaskDefinition{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "task-def",
				},
			},
			taskObj: v1alpha3.KeptnTask{},
			createAttr: CreateTaskAttributes{
				SpanName: "",
				Definition: v1alpha3.KeptnTaskDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "task-def",
					},
				},
				CheckType: apicommon.PreDeploymentCheckType,
			},
			wantStatus: []v1alpha3.ItemStatus{
				{
					DefinitionName: "task-def",
					Status:         apicommon.StatePending,
					Name:           "pre-task-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 1, Pending: 1},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 0,
		},
		{
			name: "already done task",
			object: &v1alpha3.KeptnAppVersion{
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentTasks: []string{"task-def"},
					},
				},
				Status: v1alpha3.KeptnAppVersionStatus{
					PreDeploymentStatus: apicommon.StateSucceeded,
					PreDeploymentTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "task-def",
							Status:         apicommon.StateSucceeded,
							Name:           "pre-task-def-",
						},
					},
				},
			},
			taskObj: v1alpha3.KeptnTask{},
			createAttr: CreateTaskAttributes{
				SpanName: "",
				Definition: v1alpha3.KeptnTaskDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "task-def",
					},
				},
				CheckType: apicommon.PreDeploymentCheckType,
			},
			wantStatus: []v1alpha3.ItemStatus{
				{
					DefinitionName: "task-def",
					Status:         apicommon.StateSucceeded,
					Name:           "pre-task-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 1, Succeeded: 1},
			wantErr:         nil,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name: "failed task",
			object: &v1alpha3.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentTasks: []string{"task-def"},
					},
				},
				Status: v1alpha3.KeptnAppVersionStatus{
					PreDeploymentStatus: apicommon.StateSucceeded,
					PreDeploymentTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "task-def",
							Status:         apicommon.StateProgressing,
							Name:           "pre-task-def-",
						},
					},
				},
			},
			taskObj: v1alpha3.KeptnTask{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "pre-task-def-",
				},
				Status: v1alpha3.KeptnTaskStatus{
					Status: apicommon.StateFailed,
				},
			},
			createAttr: CreateTaskAttributes{
				SpanName: "",
				Definition: v1alpha3.KeptnTaskDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "task-def",
					},
				},
				CheckType: apicommon.PreDeploymentCheckType,
			},
			wantStatus: []v1alpha3.ItemStatus{
				{
					DefinitionName: "task-def",
					Status:         apicommon.StateFailed,
					Name:           "pre-task-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 1, Failed: 1},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 1,
		},
		{
			name: "succeeded task",
			object: &v1alpha3.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentTasks: []string{"task-def"},
					},
				},
				Status: v1alpha3.KeptnAppVersionStatus{
					PreDeploymentStatus: apicommon.StateSucceeded,
					PreDeploymentTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "task-def",
							Status:         apicommon.StateProgressing,
							Name:           "pre-task-def-",
						},
					},
				},
			},
			taskObj: v1alpha3.KeptnTask{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "pre-task-def-",
				},
				Status: v1alpha3.KeptnTaskStatus{
					Status: apicommon.StateSucceeded,
				},
			},
			createAttr: CreateTaskAttributes{
				SpanName: "",
				Definition: v1alpha3.KeptnTaskDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "task-def",
					},
				},
				CheckType: apicommon.PreDeploymentCheckType,
			},
			wantStatus: []v1alpha3.ItemStatus{
				{
					DefinitionName: "task-def",
					Status:         apicommon.StateSucceeded,
					Name:           "pre-task-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 1, Succeeded: 1},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 1,
		},
	}
	config.Instance().SetDefaultNamespace(KeptnNamespace)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v1alpha3.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			spanHandlerMock := kltfake.ISpanHandlerMock{
				GetSpanFunc: func(ctx context.Context, tracer trace.Tracer, reconcileObject client.Object, phase string) (context.Context, trace.Span, error) {
					return context.TODO(), trace.SpanFromContext(context.TODO()), nil
				},
				UnbindSpanFunc: func(reconcileObject client.Object, phase string) error {
					return nil
				},
			}
			initObjs := []client.Object{&tt.taskObj}
			if tt.taskDef != nil {
				initObjs = append(initObjs, tt.taskDef)
			}
			handler := TaskHandler{
				SpanHandler: &spanHandlerMock,
				Log:         ctrl.Log.WithName("controller"),
				EventSender: NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithObjects(initObjs...).Build(),
				Tracer:      trace.NewNoopTracerProvider().Tracer("tracer"),
				Scheme:      scheme.Scheme,
			}
			status, summary, err := handler.ReconcileTasks(context.TODO(), context.TODO(), tt.object, tt.createAttr)
			if len(tt.wantStatus) == len(status) {
				for j, item := range status {
					require.Equal(t, tt.wantStatus[j].DefinitionName, item.DefinitionName)
					require.True(t, strings.Contains(item.Name, tt.wantStatus[j].Name))
					require.Equal(t, tt.wantStatus[j].Status, item.Status)
				}
			} else {
				t.Errorf("unexpected result, want %+v, got %+v", tt.wantStatus, status)
			}
			require.Equal(t, tt.wantSummary, summary)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.getSpanCalls, len(spanHandlerMock.GetSpanCalls()))
			require.Equal(t, tt.unbindSpanCalls, len(spanHandlerMock.UnbindSpanCalls()))
		})
	}
}

func TestTaskHandler_createTask(t *testing.T) {
	tests := []struct {
		name       string
		object     client.Object
		createAttr CreateTaskAttributes
		wantName   string
		wantErr    error
	}{
		{
			name:       "cannot unwrap object",
			object:     &v1alpha3.KeptnEvaluation{},
			createAttr: CreateTaskAttributes{},
			wantName:   "",
			wantErr:    controllererrors.ErrCannotWrapToPhaseItem,
		},
		{
			name: "created task",
			object: &v1alpha3.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentTasks: []string{"task-def"},
					},
				},
			},
			createAttr: CreateTaskAttributes{
				SpanName:  "",
				CheckType: apicommon.PreDeploymentCheckType,
				Definition: v1alpha3.KeptnTaskDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "task-def",
					},
				},
			},
			wantName: "pre-task-def-",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v1alpha3.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			handler := TaskHandler{
				SpanHandler: &kltfake.ISpanHandlerMock{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().Build(),
				Tracer:      trace.NewNoopTracerProvider().Tracer("tracer"),
				Scheme:      scheme.Scheme,
			}
			name, err := handler.CreateKeptnTask(context.TODO(), "namespace", tt.object, tt.createAttr)
			require.True(t, strings.Contains(name, tt.wantName))
			require.Equal(t, tt.wantErr, err)
		})
	}
}
