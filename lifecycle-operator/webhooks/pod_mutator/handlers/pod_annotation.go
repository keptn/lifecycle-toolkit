package handlers

import (
	"context"

	argov1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"github.com/go-logr/logr"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type PodAnnotationHandler struct {
	Client client.Client
	Log    logr.Logger
}

func (p *PodAnnotationHandler) IsAnnotated(ctx context.Context, req admission.Request, pod *corev1.Pod) bool {
	podIsAnnotated := isPodAnnotated(pod)
	if !podIsAnnotated {
		p.Log.Info("Pod is not annotated, check for parent annotations...")
		podIsAnnotated = p.copyAnnotationsIfParentAnnotated(ctx, &req, pod)
	}
	return podIsAnnotated
}

func (p *PodAnnotationHandler) copyAnnotationsIfParentAnnotated(ctx context.Context, req *admission.Request, pod *corev1.Pod) bool {
	podOwner := GetOwnerReference(&pod.ObjectMeta)
	if podOwner.UID == "" {
		return false
	}

	switch podOwner.Kind {
	case "ReplicaSet":
		rs := &appsv1.ReplicaSet{}
		if err := p.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: podOwner.Name}, rs); err != nil {
			return false
		}
		p.Log.Info("Done fetching RS")

		rsOwner := GetOwnerReference(&rs.ObjectMeta)
		if rsOwner.UID == "" {
			return false
		}

		if rsOwner.Kind == "Rollout" {
			ro := &argov1alpha1.Rollout{}
			return p.fetchParentObjectAndCopyLabels(ctx, podOwner.Name, req.Namespace, pod, ro)
		}
		dp := &appsv1.Deployment{}
		return p.fetchParentObjectAndCopyLabels(ctx, rsOwner.Name, req.Namespace, pod, dp)

	case "StatefulSet":
		sts := &appsv1.StatefulSet{}
		return p.fetchParentObjectAndCopyLabels(ctx, podOwner.Name, req.Namespace, pod, sts)
	case "DaemonSet":
		ds := &appsv1.DaemonSet{}
		return p.fetchParentObjectAndCopyLabels(ctx, podOwner.Name, req.Namespace, pod, ds)
	default:
		return false
	}
}

func (p *PodAnnotationHandler) fetchParentObjectAndCopyLabels(ctx context.Context, name string, namespace string, pod *corev1.Pod, objectContainer client.Object) bool {
	if err := p.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, objectContainer); err != nil {
		return false
	}
	objectContainerMetaData := metav1.ObjectMeta{
		Labels:      objectContainer.GetLabels(),
		Annotations: objectContainer.GetAnnotations(),
	}
	return copyResourceLabelsIfPresent(&objectContainerMetaData, pod)
}

func copyResourceLabelsIfPresent(sourceResource *metav1.ObjectMeta, targetPod *corev1.Pod) bool {
	var workloadName, appName, version, preDeploymentChecks, postDeploymentChecks, preEvaluationChecks, postEvaluationChecks string
	var gotWorkloadName, gotVersion bool

	workloadName, gotWorkloadName = GetLabelOrAnnotation(sourceResource, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	appName, _ = GetLabelOrAnnotation(sourceResource, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)
	version, gotVersion = GetLabelOrAnnotation(sourceResource, apicommon.VersionAnnotation, apicommon.K8sRecommendedVersionAnnotations)
	preDeploymentChecks, _ = GetLabelOrAnnotation(sourceResource, apicommon.PreDeploymentTaskAnnotation, "")
	postDeploymentChecks, _ = GetLabelOrAnnotation(sourceResource, apicommon.PostDeploymentTaskAnnotation, "")
	preEvaluationChecks, _ = GetLabelOrAnnotation(sourceResource, apicommon.PreDeploymentEvaluationAnnotation, "")
	postEvaluationChecks, _ = GetLabelOrAnnotation(sourceResource, apicommon.PostDeploymentEvaluationAnnotation, "")

	if len(targetPod.Annotations) == 0 {
		targetPod.Annotations = make(map[string]string)
	}

	if gotWorkloadName {
		setMapKey(targetPod.Annotations, apicommon.WorkloadAnnotation, workloadName)

		if !gotVersion {
			setMapKey(targetPod.Annotations, apicommon.VersionAnnotation, calculateVersion(targetPod))
		} else {
			setMapKey(targetPod.Annotations, apicommon.VersionAnnotation, version)
		}

		setMapKey(targetPod.Annotations, apicommon.AppAnnotation, appName)
		setMapKey(targetPod.Annotations, apicommon.PreDeploymentTaskAnnotation, preDeploymentChecks)
		setMapKey(targetPod.Annotations, apicommon.PostDeploymentTaskAnnotation, postDeploymentChecks)
		setMapKey(targetPod.Annotations, apicommon.PreDeploymentEvaluationAnnotation, preEvaluationChecks)
		setMapKey(targetPod.Annotations, apicommon.PostDeploymentEvaluationAnnotation, postEvaluationChecks)

		return true
	}
	return false
}

func isPodAnnotated(pod *corev1.Pod) bool {
	_, gotWorkloadAnnotation := GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	_, gotVersionAnnotation := GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.VersionAnnotation, apicommon.K8sRecommendedVersionAnnotations)

	if gotWorkloadAnnotation {
		if !gotVersionAnnotation {
			if len(pod.Annotations) == 0 {
				pod.Annotations = make(map[string]string)
			}
			pod.Annotations[apicommon.VersionAnnotation] = calculateVersion(pod)
		}
		return true
	}
	return false
}
