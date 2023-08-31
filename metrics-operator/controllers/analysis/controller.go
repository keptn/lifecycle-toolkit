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

package analysis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	common "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// AnalysisReconciler reconciles an Analysis object
type AnalysisReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Log        logr.Logger
	MaxWorkers int //maybe 2 or 4 as def
	Namespace  string
	NewWorkersPoolFactory
	common.IAnalysisEvaluator
}

//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=analyses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=analyses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=analyses/finalizers,verbs=update
// +kubebuilder:rbac:groups=metrics.keptn.sh,resources=keptnmetricsproviders,verbs=get;list;watch;
//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=analysisdefinitions,verbs=get;list;watch;
//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=analysisvaluetemplates,verbs=get;list;watch;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its AnalysisResult here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (a *AnalysisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	a.Log.Info("Reconciling Analysis")
	analysis := &metricsapi.Analysis{}

	//retrieve analysis
	if err := a.Client.Get(ctx, req.NamespacedName, analysis); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			a.Log.Info("Analysis resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		a.Log.Error(err, "Failed to get the Analysis")
		return ctrl.Result{}, nil
	}

	//find AnalysisDefinition to have the collection of Objectives
	analysisDef := &metricsapi.AnalysisDefinition{}
	if analysis.Spec.AnalysisDefinition.Namespace == "" {
		analysis.Spec.AnalysisDefinition.Namespace = a.Namespace
	}
	err := a.Client.Get(ctx,
		types.NamespacedName{
			Name:      analysis.Spec.AnalysisDefinition.Name,
			Namespace: analysis.Spec.AnalysisDefinition.Namespace},
		analysisDef,
	)

	if err != nil {
		if errors.IsNotFound(err) {
			a.Log.Info(err.Error() + ", ignoring error since object must be deleted")
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
		}
		a.Log.Error(err, "Failed to retrieve the AnalysisDefinition")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	var done map[string]metricsapi.ProviderResult
	todo := analysisDef.Spec.Objectives
	if analysis.Status.Raw != "" {
		todo, done = a.ExtractMissingObj(analysisDef, analysis.Status.StoredValues)
		if len(todo) == 0 {
			return ctrl.Result{}, nil
		}
	}

	//create multiple workers handling the Objectives
	childCtx, wp := a.NewWorkersPoolFactory(ctx, analysis, todo, a.MaxWorkers, a.Client, a.Log, a.Namespace)

	res, err := wp.DispatchAndCollect(childCtx)
	if err != nil {
		a.Log.Error(err, "Failed to Collect all SLOs, caching collected values")
		analysis.Status.StoredValues = res
		err = a.updateStatus(ctx, analysis)
		return ctrl.Result{RequeueAfter: 10 * time.Second}, err
	}

	maps.Copy(res, done)

	eval := a.Evaluate(res, analysisDef)
	analysisResultJSON, err := json.Marshal(eval)
	if err != nil {
		a.Log.Error(err, "Could not marshal status")
	}
	analysis.Status.Raw = string(analysisResultJSON)
	if eval.Warning {
		analysis.Status.Warning = true
	}
	analysis.Status.Pass = eval.Pass
	err = a.updateStatus(ctx, analysis)

	return ctrl.Result{}, err
}

func (a *AnalysisReconciler) updateStatus(ctx context.Context, analysis *metricsapi.Analysis) error {
	if err := a.Client.Status().Update(ctx, analysis); err != nil {
		a.Log.Error(err, "Failed to update the Analysis status")
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Managea.
func (a *AnalysisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metricsapi.Analysis{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(a)
}

func (a AnalysisReconciler) ExtractMissingObj(def *metricsapi.AnalysisDefinition, status map[string]metricsapi.ProviderResult) ([]metricsapi.Objective, map[string]metricsapi.ProviderResult) {
	var toDo []metricsapi.Objective
	done := make(map[string]metricsapi.ProviderResult, len(status))
	for _, obj := range def.Spec.Objectives {
		if value, ok := status[common.ComputeKey(obj.AnalysisValueTemplateRef)]; ok {
			if value.ErrMsg != "" {
				toDo = append(toDo, obj)
			} else {
				done[common.ComputeKey(obj.AnalysisValueTemplateRef)] = status[common.ComputeKey(obj.AnalysisValueTemplateRef)]
			}
		}
	}
	return toDo, done
}
