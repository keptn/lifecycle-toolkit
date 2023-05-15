package providers

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type ContainerRuntimeProvider struct {
	Log logr.Logger
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func (cp *ContainerRuntimeProvider) CreateJob(ctx context.Context, req ctrl.Request, task *klcv1alpha3.KeptnTask, definition *klcv1alpha3.KeptnTaskDefinition) (string, error) {
	job, err := cp.generateFunctionJob(ctx, task)
	if err != nil {
		return "", err
	}
	err = cp.Client.Create(ctx, job)
	if err != nil {
		cp.Log.Error(err, "could not create job")
		controllercommon.RecordEvent(cp.Recorder, apicommon.PhaseCreateTask, "Warning", task, "JobNotCreated", fmt.Sprintf("could not create Job: %s ", task.Name), "")
		return job.Name, err
	}

	controllercommon.RecordEvent(cp.Recorder, apicommon.PhaseReconcileTask, "Normal", task, "JobCreated", fmt.Sprintf("created Job: %s ", task.Name), "")
	return job.Name, nil
}

func (cp *ContainerRuntimeProvider) generateFunctionJob(ctx context.Context, task *klcv1alpha3.KeptnTask) (*batchv1.Job, error) {
	randomId := rand.Intn(99999-10000) + 10000
	jobId := fmt.Sprintf("klc-%s-%d", apicommon.TruncateString(task.Name, apicommon.MaxTaskNameLength), randomId)
	taskDefinition, err := GetTaskDefinition(ctx, cp.Client, task.Spec.TaskDefinition, task.Namespace)
	if err != nil {
		controllercommon.RecordEvent(cp.Recorder, apicommon.PhaseCreateTask, "Warning", task, "TaskDefinitionNotFound", fmt.Sprintf("could not find KeptnTaskDefinition: %s ", task.Spec.TaskDefinition), "")
		return nil, err
	}
	if taskDefinition.Spec.Container.JobNamespace == "" {
		taskDefinition.Spec.Container.JobNamespace = task.Namespace
	}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobId,
			Namespace: taskDefinition.Spec.Container.JobNamespace,
			Labels:    task.CreateKeptnLabels(),
		},
		Spec: batchv1.JobSpec{
			ActiveDeadlineSeconds: task.GetActiveDeadlineSeconds(),
			BackoffLimit:          task.Spec.Retries,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:   "OnFailure",
					SecurityContext: taskDefinition.Spec.Container.DefaultPodSecurityContext,
				},
			},
		},
	}
	err = controllerutil.SetControllerReference(task, job, cp.Scheme)
	if err != nil {
		cp.Log.Error(err, "could not set controller reference:")
	}

	container := corev1.Container{
		Name:            "keptn-container-runner",
		Image:           taskDefinition.Spec.Container.Image,
		ImagePullPolicy: taskDefinition.Spec.Container.ImagePullPolicy,
		SecurityContext: taskDefinition.Spec.Container.DefaultSecurityContext,
		Resources:       taskDefinition.Spec.Container.DefaultResourceRequirements,
		Command:         taskDefinition.Spec.Container.Command,
		Args:            taskDefinition.Spec.Container.Args,
	}

	job.Spec.Template.Spec.Containers = []corev1.Container{
		container,
	}
	return job, nil
}
