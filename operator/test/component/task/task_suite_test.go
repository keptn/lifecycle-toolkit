package task_test

import (
	"context"
	"os"
	"testing"
	"time"

	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptntask"
	"github.com/keptn/lifecycle-toolkit/operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// nolint:gci
	// +kubebuilder:scaffold:imports
)

func TestTask(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Task Suite")
}

var (
	k8sManager ctrl.Manager
	tracer     *otelsdk.TracerProvider
	k8sClient  client.Client
	ctx        context.Context
)

var _ = BeforeSuite(func() {
	var readyToStart chan struct{}
	ctx, k8sManager, tracer, _, k8sClient, readyToStart = common.InitSuite()

	_ = os.Setenv(controllercommon.FunctionRuntimeImageKey, "my-image-js")
	_ = os.Setenv(controllercommon.PythonRuntimeImageKey, "my-image-py")

	////setup controllers here
	controller := &keptntask.KeptnTaskReconciler{
		Client:        k8sManager.GetClient(),
		Scheme:        k8sManager.GetScheme(),
		EventSender:   controllercommon.NewEventSender(k8sManager.GetEventRecorderFor("test-task-controller")),
		Log:           GinkgoLogr,
		Meters:        common.InitKeptnMeters(),
		TracerFactory: &common.TracerFactory{Tracer: tracer},
	}
	Eventually(controller.SetupWithManager(k8sManager)).WithTimeout(30 * time.Second).WithPolling(time.Second).Should(Succeed())
	close(readyToStart)
})

var _ = ReportAfterSuite("custom report", func(report Report) {
	f, err := os.Create("report.task-operator")
	Expect(err).ToNot(HaveOccurred(), "failed to generate report")
	for _, specReport := range report.SpecReports {
		common.WriteReport(specReport, f)
	}
	f.Close()
})
