package handlers

import (
	"context"
	"fmt"

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

type AppCreationRequestHandler struct {
	Client      client.Client
	Log         logr.Logger
	Tracer      trace.Tracer
	EventSender controllercommon.IEvent
}

func (a *AppCreationRequestHandler) Handle(ctx context.Context, pod *corev1.Pod, namespace string) error {

	ctx, span := a.Tracer.Start(ctx, "create_appCreationRequest", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	newAppCreationRequest := generateResource(ctx, pod, namespace)
	newAppCreationRequest.SetSpanAttributes(span)

	a.Log.Info("Searching for AppCreationRequest", "appCreationRequest", newAppCreationRequest.Name, "namespace", newAppCreationRequest.Namespace)

	appCreationRequest := &klcv1alpha3.KeptnAppCreationRequest{}
	err := a.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: newAppCreationRequest.Name}, appCreationRequest)
	if errors.IsNotFound(err) {
		return a.createResource(ctx, newAppCreationRequest, span)
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return fmt.Errorf("could not fetch AppCreationRequest %w", err)
	}
	a.Log.Info("Found AppCreationRequest", "appCreationRequest", newAppCreationRequest.Name, "namespace", newAppCreationRequest.Namespace)
	return nil
}

func (a *AppCreationRequestHandler) createResource(ctx context.Context, newAppCreationRequest *klcv1alpha3.KeptnAppCreationRequest, span trace.Span) error {
	a.Log.Info("Creating app creation request", "appCreationRequest", newAppCreationRequest.Name, "namespace", newAppCreationRequest.Namespace)

	err := a.Client.Create(ctx, newAppCreationRequest)
	if err != nil {
		a.Log.Error(err, "Could not create AppCreationRequest")
		a.EventSender.Emit(apicommon.PhaseCreateAppCreationRequest, "Warning", newAppCreationRequest, apicommon.PhaseStateFailed, "could not create KeptnAppCreationRequest", newAppCreationRequest.Spec.AppName)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func generateResource(ctx context.Context, pod *corev1.Pod, namespace string) *klcv1alpha3.KeptnAppCreationRequest {

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	kacr := &klcv1alpha3.KeptnAppCreationRequest{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   namespace,
			Annotations: traceContextCarrier,
		},
	}

	if !isAppAnnotationPresent(&pod.ObjectMeta) {
		initEmptyAnnotations(&pod.ObjectMeta, 2)
		// at this point if the pod does not have an app annotation it means we create the app
		// and it will have a single workload
		pod.ObjectMeta.Annotations[apicommon.AppAnnotation] = pod.ObjectMeta.Annotations[apicommon.WorkloadAnnotation]
		// so we can mark the app request as single service type
		kacr.Annotations[apicommon.AppTypeAnnotation] = string(apicommon.AppTypeSingleService)
	}

	appName := getAppName(&pod.ObjectMeta)
	kacr.ObjectMeta.Name = appName
	kacr.Spec = klcv1alpha3.KeptnAppCreationRequestSpec{
		AppName: appName,
	}

	return kacr
}
