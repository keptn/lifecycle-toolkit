package common

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
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
func TestEvaluationHandler(t *testing.T) {
	tests := []struct {
		name            string
		object          client.Object
		createAttr      CreateEvaluationAttributes
		wantStatus      []v1alpha3.ItemStatus
		wantSummary     apicommon.StatusSummary
		evalObj         v1alpha3.KeptnEvaluation
		wantErr         error
		getSpanCalls    int
		unbindSpanCalls int
		events          []string
	}{
		{
			name:            "cannot unwrap object",
			object:          &v1alpha3.KeptnEvaluation{},
			evalObj:         v1alpha3.KeptnEvaluation{},
			createAttr:      CreateEvaluationAttributes{},
			wantStatus:      nil,
			wantSummary:     apicommon.StatusSummary{},
			wantErr:         controllererrors.ErrCannotWrapToPhaseItem,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name:    "no evaluations",
			object:  &v1alpha3.KeptnAppVersion{},
			evalObj: v1alpha3.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: v1alpha3.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus:      []v1alpha3.ItemStatus(nil),
			wantSummary:     apicommon.StatusSummary{},
			wantErr:         nil,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name: "evaluation not started",
			object: &v1alpha3.KeptnAppVersion{
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
			},
			evalObj: v1alpha3.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: v1alpha3.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []v1alpha3.ItemStatus{
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
			object: &v1alpha3.KeptnAppVersion{
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
				Status: v1alpha3.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "eval-def",
							Status:         apicommon.StateSucceeded,
							Name:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: v1alpha3.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: v1alpha3.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []v1alpha3.ItemStatus{
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
			object: &v1alpha3.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
				Status: v1alpha3.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "eval-def",
							Status:         apicommon.StateProgressing,
							Name:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: v1alpha3.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "pre-eval-eval-def-",
				},
				Status: v1alpha3.KeptnEvaluationStatus{
					OverallStatus: apicommon.StateFailed,
					EvaluationStatus: map[string]v1alpha3.EvaluationStatusItem{
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
				Definition: v1alpha3.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []v1alpha3.ItemStatus{
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
			object: &v1alpha3.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
				Status: v1alpha3.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []v1alpha3.ItemStatus{
						{
							DefinitionName: "eval-def",
							Status:         apicommon.StateProgressing,
							Name:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: v1alpha3.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "pre-eval-eval-def-",
				},
				Status: v1alpha3.KeptnEvaluationStatus{
					OverallStatus: apicommon.StateSucceeded,
				},
			},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: v1alpha3.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []v1alpha3.ItemStatus{
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
			fakeRecorder := record.NewFakeRecorder(100)
			handler := EvaluationHandler{
				SpanHandler: &spanHandlerMock,
				Log:         ctrl.Log.WithName("controller"),
				EventSender: NewEventSender(fakeRecorder),
				Client:      fake.NewClientBuilder().WithObjects(&tt.evalObj).Build(),
				Tracer:      trace.NewNoopTracerProvider().Tracer("tracer"),
				Scheme:      scheme.Scheme,
			}
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
			object:     &v1alpha3.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{},
			wantName:   "",
			wantErr:    controllererrors.ErrCannotWrapToPhaseItem,
		},
		{
			name: "created evaluation",
			object: &v1alpha3.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: v1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: v1alpha3.KeptnAppSpec{
						PreDeploymentEvaluations: []string{"eval-def"},
					},
				},
			},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: v1alpha3.KeptnEvaluationDefinition{
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
			err := v1alpha3.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			handler := EvaluationHandler{
				SpanHandler: &kltfake.ISpanHandlerMock{},
				Log:         ctrl.Log.WithName("controller"),
				EventSender: NewEventSender(record.NewFakeRecorder(100)),
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
