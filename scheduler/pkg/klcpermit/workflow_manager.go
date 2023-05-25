package klcpermit

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/keptn/lifecycle-toolkit/scheduler/pkg/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
)

var workloadInstanceResource = schema.GroupVersionResource{Group: "lifecycle.keptn.sh", Version: "v1alpha3", Resource: "keptnworkloadinstances"}

type Status string

const (
	WorkloadInstanceStatusNotSpecified Status = "Workload run status not specified"
	WorkloadInstanceNotFound           Status = "Workload run not found"
	Success                            Status = "Success"
	Failure                            Status = "Failure"
	Wait                               Status = "Wait"
)

type KeptnState string

const (
	StateProgressing KeptnState = "Progressing"
	StateSucceeded   KeptnState = "Succeeded"
	StateFailed      KeptnState = "Failed"
	StateUnknown     KeptnState = "Unknown"
	StatePending     KeptnState = "Pending"
	StateDeprecated  KeptnState = "Deprecated"
)

const WorkloadAnnotation = "keptn.sh/workload"
const VersionAnnotation = "keptn.sh/version"
const AppAnnotation = "keptn.sh/app"
const K8sRecommendedWorkloadAnnotations = "app.kubernetes.io/name"
const K8sRecommendedVersionAnnotations = "app.kubernetes.io/version"
const K8sRecommendedAppAnnotations = "app.kubernetes.io/part-of"

type Manager interface {
	Permit(context.Context, *corev1.Pod) Status
}

type WorkloadManager struct {
	dynamicClient dynamic.Interface
	Tracer        trace.Tracer
	bindCRDSpan   map[string]trace.Span
}

func NewWorkloadManager(d dynamic.Interface) *WorkloadManager {
	sMgr := &WorkloadManager{
		dynamicClient: d,
		Tracer:        otel.Tracer("keptn/scheduler"),
		bindCRDSpan:   make(map[string]trace.Span, 100),
	}
	return sMgr
}

func (sMgr *WorkloadManager) Permit(ctx context.Context, pod *corev1.Pod) Status {
	//List workloadInstance run CRDs
	name := getCRDName(pod)
	crd, err := sMgr.GetCRD(ctx, pod.Namespace, name)

	if err != nil {
		klog.Infof("[Keptn Permit Plugin] could not find workloadInstance crd %s, err:%s", name, err.Error())
		return WorkloadInstanceNotFound
	}

	_, span := sMgr.getSpan(ctx, crd, pod)

	//check CRD status
	phase, found, err := unstructured.NestedString(crd.UnstructuredContent(), "status", "preDeploymentEvaluationStatus")
	klog.Infof("[Keptn Permit Plugin] workloadInstance crd %s, found %s with phase %s ", crd, found, phase)
	if err == nil && found {
		span.AddEvent("StatusEvaluation", trace.WithAttributes(tracing.Status.String(phase)))
		switch KeptnState(phase) {
		case StateFailed, StateDeprecated:
			span.SetStatus(codes.Error, "Failed")
			span.End()
			sMgr.unbindSpan(pod)
			return Failure
		case StateSucceeded:
			span.End()
			sMgr.unbindSpan(pod)
			return Success
		case StatePending:
			return Wait
		case StateProgressing:
			return Wait
		case StateUnknown:
			return Wait
		}
	}
	return WorkloadInstanceStatusNotSpecified
}

// GetCRD returns unstructured to avoid tight coupling with the CRD resource
func (sMgr *WorkloadManager) GetCRD(ctx context.Context, namespace string, name string) (*unstructured.Unstructured, error) {
	// GET /apis/lifecycle.keptn.sh/v1/namespaces/{namespace}/workloadinstance/name
	return sMgr.dynamicClient.Resource(workloadInstanceResource).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (sMgr *WorkloadManager) getSpan(ctx context.Context, crd *unstructured.Unstructured, pod *corev1.Pod) (context.Context, trace.Span) {
	name := getCRDName(pod)
	if span, ok := sMgr.bindCRDSpan[name]; ok {
		return ctx, span
	}
	ctx, span := tracing.CreateSpan(ctx, crd, sMgr.Tracer, pod.Namespace)
	//TODO store only sampled one and cap it
	sMgr.bindCRDSpan[name] = span
	return ctx, span
}

func createResourceName(maxLen int, minSubstrLen int, str ...string) string {
	for len(str)*minSubstrLen > maxLen {
		minSubstrLen = minSubstrLen / 2
	}
	for i := 0; i < len(str)-1; i++ {
		newStr := strings.Join(str, "-")
		if len(newStr) > maxLen {
			if len(str[i]) < minSubstrLen {
				continue
			}
			cut := len(newStr) - maxLen
			if cut > len(str[i])-minSubstrLen {
				str[i] = str[i][:minSubstrLen]
			} else {
				str[i] = str[i][:len(str[i])-cut]
			}
		} else {
			return strings.ToLower(newStr)
		}
	}

	return strings.ToLower(strings.Join(str, "-"))
}

func getCRDName(pod *corev1.Pod) string {
	application, _ := getLabelOrAnnotation(pod, AppAnnotation, K8sRecommendedAppAnnotations)
	workload, _ := getLabelOrAnnotation(pod, WorkloadAnnotation, K8sRecommendedWorkloadAnnotations)
	version, versionExists := getLabelOrAnnotation(pod, VersionAnnotation, K8sRecommendedVersionAnnotations)
	if !versionExists {
		version = calculateVersion(pod)
	}
	return application + "-" + workload + "-" + version
}

func (sMgr *WorkloadManager) unbindSpan(pod *corev1.Pod) {
	name := getCRDName(pod)
	delete(sMgr.bindCRDSpan, name)
}

func getLabelOrAnnotation(pod *corev1.Pod, primaryAnnotation string, secondaryAnnotation string) (string, bool) {
	if pod.Annotations[primaryAnnotation] != "" {
		return pod.Annotations[primaryAnnotation], true
	} else if pod.Labels[primaryAnnotation] != "" {
		return pod.Labels[primaryAnnotation], true
	} else if pod.Annotations[secondaryAnnotation] != "" {
		return pod.Annotations[secondaryAnnotation], true
	} else if pod.Labels[secondaryAnnotation] != "" {
		return pod.Labels[secondaryAnnotation], true
	}
	return "", false
}

func calculateVersion(pod *corev1.Pod) string {
	name := ""

	if len(pod.Spec.Containers) == 1 {
		image := strings.Split(pod.Spec.Containers[0].Image, ":")
		if len(image) > 0 && image[1] != "" && image[1] != "latest" {
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
