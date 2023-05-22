package gating

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/operator/webhooks/common"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/gate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create,versions=v1,name=gpod.keptn.sh,admissionReviewVersions=v1,sideEffects=None
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets;daemonsets;replicasets,verbs=get

// PodGatingWebhook annotates Pods to add gates
type PodGatingWebhook struct {
	Client   client.Client
	Tracer   trace.Tracer
	decoder  *admission.Decoder
	Recorder record.EventRecorder
	Log      logr.Logger
}

func (a *PodGatingWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	logger := log.FromContext(ctx).WithValues("webhook", "/gate-v1-pod", "object", map[string]interface{}{
		"name":      req.Name,
		"namespace": req.Namespace,
		"kind":      req.Kind,
	})
	logger.Info("Gate webhook for pod called")
	ctx, span := a.Tracer.Start(ctx, "gate_pod", trace.WithNewRoot(), trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	pod := &corev1.Pod{}
	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	logger.Info(fmt.Sprintf("Pod annotations: %v", pod.Annotations))

	podIsAnnotated, err := common.IsPodOrParentAnnotated(ctx, &req, pod, a.Client)
	logger.Info("Checked if pod is annotated.")
	if err != nil {
		span.SetStatus(codes.Error, common.InvalidAnnotationMessage)
		return admission.Errored(http.StatusBadRequest, err)
	}

	if podIsAnnotated {
		logger.Info("Resource is annotated with Keptn annotations, using Keptn scheduler")
		pod.Spec.SchedulingGates = []corev1.PodSchedulingGate{
			{
				Name: "klt-gated",
			},
		}
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}
