package evaluation_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnevaluation"
	"github.com/keptn/lifecycle-toolkit/operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// nolint:gci
	// +kubebuilder:scaffold:imports
)

func TestEvaluation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Evaluation Suite")
}

var (
	k8sManager   ctrl.Manager
	tracer       *otelsdk.TracerProvider
	k8sClient    client.Client
	ctx          context.Context
	spanRecorder *sdktest.SpanRecorder
	ns           *v1.Namespace
)

const KLTnamespace = "keptnlifecycle"

var _ = BeforeSuite(func() {
	ctx, k8sManager, tracer, spanRecorder, k8sClient, _ = common.InitSuite()

	////setup controllers here
	controller := &keptnevaluation.KeptnEvaluationReconciler{
		Client:        k8sManager.GetClient(),
		Scheme:        k8sManager.GetScheme(),
		Recorder:      k8sManager.GetEventRecorderFor("test-evaluation-controller"),
		Log:           GinkgoLogr,
		Meters:        common.InitKeptnMeters(),
		TracerFactory: &common.TracerFactory{Tracer: tracer},
		Namespace:     KLTnamespace,
	}
	Eventually(controller.SetupWithManager(k8sManager)).WithTimeout(30 * time.Second).WithPolling(time.Second).Should(Succeed())

	ns = common.MakeKLTDefaultNamespace(k8sClient, KLTnamespace)

})

var _ = ReportAfterSuite("custom report", func(report Report) {
	f, err := os.Create("report.evaluation-operator")
	Expect(err).ToNot(HaveOccurred(), "failed to generate report")
	for _, specReport := range report.SpecReports {
		common.WriteReport(specReport, f)
	}
	f.Close()
})
