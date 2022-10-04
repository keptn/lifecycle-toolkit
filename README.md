# Keptn Lifecycle Controller

![build](https://img.shields.io/github/workflow/status/keptn-sandbox/lifecycle-controller/CI)
![goversion](https://img.shields.io/github/go-mod/go-version/keptn-sandbox/lifecycle-controller?filename=operator%2Fgo.mod)
![version](https://img.shields.io/badge/version-pre--alpha-green)
![status](https://img.shields.io/badge/status-not--for--production-red)
[![GitHub Discussions](https://img.shields.io/github/discussions/keptn-sandbox/lifecycle-controller)](https://github.com/keptn-sandbox/lifecycle-controller/discussions)

The purpose of this repository is to demonstrate and experiment with
a prototype of a _**Keptn Lifecycle Controller**_.
The goal of this prototype is to introduce a more “cloud-native” approach for pre- and post-deployment, as well as the concept of application health checks.
It is an experimental project, under the umbrella of the [Keptn Application Lifecycle working group](https://github.com/keptn/wg-app-lifecycle).

## Deploy the latest release

**Known Limitations**
* Kubernetes >=1.24 is needed to deploy the Lifecycle Controller
* The Lifecycle Controller is currently not compatible with [vcluster](https://github.com/loft-sh/vcluster)

**Installation**

The lifecycle controller includes a Mutating Webhook which requires TLS certificates to be mounted as a volume in its pod. The certificate creation
is handled automatically by [cert-manager](https://cert-manager.io). To install **cert-manager**, follow their [installation instructions](https://cert-manager.io/docs/installation/).

When *cert-manager* is installed, you can run

<!---x-release-please-start-version-->

```
kubectl apply -f https://github.com/keptn-sandbox/lifecycle-controller/releases/download/v0.1.2/release.yaml
```

<!---x-release-please-end-->

to install the latest release of the lifecycle controller.

## Goals

The Keptn Lifecycle Controller aims to support Cloud Native teams with:

- Pre-requisite evaluation before deploying workloads and applications
- Finding out when an application (not workload) is ready and working
- Checking the Application Health in a declarative (cloud-native) way
- Standardized way for pre- and post-deployment tasks
- Provide out-of-the-box Observability of the deployment cycle

![](./assets/operator-maturity.jpg)

The Keptn Lifecycle Controller could be seen as a general purpose and declarative [Level 3 operator](https://operatorframework.io/operator-capabilities/) for your Application.
For this reason, the Keptn Lifecycle Controller is agnostic to deployment tools that are used and works with any GitOps solution.

## How to use

The Keptn Lifecycle Controller monitors manifests that have been applied against the Kubernetes API and reacts if it finds a workload with special annotations.
For this, you should annotate your [Workload](https://kubernetes.io/docs/concepts/workloads/) with (at least) the following two annotations:

```yaml
keptn.sh/app: myAwesomeAppName
keptn.sh/workload: myAwesomeWorkload
```

In case you want to run pre- and post-deployment checks, further annotations are necessary:

```yaml
keptn.sh/pre-deployment-tasks: verify-infrastructure-problems
keptn.sh/post-deployment-tasks: slack-notification,performance-test
```

The value of these annotations are Keptn [CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
called [KeptnTaskDefinition](#keptn-task-definition)s. These CRDs contains re-usable "functions" that can
executed before and after the deployment. In this example, before the deployment starts, a check for open problems in your infrastructure
is performed. If everything is fine, the deployment continues and afterward, a slack notification is sent with the result of
the deployment and a pipeline to run performance tests is invoked. Otherwise, the deployment is kept in a pending state until
the infrastructure is capable to accept deployments again.

A more comprehensive example can be found in our [examples folder](./examples/podtatohead-deployment/) where we
use [Podtato-Head](https://github.com/podtato-head/podtato-head) to run some simple pre-deployment checks.

To run the example, use the following commands:

```bash
cd ./examples/podtatohead-deployment/
kubectl apply -f .
```

Afterward, you can monitor the status of the deployment using

```bash
kubectl get keptnworkloadinstance -n podtato-kubectl -w
```

The deployment for a Workload will stay in a `Pending` state until the respective pre-deployment check is completed. Afterward, the deployment will start and when it is `Succeeded`, the post-deployment checks will start.

## Architecture


The Keptn Lifecycle Controller is composed of the following components:

- Keptn Lifecycle Operator
- Keptn Scheduler

The Keptn Lifecycle Operator contains several controllers for Keptn CRDs and a Mutating Webhook.
The Keptn Scheduler ensures that Pods are started only after the pre-deployment checks have finished.

A Kubernetes Manifest, which is annotated with Keptn specific annotations, gets applied to the Kubernetes Cluster.
Afterward, the Keptn Scheduler gets injected (via Mutating Webhook), and Kubernetes Events for Pre-Deployment are sent to the event stream.
The Event Controller watches for events and triggers a Kubernetes Job to fullfil the Pre-Deployment.
After the Pre-Deployment has finished, the Keptn Scheduler schedules the Pod to be deployed.
The KeptnApp and KeptnWorkload Controllers watch for the workload resources to finish and then generate a Post-Deployment Event.
After the Post-Deployment checks, SLOs can be validated using an interface for retrieving SLI data from a provider, e.g, [Prometheus](https://prometheus.io/).
Finally, Keptn Lifecycle Controller exposes Metrics and Traces of the whole Deployment cycle with [OpenTelemetry](https://opentelemetry.io/).

![](./assets/architecture.png)

## How it works

The following sections will provide insights on each component of the Keptn Lifecycle controller in terms of their purpose, responsibility, and communication with other components.
Furthermore, there will be a description on what CRD they monitor and a general overview of their fields.

### Webhook

The mutating webhook works only on resources that have Keptn annotations.
The mutation consists in changing the scheduler used for the deployment with the Keptn Scheduler.
The webhook should be as fast as possible and should not create/change any resource.

When the webhook receives a request for a new pod, it will look for the following annotations:

```
keptn.sh/app
keptn.sh/workload
```

Additionally, it will compute a version string, using a hash function that takes certain properties of the pod as parameters
(e.g. the images of its containers).
Next, it will look for an existing instance of a `Workload CRD` for the given workload name:

- If it finds the `Workload`, it will update its version according to the previously computed version string.
  In addition, it will include a reference to the ReplicaSet UID of the pod (i.e. the Pods owner),
  or the pod itself, if it does not have an owner.
- If it does not find a workload instance, it will create one containing the previously computed version string.
  In addition, it will include a reference to the ReplicaSet UID of the pod (i.e. the Pods owner), or the pod itself, if it does not have an owner.

It will use the following annotations for
the specification of the pre/post deployment checks that should be executed for the `Workload`:

  - `keptn.sh/pre-deployment-tasks: task1,task2`
  - `keptn.sh/post-deployment-tasks: task1,task2`


After either one of those actions has been taken, the webhook will set the scheduler of the pod and allow the pod to be scheduled.


### Scheduler

After the Webhook mutation, the Keptn-Scheduler will handle the annotated resources. The scheduling flow follows the default scheduler behavior,
since it implements a scheduler plugin based on the [scheduling framework]( https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/).
For each pod, at the very end of the scheduling cycle, the plugin verifies whether the pre deployment checks have terminated, by retrieving the current status of the WorkloadInstance. Only if that is successful, the pod is bound to a node.


### Keptn App

tbd

### Keptn Workload

A Workload contains information about which tasks should be performed during the `preDeployment` as well as the `postDeployment`
phase of a deployment. In its state it keeps track of the currently active `Workload Instances`, which are responsible for doing those checks for
a particular instance of a Deployment/StatefulSet/ReplicaSet (e.g. a Deployment of a certain version).

### Keptn Workload Instance

A Workload Instance is responsible for executing the pre- and post deployment checks of a workload. In its state, it keeps track of the current status of all checks, as well as the overall state of
the Pre Deployment phase, which can be used by the scheduler to tell that a pod can be allowed to be placed on a node.
Workload Instances have a reference to the respective Deployment/StatefulSet/ReplicaSet, to check if it has reached the desired state. If it detects that the referenced object has reached
its desired state (e.g. all pods of a deployment are up and running), it will be able to tell that a `PostDeploymentCheck` can be triggered.

### Keptn Task Definition

A `KeptnTaskDefinition` is a CRD used to define tasks that can be run by the Keptn Lifecycle Controller
as part of pre- and post-deployment phases of a deployment.
The task definition is a [Deno](https://deno.land/) script
Please, refer to the [function runtime](./functions-runtime/) folder for more information about the runtime.
In the future, we also intend to support other runtimes, especially running a container image directly.

A task definition can be configured in three different ways:

- inline
- referring to an HTTP script
- referring to another `KeptnTaskDefinition`

An inline task definition looks like the following:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha1
kind: KeptnTaskDefinition
metadata:
  name: deployment-hello
spec:
  function:
    inline:
      code: |
        console.log("Deployment Task has been executed");
```

In the code section, it is possible to define a full-fletched Deno script.
A further example, is available [here](./examples/taskonly-hello-keptn/inline/taskdefinition.yaml).

To runtime can also fetch the script on the fly from a remote webserver. For this, the CRD should look like the following:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha1
kind: KeptnTaskDefinition
metadata:
  name: hello-keptn-http
spec:
  function:
    httpRef:
      url: <url>
```

An example is available [here](./examples/taskonly-hello-keptn/http/taskdefinition.yaml).

Finally, `KeptnTaskDefinition` can build on top of other `KeptnTaskDefinition`s.
This is a common use case where a general function can be re-used in multiple places with different parameters.

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha1
kind: KeptnTaskDefinition
metadata:
  name: slack-notification-dev
spec:
  function:
    functionRef:
      name: slack-notification
    parameters:
      map:
        textMessage: "This is my configuration"
    secureParameters:
      secret: slack-token
```

As you might have noticed, Task Definitions also have the possibility to use input parameters.
The Lifecycle Controller passes the values defined inside the `map` field as a JSON object.
At the moment, multi-level maps are not supported.
The JSON object can be read through the environment variable `DATA` using `Deno.env.get("DATA");`.
K8s secrets can also be passed to the function using the `secureParameters` field.
Here, the `secret` value is the K8s secret name that will be mounted into the runtime and made available to the function via the environment variable `SECURE_DATA`.


### Keptn Task

A Task is responsible for executing the TaskDefinition of a workload.
The execution is done spawning a K8s Job to handle a single Task.
In its state, it keeps track of the current status of the K8s Job created.


## Install a dev build

The [GitHub CLI](https://cli.github.com/) can be used to download the manifests of the latest CI build.

```bash
gh run list --repo keptn-sandbox/lifecycle-controller # find the id of a run
gh run download 3152895000 --repo keptn-sandbox/lifecycle-controller # download the artifacts
kubectl apply -f ./keptn-lifecycle-operator-manifest/release.yaml # install the operator
kubectl apply -f ./scheduler-manifest/release.yaml # install the scheduler
```

Instead, if you want to build and deploy the operator into your cluster directly from the code, you can type:

```bash
DOCKER_REGISTRY=<YOUR_DOCKER_REGISTRY>
DOCKER_TAG=<YOUR_DOCKER_TAG>

cd operator

make docker-build docker-push IMG=${DOCKER_REGISTRY}/${DOCKER_TAG}:latest
make deploy IMG=${DOCKER_REGISTRY}/${DOCKER_TAG}:latest
```


## License

Please find more information in the [LICENSE](LICENSE) file.

## Thanks to all the people who have contributed 💜

<a href="https://github.com/keptn-sandbox/lifecycle-controller/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=keptn-sandbox/lifecycle-controller" />
</a>

Made with [contrib.rocks](https://contrib.rocks).
