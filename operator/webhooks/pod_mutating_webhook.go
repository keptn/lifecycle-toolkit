package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/semconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"hash/fnv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.keptn.sh,admissionReviewVersions=v1,sideEffects=None
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch
//+kubebuilder:rbac:groups=apps,resources=deployments;statefulsets;daemonsets;replicasets,verbs=get

// PodMutatingWebhook annotates Pods
type PodMutatingWebhook struct {
	Client   client.Client
	Tracer   trace.Tracer
	decoder  *admission.Decoder
	Recorder record.EventRecorder
	Log      logr.Logger
}

var ErrTooLongAnnotations = fmt.Errorf("too long annotations, maximum length for app and workload is 25 characters, for version 12 characters")

// Handle inspects incoming Pods and injects the Keptn scheduler if they contain the Keptn lifecycle annotations.
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

	if namespace.GetAnnotations()[common.NamespaceEnabledAnnotation] != "enabled" {
		logger.Info("namespace is not enabled for lifecycle controller", "namespace", req.Namespace)
		return admission.Allowed("namespace is not enabled for lifecycle controller")
	}

	logger.Info(fmt.Sprintf("Pod annotations: %v", pod.Annotations))

	podIsAnnotated, err := a.isPodAnnotated(pod)
	logger.Info("Checked if pod is annotated.")

	if err == nil && !podIsAnnotated {
		logger.Info("Pod is not annotated, check for parent annotations...")
		podIsAnnotated, err = a.copyAnnotationsIfParentAnnotated(ctx, &req, pod)
	}

	if err != nil {
		span.SetStatus(codes.Error, "Invalid annotations")
		return admission.Errored(http.StatusBadRequest, err)
	}
	if podIsAnnotated {
		logger.Info("Resource is annotated with Keptn annotations, using Keptn scheduler")
		pod.Spec.SchedulerName = "keptn-scheduler"
		logger.Info("Annotations", "annotations", pod.Annotations)

		isAppAnnotationPresent, err := a.isAppAnnotationPresent(pod)
		if err != nil {
			span.SetStatus(codes.Error, "Invalid annotations")
			return admission.Errored(http.StatusBadRequest, err)
		}
		if !isAppAnnotationPresent {
			if err := a.handleApp(ctx, logger, pod, req.Namespace); err != nil {
				logger.Error(err, "Could not handle App")
				span.SetStatus(codes.Error, err.Error())
				return admission.Errored(http.StatusBadRequest, err)
			}
		}
		semconv.AddAttributeFromAnnotations(span, pod.Annotations)

		logger.Info("Attributes from annotations set")

		if err := a.handleWorkload(ctx, logger, pod, req.Namespace); err != nil {
			logger.Error(err, "Could not handle Workload")
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

func (a *PodMutatingWebhook) isPodAnnotated(pod *corev1.Pod) (bool, error) {
	workload, gotWorkloadAnnotation := getLabelOrAnnotation(&pod.ObjectMeta, common.WorkloadAnnotation, common.K8sRecommendedWorkloadAnnotations)
	version, gotVersionAnnotation := getLabelOrAnnotation(&pod.ObjectMeta, common.VersionAnnotation, common.K8sRecommendedVersionAnnotations)

	if len(workload) > common.MaxWorkloadNameLength || len(version) > common.MaxVersionLength {
		return false, ErrTooLongAnnotations
	}

	if gotWorkloadAnnotation {
		if !gotVersionAnnotation {
			if len(pod.Annotations) == 0 {
				pod.Annotations = make(map[string]string)
			}
			pod.Annotations[common.VersionAnnotation] = a.calculateVersion(pod)
		}
		return true, nil
	}
	return false, nil
}

func (a *PodMutatingWebhook) copyAnnotationsIfParentAnnotated(ctx context.Context, req *admission.Request, pod *corev1.Pod) (bool, error) {
	podOwner := a.getOwnerReference(&pod.ObjectMeta)
	if podOwner.UID == "" {
		return false, nil
	}

	switch podOwner.Kind {
	case "ReplicaSet":
		rs := &appsv1.ReplicaSet{}
		if err := a.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: podOwner.Name}, rs); err != nil {
			return false, nil
		}
		a.Log.Info("Done fetching RS")

		rsOwner := a.getOwnerReference(&rs.ObjectMeta)
		if rsOwner.UID == "" {
			return false, nil
		}

		dp := &appsv1.Deployment{}
		if err := a.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: rsOwner.Name}, dp); err != nil {
			return false, nil
		}

		return a.copyResourceLabelsIfPresent(&dp.ObjectMeta, pod)
	case "StatefulSet":
		sts := &appsv1.StatefulSet{}
		if err := a.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: podOwner.Name}, sts); err != nil {
			return false, nil
		}
		return a.copyResourceLabelsIfPresent(&sts.ObjectMeta, pod)
	case "Daemonset":
		ds := &appsv1.DaemonSet{}
		if err := a.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: podOwner.Name}, ds); err != nil {
			return false, nil
		}
		return a.copyResourceLabelsIfPresent(&ds.ObjectMeta, pod)
	default:
		return false, nil
	}
}

func (a *PodMutatingWebhook) copyResourceLabelsIfPresent(sourceResource *metav1.ObjectMeta, targetPod *corev1.Pod) (bool, error) {
	var workloadName, appName, version, preDeploymentChecks, postDeploymentChecks, preEvaluationChecks, postEvaluationChecks string
	var gotWorkloadName, gotAppName, gotVersion, gotPreDeploymentChecks, gotPostDeploymentChecks, gotPreEvaluationChecks, gotPostEvaluationChecks bool

	workloadName, gotWorkloadName = getLabelOrAnnotation(sourceResource, common.WorkloadAnnotation, common.K8sRecommendedWorkloadAnnotations)
	appName, gotAppName = getLabelOrAnnotation(sourceResource, common.AppAnnotation, common.K8sRecommendedAppAnnotations)
	version, gotVersion = getLabelOrAnnotation(sourceResource, common.VersionAnnotation, common.K8sRecommendedVersionAnnotations)
	preDeploymentChecks, gotPreDeploymentChecks = getLabelOrAnnotation(sourceResource, common.PreDeploymentTaskAnnotation, "")
	postDeploymentChecks, gotPostDeploymentChecks = getLabelOrAnnotation(sourceResource, common.PostDeploymentTaskAnnotation, "")
	preEvaluationChecks, gotPreEvaluationChecks = getLabelOrAnnotation(sourceResource, common.PreDeploymentEvaluationAnnotation, "")
	postEvaluationChecks, gotPostEvaluationChecks = getLabelOrAnnotation(sourceResource, common.PostDeploymentEvaluationAnnotation, "")

	if len(workloadName) > common.MaxWorkloadNameLength || len(version) > common.MaxVersionLength {
		return false, ErrTooLongAnnotations
	}

	if len(targetPod.Annotations) == 0 {
		targetPod.Annotations = make(map[string]string)
	}

	if gotWorkloadName {
		targetPod.Annotations[common.WorkloadAnnotation] = workloadName

		if !gotVersion {
			targetPod.Annotations[common.VersionAnnotation] = a.calculateVersion(targetPod)
		} else {
			targetPod.Annotations[common.VersionAnnotation] = version
		}

		if gotAppName {
			targetPod.Annotations[common.AppAnnotation] = appName
		}

		if gotPreDeploymentChecks {
			targetPod.Annotations[common.PreDeploymentTaskAnnotation] = preDeploymentChecks
		}

		if gotPostDeploymentChecks {
			targetPod.Annotations[common.PostDeploymentTaskAnnotation] = postDeploymentChecks
		}

		if gotPreEvaluationChecks {
			targetPod.Annotations[common.PreDeploymentEvaluationAnnotation] = preEvaluationChecks
		}

		if gotPostEvaluationChecks {
			targetPod.Annotations[common.PostDeploymentEvaluationAnnotation] = postEvaluationChecks
		}
		return true, nil
	}
	return false, nil
}

func (a *PodMutatingWebhook) isAppAnnotationPresent(pod *corev1.Pod) (bool, error) {
	app, gotAppAnnotation := getLabelOrAnnotation(&pod.ObjectMeta, common.AppAnnotation, common.K8sRecommendedAppAnnotations)

	if gotAppAnnotation {
		if len(app) > common.MaxAppNameLength {
			return false, ErrTooLongAnnotations
		}
		return true, nil
	}

	if len(pod.Annotations) == 0 {
		pod.Annotations = make(map[string]string)
	}
	pod.Annotations[common.AppAnnotation], _ = getLabelOrAnnotation(&pod.ObjectMeta, common.WorkloadAnnotation, common.K8sRecommendedWorkloadAnnotations)
	return false, nil
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

func (a *PodMutatingWebhook) handleWorkload(ctx context.Context, logger logr.Logger, pod *corev1.Pod, namespace string) error {

	ctx, span := a.Tracer.Start(ctx, "create_workload", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	newWorkload := a.generateWorkload(ctx, pod, namespace)

	semconv.AddAttributeFromWorkload(span, *newWorkload)

	logger.Info("Searching for workload")

	workload := &klcv1alpha1.KeptnWorkload{}
	err := a.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: newWorkload.Name}, workload)
	if errors.IsNotFound(err) {
		logger.Info("Creating workload", "workload", workload.Name)
		workload = newWorkload
		err = a.Client.Create(ctx, workload)
		if err != nil {
			logger.Error(err, "Could not create Workload")
			a.Recorder.Event(workload, "Warning", "WorkloadNotCreated", fmt.Sprintf("Could not create KeptnWorkload / Namespace: %s, Name: %s ", workload.Namespace, workload.Name))
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		a.Recorder.Event(workload, "Normal", "WorkloadCreated", fmt.Sprintf("KeptnWorkload created / Namespace: %s, Name: %s ", workload.Namespace, workload.Name))
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
		a.Recorder.Event(workload, "Warning", "WorkloadNotUpdated", fmt.Sprintf("Could not update KeptnWorkload / Namespace: %s, Name: %s ", workload.Namespace, workload.Name))
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	a.Recorder.Event(workload, "Normal", "WorkloadUpdated", fmt.Sprintf("KeptnWorkload updated / Namespace: %s, Name: %s ", workload.Namespace, workload.Name))

	return nil
}

func (a *PodMutatingWebhook) handleApp(ctx context.Context, logger logr.Logger, pod *corev1.Pod, namespace string) error {

	ctx, span := a.Tracer.Start(ctx, "create_app", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	newApp := a.generateApp(ctx, pod, namespace)

	semconv.AddAttributeFromApp(span, *newApp)

	logger.Info("Searching for app")

	app := &klcv1alpha1.KeptnApp{}
	err := a.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: newApp.Name}, app)
	if errors.IsNotFound(err) {
		logger.Info("Creating app", "app", app.Name)
		app = newApp
		err = a.Client.Create(ctx, app)
		if err != nil {
			logger.Error(err, "Could not create App")
			a.Recorder.Event(app, "Warning", "AppNotCreated", fmt.Sprintf("Could not create KeptnApp / Namespace: %s, Name: %s ", app.Namespace, app.Name))
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		a.Recorder.Event(app, "Normal", "AppCreated", fmt.Sprintf("KeptnApp created / Namespace: %s, Name: %s ", app.Namespace, app.Name))
		return nil
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("could not fetch App"+": %+v", err)
	}

	if reflect.DeepEqual(app.Spec, newApp.Spec) {
		logger.Info("Pod not changed, not updating anything")
		return nil
	}

	logger.Info("Pod changed, updating app")
	app.Spec = newApp.Spec

	err = a.Client.Update(ctx, app)
	if err != nil {
		logger.Error(err, "Could not update App")
		a.Recorder.Event(app, "Warning", "AppNotUpdated", fmt.Sprintf("Could not update KeptnApp / Namespace: %s, Name: %s ", app.Namespace, app.Name))
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	a.Recorder.Event(app, "Normal", "AppUpdated", fmt.Sprintf("KeptnApp updated / Namespace: %s, Name: %s ", app.Namespace, app.Name))

	return nil
}

func (a *PodMutatingWebhook) generateWorkload(ctx context.Context, pod *corev1.Pod, namespace string) *klcv1alpha1.KeptnWorkload {
	version, _ := getLabelOrAnnotation(&pod.ObjectMeta, common.VersionAnnotation, common.K8sRecommendedVersionAnnotations)
	applicationName, _ := getLabelOrAnnotation(&pod.ObjectMeta, common.AppAnnotation, common.K8sRecommendedAppAnnotations)

	var preDeploymentTasks []string
	var postDeploymentTasks []string
	var preDeploymentEvaluation []string
	var postDeploymentEvaluation []string

	if annotations, found := getLabelOrAnnotation(&pod.ObjectMeta, common.PreDeploymentTaskAnnotation, ""); found {
		preDeploymentTasks = strings.Split(annotations, ",")
	}

	if annotations, found := getLabelOrAnnotation(&pod.ObjectMeta, common.PostDeploymentTaskAnnotation, ""); found {
		postDeploymentTasks = strings.Split(annotations, ",")
	}

	if annotations, found := getLabelOrAnnotation(&pod.ObjectMeta, common.PreDeploymentEvaluationAnnotation, ""); found {
		preDeploymentEvaluation = strings.Split(annotations, ",")
	}

	if annotations, found := getLabelOrAnnotation(&pod.ObjectMeta, common.PostDeploymentEvaluationAnnotation, ""); found {
		postDeploymentEvaluation = strings.Split(annotations, ",")
	}

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	ownerRef := a.getOwnerReference(&pod.ObjectMeta)

	return &klcv1alpha1.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:        a.getWorkloadName(pod),
			Namespace:   namespace,
			Annotations: traceContextCarrier,
		},
		Spec: klcv1alpha1.KeptnWorkloadSpec{
			AppName:                   applicationName,
			Version:                   version,
			ResourceReference:         klcv1alpha1.ResourceReference{UID: ownerRef.UID, Kind: ownerRef.Kind},
			PreDeploymentTasks:        preDeploymentTasks,
			PostDeploymentTasks:       postDeploymentTasks,
			PreDeploymentEvaluations:  preDeploymentEvaluation,
			PostDeploymentEvaluations: postDeploymentEvaluation,
		},
	}
}

func (a *PodMutatingWebhook) generateApp(ctx context.Context, pod *corev1.Pod, namespace string) *klcv1alpha1.KeptnApp {
	version, _ := getLabelOrAnnotation(&pod.ObjectMeta, common.VersionAnnotation, common.K8sRecommendedVersionAnnotations)
	appName := a.getAppName(pod)

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	return &klcv1alpha1.KeptnApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:        appName,
			Namespace:   namespace,
			Annotations: traceContextCarrier,
		},
		Spec: klcv1alpha1.KeptnAppSpec{
			Version:                   version,
			PreDeploymentTasks:        []string{},
			PostDeploymentTasks:       []string{},
			PreDeploymentEvaluations:  []string{},
			PostDeploymentEvaluations: []string{},
			Workloads: []klcv1alpha1.KeptnWorkloadRef{
				{
					Name:    appName,
					Version: version,
				},
			},
		},
	}
}

func (a *PodMutatingWebhook) getWorkloadName(pod *corev1.Pod) string {
	workloadName, _ := getLabelOrAnnotation(&pod.ObjectMeta, common.WorkloadAnnotation, common.K8sRecommendedWorkloadAnnotations)
	applicationName, _ := getLabelOrAnnotation(&pod.ObjectMeta, common.AppAnnotation, common.K8sRecommendedAppAnnotations)
	return strings.ToLower(applicationName + "-" + workloadName)
}

func (a *PodMutatingWebhook) getAppName(pod *corev1.Pod) string {
	applicationName, _ := getLabelOrAnnotation(&pod.ObjectMeta, common.AppAnnotation, common.K8sRecommendedAppAnnotations)
	return strings.ToLower(applicationName)
}

func (a *PodMutatingWebhook) getOwnerReference(resource *metav1.ObjectMeta) metav1.OwnerReference {
	reference := metav1.OwnerReference{}
	if len(resource.OwnerReferences) != 0 {
		for _, owner := range resource.OwnerReferences {
			if owner.Kind == "ReplicaSet" || owner.Kind == "Deployment" || owner.Kind == "StatefulSet" || owner.Kind == "DaemonSet" {
				reference.UID = owner.UID
				reference.Kind = owner.Kind
				reference.Name = owner.Name
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
