package load_test

import (
	"context"
	"os"
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnapp"
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

func TestLoad(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Load Suite")
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
	controller := &keptnapp.KeptnAppReconciler{
		Client:        k8sManager.GetClient(),
		Scheme:        k8sManager.GetScheme(),
		Recorder:      k8sManager.GetEventRecorderFor("load-app-controller"),
		Log:           GinkgoLogr,
		TracerFactory: &common.TracerFactory{Tracer: tracer},
	}
	err := controller.SetupWithManager(k8sManager)
	Expect(err).To(BeNil())

})

var _ = ReportAfterSuite("custom report", func(report Report) {
	f, err := os.Create("report.load-operator")
	Expect(err).ToNot(HaveOccurred(), "failed to generate report")
	for _, specReport := range report.SpecReports {
		common.WriteReport(specReport, f)
	}
	f.Close()
})
