package component

import (
	"context"
	"os"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptntask"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/storage/names"
)

var _ = Describe("KeptnTaskController", Ordered, func() {
	var (
		name               string
		taskDefinitionName string
		namespace          string
		spanRecorder       *sdktest.SpanRecorder
		tracer             *otelsdk.TracerProvider
	)

	BeforeAll(func() {
		// setup once
		By("Waiting for Manager")
		Eventually(func() bool {
			return k8sManager != nil
		}).Should(Equal(true))

		By("Creating the Controller")
		_ = os.Setenv("FUNCTION_RUNNER_IMAGE", "my-image")

		spanRecorder = sdktest.NewSpanRecorder()
		tracer = otelsdk.NewTracerProvider(otelsdk.WithSpanProcessor(spanRecorder))

		////setup controllers here
		controllers := []interfaces.Controller{&keptntask.KeptnTaskReconciler{
			Client:        k8sManager.GetClient(),
			Scheme:        k8sManager.GetScheme(),
			Recorder:      k8sManager.GetEventRecorderFor("test-task-controller"),
			Log:           GinkgoLogr,
			Meters:        initKeptnMeters(),
			TracerFactory: &tracerFactory{tracer: tracer},
		}}
		setupManager(controllers) // we can register multiple time the same controller
	})

	BeforeEach(func() { // list var here they will be copied for every spec
		name = names.SimpleNameGenerator.GenerateName("test-task-reconciler-")
		taskDefinitionName = names.SimpleNameGenerator.GenerateName("my-taskdef-")
		namespace = "default" // namespaces are not deleted in the api so be careful
	})

	Describe("Creation of a Task", func() {
		var (
			taskDefinition *klcv1alpha3.KeptnTaskDefinition
			task           *klcv1alpha3.KeptnTask
		)
		Context("with an existing TaskDefinition", func() {
			BeforeEach(func() {
				taskDefinition = makeTaskDefinition(taskDefinitionName, namespace)
				task = makeTask(name, namespace, taskDefinition.Name)
			})

			It("should end up in a failed state if the created job fails", func() {
				By("Verifying that a job has been created")

				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespace,
						Name:      task.Name,
					}, task)
					g.Expect(err).To(BeNil())
					g.Expect(task.Status.JobName).To(Not(BeEmpty()))
				}, "10s").Should(Succeed())

				createdJob := &batchv1.Job{}

				err := k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: namespace,
					Name:      task.Status.JobName,
				}, createdJob)

				Expect(err).To(BeNil())

				By("Setting the Job Status to failed")
				createdJob.Status.Failed = 1

				err = k8sClient.Status().Update(context.TODO(), createdJob)
				Expect(err).To(BeNil())

				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespace,
						Name:      task.Name,
					}, task)
					g.Expect(err).To(BeNil())
					g.Expect(task.Status.Status).To(Equal(apicommon.StateFailed))
				}, "10s").Should(Succeed())
			})
			AfterEach(func() {
				err := k8sClient.Delete(context.TODO(), taskDefinition)
				logErrorIfPresent(err)
				err = k8sClient.Delete(context.TODO(), task)
				logErrorIfPresent(err)
			})
		})
	})
})

func makeTask(name string, namespace, taskDefinitionName string) *klcv1alpha3.KeptnTask {
	task := &klcv1alpha3.KeptnTask{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: klcv1alpha3.KeptnTaskSpec{
			Workload:       "my-workload",
			AppName:        "my-app",
			AppVersion:     "0.1.0",
			TaskDefinition: taskDefinitionName,
		},
	}

	err := k8sClient.Create(ctx, task)
	Expect(err).To(BeNil())

	return task
}

func makeTaskDefinition(taskDefinitionName, namespace string) *klcv1alpha3.KeptnTaskDefinition {
	cmName := "my-cm"
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			Namespace: namespace,
		},
		Data: map[string]string{
			"code": "console.log('hello');",
		},
	}
	err := k8sClient.Create(context.TODO(), cm)

	Expect(err).To(BeNil())

	taskDefinition := &klcv1alpha3.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      taskDefinitionName,
			Namespace: namespace,
		},
		Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
			Function: klcv1alpha3.FunctionSpec{
				ConfigMapReference: klcv1alpha3.ConfigMapReference{
					Name: cmName,
				},
			},
		},
	}

	err = k8sClient.Create(context.TODO(), taskDefinition)
	Expect(err).To(BeNil())

	return taskDefinition
}
