package common

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/interfaces"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type EvaluationHandler struct {
	client.Client
	Recorder    record.EventRecorder
	Log         logr.Logger
	Tracer      trace.Tracer
	Scheme      *runtime.Scheme
	SpanHandler ISpanHandler
}

type EvaluationCreateAttributes struct {
	SpanName             string
	EvaluationDefinition string
	CheckType            apicommon.CheckType
}

func (r EvaluationHandler) ReconcileEvaluations(ctx context.Context, phaseCtx context.Context, reconcileObject client.Object, evaluationCreateAttributes EvaluationCreateAttributes) ([]klcv1alpha1.EvaluationStatus, apicommon.StatusSummary, error) {
	piWrapper, err := interfaces.NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return nil, apicommon.StatusSummary{}, err
	}

	evaluations, statuses := r.setupEvaluations(evaluationCreateAttributes, piWrapper)

	var summary apicommon.StatusSummary
	summary.Total = len(evaluations)
	// Check current state of the PrePostEvaluationTasks
	newStatus, evaluationStatuses, statusSummary, err2 := r.handlePrePostEvaluations(ctx, phaseCtx, reconcileObject, evaluationCreateAttributes, evaluations, statuses, phase, piWrapper, summary)
	if err2 != nil {
		return evaluationStatuses, statusSummary, err2
	}

	for _, ns := range newStatus {
		summary = apicommon.UpdateStatusSummary(ns.Status, summary)
	}
	if apicommon.GetOverallState(summary) != apicommon.StateSucceeded {
		RecordEvent(r.Recorder, phase, "Warning", reconcileObject, "NotFinished", "has not finished", piWrapper.GetVersion())
	}
	return newStatus, summary, nil
}

func (r EvaluationHandler) handlePrePostEvaluations(ctx context.Context, phaseCtx context.Context, reconcileObject client.Object, evaluationCreateAttributes EvaluationCreateAttributes, evaluations []string, statuses []klcv1alpha1.EvaluationStatus, phase apicommon.KeptnPhaseType, piWrapper *interfaces.PhaseItemWrapper, summary apicommon.StatusSummary) ([]klcv1alpha1.EvaluationStatus, []klcv1alpha1.EvaluationStatus, apicommon.StatusSummary, error) {
	var newStatus []klcv1alpha1.EvaluationStatus
	for _, evaluationName := range evaluations {
		oldStatus := r.findOldEvaluationStatus(statuses, evaluationName)

		evaluationStatus := GetEvaluationStatus(evaluationName, statuses)
		evaluation := &klcv1alpha1.KeptnEvaluation{}

		if oldStatus != evaluationStatus.Status {
			RecordEvent(r.Recorder, phase, "Normal", reconcileObject, "EvaluationStatusChanged", fmt.Sprintf("evaluation status changed from %s to %s", oldStatus, evaluationStatus.Status), piWrapper.GetVersion())
		}

		// Check if evaluation has already succeeded or failed
		if evaluationStatus.Status.IsCompleted() {
			newStatus = append(newStatus, evaluationStatus)
			continue
		}

		// Check if Evaluation is already created
		evaluationExists, evaluationStatuses, i, err := r.checkAlreadyCreated(ctx, evaluationStatus, piWrapper, evaluation, &summary)
		if err != nil {
			return evaluationStatuses, i, summary, err
		}

		// Create new Evaluation if it does not exist
		if !evaluationExists {
			statusSummary, err := r.createEvaluation(ctx, phaseCtx, reconcileObject, evaluationCreateAttributes, evaluationName, piWrapper, summary, evaluationStatus, evaluation)
			if err != nil {
				return nil, nil, statusSummary, err
			}
		} else {
			_, spanEvaluationTrace, err := r.SpanHandler.GetSpan(phaseCtx, r.Tracer, evaluation, "")
			if err != nil {
				r.Log.Error(err, "could not get span")
			}
			// Update state of Evaluation if it is already created
			evaluationStatus.Status = evaluation.Status.OverallStatus
			r.updateEvaluation(evaluationStatus, spanEvaluationTrace, evaluation)
		}
		// Update state of the Check
		newStatus = append(newStatus, evaluationStatus)
	}
	return newStatus, nil, apicommon.StatusSummary{}, nil
}

func (r EvaluationHandler) updateEvaluation(evaluationStatus klcv1alpha1.EvaluationStatus, spanEvaluationTrace trace.Span, evaluation *klcv1alpha1.KeptnEvaluation) {
	if evaluationStatus.Status.IsCompleted() {
		if evaluationStatus.Status.IsSucceeded() {
			spanEvaluationTrace.AddEvent(evaluation.Name + " has finished")
			spanEvaluationTrace.SetStatus(codes.Ok, "Finished")
		} else {
			spanEvaluationTrace.AddEvent(evaluation.Name + " has failed")
			spanEvaluationTrace.SetStatus(codes.Error, "Failed")
		}
		spanEvaluationTrace.End()
		if err := r.SpanHandler.UnbindSpan(evaluation, ""); err != nil {
			r.Log.Error(err, controllererrors.ErrCouldNotUnbindSpan, evaluation.Name)
		}
		evaluationStatus.SetEndTime()
	}
}

func (r EvaluationHandler) createEvaluation(ctx context.Context, phaseCtx context.Context, reconcileObject client.Object, evaluationCreateAttributes EvaluationCreateAttributes, evaluationName string, piWrapper *interfaces.PhaseItemWrapper, summary apicommon.StatusSummary, evaluationStatus klcv1alpha1.EvaluationStatus, evaluation *klcv1alpha1.KeptnEvaluation) (apicommon.StatusSummary, error) {
	evaluationCreateAttributes.EvaluationDefinition = evaluationName
	evaluationName, err := r.CreateKeptnEvaluation(ctx, piWrapper.GetNamespace(), reconcileObject, evaluationCreateAttributes)
	if err != nil {
		return summary, err
	}
	evaluationStatus.EvaluationName = evaluationName
	evaluationStatus.SetStartTime()
	_, _, err = r.SpanHandler.GetSpan(phaseCtx, r.Tracer, evaluation, "")
	if err != nil {
		r.Log.Error(err, "could not get span")
	}
	return apicommon.StatusSummary{}, nil
}

func (r EvaluationHandler) checkAlreadyCreated(ctx context.Context, evaluationStatus klcv1alpha1.EvaluationStatus, piWrapper *interfaces.PhaseItemWrapper, evaluation *klcv1alpha1.KeptnEvaluation, summary *apicommon.StatusSummary) (bool, []klcv1alpha1.EvaluationStatus, []klcv1alpha1.EvaluationStatus, error) {
	evaluationExists := false
	if evaluationStatus.EvaluationName != "" {
		err := r.Client.Get(ctx, types.NamespacedName{Name: evaluationStatus.EvaluationName, Namespace: piWrapper.GetNamespace()}, evaluation)
		if err != nil && errors.IsNotFound(err) {
			evaluationStatus.EvaluationName = ""
		} else if err != nil {
			return false, nil, nil, err
		}
		evaluationExists = true
	}
	return evaluationExists, nil, nil, nil
}

func (r EvaluationHandler) findOldEvaluationStatus(statuses []klcv1alpha1.EvaluationStatus, evaluationName string) apicommon.KeptnState {
	var oldstatus apicommon.KeptnState
	for _, ts := range statuses {
		if ts.EvaluationDefinitionName == evaluationName {
			oldstatus = ts.Status
		}
	}
	return oldstatus
}

func (r EvaluationHandler) setupEvaluations(evaluationCreateAttributes EvaluationCreateAttributes, piWrapper *interfaces.PhaseItemWrapper) ([]string, []klcv1alpha1.EvaluationStatus) {
	var evaluations []string
	var statuses []klcv1alpha1.EvaluationStatus

	switch evaluationCreateAttributes.CheckType {
	case apicommon.PreDeploymentEvaluationCheckType:
		evaluations = piWrapper.GetPreDeploymentEvaluations()
		statuses = piWrapper.GetPreDeploymentEvaluationTaskStatus()
	case apicommon.PostDeploymentEvaluationCheckType:
		evaluations = piWrapper.GetPostDeploymentEvaluations()
		statuses = piWrapper.GetPostDeploymentEvaluationTaskStatus()
	}
	return evaluations, statuses
}

func (r EvaluationHandler) CreateKeptnEvaluation(ctx context.Context, namespace string, reconcileObject client.Object, evaluationCreateAttributes EvaluationCreateAttributes) (string, error) {
	piWrapper, err := interfaces.NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return "", err
	}

	phase := apicommon.PhaseCreateEvaluation

	newEvaluation := piWrapper.GenerateEvaluation(evaluationCreateAttributes.EvaluationDefinition, evaluationCreateAttributes.CheckType)
	err = controllerutil.SetControllerReference(reconcileObject, &newEvaluation, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}
	err = r.Client.Create(ctx, &newEvaluation)
	if err != nil {
		r.Log.Error(err, "could not create KeptnEvaluation")
		RecordEvent(r.Recorder, phase, "Warning", reconcileObject, "CreateFailed", "could not create KeptnEvaluation", piWrapper.GetVersion())
		return "", err
	}
	RecordEvent(r.Recorder, phase, "Normal", reconcileObject, "Created", "created", piWrapper.GetVersion())

	return newEvaluation.Name, nil
}

func (r EvaluationHandler) emitEvaluationFailureEvents(evaluation *klcv1alpha1.KeptnEvaluation, spanTrace trace.Span, piWrapper *interfaces.PhaseItemWrapper) {
	k8sEventMessage := "evaluation failed"
	for k, v := range evaluation.Status.EvaluationStatus {
		if v.Status == apicommon.StateFailed {
			msg := fmt.Sprintf("evaluation of '%s' failed with value: '%s' and reason: '%s'", k, v.Value, v.Message)
			spanTrace.AddEvent(msg, trace.WithTimestamp(time.Now().UTC()))
			k8sEventMessage = fmt.Sprintf("%s\n%s", k8sEventMessage, msg)
		}
	}
	RecordEvent(r.Recorder, apicommon.PhaseReconcileEvaluation, "Warning", evaluation, "Failed", k8sEventMessage, piWrapper.GetVersion())
}
