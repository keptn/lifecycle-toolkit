package app_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnapp"
	"github.com/keptn/lifecycle-toolkit/operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/types"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// nolint:gci
	// +kubebuilder:scaffold:imports
)

func TestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "App Suite")
}

var (
	cfg          *rest.Config
	k8sClient    client.Client
	testEnv      *envtest.Environment
	ctx          context.Context
	cancel       context.CancelFunc
	k8sManager   ctrl.Manager
	spanRecorder *sdktest.SpanRecorder
	tracer       *otelsdk.TracerProvider
)

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	ctx, cancel = context.WithCancel(context.TODO())
	By("bootstrapping test environment")

	if os.Getenv("USE_EXISTING_CLUSTER") == "true" {
		t := true
		testEnv = &envtest.Environment{
			UseExistingCluster: &t,
		}
	} else {
		GinkgoLogr.Info("Setting up fake test env")
		testEnv = &envtest.Environment{
			CRDDirectoryPaths: []string{
				filepath.Join("..", "..", "..", "config", "crd", "bases"),
			},
			ErrorIfCRDPathMissing: true,
		}
	}
	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	// +kubebuilder:scaffold:scheme
	err = klcv1alpha3.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
		gexec.KillAndWait(4 * time.Second)

		// Teardown the test environment once controller is finished.
		// Otherwise, from Kubernetes 1.21+, teardown timeouts waiting on
		// kube-apiserver to return
		err := testEnv.Stop()
		Expect(err).ToNot(HaveOccurred())
	}()

	By("Creating the Controller")

	spanRecorder = sdktest.NewSpanRecorder()
	tracer = otelsdk.NewTracerProvider(otelsdk.WithSpanProcessor(spanRecorder))

	////setup controllers here
	controller := &keptnapp.KeptnAppReconciler{
		Client:        k8sManager.GetClient(),
		Scheme:        k8sManager.GetScheme(),
		Recorder:      k8sManager.GetEventRecorderFor("test-app-controller"),
		Log:           GinkgoLogr,
		TracerFactory: &common.TracerFactory{Tracer: tracer},
	}
	err = controller.SetupWithManager(k8sManager)
	Expect(err).To(BeNil())

})

var _ = ReportAfterSuite("custom report", func(report Report) {
	f, err := os.Create("report.component-operator")
	Expect(err).ToNot(HaveOccurred(), "failed to generate report")
	for _, specReport := range report.SpecReports {
		path := strings.Split(specReport.FileName(), "/")
		testFile := path[len(path)-1]
		if specReport.ContainerHierarchyTexts != nil {
			testFile = specReport.ContainerHierarchyTexts[0]
		}
		fmt.Fprintf(f, "%s %s ", testFile, specReport.LeafNodeText)
		switch specReport.State {
		case types.SpecStatePassed:
			fmt.Fprintf(f, "%s\n", "✓")
		case types.SpecStateFailed:
			fmt.Fprintf(f, "%s\n", "✕")
		default:
			fmt.Fprintf(f, "%s\n", specReport.State)
		}
	}
	f.Close()
})
