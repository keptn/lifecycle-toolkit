package evaluation

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/interfaces"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/evaluationhandler_mock.go . IEvaluationHandler:MockEvaluationHandler
type IEvaluationHandler interface {
	ReconcileEvaluations(ctx context.Context, phaseCtx context.Context, reconcileObject client.Object, evaluationCreateAttributes CreateEvaluationAttributes) ([]apilifecycle.ItemStatus, apicommon.StatusSummary, error)
}

type Handler struct {
	client.Client
	EventSender eventsender.IEvent
	Log         logr.Logger
	Tracer      telemetry.ITracer
	Scheme      *runtime.Scheme
	SpanHandler telemetry.ISpanHandler
}

type CreateEvaluationAttributes struct {
	SpanName   string
	Definition apilifecycle.KeptnEvaluationDefinition
	CheckType  apicommon.CheckType
}

// NewHandler creates a new instance of the Handler.
func NewHandler(client client.Client, eventSender eventsender.IEvent, log logr.Logger, tracer telemetry.ITracer, scheme *runtime.Scheme, spanHandler telemetry.ISpanHandler) Handler {
	return Handler{
		Client:      client,
		EventSender: eventSender,
		Log:         log,
		Tracer:      tracer,
		Scheme:      scheme,
		SpanHandler: spanHandler,
	}
}

//nolint:gocognit,gocyclo
func (r Handler) ReconcileEvaluations(ctx context.Context, phaseCtx context.Context, reconcileObject client.Object, evaluationCreateAttributes CreateEvaluationAttributes) ([]apilifecycle.ItemStatus, apicommon.StatusSummary, error) {
	piWrapper, err := interfaces.NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return nil, apicommon.StatusSummary{}, err
	}

	evaluations, statuses := r.setupEvaluations(evaluationCreateAttributes, piWrapper)

	var summary apicommon.StatusSummary
	summary.Total = len(evaluations)
	// Check current state of the PrePostEvaluationTasks
	var newStatus []apilifecycle.ItemStatus
	for _, evaluationName := range evaluations {
		oldstatus := common.GetOldStatus(evaluationName, statuses)

		evaluationStatus := common.GetItemStatus(evaluationName, statuses)
		evaluation := &apilifecycle.KeptnEvaluation{}
		evaluationExists := false

		if oldstatus != evaluationStatus.Status {
			r.EventSender.Emit(apicommon.PhaseReconcileEvaluation, "Normal", reconcileObject, apicommon.PhaseStateStatusChanged, fmt.Sprintf("evaluation status changed from %s to %s", oldstatus, evaluationStatus.Status), piWrapper.GetVersion())
		}

		// Check if evaluation has already succeeded or failed
		if evaluationStatus.Status.IsCompleted() {
			newStatus = append(newStatus, evaluationStatus)
			continue
		}

		// Check if Evaluation is already created
		if evaluationStatus.Name != "" {
			err := r.Client.Get(ctx, types.NamespacedName{Name: evaluationStatus.Name, Namespace: piWrapper.GetNamespace()}, evaluation)
			if err != nil && errors.IsNotFound(err) {
				evaluationStatus.Name = ""
			} else if err != nil {
				return nil, summary, err
			}
			evaluationExists = true
		}

		// Create new Evaluation if it does not exist
		if !evaluationExists {
			err := r.handleEvaluationNotExists(
				ctx,
				phaseCtx,
				evaluationCreateAttributes,
				evaluationName,
				piWrapper,
				reconcileObject,
				evaluation,
				&evaluationStatus,
			)
			if err != nil {
				if errors.IsNotFound(err) {
					r.Log.Info("EvaluationDefinition for Evaluation not found",
						"evaluation", evaluationStatus.Name,
						"evaluationDefinition", evaluationName,
						"namespace", piWrapper.GetNamespace(),
					)
				} else {
					// log the error, but continue to proceed with other evaluations that may be created
					r.Log.Error(err, "Could not create evaluation",
						"evaluation", evaluationStatus.Name,
						"evaluationDefinition", evaluationName,
						"namespace", piWrapper.GetNamespace(),
					)
				}
				continue
			}
		} else {
			r.handleEvaluationExists(
				phaseCtx,
				piWrapper,
				evaluation,
				&evaluationStatus,
			)
		}
		// Update state of the Check
		newStatus = append(newStatus, evaluationStatus)
	}

	for _, ns := range newStatus {
		summary = apicommon.UpdateStatusSummary(ns.Status, summary)
	}

	return newStatus, summary, nil
}

//nolint:dupl
func (r Handler) CreateKeptnEvaluation(ctx context.Context, reconcileObject client.Object, evaluationCreateAttributes CreateEvaluationAttributes) (string, error) {
	piWrapper, err := interfaces.NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return "", err
	}

	phase := apicommon.PhaseCreateEvaluation

	newEvaluation := piWrapper.GenerateEvaluation(evaluationCreateAttributes.Definition, evaluationCreateAttributes.CheckType)
	err = controllerutil.SetControllerReference(reconcileObject, &newEvaluation, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}
	err = r.Client.Create(ctx, &newEvaluation)
	if err != nil {
		r.Log.Error(err, "could not create KeptnEvaluation")
		r.EventSender.Emit(phase, "Warning", reconcileObject, apicommon.PhaseStateFailed, "could not create KeptnEvaluation", piWrapper.GetVersion())
		return "", err
	}

	return newEvaluation.Name, nil
}

func (r Handler) emitEvaluationFailureEvents(evaluation *apilifecycle.KeptnEvaluation, spanTrace trace.Span, piWrapper *interfaces.PhaseItemWrapper) {
	k8sEventMessage := "evaluation failed"
	for k, v := range evaluation.Status.EvaluationStatus {
		if v.Status == apicommon.StateFailed {
			msg := fmt.Sprintf("evaluation of '%s' failed with value: '%s' and reason: '%s'", k, v.Value, v.Message)
			spanTrace.AddEvent(msg, trace.WithTimestamp(time.Now().UTC()))
			k8sEventMessage = fmt.Sprintf("%s\n%s", k8sEventMessage, msg)
		}
	}
	r.EventSender.Emit(apicommon.PhaseReconcileEvaluation, "Warning", evaluation, apicommon.PhaseStateFailed, k8sEventMessage, piWrapper.GetVersion())
}

func (r Handler) setupEvaluations(evaluationCreateAttributes CreateEvaluationAttributes, piWrapper *interfaces.PhaseItemWrapper) ([]string, []apilifecycle.ItemStatus) {
	var evaluations []string
	var statuses []apilifecycle.ItemStatus

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

func (r Handler) handleEvaluationNotExists(ctx context.Context, phaseCtx context.Context, evaluationCreateAttributes CreateEvaluationAttributes, evaluationName string, piWrapper *interfaces.PhaseItemWrapper, reconcileObject client.Object, evaluation *apilifecycle.KeptnEvaluation, evaluationStatus *apilifecycle.ItemStatus) error {
	definition, err := common.GetEvaluationDefinition(r.Client, r.Log, ctx, evaluationName, piWrapper.GetNamespace())
	if err != nil {
		return controllererrors.ErrCannotGetKeptnEvaluationDefinition
	}
	evaluationCreateAttributes.Definition = *definition
	evaluationName, err = r.CreateKeptnEvaluation(ctx, reconcileObject, evaluationCreateAttributes)
	if err != nil {
		return err
	}
	evaluationStatus.Name = evaluationName
	evaluationStatus.SetStartTime()
	_, _, err = r.SpanHandler.GetSpan(phaseCtx, r.Tracer, evaluation, "")
	if err != nil {
		r.Log.Error(err, "could not get span")
	}

	return nil
}

func (r Handler) handleEvaluationExists(phaseCtx context.Context, piWrapper *interfaces.PhaseItemWrapper, evaluation *apilifecycle.KeptnEvaluation, evaluationStatus *apilifecycle.ItemStatus) {
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
