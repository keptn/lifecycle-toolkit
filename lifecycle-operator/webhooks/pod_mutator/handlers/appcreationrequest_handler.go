package handlers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AppCreationRequestHandler struct {
	Client      client.Client
	Log         logr.Logger
	EventSender eventsender.IEvent
}

func (a *AppCreationRequestHandler) Handle(ctx context.Context, pod *corev1.Pod, namespace string) error {
	newAppCreationRequest := generateResource(ctx, pod, namespace)

	a.Log.Info("Searching for AppCreationRequest", "appCreationRequest", newAppCreationRequest.Name, "namespace", newAppCreationRequest.Namespace)

	appCreationRequest := &apilifecycle.KeptnAppCreationRequest{}
	err := a.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: newAppCreationRequest.Name}, appCreationRequest)
	if errors.IsNotFound(err) {
		return a.createResource(ctx, newAppCreationRequest)
	}

	if err != nil {
		return fmt.Errorf("could not fetch AppCreationRequest %w", err)
	}
	a.Log.Info("Found AppCreationRequest", "appCreationRequest", newAppCreationRequest.Name, "namespace", newAppCreationRequest.Namespace)
	return nil
}

func (a *AppCreationRequestHandler) createResource(ctx context.Context, newAppCreationRequest *apilifecycle.KeptnAppCreationRequest) error {
	a.Log.Info("Creating app creation request", "appCreationRequest", newAppCreationRequest.Name, "namespace", newAppCreationRequest.Namespace)

	err := a.Client.Create(ctx, newAppCreationRequest)
	if err != nil {
		a.Log.Error(err, "Could not create AppCreationRequest")
		a.EventSender.Emit(apicommon.PhaseCreateAppCreationRequest, "Warning", newAppCreationRequest, apicommon.PhaseStateFailed, "could not create KeptnAppCreationRequest", newAppCreationRequest.Spec.AppName)
		return err
	}

	return nil
}

func generateResource(ctx context.Context, pod *corev1.Pod, namespace string) *apilifecycle.KeptnAppCreationRequest {

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	kacr := &apilifecycle.KeptnAppCreationRequest{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   namespace,
			Annotations: traceContextCarrier,
		},
	}

	if !isAppAnnotationPresent(&pod.ObjectMeta) {
		initEmptyAnnotations(&pod.ObjectMeta, 2)
		// at this point if the pod does not have an app annotation it means we create the app
		// and it will have a single workload
		appName, _ := GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.WorkloadAnnotation, apicommon.K8sRecommendedWorkloadAnnotations)
		pod.Annotations[apicommon.AppAnnotation] = appName
		// so we can mark the app request as single service type
		kacr.Annotations[apicommon.AppTypeAnnotation] = string(apicommon.AppTypeSingleService)
	}

	appName := getAppName(&pod.ObjectMeta)
	kacr.ObjectMeta.Name = appName
	kacr.Spec = apilifecycle.KeptnAppCreationRequestSpec{
		AppName: appName,
	}

	return kacr
}
