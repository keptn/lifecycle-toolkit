package pod_mutator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/semconv"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/webhooks/pod_mutator/handlers"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.keptn.sh,admissionReviewVersions=v1,sideEffects=None
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets;daemonsets;replicasets,verbs=get

// PodMutatingWebhook annotates Pods

const namespaceKey = "namespace"
const podKey = "pod"

type PodMutatingWebhook struct {
	Client                 client.Client
	Tracer                 trace.Tracer
	Decoder                handlers.Decoder
	EventSender            controllercommon.IEvent
	Log                    logr.Logger
	SchedulingGatesEnabled bool
	Pod                    handlers.PodAnnotationHandler
	Workload               handlers.K8sHandler
	App                    handlers.K8sHandler
}

func NewPodMutator(client client.Client,
	tracer trace.Tracer,
	decoder *admission.Decoder,
	eventSender controllercommon.IEvent,
	log logr.Logger,
	schedulingGatesEnabled bool) *PodMutatingWebhook {
	return &PodMutatingWebhook{
		SchedulingGatesEnabled: schedulingGatesEnabled,
		Client:                 client,
		Tracer:                 tracer,
		EventSender:            eventSender,
		Decoder:                decoder,
		Log:                    log,
		Pod:                    handlers.PodAnnotationHandler{Client: client, Log: log},
		App:                    &handlers.AppCreationRequestHandler{Log: log, Client: client, EventSender: eventSender, Tracer: tracer},
		Workload:               &handlers.WorkloadHandler{Log: log, Client: client, EventSender: eventSender, Tracer: tracer},
	}
}

// Handle inspects incoming Pods and injects the Keptn scheduler if they contain the Keptn lifecycle annotations.
func (a *PodMutatingWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {

	ctx, span := a.Tracer.Start(ctx, "annotate_pod", trace.WithNewRoot(), trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	a.Log.Info("webhook for pod called")

	pod := &corev1.Pod{}
	err := a.Decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// check if Lifecycle Operator is enabled for this namespace
	namespace := &corev1.Namespace{}
	if err = a.
		Client.Get(ctx, types.NamespacedName{Name: req.Namespace}, namespace); err != nil {
		a.Log.Error(err, "could not get namespace", namespaceKey, req.Namespace)
		return admission.Errored(http.StatusInternalServerError, err)
	}

	if namespace.GetAnnotations()[apicommon.NamespaceEnabledAnnotation] != "enabled" {
		a.Log.Info("namespace is not enabled for lifecycle operator", namespaceKey, req.Namespace)
		return admission.Allowed("namespace is not enabled for lifecycle operator")
	}

	// check the OwnerReference of the pod to see if it is supported and intended to be managed by Keptn
	ownerRef := handlers.GetOwnerReference(&pod.ObjectMeta)

	if ownerRef.Kind == "" {
		msg := "owner of pod is not supported by lifecycle operator"
		a.Log.Info(msg, namespaceKey, req.Namespace, podKey, req.Name)
		return admission.Allowed(msg)
	}

	a.Log.Info(fmt.Sprintf("Pod annotations: %v", pod.Annotations))

	if a.Pod.IsAnnotated(ctx, &req, pod) {
		a.Log.Info("Resource is annotated with Keptn annotations", namespaceKey, req.Namespace, podKey, req.Name)

		if scheduled := handleScheduling(a.SchedulingGatesEnabled, a.Log, pod); scheduled {
			return admission.Allowed("gate of the pod already removed")
		}

		a.Log.Info("Annotations", "annotations", pod.Annotations)
		semconv.AddAttributeFromAnnotations(span, pod.Annotations)
		a.Log.Info("Attributes from annotations set")

		if err := a.Workload.Handle(ctx, pod, req.Namespace); err != nil {
			a.Log.Error(err, "Could not handle Workload")
			span.SetStatus(codes.Error, err.Error())
			return admission.Errored(http.StatusBadRequest, err)
		}

		if err := a.App.Handle(ctx, pod, req.Namespace); err != nil {
			a.Log.Error(err, "Could not handle App")
			span.SetStatus(codes.Error, err.Error())
			return admission.Errored(http.StatusBadRequest, err)
		}
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to marshal")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

func handleScheduling(schedulingGatesEnabled bool, logger logr.Logger, pod *corev1.Pod) bool {
	if schedulingGatesEnabled {
		logger.Info("SchedulingGates enabled")
		_, gateRemoved := handlers.GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.SchedulingGateRemoved, "")
		if gateRemoved {
			return true
		}
		pod.Spec.SchedulingGates = []corev1.PodSchedulingGate{
			{
				Name: apicommon.KeptnGate,
			},
		}
	} else {
		logger.Info("SchedulingGates disabled, using keptn-scheduler")
		pod.Spec.SchedulerName = "keptn-scheduler"
	}
	return false
}
