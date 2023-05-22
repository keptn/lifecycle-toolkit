package pod_mutator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/semconv"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/webhooks/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=update,versions=v1,name=mpod.keptn.sh,admissionReviewVersions=v1,sideEffects=None
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets;daemonsets;replicasets,verbs=get

// PodMutatingWebhook annotates Pods
type PodMutatingWebhook struct {
	Client   client.Client
	Tracer   trace.Tracer
	decoder  *admission.Decoder
	Recorder record.EventRecorder
	Log      logr.Logger
}

// Handle inspects incoming Pods and injects the Keptn scheduler if they contain the Keptn lifecycle annotations.
//
//nolint:gocyclo
func (a *PodMutatingWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {

	ctx, span := a.Tracer.Start(ctx, "annotate_pod", trace.WithNewRoot(), trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	logger := log.FromContext(ctx).WithValues("webhook", "/mutate-v1-pod", "object", map[string]interface{}{
		"name":      req.Name,
		"namespace": req.Namespace,
		"kind":      req.Kind,
	})
	logger.Info("webhook for pod called")

	pod := &corev1.Pod{}

	logger.Info("checking the request", "req obj: ", req.Object, "req raw", req.Object.Raw)
	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// check if Lifecycle Controller is enabled for this namespace
	namespace := &corev1.Namespace{}
	if err = a.Client.Get(ctx, types.NamespacedName{Name: req.Namespace}, namespace); err != nil {
		logger.Error(err, "could not get namespace", "namespace", req.Namespace)
		return admission.Errored(http.StatusInternalServerError, err)
	}

	if namespace.GetAnnotations()[apicommon.NamespaceEnabledAnnotation] != "enabled" {
		logger.Info("namespace is not enabled for lifecycle controller", "namespace", req.Namespace)
		return admission.Allowed("namespace is not enabled for lifecycle controller")
	}

	logger.Info(fmt.Sprintf("Pod annotations: %v", pod.Annotations))

	podIsAnnotated, err := common.IsPodOrParentAnnotated(ctx, &req, pod, a.Client)
	logger.Info("Checked if pod is annotated.")
	if err != nil {
		span.SetStatus(codes.Error, common.InvalidAnnotationMessage)
		return admission.Errored(http.StatusBadRequest, err)
	}

	if podIsAnnotated {

		logger.Info("Annotations", "annotations", pod.Annotations)

		isAppAnnotationPresent, err := a.isAppAnnotationPresent(pod)
		if err != nil {
			span.SetStatus(codes.Error, common.InvalidAnnotationMessage)
			return admission.Errored(http.StatusBadRequest, err)
		}
		semconv.AddAttributeFromAnnotations(span, pod.Annotations)
		logger.Info("Attributes from annotations set")

		if err := a.handleWorkload(ctx, logger, pod, req.Namespace); err != nil {
			logger.Error(err, "Could not handle Workload")
			span.SetStatus(codes.Error, err.Error())
			return admission.Errored(http.StatusBadRequest, err)
		}

		if err := a.handleApp(ctx, logger, pod, req.Namespace, isAppAnnotationPresent); err != nil {
			logger.Error(err, "Could not handle App")
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

// PodMutatingWebhook implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (a *PodMutatingWebhook) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}

func (a *PodMutatingWebhook) isAppAnnotationPresent(pod *corev1.Pod) (bool, error) {
	app, gotAppAnnotation := common.GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)

	if gotAppAnnotation {
		if len(app) > apicommon.MaxAppNameLength {
			return false, common.ErrTooLongAnnotations
		}
		return true, nil
	}

	if len(pod.Annotations) == 0 {
		pod.Annotations = make(map[string]string)
	}
	pod.Annotations[apicommon.AppAnnotation], _ = common.GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	return false, nil
}

//nolint:dupl
func (a *PodMutatingWebhook) handleWorkload(ctx context.Context, logger logr.Logger, pod *corev1.Pod, namespace string) error {

	ctx, span := a.Tracer.Start(ctx, "create_workload", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	newWorkload := a.generateWorkload(ctx, pod, namespace)

	newWorkload.SetSpanAttributes(span)

	logger.Info("Searching for workload")

	workload := &klcv1alpha3.KeptnWorkload{}
	err := a.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: newWorkload.Name}, workload)
	if errors.IsNotFound(err) {
		logger.Info("Creating workload", "workload", workload.Name)
		workload = newWorkload
		err = a.Client.Create(ctx, workload)
		if err != nil {
			logger.Error(err, "Could not create Workload")
			controllercommon.RecordEvent(a.Recorder, apicommon.PhaseCreateWorkload, "Warning", workload, "WorkloadNotCreated", "could not create KeptnWorkload", workload.Spec.Version)
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		controllercommon.RecordEvent(a.Recorder, apicommon.PhaseCreateWorkload, "Normal", workload, "WorkloadCreated", "created KeptnWorkload", workload.Spec.Version)
		return nil
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("could not fetch Workload"+": %+v", err)
	}

	if reflect.DeepEqual(workload.Spec, newWorkload.Spec) {
		logger.Info("Pod not changed, not updating anything")
		return nil
	}

	logger.Info("Pod changed, updating workload")
	workload.Spec = newWorkload.Spec

	err = a.Client.Update(ctx, workload)
	if err != nil {
		logger.Error(err, "Could not update Workload")
		controllercommon.RecordEvent(a.Recorder, apicommon.PhaseCreateWorkload, "Warning", workload, "WorkloadNotUpdated", "could not update KeptnWorkload", workload.Spec.Version)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	controllercommon.RecordEvent(a.Recorder, apicommon.PhaseCreateWorkload, "Normal", workload, "WorkloadUpdated", "updated KeptnWorkload", workload.Spec.Version)

	return nil
}

//nolint:dupl
func (a *PodMutatingWebhook) handleApp(ctx context.Context, logger logr.Logger, pod *corev1.Pod, namespace string, isAppAnnotationPresent bool) error {

	ctx, span := a.Tracer.Start(ctx, "create_app", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	newAppCreationRequest := a.generateAppCreationRequest(ctx, pod, namespace, isAppAnnotationPresent)

	newAppCreationRequest.SetSpanAttributes(span)

	logger.Info("Searching for AppCreationRequest", "appCreationRequest", newAppCreationRequest.Name)

	appCreationRequest := &klcv1alpha3.KeptnAppCreationRequest{}
	err := a.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: newAppCreationRequest.Name}, appCreationRequest)
	if errors.IsNotFound(err) {
		logger.Info("Creating app creation request", "appCreationRequest", appCreationRequest.Name)
		appCreationRequest = newAppCreationRequest
		err = a.Client.Create(ctx, appCreationRequest)
		if err != nil {
			logger.Error(err, "Could not create App")
			controllercommon.RecordEvent(a.Recorder, apicommon.PhaseCreateApp, "Warning", appCreationRequest, "AppCreationRequestNotCreated", "could not create KeptnAppCreationRequest", appCreationRequest.Spec.AppName)
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		controllercommon.RecordEvent(a.Recorder, apicommon.PhaseCreateApp, "Normal", appCreationRequest, "AppCreationRequestCreated", "created KeptnAppCreationRequest", appCreationRequest.Spec.AppName)
		return nil
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("could not fetch AppCreationRequest"+": %+v", err)
	}

	return nil
}

func (a *PodMutatingWebhook) generateWorkload(ctx context.Context, pod *corev1.Pod, namespace string) *klcv1alpha3.KeptnWorkload {
	version, _ := common.GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.VersionAnnotation, apicommon.K8sRecommendedVersionAnnotations)
	applicationName, _ := common.GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)

	var preDeploymentTasks []string
	var postDeploymentTasks []string
	var preDeploymentEvaluation []string
	var postDeploymentEvaluation []string

	if annotations, found := common.GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.PreDeploymentTaskAnnotation, ""); found {
		preDeploymentTasks = strings.Split(annotations, ",")
	}

	if annotations, found := common.GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.PostDeploymentTaskAnnotation, ""); found {
		postDeploymentTasks = strings.Split(annotations, ",")
	}

	if annotations, found := common.GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.PreDeploymentEvaluationAnnotation, ""); found {
		preDeploymentEvaluation = strings.Split(annotations, ",")
	}

	if annotations, found := common.GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.PostDeploymentEvaluationAnnotation, ""); found {
		postDeploymentEvaluation = strings.Split(annotations, ",")
	}

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	ownerRef := common.GetOwnerReference(&pod.ObjectMeta)

	return &klcv1alpha3.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:        a.getWorkloadName(pod),
			Namespace:   namespace,
			Annotations: traceContextCarrier,
			OwnerReferences: []metav1.OwnerReference{
				ownerRef,
			},
		},
		Spec: klcv1alpha3.KeptnWorkloadSpec{
			AppName:                   applicationName,
			Version:                   version,
			ResourceReference:         klcv1alpha3.ResourceReference{UID: ownerRef.UID, Kind: ownerRef.Kind, Name: ownerRef.Name},
			PreDeploymentTasks:        preDeploymentTasks,
			PostDeploymentTasks:       postDeploymentTasks,
			PreDeploymentEvaluations:  preDeploymentEvaluation,
			PostDeploymentEvaluations: postDeploymentEvaluation,
		},
	}
}

func (a *PodMutatingWebhook) generateAppCreationRequest(ctx context.Context, pod *corev1.Pod, namespace string, isAppAnnotationPresent bool) *klcv1alpha3.KeptnAppCreationRequest {
	appName := a.getAppName(pod)

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	kacr := &klcv1alpha3.KeptnAppCreationRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:        appName,
			Namespace:   namespace,
			Annotations: traceContextCarrier,
		},
		Spec: klcv1alpha3.KeptnAppCreationRequestSpec{
			AppName: appName,
		},
	}

	if !isAppAnnotationPresent {
		kacr.Annotations[apicommon.AppTypeAnnotation] = string(apicommon.AppTypeSingleService)
	}

	return kacr
}

func (a *PodMutatingWebhook) getWorkloadName(pod *corev1.Pod) string {
	workloadName, _ := common.GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	applicationName, _ := common.GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)
	return strings.ToLower(applicationName + "-" + workloadName)
}

func (a *PodMutatingWebhook) getAppName(pod *corev1.Pod) string {
	applicationName, _ := common.GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)
	return strings.ToLower(applicationName)
}
