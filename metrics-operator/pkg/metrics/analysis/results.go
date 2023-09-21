package analysis

import (
	"context"
	"fmt"
	"sync"

	analysistypes "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

const analysisResultMetricName = "keptn_analysis_result"
const objectiveResultMetricName = "keptn_objective_result"

type Metrics struct {
	AnalysisResult  *prometheus.GaugeVec
	ObjectiveResult *prometheus.GaugeVec
}

// use singleton pattern here to avoid registering the same metrics on Prometheus multiple times
var instance *resultsReporter
var once sync.Once

func GetResultsReporter(ctx context.Context, res chan analysistypes.AnalysisCompletion) *resultsReporter {
	once.Do(func() {
		instance = &resultsReporter{}
		instance.initialize(ctx, res)
	})

	return instance
}

type resultsReporter struct {
	metrics Metrics
	mtx     sync.Mutex
}

func (r *resultsReporter) initialize(ctx context.Context, res chan analysistypes.AnalysisCompletion) {
	labelNamesAnalysis := []string{"name", "namespace", "from", "to"}
	a := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: analysisResultMetricName,
		Help: "Result of Analysis",
	}, labelNamesAnalysis)
	err := prometheus.Register(a)

	if err != nil {
		klog.Errorf("Could not register Analysis results as Prometheus metric: %v", err)
	}

	labelNames := []string{"name", "namespace", "analysis_name", "analysis_namespace", "key_objective", "weight", "from", "to"}
	o := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: objectiveResultMetricName,
		Help: "Result of the Analysis Objective",
	}, labelNames)
	err = prometheus.Register(o)
	if err != nil {
		klog.Errorf("Could not register Analysis Objective results as Prometheus metric: %v", err)
	}

	r.metrics = Metrics{
		AnalysisResult:  a,
		ObjectiveResult: o,
	}

	go r.watchForResults(ctx, res)
}

func (r *resultsReporter) watchForResults(ctx context.Context, res chan analysistypes.AnalysisCompletion) {
	for {
		select {
		case <-ctx.Done():
			klog.Info("Exiting due to termination of context")
			return
		case finishedAnalysis := <-res:
			r.reportResult(finishedAnalysis)
		}
	}
}

func (r *resultsReporter) reportResult(finishedAnalysis analysistypes.AnalysisCompletion) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	f := finishedAnalysis.Analysis.Spec.From.String()
	t := finishedAnalysis.Analysis.Spec.To.String()
	labelsAnalysis := prometheus.Labels{
		"name":      finishedAnalysis.Analysis.Name,
		"namespace": finishedAnalysis.Analysis.Namespace,
		"from":      f,
		"to":        t,
	}
	if m, err := r.metrics.AnalysisResult.GetMetricWith(labelsAnalysis); err == nil {
		m.Set(finishedAnalysis.Result.GetAchievedPercentage())
	} else {
		klog.Errorf("unable to set value for analysis result metric: %v", err)
	}
	// expose also the individual objectives
	for _, o := range finishedAnalysis.Result.ObjectiveResults {
		name := o.Objective.AnalysisValueTemplateRef.Name
		ns := o.Objective.AnalysisValueTemplateRef.Namespace
		labelsObjective := prometheus.Labels{
			"name":               name,
			"namespace":          ns,
			"analysis_name":      finishedAnalysis.Analysis.Name,
			"analysis_namespace": finishedAnalysis.Analysis.Namespace,
			"key_objective":      fmt.Sprintf("%v", o.Objective.KeyObjective),
			"weight":             fmt.Sprintf("%v", o.Objective.Weight),
			"from":               f,
			"to":                 t,
		}
		if m, err := r.metrics.ObjectiveResult.GetMetricWith(labelsObjective); err == nil {
			m.Set(o.Value)
		} else {
			klog.Errorf("unable to set value for objective result metric: %v", err)
		}
	}
}
