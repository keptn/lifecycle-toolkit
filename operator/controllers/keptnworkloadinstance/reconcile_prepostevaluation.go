package keptnworkloadinstance

import (
	"context"
	"fmt"
	"time"

	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/semconv"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *KeptnWorkloadInstanceReconciler) reconcilePrePostEvaluation(ctx context.Context, workloadInstance *klcv1alpha1.KeptnWorkloadInstance, checkType common.CheckType) (common.KeptnState, error) {
	newStatus, state, err := r.reconcileEvaluations(ctx, checkType, workloadInstance)
	if err != nil {
		return common.StateUnknown, err
	}
	overallState := common.GetOverallState(state)

	switch checkType {
	case common.PreDeploymentEvaluationCheckType:
		workloadInstance.Status.PreDeploymentEvaluationStatus = overallState
		workloadInstance.Status.PreDeploymentEvaluationTaskStatus = newStatus
	case common.PostDeploymentEvaluationCheckType:
		workloadInstance.Status.PostDeploymentEvaluationStatus = overallState
		workloadInstance.Status.PostDeploymentEvaluationTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return common.StateUnknown, err
	}
	return overallState, nil
}

func (r *KeptnWorkloadInstanceReconciler) reconcileEvaluations(ctx context.Context, checkType common.CheckType, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) ([]klcv1alpha1.EvaluationStatus, common.StatusSummary, error) {
	phase := common.KeptnPhaseType{
		ShortName: "ReconcileEvaluations",
		LongName:  "Reconcile Evaluations",
	}

	var evaluations []string
	var statuses []klcv1alpha1.EvaluationStatus

	switch checkType {
	case common.PreDeploymentEvaluationCheckType:
		evaluations = workloadInstance.Spec.PreDeploymentEvaluations
		statuses = workloadInstance.Status.PreDeploymentEvaluationTaskStatus
	case common.PostDeploymentEvaluationCheckType:
		evaluations = workloadInstance.Spec.PostDeploymentEvaluations
		statuses = workloadInstance.Status.PostDeploymentEvaluationTaskStatus
	}

	var summary common.StatusSummary
	summary.Total = len(evaluations)
	// Check current state of the PrePostEvaluationTasks
	var newStatus []klcv1alpha1.EvaluationStatus
	for _, evaluationName := range evaluations {
		var oldstatus common.KeptnState
		for _, ts := range statuses {
			if ts.EvaluationDefinitionName == evaluationName {
				oldstatus = ts.Status
			}
		}

		evaluationStatus := controllercommon.GetEvaluationStatus(evaluationName, statuses)
		evaluation := &klcv1alpha1.KeptnEvaluation{}
		evaluationExists := false

		if oldstatus != evaluationStatus.Status {
			controllercommon.RecordEvent(r.Recorder, phase, "Normal", workloadInstance, "EvaluationStatusChanged", fmt.Sprintf("evaluation status changed from %s to %s", oldstatus, evaluationStatus.Status), workloadInstance.GetVersion())
		}

		// Check if evaluation has already succeeded or failed
		if evaluationStatus.Status.IsCompleted() {
			newStatus = append(newStatus, evaluationStatus)
			continue
		}

		// Check if Evaluation is already created
		if evaluationStatus.EvaluationName != "" {
			err := r.Client.Get(ctx, types.NamespacedName{Name: evaluationStatus.EvaluationName, Namespace: workloadInstance.Namespace}, evaluation)
			if err != nil && errors.IsNotFound(err) {
				evaluationStatus.EvaluationName = ""
			} else if err != nil {
				return nil, summary, err
			}
			evaluationExists = true
		}

		// Create new Evaluation if it does not exist
		if !evaluationExists {
			evaluationName, err := r.createKeptnEvaluation(ctx, workloadInstance.Namespace, workloadInstance, evaluationName, checkType)
			if err != nil {
				return nil, summary, err
			}
			evaluationStatus.EvaluationName = evaluationName
			evaluationStatus.SetStartTime()
		} else {
			// Update state of Evaluation if it is already created
			evaluationStatus.Status = evaluation.Status.OverallStatus
			if evaluationStatus.Status.IsCompleted() {
				evaluationStatus.SetEndTime()
			}
		}
		// Update state of the Check
		newStatus = append(newStatus, evaluationStatus)
	}

	for _, ns := range newStatus {
		summary = common.UpdateStatusSummary(ns.Status, summary)
	}
	if common.GetOverallState(summary) != common.StateSucceeded {
		controllercommon.RecordEvent(r.Recorder, phase, "Warning", workloadInstance, "NotFinished", "has not finished", workloadInstance.GetVersion())
	}
	return newStatus, summary, nil
}

func (r *KeptnWorkloadInstanceReconciler) createKeptnEvaluation(ctx context.Context, namespace string, workloadInstance *klcv1alpha1.KeptnWorkloadInstance, evaluationDefinition string, checkType common.CheckType) (string, error) {

	ctx, span := r.Tracer.Start(ctx, "create_workload_evaluation", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	semconv.AddAttributeFromWorkloadInstance(span, *workloadInstance)

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	phase := common.KeptnPhaseType{
		ShortName: "KeptnEvaluationCreate",
		LongName:  "Keptn Evaluation Create",
	}

	newEvaluation := &klcv1alpha1.KeptnEvaluation{
		ObjectMeta: metav1.ObjectMeta{
			Name:        common.GenerateEvaluationName(checkType, evaluationDefinition),
			Namespace:   namespace,
			Annotations: traceContextCarrier,
		},
		Spec: klcv1alpha1.KeptnEvaluationSpec{
			WorkloadVersion:      workloadInstance.Spec.Version,
			Workload:             workloadInstance.Spec.WorkloadName,
			EvaluationDefinition: evaluationDefinition,
			Type:                 checkType,
			RetryInterval: metav1.Duration{
				Duration: 5 * time.Second,
			},
		},
	}
	err := controllerutil.SetControllerReference(workloadInstance, newEvaluation, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}
	err = r.Client.Create(ctx, newEvaluation)
	if err != nil {
		r.Log.Error(err, "could not create KeptnEvaluation")
		controllercommon.RecordEvent(r.Recorder, phase, "Warning", workloadInstance, "CreateFailed", "could not create KeptnEvaluation", workloadInstance.GetVersion())
		return "", err
	}
	controllercommon.RecordEvent(r.Recorder, phase, "Normal", workloadInstance, "Created", "created", workloadInstance.GetVersion())

	return newEvaluation.Name, nil
}
