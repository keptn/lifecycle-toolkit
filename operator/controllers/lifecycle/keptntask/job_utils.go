package keptntask

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"

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

	if emptySpec(definition) {
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

func emptySpec(definition *klcv1alpha3.KeptnTaskDefinition) bool {
	//TODO when adding new builders add more logic here
	return !reflect.DeepEqual(definition.Spec.Function, klcv1alpha3.FunctionSpec{})
}

func (r *KeptnTaskReconciler) updateJob(ctx context.Context, req ctrl.Request, task *klcv1alpha3.KeptnTask) error {
	job, err := r.getJob(ctx, task.Status.JobName, req.Namespace)
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
	taskContext.AppName = task.Spec.AppName

	return taskContext
}

func (r *KeptnTaskReconciler) generateJob(ctx context.Context, task *klcv1alpha3.KeptnTask, definition *klcv1alpha3.KeptnTaskDefinition, request ctrl.Request) (*batchv1.Job, error) {
	randomId := rand.Intn(99999-10000) + 10000
	jobId := fmt.Sprintf("klc-%s-%d", apicommon.TruncateString(task.Name, apicommon.MaxTaskNameLength), randomId)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        jobId,
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
		Client:   r.Client,
		req:      request,
		Log:      r.Log,
		task:     task,
		taskDef:  definition,
		recorder: r.Recorder,
	}
	builder := getContainerBuilder(builderOpt)
	if builder == nil {
		return nil, controllererrors.ErrNoTaskDefinitionSpec
	}
	container, volumes, err := builder.CreateContainerWithVolumes(ctx)

	if err != nil {
		return nil, controllererrors.ErrCannotMarshalParams
	}

	job.Spec.Template.Spec.Containers = []corev1.Container{*container}
	job.Spec.Template.Spec.Volumes = volumes
	return job, nil
}
