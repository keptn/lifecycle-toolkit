package component

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnapp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Metric struct {
	creationTime             []time.Time
	succeededAppVersionCount int
}

const LOAD = 100

var _ = Describe("[Feature:Performance] Load KeptnAppController", Ordered, func() {
	var (
		apps         []*klcv1alpha3.KeptnApp //Shelf is declared here
		appVersions  []*klcv1alpha3.KeptnAppVersion
		spanRecorder *sdktest.SpanRecorder
		tracer       *otelsdk.TracerProvider
		metrics      Metric
	)
	BeforeAll(func() {
		//setup once
		By("Waiting for Manager")
		Eventually(func() bool {
			return k8sManager != nil
		}).Should(Equal(true))

		spanRecorder = sdktest.NewSpanRecorder()
		tracer = otelsdk.NewTracerProvider(otelsdk.WithSpanProcessor(spanRecorder))

		controllers := []interfaces.Controller{&keptnapp.KeptnAppReconciler{
			Client:        k8sManager.GetClient(),
			Scheme:        k8sManager.GetScheme(),
			Recorder:      k8sManager.GetEventRecorderFor("load-app-controller"),
			Log:           GinkgoLogr,
			TracerFactory: &tracerFactory{tracer: tracer},
		}}
		setupManager(controllers)
	})

	BeforeEach(func() {
		//		createTimes := make(map[string]metav1.Time, 0)

		for i := 0; i < LOAD; i++ {
			instance := &klcv1alpha3.KeptnApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("app-%d", i),
					Namespace: "default",
				},
				Spec: klcv1alpha3.KeptnAppSpec{
					Version: "1.2.3",
					Workloads: []klcv1alpha3.KeptnWorkloadRef{
						{
							Name:    "app-wname",
							Version: "2.0",
						},
					},
				},
			}
			apps = append(apps, instance)
			Expect(k8sClient.Create(ctx, instance)).Should(Succeed())
			metrics.creationTime = append(metrics.creationTime, time.Now())
		}
	})

	AfterAll(func() {
		generateMetricReport(metrics)
	})
	AfterEach(func() {
		for _, app := range apps {
			// Remember to clean up the cluster after each test
			deleteAppInCluster(app)
			resetSpanRecords(tracer, spanRecorder)
		}
	})
	JustAfterEach(func() { // this is an example of how to add logs to report
		if CurrentSpecReport().Failed() {
			AddReportEntry("current spans", spanRecorder.Ended())
		}
	})

	It("should create the app version CR", func() {
		for _, app := range apps {
			appVersions = append(appVersions, assertResourceUpdated(app))
			metrics.succeededAppVersionCount++
		}
	})
})

func generateMetricReport(metric Metric) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	filePath := path.Join(dir, "load-report."+time.Now().Format(time.RFC3339)+".json")
	report := []byte(fmt.Sprintf("Overall AppVersions created %d/%d \n Creation times: %+v ", metric.succeededAppVersionCount, LOAD, metric.creationTime))
	if err := os.WriteFile(filePath, report, 0644); err != nil {
		GinkgoLogr.Error(err, "error writing to ", filePath)
	}

}
