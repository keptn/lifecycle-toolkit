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

type KeptnState string

const (
	StateRunning   KeptnState = "Running"
	StateSucceeded KeptnState = "Succeeded"
	StateFailed    KeptnState = "Failed"
	StateUnknown   KeptnState = "Unknown"
	StatePending   KeptnState = "Pending"
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
	phase, found, err := unstructured.NestedString(crd.UnstructuredContent(), "status", "preDeploymentStatus")
	klog.Infof("[Keptn Permit Plugin] workloadInstance crd %s, found %s with phase %s ", crd, found, phase)
	if err == nil && found {
		switch KeptnState(phase) {
		case StatePending:
			return Wait
		case StateFailed:
			return Failure
		case StateSucceeded:
			return Success
		case StateRunning:
			return Wait
		case StateUnknown:
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
