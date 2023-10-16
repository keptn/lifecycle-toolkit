/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	ctrlcommon "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/aggregation"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const MB = 1 << (10 * 2)

// KeptnMetricReconciler reconciles a KeptnMetric object
type KeptnMetricReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
	providers.ProviderFactory
}

// clusterrole
// +kubebuilder:rbac:groups=metrics.keptn.sh,resources=providers,verbs=get;list;watch
// +kubebuilder:rbac:groups=metrics.keptn.sh,resources=keptnmetrics,verbs=get;list;watch;
// +kubebuilder:rbac:groups=metrics.keptn.sh,resources=keptnmetrics/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=metrics.keptn.sh,resources=keptnmetrics/finalizers,verbs=update
// +kubebuilder:rbac:groups=metrics.keptn.sh,resources=keptnmetricsproviders,verbs=get;list;watch;
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch

// role
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *KeptnMetricReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	requestInfo := ctrlcommon.GetRequestInfo(req)
	r.Log.Info("Reconciling Metric", "requestInfo", requestInfo)
	metric := &metricsapi.KeptnMetric{}

	if err := r.Client.Get(ctx, req.NamespacedName, metric); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			r.Log.Info("Metric resource not found. Ignoring since object must be deleted", "requestInfo", requestInfo)
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get the Metric", "requestInfo", requestInfo)
		return ctrl.Result{}, nil
	}

	fetchTime := metric.Status.LastUpdated.Add(time.Second * time.Duration(metric.Spec.FetchIntervalSeconds))
	if time.Now().Before(fetchTime) {
		diff := time.Until(fetchTime)
		r.Log.Info("Metric has not been updated for the configured interval. Skipping", "requestInfo", requestInfo)
		return ctrl.Result{Requeue: true, RequeueAfter: diff}, nil
	}

	metricProvider, err := r.fetchProvider(ctx, types.NamespacedName{Name: metric.Spec.Provider.Name, Namespace: metric.Namespace})
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info(err.Error()+", ignoring error since object must be deleted", "requestInfo", requestInfo)
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
		}
		r.Log.Error(err, "Failed to retrieve the provider", "requestInfo", requestInfo)
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}
	// load the provider
	provider, err2 := r.ProviderFactory(metricProvider.GetType(), r.Log, r.Client)
	if err2 != nil {
		r.Log.Error(err2, "Failed to get the correct Metric Provider")
		return ctrl.Result{Requeue: false}, err2
	}
	reconcile := ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}

	value, rawValue, err := r.getResults(ctx, metric, provider, metricProvider)
	if err != nil {
		r.Log.Error(err, "Failed to evaluate the query", "Response from provider was:", (string)(rawValue))
		metric.Status.ErrMsg = err.Error()
		metric.Status.Value = ""
		metric.Status.RawValue = cupSize(rawValue)
		metric.Status.LastUpdated = metav1.Time{Time: time.Now()}
		reconcile = ctrl.Result{Requeue: false}
	} else {
		metric.Status.Value = value
		metric.Status.RawValue = cupSize(rawValue)
		metric.Status.LastUpdated = metav1.Time{Time: time.Now()}
	}

	if err := r.Client.Status().Update(ctx, metric); err != nil {
		r.Log.Error(err, "Failed to update the Metric status")
		return ctrl.Result{}, err
	}

	return reconcile, err
}

func (r *KeptnMetricReconciler) getResults(ctx context.Context, metric *metricsapi.KeptnMetric, provider providers.KeptnSLIProvider, metricProvider *metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	if metric.Spec.Range != nil && metric.Spec.Range.Step != "" {
		return r.getStepQueryResults(ctx, metric, provider, metricProvider)
	}
	return r.getQueryResults(ctx, metric, provider, metricProvider)
}
func (r *KeptnMetricReconciler) getQueryResults(ctx context.Context, metric *metricsapi.KeptnMetric, provider providers.KeptnSLIProvider, metricProvider *metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	value, rawValue, err := provider.EvaluateQuery(ctx, *metric, *metricProvider)
	if err != nil {
		r.Log.Error(err, "Failed to evaluate the query", "Response from provider was:", (string)(rawValue))
		return "", cupSize(rawValue), err
	}
	return value, cupSize(rawValue), nil
}
func (r *KeptnMetricReconciler) getStepQueryResults(ctx context.Context, metric *metricsapi.KeptnMetric, provider providers.KeptnSLIProvider, metricProvider *metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	value, rawValue, err := provider.EvaluateQueryForStep(ctx, *metric, *metricProvider)
	if err != nil {
		r.Log.Error(err, "Failed to evaluate the query", "Response from provider was:", (string)(rawValue))
		return "", cupSize(rawValue), err
	}
	aggValue, err := aggregateValues(value, metric.Spec.Range.Aggregation)
	if err != nil {
		return "", nil, err
	}
	return aggValue, cupSize(rawValue), nil
}

func cupSize(value []byte) []byte {
	if len(value) == 0 {
		return []byte{}
	}
	if len(value) > MB {
		return value[:MB]
	}
	return value
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnMetricReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metricsapi.KeptnMetric{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *KeptnMetricReconciler) fetchProvider(ctx context.Context, namespacedMetric types.NamespacedName) (*metricsapi.KeptnMetricsProvider, error) {
	provider := &metricsapi.KeptnMetricsProvider{}
	if err := r.Client.Get(ctx, namespacedMetric, provider); err != nil {
		return nil, err
	}
	return provider, nil
}

func aggregateValues(stringSlice []string, aggFunc string) (string, error) {
	floatSlice, err := stringSliceToFloatSlice(stringSlice)
	if err != nil {
		return "", err
	}
	var aggValue float64
	switch aggFunc {
	case "max":
		aggValue = aggregation.CalculateMax(floatSlice)
	case "min":
		aggValue = aggregation.CalculateMin(floatSlice)
	case "median":
		aggValue = aggregation.CalculateMedian(floatSlice)
	case "avg":
		aggValue = aggregation.CalculateAverage(floatSlice)
	case "p90":
		aggValue = aggregation.CalculatePercentile(sort.Float64Slice(floatSlice), 90)
	case "p95":
		aggValue = aggregation.CalculatePercentile(sort.Float64Slice(floatSlice), 95)
	case "p99":
		aggValue = aggregation.CalculatePercentile(sort.Float64Slice(floatSlice), 99)
	}
	return fmt.Sprintf("%v", aggValue), nil
}

func stringSliceToFloatSlice(strSlice []string) ([]float64, error) {
	floatSlice := make([]float64, len(strSlice))

	for i, str := range strSlice {
		floatValue, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		floatSlice[i] = floatValue
	}

	return floatSlice, nil
}
