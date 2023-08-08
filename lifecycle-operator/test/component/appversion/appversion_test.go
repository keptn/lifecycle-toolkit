package appversion_test

import (
	"context"
	"strings"
	"time"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/storage/names"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Appversion", Ordered, func() {
	var (
		appName   string
		namespace string
		version   string
	)

	BeforeEach(func() { // list var here they avll be copied for every spec
		appName = names.SimpleNameGenerator.GenerateName("test-appversion-reconciler-")
		namespace = "default" // namespaces are not deleted in the api so be careful
		// when creating you can use ignoreAlreadyExists(err error)
		version = "1.0.0"
	})
	Describe("Creation of AppVersion", func() {
		var (
			av *klcv1alpha3.KeptnAppVersion
		)
		Context("reconcile a new AppVersions CRD", func() {

			BeforeEach(func() {
			})

			It("should be deprecated when pre-eval checks failed", func() {
				evaluation := &klcv1alpha3.KeptnEvaluation{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pre-eval-eval-def-appversion",
						Namespace: namespace,
					},
					Spec: klcv1alpha3.KeptnEvaluationSpec{
						EvaluationDefinition: "eval-def-appversion",
						AppName:              appName,
						AppVersion:           version,
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

				av = &klcv1alpha3.KeptnAppVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      appName + "-" + version,
						Namespace: namespace,
					},
					Spec: klcv1alpha3.KeptnAppVersionSpec{
						AppName: appName,
						KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
							Version:                  version,
							PreDeploymentEvaluations: []string{"eval-def-appversion"},
						},
					},
				}
				By("Creating AppVersion")
				err = k8sClient.Create(context.TODO(), av)
				Expect(err).To(BeNil())

				time.Sleep(5 * time.Second)

				av2 := &klcv1alpha3.KeptnAppVersion{}
				err = k8sClient.Get(ctx, types.NamespacedName{Namespace: av.Namespace, Name: av.Name}, av2)
				Expect(err).To(BeNil())
				Expect(av2).To(Not(BeNil()))

				av2.Status = klcv1alpha3.KeptnAppVersionStatus{
					PreDeploymentStatus:            apicommon.StateSucceeded,
					PreDeploymentEvaluationStatus:  apicommon.StateProgressing,
					WorkloadOverallStatus:          apicommon.StatePending,
					PostDeploymentStatus:           apicommon.StatePending,
					PostDeploymentEvaluationStatus: apicommon.StatePending,
					CurrentPhase:                   apicommon.PhaseWorkloadPreEvaluation.ShortName,
					Status:                         apicommon.StateProgressing,
					PreDeploymentEvaluationTaskStatus: []klcv1alpha3.ItemStatus{
						{
							Name:           "pre-eval-eval-def-appversion",
							Status:         apicommon.StateProgressing,
							DefinitionName: "eval-def-appversion",
						},
					},
				}

				err = k8sClient.Status().Update(ctx, av2)
				Expect(err).To(BeNil())

				By("Ensuring all phases after pre-eval checks are deprecated")
				avNameObj := types.NamespacedName{
					Namespace: av.Namespace,
					Name:      av.Name,
				}
				//nolint:dupl
				Eventually(func(g Gomega) {
					av := &klcv1alpha3.KeptnAppVersion{}
					err := k8sClient.Get(ctx, avNameObj, av)
					g.Expect(err).To(BeNil())
					g.Expect(av).To(Not(BeNil()))
					g.Expect(av.Status.PreDeploymentStatus).To(BeEquivalentTo(apicommon.StateSucceeded))
					g.Expect(av.Status.PreDeploymentEvaluationStatus).To(BeEquivalentTo(apicommon.StateFailed))
					g.Expect(av.Status.WorkloadOverallStatus).To(BeEquivalentTo(apicommon.StateDeprecated))
					g.Expect(av.Status.PostDeploymentStatus).To(BeEquivalentTo(apicommon.StateDeprecated))
					g.Expect(av.Status.PostDeploymentEvaluationStatus).To(BeEquivalentTo(apicommon.StateDeprecated))
					g.Expect(av.Status.Status).To(BeEquivalentTo(apicommon.StateFailed))
				}, "30s").Should(Succeed())

				By("Ensuring that a K8s Event containing the reason for the failed evaluation has been sent")

				Eventually(func(g Gomega) {
					eventList := &corev1.EventList{}
					err := k8sClient.List(ctx, eventList, client.InNamespace(namespace))
					g.Expect(err).To(BeNil())

					foundEvent := &corev1.Event{}

					for _, e := range eventList.Items {
						if strings.Contains(e.Name, av.GetName()) && e.Reason == "AppPreDeployEvaluationsFailed" {
							foundEvent = &e
							break
						}
					}
					g.Expect(foundEvent).NotTo(BeNil())
				}, "30s").Should(Succeed())
			})
			AfterEach(func() {
				// Remember to clean up the cluster after each test
				err := k8sClient.Delete(ctx, av)
				common.LogErrorIfPresent(err)
				// Reset span recorder after each spec
				common.ResetSpanRecords(tracer, spanRecorder)
			})

		})

	})
})
