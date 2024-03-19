package workloadversion_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/phase"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/keptnworkloadversion"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestWorkloadVersion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WorkloadVersion Suite")
}

var (
	k8sManager   ctrl.Manager
	tracer       *otelsdk.TracerProvider
	k8sClient    client.Client
	ctx          context.Context
	spanRecorder *sdktest.SpanRecorder
)

const (
	KeptnNamespace     = "keptnlifecycle"
	traceComponentName = "keptn/lifecycle-operator/workloadversion"
)

var _ = BeforeSuite(func() {
	var readyToStart chan struct{}
	ctx, k8sManager, tracer, spanRecorder, k8sClient, readyToStart = common.InitSuite()

	eventSender := eventsender.NewK8sSender(k8sManager.GetEventRecorderFor("test-workloadversion-controller"))

	tracerFactory := &common.TracerFactory{Tracer: tracer}

	phaseHandler := phase.NewHandler(
		k8sManager.GetClient(),
		eventSender,
		GinkgoLogr,
		&telemetry.Handler{},
	)

	// //setup controllers here
	config.Instance().SetDefaultNamespace(KeptnNamespace)
	controller := &keptnworkloadversion.KeptnWorkloadVersionReconciler{
		Client:        k8sManager.GetClient(),
		Scheme:        k8sManager.GetScheme(),
		EventSender:   eventSender,
		Log:           GinkgoLogr,
		Meters:        common.InitKeptnMeters(),
		SpanHandler:   &telemetry.Handler{},
		TracerFactory: tracerFactory,
		PhaseHandler:  phaseHandler,
		Config:        config.Instance(),
	}
	Eventually(controller.SetupWithManager(k8sManager)).WithTimeout(30 * time.Second).WithPolling(time.Second).Should(Succeed())
	close(readyToStart)
})

var _ = ReportAfterSuite("custom report", func(report Report) {
	f, err := os.Create("report.workloadversion-lifecycle-operator")
	Expect(err).ToNot(HaveOccurred(), "failed to generate report")
	for _, specReport := range report.SpecReports {
		common.WriteReport(specReport, f)
	}
	f.Close()
})
