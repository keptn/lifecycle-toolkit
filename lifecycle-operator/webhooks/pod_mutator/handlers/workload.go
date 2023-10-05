package handlers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type WorkloadHandler struct {
	Client      client.Client
	Log         logr.Logger
	Tracer      trace.Tracer
	EventSender controllercommon.IEvent
}

func (a *WorkloadHandler) Handle(ctx context.Context, pod *corev1.Pod, namespace string) error {

	ctx, span := a.Tracer.Start(ctx, "create_workload", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	newWorkload := generateWorkload(ctx, pod, namespace)
	newWorkload.SetSpanAttributes(span)

	a.Log.Info("Searching for workload")

	workload := &klcv1alpha3.KeptnWorkload{}
	err := a.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: newWorkload.Name}, workload)
	if errors.IsNotFound(err) {
		return a.createWorkload(ctx, newWorkload, span)
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("could not fetch Workload"+": %+v", err)
	}

	return a.updateWorkload(ctx, workload, newWorkload, span)
}

func (a *WorkloadHandler) updateWorkload(ctx context.Context, workload *klcv1alpha3.KeptnWorkload, newWorkload *klcv1alpha3.KeptnWorkload, span trace.Span) error {
	if reflect.DeepEqual(workload.Spec, newWorkload.Spec) {
		a.Log.Info("Pod not changed, not updating anything")
		return nil
	}

	a.Log.Info("Pod changed, updating workload")
	workload.Spec = newWorkload.Spec

	err := a.Client.Update(ctx, workload)
	if err != nil {
		a.Log.Error(err, "Could not update Workload")
		a.EventSender.Emit(apicommon.PhaseUpdateWorkload, "Warning", workload, apicommon.PhaseStateFailed, "could not update KeptnWorkload", workload.Spec.Version)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (a *WorkloadHandler) createWorkload(ctx context.Context, newWorkload *klcv1alpha3.KeptnWorkload, span trace.Span) error {
	a.Log.Info("Creating workload", "workload", newWorkload.Name)
	err := a.Client.Create(ctx, newWorkload)
	if err != nil {
		a.Log.Error(err, "Could not create Workload")
		a.EventSender.Emit(apicommon.PhaseCreateWorkload, "Warning", newWorkload, apicommon.PhaseStateFailed, "could not create KeptnWorkload", newWorkload.Spec.Version)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func generateWorkload(ctx context.Context, pod *corev1.Pod, namespace string) *klcv1alpha3.KeptnWorkload {
	version := getVersion(&pod.ObjectMeta)
	preDeploymentTasks := getAnnotations(&pod.ObjectMeta, apicommon.PreDeploymentTaskAnnotation)
	postDeploymentTasks := getAnnotations(&pod.ObjectMeta, apicommon.PostDeploymentTaskAnnotation)
	preDeploymentEvaluation := getAnnotations(&pod.ObjectMeta, apicommon.PreDeploymentEvaluationAnnotation)
	postDeploymentEvaluation := getAnnotations(&pod.ObjectMeta, apicommon.PostDeploymentEvaluationAnnotation)
	applicationName := getAppName(&pod.ObjectMeta)
	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	ownerRef := GetOwnerReference(&pod.ObjectMeta)

	return &klcv1alpha3.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getWorkloadName(&pod.ObjectMeta, applicationName),
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

func getAnnotations(objMeta *metav1.ObjectMeta, annotationKey string) []string {
	if annotations, found := GetLabelOrAnnotation(objMeta, annotationKey, ""); found {
		return strings.Split(annotations, ",")
	}
	return nil
}
