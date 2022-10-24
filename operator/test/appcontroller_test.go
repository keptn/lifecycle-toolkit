package test

import (
	klcv1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-controller/operator/api/v1alpha1/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// integration test it uses real api
var _ = Describe("Keptn APP controller", func() {
	It("should reconcile", func() {
		appname := "app-name"
		appversion := "1.0.0"
		app := &klcv1alpha1.KeptnApp{
			ObjectMeta: metav1.ObjectMeta{
				Name:      appname,
				Namespace: "default",
			},
			Spec: klcv1alpha1.KeptnAppSpec{
				Version:                   appversion,
				PreDeploymentTasks:        []string{},
				PostDeploymentTasks:       []string{},
				PreDeploymentEvaluations:  []string{},
				PostDeploymentEvaluations: []string{},
				Workloads: []klcv1alpha1.KeptnWorkloadRef{
					{
						Name:    "app-wname",
						Version: "2.0",
					},
				},
			},
		}
		By("Invoking Reconciling for Create")

		Expect(k8sClient.Create(ctx, app)).Should(Succeed())

		appVersion := &klcv1alpha1.KeptnAppVersion{}
		appvName := types.NamespacedName{
			Namespace: "default",
			Name:      appname + "-" + appversion,
		}
		By("Retrieving Created app version")
		Eventually(func() error {
			return k8sClient.Get(ctx, appvName, appVersion)
		}).Should(Succeed())

		By("Comparing expected app version")
		Expect(appVersion.Spec.AppName).To(Equal(appname))
		Expect(appVersion.Spec.Version).To(Equal(appversion))
		Expect(appVersion.Spec.Workloads[0]).To(Equal(klcv1alpha1.KeptnWorkloadRef{Name: "app-wname", Version: "2.0"}))

		By("Comparing spans")
		spans := spanRecorder.Ended()
		Expect(len(spans)).To(Equal(2))

		Expect(spans[0].Name()).To(Equal("create_app_version"))
		Expect(spans[0].Attributes()).To(ContainElement(common.AppName.String(appname)))
		Expect(spans[0].Attributes()).To(ContainElement(common.AppVersion.String(appversion)))

		Expect(spans[1].Name()).To(Equal("reconcile_app"))
		Expect(spans[1].Attributes()).To(ContainElement(common.AppName.String(appname)))
		Expect(spans[1].Attributes()).To(ContainElement(common.AppVersion.String(appversion)))

		GinkgoWriter.Printf("The attributes are %v", spans[1].Attributes())
	})
})
