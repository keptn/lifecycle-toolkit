package analysis

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/analysis/fake"
	fakeEvaluator "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/fake"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	fake2 "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAnalysisReconciler_Reconcile_BasicControlLoop(t *testing.T) {

	analysis, analysisDef, template, provider := getTestCRDs()

	tests := []struct {
		name    string
		client  client.Client
		req     controllerruntime.Request
		want    controllerruntime.Result
		wantErr bool
	}{
		{
			name:    "analysis does not exist, reconcile",
			client:  fake2.NewClient(),
			want:    controllerruntime.Result{},
			wantErr: false,
		}, {
			name:    "analysisDefinition does not exist, requeue",
			client:  fake2.NewClient(&analysis),
			want:    controllerruntime.Result{Requeue: true, RequeueAfter: 10 * time.Second},
			wantErr: false,
		}, {
			name:    "analysisValueTemplate does not exist, reconcile",
			client:  fake2.NewClient(&analysis, &analysisDef),
			want:    controllerruntime.Result{},
			wantErr: false,
		}, {
			name:    "metrics provider does not exist, reconcile",
			client:  fake2.NewClient(&analysis, &analysisDef, &template),
			want:    controllerruntime.Result{},
			wantErr: false,
		}, {
			name:   "provider exist, collect result, do nothing", //TODO: when scoring is there we need more
			client: fake2.NewClient(&analysis, &analysisDef, &template, &provider),
			req: controllerruntime.Request{
				NamespacedName: types.NamespacedName{Namespace: "default", Name: "my-analysis"},
			},
			want:    controllerruntime.Result{},
			wantErr: false,
		},
	}

	req := controllerruntime.Request{
		NamespacedName: types.NamespacedName{Namespace: "default", Name: "my-analysis"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AnalysisReconciler{
				Client:                tt.client,
				Scheme:                tt.client.Scheme(),
				Log:                   testr.New(t),
				MaxWorkers:            2,
				NewWorkersPoolFactory: NewWorkersPool,
				IAnalysisEvaluator: &fakeEvaluator.IAnalysisEvaluatorMock{
					EvaluateFunc: func(values map[string]metricstypes.ProviderResult, ad *metricsapi.AnalysisDefinition) metricstypes.AnalysisResult {
						return metricstypes.AnalysisResult{}
					}},
			}
			got, err := a.Reconcile(context.TODO(), req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reconcile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestCRDs() (metricsapi.Analysis, metricsapi.AnalysisDefinition, metricsapi.AnalysisValueTemplate, metricsapi.KeptnMetricsProvider) {
	analysis := metricsapi.Analysis{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-analysis",
			Namespace: "default",
		},
		Spec: metricsapi.AnalysisSpec{
			Timeframe: metricsapi.Timeframe{
				From: metav1.Time{
					Time: time.Now(),
				},
				To: metav1.Time{
					Time: time.Now(),
				},
			},
			AnalysisDefinition: metricsapi.ObjectReference{
				Name:      "my-analysis-def",
				Namespace: "default",
			},
		},
	}

	analysisDef := metricsapi.AnalysisDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-analysis-def",
			Namespace: "default",
		},
		Spec: metricsapi.AnalysisDefinitionSpec{
			Objectives: []metricsapi.Objective{
				{
					AnalysisValueTemplateRef: metricsapi.ObjectReference{
						Name:      "my-template",
						Namespace: "default",
					},
					Weight:       1,
					KeyObjective: false,
				},
			},
			TotalScore: metricsapi.TotalScore{
				PassPercentage:    0,
				WarningPercentage: 0,
			},
		},
	}

	template := metricsapi.AnalysisValueTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-template",
			Namespace: "default",
		},
		Spec: metricsapi.AnalysisValueTemplateSpec{
			Provider: metricsapi.ObjectReference{
				Name:      "my-provider",
				Namespace: "default",
			},
			Query: "testquery",
		},
	}

	provider := metricsapi.KeptnMetricsProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-provider",
			Namespace: "default",
		},
		Spec: metricsapi.KeptnMetricsProviderSpec{
			Type:         "prometheus",
			TargetServer: "localhost:2000",
		},
	}
	return analysis, analysisDef, template, provider
}

func TestAnalysisReconciler_Reconcile_WithMockedWorkers(t *testing.T) {
	analysis, analysisDef, template, provider := getTestCRDs()
	fclient := fake2.NewClient(&analysis, &analysisDef, &template, &provider)
	collectorCount := int32(0)
	dispatchCount := int32(0)

	mockFactory := func(analysisMoqParam *metricsapi.Analysis, definition *metricsapi.AnalysisDefinition, numWorkers int, c client.Client, log logr.Logger, namespace string) IAnalysisPool {
		mymock := fake.MyAnalysisPoolMock{
			CollectAnalysisResultsFunc: func() map[string]metricstypes.ProviderResult {
				atomic.AddInt32(&collectorCount, 1)
				return nil
			},
			DispatchObjectivesFunc: func(ctx context.Context) {
				atomic.AddInt32(&dispatchCount, 1)
			},
		}
		return &mymock
	}

	a := &AnalysisReconciler{
		Client:                fclient,
		Scheme:                fclient.Scheme(),
		Log:                   testr.New(t),
		MaxWorkers:            2,
		NewWorkersPoolFactory: mockFactory,
		IAnalysisEvaluator: &fakeEvaluator.IAnalysisEvaluatorMock{
			EvaluateFunc: func(values map[string]metricstypes.ProviderResult, ad *metricsapi.AnalysisDefinition) metricstypes.AnalysisResult {
				return metricstypes.AnalysisResult{Pass: true}
			}},
	}
	req := controllerruntime.Request{
		NamespacedName: types.NamespacedName{Namespace: "default", Name: "my-analysis"},
	}
	got, err := a.Reconcile(context.TODO(), req)
	require.Nil(t, err)
	require.Equal(t, got, controllerruntime.Result{})
	resAnalysis := metricsapi.Analysis{}
	err = fclient.Get(context.TODO(), req.NamespacedName, &resAnalysis)
	require.Nil(t, err)
	require.Equal(t, "{[] 0 0 true false}", resAnalysis.Status) //TODO change when introducing status
}
