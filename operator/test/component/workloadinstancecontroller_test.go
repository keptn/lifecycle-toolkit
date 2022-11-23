package component

import (
	"context"
	"time"

	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	keptncontroller "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/keptnworkloadinstance"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func getPodTemplateSpec() corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "nginx",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}
}

// clean example of component test (E2E test/ integration test can be achieved adding a real cluster)
// App controller creates AppVersion when a new App CRD is added
// span for creation and reconcile are correct
// container must be ordered to have the before all setup
// this way the container spec check is not randomized, so we can make
// assertions on spans number and traces
var _ = Describe("KeptnWorkloadInstanceController", Ordered, func() {
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
		controllers := []keptncontroller.Controller{&keptnworkloadinstance.KeptnWorkloadInstanceReconciler{
			Client:      k8sManager.GetClient(),
			Scheme:      k8sManager.GetScheme(),
			Recorder:    k8sManager.GetEventRecorderFor("test-app-controller"),
			Log:         GinkgoLogr,
			Meters:      initKeptnMeters(),
			SpanHandler: &keptncontroller.SpanHandler{},
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
	Describe("Creation of WorkloadInstance", func() {
		var (
			appVersion *klcv1alpha1.KeptnAppVersion
			wi         *klcv1alpha1.KeptnWorkloadInstance
		)
		Context("with a new AppVersions CRD", func() {

			BeforeEach(func() {
				appVersion = createAppVersionInCluster(name, namespace, version)
			})

			It("should fail if Workload not found in AppVersion", func() {
				wiName := "not-found"
				wi = &klcv1alpha1.KeptnWorkloadInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Spec: klcv1alpha1.KeptnWorkloadInstanceSpec{
						KeptnWorkloadSpec: klcv1alpha1.KeptnWorkloadSpec{},
						WorkloadName:      "wi-test-app-wname-" + wiName,
						TraceId:           map[string]string{"traceparent": "00-0f89f15e562489e2e171eca1cf9ba958-d2fa6dbbcbf7e29a-01"},
					},
				}
				By("Creating WorkloadInstance")
				err := k8sClient.Create(context.TODO(), wi)
				Expect(err).To(BeNil())

				By("Ensuring WorkloadInstance does not progress to next phase")
				wiNameObj := types.NamespacedName{
					Namespace: wi.Namespace,
					Name:      wi.Name,
				}
				Consistently(func(g Gomega) {
					wi := &klcv1alpha1.KeptnWorkloadInstance{}
					err := k8sClient.Get(ctx, wiNameObj, wi)
					g.Expect(err).To(BeNil())
					g.Expect(wi).To(Not(BeNil()))
					g.Expect(wi.Status.CurrentPhase).To(BeEmpty())
				}, "3s").Should(Succeed())
			})

			It("should detect that the referenced StatefulSet is progressing", func() {
				By("Deploying a StatefulSet to reference")
				repl := int32(1)
				statefulSet := &appsv1.StatefulSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-statefulset",
						Namespace: namespace,
					},
					Spec: appsv1.StatefulSetSpec{
						Replicas: &repl,
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx",
							},
						},
						Template: getPodTemplateSpec(),
					},
				}

				defer func() {
					_ = k8sClient.Delete(ctx, statefulSet)
				}()

				err := k8sClient.Create(ctx, statefulSet)
				Expect(err).To(BeNil())

				By("Setting the App PreDeploymentEvaluation Status to 'Succeeded'")
				appVersion.Status.PreDeploymentEvaluationStatus = common.StateSucceeded
				err = k8sClient.Status().Update(ctx, appVersion)
				Expect(err).To(BeNil())

				By("Bringing the StatefulSet into its ready state")
				statefulSet.Status.AvailableReplicas = 1
				statefulSet.Status.ReadyReplicas = 1
				statefulSet.Status.Replicas = 1
				err = k8sClient.Status().Update(ctx, statefulSet)
				Expect(err).To(BeNil())

				By("Looking up the StatefulSet to retrieve its UID")
				err = k8sClient.Get(ctx, types.NamespacedName{
					Namespace: namespace,
					Name:      statefulSet.Name,
				}, statefulSet)
				Expect(err).To(BeNil())

				By("Creating a WorkloadInstance that references the StatefulSet")
				wi = &klcv1alpha1.KeptnWorkloadInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Spec: klcv1alpha1.KeptnWorkloadInstanceSpec{
						KeptnWorkloadSpec: klcv1alpha1.KeptnWorkloadSpec{
							ResourceReference: klcv1alpha1.ResourceReference{
								UID:  statefulSet.UID,
								Kind: "StatefulSet",
								Name: "my-statefulset",
							},
							Version: "2.0",
							AppName: appVersion.GetAppName(),
						},
						WorkloadName: "test-app-wname",
						TraceId:      map[string]string{"traceparent": "00-0f89f15e562489e2e171eca1cf9ba958-d2fa6dbbcbf7e29a-01"},
					},
				}

				err = k8sClient.Create(context.TODO(), wi)
				Expect(err).To(BeNil())

				wiNameObj := types.NamespacedName{
					Namespace: wi.Namespace,
					Name:      wi.Name,
				}
				Eventually(func(g Gomega) {
					wi := &klcv1alpha1.KeptnWorkloadInstance{}
					err := k8sClient.Get(ctx, wiNameObj, wi)
					g.Expect(err).To(BeNil())
					g.Expect(wi.Status.DeploymentStatus).To(Equal(common.StateSucceeded))
				}, "20s").Should(Succeed())
			})
			It("should detect that the referenced DaemonSet is progressing", func() {
				By("Deploying a DaemonSet to reference")
				daemonSet := &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-daemonset",
						Namespace: namespace,
					},
					Spec: appsv1.DaemonSetSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx",
							},
						},
						Template: getPodTemplateSpec(),
					},
				}

				defer func() {
					_ = k8sClient.Delete(ctx, daemonSet)
				}()

				err := k8sClient.Create(ctx, daemonSet)
				Expect(err).To(BeNil())

				By("Setting the App PreDeploymentEvaluation Status to 'Succeeded'")
				appVersion.Status.PreDeploymentEvaluationStatus = common.StateSucceeded
				err = k8sClient.Status().Update(ctx, appVersion)
				Expect(err).To(BeNil())

				By("Bringing the DaemonSet into its ready state")
				daemonSet.Status.DesiredNumberScheduled = 1
				daemonSet.Status.NumberReady = 1
				err = k8sClient.Status().Update(ctx, daemonSet)
				Expect(err).To(BeNil())

				By("Looking up the DaemonSet to retrieve its UID")
				err = k8sClient.Get(ctx, types.NamespacedName{
					Namespace: namespace,
					Name:      daemonSet.Name,
				}, daemonSet)
				Expect(err).To(BeNil())

				By("Creating a WorkloadInstance that references the DaemonSet")
				wi = &klcv1alpha1.KeptnWorkloadInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Spec: klcv1alpha1.KeptnWorkloadInstanceSpec{
						KeptnWorkloadSpec: klcv1alpha1.KeptnWorkloadSpec{
							ResourceReference: klcv1alpha1.ResourceReference{
								UID:  daemonSet.UID,
								Kind: "DaemonSet",
								Name: "my-daemonset",
							},
							Version: "2.0",
							AppName: appVersion.GetAppName(),
						},
						WorkloadName: "test-app-wname",
						TraceId:      map[string]string{"traceparent": "00-0f89f15e562489e2e171eca1cf9ba958-d2fa6dbbcbf7e29a-01"},
					},
				}

				err = k8sClient.Create(context.TODO(), wi)
				Expect(err).To(BeNil())

				wiNameObj := types.NamespacedName{
					Namespace: wi.Namespace,
					Name:      wi.Name,
				}
				Eventually(func(g Gomega) {
					wi := &klcv1alpha1.KeptnWorkloadInstance{}
					err := k8sClient.Get(ctx, wiNameObj, wi)
					g.Expect(err).To(BeNil())
					g.Expect(wi.Status.DeploymentStatus).To(Equal(common.StateSucceeded))
				}, "20s").Should(Succeed())
			})
			It("should be cancelled when pre-eval checks failed", func() {
				evaluation := &klcv1alpha1.KeptnEvaluation{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pre-eval-eval-def",
						Namespace: namespace,
					},
					Spec: klcv1alpha1.KeptnEvaluationSpec{
						EvaluationDefinition: "eval-def",
						Workload:             "test-app-wname",
						WorkloadVersion:      "2.0",
						Type:                 common.PreDeploymentEvaluationCheckType,
						Retries:              10,
					},
				}

				defer func() {
					_ = k8sClient.Delete(ctx, evaluation)
				}()

				By("Creating Evaluation")
				err := k8sClient.Create(context.TODO(), evaluation)
				Expect(err).To(BeNil())

				evaluation.Status = klcv1alpha1.KeptnEvaluationStatus{
					OverallStatus: common.StateFailed,
					RetryCount:    10,
					EvaluationStatus: map[string]klcv1alpha1.EvaluationStatusItem{
						"something": {
							Status: common.StateFailed,
							Value:  "10",
						},
					},
					StartTime: metav1.Time{Time: time.Now().UTC()},
					EndTime:   metav1.Time{Time: time.Now().UTC().Add(5 * time.Second)},
				}

				err = k8sClient.Status().Update(ctx, evaluation)
				Expect(err).To(BeNil())

				wi = &klcv1alpha1.KeptnWorkloadInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-app-wname-2.0",
						Namespace: namespace,
					},
					Spec: klcv1alpha1.KeptnWorkloadInstanceSpec{
						KeptnWorkloadSpec: klcv1alpha1.KeptnWorkloadSpec{
							Version:                  "2.0",
							AppName:                  appVersion.GetAppName(),
							PreDeploymentEvaluations: []string{"eval-def"},
						},
						WorkloadName: "test-app-wname",
					},
				}
				By("Creating WorkloadInstance")
				err = k8sClient.Create(context.TODO(), wi)
				Expect(err).To(BeNil())

				wi.Status = klcv1alpha1.KeptnWorkloadInstanceStatus{
					PreDeploymentStatus:            common.StateSucceeded,
					PreDeploymentEvaluationStatus:  common.StateProgressing,
					DeploymentStatus:               common.StatePending,
					PostDeploymentStatus:           common.StatePending,
					PostDeploymentEvaluationStatus: common.StatePending,
					CurrentPhase:                   common.PhaseWorkloadPreEvaluation.ShortName,
					Status:                         common.StateProgressing,
					PreDeploymentEvaluationTaskStatus: []klcv1alpha1.EvaluationStatus{
						{
							EvaluationName:           "pre-eval-eval-def",
							Status:                   common.StateProgressing,
							EvaluationDefinitionName: "eval-def",
						},
					},
				}

				err = k8sClient.Status().Update(ctx, wi)
				Expect(err).To(BeNil())

				By("Ensuring all phases after pre-eval checks are cancelled")
				wiNameObj := types.NamespacedName{
					Namespace: wi.Namespace,
					Name:      wi.Name,
				}
				Eventually(func(g Gomega) {
					wi := &klcv1alpha1.KeptnWorkloadInstance{}
					err := k8sClient.Get(ctx, wiNameObj, wi)
					g.Expect(err).To(BeNil())
					g.Expect(wi).To(Not(BeNil()))
					g.Expect(wi.Status.PreDeploymentStatus).To(BeEquivalentTo(common.StateSucceeded))
					g.Expect(wi.Status.PreDeploymentEvaluationStatus).To(BeEquivalentTo(common.StateFailed))
					g.Expect(wi.Status.DeploymentStatus).To(BeEquivalentTo(common.StateCancelled))
					g.Expect(wi.Status.PostDeploymentStatus).To(BeEquivalentTo(common.StateCancelled))
					g.Expect(wi.Status.PostDeploymentEvaluationStatus).To(BeEquivalentTo(common.StateCancelled))
					g.Expect(wi.Status.Status).To(BeEquivalentTo(common.StateFailed))
				}, "30s").Should(Succeed())
			})
			AfterEach(func() {
				// Remember to clean up the cluster after each test
				err := k8sClient.Delete(ctx, appVersion)
				logErrorIfPresent(err)
				err = k8sClient.Delete(ctx, wi)
				logErrorIfPresent(err)
				// Reset span recorder after each spec
				resetSpanRecords(tracer, spanRecorder)
			})

		})

	})
})

func createAppVersionInCluster(name string, namespace string, version string) *klcv1alpha1.KeptnAppVersion {
	instance := &klcv1alpha1.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: klcv1alpha1.KeptnAppVersionSpec{
			AppName: name,
			KeptnAppSpec: klcv1alpha1.KeptnAppSpec{
				Version: version,
				Workloads: []klcv1alpha1.KeptnWorkloadRef{
					{
						Name:    "wname",
						Version: "2.0",
					},
				},
			},
		},
	}
	By("Invoking Reconciling for Create")

	Expect(ignoreAlreadyExists(k8sClient.Create(ctx, instance))).Should(Succeed())
	instance.Status.PreDeploymentEvaluationStatus = common.StateSucceeded
	_ = k8sClient.Status().Update(ctx, instance)
	return instance
}
