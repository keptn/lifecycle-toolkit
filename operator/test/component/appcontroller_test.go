package component

import (
	"fmt"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnapp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/storage/names"
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
		controllers := []interfaces.Controller{&keptnapp.KeptnAppReconciler{
			Client:        k8sManager.GetClient(),
			Scheme:        k8sManager.GetScheme(),
			Recorder:      k8sManager.GetEventRecorderFor("test-app-controller"),
			Log:           GinkgoLogr,
			TracerFactory: &tracerFactory{tracer: tracer},
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
		name = names.SimpleNameGenerator.GenerateName("my-app-")
		namespace = "default" // namespaces are not deleted in the api so be careful
		// when creating you can use ignoreAlreadyExists(err error)
		version = "1.0.0"
	})
	Describe("Creation of AppVersion from a new App", func() {
		var (
			instance *klcv1alpha3.KeptnApp
		)

		BeforeEach(func() {
			instance = createInstanceInCluster(name, namespace, version)
			fmt.Println("created ", instance.Name)
		})

		Context("with a new App CRD", func() {

			It("should update the spans", func() {
				By("creating a new app version")
				assertResourceUpdated(instance)
				assertAppSpan(instance, spanRecorder)
				fmt.Println("spanned ", instance.Name)
			})

		})
		AfterEach(func() {
			// Remember to clean up the cluster after each test
			deleteAppInCluster(instance)
			// Reset span recorder after each spec
			resetSpanRecords(tracer, spanRecorder)
		})

	})
})

func deleteAppInCluster(instance *klcv1alpha3.KeptnApp) {
	By("Cleaning Up KeptnApp CRD ")
	err := k8sClient.Delete(ctx, instance)
	logErrorIfPresent(err)
}

func assertResourceUpdated(instance *klcv1alpha3.KeptnApp) *klcv1alpha3.KeptnAppVersion {

	appVersion := getAppVersion(instance)

	By("Comparing expected app version")
	Expect(appVersion.Spec.AppName).To(Equal(instance.Name))
	Expect(appVersion.Spec.Version).To(Equal(instance.Spec.Version))
	Expect(appVersion.Spec.Workloads).To(Equal(instance.Spec.Workloads))

	return appVersion
}

func getAppVersion(instance *klcv1alpha3.KeptnApp) *klcv1alpha3.KeptnAppVersion {
	appvName := types.NamespacedName{
		Namespace: instance.Namespace,
		Name:      fmt.Sprintf("%s-%s-%d", instance.Name, instance.Spec.Version, instance.Generation),
	}

	appVersion := &klcv1alpha3.KeptnAppVersion{}
	By("Retrieving Created app version")
	Eventually(func() error {
		return k8sClient.Get(ctx, appvName, appVersion)
	}, "20s").Should(Succeed())

	return appVersion
}

func assertAppSpan(instance *klcv1alpha3.KeptnApp, spanRecorder *sdktest.SpanRecorder) {
	By("Comparing spans")
	var spans []otelsdk.ReadOnlySpan
	Eventually(func() bool {
		spans = spanRecorder.Ended()
		return len(spans) >= 3
	}, "10s").Should(BeTrue())

	Expect(spans[0].Name()).To(Equal(fmt.Sprintf("%s-%s-%d", instance.Name, instance.Spec.Version, instance.Generation)))
	Expect(spans[0].Attributes()).To(ContainElement(apicommon.AppName.String(instance.Name)))
	Expect(spans[0].Attributes()).To(ContainElement(apicommon.AppVersion.String(instance.Spec.Version)))

	Expect(spans[1].Name()).To(Equal("create_app_version"))
	Expect(spans[1].Attributes()).To(ContainElement(apicommon.AppName.String(instance.Name)))
	Expect(spans[1].Attributes()).To(ContainElement(apicommon.AppVersion.String(instance.Spec.Version)))

	Expect(spans[2].Name()).To(Equal("reconcile_app"))
	Expect(spans[2].Attributes()).To(ContainElement(apicommon.AppName.String(instance.Name)))
	Expect(spans[2].Attributes()).To(ContainElement(apicommon.AppVersion.String(instance.Spec.Version)))
}

func createInstanceInCluster(name string, namespace string, version string) *klcv1alpha3.KeptnApp {
	instance := &klcv1alpha3.KeptnApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:       name,
			Namespace:  namespace,
			Generation: 1,
		},
		Spec: klcv1alpha3.KeptnAppSpec{
			Version: version,
			Workloads: []klcv1alpha3.KeptnWorkloadRef{
				{
					Name:    "app-wname",
					Version: "2.0",
				},
			},
		},
	}
	By("Invoking Reconciling for Create")

	Expect(k8sClient.Create(ctx, instance)).Should(Succeed())
	return instance
}
