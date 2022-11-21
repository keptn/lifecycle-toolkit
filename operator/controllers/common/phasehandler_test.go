package common

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestPhaseHandler(t *testing.T) {
	//phase, state, endtimeset
	requeueResult := ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}
	tests := []struct {
		name           string
		handler        PhaseHandler
		object         *v1alpha1.KeptnAppVersion
		phase          common.KeptnPhaseType
		reconcilePhase func() (common.KeptnState, error)
		wantObject     *v1alpha1.KeptnAppVersion
		want           *PhaseResult
		wantErr        error
		endTimeSet     bool
	}{
		{
			name: "cancelled",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
			},
			object: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status: common.StateCancelled,
				},
			},
			want:    &PhaseResult{Continue: false, Result: ctrl.Result{}},
			wantErr: nil,
			wantObject: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status: common.StateCancelled,
				},
			},
		},
		{
			name: "reconcilePhase error",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StatePending,
					CurrentPhase: common.PhaseAppDeployment.LongName,
				},
			},
			phase: common.PhaseAppDeployment,
			reconcilePhase: func() (common.KeptnState, error) {
				return "", fmt.Errorf("some err")
			},
			want:    &PhaseResult{Continue: false, Result: requeueResult},
			wantErr: fmt.Errorf("some err"),
			wantObject: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StatePending,
					CurrentPhase: common.PhaseAppDeployment.ShortName,
				},
			},
		},
		{
			name: "reconcilePhase pending state",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StatePending,
					CurrentPhase: common.PhaseAppDeployment.LongName,
				},
			},
			phase: common.PhaseAppDeployment,
			reconcilePhase: func() (common.KeptnState, error) {
				return common.StatePending, nil
			},
			want:    &PhaseResult{Continue: false, Result: requeueResult},
			wantErr: nil,
			wantObject: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateProgressing,
					CurrentPhase: common.PhaseAppDeployment.ShortName,
				},
			},
		},
		{
			name: "reconcilePhase progressing state",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StatePending,
					CurrentPhase: common.PhaseAppDeployment.LongName,
				},
			},
			phase: common.PhaseAppDeployment,
			reconcilePhase: func() (common.KeptnState, error) {
				return common.StateProgressing, nil
			},
			want:    &PhaseResult{Continue: false, Result: requeueResult},
			wantErr: nil,
			wantObject: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateProgressing,
					CurrentPhase: common.PhaseAppDeployment.ShortName,
				},
			},
		},
		{
			name: "reconcilePhase succeeded state",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StatePending,
					CurrentPhase: common.PhaseAppDeployment.LongName,
				},
			},
			phase: common.PhaseAppDeployment,
			reconcilePhase: func() (common.KeptnState, error) {
				return common.StateSucceeded, nil
			},
			want:    &PhaseResult{Continue: true, Result: requeueResult},
			wantErr: nil,
			wantObject: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateSucceeded,
					CurrentPhase: common.PhaseAppDeployment.ShortName,
				},
			},
		},
		{
			name: "reconcilePhase failed state",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateProgressing,
					CurrentPhase: common.PhaseAppPreEvaluation.LongName,
				},
			},
			phase: common.PhaseAppPreEvaluation,
			reconcilePhase: func() (common.KeptnState, error) {
				return common.StateFailed, nil
			},
			want:    &PhaseResult{Continue: false, Result: ctrl.Result{}},
			wantErr: nil,
			wantObject: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateFailed,
					CurrentPhase: common.PhaseAppPreEvaluation.ShortName,
					EndTime:      v1.Time{Time: time.Now().UTC()},
				},
			},
		},
		{
			name: "reconcilePhase unknown state",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateProgressing,
					CurrentPhase: common.PhaseAppPreEvaluation.LongName,
				},
			},
			phase: common.PhaseAppPreEvaluation,
			reconcilePhase: func() (common.KeptnState, error) {
				return common.StateUnknown, nil
			},
			want:    &PhaseResult{Continue: false, Result: requeueResult},
			wantErr: nil,
			wantObject: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateProgressing,
					CurrentPhase: common.PhaseAppPreEvaluation.ShortName,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.handler.HandlePhase(context.TODO(), context.TODO(), trace.NewNoopTracerProvider().Tracer("tracer"), tt.object, tt.phase, trace.SpanFromContext(context.TODO()), tt.reconcilePhase)
			require.Equal(t, tt.want, result)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.wantObject.Status.Status, tt.object.Status.Status)
			require.Equal(t, tt.wantObject.Status.CurrentPhase, tt.object.Status.CurrentPhase)
			require.Equal(t, tt.wantObject.Status.Status.IsFailed(), tt.object.IsEndTimeSet())
		})
	}
}

func TestPhaseHandler_GetEvaluationFailureReasons(t *testing.T) {
	tests := []struct {
		name         string
		handler      PhaseHandler
		object       *v1alpha1.KeptnAppVersion
		clientObject *v1alpha1.KeptnEvaluation
		phase        common.KeptnPhaseType
		want         []failedCheckReason
		wantErr      error
	}{
		{
			name: "status len is 0",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
			},
			object: &v1alpha1.KeptnAppVersion{
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateFailed,
					CurrentPhase: common.PhaseAppPreEvaluation.LongName,
				},
			},
			clientObject: &v1alpha1.KeptnEvaluation{},
			phase:        common.PhaseAppPreEvaluation,
			want:         nil,
			wantErr:      fmt.Errorf("evaluation status not found for /"),
		},
		{
			name: "cannot get evaluation",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
			},
			object: &v1alpha1.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Name:      "appversion",
					Namespace: "namespace",
				},
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateFailed,
					CurrentPhase: common.PhaseAppPreEvaluation.LongName,
					PreDeploymentEvaluationTaskStatus: []v1alpha1.EvaluationStatus{
						{
							EvaluationName: "eval-name",
						},
					},
				},
			},
			clientObject: &v1alpha1.KeptnEvaluation{},
			phase:        common.PhaseAppPreEvaluation,
			want:         nil,
			wantErr:      fmt.Errorf("evaluation eval-name not found for /"),
		},
		{
			name: "evaluation failed",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
			},
			object: &v1alpha1.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Name:      "appversion",
					Namespace: "namespace",
				},
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateFailed,
					CurrentPhase: common.PhaseAppPreEvaluation.LongName,
					PreDeploymentEvaluationTaskStatus: []v1alpha1.EvaluationStatus{
						{
							EvaluationName: "eval-name",
						},
					},
				},
			},
			clientObject: &v1alpha1.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Name:      "eval-name",
					Namespace: "namespace",
				},
				Status: v1alpha1.KeptnEvaluationStatus{
					EndTime: v1.Time{Time: time.Date(1, 1, 1, 1, 1, 1, 0, time.Local)},
					EvaluationStatus: map[string]v1alpha1.EvaluationStatusItem{
						"cpu": {
							Value:   "10",
							Status:  common.StateFailed,
							Message: "cpu failed",
						},
						"mem": {
							Value:   "10",
							Status:  common.StateSucceeded,
							Message: "mem passed",
						},
					},
				},
			},
			phase: common.PhaseAppPreEvaluation,
			want: []failedCheckReason{
				{
					Message: "evaluation of 'cpu' failed with value: '10' and reason: 'cpu failed'\n",
					Time:    time.Date(1, 1, 1, 1, 1, 1, 0, time.Local),
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1alpha1.AddToScheme(scheme.Scheme)
			client := fake.NewClientBuilder().WithObjects(tt.clientObject).Build()
			tt.handler.Client = client
			result, err := tt.handler.GetEvaluationFailureReasons(context.TODO(), tt.phase, tt.object)
			require.Equal(t, tt.want, result)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestPhaseHandler_GetTaskFailureReasons(t *testing.T) {
	tests := []struct {
		name         string
		handler      PhaseHandler
		object       *v1alpha1.KeptnAppVersion
		clientObject *v1alpha1.KeptnTask
		phase        common.KeptnPhaseType
		want         []failedCheckReason
		wantErr      error
	}{
		{
			name: "no state failed",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
			},
			object: &v1alpha1.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Name:      "appversion",
					Namespace: "namespace",
				},
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateSucceeded,
					CurrentPhase: common.PhaseAppPreDeployment.LongName,
					PreDeploymentTaskStatus: []v1alpha1.TaskStatus{
						{
							TaskName: "task-name",
							Status:   common.StateSucceeded,
						},
					},
				},
			},
			clientObject: &v1alpha1.KeptnTask{},
			phase:        common.PhaseAppPreDeployment,
			want:         []failedCheckReason{},
			wantErr:      nil,
		},
		{
			name: "cannot get task",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
			},
			object: &v1alpha1.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Name:      "appversion",
					Namespace: "namespace",
				},
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateFailed,
					CurrentPhase: common.PhaseAppPreDeployment.LongName,
					PreDeploymentTaskStatus: []v1alpha1.TaskStatus{
						{
							TaskName: "task-name",
							Status:   common.StateFailed,
						},
					},
				},
			},
			clientObject: &v1alpha1.KeptnTask{},
			phase:        common.PhaseAppPreDeployment,
			want:         nil,
			wantErr:      fmt.Errorf("task task-name not found for /"),
		},
		{
			name: "task failed",
			handler: PhaseHandler{
				SpanHandler: &SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
			},
			object: &v1alpha1.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Name:      "appversion",
					Namespace: "namespace",
				},
				Status: v1alpha1.KeptnAppVersionStatus{
					Status:       common.StateFailed,
					CurrentPhase: common.PhaseAppPreDeployment.LongName,
					PreDeploymentTaskStatus: []v1alpha1.TaskStatus{
						{
							TaskName: "task-name",
							Status:   common.StateFailed,
						},
					},
				},
			},
			clientObject: &v1alpha1.KeptnTask{
				ObjectMeta: v1.ObjectMeta{
					Name:      "task-name",
					Namespace: "namespace",
				},
				Status: v1alpha1.KeptnTaskStatus{
					Status:  common.StateFailed,
					Message: "task failed",
					EndTime: v1.Time{Time: time.Date(1, 1, 1, 1, 1, 1, 0, time.Local)},
				},
			},
			phase: common.PhaseAppPreDeployment,
			want: []failedCheckReason{
				{
					Message: "task 'task-name' failed with reason: 'task failed'\n",
					Time:    time.Date(1, 1, 1, 1, 1, 1, 0, time.Local),
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1alpha1.AddToScheme(scheme.Scheme)
			client := fake.NewClientBuilder().WithObjects(tt.clientObject).Build()
			tt.handler.Client = client
			result, err := tt.handler.GetTaskFailureReasons(context.TODO(), tt.phase, tt.object)
			require.Equal(t, tt.want, result)
			require.Equal(t, tt.wantErr, err)
		})
	}
}
