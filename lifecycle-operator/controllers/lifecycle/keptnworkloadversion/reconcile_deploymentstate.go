package keptnworkloadversion

import (
	"context"
	"time"

	argov1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *KeptnWorkloadVersionReconciler) reconcileDeployment(ctx context.Context, workloadVersion *apilifecycle.KeptnWorkloadVersion) (apicommon.KeptnState, error) {
	var isRunning bool
	var err error

	if r.isDeploymentTimedOut(workloadVersion) {
		workloadVersion.Status.DeploymentStatus = apicommon.StateFailed
		err = r.Client.Status().Update(ctx, workloadVersion)
		if err != nil {
			return apicommon.StateUnknown, err
		}
		r.EventSender.Emit(apicommon.PhaseWorkloadDeployment, "Warning", workloadVersion, apicommon.PhaseStateFinished, "has reached timeout", workloadVersion.GetVersion())
		return workloadVersion.Status.DeploymentStatus, nil
	}

	switch workloadVersion.Spec.ResourceReference.Kind {
	case "ReplicaSet":
		isRunning, err = r.isReplicaSetRunning(ctx, workloadVersion.Spec.ResourceReference, workloadVersion.Namespace)
	case "StatefulSet":
		isRunning, err = r.isStatefulSetRunning(ctx, workloadVersion.Spec.ResourceReference, workloadVersion.Namespace)
	case "DaemonSet":
		isRunning, err = r.isDaemonSetRunning(ctx, workloadVersion.Spec.ResourceReference, workloadVersion.Namespace)
	default:
		isRunning, err = false, controllererrors.ErrUnsupportedWorkloadVersionResourceReference
	}

	if err != nil {
		return apicommon.StateUnknown, err
	}

	if !workloadVersion.IsDeploymentStartTimeSet() {
		workloadVersion.SetDeploymentStartTime()
		workloadVersion.Status.DeploymentStatus = apicommon.StateProgressing
	}

	if isRunning {
		workloadVersion.Status.DeploymentStatus = apicommon.StateSucceeded
	}

	err = r.Client.Status().Update(ctx, workloadVersion)
	if err != nil {
		return apicommon.StateUnknown, err
	}
	return workloadVersion.Status.DeploymentStatus, nil
}

func (r *KeptnWorkloadVersionReconciler) isDeploymentTimedOut(workloadVersion *apilifecycle.KeptnWorkloadVersion) bool {
	if !workloadVersion.IsDeploymentStartTimeSet() {
		return false
	}

	deploymentDeadline := workloadVersion.Status.DeploymentStartTime.Add(r.Config.GetObservabilityTimeout().Duration)
	currentTime := time.Now().UTC()
	return currentTime.After(deploymentDeadline)
}

func (r *KeptnWorkloadVersionReconciler) isReplicaSetRunning(ctx context.Context, resource apilifecycle.ResourceReference, namespace string) (bool, error) {
	rep := appsv1.ReplicaSet{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: resource.Name, Namespace: namespace}, &rep)
	if err != nil {
		return false, err
	}

	for _, ownerRef := range rep.OwnerReferences {
		if ownerRef.Kind == "Rollout" {
			return r.isRolloutRunning(ctx, apilifecycle.ResourceReference{Name: ownerRef.Name, UID: ownerRef.UID}, namespace)
		}
	}

	return *rep.Spec.Replicas == rep.Status.AvailableReplicas, nil
}

func (r *KeptnWorkloadVersionReconciler) isDaemonSetRunning(ctx context.Context, resource apilifecycle.ResourceReference, namespace string) (bool, error) {
	daemonSet := &appsv1.DaemonSet{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: resource.Name, Namespace: namespace}, daemonSet)
	if err != nil {
		return false, err
	}
	return daemonSet.Status.DesiredNumberScheduled == daemonSet.Status.NumberReady, nil
}

func (r *KeptnWorkloadVersionReconciler) isStatefulSetRunning(ctx context.Context, resource apilifecycle.ResourceReference, namespace string) (bool, error) {
	sts := appsv1.StatefulSet{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: resource.Name, Namespace: namespace}, &sts)
	if err != nil {
		return false, err
	}
	return *sts.Spec.Replicas == sts.Status.AvailableReplicas, nil
}

func (r *KeptnWorkloadVersionReconciler) isRolloutRunning(ctx context.Context, resource apilifecycle.ResourceReference, namespace string) (bool, error) {
	rollout := argov1alpha1.Rollout{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: resource.Name, Namespace: namespace}, &rollout)
	if err != nil {
		return false, err
	}
	return rollout.Status.Replicas == rollout.Status.UpdatedReplicas && rollout.Status.Phase == argov1alpha1.RolloutPhaseHealthy, nil
}
