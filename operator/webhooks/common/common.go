package common

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"

	argov1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	InvalidAnnotationMessage = "Invalid annotations"

	// SecretCertsName is the name of the secret where the webhook certificates are stored.
	SecretCertsName = "klt-certs"
)

var ErrTooLongAnnotations = fmt.Errorf("too long annotations, maximum length for app and workload is 25 characters, for version 12 characters")

func IsPodOrParentAnnotated(ctx context.Context, req *admission.Request, pod *corev1.Pod, k8sclient client.Client) (bool, error) {
	podIsAnnotated, err := isPodAnnotated(pod)

	if err != nil {
		return false, err
	}

	if !podIsAnnotated {
		podIsAnnotated, err = copyAnnotationsIfParentAnnotated(ctx, req, pod, k8sclient)
	}
	return podIsAnnotated, err
}

func isPodAnnotated(pod *corev1.Pod) (bool, error) {
	workload, gotWorkloadAnnotation := GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	version, gotVersionAnnotation := GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.VersionAnnotation, apicommon.K8sRecommendedVersionAnnotations)

	if len(workload) > apicommon.MaxWorkloadNameLength || len(version) > apicommon.MaxVersionLength {
		return false, ErrTooLongAnnotations
	}

	if gotWorkloadAnnotation {
		if !gotVersionAnnotation {
			if len(pod.Annotations) == 0 {
				pod.Annotations = make(map[string]string)
			}
			pod.Annotations[apicommon.VersionAnnotation] = calculateVersion(pod)
		}
		return true, nil
	}
	return false, nil
}

func GetLabelOrAnnotation(resource *metav1.ObjectMeta, primaryAnnotation string, secondaryAnnotation string) (string, bool) {
	if resource.Annotations[primaryAnnotation] != "" {
		return resource.Annotations[primaryAnnotation], true
	}

	if resource.Labels[primaryAnnotation] != "" {
		return resource.Labels[primaryAnnotation], true
	}

	if secondaryAnnotation == "" {
		return "", false
	}

	if resource.Annotations[secondaryAnnotation] != "" {
		return resource.Annotations[secondaryAnnotation], true
	}

	if resource.Labels[secondaryAnnotation] != "" {
		return resource.Labels[secondaryAnnotation], true
	}
	return "", false
}

func copyAnnotationsIfParentAnnotated(ctx context.Context, req *admission.Request, pod *corev1.Pod, k8sclient client.Client) (bool, error) {
	podOwner := GetOwnerReference(&pod.ObjectMeta)
	if podOwner.UID == "" {
		return false, nil
	}

	switch podOwner.Kind {
	case "ReplicaSet":
		rs := &appsv1.ReplicaSet{}
		if err := k8sclient.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: podOwner.Name}, rs); err != nil {
			return false, nil
		}

		rsOwner := GetOwnerReference(&rs.ObjectMeta)
		if rsOwner.UID == "" {
			return false, nil
		}

		if rsOwner.Kind == "Rollout" {
			ro := &argov1alpha1.Rollout{}
			return fetchParentObjectAndCopyLabels(ctx, podOwner.Name, req.Namespace, pod, ro, k8sclient)
		}
		dp := &appsv1.Deployment{}
		return fetchParentObjectAndCopyLabels(ctx, rsOwner.Name, req.Namespace, pod, dp, k8sclient)

	case "StatefulSet":
		sts := &appsv1.StatefulSet{}
		return fetchParentObjectAndCopyLabels(ctx, podOwner.Name, req.Namespace, pod, sts, k8sclient)
	case "DaemonSet":
		ds := &appsv1.DaemonSet{}
		return fetchParentObjectAndCopyLabels(ctx, podOwner.Name, req.Namespace, pod, ds, k8sclient)
	default:
		return false, nil
	}
}

func fetchParentObjectAndCopyLabels(ctx context.Context, name string, namespace string, pod *corev1.Pod, objectContainer client.Object, k8sclient client.Client) (bool, error) {
	if err := k8sclient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, objectContainer); err != nil {
		return false, nil
	}
	objectContainerMetaData := metav1.ObjectMeta{
		Labels:      objectContainer.GetLabels(),
		Annotations: objectContainer.GetAnnotations(),
	}
	return copyResourceLabelsIfPresent(&objectContainerMetaData, pod)
}

func copyResourceLabelsIfPresent(sourceResource *metav1.ObjectMeta, targetPod *corev1.Pod) (bool, error) {
	var workloadName, appName, version, preDeploymentChecks, postDeploymentChecks, preEvaluationChecks, postEvaluationChecks string
	var gotWorkloadName, gotVersion bool

	workloadName, gotWorkloadName = GetLabelOrAnnotation(sourceResource, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	appName, _ = GetLabelOrAnnotation(sourceResource, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)
	version, gotVersion = GetLabelOrAnnotation(sourceResource, apicommon.VersionAnnotation, apicommon.K8sRecommendedVersionAnnotations)
	preDeploymentChecks, _ = GetLabelOrAnnotation(sourceResource, apicommon.PreDeploymentTaskAnnotation, "")
	postDeploymentChecks, _ = GetLabelOrAnnotation(sourceResource, apicommon.PostDeploymentTaskAnnotation, "")
	preEvaluationChecks, _ = GetLabelOrAnnotation(sourceResource, apicommon.PreDeploymentEvaluationAnnotation, "")
	postEvaluationChecks, _ = GetLabelOrAnnotation(sourceResource, apicommon.PostDeploymentEvaluationAnnotation, "")

	if len(workloadName) > apicommon.MaxWorkloadNameLength || len(version) > apicommon.MaxVersionLength {
		return false, ErrTooLongAnnotations
	}

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

		return true, nil
	}
	return false, nil
}

func GetOwnerReference(resource *metav1.ObjectMeta) metav1.OwnerReference {
	reference := metav1.OwnerReference{}
	if len(resource.OwnerReferences) != 0 {
		for _, owner := range resource.OwnerReferences {
			if owner.Kind == "ReplicaSet" || owner.Kind == "Deployment" || owner.Kind == "StatefulSet" || owner.Kind == "DaemonSet" || owner.Kind == "Rollout" {
				reference.UID = owner.UID
				reference.Kind = owner.Kind
				reference.Name = owner.Name
				reference.APIVersion = owner.APIVersion
			}
		}
	}
	return reference
}

func calculateVersion(pod *corev1.Pod) string {
	name := ""

	if len(pod.Spec.Containers) == 1 {
		image := strings.Split(pod.Spec.Containers[0].Image, ":")
		if len(image) > 1 && image[1] != "" && image[1] != "latest" {
			return image[1]
		}
	}

	for _, item := range pod.Spec.Containers {
		name = name + item.Name + item.Image
		for _, e := range item.Env {
			name = name + e.Name + e.Value
		}
	}

	h := fnv.New32a()
	h.Write([]byte(name))
	return fmt.Sprint(h.Sum32())
}

func setMapKey(myMap map[string]string, key, value string) {
	if myMap == nil {
		return
	}
	if value != "" {
		myMap[key] = value
	}
}
