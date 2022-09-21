package klcpermit

import (
	"context"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
)

var serviceResource = schema.GroupVersionResource{Group: "lifecycle.keptn.sh", Version: "v1alpha1", Resource: "serviceruns"}

type Status string

const (
	ServiceRunNotSpecified Status = "Service run not specified"
	ServiceRunNotFound     Status = "Service run not found"
	Success                Status = "Success"
	Failure                Status = "Failure"
	Wait                   Status = "Wait"
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

func (sMgr *ServiceManager) Permit(ctx context.Context, pod *corev1.Pod) Status { //This is to not have tight coupling with CRD resource
	//List service run CRDs
	services, err := sMgr.ListServiceRun(ctx, metav1.NamespaceAll)

	if err != nil {
		klog.Infof("[Keptn Permit Plugin] could not find service crd err:%s svc:%+v", err.Error(), services)
		return ServiceRunNotSpecified
	}
	owner := GetPodOwner(pod)
	for _, s := range services {

		crd, err := GetCRD(s)
		if err != nil {
			klog.Infof("[Keptn Permit Plugin] spec error is %+v", err)
			return ServiceRunNotSpecified
		}

		replicasetUID, found, err := unstructured.NestedString(crd, "spec", "replicasetUID")

		if err != nil {
			klog.Infof("[Keptn Permit Plugin] spec error is %+v", err)
			return ServiceRunNotSpecified
		}
		//match pod to CRD
		if err == nil && found && types.UID(replicasetUID) == owner {

			//check CRD status
			phase, found, err := unstructured.NestedString(crd, "status", "phase")
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
		}
	}

	return ServiceRunNotFound
}

func (sMgr *ServiceManager) ListServiceRun(ctx context.Context, namespace string) ([]unstructured.Unstructured, error) {
	// GET /apis/lifecycle.keptn.sh/v1/namespaces/{namespace}/servicerun/
	list, err := sMgr.dynamicClient.Resource(serviceResource).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func GetPodOwner(pod *corev1.Pod) types.UID {
	for _, owner := range pod.OwnerReferences {
		if owner.Kind == "ReplicaSet" {
			return owner.UID
		}
	}
	return ""
}

// GetCRD takes an unstructured object from a dynamic lister and returns a maps to access the last applied changes of a CRD
func GetCRD(u unstructured.Unstructured) (map[string]interface{}, error) {
	spec := map[string]interface{}{}
	replicaRef, found, err := unstructured.NestedString(u.UnstructuredContent(), "metadata", "annotations", "kubectl.kubernetes.io/last-applied-configuration")
	if !found || err != nil {
		return spec, err
	}

	data := []byte(replicaRef)
	err = json.Unmarshal(data, &spec)
	if err != nil {
		return spec, err
	}

	return spec, err
}
