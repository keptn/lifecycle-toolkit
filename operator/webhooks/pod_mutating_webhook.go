package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"hash/fnv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.keptn.sh,admissionReviewVersions=v1,sideEffects=None

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
		logger.Info("Resource is annotated with Keptn annotations, using Keptn scheduler")
		//TODO uncomment this
		pod.Spec.SchedulerName = "keptn-scheduler"
		logger.Info("Pod annotaded, creating ServiceRun")
		if err := a.handleServiceRun(ctx, logger, pod, req.Namespace); err != nil {
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
	_, gotServiceAnnotation := pod.Annotations["keptn.sh/service"]
	_, gotVersionAnnotation := pod.Annotations["keptn.sh/version"]

	if gotApplicationAnnotation && gotServiceAnnotation && gotVersionAnnotation {
		return true
	}
	return false
}

func (a *PodMutatingWebhook) calculateVersion(pod *corev1.Pod) string {
	name := ""
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

func (r *PodMutatingWebhook) handleServiceRun(ctx context.Context, logger logr.Logger, pod *corev1.Pod, namespace string) error {
	serviceName, _ := pod.Annotations["keptn.sh/service"]

	logger.Info("Service name", "service", serviceName)

	service := &v1alpha1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: serviceName}, service)
	if errors.IsNotFound(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("could not fetch Service: %+v", err)
	}

	serviceRun := &v1alpha1.ServiceRun{}
	err = r.Client.Get(ctx, types.NamespacedName{Namespace: service.Namespace, Name: service.GetServiceRunName()}, serviceRun)
	if errors.IsNotFound(err) {
		logger.Info("Creating serviceRun from service", "service", service.Name)
		serviceRun, err := r.createServiceRun(ctx, service)
		if err != nil {
			logger.Error(err, "Could not create ServiceRun")
			return err
		}

		k8sEvent := r.generateK8sEvent(service, serviceRun)
		if err := r.Client.Create(ctx, k8sEvent); err != nil {
			logger.Error(err, "Could not send serviceRun created K8s event")
			return err
		}

		if err := r.Client.Status().Update(ctx, service); err != nil {
			logger.Error(err, "Could not update Service")
			return err
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("could not fetch ServiceRun: %+v", err)
	}

	return nil
}

func (r *PodMutatingWebhook) createServiceRun(ctx context.Context, service *v1alpha1.Service) (*v1alpha1.ServiceRun, error) {
	serviceRun := &v1alpha1.ServiceRun{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"keptn.sh/application": service.Spec.ApplicationName,
				"keptn.sh/service":     service.Name,
			},
			Name:      service.GetServiceRunName(),
			Namespace: service.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: service.APIVersion,
					Kind:       service.Kind,
					Name:       service.Name,
					UID:        service.UID,
				},
			},
		},
	}
	return serviceRun, r.Client.Create(ctx, serviceRun)
}

func (r *PodMutatingWebhook) generateSuffix() string {
	uid := uuid.New().String()
	return uid[:10]
}

func (r *PodMutatingWebhook) generateK8sEvent(service *v1alpha1.Service, serviceRun *v1alpha1.ServiceRun) *corev1.Event {
	return &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName:    serviceRun.Name + "-created-",
			Namespace:       serviceRun.Namespace,
			ResourceVersion: "v1alpha1",
			Labels: map[string]string{
				"keptn.sh/application": service.Spec.ApplicationName,
				"keptn.sh/service":     serviceRun.Name,
			},
		},
		InvolvedObject: corev1.ObjectReference{
			Kind:      serviceRun.Kind,
			Namespace: serviceRun.Namespace,
			Name:      serviceRun.Name,
		},
		Reason:  "created",
		Message: "serviceRun " + serviceRun.Name + " was created",
		Source: corev1.EventSource{
			Component: serviceRun.Kind,
		},
		Type: "Normal",
		EventTime: metav1.MicroTime{
			Time: time.Now().UTC(),
		},
		FirstTimestamp: metav1.Time{
			Time: time.Now().UTC(),
		},
		LastTimestamp: metav1.Time{
			Time: time.Now().UTC(),
		},
		Action:              "created",
		ReportingController: "serviceRun-controller",
		ReportingInstance:   "serviceRun-controller",
	}
}
