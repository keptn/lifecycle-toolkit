package phase

import (
	"context"
	"fmt"
	"testing"
	"time"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestHandler(t *testing.T) {
	requeueResult := ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}
	tests := []struct {
		name           string
		handler        Handler
		object         *apilifecycle.KeptnAppVersion
		phase          apicommon.KeptnPhaseType
		reconcilePhase func(phaseCtx context.Context) (apicommon.KeptnState, error)
		wantObject     *apilifecycle.KeptnAppVersion
		want           PhaseResult
		wantErr        error
		endTimeSet     bool
	}{
		{
			name: "deprecated",
			handler: Handler{
				SpanHandler: &telemetry.Handler{},
			},
			object: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status: apicommon.StateDeprecated,
				},
			},
			want:    PhaseResult{Continue: false, Result: ctrl.Result{}},
			wantErr: nil,
			wantObject: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status: apicommon.StateDeprecated,
				},
			},
		},
		{
			name: "reconcilePhase error",
			handler: Handler{
				SpanHandler: &telemetry.Handler{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status:       apicommon.StatePending,
					CurrentPhase: apicommon.PhaseAppDeployment.LongName,
				},
			},
			phase: apicommon.PhaseAppDeployment,
			reconcilePhase: func(phaseCtx context.Context) (apicommon.KeptnState, error) {
				return "", fmt.Errorf("some err")
			},
			want:    PhaseResult{Continue: false, Result: requeueResult},
			wantErr: fmt.Errorf("some err"),
			wantObject: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status:       apicommon.StatePending,
					CurrentPhase: apicommon.PhaseAppDeployment.ShortName,
				},
			},
		},
		{
			name: "reconcilePhase pending state",
			handler: Handler{
				SpanHandler: &telemetry.Handler{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status:       apicommon.StatePending,
					CurrentPhase: apicommon.PhaseAppDeployment.LongName,
				},
			},
			phase: apicommon.PhaseAppDeployment,
			reconcilePhase: func(phaseCtx context.Context) (apicommon.KeptnState, error) {
				return apicommon.StatePending, nil
			},
			want:    PhaseResult{Continue: false, Result: requeueResult},
			wantErr: nil,
			wantObject: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status:       apicommon.StateProgressing,
					CurrentPhase: apicommon.PhaseAppDeployment.ShortName,
				},
			},
		},
		{
			name: "reconcilePhase progressing state",
			handler: Handler{
				SpanHandler: &telemetry.Handler{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status:       apicommon.StatePending,
					CurrentPhase: apicommon.PhaseAppDeployment.LongName,
				},
			},
			phase: apicommon.PhaseAppDeployment,
			reconcilePhase: func(phaseCtx context.Context) (apicommon.KeptnState, error) {
				return apicommon.StateProgressing, nil
			},
			want:    PhaseResult{Continue: false, Result: requeueResult},
			wantErr: nil,
			wantObject: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status:       apicommon.StateProgressing,
					CurrentPhase: apicommon.PhaseAppDeployment.ShortName,
				},
			},
		},
		{
			name: "reconcilePhase succeeded state",
			handler: Handler{
				SpanHandler: &telemetry.Handler{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status:       apicommon.StatePending,
					CurrentPhase: apicommon.PhaseAppDeployment.LongName,
				},
			},
			phase: apicommon.PhaseAppDeployment,
			reconcilePhase: func(phaseCtx context.Context) (apicommon.KeptnState, error) {
				return apicommon.StateSucceeded, nil
			},
			want:    PhaseResult{Continue: true, Result: requeueResult},
			wantErr: nil,
			wantObject: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status:       apicommon.StatePending,
					CurrentPhase: apicommon.PhaseAppDeployment.ShortName,
				},
			},
		},
		{
			name: "reconcilePhase failed state",
			handler: Handler{
				SpanHandler: &telemetry.Handler{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status:       apicommon.StateProgressing,
					CurrentPhase: apicommon.PhaseAppPreEvaluation.LongName,
				},
			},
			phase: apicommon.PhaseAppPreEvaluation,
			reconcilePhase: func(phaseCtx context.Context) (apicommon.KeptnState, error) {
				return apicommon.StateFailed, nil
			},
			want:    PhaseResult{Continue: false, Result: ctrl.Result{}},
			wantErr: nil,
			wantObject: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status:       apicommon.StateFailed,
					CurrentPhase: apicommon.PhaseAppPreEvaluation.ShortName,
					EndTime:      v1.Time{Time: time.Now().UTC()},
				},
			},
		},
		{
			name: "reconcilePhase unknown state",
			handler: Handler{
				SpanHandler: &telemetry.Handler{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status:       apicommon.StateProgressing,
					CurrentPhase: apicommon.PhaseAppPreEvaluation.LongName,
				},
			},
			phase: apicommon.PhaseAppPreEvaluation,
			reconcilePhase: func(phaseCtx context.Context) (apicommon.KeptnState, error) {
				return apicommon.StateUnknown, nil
			},
			want:    PhaseResult{Continue: false, Result: requeueResult},
			wantErr: nil,
			wantObject: &apilifecycle.KeptnAppVersion{
				Status: apilifecycle.KeptnAppVersionStatus{
					Status:       apicommon.StateProgressing,
					CurrentPhase: apicommon.PhaseAppPreEvaluation.ShortName,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.handler.HandlePhase(context.TODO(), context.TODO(), noop.NewTracerProvider().Tracer("tracer"), tt.object, tt.phase, tt.reconcilePhase)
			require.Equal(t, tt.want, result)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.wantObject.Status.Status, tt.object.Status.Status)
			require.Equal(t, tt.wantObject.Status.CurrentPhase, tt.object.Status.CurrentPhase)
			require.Equal(t, tt.wantObject.Status.Status.IsFailed(), tt.object.IsEndTimeSet())
		})
	}
}

func TestNewHandler(t *testing.T) {
	spanHandler := &telemetry.Handler{}
	log := ctrl.Log.WithName("controller")
	eventSender := eventsender.NewK8sSender(record.NewFakeRecorder(100))
	client := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

	handler := NewHandler(client, eventSender, log, spanHandler)

	require.NotNil(t, handler)
	require.NotNil(t, handler.Client)
	require.NotNil(t, handler.EventSender)
	require.NotNil(t, handler.Log)
	require.NotNil(t, handler.SpanHandler)
}
