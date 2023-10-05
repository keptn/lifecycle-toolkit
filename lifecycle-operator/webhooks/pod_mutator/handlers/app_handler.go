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

type AppHandler struct {
	Client      client.Client
	Log         logr.Logger
	Tracer      trace.Tracer
	EventSender controllercommon.IEvent
}

func (a *AppHandler) Handle(ctx context.Context, pod *corev1.Pod, namespace string) error {

	ctx, span := a.Tracer.Start(ctx, "create_app", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	newAppCreationRequest := generateAppCreationRequest(ctx, pod, namespace)
	newAppCreationRequest.SetSpanAttributes(span)

	a.Log.Info("Searching for AppCreationRequest", "appCreationRequest", newAppCreationRequest.Name)

	appCreationRequest := &klcv1alpha3.KeptnAppCreationRequest{}
	err := a.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: newAppCreationRequest.Name}, appCreationRequest)
	if errors.IsNotFound(err) {
		return a.createApp(ctx, newAppCreationRequest, span)
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("could not fetch AppCreationRequest"+": %+v", err)
	}

	return nil
}

func (a *AppHandler) createApp(ctx context.Context, newAppCreationRequest *klcv1alpha3.KeptnAppCreationRequest, span trace.Span) error {
	a.Log.Info("Creating app creation request", "appCreationRequest", newAppCreationRequest.Name)
	err := a.Client.Create(ctx, newAppCreationRequest)
	if err != nil {
		a.Log.Error(err, "Could not create AppCreationRequest")
		a.EventSender.Emit(apicommon.PhaseCreateAppCreationRequest, "Warning", newAppCreationRequest, apicommon.PhaseStateFailed, "could not create KeptnAppCreationRequest", newAppCreationRequest.Spec.AppName)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func generateAppCreationRequest(ctx context.Context, pod *corev1.Pod, namespace string) *klcv1alpha3.KeptnAppCreationRequest {

	appName := getAppName(&pod.ObjectMeta)

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

	if !isAppAnnotationPresent(pod) {
		kacr.Annotations[apicommon.AppTypeAnnotation] = string(apicommon.AppTypeSingleService)
	}

	return kacr
}

func isAppAnnotationPresent(pod *corev1.Pod) bool {
	_, gotAppAnnotation := GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.AppAnnotation, apicommon.K8sRecommendedAppAnnotations)

	if gotAppAnnotation {
		return true
	}

	if len(pod.Annotations) == 0 {
		pod.Annotations = make(map[string]string)
	}
	pod.Annotations[apicommon.AppAnnotation], _ = GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
	return false
}
