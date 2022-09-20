package klcpermit

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

var serviceResource = schema.GroupVersionResource{Group: "lifecycle.keptn.sh", Version: "v1alpha1", Resource: "servicerun"}

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
	// ServiceRunRunning means that all of the serviceRuns have been started.
	ServiceRunRunning string = "Running"
	// ServiceRunSucceeded means that all of the serviceRuns have been finished successfully.
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

func NewServiceManager(dy dynamic.Interface) *ServiceManager {
	sMgr := &ServiceManager{
		dynamicClient: dy,
	}
	return sMgr
}

func (sMgr *ServiceManager) Permit(ctx context.Context, pod *corev1.Pod) Status { //This is to not have tight coupling with CRD resource

	// list resources here
	services, err := sMgr.ListServiceRun(ctx, pod.Namespace)

	if err != nil {
		return ServiceRunNotSpecified
	}
	owner := GetPodOwner(pod)
	for _, s := range services {
		//get spec!
		replicaRef, found, err := unstructured.NestedFieldCopy(s.UnstructuredContent(), "spec", "replicasetref")

		if err == nil && found && replicaRef == owner {
			phase, found, err := unstructured.NestedString(s.UnstructuredContent(), "status", "phase")
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
