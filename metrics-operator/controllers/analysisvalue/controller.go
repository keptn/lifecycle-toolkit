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

package analysisvalue

import (
	"context"
	"strings"
	"time"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AnalysisValueReconciler reconciles a AnalysisValue object
type AnalysisValueReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

const MB = 1 << (10 * 2)

//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=analysisvalues,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=analysisvalues/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=analysisvalues/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AnalysisValue object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *AnalysisValueReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("Reconciling AnalysisValue")
	analysisValue := &metricsapi.AnalysisValue{}

	if err := r.Client.Get(ctx, req.NamespacedName, analysisValue); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			r.Log.Info("AnalysisValue resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get the AnalysisValue")
		return ctrl.Result{}, nil
	}

	analysisTemplate := &metricsapi.AnalysisTemplate{}
	if err := r.Client.Get(ctx, types.NamespacedName{Name: analysisValue.Spec.AnalysisTemplate.Name, Namespace: req.Namespace}, analysisTemplate); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			r.Log.Info("AnalysisTemplate resource not found. Ignoring since object must be deleted")
			return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		}
		r.Log.Error(err, "Failed to get the AnalysisTemplate")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	analysisValue.Status.Query = generateQuery(analysisTemplate.Spec.Query, analysisValue.Spec.Selectors)

	metricProvider, err := r.fetchProvider(ctx, types.NamespacedName{Name: analysisTemplate.Spec.ProviderRef.Name, Namespace: req.Namespace})
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info(err.Error() + ", ignoring error since object must be deleted")
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, err
		}
		r.Log.Error(err, "Failed to retrieve the provider")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, err
	}

	// load the provider
	provider, err := providers.NewProvider(metricProvider.GetType(), r.Log, r.Client)
	if err != nil {
		r.Log.Error(err, "Failed to get the correct Metric Provider")
		return ctrl.Result{Requeue: false}, err
	}

	value, rawValue, err := provider.EvaluateQuery(ctx, *analysisValue, *metricProvider)
	if err != nil {
		r.Log.Error(err, "Failed to evaluate the query", "Response from provider was:", (string)(rawValue))
		analysisValue.Status.ErrMsg = err.Error()
		analysisValue.Status.Value = ""
		analysisValue.Status.RawValue = cupSize(rawValue)
	} else {
		analysisValue.Status.Value = value
		analysisValue.Status.RawValue = cupSize(rawValue)
	}

	if err := r.Client.Status().Update(ctx, analysisValue); err != nil {
		r.Log.Error(err, "Failed to update the AnalysisValue status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.RequeueAfter: 10 * time.Second
func (r *AnalysisValueReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metricsapi.AnalysisValue{}).
		Complete(r)
}

func generateQuery(query string, selectors map[string]string) string {
	for key, value := range selectors {
		query = strings.Replace(query, "$"+strings.ToUpper(key), value, -1)
	}
	return query
}

func (r *AnalysisValueReconciler) fetchProvider(ctx context.Context, namespacedMetric types.NamespacedName) (*metricsapi.KeptnMetricsProvider, error) {
	provider := &metricsapi.KeptnMetricsProvider{}
	if err := r.Client.Get(ctx, namespacedMetric, provider); err != nil {
		return nil, err
	}
	return provider, nil
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
