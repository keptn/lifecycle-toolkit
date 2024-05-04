package evaluation

import (
	"context"
	"fmt"
	"strings"
	"testing"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry"
	telemetryfake "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry/fake"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/testcommon"
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
		wantStatus      []apilifecycle.ItemStatus
		wantSummary     apicommon.StatusSummary
		evalObj         apilifecycle.KeptnEvaluation
		evalDef         *apilifecycle.KeptnEvaluationDefinition
		wantErr         error
		getSpanCalls    int
		unbindSpanCalls int
		events          []string
	}{
		{
			name:            "cannot unwrap object",
			object:          &apilifecycle.KeptnEvaluation{},
			evalObj:         apilifecycle.KeptnEvaluation{},
			createAttr:      CreateEvaluationAttributes{},
			wantStatus:      nil,
			wantSummary:     apicommon.StatusSummary{},
			wantErr:         controllererrors.ErrCannotWrapToPhaseItem,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name:    "no evaluations",
			object:  &apilifecycle.KeptnAppVersion{},
			evalObj: apilifecycle.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: apilifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus:      []apilifecycle.ItemStatus(nil),
			wantSummary:     apicommon.StatusSummary{},
			wantErr:         nil,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name: "evaluation not started - could not find evaluationDefinition",
			object: &apilifecycle.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppContextSpec: apilifecycle.KeptnAppContextSpec{
						DeploymentTaskSpec: apilifecycle.DeploymentTaskSpec{
							PreDeploymentEvaluations: []string{"eval-def"},
						},
					},
				},
			},
			evalObj: apilifecycle.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: apilifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus:      nil,
			wantSummary:     apicommon.StatusSummary{Total: 1, Pending: 0},
			wantErr:         nil,
			getSpanCalls:    0,
			unbindSpanCalls: 0,
		},
		{
			name: "evaluations not started - could not find evaluationDefinition of one evaluation",
			object: &apilifecycle.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppContextSpec: apilifecycle.KeptnAppContextSpec{
						DeploymentTaskSpec: apilifecycle.DeploymentTaskSpec{
							PreDeploymentEvaluations: []string{"eval-def", "other-eval-def"},
						},
					},
				},
			},
			evalDef: &apilifecycle.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Namespace: testcommon.KeptnNamespace,
					Name:      "eval-def",
				},
			},
			evalObj: apilifecycle.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: apilifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []apilifecycle.ItemStatus{
				{
					DefinitionName: "eval-def",
					Status:         apicommon.StatePending,
					Name:           "pre-eval-eval-def-",
				},
			},
			wantSummary:     apicommon.StatusSummary{Total: 2, Pending: 1},
			wantErr:         nil,
			getSpanCalls:    1,
			unbindSpanCalls: 0,
		},
		{
			name: "evaluation not started - evaluationDefinition in default Keptn namespace",
			object: &apilifecycle.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppContextSpec: apilifecycle.KeptnAppContextSpec{
						DeploymentTaskSpec: apilifecycle.DeploymentTaskSpec{
							PreDeploymentEvaluations: []string{"eval-def"},
						},
					},
				},
			},
			evalDef: &apilifecycle.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Namespace: testcommon.KeptnNamespace,
					Name:      "eval-def",
				},
			},
			evalObj: apilifecycle.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: apilifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []apilifecycle.ItemStatus{
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
			name: "evaluation not started",
			object: &apilifecycle.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppContextSpec: apilifecycle.KeptnAppContextSpec{
						DeploymentTaskSpec: apilifecycle.DeploymentTaskSpec{
							PreDeploymentEvaluations: []string{"eval-def"},
						},
					},
				},
			},
			evalDef: &apilifecycle.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "eval-def",
				},
			},
			evalObj: apilifecycle.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: apilifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []apilifecycle.ItemStatus{
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
			object: &apilifecycle.KeptnAppVersion{
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppContextSpec: apilifecycle.KeptnAppContextSpec{
						DeploymentTaskSpec: apilifecycle.DeploymentTaskSpec{
							PreDeploymentEvaluations: []string{"eval-def"},
						},
					},
				},
				Status: apilifecycle.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []apilifecycle.ItemStatus{
						{
							DefinitionName: "eval-def",
							Status:         apicommon.StateSucceeded,
							Name:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: apilifecycle.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: apilifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []apilifecycle.ItemStatus{
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
			object: &apilifecycle.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppContextSpec: apilifecycle.KeptnAppContextSpec{
						DeploymentTaskSpec: apilifecycle.DeploymentTaskSpec{
							PreDeploymentEvaluations: []string{"eval-def"},
						},
					},
				},
				Status: apilifecycle.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []apilifecycle.ItemStatus{
						{
							DefinitionName: "eval-def",
							Status:         apicommon.StateProgressing,
							Name:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: apilifecycle.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "pre-eval-eval-def-",
				},
				Status: apilifecycle.KeptnEvaluationStatus{
					OverallStatus: apicommon.StateFailed,
					EvaluationStatus: map[string]apilifecycle.EvaluationStatusItem{
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
				Definition: apilifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []apilifecycle.ItemStatus{
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
			object: &apilifecycle.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppContextSpec: apilifecycle.KeptnAppContextSpec{
						DeploymentTaskSpec: apilifecycle.DeploymentTaskSpec{
							PreDeploymentEvaluations: []string{"eval-def"},
						},
					},
				},
				Status: apilifecycle.KeptnAppVersionStatus{
					PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
					PreDeploymentEvaluationTaskStatus: []apilifecycle.ItemStatus{
						{
							DefinitionName: "eval-def",
							Status:         apicommon.StateProgressing,
							Name:           "pre-eval-eval-def-",
						},
					},
				},
			},
			evalObj: apilifecycle.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
					Name:      "pre-eval-eval-def-",
				},
				Status: apilifecycle.KeptnEvaluationStatus{
					OverallStatus: apicommon.StateSucceeded,
				},
			},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: apilifecycle.KeptnEvaluationDefinition{
					ObjectMeta: v1.ObjectMeta{
						Name: "eval-def",
					},
				},
				CheckType: apicommon.PreDeploymentEvaluationCheckType,
			},
			wantStatus: []apilifecycle.ItemStatus{
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

	config.Instance().SetDefaultNamespace(testcommon.KeptnNamespace)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apilifecycle.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			spanHandlerMock := telemetryfake.ISpanHandlerMock{
				GetSpanFunc: func(ctx context.Context, tracer telemetry.ITracer, reconcileObject client.Object, phase string, links ...trace.Link) (context.Context, trace.Span, error) {
					return context.TODO(), trace.SpanFromContext(context.TODO()), nil
				},
				UnbindSpanFunc: func(reconcileObject client.Object, phase string) error {
					return nil
				},
			}
			initObjs := []client.Object{&tt.evalObj}
			if tt.evalDef != nil {
				initObjs = append(initObjs, tt.evalDef)
			}
			fakeRecorder := record.NewFakeRecorder(100)
			handler := NewHandler(
				fake.NewClientBuilder().WithObjects(initObjs...).Build(),
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
			object:     &apilifecycle.KeptnEvaluation{},
			createAttr: CreateEvaluationAttributes{},
			wantName:   "",
			wantErr:    controllererrors.ErrCannotWrapToPhaseItem,
		},
		{
			name: "created evaluation",
			object: &apilifecycle.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Namespace: "namespace",
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppContextSpec: apilifecycle.KeptnAppContextSpec{
						DeploymentTaskSpec: apilifecycle.DeploymentTaskSpec{
							PreDeploymentEvaluations: []string{"eval-def"},
						},
					},
				},
			},
			createAttr: CreateEvaluationAttributes{
				SpanName: "",
				Definition: apilifecycle.KeptnEvaluationDefinition{
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
			err := apilifecycle.AddToScheme(scheme.Scheme)
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
