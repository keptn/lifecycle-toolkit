# Integration/ E2E tests

This test suite can run test verifying the operator

## Running on kind cluster

```shell
kind create cluster
cd lifecycle-toolkit
make build-deploy-operator RELEASE_REGISTRY=yourregistry

```

wait for everything to be up and running, then cd to operator folder and run

```make e2e-test```

If you would like more info on kubebuilder envtest or to set up a real cluster behind the test have a
look [here](https://book.kubebuilder.io/reference/envtest.html)

After the run a ```report.E2E-operator``` file will be generated with the results of each test:

```text
2022-11-04 12:46:05.2373262 +0000 UTC
If annotated for keptn, a new Pod should stay pending | passed
If annotated for keptn, a new Pod should be assigned to keptn scheduler | passed
```

## Contributing

## Load Tests

You can append ```[Feature:Performance]``` to any spec you would like to execute during performance test
with ```make performance-test``` the file
"load_test.go" contains examples of such tests, including a simple reporter. The report "MetricForLoadTestSuite" is
generated for every run of the load test.

## Contributing Tips

1. Keep in mind to clean up after each test since the environment is shared. E.g. if you plan assertions on events or
   spans, make sure your specs are either ordered or assigned to their own controller
2. You should generally try to use Gomega’s Eventually to make asynchronous assertions, especially in the case of Get
   and Update calls to API Server.
3. Use ginkgo --until-it-fails to identify flaky tests.
4. Avoid general utility packages. Packages called "util" are suspect. Instead, derive a name that describes your
   desired function. For example, the utility functions dealing with waiting for operations are in the wait package and
   include functionality like Poll. The full name is wait.Poll.
5. All filenames should be lowercase.
6. Go source files and directories use underscores, not dashes.
7. Package directories should generally avoid using separators as much as possible. When package names are multiple
   words, they usually should be in nested subdirectories.
8. Document directories and filenames should use dashes rather than underscores.
9. Examples should also illustrate best practices for configuration and using
   the [system](https://kubernetes.io/docs/concepts/configuration/overview/).
