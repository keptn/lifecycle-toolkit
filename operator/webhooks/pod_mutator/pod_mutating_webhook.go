package pod_mutator

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"reflect"
	"strings"

	argov1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/semconv"
	operatorcommon "github.com/keptn/lifecycle-toolkit/operator/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.keptn.sh,admissionReviewVersions=v1,sideEffects=None
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets;daemonsets;replicasets,verbs=get

// PodMutatingWebhook annotates Pods
type PodMutatingWebhook struct {
	Client   client.Client
	Tracer   trace.Tracer
	Decoder  *admission.Decoder
	Recorder record.EventRecorder
	Log      logr.Logger
}

const InvalidAnnotationMessage = "Invalid annotations"

var ErrTooLongAnnotations = fmt.Errorf("too long annotations, maximum length for app and workload is 25 characters, for version 12 characters")

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

	err := a.Decoder.Decode(req, pod)
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

	podIsAnnotated := a.isPodAnnotated(pod)
	logger.Info("Checked if pod is annotated.")

	if !podIsAnnotated {
		logger.Info("Pod is not annotated, check for parent annotations...")
		podIsAnnotated = a.copyAnnotationsIfParentAnnotated(ctx, &req, pod)
	}

	if podIsAnnotated {
		logger.Info("Resource is annotated with Keptn annotations, using Keptn scheduler")
		pod.Spec.SchedulerName = "keptn-scheduler"
		logger.Info("Annotations", "annotations", pod.Annotations)

		isAppAnnotationPresent := a.isAppAnnotationPresent(pod)
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

func (a *PodMutatingWebhook) isPodAnnotated(pod *corev1.Pod) bool {
	_, gotWorkloadAnnotation := getLabelOrAnnotation(&pod.ObjectMeta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	_, gotVersionAnnotation := getLabelOrAnnotation(&pod.ObjectMeta, apicommon.VersionAnnotation, apicommon.K8sRecommendedVersionAnnotations)

	if gotWorkloadAnnotation {
		if !gotVersionAnnotation {
			if len(pod.Annotations) == 0 {
				pod.Annotations = make(map[string]string)
			}
			pod.Annotations[apicommon.VersionAnnotation] = a.calculateVersion(pod)
		}
		return true
	}
	return false
}

func (a *PodMutatingWebhook) copyAnnotationsIfParentAnnotated(ctx context.Context, req *admission.Request, pod *corev1.Pod) bool {
	podOwner := a.getOwnerReference(&pod.ObjectMeta)
	if podOwner.UID == "" {
		return false
	}

	switch podOwner.Kind {
	case "ReplicaSet":
		rs := &appsv1.ReplicaSet{}
		if err := a.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: podOwner.Name}, rs); err != nil {
			return false
		}
		a.Log.Info("Done fetching RS")

		rsOwner := a.getOwnerReference(&rs.ObjectMeta)
		if rsOwner.UID == "" {
			return false
		}

		if rsOwner.Kind == "Rollout" {
			ro := &argov1alpha1.Rollout{}
			return a.fetchParentObjectAndCopyLabels(ctx, podOwner.Name, req.Namespace, pod, ro)
		}
		dp := &appsv1.Deployment{}
		return a.fetchParentObjectAndCopyLabels(ctx, rsOwner.Name, req.Namespace, pod, dp)

	case "StatefulSet":
		sts := &appsv1.StatefulSet{}
		return a.fetchParentObjectAndCopyLabels(ctx, podOwner.Name, req.Namespace, pod, sts)
	case "DaemonSet":
		ds := &appsv1.DaemonSet{}
		return a.fetchParentObjectAndCopyLabels(ctx, podOwner.Name, req.Namespace, pod, ds)
	default:
		return false
	}
}

func (a *PodMutatingWebhook) fetchParentObjectAndCopyLabels(ctx context.Context, name string, namespace string, pod *corev1.Pod, objectContainer client.Object) bool {
	if err := a.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, objectContainer); err != nil {
		return false
	}
	objectContainerMetaData := metav1.ObjectMeta{
		Labels:      objectContainer.GetLabels(),
		Annotations: objectContainer.GetAnnotations(),
	}
	return a.copyResourceLabelsIfPresent(&objectContainerMetaData, pod)
}

func (a *PodMutatingWebhook) copyResourceLabelsIfPresent(sourceResource *metav1.ObjectMeta, targetPod *corev1.Pod) bool {
	var workloadName, appName, version, preDeploymentChecks, postDeploymentChecks, preEvaluationChecks, postEvaluationChecks string
	var gotWorkloadName, gotVersion bool

	workloadName, gotWorkloadName = getLabelOrAnnotation(sourceResource, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	appName, _ = getLabelOrAnnotation(sourceResource, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)
	version, gotVersion = getLabelOrAnnotation(sourceResource, apicommon.VersionAnnotation, apicommon.K8sRecommendedVersionAnnotations)
	preDeploymentChecks, _ = getLabelOrAnnotation(sourceResource, apicommon.PreDeploymentTaskAnnotation, "")
	postDeploymentChecks, _ = getLabelOrAnnotation(sourceResource, apicommon.PostDeploymentTaskAnnotation, "")
	preEvaluationChecks, _ = getLabelOrAnnotation(sourceResource, apicommon.PreDeploymentEvaluationAnnotation, "")
	postEvaluationChecks, _ = getLabelOrAnnotation(sourceResource, apicommon.PostDeploymentEvaluationAnnotation, "")

	if len(targetPod.Annotations) == 0 {
		targetPod.Annotations = make(map[string]string)
	}

	if gotWorkloadName {
		setMapKey(targetPod.Annotations, apicommon.WorkloadAnnotation, workloadName)

		if !gotVersion {
			setMapKey(targetPod.Annotations, apicommon.VersionAnnotation, a.calculateVersion(targetPod))
		} else {
			setMapKey(targetPod.Annotations, apicommon.VersionAnnotation, version)
		}

		setMapKey(targetPod.Annotations, apicommon.AppAnnotation, appName)
		setMapKey(targetPod.Annotations, apicommon.PreDeploymentTaskAnnotation, preDeploymentChecks)
		setMapKey(targetPod.Annotations, apicommon.PostDeploymentTaskAnnotation, postDeploymentChecks)
		setMapKey(targetPod.Annotations, apicommon.PreDeploymentEvaluationAnnotation, preEvaluationChecks)
		setMapKey(targetPod.Annotations, apicommon.PostDeploymentEvaluationAnnotation, postEvaluationChecks)

		return true
	}
	return false
}

func (a *PodMutatingWebhook) isAppAnnotationPresent(pod *corev1.Pod) bool {
	_, gotAppAnnotation := getLabelOrAnnotation(&pod.ObjectMeta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)

	if gotAppAnnotation {
		return true
	}

	if len(pod.Annotations) == 0 {
		pod.Annotations = make(map[string]string)
	}
	pod.Annotations[apicommon.AppAnnotation], _ = getLabelOrAnnotation(&pod.ObjectMeta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	return false
}

func (a *PodMutatingWebhook) calculateVersion(pod *corev1.Pod) string {
	name := ""

	if len(pod.Spec.Containers) == 1 {
		image := strings.Split(pod.Spec.Containers[0].Image, ":")
		if len(image) > 1 && image[1] != "" && image[1] != "latest" {
			return image[1]
		}
	}

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
	version, _ := getLabelOrAnnotation(&pod.ObjectMeta, apicommon.VersionAnnotation, apicommon.K8sRecommendedVersionAnnotations)
	applicationName, _ := getLabelOrAnnotation(&pod.ObjectMeta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)

	var preDeploymentTasks []string
	var postDeploymentTasks []string
	var preDeploymentEvaluation []string
	var postDeploymentEvaluation []string

	if annotations, found := getLabelOrAnnotation(&pod.ObjectMeta, apicommon.PreDeploymentTaskAnnotation, ""); found {
		preDeploymentTasks = strings.Split(annotations, ",")
	}

	if annotations, found := getLabelOrAnnotation(&pod.ObjectMeta, apicommon.PostDeploymentTaskAnnotation, ""); found {
		postDeploymentTasks = strings.Split(annotations, ",")
	}

	if annotations, found := getLabelOrAnnotation(&pod.ObjectMeta, apicommon.PreDeploymentEvaluationAnnotation, ""); found {
		preDeploymentEvaluation = strings.Split(annotations, ",")
	}

	if annotations, found := getLabelOrAnnotation(&pod.ObjectMeta, apicommon.PostDeploymentEvaluationAnnotation, ""); found {
		postDeploymentEvaluation = strings.Split(annotations, ",")
	}

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	ownerRef := a.getOwnerReference(&pod.ObjectMeta)

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
	workloadName, _ := getLabelOrAnnotation(&pod.ObjectMeta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	applicationName, _ := getLabelOrAnnotation(&pod.ObjectMeta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)
	return operatorcommon.CreateResourceName(apicommon.MaxK8sObjectLength, apicommon.MinKLTNameLen, applicationName, workloadName)
}

func (a *PodMutatingWebhook) getAppName(pod *corev1.Pod) string {
	applicationName, _ := getLabelOrAnnotation(&pod.ObjectMeta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)
	return strings.ToLower(applicationName)
}

func (a *PodMutatingWebhook) getOwnerReference(resource *metav1.ObjectMeta) metav1.OwnerReference {
	reference := metav1.OwnerReference{}
	if len(resource.OwnerReferences) != 0 {
		for _, owner := range resource.OwnerReferences {
			if owner.Kind == "ReplicaSet" || owner.Kind == "Deployment" || owner.Kind == "StatefulSet" || owner.Kind == "DaemonSet" || owner.Kind == "Rollout" {
				reference.UID = owner.UID
				reference.Kind = owner.Kind
				reference.Name = owner.Name
				reference.APIVersion = owner.APIVersion
			}
		}
	}
	return reference
}

func getLabelOrAnnotation(resource *metav1.ObjectMeta, primaryAnnotation string, secondaryAnnotation string) (string, bool) {
	if resource.Annotations[primaryAnnotation] != "" {
		return resource.Annotations[primaryAnnotation], true
	}

	if resource.Labels[primaryAnnotation] != "" {
		return resource.Labels[primaryAnnotation], true
	}

	if secondaryAnnotation == "" {
		return "", false
	}

	if resource.Annotations[secondaryAnnotation] != "" {
		return resource.Annotations[secondaryAnnotation], true
	}

	if resource.Labels[secondaryAnnotation] != "" {
		return resource.Labels[secondaryAnnotation], true
	}
	return "", false
}

func setMapKey(myMap map[string]string, key, value string) {
	if myMap == nil {
		return
	}
	if value != "" {
		myMap[key] = value
	}
}
