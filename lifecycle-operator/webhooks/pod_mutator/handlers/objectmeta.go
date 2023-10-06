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

func getWorkloadName(meta *metav1.ObjectMeta, applicationName string) string {
	workloadName, _ := GetLabelOrAnnotation(meta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	return operatorcommon.CreateResourceName(apicommon.MaxK8sObjectLength, apicommon.MinKeptnNameLen, applicationName, workloadName)
}

func getAppName(meta *metav1.ObjectMeta) string {
	var applicationName string
	if !isAppAnnotationPresent(meta) {
		applicationName, _ = GetLabelOrAnnotation(meta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	} else {
		applicationName, _ = GetLabelOrAnnotation(meta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)
	}
	return strings.ToLower(applicationName)
}

func getAnnotations(objMeta *metav1.ObjectMeta, annotationKey string) []string {
	if annotations, found := GetLabelOrAnnotation(objMeta, annotationKey, ""); found {
		return strings.Split(annotations, ",")
	}
	return nil
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

func isAppAnnotationPresent(meta *metav1.ObjectMeta) bool {
	_, gotAppAnnotation := GetLabelOrAnnotation(meta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)
	return gotAppAnnotation
}

func initEmptyAnnotations(meta *metav1.ObjectMeta) {
	if len(meta.Annotations) == 0 {
		meta.Annotations = make(map[string]string)
	}
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
