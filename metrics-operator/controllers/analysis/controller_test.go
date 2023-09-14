package analysis

import (
	"context"
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

func TestAnalysisReconciler_Reconcile_BasicControlLoop(t *testing.T) {

	analysis, analysisDef, template, _ := getTestCRDs()

	tests := []struct {
		name    string
		client  client.Client
		req     controllerruntime.Request
		want    controllerruntime.Result
		wantErr bool
		status  *metricsapi.AnalysisStatus
		res     metricstypes.AnalysisResult
	}{
		{
			name:    "analysis does not exist, reconcile no status update",
			client:  fake2.NewClient(),
			want:    controllerruntime.Result{},
			wantErr: false,
			status:  nil,
			res:     metricstypes.AnalysisResult{},
		}, {
			name:    "analysisDefinition does not exist, requeue no status update",
			client:  fake2.NewClient(&analysis),
			want:    controllerruntime.Result{Requeue: true, RequeueAfter: 10 * time.Second},
			wantErr: false,
			status:  &metricsapi.AnalysisStatus{},
			res:     metricstypes.AnalysisResult{Pass: false},
		}, {
			name:    "succeeded, status updated",
			client:  fake2.NewClient(&analysis, &analysisDef, &template),
			want:    controllerruntime.Result{},
			wantErr: false,
			status:  &metricsapi.AnalysisStatus{Raw: "{\"objectiveResults\":null,\"totalScore\":0,\"maximumScore\":0,\"pass\":true,\"warning\":false}", Pass: true},
			res:     metricstypes.AnalysisResult{Pass: true},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AnalysisReconciler{
				Client:                tt.client,
				Scheme:                tt.client.Scheme(),
				Log:                   testr.New(t),
				MaxWorkers:            2,
				NewWorkersPoolFactory: mockFactory,
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
			Args: map[string]string{
				"good": "good",
				"dot":  ".",
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
