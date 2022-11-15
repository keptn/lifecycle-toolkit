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
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestPhaseHandler(t *testing.T) {
	requeueResult := ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}
	tests := []struct {
		name           string
		handler        PhaseHandler
		object         client.Object
		phase          common.KeptnPhaseType
		reconcilePhase func() (common.KeptnState, error)
		want           *PhaseResult
		wantErr        error
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.handler.HandlePhase(context.TODO(), context.TODO(), trace.NewNoopTracerProvider().Tracer("tracer"), tt.object, tt.phase, trace.SpanFromContext(context.TODO()), tt.reconcilePhase)
			require.Equal(t, tt.want, result)
			require.Equal(t, tt.wantErr, err)
		})
	}
}
