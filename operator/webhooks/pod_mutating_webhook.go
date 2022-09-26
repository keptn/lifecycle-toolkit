package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"hash/fnv"

	v1 "k8s.io/api/batch/v1"
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
		pod.Spec.SchedulerName = "keptn-scheduler"
		logger.Info("Annotations", "annotations", pod.Annotations)
		if err := a.handleService(ctx, logger, pod, req.Namespace); err != nil {
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
	_, gotApplicationAnnotation := pod.Annotations[common.ApplicationAnnotation]
	_, gotServiceAnnotation := pod.Annotations[common.ServiceAnnotation]
	_, gotVersionAnnotation := pod.Annotations[common.VersionAnnotation]

	if gotApplicationAnnotation && gotServiceAnnotation {
		if !gotVersionAnnotation {
			pod.Annotations[common.VersionAnnotation] = a.calculateVersion(pod)
		}
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

func (r *PodMutatingWebhook) handleService(ctx context.Context, logger logr.Logger, pod *corev1.Pod, namespace string) error {
	serviceName, _ := pod.Annotations[common.ServiceAnnotation]

	service := &v1alpha1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: r.getServiceName(pod)}, service)
	if errors.IsNotFound(err) {
		logger.Info("Service name", "service", serviceName)

		logger.Info("Creating service service", "service", service.Name)
		service = r.generateService(ctx, pod)
		err = r.Client.Create(ctx, service)
		if err != nil {
			logger.Error(err, "Could not create Service")
			return err
		}

		k8sEvent := r.generateK8sEvent(service, "created")
		if err := r.Client.Create(ctx, k8sEvent); err != nil {
			logger.Error(err, "Could not send service created K8s event")
			return err
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("could not fetch ServiceRun: %+v", err)
	}

	if service.Spec.Version == pod.Annotations[common.VersionAnnotation] {
		return nil
	}

	service.Spec.Version = pod.Annotations[common.VersionAnnotation]
	service.Spec.ResourceReference = r.getResourceReference(pod)
	service.Annotations[common.VersionAnnotation] = pod.Annotations[common.VersionAnnotation]
	err = r.Client.Update(ctx, service)
	if err != nil {
		logger.Error(err, "Could not update Service")
		return err
	}

	k8sEvent := r.generateK8sEvent(service, "updated")
	if err := r.Client.Create(ctx, k8sEvent); err != nil {
		logger.Error(err, "Could not send service updated K8s event")
		return err
	}

	return nil
}

func (r *PodMutatingWebhook) generateService(ctx context.Context, pod *corev1.Pod) *v1alpha1.Service {
	version, _ := pod.Annotations[common.VersionAnnotation]
	applicationName, _ := pod.Annotations[common.ApplicationAnnotation]
	return &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: pod.Annotations,
			Name:        r.getServiceName(pod),
			Namespace:   pod.Namespace,
		},
		Spec: v1alpha1.ServiceSpec{
			ApplicationName:   applicationName,
			Version:           version,
			ResourceReference: r.getResourceReference(pod),
			//for now hardcoded, will be changed in future
			PreDeploymentCheck: v1alpha1.EventSpec{
				Service:     r.getServiceName(pod),
				Application: applicationName,
				JobSpec: v1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:    "hello-world",
									Image:   "ubuntu:latest",
									Command: []string{"echo", "Hello from Keptn"},
								},
							},
							RestartPolicy: corev1.RestartPolicyNever,
						},
					},
				},
			},
		},
	}
}

func (r *PodMutatingWebhook) generateK8sEvent(service *v1alpha1.Service, eventType string) *corev1.Event {
	return &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName:    service.Name + "-" + eventType + "-",
			Namespace:       service.Namespace,
			ResourceVersion: "v1alpha1",
			Labels: map[string]string{
				common.ApplicationAnnotation: service.Spec.ApplicationName,
				common.ServiceAnnotation:     service.Name,
			},
		},
		InvolvedObject: corev1.ObjectReference{
			Kind:      service.Kind,
			Namespace: service.Namespace,
			Name:      service.Name,
		},
		Reason:  eventType,
		Message: "serviceRun " + service.Name + " was " + eventType,
		Source: corev1.EventSource{
			Component: service.Kind,
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
		Action:              eventType,
		ReportingController: "webhook",
		ReportingInstance:   "webhook",
	}
}

func (r *PodMutatingWebhook) getServiceName(pod *corev1.Pod) string {
	serviceName, _ := pod.Annotations[common.ServiceAnnotation]
	applicationName, _ := pod.Annotations[common.ApplicationAnnotation]
	return strings.ToLower(applicationName + "-" + serviceName)
}

func (r *PodMutatingWebhook) getResourceReference(pod *corev1.Pod) v1alpha1.ResourceReference {
	reference := v1alpha1.ResourceReference{
		UID:  pod.UID,
		Kind: pod.Kind,
	}
	if len(pod.OwnerReferences) != 0 {
		for _, o := range pod.OwnerReferences {
			if o.Kind == "ReplicaSet" {
				reference.UID = o.UID
				reference.Kind = o.Kind
			}
		}
	}
	return reference
}
