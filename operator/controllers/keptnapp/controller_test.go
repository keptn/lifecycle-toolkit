package keptnapp

import (
	klcv1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

//component test using fake api, or integration test if it uses real api
var _ = Describe("Keptn APP controller", func() {
	It("should reconcile", func() {

		app := &klcv1alpha1.KeptnApp{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "app-name",
				Namespace: "default",
			},
			Spec: klcv1alpha1.KeptnAppSpec{
				Version:                   "1.0.0",
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
			Name:      "app-name-1.0.0",
		}
		By("Retrieving Created app version")
		Eventually(func() error {
			return k8sClient.Get(ctx, appvName, appVersion)
		}).Should(Succeed())

		By("Comparing expected app version")
		Expect(appVersion.Spec.AppName).To(Equal("app-name"))
		Expect(appVersion.Spec.Version).To(Equal("1.0.0"))
		Expect(appVersion.Spec.Workloads[0]).To(Equal(klcv1alpha1.KeptnWorkloadRef{Name: "app-wname", Version: "2.0"}))

		By("Comparing spans")
		spans := spanRecorder.Ended()
		Expect(len(spans)).To(Equal(2))

		//span := spans[0]
		// exampleasserts spans
		//assert.Equal(t, "Route 53", span.Name())
		//assert.Equal(t, trace.SpanKindClient, span.SpanKind())
		//assert.Equal(t, c.expectedError, span.Status().Code)
		//attrs := span.Attributes()
		//assert.Contains(t, attrs, attribute.Int("http.status_code", c.expectedStatusCode))
		//if c.expectedRequestID != "" {
		//	assert.Contains(t, attrs, attribute.String("aws.request_id", c.expectedRequestID))
		//}
		//assert.Contains(t, attrs, attribute.String("aws.service", "Route 53"))
		//assert.Contains(t, attrs, attribute.String("aws.region", c.expectedRegion))
		//assert.Contains(t, attrs, attribute.String("aws.operation", "ChangeResourceRecordSets"))

	})
})
