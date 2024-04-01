package keptntask

import (
	"context"
	"testing"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/testcommon"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const KeptnNamespace = "mynamespace"

func TestKeptnTaskReconciler_createJob(t *testing.T) {
	namespace := "default"
	cmName := "my-cmd"
	taskDefinitionName := "my-task-definition"

	cm := makeConfigMap(cmName, namespace)

	taskDefinition := makeTaskDefinitionWithConfigmapRef(taskDefinitionName, namespace, cmName)
	fakeClient := testcommon.NewTestClient(cm, taskDefinition)

	taskDefinition.Status.Function.ConfigMap = cmName
	err := fakeClient.Status().Update(context.TODO(), taskDefinition)
	require.Nil(t, err)

	r := &KeptnTaskReconciler{
		Client:      fakeClient,
		EventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
		Log:         ctrl.Log.WithName("task-controller"),
		Scheme:      fakeClient.Scheme(),
	}

	task := makeTask("my-task", namespace, taskDefinitionName)

	err = fakeClient.Create(context.TODO(), task)
	require.Nil(t, err)

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: namespace,
		},
	}

	// retrieve the task again to verify its status
	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      task.Name,
	}, task)

	require.Nil(t, err)

	err = r.createJob(context.TODO(), req, task)
	require.Nil(t, err)

	require.NotEmpty(t, task.Status.JobName)

	resultingJob := &batchv1.Job{}
	err = fakeClient.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: task.Status.JobName}, resultingJob)
	require.Nil(t, err)

	require.Equal(t, namespace, resultingJob.Namespace)
	require.NotEmpty(t, resultingJob.OwnerReferences)
	require.Len(t, resultingJob.Spec.Template.Spec.Containers, 1)
	require.Len(t, resultingJob.Spec.Template.Spec.Containers[0].Env, 5)
	require.Equal(t, map[string]string{
		"label1": "label2",
	}, resultingJob.Labels)
	require.Equal(t, map[string]string{
		"annotation1":        "annotation2",
		"keptn.sh/app":       "my-app",
		"keptn.sh/task-name": "my-task",
		"keptn.sh/version":   "",
		"keptn.sh/workload":  "my-workload",
	}, resultingJob.Annotations)
}

func TestKeptnTaskReconciler_createJob_withTaskDefInDefaultNamespace(t *testing.T) {
	namespace := "default"
	cmName := "my-cmd"
	taskDefinitionName := "my-task-definition"

	cm := makeConfigMap(cmName, namespace)
	taskDefinition := makeTaskDefinitionWithConfigmapRef(taskDefinitionName, KeptnNamespace, cmName)

	fakeClient := testcommon.NewTestClient(cm, taskDefinition)

	taskDefinition.Status.Function.ConfigMap = cmName
	err := fakeClient.Status().Update(context.TODO(), taskDefinition)
	require.Nil(t, err)

	config.Instance().SetDefaultNamespace(KeptnNamespace)
	r := &KeptnTaskReconciler{
		Client:      fakeClient,
		EventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
		Log:         ctrl.Log.WithName("task-controller"),
		Scheme:      fakeClient.Scheme(),
	}

	task := makeTask("my-task", namespace, taskDefinitionName)

	err = fakeClient.Create(context.TODO(), task)
	require.Nil(t, err)

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: namespace,
		},
	}

	// retrieve the task again to verify its status
	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      task.Name,
	}, task)

	require.Nil(t, err)

	err = r.createJob(context.TODO(), req, task)
	require.Nil(t, err)

	require.NotEmpty(t, task.Status.JobName)

	resultingJob := &batchv1.Job{}
	err = fakeClient.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: task.Status.JobName}, resultingJob)
	require.Nil(t, err)

	require.Equal(t, namespace, resultingJob.Namespace)
	require.NotEmpty(t, resultingJob.OwnerReferences)
	require.Len(t, resultingJob.Spec.Template.Spec.Containers, 1)
	require.Len(t, resultingJob.Spec.Template.Spec.Containers[0].Env, 5)
	require.Equal(t, map[string]string{
		"label1": "label2",
	}, resultingJob.Labels)
	require.Equal(t, map[string]string{
		"annotation1":        "annotation2",
		"keptn.sh/app":       "my-app",
		"keptn.sh/task-name": "my-task",
		"keptn.sh/version":   "",
		"keptn.sh/workload":  "my-workload",
	}, resultingJob.Annotations)
}

func TestKeptnTaskReconciler_updateTaskStatus(t *testing.T) {
	namespace := "default"
	taskDefinitionName := "my-task-definition"

	jobStatus := batchv1.JobStatus{
		Conditions: []batchv1.JobCondition{
			{
				Type: batchv1.JobFailed,
			},
		},
	}

	job := makeJob("my.job", namespace, jobStatus)

	fakeClient := fake.NewClientBuilder().WithObjects(job).Build()

	err := apilifecycle.AddToScheme(fakeClient.Scheme())
	require.Nil(t, err)

	r := &KeptnTaskReconciler{
		Client:      fakeClient,
		EventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
		Log:         ctrl.Log.WithName("task-controller"),
		Scheme:      fakeClient.Scheme(),
	}

	task := makeTask("my-task", namespace, taskDefinitionName)

	err = fakeClient.Create(context.TODO(), task)
	require.Nil(t, err)

	task.Status.JobName = job.Name

	r.updateTaskStatus(job, task)

	require.Equal(t, apicommon.StateFailed, task.Status.Status)

	// now, set the job to succeeded
	job.Status.Conditions = []batchv1.JobCondition{
		{
			Type: batchv1.JobComplete,
		},
	}

	r.updateTaskStatus(job, task)

	require.Equal(t, apicommon.StateSucceeded, task.Status.Status)
}

func TestKeptnTaskReconciler_generateJob(t *testing.T) {
	namespace := "default"
	taskName := "my-task"
	svcAccname := "svcAccname"
	taskDefinitionName := "my-task-definition"
	token := true
	var ttlSecondsAfterFinished int32 = 100
	imagePullSecret := []v1.LocalObjectReference{{
		Name: "my-docker-secret",
	}}

	taskDefinition := makeTaskDefinitionWithServiceAccount(taskDefinitionName, namespace, svcAccname, &token, &ttlSecondsAfterFinished, imagePullSecret)
	taskDefinition.Spec.ServiceAccount.Name = svcAccname
	fakeClient := testcommon.NewTestClient(taskDefinition)
	task := makeTask(taskName, namespace, taskDefinitionName)

	r := &KeptnTaskReconciler{
		Client:      fakeClient,
		EventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
		Log:         ctrl.Log.WithName("task-controller"),
		Scheme:      fakeClient.Scheme(),
	}

	err := fakeClient.Create(context.TODO(), task)
	require.Nil(t, err)

	ctx := context.TODO()
	request := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: namespace,
		},
	}

	errTask := fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      task.Name,
	}, task)
	require.Nil(t, errTask)

	errTaskDefinition := fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      taskDefinition.Name,
	}, taskDefinition)
	require.Nil(t, errTaskDefinition)

	resultingJob, err := r.generateJob(ctx, task, taskDefinition, request)
	require.Nil(t, err)
	require.NotNil(t, resultingJob, "generateJob function return a valid Job")

	require.NotNil(t, resultingJob.Spec.Template.Spec.Containers)
	require.Equal(t, resultingJob.Spec.Template.Spec.ImagePullSecrets[0].Name, imagePullSecret[0].Name, "ImagePullSecret is not assigned correctly")
	require.Equal(t, resultingJob.Spec.Template.Spec.ServiceAccountName, svcAccname)
	require.Equal(t, resultingJob.Spec.Template.Spec.AutomountServiceAccountToken, &token)
	require.Equal(t, resultingJob.Spec.TTLSecondsAfterFinished, &ttlSecondsAfterFinished)
	require.Equal(t, map[string]string{
		"label1": "label2",
	}, resultingJob.Labels)
	require.Equal(t, map[string]string{
		"annotation1":        "annotation2",
		"keptn.sh/app":       "my-app",
		"keptn.sh/task-name": "my-task",
		"keptn.sh/version":   "",
		"keptn.sh/workload":  "my-workload",
	}, resultingJob.Annotations)
}

func makeJob(name, namespace string, status batchv1.JobStatus) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec:   batchv1.JobSpec{},
		Status: status,
	}
}

func makeTask(name, namespace string, taskDefinitionName string) *apilifecycle.KeptnTask {
	return &apilifecycle.KeptnTask{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"label1": "label2",
			},
			Annotations: map[string]string{
				"annotation1": "annotation2",
			},
		},
		Spec: apilifecycle.KeptnTaskSpec{
			Context: apilifecycle.TaskContext{
				WorkloadName: "my-workload",
				AppName:      "my-app",
				AppVersion:   "0.1.0",
				ObjectType:   "Workload",
				TaskType:     string(apicommon.PostDeploymentCheckType),
			},
			TaskDefinition: taskDefinitionName,
			Type:           apicommon.PostDeploymentCheckType,
		},
	}
}

func makeTaskDefinitionWithConfigmapRef(name, namespace, configMapName string) *apilifecycle.KeptnTaskDefinition {
	return &apilifecycle.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"label1": "label2",
			},
			Annotations: map[string]string{
				"annotation1": "annotation2",
			},
		},
		Spec: apilifecycle.KeptnTaskDefinitionSpec{
			Deno: &apilifecycle.RuntimeSpec{
				ConfigMapReference: apilifecycle.ConfigMapReference{
					Name: configMapName,
				},
				Parameters:       apilifecycle.TaskParameters{Inline: map[string]string{"foo": "bar"}},
				SecureParameters: apilifecycle.SecureParameters{Secret: "my-secret"},
			},
		},
	}
}

func makeConfigMap(name, namespace string) *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{
			"code": "console.log('hello');",
		},
	}
}

func makeTaskDefinitionWithServiceAccount(name, namespace, serviceAccountName string, token *bool, ttlSeconds *int32, imagePullSecrets []v1.LocalObjectReference) *apilifecycle.KeptnTaskDefinition {
	return &apilifecycle.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"label1": "label2",
			},
			Annotations: map[string]string{
				"annotation1": "annotation2",
			},
		},
		Spec: apilifecycle.KeptnTaskDefinitionSpec{
			Container: &apilifecycle.ContainerSpec{
				Container: &v1.Container{},
			},
			ImagePullSecrets: imagePullSecrets,
			ServiceAccount: &apilifecycle.ServiceAccountSpec{
				Name: serviceAccountName,
			},
			AutomountServiceAccountToken: &apilifecycle.AutomountServiceAccountTokenSpec{
				Type: token,
			},
			TTLSecondsAfterFinished: ttlSeconds,
		},
	}
}
