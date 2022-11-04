# E2E tests
This test suite can run test verifying the scheduler, this rely on a real cluster with installed scheduler

### Running on kind cluster

```
kind create cluster
cd lifecycle-toolkit
make build-deploy-scheduler RELEASE_REGISTRY=yourregistry

```

wait for everything to be up and running, then cd to scheduler folder and run 
```make test```
Make test is the one-stop shop for downloading the binaries, setting up the test environment, and running the tests.

If you would like to run the generated bin for apiserver etcd etc. from your IDE copy them to the default path "/usr/local/kubebuilder/bin"
This way the default test setup will pick them up without specifying any ENVVAR.
For more info on kubebuilder envtest or to set up a real cluster behind the test have a look [here](https://book.kubebuilder.io/reference/envtest.html)

After run a ```report.custom``` file will be generated with the results of each test


## Contributing



## Contributing Tips

1. Keep in mind to clean up after each test since the environment is shared. E.g. if you plan assertions on events or spans, make sure your specs are either ordered or assigned to their own controller
2. Namespaces do not get cleaned up by EnvTest, so do not make assertion based on the idea that the namespace has been deleted, and make sure to use `ignoreAlreadyExists(err error)` when creating a new one
3. EnvTest is a lightweight control plane only meant for testing purposes. This means it does not contain inbuilt Kubernetes controllers like deployment controllers, ReplicaSet controllers, etc. You cannot assert/verify for pods being created or not for created deployment. 
4. You should generally try to use Gomega’s Eventually to make asynchronous assertions, especially in the case of Get and Update calls to API Server.
5. Use ginkgo --until-it-fails to identify flaky tests.
6. Avoid general utility packages. Packages called "util" are suspect. Instead, derive a name that describes your desired function. For example, the utility functions dealing with waiting for operations are in the wait package and include functionality like Poll. The full name is wait.Poll.
7. All filenames should be lowercase. 
8. Go source files and directories use underscores, not dashes.
9. Package directories should generally avoid using separators as much as possible. When package names are multiple words, they usually should be in nested subdirectories. 
10. Document directories and filenames should use dashes rather than underscores. 
11. Examples should also illustrate best practices for configuration and using the [system](https://kubernetes.io/docs/concepts/configuration/overview/).
