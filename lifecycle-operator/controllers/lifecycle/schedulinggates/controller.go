package schedulinggates

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// SchedulingGatesReconciler reconciles a KeptnWorkloadVersion object
type SchedulingGatesReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadversions,verbs=get;list;watch;
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;update

func (r *SchedulingGatesReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	requestInfo := controllercommon.GetRequestInfo(req)
	r.Log.Info("Searching for pod", "requestInfo", requestInfo)

	pod := &v1.Pod{}

	err := r.Get(ctx, req.NamespacedName, pod)
	if errors.IsNotFound(err) {
		return ctrl.Result{}, nil
	}

	if err != nil {
		r.Log.Error(err, "Could not retrieve pod", "requestInfo", requestInfo)
		return ctrl.Result{}, fmt.Errorf("could not retrieve pod, %w", err)
	}

	// check if the owner of the pod is the one that the KeptnWorkloadVersion is referring to
	owner := pod.GetOwnerReferences()
	if len(owner) == 0 {
		return ctrl.Result{}, nil
	}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(".spec.resourceReference.uid", string(owner[0].UID)),
		Namespace:     pod.GetNamespace(),
	}

	attachedWorkloadVersions := &apilifecycle.KeptnWorkloadVersionList{}

	if err := r.List(ctx, attachedWorkloadVersions, listOps); err != nil {
		r.Log.Error(err, "Could not list WorkloadVersions related to pod", "pod", pod.GetName(), "namespace", pod.GetNamespace())
		return ctrl.Result{}, err
	}

	for _, workloadVersion := range attachedWorkloadVersions.Items {
		if workloadVersion.Status.DeploymentStatus.IsCompleted() || workloadVersion.Status.DeploymentStatus == apicommon.StateProgressing {
			return r.removeGate(ctx, pod)
		}
	}
	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil

}

func (r *SchedulingGatesReconciler) removeGate(ctx context.Context, pod *v1.Pod) (ctrl.Result, error) {
	pod.Spec.SchedulingGates = nil
	if len(pod.Annotations) == 0 {
		pod.Annotations = make(map[string]string, 1)
	}
	pod.Annotations[apicommon.SchedulingGateRemoved] = "true"
	r.Log.Info("removing scheduling gate of pod", "pod", pod.Name, "uid", pod.UID)

	if err := r.Update(ctx, pod); err != nil {
		r.Log.Error(err, "Could not remove pod scheduling gate", "namespace", pod.Namespace, "pod", pod.Name)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SchedulingGatesReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(
			&v1.Pod{},
			builder.WithPredicates(
				predicate.NewPredicateFuncs(func(object client.Object) bool {
					pod, ok := object.(*v1.Pod)
					if !ok {
						return false
					}
					return hasKeptnSchedulingGate(pod)
				}),
			),
		).
		Complete(r)
}

func hasKeptnSchedulingGate(pod *v1.Pod) bool {
	for _, gate := range pod.Spec.SchedulingGates {
		if gate.Name == apicommon.KeptnGate {
			return true
		}
	}
	return false
}
