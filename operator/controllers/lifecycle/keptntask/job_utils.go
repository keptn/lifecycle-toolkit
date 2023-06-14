package keptntask

import (
	"context"
	"fmt"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *KeptnTaskReconciler) createJob(ctx context.Context, req ctrl.Request, task *klcv1alpha3.KeptnTask) error {
	jobName := ""
	definition, err := controllercommon.GetTaskDefinition(r.Client, r.Log, ctx, task.Spec.TaskDefinition, req.Namespace)
	if err != nil {
		controllercommon.RecordEvent(r.Recorder, apicommon.PhaseCreateTask, "Warning", task, "TaskDefinitionNotFound", fmt.Sprintf("could not find KeptnTaskDefinition: %s ", task.Spec.TaskDefinition), "")
		return err
	}

	if controllercommon.SpecExists(definition) {
		jobName, err = r.createFunctionJob(ctx, req, task, definition)
		if err != nil {
			return err
		}
	}

	task.Status.JobName = jobName
	task.Status.Status = apicommon.StatePending

	return nil
}

func (r *KeptnTaskReconciler) createFunctionJob(ctx context.Context, req ctrl.Request, task *klcv1alpha3.KeptnTask, definition *klcv1alpha3.KeptnTaskDefinition) (string, error) {

	job, err := r.generateJob(ctx, task, definition, req)
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

func (r *KeptnTaskReconciler) updateTaskStatus(job *batchv1.Job, task *klcv1alpha3.KeptnTask) {
	if len(job.Status.Conditions) > 0 {
		if job.Status.Conditions[0].Type == batchv1.JobComplete {
			task.Status.Status = apicommon.StateSucceeded
		} else if job.Status.Conditions[0].Type == batchv1.JobFailed {
			task.Status.Status = apicommon.StateFailed
			task.Status.Message = job.Status.Conditions[0].Message
			task.Status.Reason = job.Status.Conditions[0].Reason
		}
	}
}

func (r *KeptnTaskReconciler) getJob(ctx context.Context, jobName string, namespace string) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: jobName, Namespace: namespace}, job)
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
	taskContext.TaskType = string(task.Spec.Type)
	taskContext.AppName = task.Spec.AppName

	return taskContext
}

func (r *KeptnTaskReconciler) generateJob(ctx context.Context, task *klcv1alpha3.KeptnTask, definition *klcv1alpha3.KeptnTaskDefinition, request ctrl.Request) (*batchv1.Job, error) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        apicommon.GenerateJobName(task.Name),
			Namespace:   task.Namespace,
			Labels:      task.Labels,
			Annotations: task.CreateKeptnAnnotations(),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      task.Labels,
					Annotations: task.Annotations,
				},
				Spec: corev1.PodSpec{
					RestartPolicy: "OnFailure",
				},
			},
			BackoffLimit:          task.Spec.Retries,
			ActiveDeadlineSeconds: task.GetActiveDeadlineSeconds(),
		},
	}
	err := controllerutil.SetControllerReference(task, job, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}

	builderOpt := BuilderOptions{
		Client:        r.Client,
		req:           request,
		Log:           r.Log,
		task:          task,
		containerSpec: definition.Spec.Container,
		funcSpec:      controllercommon.GetRuntimeSpec(definition),
		recorder:      r.Recorder,
		Image:         controllercommon.GetRuntimeImage(definition),
		MountPath:     controllercommon.GetRuntimeMountPath(definition),
		ConfigMap:     definition.Status.Function.ConfigMap,
	}

	builder := NewJobRunnerBuilder(builderOpt)
	if builder == nil {
		return nil, controllererrors.ErrNoTaskDefinitionSpec
	}

	container, err := builder.CreateContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create container for Job: %w", err)
	}

	volume, err := builder.CreateVolume(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create volume for Job: %w", err)
	}

	if volume != nil {
		job.Spec.Template.Spec.Volumes = []corev1.Volume{*volume}
	}

	job.Spec.Template.Spec.Containers = []corev1.Container{*container}

	return job, nil
}
