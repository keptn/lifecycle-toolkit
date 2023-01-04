/*
Copyright 2023.

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

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"time"

	metricsv1alpha1 "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha1"
	providers "github.com/keptn/lifecycle-toolkit/metrics-operator/pkg/metrics/providers"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MetricReconciler reconciles a Metric object
type MetricReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger

	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=providers,verbs=get;list;watch
//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=metrics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=metrics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=metrics/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Metric object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *MetricReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("Reconciling Metric")
	metric := &metricsv1alpha1.Metric{}

	if err := r.Client.Get(ctx, req.NamespacedName, metric); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			r.Log.Info("Metric resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get the Metric")
		return ctrl.Result{}, nil
	}

	if time.Now().Before(metric.Status.LastUpdated.Add(metric.Spec.FetchIntervalSeconds * time.Second)) {
		r.Log.Info("Metric has not been updated for the configured interval. Skipping")
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}

	evaluationProvider, err := r.fetchProvider(ctx, types.NamespacedName{Name: metric.Spec.Source, Namespace: metric.Namespace})
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info(err.Error() + ", ignoring error since object must be deleted")
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
		}
		r.Log.Error(err, "Failed to retrieve the provider")
		return ctrl.Result{}, nil
	}
	// load the provider
	provider, err2 := providers.NewProvider(metric.Spec.Source, r.Log, r.Client)
	if err2 != nil {
		r.recordEvent("Error", metric, "ProviderNotFound", "provider was not found")
		r.Log.Error(err2, "Failed to get the correct Metric Provider")
		return ctrl.Result{Requeue: false}, err2
	}

	value, err := provider.EvaluateQuery(ctx, *metric, *evaluationProvider)
	if err != nil {
		r.recordEvent("Error", metric, "EvaluationFailed", "evaluation failed")
		r.Log.Error(err, "Failed to evaluate the query")
		return ctrl.Result{Requeue: false}, err
	}
	metric.Status.Value = value
	metric.Status.LastUpdated = metav1.Time{Time: time.Now()}

	if err := r.Client.Status().Update(ctx, metric); err != nil {
		r.Log.Error(err, "Failed to update the Metric status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MetricReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metricsv1alpha1.Metric{}).
		Complete(r)
}

func (r *MetricReconciler) fetchProvider(ctx context.Context, namespacedMetric types.NamespacedName) (*metricsv1alpha1.Provider, error) {
	provider := &metricsv1alpha1.Provider{}
	if err := r.Client.Get(ctx, namespacedMetric, provider); err != nil {
		return nil, err
	}
	return provider, nil
}

func (r *MetricReconciler) recordEvent(eventType string, metric *metricsv1alpha1.Metric, shortReason string, longReason string) {
	r.Recorder.Event(metric, eventType, shortReason, fmt.Sprintf("%s / Namespace: %s, Name: %s, WorkloadVersion: %s ", longReason, metric.Namespace, metric.Name, metric.Namespace))
}
