package evaluation_test

import (
	"context"
	"fmt"
	"time"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	metricsapi "github.com/keptn/lifecycle-toolkit/lifecycle-operator/test/api/metrics/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/storage/names"
)

var _ = Describe("Evaluation", Ordered, func() {
	var (
		evaluationName           string
		evaluationDefinitionName string
		metricName               string
		namespaceName            string
		ns                       *v1.Namespace
	)

	BeforeEach(func() { // list var here they will be copied for every spec
		evaluationName = names.SimpleNameGenerator.GenerateName("test-evaluation-")
		evaluationDefinitionName = names.SimpleNameGenerator.GenerateName("my-evaldefinition-")
		metricName = names.SimpleNameGenerator.GenerateName("metric1-")
		namespaceName = "default" // namespaces are not deleted in the api so be careful
	})

	Describe("Testing reconcile Evaluation scenario when using KeptnMetric instead of provider directly", func() {
		var (
			evaluationDefinition *klcv1alpha3.KeptnEvaluationDefinition
			evaluation           *klcv1alpha3.KeptnEvaluation
		)
		Context("With an existing EvaluationDefinition pointing to KeptnMetric", func() {
			It("KeptnEvaluationController Should succeed, as it finds valid values in KeptnMetric", func() {
				By("Create EvaluationDefiniton")

				evaluationDefinition = makeEvaluationDefinition(evaluationDefinitionName, namespaceName, metricName)

				By("Create KeptnMetric")

				metric := makeKeptnMetric(metricName, namespaceName)

				By("Update KeptnMetric to have status")

				metric2 := &metricsapi.KeptnMetric{}
				err := k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: namespaceName,
					Name:      metric.Name,
				}, metric2)
				Expect(err).To(BeNil())

				metric2.Status = metricsapi.KeptnMetricStatus{
					Value:       "5",
					RawValue:    []byte("5"),
					LastUpdated: metav1.NewTime(time.Now().UTC()),
				}

				err = k8sClient.Status().Update(context.TODO(), metric2)
				Expect(err).To(BeNil())

				evaluationdef := &klcv1alpha3.KeptnEvaluationDefinition{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespaceName,
						Name:      evaluationDefinitionName,
					}, evaluationdef)
					g.Expect(err).To(BeNil())
					g.Expect(evaluationdef.Spec.Objectives[0]).To(Equal(klcv1alpha3.Objective{
						KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
							Name:      metricName,
							Namespace: namespaceName,
						},
						EvaluationTarget: "<10",
					}))
				}, "5s").Should(Succeed())

				By("Create evaluation to start the process")

				evaluation = makeEvaluation(evaluationName, namespaceName, evaluationDefinitionName)

				By("Check that the evaluation passed")

				evaluation2 := &klcv1alpha3.KeptnEvaluation{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespaceName,
						Name:      evaluation.Name,
					}, evaluation2)
					g.Expect(err).To(BeNil())
					g.Expect(evaluation2.Status.OverallStatus).To(Equal(apicommon.StateSucceeded))
					g.Expect(evaluation2.Status.EvaluationStatus).To(Equal(map[string]klcv1alpha3.EvaluationStatusItem{
						metricName: {
							Value:   "5",
							Status:  apicommon.StateSucceeded,
							Message: "value '5' met objective '<10'",
						},
					}))
				}, "30s").Should(Succeed())

				err = k8sClient.Delete(context.TODO(), metric)
				common.LogErrorIfPresent(err)
			})
			It("KeptnEvaluationController Should succeed, as it finds KeptnEvaluationDefinition in default KLT namespace", func() {
				By("Create EvaluationDefiniton")

				evaluationDefinition = makeEvaluationDefinition(evaluationDefinitionName, KLTnamespace, metricName)

				By("Create KeptnMetric")

				metric := makeKeptnMetric(metricName, KLTnamespace)

				By("Update KeptnMetric to have status")

				metric2 := &metricsapi.KeptnMetric{}
				err := k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: KLTnamespace,
					Name:      metric.Name,
				}, metric2)
				Expect(err).To(BeNil())

				metric2.Status = metricsapi.KeptnMetricStatus{
					Value:       "5",
					RawValue:    []byte("5"),
					LastUpdated: metav1.NewTime(time.Now().UTC()),
				}

				err = k8sClient.Status().Update(context.TODO(), metric2)
				Expect(err).To(BeNil())

				evaluationdef := &klcv1alpha3.KeptnEvaluationDefinition{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: KLTnamespace,
						Name:      evaluationDefinitionName,
					}, evaluationdef)
					g.Expect(err).To(BeNil())
					g.Expect(evaluationdef.Spec.Objectives[0]).To(Equal(klcv1alpha3.Objective{
						KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
							Name:      metricName,
							Namespace: KLTnamespace,
						},
						EvaluationTarget: "<10",
					}))
				}, "5s").Should(Succeed())

				By("Create evaluation to start the process")

				evaluation = makeEvaluation(evaluationName, namespaceName, evaluationDefinitionName)

				By("Check that the evaluation passed")

				evaluation2 := &klcv1alpha3.KeptnEvaluation{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespaceName,
						Name:      evaluation.Name,
					}, evaluation2)
					g.Expect(err).To(BeNil())
					g.Expect(evaluation2.Status.OverallStatus).To(Equal(apicommon.StateSucceeded))
					g.Expect(evaluation2.Status.EvaluationStatus).To(Equal(map[string]klcv1alpha3.EvaluationStatusItem{
						metricName: {
							Value:   "5",
							Status:  apicommon.StateSucceeded,
							Message: "value '5' met objective '<10'",
						},
					}))
				}, "30s").Should(Succeed())

				err = k8sClient.Delete(context.TODO(), metric)
				common.LogErrorIfPresent(err)

				err = k8sClient.Delete(context.TODO(), ns)
				common.LogErrorIfPresent(err)
			})
			It("KeptnEvaluationController Metric status does not exist", func() {
				By("Create EvaluationDefiniton")

				evaluationDefinition = makeEvaluationDefinition(evaluationDefinitionName, namespaceName, metricName)

				By("Create KeptnMetric")

				metric := makeKeptnMetric(metricName, namespaceName)

				By("Create evaluation to start the process")

				evaluation = makeEvaluation(evaluationName, namespaceName, evaluationDefinitionName)

				By("Check that the evaluation failed")

				evaluation2 := &klcv1alpha3.KeptnEvaluation{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespaceName,
						Name:      evaluation.Name,
					}, evaluation2)
					g.Expect(err).To(BeNil())
					g.Expect(evaluation2.Status.OverallStatus).To(Equal(apicommon.StateFailed))
					g.Expect(evaluation2.Status.EvaluationStatus).To(Equal(map[string]klcv1alpha3.EvaluationStatusItem{
						metricName: {
							Value:   "",
							Status:  apicommon.StateFailed,
							Message: fmt.Sprintf("empty value for: %s", metric.Name),
						},
					}))
				}, "30s").Should(Succeed())

				err := k8sClient.Delete(context.TODO(), metric)
				common.LogErrorIfPresent(err)
			})
			It("KeptnEvaluationController Metric does not exist", func() {
				By("Create EvaluationDefiniton")

				evaluationDefinition = makeEvaluationDefinition(evaluationDefinitionName, namespaceName, metricName)

				By("Create evaluation to start the process")

				evaluation = makeEvaluation(evaluationName, namespaceName, evaluationDefinitionName)

				By("Check that the evaluation failed")

				evaluation2 := &klcv1alpha3.KeptnEvaluation{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespaceName,
						Name:      evaluation.Name,
					}, evaluation2)
					g.Expect(err).To(BeNil())
					g.Expect(evaluation2.Status.OverallStatus).To(Equal(apicommon.StateFailed))
					g.Expect(evaluation2.Status.EvaluationStatus).To(Equal(map[string]klcv1alpha3.EvaluationStatusItem{
						metricName: {
							Value:   "",
							Status:  apicommon.StateFailed,
							Message: fmt.Sprintf("keptnmetrics.metrics.keptn.sh \"%s\" not found", evaluationDefinition.Spec.Objectives[0].KeptnMetricRef.Name),
						},
					}))
				}, "30s").Should(Succeed())
			})
			AfterEach(func() {
				err := k8sClient.Delete(context.TODO(), evaluationDefinition)
				common.LogErrorIfPresent(err)
				err = k8sClient.Delete(context.TODO(), evaluation)
				common.LogErrorIfPresent(err)
			})
			AfterAll(func() {
				err := k8sClient.Delete(context.TODO(), ns)
				common.LogErrorIfPresent(err)
			})
		})
	})
})

func makeEvaluationDefinition(name string, namespaceName string, objectiveName string) *klcv1alpha3.KeptnEvaluationDefinition {
	evalDef := &klcv1alpha3.KeptnEvaluationDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
		},
		Spec: klcv1alpha3.KeptnEvaluationDefinitionSpec{
			Objectives: []klcv1alpha3.Objective{
				{
					KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
						Name:      objectiveName,
						Namespace: namespaceName,
					},
					EvaluationTarget: "<10",
				},
			},
		},
	}

	err := k8sClient.Create(context.TODO(), evalDef)
	Expect(err).To(BeNil())

	return evalDef
}

func makeKeptnMetric(name string, namespaceName string) *metricsapi.KeptnMetric {
	metric := &metricsapi.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
		},
		Spec: metricsapi.KeptnMetricSpec{
			Provider: metricsapi.ProviderRef{
				Name: "provider",
			},
			Query:                "query",
			FetchIntervalSeconds: 5,
		},
	}

	err := k8sClient.Create(context.TODO(), metric)
	Expect(err).To(BeNil())

	return metric
}

func makeEvaluation(name string, namespaceName string, evaluationDefinition string) *klcv1alpha3.KeptnEvaluation {
	eval := &klcv1alpha3.KeptnEvaluation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
		},
		Spec: klcv1alpha3.KeptnEvaluationSpec{
			AppVersion:           "1",
			AppName:              "app",
			EvaluationDefinition: evaluationDefinition,
			Type:                 apicommon.PreDeploymentEvaluationCheckType,
			Retries:              3,
		},
	}

	err := k8sClient.Create(context.TODO(), eval)
	Expect(err).To(BeNil())

	return eval
}
