package component

import (
	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	keptncontroller "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/keptnapp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// clean example of component test (E2E test/ integration test can be achieved adding a real cluster)
// App controller creates AppVersion when a new App CRD is added
// span for creation and reconcile are correct
// container must be ordered to have the before all setup
// this way the container spec check is not randomized, so we can make
// assertions on spans number and traces
var _ = Describe("KeptnAppController", Ordered, func() {
	var (
		name         string
		namespace    string
		version      string
		spanRecorder *sdktest.SpanRecorder
		tracer       *otelsdk.TracerProvider
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
		controllers := []keptncontroller.Controller{&keptnapp.KeptnAppReconciler{
			Client:   k8sManager.GetClient(),
			Scheme:   k8sManager.GetScheme(),
			Recorder: k8sManager.GetEventRecorderFor("test-app-controller"),
			Log:      GinkgoLogr,
			Tracer:   tracer.Tracer("test-app-tracer"),
		}}
		setupManager(controllers) // we can register multiple time the same controller
		// so that they have a different span/trace

		//for a fake controller you can also use
		//controller, err := controller.New("app-controller", cm, controller.Options{
		//	Reconciler: reconcile.Func(
		//		func(_ context.Context, request reconcile.Request) (reconcile.Result, error) {
		//			reconciled <- request
		//			return reconcile.Result{}, nil
		//		}),
		//})
		//Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() { // list var here they will be copied for every spec
		name = "test-app"
		namespace = "default" // namespaces are not deleted in the api so be careful
		// when creating you can use ignoreAlreadyExists(err error)
		version = "1.0.0"
	})
	Describe("Creation of AppVersion from a new App", func() {
		var (
			instance   *klcv1alpha1.KeptnApp
			appVersion *klcv1alpha1.KeptnAppVersion
		)
		Context("with a new App CRD", func() {

			BeforeEach(func() {
				instance = createInstanceInCluster(name, namespace, version, instance)
			})

			It("should update the status of the CR ", func() {
				appVersion = assertResourceUpdated(instance)
			})
			It("should update the spans", func() {
				assertAppSpan(instance, spanRecorder)
			})

		})
		AfterEach(func() {
			// Remember to clean up the cluster after each test
			deleteAppInCluster(instance)
			deleteAppVersionInCluster(appVersion)
			// Reset span recorder after each spec
			resetSpanRecords(tracer, spanRecorder)
		})
	})
})

func deleteAppVersionInCluster(version *klcv1alpha1.KeptnAppVersion) {
	By("Cleaning Up Keptn AppVersion CRD")
	Expect(k8sClient.Delete(ctx, version)).Should(Succeed())
}

func deleteAppInCluster(instance *klcv1alpha1.KeptnApp) {
	By("Cleaning Up KeptnApp CRD ")
	Expect(k8sClient.Delete(ctx, instance)).Should(Succeed())

}

func assertResourceUpdated(instance *klcv1alpha1.KeptnApp) *klcv1alpha1.KeptnAppVersion {

	appVersion := &klcv1alpha1.KeptnAppVersion{}
	appvName := types.NamespacedName{
		Namespace: instance.Namespace,
		Name:      instance.Name + "-" + instance.Spec.Version,
	}
	By("Retrieving Created app version")
	Eventually(func() error {
		return k8sClient.Get(ctx, appvName, appVersion)
	}).Should(Succeed())

	By("Comparing expected app version")
	Expect(appVersion.Spec.AppName).To(Equal(instance.Name))
	Expect(appVersion.Spec.Version).To(Equal(instance.Spec.Version))
	Expect(appVersion.Spec.Workloads).To(Equal(instance.Spec.Workloads))

	return appVersion
}

func assertAppSpan(instance *klcv1alpha1.KeptnApp, spanRecorder *sdktest.SpanRecorder) {
	By("Comparing spans")
	var spans []otelsdk.ReadOnlySpan
	Eventually(func() int {
		spans = spanRecorder.Ended()
		return len(spans)
	}).Should(Equal(3))

	Expect(spans[0].Name()).To(Equal("appversion_deployment"))
	Expect(spans[0].Attributes()).To(ContainElement(common.AppName.String(instance.Name)))
	Expect(spans[0].Attributes()).To(ContainElement(common.AppVersion.String(instance.Spec.Version)))

	Expect(spans[1].Name()).To(Equal("create_app_version"))
	Expect(spans[1].Attributes()).To(ContainElement(common.AppName.String(instance.Name)))
	Expect(spans[1].Attributes()).To(ContainElement(common.AppVersion.String(instance.Spec.Version)))

	Expect(spans[2].Name()).To(Equal("reconcile_app"))
	Expect(spans[2].Attributes()).To(ContainElement(common.AppName.String(instance.Name)))
	Expect(spans[2].Attributes()).To(ContainElement(common.AppVersion.String(instance.Spec.Version)))

}

func createInstanceInCluster(name string, namespace string, version string, instance *klcv1alpha1.KeptnApp) *klcv1alpha1.KeptnApp {
	instance = &klcv1alpha1.KeptnApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: klcv1alpha1.KeptnAppSpec{
			Version: version,
			Workloads: []klcv1alpha1.KeptnWorkloadRef{
				{
					Name:    "app-wname",
					Version: "2.0",
				},
			},
		},
	}
	By("Invoking Reconciling for Create")

	Expect(ignoreAlreadyExists(k8sClient.Create(ctx, instance))).Should(Succeed())
	return instance
}
