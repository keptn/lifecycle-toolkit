package klcpermit

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/keptn/lifecycle-toolkit/scheduler/pkg/tracing"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

var workloadInstanceResource = schema.GroupVersionResource{Group: "lifecycle.keptn.sh", Version: "v1alpha2", Resource: "keptnworkloadinstances"}

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
	StateRunning   KeptnState = "Running"
	StateSucceeded KeptnState = "Succeeded"
	StateFailed    KeptnState = "Failed"
	StateUnknown   KeptnState = "Unknown"
	StatePending   KeptnState = "Pending"
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
}

func NewWorkloadManager(d dynamic.Interface) *WorkloadManager {
	sMgr := &WorkloadManager{
		dynamicClient: d,
		Tracer:        otel.Tracer("keptn/scheduler"),
	}
	return sMgr
}

var bindCRDSpan = make(map[string]trace.Span, 100)

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
		case StateFailed:
			span.SetStatus(codes.Error, "Failed")
			span.End()
			unbindSpan(pod)
			return Failure
		case StateSucceeded:
			span.End()
			unbindSpan(pod)
			return Success
		case StatePending:
			return Wait
		case StateRunning:
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
	if span, ok := bindCRDSpan[name]; ok {
		return ctx, span
	}
	ctx, span := tracing.CreateSpan(ctx, crd, sMgr.Tracer, pod.Namespace)
	//TODO store only sampled one and cap it
	bindCRDSpan[name] = span
	return ctx, span
}

func (sMgr *WorkloadManager) ObserveWorkloadForPod(ctx context.Context, handler framework.WaitingPod, pod *corev1.Pod) {
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(sMgr.dynamicClient, 0, pod.GetNamespace(), nil)

	gvr, _ := schema.ParseResourceArg("keptnworkloadinstances.v1alpha1.lifecycle.keptn.sh")

	informer := factory.ForResource(*gvr)

	sMgr.startWatching(ctx, informer.Informer(), pod, handler)

}

func (sMgr *WorkloadManager) startWatching(ctx context.Context, s cache.SharedIndexInformer, pod *corev1.Pod, handler framework.WaitingPod) {
	stopCh := make(chan struct{})
	workloadInstanceName := getCRDName(pod)

	checkWorkloadInstance := func(obj interface{}) {
		unstructuredWI := obj.(*unstructured.Unstructured)

		if unstructuredWI.GetName() != workloadInstanceName {
			return
		}

		_, span := sMgr.getSpan(ctx, unstructuredWI, pod)

		phase, found, err := unstructured.NestedString(unstructuredWI.UnstructuredContent(), "status", "preDeploymentEvaluationStatus")
		klog.Infof("[Keptn Permit Plugin] workloadInstance crd %s, found %s with phase %s ", unstructuredWI, found, phase)
		if err == nil && found {
			span.AddEvent("StatusEvaluation", trace.WithAttributes(tracing.Status.String(phase)))
			switch KeptnState(phase) {
			case StateFailed:
				span.End()
				handler.Reject(PluginName, "Pre Deployment Check failed")
				unbindSpan(pod)
				stopCh <- struct{}{}
			case StateSucceeded:
				handler.Allow(PluginName)
				span.End()
				unbindSpan(pod)
				stopCh <- struct{}{}
			case StatePending:
			case StateRunning:
			case StateUnknown:
			}
		}
	}
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			checkWorkloadInstance(obj)
		},
		UpdateFunc: func(oldObj, obj interface{}) {
			checkWorkloadInstance(obj)
		},
		DeleteFunc: func(obj interface{}) {
			logrus.Info("received update event!")
		},
	}
	s.AddEventHandler(handlers)
	s.Run(stopCh)
}

func getCRDName(pod *corev1.Pod) string {
	application, _ := getLabelOrAnnotation(pod, AppAnnotation, K8sRecommendedAppAnnotations)
	workloadInstance, _ := getLabelOrAnnotation(pod, WorkloadAnnotation, K8sRecommendedWorkloadAnnotations)
	version, versionExists := getLabelOrAnnotation(pod, VersionAnnotation, K8sRecommendedVersionAnnotations)
	if !versionExists {
		version = calculateVersion(pod)
	}
	return application + "-" + workloadInstance + "-" + version
}

func unbindSpan(pod *corev1.Pod) {
	name := getCRDName(pod)
	delete(bindCRDSpan, name)
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
