package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.keptn.sh,admissionReviewVersions=v1,sideEffects=None
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=serviceruns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=serviceruns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=serviceruns/finalizers,verbs=update
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=services,verbs=get;list;watch

// PodMutatingWebhook annotates Pods
type PodMutatingWebhook struct {
	Client  client.Client
	decoder *admission.Decoder
}

// Handle inspects incoming Pods and injects the Keptn scheduler if they contain the Keptn lifecycle annotations.
func (a *PodMutatingWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	logger := log.FromContext(ctx).WithValues("webhook", "/mutate-v1-pod", "object", map[string]interface{}{
		"name":      req.Name,
		"namespace": req.Namespace,
		"kind":      req.Kind,
	})
	logger.Info("webhook for pod called")

	pod := &corev1.Pod{}

	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	logger.Info(fmt.Sprintf("Pod annotations: %v", pod.Annotations))

	if a.isKeptnAnnotated(pod) {
		//pod.Spec.SchedulerName = "keptn-scheduler"
		logger.Info("Pod annotaded, creating ServiceRun")
		if err := a.handleServiceRun(ctx, logger, pod); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// PodMutatingWebhook implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (a *PodMutatingWebhook) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}

func (a *PodMutatingWebhook) isKeptnAnnotated(pod *corev1.Pod) bool {
	_, gotApplicationAnnotation := pod.Annotations["keptn.sh/application"]
	_, gotServiceAnnotation := pod.Annotations["keptn.sh/application"]

	if gotApplicationAnnotation && gotServiceAnnotation {
		return true
	}
	return false
}

func (r *PodMutatingWebhook) handleServiceRun(ctx context.Context, logger logr.Logger, pod *corev1.Pod) error {
	serviceName, _ := pod.Annotations["keptn.sh/service"]

	logger.Info("Service name", "service", serviceName)

	var replicaSetUID types.UID
	podRefs := pod.GetOwnerReferences()
	for _, ref := range podRefs {
		if ref.Kind == "ReplicaSet" {
			replicaSetUID = ref.UID
		}
	}

	logger.Info("ResplicaSerUID", "uid", replicaSetUID)

	serviceRunList := &v1alpha1.ServiceRunList{}
	_ = r.Client.List(ctx, serviceRunList, &client.ListOptions{Namespace: pod.Namespace, FieldSelector: fields.OneTermEqualSelector("spec.replicaSetUID", string(replicaSetUID))})
	// if err != nil {
	// 	logger.Error(err, "Cannot fetch ServiceRunList")
	// 	return fmt.Errorf("could not fetch ServiceRunList: %+v", err)
	// }

	logger.Info("ServiceRunList", "list", serviceRunList)

	if len(serviceRunList.Items) == 0 {
		logger.Info("Searching for service")

		service := &v1alpha1.Service{}
		err := r.Client.Get(ctx, types.NamespacedName{Name: serviceName, Namespace: pod.Namespace}, service)
		if err != nil {
			logger.Error(err, "Cannot fetch Service")
			return fmt.Errorf("could not fetch Service: %+v", err)
		}

		logger.Info("ServiceRun does not exist, creating")

		_, err = r.createServiceRun(ctx, service, replicaSetUID)
		if err != nil {
			logger.Error(err, "Could not create ServiceRun")
			return err
		}
		return nil
	}

	return nil
}

func (r *PodMutatingWebhook) createServiceRun(ctx context.Context, service *v1alpha1.Service, replicaSetUID types.UID) (string, error) {
	serviceRun := &v1alpha1.ServiceRun{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"keptn.sh/application": service.Spec.ApplicationName,
				"keptn.sh/service":     service.Name,
			},
			Name:      service.Name + "-" + r.generateSuffix(),
			Namespace: service.Namespace,
		},
		Spec: v1alpha1.ServiceRunSpec{
			ServiceName:   service.Name,
			ReplicaSetUID: replicaSetUID,
		},
	}
	for i := 0; i < 5; i++ {
		if err := r.Client.Create(ctx, serviceRun); err != nil {
			if errors.IsAlreadyExists(err) {
				serviceRun.Name = service.Name + "-" + r.generateSuffix()
				continue
			}
			return "", err
		}
		break
	}
	return serviceRun.Name, nil
}

func (r *PodMutatingWebhook) generateSuffix() string {
	uid := uuid.New().String()
	return uid[:10]
}
