package analysis

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/analysis/fake"
	common "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis"
	fakeEvaluator "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/fake"
	metricstypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	fake2 "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAnalysisReconciler_SendResultToChannel(t *testing.T) {
	analysis, analysisDef, template, _ := getTestCRDs()
	fakeclient := fake2.NewClient(&analysis, &analysisDef, &template)
	res := metricstypes.AnalysisResult{
		Pass: true,
		ObjectiveResults: []metricstypes.ObjectiveResult{
			{
				Objective: analysisDef.Spec.Objectives[0],
			},
		},
	}

	req := controllerruntime.Request{
		NamespacedName: types.NamespacedName{Namespace: "default", Name: "my-analysis"},
	}
	mockFactory := func(ctx context.Context, analysisMoqParam *metricsapi.Analysis, obj []metricsapi.Objective, numWorkers int, c client.Client, log logr.Logger, namespace string) (context.Context, IAnalysisPool) {
		mymock := fake.IAnalysisPoolMock{
			DispatchAndCollectFunc: func(ctx context.Context) (map[string]metricsapi.ProviderResult, error) {
				return map[string]metricsapi.ProviderResult{}, nil
			},
		}
		return ctx, &mymock
	}

	a := &AnalysisReconciler{
		Client:                fakeclient,
		Scheme:                fakeclient.Scheme(),
		Log:                   testr.New(t),
		MaxWorkers:            2,
		NewWorkersPoolFactory: mockFactory,
		IAnalysisEvaluator: &fakeEvaluator.IAnalysisEvaluatorMock{
			EvaluateFunc: func(values map[string]metricsapi.ProviderResult, ad *metricsapi.AnalysisDefinition) metricstypes.AnalysisResult {
				return res
			}},
	}

	resChan := make(chan metricstypes.AnalysisCompletion)
	a.SetAnalysisResultsChannel(resChan)

	_, err := a.Reconcile(context.TODO(), req)
	require.Nil(t, err)

	select {
	case <-time.After(5 * time.Second):
		t.Error("timed out waiting for the analysis result to be reported")
	case analysisResult := <-resChan:
		require.Equal(t, "my-analysis", analysisResult.Analysis.Name)
	}
}

func TestAnalysisReconciler_Reconcile_BasicControlLoop(t *testing.T) {

	analysis, analysisDef, template, _ := getTestCRDs()

	currentTime := time.Now().Round(time.Minute)
	analysis2 := metricsapi.Analysis{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-analysis",
			Namespace: "default",
		},
		Spec: metricsapi.AnalysisSpec{
			Timeframe: metricsapi.Timeframe{
				From: metav1.Time{
					Time: currentTime,
				},
				To: metav1.Time{
					Time: currentTime,
				},
			},
			Args: map[string]string{
				"good": "good",
				"dot":  ".",
			},
			AnalysisDefinition: metricsapi.ObjectReference{
				Name:      "my-analysis-def",
				Namespace: "default2",
			},
		},
	}

	analysisCompleted := metricsapi.Analysis{
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
			Args: map[string]string{
				"good": "good",
				"dot":  ".",
			},
			AnalysisDefinition: metricsapi.ObjectReference{
				Name:      "my-analysis-def",
				Namespace: "default",
			},
		},
		Status: metricsapi.AnalysisStatus{
			State: metricsapi.StateCompleted,
			Raw:   "{\"objectiveResults\":null,\"totalScore\":0,\"maximumScore\":0,\"pass\":true,\"warning\":false}",
			Pass:  true,
		},
	}

	analysisDef2 := metricsapi.AnalysisDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-analysis-def",
			Namespace: "default2",
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

	req := controllerruntime.Request{
		NamespacedName: types.NamespacedName{Namespace: "default", Name: "my-analysis"},
	}
	mockFactory := func(ctx context.Context, analysisMoqParam *metricsapi.Analysis, obj []metricsapi.Objective, numWorkers int, c client.Client, log logr.Logger, namespace string) (context.Context, IAnalysisPool) {
		mymock := fake.IAnalysisPoolMock{
			DispatchAndCollectFunc: func(ctx context.Context) (map[string]metricsapi.ProviderResult, error) {
				return map[string]metricsapi.ProviderResult{}, nil
			},
		}
		return ctx, &mymock
	}

	tests := []struct {
		name        string
		client      client.Client
		req         controllerruntime.Request
		want        controllerruntime.Result
		wantErr     bool
		status      *metricsapi.AnalysisStatus
		res         metricstypes.AnalysisResult
		mockFactory NewWorkersPoolFactory
	}{
		{
			name:        "analysis does not exist, reconcile no status update",
			client:      fake2.NewClient(),
			want:        controllerruntime.Result{},
			wantErr:     false,
			status:      nil,
			res:         metricstypes.AnalysisResult{},
			mockFactory: mockFactory,
		}, {
			name:    "analysisDefinition does not exist, requeue no status update",
			client:  fake2.NewClient(&analysis),
			want:    controllerruntime.Result{Requeue: true, RequeueAfter: 10 * time.Second},
			wantErr: false,
			status: &metricsapi.AnalysisStatus{
				Timeframe: metricsapi.Timeframe{
					From: analysis.Spec.From,
					To:   analysis.Spec.To,
				},
				State: metricsapi.StatePending,
			},
			res:         metricstypes.AnalysisResult{Pass: false},
			mockFactory: mockFactory,
		}, {
			name:    "mockfactory failed",
			client:  fake2.NewClient(&analysis, &analysisDef, &template),
			want:    controllerruntime.Result{Requeue: true, RequeueAfter: 10 * time.Second},
			wantErr: false,
			status: &metricsapi.AnalysisStatus{
				Timeframe: metricsapi.Timeframe{
					From: analysis.Spec.From,
					To:   analysis.Spec.To,
				},
				State: metricsapi.StateProgressing,
			},
			res: metricstypes.AnalysisResult{Pass: false},
			mockFactory: func(ctx context.Context, analysisMoqParam *metricsapi.Analysis, obj []metricsapi.Objective, numWorkers int, c client.Client, log logr.Logger, namespace string) (context.Context, IAnalysisPool) {
				mymock := fake.IAnalysisPoolMock{
					DispatchAndCollectFunc: func(ctx context.Context) (map[string]metricsapi.ProviderResult, error) {
						return map[string]metricsapi.ProviderResult{}, fmt.Errorf("error")
					},
				}
				return ctx, &mymock
			},
		}, {
			name:    "succeeded, status updated",
			client:  fake2.NewClient(&analysis, &analysisDef, &template),
			want:    controllerruntime.Result{},
			wantErr: false,
			status: &metricsapi.AnalysisStatus{
				Timeframe: metricsapi.Timeframe{
					From: analysis.Spec.From,
					To:   analysis.Spec.To,
				},
				Raw:   "{\"objectiveResults\":null,\"totalScore\":0,\"maximumScore\":0,\"pass\":true,\"warning\":false}",
				Pass:  true,
				State: metricsapi.StateCompleted,
			},
			res:         metricstypes.AnalysisResult{Pass: true},
			mockFactory: mockFactory,
		}, {
			name:    "already completed analysis",
			client:  fake2.NewClient(&analysisCompleted),
			want:    controllerruntime.Result{},
			wantErr: false,
		}, {
			name:    "succeeded - analysis in different namespace, status updated",
			client:  fake2.NewClient(&analysis2, &analysisDef2, &template),
			want:    controllerruntime.Result{},
			wantErr: false,
			status: &metricsapi.AnalysisStatus{
				Timeframe: metricsapi.Timeframe{
					From: analysis.Spec.From,
					To:   analysis.Spec.To,
				},
				Raw:   "{\"objectiveResults\":null,\"totalScore\":0,\"maximumScore\":0,\"pass\":true,\"warning\":false}",
				Pass:  true,
				State: metricsapi.StateCompleted,
			},
			res:         metricstypes.AnalysisResult{Pass: true},
			mockFactory: mockFactory,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AnalysisReconciler{
				Client:                tt.client,
				Scheme:                tt.client.Scheme(),
				Log:                   testr.New(t),
				MaxWorkers:            2,
				NewWorkersPoolFactory: tt.mockFactory,
				IAnalysisEvaluator: &fakeEvaluator.IAnalysisEvaluatorMock{
					EvaluateFunc: func(values map[string]metricsapi.ProviderResult, ad *metricsapi.AnalysisDefinition) metricstypes.AnalysisResult {
						return tt.res
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
			if tt.status != nil {
				resAnalysis := metricsapi.Analysis{}
				err = tt.client.Get(context.TODO(), req.NamespacedName, &resAnalysis)
				require.Nil(t, err)
				require.Equal(t, *tt.status, resAnalysis.Status)
			}
		})
	}
}

func TestAnalysisReconciler_ExistingAnalysisStatusIsFlushedWhenEvaluationFinishes(t *testing.T) {
	analysis, analysisDef, template, _ := getTestCRDs()

	analysis.Status = metricsapi.AnalysisStatus{
		StoredValues: map[string]metricsapi.ProviderResult{
			"default": {
				Objective: metricsapi.ObjectReference{
					Name:      "my-analysis-def",
					Namespace: "default",
				},
				Value: "1",
			},
		},
	}

	mockFactory := func(ctx context.Context, analysisMoqParam *metricsapi.Analysis, obj []metricsapi.Objective, numWorkers int, c client.Client, log logr.Logger, namespace string) (context.Context, IAnalysisPool) {
		mymock := fake.IAnalysisPoolMock{
			DispatchAndCollectFunc: func(ctx context.Context) (map[string]metricsapi.ProviderResult, error) {
				return map[string]metricsapi.ProviderResult{}, nil
			},
		}
		return ctx, &mymock
	}

	fclient := fake2.NewClient(&analysis, &analysisDef, &template)
	a := &AnalysisReconciler{
		Client:                fclient,
		Scheme:                fclient.Scheme(),
		Log:                   testr.New(t),
		MaxWorkers:            2,
		NewWorkersPoolFactory: mockFactory,
		IAnalysisEvaluator: &fakeEvaluator.IAnalysisEvaluatorMock{
			EvaluateFunc: func(values map[string]metricsapi.ProviderResult, ad *metricsapi.AnalysisDefinition) metricstypes.AnalysisResult {
				return metricstypes.AnalysisResult{Pass: true}
			}},
	}

	req := controllerruntime.Request{
		NamespacedName: types.NamespacedName{Namespace: "default", Name: "my-analysis"},
	}

	status := &metricsapi.AnalysisStatus{
		Timeframe: metricsapi.Timeframe{
			From: analysis.Spec.From,
			To:   analysis.Spec.To,
		},
		Raw:   "{\"objectiveResults\":null,\"totalScore\":0,\"maximumScore\":0,\"pass\":true,\"warning\":false}",
		Pass:  true,
		State: metricsapi.StateCompleted,
	}

	got, err := a.Reconcile(context.TODO(), req)

	require.Nil(t, err)
	require.Equal(t, controllerruntime.Result{}, got)
	resAnalysis := metricsapi.Analysis{}
	err = fclient.Get(context.TODO(), req.NamespacedName, &resAnalysis)
	require.Nil(t, err)
	require.Nil(t, resAnalysis.Status.StoredValues)
	require.Equal(t, *status, resAnalysis.Status)

}

func TestAnalysisReconciler_AnalysisTimeframeIsDerivedFromDurationString(t *testing.T) {
	analysis, _, _, _ := getTestCRDs()

	analysis.Spec.Timeframe = metricsapi.Timeframe{Recent: metav1.Duration{Duration: 5 * time.Minute}}

	mockFactory := func(ctx context.Context, analysisMoqParam *metricsapi.Analysis, obj []metricsapi.Objective, numWorkers int, c client.Client, log logr.Logger, namespace string) (context.Context, IAnalysisPool) {
		mymock := fake.IAnalysisPoolMock{
			DispatchAndCollectFunc: func(ctx context.Context) (map[string]metricsapi.ProviderResult, error) {
				return map[string]metricsapi.ProviderResult{}, nil
			},
		}
		return ctx, &mymock
	}

	fclient := fake2.NewClient(&analysis)
	a := &AnalysisReconciler{
		Client:                fclient,
		Scheme:                fclient.Scheme(),
		Log:                   testr.New(t),
		MaxWorkers:            2,
		NewWorkersPoolFactory: mockFactory,
		IAnalysisEvaluator: &fakeEvaluator.IAnalysisEvaluatorMock{
			EvaluateFunc: func(values map[string]metricsapi.ProviderResult, ad *metricsapi.AnalysisDefinition) metricstypes.AnalysisResult {
				return metricstypes.AnalysisResult{Pass: true}
			}},
	}

	req := controllerruntime.Request{
		NamespacedName: types.NamespacedName{Namespace: "default", Name: "my-analysis"},
	}

	got, err := a.Reconcile(context.TODO(), req)

	// expect to be re-queued, since the AnalysisDefinition was not there, but the from/to timestamps should be set
	// as soon as the reconciliation has started
	require.Nil(t, err)
	require.True(t, got.Requeue)
	resAnalysis := metricsapi.Analysis{}
	err = fclient.Get(context.TODO(), req.NamespacedName, &resAnalysis)
	require.Nil(t, err)

	currentTime := time.Now().UTC()
	require.WithinDuration(t, currentTime, resAnalysis.Status.Timeframe.GetTo(), time.Minute)
	require.WithinDuration(t, currentTime.Add(-5*time.Minute), resAnalysis.Status.Timeframe.GetFrom(), time.Minute)

}

func getTestCRDs() (metricsapi.Analysis, metricsapi.AnalysisDefinition, metricsapi.AnalysisValueTemplate, metricsapi.KeptnMetricsProvider) {
	currentTime := time.Now().Round(time.Minute)
	analysis := metricsapi.Analysis{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-analysis",
			Namespace: "default",
		},
		Spec: metricsapi.AnalysisSpec{
			Timeframe: metricsapi.Timeframe{
				From: metav1.Time{
					Time: currentTime,
				},
				To: metav1.Time{
					Time: currentTime,
				},
			},
			Args: map[string]string{
				"good": "good",
				"dot":  ".",
			},
			AnalysisDefinition: metricsapi.ObjectReference{
				Name:      "my-analysis-def",
				Namespace: "default",
			},
		},
		Status: metricsapi.AnalysisStatus{
			State: metricsapi.StatePending,
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
			Query: "this is a {{.good}} query{{.dot}}",
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

func Test_extractMissingObjectives(t *testing.T) {

	missing := metricsapi.ObjectReference{
		Name:      "missing",
		Namespace: "test",
	}

	done := metricsapi.ObjectReference{
		Name:      "done",
		Namespace: "test",
	}

	needToRetry := metricsapi.ObjectReference{
		Name:      "need-to-retry",
		Namespace: "test",
	}

	ad := &metricsapi.AnalysisDefinition{Spec: metricsapi.AnalysisDefinitionSpec{Objectives: []metricsapi.Objective{
		{
			AnalysisValueTemplateRef: missing,
			Target:                   metricsapi.Target{},
			Weight:                   1,
			KeyObjective:             false,
		},
		{
			AnalysisValueTemplateRef: done,
			Target:                   metricsapi.Target{},
			Weight:                   1,
			KeyObjective:             false,
		},
		{
			AnalysisValueTemplateRef: needToRetry,
			Target:                   metricsapi.Target{},
			Weight:                   1,
			KeyObjective:             false,
		},
	}}}

	existingValues := map[string]metricsapi.ProviderResult{
		common.ComputeKey(ad.Spec.Objectives[1].AnalysisValueTemplateRef): {
			Value: "1.0",
		},
		common.ComputeKey(ad.Spec.Objectives[2].AnalysisValueTemplateRef): {
			ErrMsg: "error",
		},
	}
	todo, existing := extractMissingObjectives(ad.Spec.Objectives, existingValues)

	require.Len(t, todo, 2)
	require.Equal(t, missing, todo[0].AnalysisValueTemplateRef)
	require.Equal(t, needToRetry, todo[1].AnalysisValueTemplateRef)
	require.Len(t, existing, 1)
	require.Equal(t, "1.0", existing[common.ComputeKey(done)].Value)

	// verify that the analysisDefinition has not been changed
	require.Len(t, ad.Spec.Objectives, 3)
}
