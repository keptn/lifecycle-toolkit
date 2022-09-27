package klcpermit

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
)

var workloadInstanceResource = schema.GroupVersionResource{Group: "lifecycle.keptn.sh", Version: "v1alpha1", Resource: "keptnworkloadinstances"}

type Status string

const (
	WorkloadInstanceStatusNotSpecified Status = "Workload run status not specified"
	WorkloadInstanceNotFound           Status = "Workload run not found"
	Success                            Status = "Success"
	Failure                            Status = "Failure"
	Wait                               Status = "Wait"
)

type WorkloadInstancePhase string

const (
	// WorkloadInstancePending means the application has been accepted by the system, but one or more of its
	// workloadInstances has not been started.
	WorkloadInstancePending WorkloadInstancePhase = "Pending"
	// WorkloadInstanceRunning means that workloadInstance has been started.
	WorkloadInstanceRunning WorkloadInstancePhase = "Running"
	// WorkloadInstanceSucceeded means that workloadInstance has been finished successfully.
	WorkloadInstanceSucceeded WorkloadInstancePhase = "Succeeded"
	// WorkloadInstanceFailed means that one or more pre-deployment checks was not successful and terminated.
	WorkloadInstanceFailed WorkloadInstancePhase = "Failed"
	// WorkloadInstanceUnknown means that for some reason the state of the application could not be obtained.
	WorkloadInstanceUnknown WorkloadInstancePhase = "Unknown"
)

type Manager interface {
	Permit(context.Context, *corev1.Pod) Status
}

type WorkloadManager struct {
	dynamicClient dynamic.Interface
}

func NewWorkloadManager(d dynamic.Interface) *WorkloadManager {
	sMgr := &WorkloadManager{
		dynamicClient: d,
	}
	return sMgr
}

func (sMgr *WorkloadManager) Permit(ctx context.Context, pod *corev1.Pod) Status {
	//List workloadInstance run CRDs
	name := GetCRDName(pod)
	crd, err := sMgr.GetCRD(ctx, metav1.NamespaceDefault, name)

	if err != nil {
		klog.Infof("[Keptn Permit Plugin] could not find workloadInstance crd %s, err:%s", name, err.Error())
		return WorkloadInstanceNotFound
	}
	//check CRD status
	phase, found, err := unstructured.NestedString(crd.UnstructuredContent(), "status", "phase")
	klog.Infof("[Keptn Permit Plugin] workloadInstance crd %s, found %s with phase %s ", crd, found, phase)
	if err == nil && found {
		switch WorkloadInstancePhase(phase) {
		case WorkloadInstancePending:
			return Wait
		case WorkloadInstanceFailed:
			return Failure
		case WorkloadInstanceSucceeded:
			return Success
		case WorkloadInstanceRunning:
			return Wait
		case WorkloadInstanceUnknown:
			return Wait
		}

	}
	return WorkloadInstanceStatusNotSpecified
}

//GetCRD returns unstructured to avoid tight coupling with the CRD resource
func (sMgr *WorkloadManager) GetCRD(ctx context.Context, namespace string, name string) (*unstructured.Unstructured, error) {
	// GET /apis/lifecycle.keptn.sh/v1/namespaces/{namespace}/workloadinstance/name
	return sMgr.dynamicClient.Resource(workloadInstanceResource).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

func GetCRDName(pod *corev1.Pod) string {
	application := pod.Annotations["keptn.sh/app"]
	workloadInstance := pod.Annotations["keptn.sh/workload"]
	version := pod.Annotations["keptn.sh/version"]
	return application + "-" + workloadInstance + "-" + version
}
