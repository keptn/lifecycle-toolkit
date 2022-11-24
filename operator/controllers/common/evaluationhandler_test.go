package common

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	kltfake "github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
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
		wantStatus      []v1alpha1.EvaluationStatus
		wantSummary     common.StatusSummary
		evalObj         v1alpha1.KeptnEvaluation
		wantErr         error
		getSpanCalls    int
		unbindSpanCalls int
	}{
		{
			name:            "cannot unwrap object",
			object:          &v1alpha1.KeptnEvaluation{},
			evalObj:         v1alpha1.KeptnEvaluation{},
			createAttr:      EvaluationCreateAttributes{},
			wantStatus:      nil,
			wantSummary:     common.StatusSummary{},
			wantErr:         ErrCannotWrapToPhaseItem,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name:    "no evaluations",
			object:  &v1alpha1.KeptnAppVersion{},
			evalObj: v1alpha1.KeptnEvaluation{},
			createAttr: EvaluationCreateAttributes{
				SpanName:             "",
				EvaluationDefinition: "",
				CheckType:            common.PreDeploymentEvaluationCheckType,
			},
			wantStatus:      []v1alpha1.EvaluationStatus(nil),
			wantSummary:     common.StatusSummary{},
			wantErr:         nil,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name: "evaluation not started",
			object: &v1alpha1.KeptnAppVersion{
				Spec: v1alpha1.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha1.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
			},
			evalObj: v1alpha1.KeptnEvaluation{},
			createAttr: EvaluationCreateAttributes{
				SpanName:             "",
				EvaluationDefinition: "eval-def",
				CheckType:            common.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []v1alpha1.EvaluationStatus{
				{
					EvaluationDefinitionName: "eval-def",
					Status:                   common.StatePending,
					EvaluationName:           "pre-eval-eval-def-",
				},
			},
			wantSummary:     common.StatusSummary{1, 0, 0, 0, 1, 0, 0},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 0,
		},
		{
			name: "already done evaluation",
			object: &v1alpha1.KeptnAppVersion{
				Spec: v1alpha1.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha1.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
				Status: v1alpha1.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: common.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []v1alpha1.EvaluationStatus{
						{
							EvaluationDefinitionName: "eval-def",
							Status:                   common.StateSucceeded,
							EvaluationName:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: v1alpha1.KeptnEvaluation{},
			createAttr: EvaluationCreateAttributes{
				SpanName:             "",
				EvaluationDefinition: "eval-def",
				CheckType:            common.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []v1alpha1.EvaluationStatus{
				{
					EvaluationDefinitionName: "eval-def",
					Status:                   common.StateSucceeded,
					EvaluationName:           "pre-eval-eval-def-",
				},
			},
			wantSummary:     common.StatusSummary{1, 0, 0, 1, 0, 0, 0},
			wantErr:         nil,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name: "failed evaluation",
			object: &v1alpha1.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha1.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha1.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
				Status: v1alpha1.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: common.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []v1alpha1.EvaluationStatus{
						{
							EvaluationDefinitionName: "eval-def",
							Status:                   common.StateProgressing,
							EvaluationName:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: v1alpha1.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "pre-eval-eval-def-",
				},
				Status: v1alpha1.KeptnEvaluationStatus{
					OverallStatus: common.StateFailed,
				},
			},
			createAttr: EvaluationCreateAttributes{
				SpanName:             "",
				EvaluationDefinition: "eval-def",
				CheckType:            common.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []v1alpha1.EvaluationStatus{
				{
					EvaluationDefinitionName: "eval-def",
					Status:                   common.StateFailed,
					EvaluationName:           "pre-eval-eval-def-",
				},
			},
			wantSummary:     common.StatusSummary{1, 0, 1, 0, 0, 0, 0},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 1,
		},
		{
			name: "succeeded evaluation",
			object: &v1alpha1.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha1.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha1.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
				Status: v1alpha1.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: common.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []v1alpha1.EvaluationStatus{
						{
							EvaluationDefinitionName: "eval-def",
							Status:                   common.StateProgressing,
							EvaluationName:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: v1alpha1.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "pre-eval-eval-def-",
				},
				Status: v1alpha1.KeptnEvaluationStatus{
					OverallStatus: common.StateSucceeded,
				},
			},
			createAttr: EvaluationCreateAttributes{
				SpanName:             "",
				EvaluationDefinition: "eval-def",
				CheckType:            common.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []v1alpha1.EvaluationStatus{
				{
					EvaluationDefinitionName: "eval-def",
					Status:                   common.StateSucceeded,
					EvaluationName:           "pre-eval-eval-def-",
				},
			},
			wantSummary:     common.StatusSummary{1, 0, 0, 1, 0, 0, 0},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v1alpha1.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			spanHandlerMock := kltfake.ISpanHandlerMock{
				GetSpanFunc: func(ctx context.Context, tracer trace.Tracer, reconcileObject client.Object, phase string) (context.Context, trace.Span, error) {
					return context.TODO(), trace.SpanFromContext(context.TODO()), nil
				},
				UnbindSpanFunc: func(reconcileObject client.Object, phase string) error {
					return nil
				},
			}
			handler := EvaluationHandler{
				SpanHandler: &spanHandlerMock,
				Log:         ctrl.Log.WithName("controller"),
				Recorder:    record.NewFakeRecorder(100),
				Client:      fake.NewClientBuilder().WithObjects(&tt.evalObj).Build(),
				Tracer:      trace.NewNoopTracerProvider().Tracer("tracer"),
				Scheme:      scheme.Scheme,
			}
			status, summary, err := handler.ReconcileEvaluations(context.TODO(), context.TODO(), tt.object, tt.createAttr)
			if len(tt.wantStatus) == len(status) {
				for j, item := range tt.wantStatus {
					require.Equal(t, tt.wantStatus[j].EvaluationDefinitionName, item.EvaluationDefinitionName)
					require.True(t, strings.Contains(item.EvaluationName, tt.wantStatus[j].EvaluationName))
					require.Equal(t, tt.wantStatus[j].Status, item.Status)
				}
			} else {
				fmt.Errorf("unexpected result, want %+v, got %+v", tt.wantStatus, status)
			}
			require.Equal(t, tt.wantSummary, summary)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.getSpanCalls, len(spanHandlerMock.GetSpanCalls()))
			require.Equal(t, tt.unbindSpanCalls, len(spanHandlerMock.UnbindSpanCalls()))
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
			object:     &v1alpha1.KeptnEvaluation{},
			createAttr: EvaluationCreateAttributes{},
			wantName:   "",
			wantErr:    ErrCannotWrapToPhaseItem,
		},
		{
			name: "created evaluation",
			object: &v1alpha1.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha1.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha1.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
			},
			createAttr: EvaluationCreateAttributes{
				SpanName:             "",
				EvaluationDefinition: "eval-def",
				CheckType:            common.PreDeploymentEvaluationCheckType,
			},
			wantName: "pre-eval-eval-def-",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v1alpha1.AddToScheme(scheme.Scheme)
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
