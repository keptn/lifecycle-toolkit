package taskdefinition_test

import (
	"os"
	"testing"
	"time"

	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/keptntaskdefinition"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// nolint:gci
	// +kubebuilder:scaffold:imports
)

func TestTaskdefinition(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Taskdefinition Suite")
}

var (
	k8sManager ctrl.Manager
	k8sClient  client.Client
)

var _ = BeforeSuite(func() {
	var readyToStart chan struct{}
	_, k8sManager, _, _, k8sClient, readyToStart = common.InitSuite()

	// //setup controllers here
	controller := &keptntaskdefinition.KeptnTaskDefinitionReconciler{
		Client:      k8sManager.GetClient(),
		Scheme:      k8sManager.GetScheme(),
		EventSender: controllercommon.NewK8sSender(k8sManager.GetEventRecorderFor("test-taskdefinition-controller")),
		Log:         GinkgoLogr,
	}
	Eventually(controller.SetupWithManager(k8sManager)).WithTimeout(30 * time.Second).WithPolling(time.Second).Should(Succeed())
	close(readyToStart)
})

var _ = ReportAfterSuite("custom report", func(report Report) {
	f, err := os.Create("report.taskdefinition-lifecycle-operator")
	Expect(err).ToNot(HaveOccurred(), "failed to generate report")
	for _, specReport := range report.SpecReports {
		common.WriteReport(specReport, f)
	}
	f.Close()
})
