package task_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/taskdefinition"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/keptntask"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
	k8sClient  client.Client
	ctx        context.Context
)

const KeptnNamespace = "keptnlifecycle"

var _ = BeforeSuite(func() {
	var readyToStart chan struct{}
	ctx, k8sManager, _, _, k8sClient, readyToStart = common.InitSuite()

	_ = os.Setenv(taskdefinition.FunctionRuntimeImageKey, "my-image-js")
	_ = os.Setenv(taskdefinition.PythonRuntimeImageKey, "my-image-py")

	config.Instance().SetDefaultNamespace(KeptnNamespace)

	// //setup controllers here
	controller := &keptntask.KeptnTaskReconciler{
		Client:      k8sManager.GetClient(),
		Scheme:      k8sManager.GetScheme(),
		EventSender: eventsender.NewK8sSender(k8sManager.GetEventRecorderFor("test-task-controller")),
		Log:         GinkgoLogr,
		Meters:      common.InitKeptnMeters(),
	}
	Eventually(controller.SetupWithManager(k8sManager)).WithTimeout(30 * time.Second).WithPolling(time.Second).Should(Succeed())
	close(readyToStart)
})

var _ = ReportAfterSuite("custom report", func(report Report) {
	f, err := os.Create("report.task-lifecycle-operator")
	Expect(err).ToNot(HaveOccurred(), "failed to generate report")
	for _, specReport := range report.SpecReports {
		common.WriteReport(specReport, f)
	}
	f.Close()
})
