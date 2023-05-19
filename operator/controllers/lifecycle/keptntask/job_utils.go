package keptntask

import (
	"context"
	"fmt"
	"reflect"

	"github.com/imdario/mergo"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *KeptnTaskReconciler) createJob(ctx context.Context, req ctrl.Request, task *klcv1alpha3.KeptnTask) error {
	jobName := ""
	definition, err := controllercommon.GetTaskDefinition(r.Client, r.Log, ctx, task.Spec.TaskDefinition, req.Namespace)
	if err != nil {
		controllercommon.RecordEvent(r.Recorder, apicommon.PhaseCreateTask, "Warning", task, "TaskDefinitionNotFound", fmt.Sprintf("could not find KeptnTaskDefinition: %s ", task.Spec.TaskDefinition), "")
		return err
	}

	if !reflect.DeepEqual(definition.Spec.Function, klcv1alpha3.FunctionSpec{}) {
		jobName, err = r.createFunctionJob(ctx, req, task, definition)
		if err != nil {
			return err
		}
	}

	task.Status.JobName = jobName
	task.Status.Status = apicommon.StatePending
	err = r.Client.Status().Update(ctx, task)
	if err != nil {
		r.Log.Error(err, "could not update KeptnTask status reference for: "+task.Name)
	}
	r.Log.Info("updated configmap status reference for: " + definition.Name)
	return nil
}

func (r *KeptnTaskReconciler) createFunctionJob(ctx context.Context, req ctrl.Request, task *klcv1alpha3.KeptnTask, definition *klcv1alpha3.KeptnTaskDefinition) (string, error) {
	params, hasParent, err := r.parseFunctionTaskDefinition(definition)
	if err != nil {
		return "", err
	}
	if hasParent {
		if err := r.handleParent(ctx, req, task, definition, params); err != nil {
			return "", err
		}
	}

	params.Context = setupTaskContext(task)

	if len(task.Spec.Parameters.Inline) > 0 {
		err = mergo.Merge(&params.Parameters, task.Spec.Parameters.Inline)
		if err != nil {
			controllercommon.RecordEvent(r.Recorder, apicommon.PhaseCreateTask, "Warning", task, "TaskDefinitionMergeFailure", fmt.Sprintf("could not merge KeptnTaskDefinition: %s ", task.Spec.TaskDefinition), "")
			return "", err
		}
	}

	if task.Spec.SecureParameters.Secret != "" {
		params.SecureParameters = task.Spec.SecureParameters.Secret
	}

	job, err := r.generateFunctionJob(task, params)
	if err != nil {
		return "", err
	}
	err = r.Client.Create(ctx, job)
	if err != nil {
		r.Log.Error(err, "could not create job")
		controllercommon.RecordEvent(r.Recorder, apicommon.PhaseCreateTask, "Warning", task, "JobNotCreated", fmt.Sprintf("could not create Job: %s ", task.Name), "")
		return job.Name, err
	}

	controllercommon.RecordEvent(r.Recorder, apicommon.PhaseReconcileTask, "Normal", task, "JobCreated", fmt.Sprintf("created Job: %s ", task.Name), "")
	return job.Name, nil
}

func (r *KeptnTaskReconciler) updateJob(ctx context.Context, req ctrl.Request, task *klcv1alpha3.KeptnTask) error {
	job, err := r.getJob(ctx, task, req.Namespace)
	if err != nil {
		task.Status.JobName = ""
		controllercommon.RecordEvent(r.Recorder, apicommon.PhaseReconcileTask, "Warning", task, "JobReferenceRemoved", "removed Job Reference as Job could not be found", "")
		err = r.Client.Status().Update(ctx, task)
		if err != nil {
			r.Log.Error(err, "could not remove job reference for: "+task.Name)
		}
		return err
	}
	if len(job.Status.Conditions) > 0 {
		if job.Status.Conditions[0].Type == batchv1.JobComplete {
			task.Status.Status = apicommon.StateSucceeded
		} else if job.Status.Conditions[0].Type == batchv1.JobFailed {
			task.Status.Status = apicommon.StateFailed
			task.Status.Message = job.Status.Conditions[0].Message
			task.Status.Reason = job.Status.Conditions[0].Reason
		}
	}
	return nil
}
func (r *KeptnTaskReconciler) getJob(ctx context.Context, task *v1alpha3.KeptnTask, namespace string) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: task.Status.JobName, Namespace: namespace}, job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func setupTaskContext(task *klcv1alpha3.KeptnTask) klcv1alpha3.TaskContext {
	taskContext := klcv1alpha3.TaskContext{}

	if task.Spec.Workload != "" {
		taskContext.WorkloadName = task.Spec.Workload
		taskContext.WorkloadVersion = task.Spec.WorkloadVersion
		taskContext.ObjectType = "Workload"

	} else {
		taskContext.ObjectType = "Application"
		taskContext.AppVersion = task.Spec.AppVersion
	}
	taskContext.AppName = task.Spec.AppName

	return taskContext
}

func (r *KeptnTaskReconciler) handleParent(ctx context.Context, req ctrl.Request, task *klcv1alpha3.KeptnTask, definition *klcv1alpha3.KeptnTaskDefinition, params FunctionExecutionParams) error {
	var parentJobParams FunctionExecutionParams
	parentDefinition, err := controllercommon.GetTaskDefinition(r.Client, r.Log, ctx, definition.Spec.Function.FunctionReference.Name, req.Namespace)
	if err != nil {
		controllercommon.RecordEvent(r.Recorder, apicommon.PhaseCreateTask, "Warning", task, "TaskDefinitionNotFound", fmt.Sprintf("could not find KeptnTaskDefinition: %s ", task.Spec.TaskDefinition), "")
		return err
	}
	parentJobParams, _, err = r.parseFunctionTaskDefinition(parentDefinition)
	if err != nil {
		return err
	}
	err = mergo.Merge(&params, parentJobParams)
	if err != nil {
		controllercommon.RecordEvent(r.Recorder, apicommon.PhaseCreateTask, "Warning", task, "TaskDefinitionMergeFailure", fmt.Sprintf("could not merge KeptnTaskDefinition: %s ", task.Spec.TaskDefinition), "")
		return err
	}
	return nil
}
