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

var serviceResource = schema.GroupVersionResource{Group: "lifecycle.keptn.sh", Version: "v1alpha1", Resource: "serviceruns"} //TODO change this resource name with workloadinstance and eventually appinstance :)

type Status string

const (
	ServiceRunStatusNotSpecified Status = "Service run status not specified"
	ServiceRunNotFound           Status = "Service run not found"
	Success                      Status = "Success"
	Failure                      Status = "Failure"
	Wait                         Status = "Wait"
)

const (
	// ServiceRunPending means the application has been accepted by the system, but one or more of its
	// serviceRuns has not been started.
	ServiceRunPending string = "Pending"
	// ServiceRunRunning means that serviceRun has been started.
	ServiceRunRunning string = "Running"
	// ServiceRunSucceeded means that serviceRun has been finished successfully.
	ServiceRunSucceeded string = "Succeeded"
	// ServiceRunFailed means that one or more pre-deployment checks was not successful and terminated.
	ServiceRunFailed string = "Failed"
	// ServiceRunUnknown means that for some reason the state of the application could not be obtained.
	ServiceRunUnknown string = "Unknown"
)

type Manager interface {
	Permit(context.Context, *corev1.Pod) Status
}

type ServiceManager struct {
	dynamicClient dynamic.Interface
}

func NewServiceManager(d dynamic.Interface) *ServiceManager {
	sMgr := &ServiceManager{
		dynamicClient: d,
	}
	return sMgr
}

func (sMgr *ServiceManager) Permit(ctx context.Context, pod *corev1.Pod) Status {
	//List service run CRDs
	name := GetCRDName(pod)
	crd, err := sMgr.GetCRD(ctx, metav1.NamespaceDefault, name)

	if err != nil {
		klog.Infof("[Keptn Permit Plugin] could not find service crd %s, err:%s", name, err.Error())
		return ServiceRunNotFound
	}
	//check CRD status
	phase, found, err := unstructured.NestedString(crd.UnstructuredContent(), "status", "phase")
	klog.Infof("[Keptn Permit Plugin] service crd %s, found %s with phase %s ", crd, found, phase)
	if err == nil && found {
		switch phase {
		case ServiceRunPending:
			return Wait
		case ServiceRunFailed:
			return Failure
		case ServiceRunSucceeded:
			return Success
		case ServiceRunRunning:
			return Wait
		case ServiceRunUnknown:
			return Wait
		}

	}
	return ServiceRunStatusNotSpecified
}

//GetCRD returns unstructured to avoid tight coupling with the CRD resource
func (sMgr *ServiceManager) GetCRD(ctx context.Context, namespace string, name string) (*unstructured.Unstructured, error) {
	// GET /apis/lifecycle.keptn.sh/v1/namespaces/{namespace}/servicerun/name
	return sMgr.dynamicClient.Resource(serviceResource).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

func GetCRDName(pod *corev1.Pod) string {
	application := pod.Annotations["keptn.sh/application"]
	service := pod.Annotations["keptn.sh/service"]
	version := pod.Annotations["keptn.sh/version"]
	return application + "-" + service + "-" + version
}
