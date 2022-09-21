package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
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
