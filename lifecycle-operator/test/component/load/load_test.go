package load_test

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Metric struct {
	creationTime             []time.Time
	succeededAppVersionCount int
}

const LOAD = 100

var _ = Describe("Load", Ordered, func() {
	var (
		apps        []*apilifecycle.KeptnApp // Shelf is declared here
		appVersions []*apilifecycle.KeptnAppVersion
		metrics     Metric
	)

	BeforeEach(func() {
		for i := 0; i < LOAD; i++ {
			instance := &apilifecycle.KeptnApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("app-%d", i),
					Namespace: "default",
				},
				Spec: apilifecycle.KeptnAppSpec{
					Version: "1.2.3",
					Workloads: []apilifecycle.KeptnWorkloadRef{
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
			common.DeleteAppInCluster(ctx, k8sClient, app)
		}
	})
	JustAfterEach(func() { // this is an example of how to add logs to report
		if CurrentSpecReport().Failed() {
			AddReportEntry("current spans", spanRecorder.Ended())
		}
	})

	It("should create the app version CR", func() {
		for _, app := range apps {
			appVersions = append(appVersions, common.AssertResourceUpdated(ctx, k8sClient, app))
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
