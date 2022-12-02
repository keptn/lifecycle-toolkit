package klcpermit

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/keptn/lifecycle-toolkit/scheduler/pkg/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
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
	StateCancelled KeptnState = "Cancelled"
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

func (sMgr *WorkloadManager) ObserveWorkloadForPod(ctx context.Context, handler framework.WaitingPod, pod *corev1.Pod) {
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(sMgr.dynamicClient, 0, pod.GetNamespace(), nil)

	informer := factory.ForResource(workloadInstanceResource)

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
		if err != nil {
			klog.Errorf("[Keptn Permit Plugin] cannot fetch workloadInstance crd %s preDeploymentStatus: %s", unstructuredWI, err.Error())
			return
		}
		klog.Infof("[Keptn Permit Plugin] workloadInstance crd %s, found %s with phase %s ", unstructuredWI, found, phase)
		if found {
			span.AddEvent("StatusEvaluation", trace.WithAttributes(tracing.Status.String(phase)))
			switch KeptnState(phase) {
			case StateFailed, StateCancelled:
				handler.Reject(PluginName, "Pre Deployment Check failed")
				span.End()
				sMgr.unbindSpan(pod)
				stopCh <- struct{}{}
			case StateSucceeded:
				handler.Allow(PluginName)
				span.End()
				sMgr.unbindSpan(pod)
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
			klog.Info("received delete event!")
		},
	}
	s.AddEventHandler(handlers)
	s.Run(stopCh)
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

func (sMgr *WorkloadManager) unbindSpan(pod *corev1.Pod) {
	name := getCRDName(pod)
	delete(sMgr.bindCRDSpan, name)
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
