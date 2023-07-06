package keptntask

import (
	"context"
	"testing"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestKeptnTaskReconciler_createJob(t *testing.T) {
	namespace := "default"
	cmName := "my-cmd"
	taskDefinitionName := "my-task-definition"

	cm := makeConfigMap(cmName, namespace)

	fakeClient := fake.NewClientBuilder().WithObjects(cm).Build()

	err := klcv1alpha3.AddToScheme(fakeClient.Scheme())
	require.Nil(t, err)

	taskDefinition := makeTaskDefinitionWithConfigmapRef(taskDefinitionName, namespace, cmName)

	err = fakeClient.Create(context.TODO(), taskDefinition)
	require.Nil(t, err)

	taskDefinition.Status.Function.ConfigMap = cmName
	err = fakeClient.Status().Update(context.TODO(), taskDefinition)
	require.Nil(t, err)

	r := &KeptnTaskReconciler{
		Client:      fakeClient,
		EventSender: controllercommon.NewEventSender(record.NewFakeRecorder(100)),
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

	fakeClient := fake.NewClientBuilder().WithObjects(cm).Build()

	err := klcv1alpha3.AddToScheme(fakeClient.Scheme())
	require.Nil(t, err)

	taskDefinition := makeTaskDefinitionWithConfigmapRef(taskDefinitionName, common.KLTNamespace, cmName)

	err = fakeClient.Create(context.TODO(), taskDefinition)
	require.Nil(t, err)

	taskDefinition.Status.Function.ConfigMap = cmName
	err = fakeClient.Status().Update(context.TODO(), taskDefinition)
	require.Nil(t, err)

	r := &KeptnTaskReconciler{
		Client:      fakeClient,
		EventSender: controllercommon.NewEventSender(record.NewFakeRecorder(100)),
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

	job := makeJob("my.job", namespace)

	fakeClient := fake.NewClientBuilder().WithObjects(job).Build()

	err := klcv1alpha3.AddToScheme(fakeClient.Scheme())
	require.Nil(t, err)

	job.Status.Conditions = []batchv1.JobCondition{
		{
			Type: batchv1.JobFailed,
		},
	}

	err = fakeClient.Status().Update(context.TODO(), job)
	require.Nil(t, err)

	r := &KeptnTaskReconciler{
		Client:      fakeClient,
		EventSender: controllercommon.NewEventSender(record.NewFakeRecorder(100)),
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

	err = fakeClient.Status().Update(context.TODO(), job)
	require.Nil(t, err)

	r.updateTaskStatus(job, task)

	require.Equal(t, apicommon.StateSucceeded, task.Status.Status)
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

func makeTask(name, namespace string, taskDefinitionName string) *klcv1alpha3.KeptnTask {
	return &klcv1alpha3.KeptnTask{
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
		Spec: klcv1alpha3.KeptnTaskSpec{
			Workload:       "my-workload",
			AppName:        "my-app",
			AppVersion:     "0.1.0",
			TaskDefinition: taskDefinitionName,
			Type:           apicommon.PostDeploymentCheckType,
		},
	}
}

func makeTaskDefinitionWithConfigmapRef(name, namespace, configMapName string) *klcv1alpha3.KeptnTaskDefinition {
	return &klcv1alpha3.KeptnTaskDefinition{
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
		Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
			Function: &klcv1alpha3.RuntimeSpec{
				ConfigMapReference: klcv1alpha3.ConfigMapReference{
					Name: configMapName,
				},
				Parameters:       klcv1alpha3.TaskParameters{Inline: map[string]string{"foo": "bar"}},
				SecureParameters: klcv1alpha3.SecureParameters{Secret: "my-secret"},
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
