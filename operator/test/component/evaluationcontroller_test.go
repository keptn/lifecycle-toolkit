package component

import (
	"context"
	"os"
	"time"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	metricsv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/apis/metrics/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/providers"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnevaluation"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("KeptnEvaluationController", Ordered, func() {
	var (
		evaluationName           string
		evaluationDefinitionName string
		metricName               string
		namespace                string
		spanRecorder             *sdktest.SpanRecorder
		tracer                   *otelsdk.TracerProvider
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
			Client:   k8sManager.GetClient(),
			Scheme:   k8sManager.GetScheme(),
			Recorder: k8sManager.GetEventRecorderFor("test-evaluation-controller"),
			Log:      GinkgoLogr,
			Meters:   initKeptnMeters(),
			Tracer:   tracer.Tracer("test-evaluation-tracer"),
		}}
		setupManager(controllers) // we can register multiple time the same controller
	})

	BeforeEach(func() { // list var here they will be copied for every spec
		evaluationName = "test-evaluation"
		evaluationDefinitionName = "my-evaldefinition"
		metricName = "metric1"
		namespace = "default" // namespaces are not deleted in the api so be careful
	})

	Describe("Testing reconcile Evaluation scenario when using KeptnMetric instead of provider directly", func() {
		var (
			evaluationDefinition *klcv1alpha2.KeptnEvaluationDefinition
			evaluation           *klcv1alpha2.KeptnEvaluation
			metric               *metricsv1alpha1.KeptnMetric
		)
		Context("With an existing EvaluationDefinition pointing to KeptnMetric", func() {
			BeforeEach(func() {
				evaluationDefinition = makeEvaluationDefinition(evaluationDefinitionName, namespace, metricName)
				metric = makeKeptnMetric(metricName, namespace)
			})

			It("Should succeed, as it finds valid values in KeptnMetric", func() {
				By("Update KeptnMetric to have status")

				metric2 := &metricsv1alpha1.KeptnMetric{}
				err := k8sClient.Get(context.TODO(), types.NamespacedName{
					Namespace: namespace,
					Name:      metric.Name,
				}, metric2)
				Expect(err).To(BeNil())

				metric2.Status = metricsv1alpha1.KeptnMetricStatus{
					Value:       "5",
					RawValue:    []byte("5"),
					LastUpdated: metav1.NewTime(time.Now().UTC()),
				}

				err = k8sClient.Status().Update(context.TODO(), metric2)
				Expect(err).To(BeNil())

				By("Create evaluation to start the process")

				evaluation = makeEvaluation(evaluationName, namespace, evaluationDefinitionName)

				By("Check that the evaluation passed")

				evaluation2 := &klcv1alpha2.KeptnEvaluation{}
				Eventually(func(g Gomega) {
					err := k8sClient.Get(context.TODO(), types.NamespacedName{
						Namespace: namespace,
						Name:      evaluation.Name,
					}, evaluation2)
					g.Expect(err).To(BeNil())
					g.Expect(evaluation2.Status.OverallStatus).To(Equal(apicommon.StateSucceeded))
				}, "30s").Should(Succeed())
			})
			AfterEach(func() {
				err := k8sClient.Delete(context.TODO(), evaluationDefinition)
				logErrorIfPresent(err)
				err = k8sClient.Delete(context.TODO(), metric)
				logErrorIfPresent(err)
				err = k8sClient.Delete(context.TODO(), evaluation)
				logErrorIfPresent(err)
			})
		})
	})
})

func makeEvaluationDefinition(name string, namespace string, metricName string) *klcv1alpha2.KeptnEvaluationDefinition {
	evalDef := &klcv1alpha2.KeptnEvaluationDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: klcv1alpha2.KeptnEvaluationDefinitionSpec{
			Source: providers.KeptnMetricProviderName,
			Objectives: []klcv1alpha2.Objective{
				{
					Name:             metricName,
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

func makeKeptnMetric(name string, namespace string) *metricsv1alpha1.KeptnMetric {
	metric := &metricsv1alpha1.KeptnMetric{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: metricsv1alpha1.KeptnMetricSpec{
			Provider: metricsv1alpha1.ProviderRef{
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

func makeEvaluation(name string, namespace string, evaluationDefinition string) *klcv1alpha2.KeptnEvaluation {
	eval := &klcv1alpha2.KeptnEvaluation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: klcv1alpha2.KeptnEvaluationSpec{
			AppVersion:           "1",
			AppName:              "app",
			EvaluationDefinition: evaluationDefinition,
			Type:                 apicommon.PreDeploymentEvaluationCheckType,
		},
	}

	err := k8sClient.Create(context.TODO(), eval)
	Expect(err).To(BeNil())

	return eval
}
