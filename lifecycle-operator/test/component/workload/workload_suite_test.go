package workload_test

import (
	"context"
	"os"
	"testing"
	"time"

	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/keptnworkload"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// nolint:gci
	// +kubebuilder:scaffold:imports
)

func TestWorkload(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Workload Suite")
}

var (
	k8sManager   ctrl.Manager
	tracer       *otelsdk.TracerProvider
	k8sClient    client.Client
	ctx          context.Context
	spanRecorder *sdktest.SpanRecorder
)

var _ = BeforeSuite(func() {
	var readyToStart chan struct{}
	ctx, k8sManager, tracer, spanRecorder, k8sClient, readyToStart = common.InitSuite()

	// //setup controllers here
	controller := &keptnworkload.KeptnWorkloadReconciler{
		Client:        k8sManager.GetClient(),
		Scheme:        k8sManager.GetScheme(),
		EventSender:   controllercommon.NewK8sSender(k8sManager.GetEventRecorderFor("test-workload-controller")),
		Log:           GinkgoLogr,
		TracerFactory: &common.TracerFactory{Tracer: tracer},
	}
	Eventually(controller.SetupWithManager(k8sManager)).WithTimeout(30 * time.Second).WithPolling(time.Second).Should(Succeed())
	close(readyToStart)
})

var _ = ReportAfterSuite("custom report", func(report Report) {
	f, err := os.Create("report.workload-lifecycle-operator")
	Expect(err).ToNot(HaveOccurred(), "failed to generate report")
	for _, specReport := range report.SpecReports {
		common.WriteReport(specReport, f)
	}
	f.Close()
})
