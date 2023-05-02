//nolint:dupl
package keptnevaluation

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/providers"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const KltNamespace = "klt-namespace"

func TestKeptnEvaluationReconciler_fetchDefinition(t *testing.T) {

	metricEvalDef, EvalDef := setupEvalDefinitions()
	DTProv, PromProv := setupProviders()
	client := fake.NewClient(metricEvalDef, EvalDef, DTProv, PromProv)

	r := &KeptnEvaluationReconciler{
		Client: client,
		Scheme: client.Scheme(),
		Log:    testr.New(t),
	}

	tests := []struct {
		name                 string
		namespacedDefinition types.NamespacedName
		wantDef              *klcv1alpha3.KeptnEvaluationDefinition
		wantErr              bool
	}{
		{
			name: "keptn metrics",
			namespacedDefinition: types.NamespacedName{
				Namespace: KltNamespace,
				Name:      "myKeptn",
			},
			wantDef: metricEvalDef,
		},

		{
			name: "Unexisting Evaluation Def",
			namespacedDefinition: types.NamespacedName{
				Namespace: KltNamespace,
				Name:      "whatever",
			},
			wantDef: nil,
			wantErr: true,
		},
		{
			name: "Unexisting Provider",
			namespacedDefinition: types.NamespacedName{
				Namespace: KltNamespace,
				Name:      "mydef",
			},
			wantDef: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := r.fetchDefinition(context.TODO(), tt.namespacedDefinition)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchDefinitionAndProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantDef != nil {
				require.Equal(t, got.Name, tt.wantDef.Name)
			} else {
				require.Nil(t, got)
			}
		})
	}
}

func setupEvalDefinitions() (*klcv1alpha3.KeptnEvaluationDefinition, *klcv1alpha3.KeptnEvaluationDefinition) {
	metricEvalDef := &klcv1alpha3.KeptnEvaluationDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: KltNamespace,
			Name:      "myKeptn",
		},
		Spec: klcv1alpha3.KeptnEvaluationDefinitionSpec{
			Objectives: nil,
		},
		Status: klcv1alpha3.KeptnEvaluationDefinitionStatus{},
	}

	EvalDef := &klcv1alpha3.KeptnEvaluationDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: KltNamespace,
			Name:      "mdef",
		},
		Spec: klcv1alpha3.KeptnEvaluationDefinitionSpec{
			Objectives: nil,
		},
		Status: klcv1alpha3.KeptnEvaluationDefinitionStatus{},
	}

	return metricEvalDef, EvalDef
}

func setupProviders() (*metricsapi.KeptnMetricsProvider, *metricsapi.KeptnMetricsProvider) {
	DTProv := &metricsapi.KeptnMetricsProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      providers.DynatraceProviderName,
			Namespace: KltNamespace,
		},
	}

	PromProv := &metricsapi.KeptnMetricsProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      providers.PrometheusProviderName,
			Namespace: KltNamespace,
		},
	}

	return DTProv, PromProv
}

func TestKeptnEvaluationReconciler_Reconcile_FailEvaluation(t *testing.T) {

	const namespace = "my-namespace"
	metric := &metricsapi.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-metric",
			Namespace: namespace,
		},
		Status: metricsapi.KeptnMetricStatus{
			Value: "10",
		},
	}

	evaluationDefinition := &klcv1alpha3.KeptnEvaluationDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-definition",
			Namespace: namespace,
		},
		Spec: klcv1alpha3.KeptnEvaluationDefinitionSpec{
			Objectives: []klcv1alpha3.Objective{
				{
					KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
						Name:      metric.Name,
						Namespace: namespace,
					},
					EvaluationTarget: "<5",
				},
			},
		},
	}

	evaluation := &klcv1alpha3.KeptnEvaluation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-evaluation",
			Namespace: namespace,
		},
		Spec: klcv1alpha3.KeptnEvaluationSpec{
			EvaluationDefinition: evaluationDefinition.Name,
			Retries:              1,
		},
	}

	reconciler, fakeClient := setupReconcilerAndClient(t, metric, evaluationDefinition, evaluation)

	request := controllerruntime.Request{
		NamespacedName: types.NamespacedName{
			Namespace: namespace,
			Name:      evaluation.Name,
		},
	}

	reconcile, err := reconciler.Reconcile(context.TODO(), request)

	require.Nil(t, err)
	require.True(t, reconcile.Requeue)

	updatedEvaluation := &klcv1alpha3.KeptnEvaluation{}
	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      evaluation.Name,
	}, updatedEvaluation)

	require.Nil(t, err)

	require.Equal(t, common.StateFailed, updatedEvaluation.Status.EvaluationStatus[metric.Name].Status)
	require.Equal(t, "value '10' did not meet objective '<5'", updatedEvaluation.Status.EvaluationStatus[metric.Name].Message)
}

func TestKeptnEvaluationReconciler_Reconcile_SucceedEvaluation(t *testing.T) {

	const namespace = "my-namespace"
	metric := &metricsapi.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-metric",
			Namespace: namespace,
		},
		Status: metricsapi.KeptnMetricStatus{
			Value: "10",
		},
	}

	evaluationDefinition := &klcv1alpha3.KeptnEvaluationDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-definition",
			Namespace: namespace,
		},
		Spec: klcv1alpha3.KeptnEvaluationDefinitionSpec{
			Objectives: []klcv1alpha3.Objective{
				{
					KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
						Name:      metric.Name,
						Namespace: namespace,
					},
					EvaluationTarget: "<11",
				},
			},
		},
	}

	evaluation := &klcv1alpha3.KeptnEvaluation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-evaluation",
			Namespace: namespace,
		},
		Spec: klcv1alpha3.KeptnEvaluationSpec{
			EvaluationDefinition: evaluationDefinition.Name,
			Retries:              1,
		},
	}

	reconciler, fakeClient := setupReconcilerAndClient(t, metric, evaluationDefinition, evaluation)

	request := controllerruntime.Request{
		NamespacedName: types.NamespacedName{
			Namespace: namespace,
			Name:      evaluation.Name,
		},
	}

	reconcile, err := reconciler.Reconcile(context.TODO(), request)

	require.Nil(t, err)
	require.False(t, reconcile.Requeue)

	updatedEvaluation := &klcv1alpha3.KeptnEvaluation{}
	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      evaluation.Name,
	}, updatedEvaluation)

	require.Nil(t, err)

	require.Equal(t, common.StateSucceeded, updatedEvaluation.Status.EvaluationStatus[metric.Name].Status)
	require.Equal(t, "value '10' met objective '<11'", updatedEvaluation.Status.EvaluationStatus[metric.Name].Message)
}

func setupReconcilerAndClient(t *testing.T, objects ...client.Object) (*KeptnEvaluationReconciler, client.Client) {
	scheme := runtime.NewScheme()

	err := klcv1alpha3.AddToScheme(scheme)
	require.Nil(t, err)

	err = metricsapi.AddToScheme(scheme)
	require.Nil(t, err)

	// fake a tracer
	tr := &fake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	tf := &fake.TracerFactoryMock{GetTracerFunc: func(name string) trace.Tracer {
		return tr
	}}

	recorder := record.NewFakeRecorder(100)

	fakeClient := k8sfake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build()

	provider := metric.NewMeterProvider()
	meter := provider.Meter("keptn/task")

	r := &KeptnEvaluationReconciler{
		Client:        fakeClient,
		Scheme:        fakeClient.Scheme(),
		Log:           logr.Logger{},
		Recorder:      recorder,
		Meters:        controllercommon.SetUpKeptnTaskMeters(meter),
		TracerFactory: tf,
	}
	return r, fakeClient
}
