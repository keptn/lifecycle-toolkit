package common

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestPhaseHandler(t *testing.T) {
	requeueResult := ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}
	tests := []struct {
		name           string
		handler        PhaseHandler
		object         *v1alpha3.KeptnAppVersion
		phase          apicommon.KeptnPhaseType
		reconcilePhase func(phaseCtx context.Context) (apicommon.KeptnState, error)
		wantObject     *v1alpha3.KeptnAppVersion
		want           *PhaseResult
		wantErr        error
		endTimeSet     bool
	}{
		{
			name: "deprecated",
			handler: PhaseHandler{
				SpanHandler: &telemetry.SpanHandler{},
			},
			object: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status: apicommon.StateDeprecated,
				},
			},
			want:    &PhaseResult{Continue: false, Result: ctrl.Result{}},
			wantErr: nil,
			wantObject: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status: apicommon.StateDeprecated,
				},
			},
		},
		{
			name: "reconcilePhase error",
			handler: PhaseHandler{
				SpanHandler: &telemetry.SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StatePending,
					CurrentPhase: apicommon.PhaseAppDeployment.LongName,
				},
			},
			phase: apicommon.PhaseAppDeployment,
			reconcilePhase: func(phaseCtx context.Context) (apicommon.KeptnState, error) {
				return "", fmt.Errorf("some err")
			},
			want:    &PhaseResult{Continue: false, Result: requeueResult},
			wantErr: fmt.Errorf("some err"),
			wantObject: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StatePending,
					CurrentPhase: apicommon.PhaseAppDeployment.ShortName,
				},
			},
		},
		{
			name: "reconcilePhase pending state",
			handler: PhaseHandler{
				SpanHandler: &telemetry.SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StatePending,
					CurrentPhase: apicommon.PhaseAppDeployment.LongName,
				},
			},
			phase: apicommon.PhaseAppDeployment,
			reconcilePhase: func(phaseCtx context.Context) (apicommon.KeptnState, error) {
				return apicommon.StatePending, nil
			},
			want:    &PhaseResult{Continue: false, Result: requeueResult},
			wantErr: nil,
			wantObject: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StateProgressing,
					CurrentPhase: apicommon.PhaseAppDeployment.ShortName,
				},
			},
		},
		{
			name: "reconcilePhase progressing state",
			handler: PhaseHandler{
				SpanHandler: &telemetry.SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StatePending,
					CurrentPhase: apicommon.PhaseAppDeployment.LongName,
				},
			},
			phase: apicommon.PhaseAppDeployment,
			reconcilePhase: func(phaseCtx context.Context) (apicommon.KeptnState, error) {
				return apicommon.StateProgressing, nil
			},
			want:    &PhaseResult{Continue: false, Result: requeueResult},
			wantErr: nil,
			wantObject: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StateProgressing,
					CurrentPhase: apicommon.PhaseAppDeployment.ShortName,
				},
			},
		},
		{
			name: "reconcilePhase succeeded state",
			handler: PhaseHandler{
				SpanHandler: &telemetry.SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StatePending,
					CurrentPhase: apicommon.PhaseAppDeployment.LongName,
				},
			},
			phase: apicommon.PhaseAppDeployment,
			reconcilePhase: func(phaseCtx context.Context) (apicommon.KeptnState, error) {
				return apicommon.StateSucceeded, nil
			},
			want:    &PhaseResult{Continue: true, Result: requeueResult},
			wantErr: nil,
			wantObject: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StatePending,
					CurrentPhase: apicommon.PhaseAppDeployment.ShortName,
				},
			},
		},
		{
			name: "reconcilePhase failed state",
			handler: PhaseHandler{
				SpanHandler: &telemetry.SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StateProgressing,
					CurrentPhase: apicommon.PhaseAppPreEvaluation.LongName,
				},
			},
			phase: apicommon.PhaseAppPreEvaluation,
			reconcilePhase: func(phaseCtx context.Context) (apicommon.KeptnState, error) {
				return apicommon.StateFailed, nil
			},
			want:    &PhaseResult{Continue: false, Result: ctrl.Result{}},
			wantErr: nil,
			wantObject: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StateFailed,
					CurrentPhase: apicommon.PhaseAppPreEvaluation.ShortName,
					EndTime:      v1.Time{Time: time.Now().UTC()},
				},
			},
		},
		{
			name: "reconcilePhase unknown state",
			handler: PhaseHandler{
				SpanHandler: &telemetry.SpanHandler{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: NewK8sSender(record.NewFakeRecorder(100)),
				Client:      fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
			},
			object: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StateProgressing,
					CurrentPhase: apicommon.PhaseAppPreEvaluation.LongName,
				},
			},
			phase: apicommon.PhaseAppPreEvaluation,
			reconcilePhase: func(phaseCtx context.Context) (apicommon.KeptnState, error) {
				return apicommon.StateUnknown, nil
			},
			want:    &PhaseResult{Continue: false, Result: requeueResult},
			wantErr: nil,
			wantObject: &v1alpha3.KeptnAppVersion{
				Status: v1alpha3.KeptnAppVersionStatus{
					Status:       apicommon.StateProgressing,
					CurrentPhase: apicommon.PhaseAppPreEvaluation.ShortName,
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
