package component

import (
	"strings"
	"time"

	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/interfaces"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/keptnappversion"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// clean example of component test (E2E test/ integration test can be achieved adding a real cluster)
// App controller creates AppVersion when a new App CRD is added
// span for creation and reconcile are correct
// container must be ordered to have the before all setup
// this way the container spec check is not randomized, so we can make
// assertions on spans number and traces
var _ = Describe("KeptnAppVersionController", Ordered, func() {
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
		controllers := []interfaces.Controller{&keptnappversion.KeptnAppVersionReconciler{
			Client:      k8sManager.GetClient(),
			Scheme:      k8sManager.GetScheme(),
			Recorder:    k8sManager.GetEventRecorderFor("test-app-controller"),
			Log:         GinkgoLogr,
			Meters:      initKeptnMeters(),
			SpanHandler: &controllercommon.SpanHandler{},
			Tracer:      tracer.Tracer("test-app-tracer"),
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
	Describe("Creation of AppVersion", func() {
		var (
			appVersion *klcv1alpha1.KeptnAppVersion
		)
		Context("with a new AppVersions CRD", func() {

			It("should be cancelled when pre-eval checks failed", func() {

				By("Creating a new App Version")

				appVersion = &klcv1alpha1.KeptnAppVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Spec: klcv1alpha1.KeptnAppVersionSpec{
						AppName: name,
						KeptnAppSpec: klcv1alpha1.KeptnAppSpec{
							Version:                  version,
							PreDeploymentEvaluations: []string{"eval-def"},
							Workloads: []klcv1alpha1.KeptnWorkloadRef{
								{
									Name:    "wname",
									Version: "2.0",
								},
							},
						},
					},
				}
				Expect(ignoreAlreadyExists(k8sClient.Create(ctx, appVersion))).Should(Succeed())

				appVersionNameObj := types.NamespacedName{
					Namespace: appVersion.Namespace,
					Name:      appVersion.Name,
				}

				By("Ensuring an evaluation has been created")

				evaluation := &klcv1alpha1.KeptnEvaluation{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(ctx, appVersionNameObj, appVersion)
					g.Expect(err).To(BeNil())
					g.Expect(appVersion.Status.PreDeploymentEvaluationTaskStatus).To(Not(BeEmpty()))
					err = k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: appVersion.Status.PreDeploymentEvaluationTaskStatus[0].EvaluationName}, evaluation)
					g.Expect(err).ToNot(BeNil())
				}, "30s").Should(Succeed())

				By("Updating Evaluation status")

				evaluation.Status = klcv1alpha1.KeptnEvaluationStatus{
					OverallStatus: apicommon.StateFailed,
					RetryCount:    10,
					EvaluationStatus: map[string]klcv1alpha1.EvaluationStatusItem{
						"something": {
							Status: apicommon.StateFailed,
							Value:  "10",
						},
					},
					StartTime: metav1.Time{Time: time.Now().UTC()},
					EndTime:   metav1.Time{Time: time.Now().UTC().Add(5 * time.Second)},
				}

				err := k8sClient.Status().Update(ctx, evaluation)
				Expect(err).To(BeNil())

				By("Ensuring all phases after pre-eval checks are cancelled")

				Eventually(func(g Gomega) {
					appVersion := &klcv1alpha1.KeptnAppVersion{}
					err := k8sClient.Get(ctx, appVersionNameObj, appVersion)
					g.Expect(err).To(BeNil())
					g.Expect(appVersion).To(Not(BeNil()))
					g.Expect(appVersion.Status.PreDeploymentStatus).To(BeEquivalentTo(apicommon.StateSucceeded))
					g.Expect(appVersion.Status.PreDeploymentEvaluationStatus).To(BeEquivalentTo(apicommon.StateFailed))
					g.Expect(appVersion.Status.PostDeploymentStatus).To(BeEquivalentTo(apicommon.StateCancelled))
					g.Expect(appVersion.Status.PostDeploymentEvaluationStatus).To(BeEquivalentTo(apicommon.StateCancelled))
					g.Expect(appVersion.Status.Status).To(BeEquivalentTo(apicommon.StateFailed))
				}, "30s").Should(Succeed())

				By("Ensuring that a K8s Event containing the reason for the failed evaluation has been sent")

				Eventually(func(g Gomega) {
					eventList := &corev1.EventList{}
					err := k8sClient.List(ctx, eventList, client.InNamespace(namespace))
					g.Expect(err).To(BeNil())

					foundEvent := &corev1.Event{}

					for _, e := range eventList.Items {
						if strings.Contains(e.Name, appVersion.GetName()) && e.Reason == "AppPreDeployEvaluationsFailed" {
							foundEvent = &e
							break
						}
					}
					g.Expect(foundEvent).NotTo(BeNil())
				}, "30s").Should(Succeed())
			})
			AfterEach(func() {
				// Remember to clean up the cluster after each test
				err := k8sClient.Delete(ctx, appVersion)
				logErrorIfPresent(err)
				// Reset span recorder after each spec
				resetSpanRecords(tracer, spanRecorder)
			})

		})

	})
})
