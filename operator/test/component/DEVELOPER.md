# Component tests

This test suite can run test verifying multiple Controllers

## Running on envtest cluster

cd to operator folder, run
```make test```
Make test is the one-stop shop for downloading the binaries, setting up the test environment, and running the tests.

If you would like to run the generated bin for apiserver etcd etc. from your IDE copy them to the default path "
/usr/local/kubebuilder/bin"
This way the default test setup will pick them up without specifying any ENVVAR.
For more info on kubebuilder envtest or to set up a real cluster behind the test have a
look [here](https://book.kubebuilder.io/reference/envtest.html)

After run a ```report.component-operator``` file will be generated with the results of each test:

```text
suite_test.go  | passed
[Feature:Performance] Load KeptnAppController should create the app version CR | passed
KeptnAppController should update the status of the CR  | passed
KeptnAppController should update the spans | failed
```

## Contributing

In each test you can add one or more new controllers to the suite_test similarly as follows:

```go
controllers := []keptncontroller.Controller{&keptnapp.KeptnAppReconciler{

Client:   k8sManager.GetClient(),
Scheme:   k8sManager.GetScheme(),
Recorder: k8sManager.GetEventRecorderFor("test-app-controller"),
Log:      GinkgoLogr,
Tracer:   tracer.Tracer("test-app-tracer"),
}}
setupManager(controllers)
```

After that the k8s API from kubebuilder will handle its CRD
Each Ginkgo test should be structured following the [spec bestpractices](https://onsi.github.io/ginkgo/#writing-specs)

As a minimum example, a test could be:

<!-- markdownlint-disable MD010 -->

```go
package component

var _ = Describe("KeptnAppController", func() {
	var ( //setup needed var
		name string
	)
	BeforeEach(func() { // init them
		name = "test-app"
	})
	AfterEach(ResetSpanRecords) //you must clean up spans each time 

	Describe("Creation of AppVersion from a new App", func() {
		var (
			instance *klcv1alpha3.KeptnApp // declare CRD
		)
		Context("with one App", func() {
			BeforeEach(func() {
				//create it using the client eg. Expect(k8sClient.Create(ctx, instance)).Should(Succeed())
				instance = createInstanceInCluster(name, namespace, version)
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

<!-- markdownlint-enable MD010 -->

## Load Tests

You can append ```[Feature:Performance]``` to any spec you would like to execute during performance test
with ```make performance-test``` the file
"load_test.go" contains examples of such tests, including a simple reporter. The report "MetricForLoadTestSuite" is
generated for every run of the load test.

## Contributing Tips

1. Keep in mind to clean up after each test since the environment is shared. E.g. if you plan assertions on events or
   spans, make sure your specs are either ordered or assigned to their own controller
2. Namespaces do not get cleaned up by EnvTest, so do not make assertion based on the idea that the namespace has been
   deleted, and make sure to use `ignoreAlreadyExists(err error)` when creating a new one
3. EnvTest is a lightweight control plane only meant for testing purposes. This means it does not contain inbuilt
   Kubernetes controllers like deployment controllers, ReplicaSet controllers, etc. You cannot assert/verify for pods
   being created or not for created deployment.
4. You should generally try to use Gomega’s Eventually to make asynchronous assertions, especially in the case of Get
   and Update calls to API Server.
5. Use ginkgo --until-it-fails to identify flaky tests.
6. Avoid general utility packages. Packages called "util" are suspect. Instead, derive a name that describes your
   desired function. For example, the utility functions dealing with waiting for operations are in the wait package and
   include functionality like Poll. The full name is wait.Poll.
7. All filenames should be lowercase.
8. Go source files and directories use underscores, not dashes.
9. Package directories should generally avoid using separators as much as possible. When package names are multiple
   words, they usually should be in nested subdirectories.
10. Document directories and filenames should use dashes rather than underscores.
11. Examples should also illustrate best practices for configuration and using
    the [system](https://kubernetes.io/docs/concepts/configuration/overview/).
