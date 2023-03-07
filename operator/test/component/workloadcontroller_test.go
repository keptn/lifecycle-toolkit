package component

import (
	"context"
	"fmt"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnworkload"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// clean example of component test (E2E test/ integration test can be achieved adding a real cluster)
// Workload controller creates WorkloadVersion when a new Workload CRD is added
// span for creation and reconcile are correct
// container must be ordered to have the before all setup
// this way the container spec check is not randomized, so we can make
// assertions on spans number and traces
var _ = Describe("KeptnWorkloadController", Ordered, func() {
	var (
		name            string
		namespace       string
		version         string
		applicationName string
		spanRecorder    *sdktest.SpanRecorder
		tracer          *otelsdk.TracerProvider
	)

	BeforeAll(func() {
		//setup once
		By("Waiting for Manager")
		Eventually(func() bool {
			return k8sManager != nil
		}).Should(Equal(true))

		By("Creating the Controller")

		spanRecorder = sdktest.NewSpanRecorder()
		tracer = otelsdk.NewTracerProvider(otelsdk.WithSpanProcessor(spanRecorder))

		////setup controllers here
		controllers := []interfaces.Controller{&keptnworkload.KeptnWorkloadReconciler{
			Client:        k8sManager.GetClient(),
			Scheme:        k8sManager.GetScheme(),
			Recorder:      k8sManager.GetEventRecorderFor("test-workload-controller"),
			Log:           GinkgoLogr,
			TracerFactory: &tracerFactory{tracer: tracer},
		}}
		setupManager(controllers) // we can register multiple time the same controller
		// so that they have a different span/trace
	})

	BeforeEach(func() { // list var here they will be copied for every spec
		name = "my-workload"
		applicationName = "my-app"
		namespace = "default" // namespaces are not deleted in the api so be careful
		// when creating you can use ignoreAlreadyExists(err error)
		version = "1.0.0"
	})
	Describe("Creation of WorkloadInstance from a new Workload", func() {
		var (
			workload         *klcv1alpha3.KeptnWorkload
			workloadInstance *klcv1alpha3.KeptnWorkloadInstance
		)

		BeforeEach(func() {
			workload = createWorkloadInCluster(name, namespace, version, applicationName)
		})

		Context("with a new Workload CRD", func() {
			It("should update the spans and create WorkloadInstance", func() {
				By("Check if WorkloadInstance was created")

				workloadInstance = &klcv1alpha3.KeptnWorkloadInstance{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespace,
						Name:      fmt.Sprintf("%s-%s", workload.Name, workload.Spec.Version),
					}, workloadInstance)
					g.Expect(err).To(BeNil())
					g.Expect(workloadInstance.Spec.WorkloadName).To(Equal(workload.Name))
					g.Expect(workloadInstance.Spec.KeptnWorkloadSpec).To(Equal(workload.Spec))

				}, "30s").Should(Succeed())

				By("Comparing spans")
				var spans []otelsdk.ReadWriteSpan
				Eventually(func() bool {
					spans = spanRecorder.Started()
					return len(spans) >= 2
				}, "10s").Should(BeTrue())

				Expect(spans[0].Name()).To(Equal("reconcile_workload"))
				Expect(spans[0].Attributes()).To(ContainElement(apicommon.WorkloadName.String(workload.Name)))
				Expect(spans[0].Attributes()).To(ContainElement(apicommon.WorkloadVersion.String(workload.Spec.Version)))
				Expect(spans[0].Attributes()).To(ContainElement(apicommon.AppName.String(workload.Spec.AppName)))

				Expect(spans[1].Name()).To(Equal("create_workload_instance"))
				Expect(spans[1].Attributes()).To(ContainElement(apicommon.WorkloadName.String(workload.Name)))
				Expect(spans[1].Attributes()).To(ContainElement(apicommon.WorkloadVersion.String(workload.Spec.Version)))
				Expect(spans[0].Attributes()).To(ContainElement(apicommon.AppName.String(workload.Spec.AppName)))
			})

		})
		AfterEach(func() {
			By("Cleaning Up KeptnWorkload CRD")
			err := k8sClient.Delete(ctx, workload)
			logErrorIfPresent(err)
			By("Cleaning Up KeptnWorkloadInstance CRD")
			err = k8sClient.Delete(ctx, workloadInstance)
			logErrorIfPresent(err)
		})

	})
})

func createWorkloadInCluster(name string, namespace string, version string, applicationName string) *klcv1alpha3.KeptnWorkload {
	workload := &klcv1alpha3.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: klcv1alpha3.KeptnWorkloadSpec{
			AppName:           applicationName,
			Version:           version,
			ResourceReference: klcv1alpha3.ResourceReference{UID: types.UID("uid"), Kind: "Pod", Name: "pod1"},
		},
	}
	By("Invoking Reconciling for Create")

	Expect(k8sClient.Create(ctx, workload)).Should(Succeed())
	return workload
}
