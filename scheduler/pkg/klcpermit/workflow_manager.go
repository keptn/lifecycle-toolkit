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

var workloadVersionResource = schema.GroupVersionResource{Group: "lifecycle.keptn.sh", Version: "v1alpha4", Resource: "keptnworkloadversions"}

type Status string

const (
	WorkloadVersionStatusNotSpecified Status = "Workload run status not specified"
	WorkloadVersionNotFound           Status = "Workload run not found"
	Success                           Status = "Success"
	Failure                           Status = "Failure"
	Wait                              Status = "Wait"
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

const MinKLTNameLen = 80
const MaxK8sObjectLength = 253

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
	//List workloadVersion run CRDs
	name := getCRDName(pod)
	crd, err := sMgr.GetCRD(ctx, pod.Namespace, name)

	if err != nil {
		klog.Infof("[Keptn Permit Plugin] could not find workloadVersion crd %s, err:%s", name, err.Error())
		return WorkloadVersionNotFound
	}

	_, span := sMgr.getSpan(ctx, crd, pod)

	//check CRD status
	phase, found, err := unstructured.NestedString(crd.UnstructuredContent(), "status", "preDeploymentEvaluationStatus")
	klog.Infof("[Keptn Permit Plugin] workloadVersion crd %s, found %s with phase %s ", crd, found, phase)
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
	return WorkloadVersionStatusNotSpecified
}

// GetCRD returns unstructured to avoid tight coupling with the CRD resource
func (sMgr *WorkloadManager) GetCRD(ctx context.Context, namespace string, name string) (*unstructured.Unstructured, error) {
	// GET /apis/lifecycle.keptn.sh/v1/namespaces/{namespace}/workloadversion/name
	return sMgr.dynamicClient.Resource(workloadVersionResource).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
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

// CreateResourceName is a function that concatenates the parts from the
// input and checks, if the resulting string matches the maxLen condition.
// If it does not match, it reduces the subparts, starting with the first
// one (but leaving its length at least in minSubstrLen so it's not deleted
// completely) adn continuing with the rest if needed.
// Let's take WorkloadVersion as an example (3 parts: app, workload, version).
// First the app name is reduced if needed (only to minSubstrLen),
// afterwards workload and the version is not reduced at all. This pattern is
// chosen to not reduce only one part of the name (that can be completely gone
// afterwards), but to include all of the parts in the resulting string.
func createResourceName(maxLen int, minSubstrLen int, str ...string) string {
	// if the minSubstrLen is too long for the number of parts,
	// needs to be reduced
	for len(str)*minSubstrLen > maxLen {
		minSubstrLen = minSubstrLen / 2
	}
	// looping through the subparts and checking if the resulting string
	// matches the maxLen condition
	for i := 0; i < len(str)-1; i++ {
		newStr := strings.Join(str, "-")
		if len(newStr) > maxLen {
			// if the part is already smaller than the minSubstrLen,
			// this part cannot be reduced anymore, so we continue
			if len(str[i]) <= minSubstrLen {
				continue
			}
			// part needs to be reduced
			cut := len(newStr) - maxLen
			// if the needed reduction is bigger than the allowed
			// reduction on the part, it's reduced to the minimum
			if cut > len(str[i])-minSubstrLen {
				str[i] = str[i][:minSubstrLen]
			} else {
				// the needed reduction can be completed fully on this
				// part, so it's reduced accordingly
				str[i] = str[i][:len(str[i])-cut]
			}
		} else {
			return sanitizeResourceNameString(newStr)
		}
	}

	return sanitizeResourceNameString(strings.Join(str, "-"))
}

func sanitizeResourceNameString(name string) string {
	// ensure lower case and replace '_' with '-'
	// the replacement is done because names for resources generated by KLT
	// can be derived from label values, which can include '_' characters (https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set).
	// However, '_'  is not an allowed character for resource names (https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names).
	return strings.ReplaceAll(strings.ToLower(name), "_", "-")
}

func getCRDName(pod *corev1.Pod) string {
	application, _ := getLabelOrAnnotation(pod, AppAnnotation, K8sRecommendedAppAnnotations)
	workload, _ := getLabelOrAnnotation(pod, WorkloadAnnotation, K8sRecommendedWorkloadAnnotations)
	version, versionExists := getLabelOrAnnotation(pod, VersionAnnotation, K8sRecommendedVersionAnnotations)
	if !versionExists {
		version = calculateVersion(pod)
	}
	return createResourceName(MaxK8sObjectLength, MinKLTNameLen, application, workload, version)
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
