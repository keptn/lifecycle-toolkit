package test

import (
	klcv1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-controller/operator/api/v1alpha1/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// clean example of E2E test/ integration test --
// App controller creates AppVersion when a new App CRD is added
// span for creation and reconcile are correct
var _ = Describe("KeptnAppController", func() {
	var (
		name      string
		namespace string
		version   string
	)
	BeforeEach(func() { // list var here
		name = "test-app"
		namespace = "default" // namespaces are not deleted in the api so be careful when creating new ones
		version = "1.0.0"
	})
	AfterEach(ResetSpanRecords) //you must clean up spans each time

	Describe("Creation of AppVersion from a new App", func() {
		var (
			instance   *klcv1alpha1.KeptnApp
			appVersion *klcv1alpha1.KeptnAppVersion
		)
		Context("with one App", func() {
			BeforeEach(func() {
				instance = createInstanceInCluster(name, namespace, version, instance)
			})
			AfterEach(func() {
				// Remember to clean up the cluster after each test
				deleteAppInCluster(instance)
				deleteAppVersionInCluster(appVersion)
			})
			It("should update the status of the CR", func() {
				appVersion = assertResourceUpdated(instance)
			})
		})
	})

})

func deleteAppVersionInCluster(version *klcv1alpha1.KeptnAppVersion) {
	By("Cleaning Up Keptn AppVersion CRD")
	Expect(k8sClient.Delete(ctx, version)).Should(Succeed())
}

func deleteAppInCluster(instance *klcv1alpha1.KeptnApp) {
	By("Cleaning Up KeptnApp CRD ")
	Expect(k8sClient.Delete(ctx, instance)).Should(Succeed())

}

func assertResourceUpdated(instance *klcv1alpha1.KeptnApp) *klcv1alpha1.KeptnAppVersion {

	appVersion := &klcv1alpha1.KeptnAppVersion{}
	appvName := types.NamespacedName{
		Namespace: instance.Namespace,
		Name:      instance.Name + "-" + instance.Spec.Version,
	}
	By("Retrieving Created app version")
	Eventually(func() error {
		return k8sClient.Get(ctx, appvName, appVersion)
	}).Should(Succeed())

	By("Comparing expected app version")
	Expect(appVersion.Spec.AppName).To(Equal(instance.Name))
	Expect(appVersion.Spec.Version).To(Equal(instance.Spec.Version))
	Expect(appVersion.Spec.Workloads).To(Equal(instance.Spec.Workloads))

	return appVersion
}

func assertAppSpan(instance *klcv1alpha1.KeptnApp) {
	By("Comparing spans")
	spans := spanRecorder.Ended()
	Expect(len(spans)).To(Equal(2)) //this works only if we do not run tests in parallel

	Expect(spans[0].Name()).To(Equal("create_app_version"))
	Expect(spans[0].Attributes()).To(ContainElement(common.AppName.String(instance.Name)))
	Expect(spans[0].Attributes()).To(ContainElement(common.AppVersion.String(instance.Spec.Version)))

	Expect(spans[1].Name()).To(Equal("reconcile_app"))
	Expect(spans[1].Attributes()).To(ContainElement(common.AppName.String(instance.Name)))
	Expect(spans[1].Attributes()).To(ContainElement(common.AppVersion.String(instance.Spec.Version)))
}

func createInstanceInCluster(name string, namespace string, version string, instance *klcv1alpha1.KeptnApp) *klcv1alpha1.KeptnApp {
	instance = &klcv1alpha1.KeptnApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: klcv1alpha1.KeptnAppSpec{
			Version: version,
			Workloads: []klcv1alpha1.KeptnWorkloadRef{
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
