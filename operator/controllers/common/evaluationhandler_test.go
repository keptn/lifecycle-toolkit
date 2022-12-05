package common

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2/common"
	kltfake "github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestEvaluationHandler(t *testing.T) {
	tests := []struct {
		name            string
		object          client.Object
		createAttr      EvaluationCreateAttributes
		wantStatus      []v1alpha2.EvaluationStatus
		wantSummary     apicommon.StatusSummary
		evalObj         v1alpha2.KeptnEvaluation
		wantErr         error
		getSpanCalls    int
		unbindSpanCalls int
		events          []string
	}{
		{
			name:            "cannot unwrap object",
			object:          &v1alpha2.KeptnEvaluation{},
			evalObj:         v1alpha2.KeptnEvaluation{},
			createAttr:      EvaluationCreateAttributes{},
			wantStatus:      nil,
			wantSummary:     apicommon.StatusSummary{},
			wantErr:         controllererrors.ErrCannotWrapToPhaseItem,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name:    "no evaluations",
			object:  &v1alpha2.KeptnAppVersion{},
			evalObj: v1alpha2.KeptnEvaluation{},
			createAttr: EvaluationCreateAttributes{
				SpanName:             "",
				EvaluationDefinition: "",
				CheckType:            apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus:      []v1alpha2.EvaluationStatus(nil),
			wantSummary:     apicommon.StatusSummary{},
			wantErr:         nil,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name: "evaluation not started",
			object: &v1alpha2.KeptnAppVersion{
				Spec: v1alpha2.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha2.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
			},
			evalObj: v1alpha2.KeptnEvaluation{},
			createAttr: EvaluationCreateAttributes{
				SpanName:             "",
				EvaluationDefinition: "eval-def",
				CheckType:            apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []v1alpha2.EvaluationStatus{
				{
					EvaluationDefinitionName: "eval-def",
					Status:                   apicommon.StatePending,
					EvaluationName:           "pre-eval-eval-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 1, Pending: 1},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 0,
		},
		{
			name: "already done evaluation",
			object: &v1alpha2.KeptnAppVersion{
				Spec: v1alpha2.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha2.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
				Status: v1alpha2.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []v1alpha2.EvaluationStatus{
						{
							EvaluationDefinitionName: "eval-def",
							Status:                   apicommon.StateSucceeded,
							EvaluationName:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: v1alpha2.KeptnEvaluation{},
			createAttr: EvaluationCreateAttributes{
				SpanName:             "",
				EvaluationDefinition: "eval-def",
				CheckType:            apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []v1alpha2.EvaluationStatus{
				{
					EvaluationDefinitionName: "eval-def",
					Status:                   apicommon.StateSucceeded,
					EvaluationName:           "pre-eval-eval-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 1, Succeeded: 1},
			wantErr:         nil,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name: "failed evaluation",
			object: &v1alpha2.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha2.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha2.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
				Status: v1alpha2.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []v1alpha2.EvaluationStatus{
						{
							EvaluationDefinitionName: "eval-def",
							Status:                   apicommon.StateProgressing,
							EvaluationName:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: v1alpha2.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "pre-eval-eval-def-",
				},
				Status: v1alpha2.KeptnEvaluationStatus{
					OverallStatus: apicommon.StateFailed,
					EvaluationStatus: map[string]v1alpha2.EvaluationStatusItem{
						"my-target": {
							Value:   "1",
							Status:  apicommon.StateFailed,
							Message: "failed",
						},
					},
				},
			},
			createAttr: EvaluationCreateAttributes{
				SpanName:             "",
				EvaluationDefinition: "eval-def",
				CheckType:            apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []v1alpha2.EvaluationStatus{
				{
					EvaluationDefinitionName: "eval-def",
					Status:                   apicommon.StateFailed,
					EvaluationName:           "pre-eval-eval-def-",
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
			object: &v1alpha2.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha2.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha2.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
				Status: v1alpha2.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []v1alpha2.EvaluationStatus{
						{
							EvaluationDefinitionName: "eval-def",
							Status:                   apicommon.StateProgressing,
							EvaluationName:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: v1alpha2.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "pre-eval-eval-def-",
				},
				Status: v1alpha2.KeptnEvaluationStatus{
					OverallStatus: apicommon.StateSucceeded,
				},
			},
			createAttr: EvaluationCreateAttributes{
				SpanName:             "",
				EvaluationDefinition: "eval-def",
				CheckType:            apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []v1alpha2.EvaluationStatus{
				{
					EvaluationDefinitionName: "eval-def",
					Status:                   apicommon.StateSucceeded,
					EvaluationName:           "pre-eval-eval-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 1, Succeeded: 1},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 1,
			events: []string{
				"ReconcileEvaluationSucceeded",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v1alpha2.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			spanHandlerMock := kltfake.ISpanHandlerMock{
				GetSpanFunc: func(ctx context.Context, tracer trace.Tracer, reconcileObject client.Object, phase string) (context.Context, trace.Span, error) {
					return context.TODO(), trace.SpanFromContext(context.TODO()), nil
				},
				UnbindSpanFunc: func(reconcileObject client.Object, phase string) error {
					return nil
				},
			}
			fakeRecorder := record.NewFakeRecorder(100)
			handler := EvaluationHandler{
				SpanHandler: &spanHandlerMock,
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    fakeRecorder,
				Client:      fake.NewClientBuilder().WithObjects(&tt.evalObj).Build(),
				Tracer:      trace.NewNoopTracerProvider().Tracer("tracer"),
				Scheme:      scheme.Scheme,
			}
			status, summary, err := handler.ReconcileEvaluations(context.TODO(), context.TODO(), tt.object, tt.createAttr)
			if len(tt.wantStatus) == len(status) {
				for j, item := range status {
					require.Equal(t, tt.wantStatus[j].EvaluationDefinitionName, item.EvaluationDefinitionName)
					require.True(t, strings.Contains(item.EvaluationName, tt.wantStatus[j].EvaluationName))
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
		createAttr EvaluationCreateAttributes
		wantName   string
		wantErr    error
	}{
		{
			name:       "cannot unwrap object",
			object:     &v1alpha2.KeptnEvaluation{},
			createAttr: EvaluationCreateAttributes{},
			wantName:   "",
			wantErr:    controllererrors.ErrCannotWrapToPhaseItem,
		},
		{
			name: "created evaluation",
			object: &v1alpha2.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha2.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha2.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
			},
			createAttr: EvaluationCreateAttributes{
				SpanName:             "",
				EvaluationDefinition: "eval-def",
				CheckType:            apicommon.PreDeploymentEvaluationCheckType,
			},
			wantName: "pre-eval-eval-def-",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v1alpha2.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			handler := EvaluationHandler{
				SpanHandler: &kltfake.ISpanHandlerMock{},
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
				Client:      fake.NewClientBuilder().Build(),
				Tracer:      trace.NewNoopTracerProvider().Tracer("tracer"),
				Scheme:      scheme.Scheme,
			}
			name, err := handler.CreateKeptnEvaluation(context.TODO(), "namespace", tt.object, tt.createAttr)
			require.True(t, strings.Contains(name, tt.wantName))
			require.Equal(t, tt.wantErr, err)
		})
	}
}
