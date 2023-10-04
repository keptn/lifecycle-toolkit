package handlers

import (
	"fmt"
	"hash/fnv"
	"strings"

	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	operatorcommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func getWorkloadName(meta *metav1.ObjectMeta) string {
	workloadName, _ := GetLabelOrAnnotation(meta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	applicationName, _ := GetLabelOrAnnotation(meta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)
	return operatorcommon.CreateResourceName(apicommon.MaxK8sObjectLength, apicommon.MinKeptnNameLen, applicationName, workloadName)
}

func getAppName(meta *metav1.ObjectMeta) string {
	applicationName, _ := GetLabelOrAnnotation(meta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)
	return strings.ToLower(applicationName)
}

func getVersion(meta *metav1.ObjectMeta) string {
	version, _ := GetLabelOrAnnotation(meta, apicommon.VersionAnnotation, apicommon.K8sRecommendedVersionAnnotations)
	return strings.ToLower(version)
}

func GetOwnerReference(resource *metav1.ObjectMeta) metav1.OwnerReference {
	reference := metav1.OwnerReference{}
	if len(resource.OwnerReferences) != 0 {
		for _, owner := range resource.OwnerReferences {
			if apicommon.IsOwnerSupported(owner) {
				reference.UID = owner.UID
				reference.Kind = owner.Kind
				reference.Name = owner.Name
				reference.APIVersion = owner.APIVersion
			}
		}
	}
	return reference
}

func setMapKey(myMap map[string]string, key, value string) {
	if myMap == nil {
		return
	}
	if value != "" {
		myMap[key] = value
	}
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

func isAppAnnotationPresent(pod *corev1.Pod) bool {
	_, gotAppAnnotation := GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)

	if gotAppAnnotation {
		return true
	}

	if len(pod.Annotations) == 0 {
		pod.Annotations = make(map[string]string)
	}
	pod.Annotations[apicommon.AppAnnotation], _ = GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	return false
}

func calculateVersion(pod *corev1.Pod) string {
	name := ""

	if len(pod.Spec.Containers) == 1 {
		image := strings.Split(pod.Spec.Containers[0].Image, ":")
		lenImg := len(image) - 1
		if lenImg >= 1 && image[lenImg] != "" && image[lenImg] != "latest" {
			return image[lenImg]
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
