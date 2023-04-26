package workloadinstance_test

import (
	"context"
	"os"
	"testing"
	"time"

	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnworkloadinstance"
	"github.com/keptn/lifecycle-toolkit/operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// nolint:gci
	// +kubebuilder:scaffold:imports
)

func TestWorkloadinstance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Workloadinstance Suite")
}

var (
	k8sManager   ctrl.Manager
	tracer       *otelsdk.TracerProvider
	k8sClient    client.Client
	ctx          context.Context
	spanRecorder *sdktest.SpanRecorder
)

var _ = BeforeSuite(func() {
	ctx, k8sManager, tracer, spanRecorder, k8sClient, _ = common.InitSuite()

	////setup controllers here
	controller := &keptnworkloadinstance.KeptnWorkloadInstanceReconciler{
		Client:        k8sManager.GetClient(),
		Scheme:        k8sManager.GetScheme(),
		Recorder:      k8sManager.GetEventRecorderFor("test-workloadinstance-controller"),
		Log:           GinkgoLogr,
		Meters:        common.InitKeptnMeters(),
		SpanHandler:   &controllercommon.SpanHandler{},
		TracerFactory: &common.TracerFactory{Tracer: tracer},
	}
	Eventually(controller.SetupWithManager(k8sManager)).WithTimeout(30 * time.Second).WithPolling(time.Second).Should(Succeed())

})

var _ = ReportAfterSuite("custom report", func(report Report) {
	f, err := os.Create("report.workloadinstance-operator")
	Expect(err).ToNot(HaveOccurred(), "failed to generate report")
	for _, specReport := range report.SpecReports {
		common.WriteReport(specReport, f)
	}
	f.Close()
})
