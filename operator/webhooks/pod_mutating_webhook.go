package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
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

	isAnnotated, err := a.isKeptnAnnotated(pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	if isAnnotated {
		logger.Info("Resource is annotated with Keptn annotations, using Keptn scheduler")
		pod.Spec.SchedulerName = "keptn-scheduler"
		logger.Info("Annotations", "annotations", pod.Annotations)
		if err := a.handleWorkload(ctx, logger, pod, req.Namespace); err != nil {
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

func (a *PodMutatingWebhook) isKeptnAnnotated(pod *corev1.Pod) (bool, error) {
	app, gotAppAnnotation := pod.Annotations[common.AppAnnotation]
	workload, gotWorkloadAnnotation := pod.Annotations[common.WorkloadAnnotation]
	version, gotVersionAnnotation := pod.Annotations[common.VersionAnnotation]

	if len(app) > common.MaxAppNameLength || len(workload) > common.MaxWorkloadNameLength || len(version) > common.MaxVersionLength {
		return false, common.ErrTooLongAnnotations
	}

	if gotAppAnnotation && gotWorkloadAnnotation {
		if !gotVersionAnnotation {
			pod.Annotations[common.VersionAnnotation] = a.calculateVersion(pod)
		}
		return true, nil
	}
	return false, nil
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

func (a *PodMutatingWebhook) handleWorkload(ctx context.Context, logger logr.Logger, pod *corev1.Pod, namespace string) error {
	workloadName, _ := pod.Annotations[common.WorkloadAnnotation]

	workload := &klcv1alpha1.KeptnWorkload{}
	err := a.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: a.getWorkloadName(pod)}, workload)
	if errors.IsNotFound(err) {
		logger.Info("Workload name", "workload", workloadName)

		logger.Info("Creating workload workload", "workload", workload.Name)
		workload = a.generateWorkload(pod, namespace)
		err = a.Client.Create(ctx, workload)
		if err != nil {
			logger.Error(err, "Could not create Workload")
			return err
		}

		k8sEvent := a.generateK8sEvent(workload, "created")
		if err := a.Client.Create(ctx, k8sEvent); err != nil {
			logger.Error(err, "Could not send workload created K8s event")
			return err
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("could not fetch WorkloadInstance"+": %+v", err)
	}

	if workload.Spec.Version == pod.Annotations[common.VersionAnnotation] {
		return nil
	}

	workload.Spec.Version = pod.Annotations[common.VersionAnnotation]
	workload.Spec.ResourceReference = a.getResourceReference(pod)
	workload.Annotations[common.VersionAnnotation] = pod.Annotations[common.VersionAnnotation]
	err = a.Client.Update(ctx, workload)
	if err != nil {
		logger.Error(err, "Could not update Workload")
		return err
	}

	k8sEvent := a.generateK8sEvent(workload, "updated")
	if err := a.Client.Create(ctx, k8sEvent); err != nil {
		logger.Error(err, "Could not send workload updated K8s event")
		return err
	}

	return nil
}

func (a *PodMutatingWebhook) generateWorkload(pod *corev1.Pod, namespace string) *klcv1alpha1.KeptnWorkload {
	version, _ := pod.Annotations[common.VersionAnnotation]
	applicationName, _ := pod.Annotations[common.AppAnnotation]

	var preDeploymentTasks []string
	var postDeploymentTasks []string
	var preDeploymentAnalysis []string
	var postDeploymentAnalysis []string

	if pod.Annotations[common.PreDeploymentTaskAnnotation] != "" {
		preDeploymentTasks = strings.Split(pod.Annotations[common.PreDeploymentTaskAnnotation], ",")
	}

	if pod.Annotations[common.PostDeploymentTaskAnnotation] != "" {
		postDeploymentTasks = strings.Split(pod.Annotations[common.PostDeploymentTaskAnnotation], ",")
	}

	if pod.Annotations[common.PreDeploymentAnalysisAnnotation] != "" {
		preDeploymentAnalysis = strings.Split(pod.Annotations[common.PreDeploymentAnalysisAnnotation], ",")
	}

	if pod.Annotations[common.PostDeploymentAnalysisAnnotation] != "" {
		postDeploymentAnalysis = strings.Split(pod.Annotations[common.PostDeploymentAnalysisAnnotation], ",")
	}

	return &klcv1alpha1.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: pod.Annotations,
			Name:        a.getWorkloadName(pod),
			Namespace:   namespace,
		},
		Spec: klcv1alpha1.KeptnWorkloadSpec{
			AppName:                applicationName,
			Version:                version,
			ResourceReference:      a.getResourceReference(pod),
			PreDeploymentTasks:     preDeploymentTasks,
			PostDeploymentTasks:    postDeploymentTasks,
			PreDeploymentAnalysis:  preDeploymentAnalysis,
			PostDeploymentAnalysis: postDeploymentAnalysis,
		},
	}
}

func (a *PodMutatingWebhook) generateK8sEvent(workload *klcv1alpha1.KeptnWorkload, eventType string) *corev1.Event {
	return &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName:    workload.Name + "-" + eventType + "-",
			Namespace:       workload.Namespace,
			ResourceVersion: "v1alpha1",
			Labels: map[string]string{
				common.AppAnnotation:      workload.Spec.AppName,
				common.WorkloadAnnotation: workload.Name,
			},
		},
		InvolvedObject: corev1.ObjectReference{
			Kind:      workload.Kind,
			Namespace: workload.Namespace,
			Name:      workload.Name,
		},
		Reason:  eventType,
		Message: "WorkloadInstance " + workload.Name + " was " + eventType,
		Source: corev1.EventSource{
			Component: workload.Kind,
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

func (a *PodMutatingWebhook) getWorkloadName(pod *corev1.Pod) string {
	workloadName, _ := pod.Annotations[common.WorkloadAnnotation]
	applicationName, _ := pod.Annotations[common.AppAnnotation]
	return strings.ToLower(applicationName + "-" + workloadName)
}

func (a *PodMutatingWebhook) getResourceReference(pod *corev1.Pod) klcv1alpha1.ResourceReference {
	reference := klcv1alpha1.ResourceReference{
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
