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

	var summary apicommon.StatusSummary
	summary.Total = len(evaluations)
	// Check current state of the PrePostEvaluationTasks
	var newStatus []klcv1alpha1.EvaluationStatus
	for _, evaluationName := range evaluations {
		var oldstatus apicommon.KeptnState
		for _, ts := range statuses {
			if ts.EvaluationDefinitionName == evaluationName {
				oldstatus = ts.Status
			}
		}

		evaluationStatus := GetEvaluationStatus(evaluationName, statuses)
		evaluation := &klcv1alpha1.KeptnEvaluation{}
		evaluationExists := false

		if oldstatus != evaluationStatus.Status {
			RecordEvent(r.Recorder, apicommon.PhaseReconcileEvaluation, "Normal", reconcileObject, "EvaluationStatusChanged", fmt.Sprintf("evaluation status changed from %s to %s", oldstatus, evaluationStatus.Status), piWrapper.GetVersion())
		}

		// Check if evaluation has already succeeded or failed
		if evaluationStatus.Status.IsCompleted() {
			newStatus = append(newStatus, evaluationStatus)
			continue
		}

		// Check if Evaluation is already created
		if evaluationStatus.EvaluationName != "" {
			err := r.Client.Get(ctx, types.NamespacedName{Name: evaluationStatus.EvaluationName, Namespace: piWrapper.GetNamespace()}, evaluation)
			if err != nil && errors.IsNotFound(err) {
				evaluationStatus.EvaluationName = ""
			} else if err != nil {
				return nil, summary, err
			}
			evaluationExists = true
		}

		// Create new Evaluation if it does not exist
		if !evaluationExists {
			evaluationCreateAttributes.EvaluationDefinition = evaluationName
			evaluationName, err := r.CreateKeptnEvaluation(ctx, piWrapper.GetNamespace(), reconcileObject, evaluationCreateAttributes)
			if err != nil {
				return nil, summary, err
			}
			evaluationStatus.EvaluationName = evaluationName
			evaluationStatus.SetStartTime()
			_, _, err = r.SpanHandler.GetSpan(phaseCtx, r.Tracer, evaluation, "")
			if err != nil {
				r.Log.Error(err, "could not get span")
			}
		} else {
			_, spanEvaluationTrace, err := r.SpanHandler.GetSpan(phaseCtx, r.Tracer, evaluation, "")
			if err != nil {
				r.Log.Error(err, "could not get span")
			}
			// Update state of Evaluation if it is already created
			evaluationStatus.Status = evaluation.Status.OverallStatus
			if evaluationStatus.Status.IsCompleted() {
				if evaluationStatus.Status.IsSucceeded() {
					spanEvaluationTrace.AddEvent(evaluation.Name + " has finished")
					spanEvaluationTrace.SetStatus(codes.Ok, "Finished")
					RecordEvent(r.Recorder, apicommon.PhaseReconcileEvaluation, "Normal", evaluation, "Succeeded", "evaluation succeeded", piWrapper.GetVersion())
				} else {
					spanEvaluationTrace.AddEvent(evaluation.Name + " has failed")
					r.emitEvaluationFailureEvents(evaluation, spanEvaluationTrace, piWrapper)
					spanEvaluationTrace.SetStatus(codes.Error, "Failed")
				}
				spanEvaluationTrace.End()
				if err := r.SpanHandler.UnbindSpan(evaluation, ""); err != nil {
					r.Log.Error(err, controllererrors.ErrCouldNotUnbindSpan, evaluation.Name)
				}
				evaluationStatus.SetEndTime()
			}
		}
		// Update state of the Check
		newStatus = append(newStatus, evaluationStatus)
	}

	for _, ns := range newStatus {
		summary = apicommon.UpdateStatusSummary(ns.Status, summary)
	}
	if apicommon.GetOverallState(summary) != apicommon.StateSucceeded {
		RecordEvent(r.Recorder, apicommon.PhaseReconcileEvaluation, "Warning", reconcileObject, "NotFinished", "has not finished", piWrapper.GetVersion())
	}
	return newStatus, summary, nil
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
