package keptnappversion

import (
	"context"
	"fmt"
	"time"

	klcv1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-controller/operator/api/v1alpha1/common"
	"github.com/keptn/lifecycle-controller/operator/api/v1alpha1/semconv"
	controllercommon "github.com/keptn/lifecycle-controller/operator/controllers/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *KeptnAppVersionReconciler) reconcilePrePostEvaluation(ctx context.Context, appVersion *klcv1alpha1.KeptnAppVersion, checkType common.CheckType) (common.KeptnState, error) {
	newStatus, state, err := r.reconcileEvaluations(ctx, checkType, appVersion)
	if err != nil {
		return common.StateUnknown, err
	}
	overallState := common.GetOverallState(state)

	switch checkType {
	case common.PreDeploymentEvaluationCheckType:
		appVersion.Status.PreDeploymentEvaluationStatus = overallState
		appVersion.Status.PreDeploymentEvaluationTaskStatus = newStatus
	case common.PostDeploymentEvaluationCheckType:
		appVersion.Status.PostDeploymentEvaluationStatus = overallState
		appVersion.Status.PostDeploymentEvaluationTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, appVersion)
	if err != nil {
		return common.StateUnknown, err
	}
	return overallState, nil
}

func (r *KeptnAppVersionReconciler) reconcileEvaluations(ctx context.Context, checkType common.CheckType, appVersion *klcv1alpha1.KeptnAppVersion) ([]klcv1alpha1.EvaluationStatus, common.StatusSummary, error) {
	phase := common.KeptnPhaseType{
		ShortName: "ReconcileEvaluations",
		LongName:  "Reconcile Evaluations",
	}

	var evaluations []string
	var statuses []klcv1alpha1.EvaluationStatus

	switch checkType {
	case common.PreDeploymentEvaluationCheckType:
		evaluations = appVersion.Spec.PreDeploymentEvaluations
		statuses = appVersion.Status.PreDeploymentEvaluationTaskStatus
	case common.PostDeploymentEvaluationCheckType:
		evaluations = appVersion.Spec.PostDeploymentEvaluations
		statuses = appVersion.Status.PostDeploymentEvaluationTaskStatus
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
			controllercommon.RecordEvent(r.Recorder, phase, "Normal", appVersion, "EvaluationStatusChanged", fmt.Sprintf("evaluation status changed from %s to %s", oldstatus, evaluationStatus.Status), appVersion.GetVersion())
		}

		// Check if evaluation has already succeeded or failed
		if evaluationStatus.Status.IsCompleted() {
			newStatus = append(newStatus, evaluationStatus)
			continue
		}

		// Check if Evaluation is already created
		if evaluationStatus.EvaluationName != "" {
			err := r.Client.Get(ctx, types.NamespacedName{Name: evaluationStatus.EvaluationName, Namespace: appVersion.Namespace}, evaluation)
			if err != nil && errors.IsNotFound(err) {
				evaluationStatus.EvaluationName = ""
			} else if err != nil {
				return nil, summary, err
			}
			evaluationExists = true
		}

		// Create new Evaluation if it does not exist
		if !evaluationExists {
			evaluationName, err := r.createKeptnEvaluation(ctx, appVersion.Namespace, appVersion, evaluationName, checkType)
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
		controllercommon.RecordEvent(r.Recorder, phase, "Warning", appVersion, "NotFinished", "has not finished", appVersion.GetVersion())
	}
	return newStatus, summary, nil
}

func (r *KeptnAppVersionReconciler) createKeptnEvaluation(ctx context.Context, namespace string, appVersion *klcv1alpha1.KeptnAppVersion, evaluationDefinition string, checkType common.CheckType) (string, error) {

	ctx, span := r.Tracer.Start(ctx, "create_app_evaluation", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	semconv.AddAttributeFromAppVersion(span, *appVersion)

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
			AppVersion:           appVersion.Spec.Version,
			AppName:              appVersion.Spec.AppName,
			EvaluationDefinition: evaluationDefinition,
			Type:                 checkType,
			RetryInterval: metav1.Duration{
				Duration: 5 * time.Second,
			},
		},
	}
	err := controllerutil.SetControllerReference(appVersion, newEvaluation, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}
	err = r.Client.Create(ctx, newEvaluation)
	if err != nil {
		r.Log.Error(err, "could not create KeptnEvaluation")
		controllercommon.RecordEvent(r.Recorder, phase, "Warning", appVersion, "CreateFailed", "could not create KeptnEvaluation", appVersion.GetVersion())
		return "", err
	}
	controllercommon.RecordEvent(r.Recorder, phase, "Normal", appVersion, "Created", "created", appVersion.GetVersion())

	return newEvaluation.Name, nil
}
