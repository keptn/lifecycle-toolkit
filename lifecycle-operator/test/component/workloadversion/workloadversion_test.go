package workloadversion_test

import (
	"context"
	"strings"
	"time"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/storage/names"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("WorkloadVersion", Ordered, func() {
	var (
		appName   string
		namespace string
		version   string
	)

	BeforeEach(func() { // list var here they will be copied for every spec
		namespace = "default" // namespaces are not deleted in the api so be careful
		// when creating you can use ignoreAlreadyExists(err error)
		version = "1.0.0"
	})
	Describe("Creation of WorkloadVersion", func() {
		var (
			appVersion *klcv1alpha3.KeptnAppVersion
			wi         *klcv1alpha3.KeptnWorkloadVersion
		)
		Context("with a new AppVersions CRD", func() {

			BeforeEach(func() {
				appName = names.SimpleNameGenerator.GenerateName("test-app-")
				appVersion = createAppVersionInCluster(appName, namespace, version)
			})

			It("should fail if Workload not found in AppVersion", func() {
				wiName := "not-found"
				wi = &klcv1alpha3.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      appName,
						Namespace: namespace,
					},
					Spec: klcv1alpha3.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{},
						WorkloadName:      appName + "-wname-" + wiName,
						TraceId:           map[string]string{"traceparent": "00-0f89f15e562489e2e171eca1cf9ba958-d2fa6dbbcbf7e29a-01"},
					},
				}
				By("Creating WorkloadVersion")
				err := k8sClient.Create(context.TODO(), wi)
				Expect(err).To(BeNil())

				By("Ensuring WorkloadVersion does not progress to next phase")
				wiNameObj := types.NamespacedName{
					Namespace: wi.Namespace,
					Name:      wi.Name,
				}
				Consistently(func(g Gomega) {
					wi := &klcv1alpha3.KeptnWorkloadVersion{}
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

				av := &klcv1alpha3.KeptnAppVersion{}
				err = k8sClient.Get(ctx, types.NamespacedName{
					Namespace: namespace,
					Name:      appName,
				}, av)
				Expect(err).To(BeNil())

				av.Status.PreDeploymentEvaluationStatus = apicommon.StateSucceeded
				err = k8sClient.Status().Update(ctx, av)
				Expect(err).To(BeNil())

				By("Looking up the StatefulSet to retrieve its UID")
				err = k8sClient.Get(ctx, types.NamespacedName{
					Namespace: namespace,
					Name:      statefulSet.Name,
				}, statefulSet)
				Expect(err).To(BeNil())

				By("Bringing the StatefulSet into its ready state")
				statefulSet.Status.AvailableReplicas = 1
				statefulSet.Status.ReadyReplicas = 1
				statefulSet.Status.Replicas = 1
				err = k8sClient.Status().Update(ctx, statefulSet)
				Expect(err).To(BeNil())

				By("Creating a WorkloadVersion that references the StatefulSet")
				wi = &klcv1alpha3.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      appName,
						Namespace: namespace,
					},
					Spec: klcv1alpha3.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
							ResourceReference: klcv1alpha3.ResourceReference{
								UID:  statefulSet.UID,
								Kind: "StatefulSet",
								Name: "my-statefulset",
							},
							Version: "2.0",
							AppName: appVersion.GetAppName(),
						},
						WorkloadName: appName + "-wname",
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
					wi := &klcv1alpha3.KeptnWorkloadVersion{}
					err := k8sClient.Get(ctx, wiNameObj, wi)
					g.Expect(err).To(BeNil())
					g.Expect(wi.Status.DeploymentStatus).To(Equal(apicommon.StateSucceeded))
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

				av := &klcv1alpha3.KeptnAppVersion{}
				err = k8sClient.Get(ctx, types.NamespacedName{
					Namespace: namespace,
					Name:      appName,
				}, av)
				Expect(err).To(BeNil())

				av.Status.PreDeploymentEvaluationStatus = apicommon.StateSucceeded
				err = k8sClient.Status().Update(ctx, av)
				Expect(err).To(BeNil())

				By("Looking up the DaemonSet to retrieve its UID")
				err = k8sClient.Get(ctx, types.NamespacedName{
					Namespace: namespace,
					Name:      daemonSet.Name,
				}, daemonSet)
				Expect(err).To(BeNil())

				By("Bringing the DaemonSet into its ready state")
				daemonSet.Status.DesiredNumberScheduled = 1
				daemonSet.Status.NumberReady = 1
				err = k8sClient.Status().Update(ctx, daemonSet)
				Expect(err).To(BeNil())

				By("Creating a WorkloadVersion that references the DaemonSet")
				wi = &klcv1alpha3.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      appName,
						Namespace: namespace,
					},
					Spec: klcv1alpha3.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
							ResourceReference: klcv1alpha3.ResourceReference{
								UID:  daemonSet.UID,
								Kind: "DaemonSet",
								Name: "my-daemonset",
							},
							Version: "2.0",
							AppName: appVersion.GetAppName(),
						},
						WorkloadName: appName + "-wname",
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
					wi := &klcv1alpha3.KeptnWorkloadVersion{}
					err := k8sClient.Get(ctx, wiNameObj, wi)
					g.Expect(err).To(BeNil())
					g.Expect(wi.Status.DeploymentStatus).To(Equal(apicommon.StateSucceeded))
				}, "20s").Should(Succeed())
			})
			It("should be deprecated when pre-eval checks failed", func() {
				evaluation := &klcv1alpha3.KeptnEvaluation{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pre-eval-eval-def",
						Namespace: namespace,
					},
					Spec: klcv1alpha3.KeptnEvaluationSpec{
						EvaluationDefinition: "eval-def",
						Workload:             appName + "-wname",
						WorkloadVersion:      "2.0",
						Type:                 apicommon.PreDeploymentEvaluationCheckType,
						Retries:              10,
					},
				}

				defer func() {
					_ = k8sClient.Delete(ctx, evaluation)
				}()

				By("Creating Evaluation")
				err := k8sClient.Create(context.TODO(), evaluation)
				Expect(err).To(BeNil())

				err = k8sClient.Get(ctx, types.NamespacedName{
					Namespace: namespace,
					Name:      evaluation.Name,
				}, evaluation)
				Expect(err).To(BeNil())

				evaluation.Status = klcv1alpha3.KeptnEvaluationStatus{
					OverallStatus: apicommon.StateFailed,
					RetryCount:    10,
					EvaluationStatus: map[string]klcv1alpha3.EvaluationStatusItem{
						"something": {
							Status: apicommon.StateFailed,
							Value:  "10",
						},
					},
					StartTime: metav1.Time{Time: time.Now().UTC()},
					EndTime:   metav1.Time{Time: time.Now().UTC().Add(5 * time.Second)},
				}

				err = k8sClient.Status().Update(ctx, evaluation)
				Expect(err).To(BeNil())

				wi = &klcv1alpha3.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      appName + "-wname-2.0",
						Namespace: namespace,
					},
					Spec: klcv1alpha3.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
							Version:                  "2.0",
							AppName:                  appVersion.GetAppName(),
							PreDeploymentEvaluations: []string{"eval-def"},
						},
						WorkloadName: appName + "-wname",
					},
				}
				By("Creating WorkloadVersion")
				err = k8sClient.Create(context.TODO(), wi)
				Expect(err).To(BeNil())

				time.Sleep(5 * time.Second)

				wi2 := &klcv1alpha3.KeptnWorkloadVersion{}
				err = k8sClient.Get(ctx, types.NamespacedName{Namespace: wi.Namespace, Name: wi.Name}, wi2)
				Expect(err).To(BeNil())
				Expect(wi2).To(Not(BeNil()))

				wi2.Status = klcv1alpha3.KeptnWorkloadVersionStatus{
					PreDeploymentStatus:            apicommon.StateSucceeded,
					PreDeploymentEvaluationStatus:  apicommon.StateProgressing,
					DeploymentStatus:               apicommon.StatePending,
					PostDeploymentStatus:           apicommon.StatePending,
					PostDeploymentEvaluationStatus: apicommon.StatePending,
					CurrentPhase:                   apicommon.PhaseWorkloadPreEvaluation.ShortName,
					Status:                         apicommon.StateProgressing,
					PreDeploymentEvaluationTaskStatus: []klcv1alpha3.ItemStatus{
						{
							Name:           "pre-eval-eval-def",
							Status:         apicommon.StateProgressing,
							DefinitionName: "eval-def",
						},
					},
				}

				err = k8sClient.Status().Update(ctx, wi2)
				Expect(err).To(BeNil())

				By("Ensuring all phases after pre-eval checks are deprecated")
				wiNameObj := types.NamespacedName{
					Namespace: wi.Namespace,
					Name:      wi.Name,
				}
				//nolint:dupl
				Eventually(func(g Gomega) {
					wi := &klcv1alpha3.KeptnWorkloadVersion{}
					err := k8sClient.Get(ctx, wiNameObj, wi)
					g.Expect(err).To(BeNil())
					g.Expect(wi).To(Not(BeNil()))
					g.Expect(wi.Status.PreDeploymentStatus).To(BeEquivalentTo(apicommon.StateSucceeded))
					g.Expect(wi.Status.PreDeploymentEvaluationStatus).To(BeEquivalentTo(apicommon.StateFailed))
					g.Expect(wi.Status.DeploymentStatus).To(BeEquivalentTo(apicommon.StateDeprecated))
					g.Expect(wi.Status.PostDeploymentStatus).To(BeEquivalentTo(apicommon.StateDeprecated))
					g.Expect(wi.Status.PostDeploymentEvaluationStatus).To(BeEquivalentTo(apicommon.StateDeprecated))
					g.Expect(wi.Status.Status).To(BeEquivalentTo(apicommon.StateFailed))
				}, "30s").Should(Succeed())

				By("Ensuring that a K8s Event containing the reason for the failed evaluation has been sent")

				Eventually(func(g Gomega) {
					eventList := &corev1.EventList{}
					err := k8sClient.List(ctx, eventList, client.InNamespace(namespace))
					g.Expect(err).To(BeNil())

					foundEvent := &corev1.Event{}

					for _, e := range eventList.Items {
						if strings.Contains(e.Name, wi.GetName()) && e.Reason == "AppPreDeployEvaluationsFailed" {
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
				common.LogErrorIfPresent(err)
				err = k8sClient.Delete(ctx, wi)
				common.LogErrorIfPresent(err)
				// Reset span recorder after each spec
				common.ResetSpanRecords(tracer, spanRecorder)
			})

		})

	})
})

func createAppVersionInCluster(name string, namespace string, version string) *klcv1alpha3.KeptnAppVersion {
	instance := &klcv1alpha3.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:       name,
			Namespace:  namespace,
			Generation: 1,
		},
		Spec: klcv1alpha3.KeptnAppVersionSpec{
			AppName: name,
			KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
				Version: version,
				Workloads: []klcv1alpha3.KeptnWorkloadRef{
					{
						Name:    "wname",
						Version: "2.0",
					},
				},
			},
		},
	}
	By("Invoking Reconciling for Create")

	Expect(common.IgnoreAlreadyExists(k8sClient.Create(ctx, instance))).Should(Succeed())

	av := &klcv1alpha3.KeptnAppVersion{}
	err := k8sClient.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, av)
	Expect(err).To(BeNil())

	av.Status.PreDeploymentEvaluationStatus = apicommon.StateSucceeded
	_ = k8sClient.Status().Update(ctx, av)
	return av
}

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
