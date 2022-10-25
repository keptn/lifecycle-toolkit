# Integration tests
This test suite can run test verifying multiple Controllers

### Running on your cluster
1. 


## Contributing

Each new controller should be added to the suite_test similarly as follows:

	```err = (&keptnapp.KeptnAppReconciler{
		Client:   k8sManager.GetClient(),
		Scheme:   k8sManager.GetScheme(),
		Recorder: k8sManager.GetEventRecorderFor("test-app-controller"),
		Log:      GinkgoLogr,
		Tracer:   tr.Tracer("test-app-tracer"),
	}).SetupWithManager(k8sManager)```
	
After that the k8s API from kubebuilder will handle its CRD 

Each Ginkgo test should be structured following the [spec bestpractices](https://onsi.github.io/ginkgo/#writing-specs)

As a minimum example a test could be

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
			instance *klcv1alpha1.KeptnApp
		)
		Context("with one App", func() {
			BeforeEach(func() {
				instance = createInstanceInCluster(name, namespace, version, instance)
			})
			AfterEach(func() {
				deleteAppInCluster(instance)
			})
			It("should update the status of the CR", func() {
				assertResourceUpdated(instance)	
			})
		})
	})

})



### How it works


