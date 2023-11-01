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

package keptnworkload

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	klcv1alpha4 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha4"
	operatorcommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// KeptnWorkloadReconciler reconciles a KeptnWorkload object
type KeptnWorkloadReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	EventSender   controllercommon.IEvent
	Log           logr.Logger
	TracerFactory telemetry.TracerFactory
}

// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloads,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloads/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloads/finalizers,verbs=update
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadversions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadversions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadversions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KeptnWorkload object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *KeptnWorkloadReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	requestInfo := controllercommon.GetRequestInfo(req)
	r.Log.Info("Searching for workload", "requestInfo", requestInfo)

	workload := &klcv1alpha3.KeptnWorkload{}
	err := r.Get(ctx, req.NamespacedName, workload)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}
	if err != nil {
		return reconcile.Result{}, fmt.Errorf(controllererrors.ErrCannotRetrieveWorkloadMsg, err)
	}

	traceContextCarrier := propagation.MapCarrier(workload.Annotations)
	ctx = otel.GetTextMapPropagator().Extract(ctx, traceContextCarrier)

	r.Log.Info("Reconciling Keptn Workload", "workload", workload.Name, "requestInfo", requestInfo)

	workloadVersion := &klcv1alpha4.KeptnWorkloadVersion{}
	workloadVersionName := operatorcommon.CreateResourceName(common.MaxK8sObjectLength, common.MinKeptnNameLen, workload.Name, workload.Spec.Version)

	// Try to find the workload instance
	err = r.Get(ctx, types.NamespacedName{Namespace: workload.Namespace, Name: workloadVersionName}, workloadVersion)
	// If the workload instance does not exist, create it
	if errors.IsNotFound(err) {
		workloadVersion, err := r.createWorkloadVersion(ctx, workload)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = r.Client.Create(ctx, workloadVersion)
		if err != nil {
			r.Log.Error(err, "could not create WorkloadVersion")
			r.EventSender.Emit(apicommon.PhaseCreateWorkloadVersion, "Warning", workloadVersion, apicommon.PhaseStateFailed, "could not create KeptnWorkloadVersion ", workloadVersion.Spec.Version)
			return ctrl.Result{}, err
		}
		workload.Status.CurrentVersion = workload.Spec.Version
		if err := r.Client.Status().Update(ctx, workload); err != nil {
			r.Log.Error(err, "could not update Current Version of Workload")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		r.Log.Error(err, "could not get WorkloadVersion")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnWorkloadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &klcv1alpha3.KeptnWorkload{}, "spec.app", func(rawObj client.Object) []string {
		workload := rawObj.(*klcv1alpha3.KeptnWorkload)
		return []string{workload.Spec.AppName}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&klcv1alpha3.KeptnWorkload{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *KeptnWorkloadReconciler) createWorkloadVersion(ctx context.Context, workload *klcv1alpha3.KeptnWorkload) (*klcv1alpha4.KeptnWorkloadVersion, error) {
	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	previousVersion := ""
	if workload.Spec.Version != workload.Status.CurrentVersion {
		previousVersion = workload.Status.CurrentVersion
	}

	workloadVersion := generateWorkloadVersion(previousVersion, traceContextCarrier, workload)
	err := controllerutil.SetControllerReference(workload, &workloadVersion, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference for WorkloadVersion: "+workloadVersion.Name)
	}

	return &workloadVersion, err
}

func generateWorkloadVersion(previousVersion string, traceContextCarrier map[string]string, w *klcv1alpha3.KeptnWorkload) klcv1alpha4.KeptnWorkloadVersion {
	return klcv1alpha4.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: traceContextCarrier,
			Name:        operatorcommon.CreateResourceName(common.MaxK8sObjectLength, common.MinKeptnNameLen, w.Name, w.Spec.Version),
			Namespace:   w.Namespace,
		},
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: w.Spec,
			WorkloadName:      w.Name,
			PreviousVersion:   previousVersion,
		},
	}
}
