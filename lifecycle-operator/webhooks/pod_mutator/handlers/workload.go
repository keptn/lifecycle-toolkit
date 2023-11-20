package handlers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	lifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type WorkloadHandler struct {
	Client      client.Client
	Log         logr.Logger
	EventSender eventsender.IEvent
}

func (a *WorkloadHandler) Handle(ctx context.Context, pod *corev1.Pod, namespace string) error {

	newWorkload := generateWorkload(ctx, pod, namespace)

	a.Log.Info("Searching for workload")

	workload := &lifecycle.KeptnWorkload{}
	err := a.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: newWorkload.Name}, workload)
	if errors.IsNotFound(err) {
		return a.createWorkload(ctx, newWorkload)
	}

	if err != nil {
		return fmt.Errorf("could not fetch Workload %w", err)
	}

	return a.updateWorkload(ctx, workload, newWorkload)
}

func (a *WorkloadHandler) updateWorkload(ctx context.Context, workload *lifecycle.KeptnWorkload, newWorkload *lifecycle.KeptnWorkload) error {
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
		return err
	}

	return nil
}

func (a *WorkloadHandler) createWorkload(ctx context.Context, newWorkload *lifecycle.KeptnWorkload) error {
	a.Log.Info("Creating workload", "workload", newWorkload.Name)
	err := a.Client.Create(ctx, newWorkload)
	if err != nil {
		a.Log.Error(err, "Could not create Workload")
		a.EventSender.Emit(apicommon.PhaseCreateWorkload, "Warning", newWorkload, apicommon.PhaseStateFailed, "could not create KeptnWorkload", newWorkload.Spec.Version)
		return err
	}

	return nil
}

func generateWorkload(ctx context.Context, pod *corev1.Pod, namespace string) *lifecycle.KeptnWorkload {
	version, _ := GetLabelOrAnnotation(&pod.ObjectMeta, apicommon.VersionAnnotation, apicommon.K8sRecommendedVersionAnnotations)
	version = strings.ToLower(version)
	preDeploymentTasks := getValuesForAnnotations(&pod.ObjectMeta, apicommon.PreDeploymentTaskAnnotation)
	postDeploymentTasks := getValuesForAnnotations(&pod.ObjectMeta, apicommon.PostDeploymentTaskAnnotation)
	preDeploymentEvaluation := getValuesForAnnotations(&pod.ObjectMeta, apicommon.PreDeploymentEvaluationAnnotation)
	postDeploymentEvaluation := getValuesForAnnotations(&pod.ObjectMeta, apicommon.PostDeploymentEvaluationAnnotation)
	applicationName := getAppName(&pod.ObjectMeta)
	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	ownerRef := GetOwnerReference(&pod.ObjectMeta)

	return &lifecycle.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:        getWorkloadName(&pod.ObjectMeta, applicationName),
			Namespace:   namespace,
			Annotations: traceContextCarrier,
			OwnerReferences: []metav1.OwnerReference{
				ownerRef,
			},
		},
		Spec: lifecycle.KeptnWorkloadSpec{
			AppName:                   applicationName,
			Version:                   version,
			ResourceReference:         lifecycle.ResourceReference{UID: ownerRef.UID, Kind: ownerRef.Kind, Name: ownerRef.Name},
			PreDeploymentTasks:        preDeploymentTasks,
			PostDeploymentTasks:       postDeploymentTasks,
			PreDeploymentEvaluations:  preDeploymentEvaluation,
			PostDeploymentEvaluations: postDeploymentEvaluation,
		},
	}
}
