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
	r.Log.Info("Reconciling Metric")
	metric := &metricsapi.KeptnMetric{}

	if err := r.Client.Get(ctx, req.NamespacedName, metric); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			r.Log.Info("Metric resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get the Metric")
		return ctrl.Result{}, nil
	}

	fetchTime := metric.Status.LastUpdated.Add(time.Second * time.Duration(metric.Spec.FetchIntervalSeconds))
	if time.Now().Before(fetchTime) {
		diff := time.Until(fetchTime)
		r.Log.Info("Metric has not been updated for the configured interval. Skipping")
		return ctrl.Result{Requeue: true, RequeueAfter: diff}, nil
	}

	metricProvider, err := r.fetchProvider(ctx, types.NamespacedName{Name: metric.Spec.Provider.Name, Namespace: metric.Namespace})
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info(err.Error() + ", ignoring error since object must be deleted")
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
		}
		r.Log.Error(err, "Failed to retrieve the provider")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}
	// load the provider
	provider, err2 := providers.NewProvider(metricProvider.GetType(), r.Log, r.Client)
	if err2 != nil {
		r.Log.Error(err2, "Failed to get the correct Metric Provider")
		return ctrl.Result{Requeue: false}, err2
	}
	reconcile := ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}

	if metric.Spec.Range.Step != "" {
		value, rawValue, err := provider.EvaluateQueryForStep(ctx, *metric, *metricProvider)
		if err != nil {
			r.Log.Error(err, "Failed to evaluate the query", "Response from provider was:", (string)(rawValue))
			metric.Status.ErrMsg = err.Error()
			metric.Status.Value = ""
			metric.Status.RawValue = cupSize(rawValue)
			metric.Status.LastUpdated = metav1.Time{Time: time.Now()}
			reconcile = ctrl.Result{Requeue: false}
		} else {
			aggValue, err := aggregateValues(value)
			if err != nil {
				return ctrl.Result{}, err
			}
			metric.Status.Value = aggValue
			metric.Status.RawValue = cupSize(rawValue)
			metric.Status.LastUpdated = metav1.Time{Time: time.Now()}
		}

		if err := r.Client.Status().Update(ctx, metric); err != nil {
			r.Log.Error(err, "Failed to update the Metric status")
			return ctrl.Result{}, err
		}
	} else {
		value, rawValue, err := provider.EvaluateQuery(ctx, *metric, *metricProvider)
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
	}

	return reconcile, err
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
	var aggValue string
	switch aggFunc {
	case "max":
		aggValue = fmt.Sprintf("%v", calculateMax(floatSlice))
	case "min":
		aggValue = fmt.Sprintf("%v", calculateMin(floatSlice))
	case "median":
		aggValue = fmt.Sprintf("%v", calculateMedian(floatSlice))
	case "avg":
		aggValue = fmt.Sprintf("%v", calculateAverage(floatSlice))
	case "p90":
		aggValue = fmt.Sprintf("%v", calculatePercentile(sort.Float64Slice(floatSlice), 0.90))
	case "p95":
		aggValue = fmt.Sprintf("%v", calculatePercentile(sort.Float64Slice(floatSlice), 0.95))
	case "p99":
		aggValue = fmt.Sprintf("%v", calculatePercentile(sort.Float64Slice(floatSlice), 0.99))
	}
	return aggValue, nil
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

func calculateMax(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	max := values[0]
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
}

func calculateMin(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	min := values[0]
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min
}

func calculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	// Sort the values
	sortedValues := make([]float64, len(values))
	copy(sortedValues, values)
	sort.Float64s(sortedValues)

	// Calculate the median
	middle := len(sortedValues) / 2
	if len(sortedValues)%2 == 0 {
		return (sortedValues[middle-1] + sortedValues[middle]) / 2
	}
	return sortedValues[middle]
}

func calculateAverage(values []float64) float64 {
	sum := 0.0

	for _, value := range values {
		sum += value
	}
	if len(values) > 0 {
		return sum / float64(len(values))
	}

	return 0.0
}

func calculatePercentile(values sort.Float64Slice, perc float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	// Calculate the index for the requested percentile
	i := int(float64(len(values)) * perc / 100)

	return values[i]
}
