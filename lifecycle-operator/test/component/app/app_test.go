package app_test

import (
	"fmt"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/test/component/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/storage/names"
)

var _ = Describe("App", Ordered, func() {
	var (
		name      string
		namespace string
		version   string
	)
	BeforeEach(func() { // list var here they will be copied for every spec
		name = names.SimpleNameGenerator.GenerateName("my-app-")
		namespace = "default" // namespaces are not deleted in the api so be careful
		// when creating you can use ignoreAlreadyExists(err error)
		version = "1.0.0"
	})
	Describe("Creation of AppVersion from a new App", func() {
		var (
			instance *klcv1alpha3.KeptnApp
		)

		BeforeEach(func() {
			instance = createInstanceInCluster(name, namespace, version)
			fmt.Println("created ", instance.Name)
		})

		Context("with a new App CRD", func() {

			It("should update the spans", func() {
				By("creating a new app version")
				common.AssertResourceUpdated(ctx, k8sClient, instance)
				fmt.Println("spanned ", instance.Name)
			})

		})
		AfterEach(func() {
			// Remember to clean up the cluster after each test
			common.DeleteAppInCluster(ctx, k8sClient, instance)
			// Reset span recorder after each spec
			common.ResetSpanRecords(tracer, spanRecorder)
		})

	})
})

func createInstanceInCluster(name string, namespace string, version string) *klcv1alpha3.KeptnApp {
	instance := &klcv1alpha3.KeptnApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:       name,
			Namespace:  namespace,
			Generation: 1,
		},
		Spec: klcv1alpha3.KeptnAppSpec{
			Version: version,
			Workloads: []klcv1alpha3.KeptnWorkloadRef{
				{
					Name:    "app-wname",
					Version: "2.0",
				},
			},
		},
	}
	By("Invoking Reconciling for Create")

	Expect(k8sClient.Create(ctx, instance)).Should(Succeed())
	return instance
}
