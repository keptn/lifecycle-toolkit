package evaluation

import (
	"context"
	"fmt"
	"strings"
	"testing"

	lifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry"
	telemetryfake "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry/fake"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

//nolint:dupl
func TestEvaluationHandler(t *testing.T) {
	tests := []struct {
		name            string
		object          client.Object
		createAttr      CreateEvaluationAttributes
		wantStatus      []lifecycle.ItemStatus
		wantSummary     apicommon.StatusSummary
		evalObj         lifecycle.KeptnEvaluation
		wantErr         error
		getSpanCalls    int
		unbindSpanCalls int
		events          []string
	}{
		{
			name:            "cannot unwrap object",
			object:          &lifecycle.KeptnEvaluation{},
			evalObj:         lifecycle.KeptnEvaluation{},
			createAttr:      CreateEvaluationAttributes{},
			wantStatus:      nil,
			wantSummary:     apicommon.StatusSummary{},
			wantErr:         controllererrors.ErrCannotWrapToPhaseItem,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name:    "no evaluations",
			object:  &lifecycle.KeptnAppVersion{},
			evalObj: lifecycle.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: lifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus:      []lifecycle.ItemStatus(nil),
			wantSummary:     apicommon.StatusSummary{},
			wantErr:         nil,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name: "evaluation not started",
			object: &lifecycle.KeptnAppVersion{
				Spec: lifecycle.KeptnAppVersionSpec{
					DeploymentTaskSpec: lifecycle.DeploymentTaskSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
			},
			evalObj: lifecycle.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: lifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []lifecycle.ItemStatus{
				{
					DefinitionName: "eval-def",
					Status:         apicommon.StatePending,
					Name:           "pre-eval-eval-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 1, Pending: 1},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 0,
		},
		{
			name: "already done evaluation",
			object: &lifecycle.KeptnAppVersion{
				Spec: lifecycle.KeptnAppVersionSpec{
					DeploymentTaskSpec: lifecycle.DeploymentTaskSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
				Status: lifecycle.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []lifecycle.ItemStatus{
						{
							DefinitionName: "eval-def",
							Status:         apicommon.StateSucceeded,
							Name:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: lifecycle.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: lifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []lifecycle.ItemStatus{
				{
					DefinitionName: "eval-def",
					Status:         apicommon.StateSucceeded,
					Name:           "pre-eval-eval-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 1, Succeeded: 1},
			wantErr:         nil,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name: "failed evaluation",
			object: &lifecycle.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: lifecycle.KeptnAppVersionSpec{
					DeploymentTaskSpec: lifecycle.DeploymentTaskSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
				Status: lifecycle.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []lifecycle.ItemStatus{
						{
							DefinitionName: "eval-def",
							Status:         apicommon.StateProgressing,
							Name:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: lifecycle.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "pre-eval-eval-def-",
				},
				Status: lifecycle.KeptnEvaluationStatus{
					OverallStatus: apicommon.StateFailed,
					EvaluationStatus: map[string]lifecycle.EvaluationStatusItem{
						"my-target": {
							Value:   "1",
							Status:  apicommon.StateFailed,
							Message: "failed",
						},
					},
				},
			},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: lifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []lifecycle.ItemStatus{
				{
					DefinitionName: "eval-def",
					Status:         apicommon.StateFailed,
					Name:           "pre-eval-eval-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 1, Failed: 1},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 1,
			events: []string{
				"evaluation of 'my-target' failed with value: '1' and reason: 'failed'",
			},
		},
		{
			name: "succeeded evaluation",
			object: &lifecycle.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: lifecycle.KeptnAppVersionSpec{
					DeploymentTaskSpec: lifecycle.DeploymentTaskSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
				Status: lifecycle.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []lifecycle.ItemStatus{
						{
							DefinitionName: "eval-def",
							Status:         apicommon.StateProgressing,
							Name:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: lifecycle.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "pre-eval-eval-def-",
				},
				Status: lifecycle.KeptnEvaluationStatus{
					OverallStatus: apicommon.StateSucceeded,
				},
			},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: lifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []lifecycle.ItemStatus{
				{
					DefinitionName: "eval-def",
					Status:         apicommon.StateSucceeded,
					Name:           "pre-eval-eval-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 1, Succeeded: 1},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lifecycle.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			spanHandlerMock := telemetryfake.ISpanHandlerMock{
				GetSpanFunc: func(ctx context.Context, tracer telemetry.ITracer, reconcileObject client.Object, phase string) (context.Context, trace.Span, error) {
					return context.TODO(), trace.SpanFromContext(context.TODO()), nil
				},
				UnbindSpanFunc: func(reconcileObject client.Object, phase string) error {
					return nil
				},
			}
			fakeRecorder := record.NewFakeRecorder(100)
			handler := NewHandler(
				fake.NewClientBuilder().WithObjects(&tt.evalObj).Build(),
				eventsender.NewK8sSender(fakeRecorder),
				ctrl.Log.WithName("controller"),
				noop.NewTracerProvider().Tracer("tracer"),
				scheme.Scheme,
				&spanHandlerMock)
			status, summary, err := handler.ReconcileEvaluations(context.TODO(), context.TODO(), tt.object, tt.createAttr)
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

			if tt.events != nil {
				for _, e := range tt.events {
					event := <-fakeRecorder.Events
					require.Equal(t, strings.Contains(event, tt.object.GetName()), true, "wrong appversion")
					require.Equal(t, strings.Contains(event, tt.object.GetNamespace()), true, "wrong namespace")
					require.Equal(t, strings.Contains(event, e), true, fmt.Sprintf("no %s found in %s", e, event))
				}

			}
		})
	}
}

func TestEvaluationHandler_createEvaluation(t *testing.T) {
	tests := []struct {
		name       string
		object     client.Object
		createAttr CreateEvaluationAttributes
		wantName   string
		wantErr    error
	}{
		{
			name:       "cannot unwrap object",
			object:     &lifecycle.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{},
			wantName:   "",
			wantErr:    controllererrors.ErrCannotWrapToPhaseItem,
		},
		{
			name: "created evaluation",
			object: &lifecycle.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: lifecycle.KeptnAppVersionSpec{
					DeploymentTaskSpec: lifecycle.DeploymentTaskSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
			},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: lifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantName: "pre-eval-eval-def-",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lifecycle.AddToScheme(scheme.Scheme)
			require.Nil(t, err)

			handler := NewHandler(
				fake.NewClientBuilder().Build(),
				eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				ctrl.Log.WithName("controller"),
				noop.NewTracerProvider().Tracer("tracer"),
				scheme.Scheme,
				&telemetryfake.ISpanHandlerMock{})

			name, err := handler.CreateKeptnEvaluation(context.TODO(), tt.object, tt.createAttr)
			require.True(t, strings.Contains(name, tt.wantName))
			require.Equal(t, tt.wantErr, err)
		})
	}
}
