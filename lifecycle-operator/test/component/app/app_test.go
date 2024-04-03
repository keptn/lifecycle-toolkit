package app_test

import (
	"fmt"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
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
			instance *apilifecycle.KeptnApp
		)

		BeforeEach(func() {
			instance = createInstanceInCluster(name, namespace, version)
			fmt.Println("created ", instance.Name)
		})

		Context("with a new App CRD", func() {

			It("should update the spans", func() {
				By("creating a new app version")
				common.AssertResourceUpdated(ctx, k8sClient, instance)
			})

		})
		AfterEach(func() {
			// Remember to clean up the cluster after each test
			common.DeleteAppInCluster(ctx, k8sClient, instance)
		})

	})
})

func createInstanceInCluster(name string, namespace string, version string) *apilifecycle.KeptnApp {
	instance := &apilifecycle.KeptnApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:       name,
			Namespace:  namespace,
			Generation: 1,
		},
		Spec: apilifecycle.KeptnAppSpec{
			Version: version,
			Workloads: []apilifecycle.KeptnWorkloadRef{
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
