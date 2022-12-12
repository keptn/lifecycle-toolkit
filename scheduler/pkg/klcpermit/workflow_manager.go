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
	handler       framework.Handle
}

func NewWorkloadManager(handler framework.Handle) (*WorkloadManager, error) {
	client, err := newClient(handler)
	sMgr := &WorkloadManager{
		dynamicClient: client,
		Tracer:        otel.Tracer("keptn/scheduler"),
		bindCRDSpan:   make(map[string]trace.Span, 100),
		handler:       handler,
	}
	return sMgr, err
}

func newClient(handle framework.Handle) (dynamic.Interface, error) {
	config := handle.KubeConfig()

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return dynClient, nil
}

func (sMgr *WorkloadManager) ObserveWorkloadForPod(ctx context.Context, pod *corev1.Pod) {
	// TODO investigate how long the defaultRSync interval should be
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(sMgr.dynamicClient, 10, pod.GetNamespace(), nil)

	sMgr.startWatching(ctx, factory, pod)
}

func (sMgr *WorkloadManager) startWatching(ctx context.Context, factory dynamicinformer.DynamicSharedInformerFactory, pod *corev1.Pod) {
	workloadInstanceName := getCRDName(pod)

	stopCh := make(chan struct{})

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
		waitingPodHandler := sMgr.handler.GetWaitingPod(pod.UID)
		if waitingPodHandler == nil {
			return
		}
		klog.Infof("[Keptn Permit Plugin] workloadInstance crd %s, found %s with phase %s ", unstructuredWI, found, phase)
		if found {
			span.AddEvent("StatusEvaluation", trace.WithAttributes(tracing.Status.String(phase)))
			switch KeptnState(phase) {
			case StateFailed, StateCancelled:
				waitingPodHandler.Reject(PluginName, "Pre Deployment Check failed")
				span.End()
				sMgr.unbindSpan(pod)
				stopCh <- struct{}{}
			case StateSucceeded:
				waitingPodHandler.Allow(PluginName)
				span.End()
				sMgr.unbindSpan(pod)
				stopCh <- struct{}{}
			case StatePending:
			case StateRunning:
			case StateUnknown:
			}
		}
	}
	handlers := sMgr.createHandler(checkWorkloadInstance)

	informer := factory.ForResource(workloadInstanceResource).Informer()
	informer.AddEventHandler(handlers)
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)
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

func (sMgr *WorkloadManager) createHandler(checkFunc func(obj interface{})) cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			checkFunc(obj)
		},
		UpdateFunc: func(oldObj, obj interface{}) {
			checkFunc(obj)
		},
		DeleteFunc: func(obj interface{}) {
			klog.Info("received delete event!")
		},
	}
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
