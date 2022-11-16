package keptntask

import (
	"context"
	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestKeptnTaskReconciler_createJob(t *testing.T) {
	namespace := "default"
	cmName := "my-cmd"
	taskDefinitionName := "my-task-definition"

	cm := makeConfigMap(cmName, namespace)

	fakeClient := fake.NewClientBuilder().WithObjects(cm).Build()

	fakeRecorder := &record.FakeRecorder{}

	err := klcv1alpha1.AddToScheme(fakeClient.Scheme())
	require.Nil(t, err)

	taskDefinition := makeTaskDefinitionWithConfigmapRef(taskDefinitionName, namespace, cmName)

	err = fakeClient.Create(context.TODO(), taskDefinition)
	require.Nil(t, err)

	taskDefinition.Status.Function.ConfigMap = cmName
	err = fakeClient.Status().Update(context.TODO(), taskDefinition)
	require.Nil(t, err)

	r := &KeptnTaskReconciler{
		Client:   fakeClient,
		Recorder: fakeRecorder,
		Log:      ctrl.Log.WithName("task-controller"),
		Scheme:   fakeClient.Scheme(),
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
}

func TestKeptnTaskReconciler_updateJob(t *testing.T) {
	namespace := "default"
	taskDefinitionName := "my-task-definition"

	job := makeJob("my.job", namespace)

	fakeClient := fake.NewClientBuilder().WithObjects(job).Build()

	fakeRecorder := &record.FakeRecorder{}

	err := klcv1alpha1.AddToScheme(fakeClient.Scheme())
	require.Nil(t, err)

	job.Status.Failed = 1

	err = fakeClient.Status().Update(context.TODO(), job)
	require.Nil(t, err)

	r := &KeptnTaskReconciler{
		Client:   fakeClient,
		Recorder: fakeRecorder,
		Log:      ctrl.Log.WithName("task-controller"),
		Scheme:   fakeClient.Scheme(),
	}

	task := makeTask("my-task", namespace, taskDefinitionName)

	err = fakeClient.Create(context.TODO(), task)
	require.Nil(t, err)

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: namespace,
		},
	}

	task.Status.JobName = job.Name

	err = r.updateJob(context.TODO(), req, task)
	require.Nil(t, err)

	require.Equal(t, common.StateFailed, task.Status.Status)

	// now, set the job to succeeded
	job.Status.Succeeded = 1
	job.Status.Failed = 0

	err = fakeClient.Status().Update(context.TODO(), job)
	require.Nil(t, err)

	err = r.updateJob(context.TODO(), req, task)
	require.Nil(t, err)

	require.Equal(t, common.StateSucceeded, task.Status.Status)
}

func makeJob(name, namespace string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{},
	}
}

func makeTask(name, namespace string, taskDefinitionName string) *klcv1alpha1.KeptnTask {
	return &klcv1alpha1.KeptnTask{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: klcv1alpha1.KeptnTaskSpec{
			Workload:       "my-workload",
			AppName:        "my-app",
			AppVersion:     "0.1.0",
			TaskDefinition: taskDefinitionName,
		},
	}
}

func makeTaskDefinitionWithConfigmapRef(name, namespace, configMapName string) *klcv1alpha1.KeptnTaskDefinition {
	return &klcv1alpha1.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: klcv1alpha1.KeptnTaskDefinitionSpec{
			Function: klcv1alpha1.FunctionSpec{
				ConfigMapReference: klcv1alpha1.ConfigMapReference{
					Name: configMapName,
				},
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
