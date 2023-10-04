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

func generateWorkload(ctx context.Context, pod *corev1.Pod, namespace string) *klcv1alpha3.KeptnWorkload {
	version := getVersion(&pod.ObjectMeta)
	applicationName := getAppName(&pod.ObjectMeta)

	var preDeploymentTasks []string
	var postDeploymentTasks []string
	var preDeploymentEvaluation []string
	var postDeploymentEvaluation []string

	if annotations, found := GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.PreDeploymentTaskAnnotation, ""); found {
		preDeploymentTasks = strings.Split(annotations, ",")
	}

	if annotations, found := GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.PostDeploymentTaskAnnotation, ""); found {
		postDeploymentTasks = strings.Split(annotations, ",")
	}

	if annotations, found := GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.PreDeploymentEvaluationAnnotation, ""); found {
		preDeploymentEvaluation = strings.Split(annotations, ",")
	}

	if annotations, found := GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.PostDeploymentEvaluationAnnotation, ""); found {
		postDeploymentEvaluation = strings.Split(annotations, ",")
	}

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	ownerRef := GetOwnerReference(&pod.ObjectMeta)

	return &klcv1alpha3.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getWorkloadName(&pod.ObjectMeta),
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

func (a *WorkloadHandler) HandleWorkload(ctx context.Context, pod *corev1.Pod, namespace string) error {

	ctx, span := a.Tracer.Start(ctx, "create_workload", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	newWorkload := generateWorkload(ctx, pod, namespace)

	newWorkload.SetSpanAttributes(span)

	a.Log.Info("Searching for workload")

	workload := &klcv1alpha3.KeptnWorkload{}
	err := a.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: newWorkload.Name}, workload)
	if errors.IsNotFound(err) {
		a.Log.Info("Creating workload", "workload", workload.Name)
		workload = newWorkload
		err = a.Client.Create(ctx, workload)
		if err != nil {
			a.Log.Error(err, "Could not create Workload")
			a.EventSender.Emit(apicommon.PhaseCreateWorkload, "Warning", workload, apicommon.PhaseStateFailed, "could not create KeptnWorkload", workload.Spec.Version)
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		return nil
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("could not fetch Workload"+": %+v", err)
	}

	if reflect.DeepEqual(workload.Spec, newWorkload.Spec) {
		a.Log.Info("Pod not changed, not updating anything")
		return nil
	}

	a.Log.Info("Pod changed, updating workload")
	workload.Spec = newWorkload.Spec

	err = a.Client.Update(ctx, workload)
	if err != nil {
		a.Log.Error(err, "Could not update Workload")
		a.EventSender.Emit(apicommon.PhaseUpdateWorkload, "Warning", workload, apicommon.PhaseStateFailed, "could not update KeptnWorkload", workload.Spec.Version)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
