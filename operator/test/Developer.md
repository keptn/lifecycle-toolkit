# Integration tests
This test suite can run test verifying multiple Controllers

### Running on envtest cluster

cd to operator folder, run 
```make test```
Make test is the one-stop shop for downloading the binaries, setting up the test environment, and running the tests.

If you would like to run the generated bin for apiserver etcd etc. from your IDE copy them to the default path "/usr/local/kubebuilder/bin"
This way the default test setup will pick them up without specifying any ENVVAR.
For more info on kubebuilder envtest or to set up a real cluster behind the test have a look [here](https://book.kubebuilder.io/reference/envtest.html)

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

As a minimum example a test could be:
```
var _ = Describe("KeptnAppController", func() {
    var ( //setup needed var
    name      string
    )
    BeforeEach(func() { // init them
    name = "test-app"
    })
    AfterEach(ResetSpanRecords) //you must clean up spans each time 
    
        Describe("Creation of AppVersion from a new App", func() {
            var (
                instance   *klcv1alpha1.KeptnApp // declare CRD
            )
            Context("with one App", func() {
                BeforeEach(func() {  
                //create it using the client eg. Expect(k8sClient.Create(ctx, instance)).Should(Succeed())
                    instance = createInstanceInCluster(name, namespace, version, instance)
                })
                AfterEach(func() {
                    // Remember to clean up the cluster after each test
                    deleteAppInCluster(instance)
                })
                It("should update the status of the CR", func() {
                    assertResourceUpdated(instance)
                })
            })
        })

})
```


## Contributing Best Practice

1. Keep in mind to clean up after each test
2. Namespaces do not get cleaned up by kubebuilder testenv so be careful on that
3. Make sure not to mik up gomega patter with other assertion packages
