package keptntask

import (
	"context"
	"fmt"
	"github.com/imdario/mergo"
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *KeptnTaskReconciler) createJob(ctx context.Context, req ctrl.Request, task *klcv1alpha1.KeptnTask) error {
	job := ""
	definition, err := r.getTaskDefinition(ctx, task.Spec.TaskDefinition, req.Namespace)
	if err != nil {
		r.Recorder.Event(task, "Warning", "TaskDefinitionNotFound", fmt.Sprintf("Could not find KeptnTaskDefinition / Namespace: %s, Name: %s ", task.Namespace, task.Spec.TaskDefinition))
		return err
	}

	if !reflect.DeepEqual(definition.Spec.Function, klcv1alpha1.FunctionSpec{}) {
		job, err = r.createFunctionJob(ctx, req, task, definition)
		if err != nil {
			return err
		}
	}

	task.Status.JobName = job
	task.Status.Status = klcv1alpha1.TaskPending
	err = r.Client.Status().Update(ctx, task)
	if err != nil {
		r.Log.Error(err, "could not update configmap status reference for: "+definition.Name)
	}
	r.Log.Info("updated configmap status reference for: " + definition.Name)
	return nil
}

func (r *KeptnTaskReconciler) createFunctionJob(ctx context.Context, req ctrl.Request, task *klcv1alpha1.KeptnTask, definition *klcv1alpha1.KeptnTaskDefinition) (string, error) {
	params, hasParent, err := parseFunctionTaskDefinition(definition)
	var parentJobParams FunctionExecutionParams
	if err != nil {
		return "", err
	}
	if hasParent {
		parentDefinition, err := r.getTaskDefinition(ctx, definition.Spec.Function.FunctionReference.Name, req.Namespace)
		if err != nil {
			r.Recorder.Event(task, "Warning", "TaskDefinitionNotFound", fmt.Sprintf("Could not find KeptnTaskDefinition / Namespace: %s, Name: %s ", task.Namespace, task.Spec.TaskDefinition))
			return "", err
		}
		parentJobParams, _, err = parseFunctionTaskDefinition(parentDefinition)
		if err != nil {
			return "", err
		}
		err = mergo.Merge(&params, parentJobParams)
		if err != nil {
			r.Recorder.Event(task, "Warning", "TaskDefinitionMergeFailure", fmt.Sprintf("Could not merge KeptnTaskDefinition / Namespace: %s, Name: %s ", task.Namespace, task.Spec.TaskDefinition))
			return "", err
		}
	}

	if len(task.Spec.Parameters.Inline) > 0 {
		err = mergo.Merge(&params.Parameters, task.Spec.Parameters.Inline)
		if err != nil {
			r.Recorder.Event(task, "Warning", "TaskDefinitionMergeFailure", fmt.Sprintf("Could not merge KeptnTaskDefinition / Namespace: %s, Name: %s ", task.Namespace, task.Spec.TaskDefinition))
			return "", err
		}
	}

	if task.Spec.SecureParameters.Secret != "" {
		params.SecureParameters = task.Spec.SecureParameters.Secret
	}

	job, _ := r.generateFunctionJob(task, params)
	if err != nil {
		return "", err
	}
	err = r.Client.Create(ctx, job)
	if err != nil {
		r.Log.Error(err, "could not create job")
		r.Recorder.Event(task, "Warning", "JobNotCreated", fmt.Sprintf("Could not create Job / Namespace: %s, Name: %s ", task.Namespace, task.Name))
		return job.Name, err
	}
	r.Recorder.Event(task, "Normal", "JobCreated", fmt.Sprintf("Created Job / Namespace: %s, Name: %s ", task.Namespace, task.Name))
	return job.Name, nil
}

func (r *KeptnTaskReconciler) updateJob(ctx context.Context, req ctrl.Request, task *klcv1alpha1.KeptnTask) error {
	job, err := r.getJob(ctx, task.Status.JobName, req.Namespace)
	if err != nil {
		task.Status.JobName = ""
		r.Recorder.Event(task, "Warning", "JobReferenceRemoved", fmt.Sprintf("Removed Job Reference as Job could not be found / Namespace: %s, Name: %s ", task.Namespace, task.Name))
		err = r.Client.Status().Update(ctx, task)
		if err != nil {
			r.Log.Error(err, "could not update job reference reference for: "+task.Name)
		}
		return err
	}
	if job.Status.Succeeded > 0 {
		task.Status.Status = klcv1alpha1.TaskSucceeded
		err = r.Client.Status().Update(ctx, task)
		if err != nil {
			r.Log.Error(err, "could not update job reference reference for: "+task.Name)
		}
	}
	return nil
}
func (r *KeptnTaskReconciler) getJob(ctx context.Context, jobName string, namespace string) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: jobName, Namespace: namespace}, job)
	if err != nil {
		return job, err
	}
	return job, nil
}
