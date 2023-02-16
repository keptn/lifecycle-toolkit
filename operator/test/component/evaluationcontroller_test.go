package component

import (
	"context"
	"os"
	"time"

	metricsv1alpha2 "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/providers"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnevaluation"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const KLTnamespace = "keptnlifecycle"

var _ = Describe("KeptnEvaluationController", Ordered, func() {
	var (
		evaluationName           string
		evaluationDefinitionName string
		metricName               string
		namespaceName            string
		spanRecorder             *sdktest.SpanRecorder
		tracer                   *otelsdk.TracerProvider
		ns                       *v1.Namespace
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
		controllers := []interfaces.Controller{&keptnevaluation.KeptnEvaluationReconciler{
			Client:        k8sManager.GetClient(),
			Scheme:        k8sManager.GetScheme(),
			Recorder:      k8sManager.GetEventRecorderFor("test-evaluation-controller"),
			Log:           GinkgoLogr,
			Meters:        initKeptnMeters(),
			TracerFactory: &tracerFactory{tracer: tracer},
			Namespace:     KLTnamespace,
		}}
		setupManager(controllers) // we can register multiple time the same controller
		ns = makeKLTDefaultNamespace(KLTnamespace)
	})

	BeforeEach(func() { // list var here they will be copied for every spec
		evaluationName = "test-evaluation"
		evaluationDefinitionName = "my-evaldefinition"
		metricName = "metric1"
		namespaceName = "default" // namespaces are not deleted in the api so be careful
	})

	Describe("Testing reconcile Evaluation scenario when using KeptnMetric instead of provider directly", func() {
		var (
			evaluationDefinition *klcv1alpha2.KeptnEvaluationDefinition
			evaluation           *klcv1alpha2.KeptnEvaluation
		)
		Context("With an existing EvaluationDefinition pointing to KeptnMetric", func() {
			// It("KeptnEvaluationController Should succeed, as it finds valid values in KeptnMetric", func() {
			// 	By("Create EvaluationDefiniton")

			// 	evaluationDefinition = makeEvaluationDefinition(evaluationDefinitionName, namespaceName, metricName, providers.KeptnMetricProviderName)

			// 	By("Create KeptnMetric")

			// 	metric := makeKeptnMetric(metricName)

			// 	By("Update KeptnMetric to have status")

			// 	metric2 := &metricsv1alpha2.KeptnMetric{}
			// 	err := k8sClient.Get(context.TODO(), types.NamespacedName{
			// 		Namespace: KLTnamespace,
			// 		Name:      metric.Name,
			// 	}, metric2)
			// 	Expect(err).To(BeNil())

			// 	metric2.Status = metricsv1alpha2.KeptnMetricStatus{
			// 		Value:       "5",
			// 		RawValue:    []byte("5"),
			// 		LastUpdated: metav1.NewTime(time.Now().UTC()),
			// 	}

			// 	err = k8sClient.Status().Update(context.TODO(), metric2)
			// 	Expect(err).To(BeNil())

			// 	By("Create evaluation to start the process")

			// 	evaluation = makeEvaluation(evaluationName, namespaceName, evaluationDefinitionName)

			// 	By("Check that the evaluation passed")

			// 	evaluation2 := &klcv1alpha2.KeptnEvaluation{}
			// 	Eventually(func(g Gomega) {
			// 		err := k8sClient.Get(context.TODO(), types.NamespacedName{
			// 			Namespace: namespaceName,
			// 			Name:      evaluation.Name,
			// 		}, evaluation2)
			// 		g.Expect(err).To(BeNil())
			// 		g.Expect(evaluation2.Status.OverallStatus).To(Equal(apicommon.StateSucceeded))
			// 		g.Expect(evaluation2.Status.EvaluationStatus).To(Equal(map[string]klcv1alpha2.EvaluationStatusItem{
			// 			metricName: {
			// 				Value:  "5",
			// 				Status: apicommon.StateSucceeded,
			// 			},
			// 		}))
			// 	}, "30s").Should(Succeed())

			// 	err = k8sClient.Delete(context.TODO(), metric)
			// 	logErrorIfPresent(err)
			// })

			// It("KeptnEvaluationController Metric status does not exist", func() {
			// 	By("Create EvaluationDefiniton")

			// 	evaluationDefinition = makeEvaluationDefinition(evaluationDefinitionName, namespaceName, metricName, providers.KeptnMetricProviderName)

			// 	By("Create KeptnMetric")

			// 	metric := makeKeptnMetric(metricName)

			// 	By("Create evaluation to start the process")

			// 	evaluation = makeEvaluation(evaluationName, namespaceName, evaluationDefinitionName)

			// 	By("Check that the evaluation failed")

			// 	evaluation2 := &klcv1alpha2.KeptnEvaluation{}
			// 	Eventually(func(g Gomega) {
			// 		err := k8sClient.Get(context.TODO(), types.NamespacedName{
			// 			Namespace: namespaceName,
			// 			Name:      evaluation.Name,
			// 		}, evaluation2)
			// 		g.Expect(err).To(BeNil())
			// 		g.Expect(evaluation2.Status.OverallStatus).To(Equal(apicommon.StateFailed))
			// 		g.Expect(evaluation2.Status.EvaluationStatus).To(Equal(map[string]klcv1alpha2.EvaluationStatusItem{
			// 			metricName: {
			// 				Value:   "",
			// 				Status:  apicommon.StateFailed,
			// 				Message: "no values",
			// 			},
			// 		}))
			// 	}, "30s").Should(Succeed())

			// 	err := k8sClient.Delete(context.TODO(), metric)
			// 	logErrorIfPresent(err)
			// })
			It("KeptnEvaluationController Metric does not exist", func() {
				By("Create EvaluationDefiniton")

				evaluationDefinition = makeEvaluationDefinition(evaluationDefinitionName, namespaceName, metricName, providers.KeptnMetricProviderName)

				By("Create evaluation to start the process")

				evaluation = makeEvaluation(evaluationName, namespaceName, evaluationDefinitionName)

				By("Check that the evaluation failed")

				evaluation2 := &klcv1alpha2.KeptnEvaluation{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespaceName,
						Name:      evaluation.Name,
					}, evaluation2)
					g.Expect(err).To(BeNil())
					g.Expect(evaluation2.Status.OverallStatus).To(Equal(apicommon.StateFailed))
					g.Expect(evaluation2.Status.EvaluationStatus).To(Equal(map[string]klcv1alpha2.EvaluationStatusItem{
						metricName: {
							Value:   "",
							Status:  apicommon.StateFailed,
							Message: "no values",
						},
					}))
				}, "30s").Should(Succeed())
			})
			It("KeptnEvaluationController Invalid provider", func() {
				By("Create EvaluationDefiniton")

				evaluationDefinition = makeEvaluationDefinition(evaluationDefinitionName, namespaceName, "invalid", "invalid")

				By("Create evaluation to start the process")

				evaluation = makeEvaluation(evaluationName, namespaceName, evaluationDefinitionName)

				By("Check that the evaluation failed")

				time.Sleep(15 * time.Second)

				evaluation2 := &klcv1alpha2.KeptnEvaluation{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespaceName,
						Name:      evaluation.Name,
					}, evaluation2)
					g.Expect(err).To(BeNil())
					g.Expect(evaluation2.Status.OverallStatus).To(BeEmpty())

				}, "10s").Should(Succeed())
			})
			AfterEach(func() {
				err := k8sClient.Delete(context.TODO(), evaluationDefinition)
				logErrorIfPresent(err)
				err = k8sClient.Delete(context.TODO(), evaluation)
				logErrorIfPresent(err)
			})
			AfterAll(func() {
				err := k8sClient.Delete(context.TODO(), ns)
				logErrorIfPresent(err)
			})
		})
	})
})

func makeEvaluationDefinition(name string, namespaceName string, objectiveName string, source string) *klcv1alpha2.KeptnEvaluationDefinition {
	evalDef := &klcv1alpha2.KeptnEvaluationDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
		},
		Spec: klcv1alpha2.KeptnEvaluationDefinitionSpec{
			Source: source,
			Objectives: []klcv1alpha2.Objective{
				{
					Name:             objectiveName,
					Query:            "",
					EvaluationTarget: "<10",
				},
			},
		},
	}

	err := k8sClient.Create(context.TODO(), evalDef)
	Expect(err).To(BeNil())

	return evalDef
}

func makeKeptnMetric(name string) *metricsv1alpha2.KeptnMetric {
	metric := &metricsv1alpha2.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: KLTnamespace,
		},
		Spec: metricsv1alpha2.KeptnMetricSpec{
			Provider: metricsv1alpha2.ProviderRef{
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

func makeEvaluation(name string, namespaceName string, evaluationDefinition string) *klcv1alpha2.KeptnEvaluation {
	eval := &klcv1alpha2.KeptnEvaluation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
		},
		Spec: klcv1alpha2.KeptnEvaluationSpec{
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

func makeKLTDefaultNamespace(name string) *v1.Namespace {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	err := k8sClient.Create(context.TODO(), ns)
	Expect(err).To(BeNil())

	return ns
}
