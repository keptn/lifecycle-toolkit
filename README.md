# Keptn Lifecycle Controller

[![experimental](http://badges.github.io/stability-badges/dist/experimental.svg)](http://github.com/badges/stability-badges)
[![GitHub Discussions](https://img.shields.io/github/discussions/keptn-sandbox/lifecycle-controller)](https://github.com/keptn-sandbox/lifecycle-controller/discussions)

The purpose of this repository is to demonstrate and experiment with
a prototype of a _**Keptn Lifecycle Controller**_.
The goal of this prototype is to introduce a more “cloud-native” approach for pre- and post-deployment, as well as the concept of application health checks.
It is an experimental project, under the umbrella of the [Keptn Application Lifecycle working group](https://github.com/keptn/wg-app-lifecycle).

The Keptn Lifecycle Controller is composed of the following components:

- Keptn Lifecycle Operator
- Keptn Scheduler

The Keptn Lifecycle Operator contains several controllers for Keptn CRDs and a Mutating Webhook.
The Keptn Scheduler ensures that Pods are started only after the pre-deployment checks have finished.

## Architecture

![](./assets/architecture.png)

A Kubernetes Manifest, which is annotated with Keptn specific annotations, gets applied to the Kubernetes Cluster.
Afterward, the Keptn Scheduler gets injected (via Mutating Webhook), and Kubernetes Events for Pre-Deployment are sent to the event stream.
The Event Controller watches for events and triggers a Kubernetes Job to fullfil the Pre-Deployment.
After the Pre-Deployment has finished, the Keptn Scheduler schedules the Pod to be deployed.
The KeptnApp and KeptnWorkload Controllers watchfor the workload resources to finish and then generate a Post-Deployment Event.
After the Post-Deployment checks, SLOs can be validated using an interface for retrieving SLI data from a provider, e.g, [Prometheus](https://prometheus.io/).
Finally, Keptn Lifecycle Controller exposes Metrics and Traces of the whole Deployment cycle with [OpenTelemetry](https://opentelemetry.io/).

## How it works

The following sections will provide insights on each component of the Keptn Lifecycle controller in terms of their purpose, responsibility, and communication with other components.
Furthermore, there will be a description on what CRD they monitor and a general overview of their fields.

### Webhook

The mutating webhook works only on resources that have Keptn annotations.
The mutation consists in changing the scheduler used for the deployment with the Keptn Scheduler.
The webhook should be as fast as possible and should not create/change any resource.

When the webhook receives a request for a new pod, it will look for the following annotations:

```
keptn.sh/app (optional)
keptn.sh/workload
```

Additionally, it will compute a version string, using a hash function that takes certain properties of the pod as parameters
(e.g. the images of its containers).
Next, it will look for an existing instance of a `Workload CRD` for the given workload name:

- If it finds the `Workload`, it will update its version according to the previously computed version string. In addition, it will include a reference to the ReplicaSet UID of the pod (i.e. the Pods owner), or the pod itself, if it does not have an owner.
- If it does not find a workload instance, it will create one containing the previously computed version string. In addition, it will include a reference to the ReplicaSet UID of the pod (i.e. the Pods owner), or the pod itself, if it does not have an owner. 
It will use the following annotations for
the specification of the pre/post deployment checks that should be executed for the `Workload`:
  - `keptn.sh/pre-deployment-tasks: task1,task2`
  - `keptn.sh/post-deployment-tasks: task1,task2`


After either one of those actions has been taken, the webhook will set the scheduler of the pod and allow the pod to be scheduled.


### Scheduler

After the Webhook mutation, the Keptn-Scheduler will handle the annotated resources. The scheduling flow follows that of the default scheduler, 
since it implements a scheduler plugin based on the [scheduling framework]( https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/).
For each pod, at the very end of the scheduling cycle, the plugin verifies whether the pre deployment checks have terminated, by retrieving the current status of the WorkloadInstance. Only if that is successful, the pod is binded to a node.


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

### Keptn Task

tbd

### Keptn Task Definition

tbd

## How to use

TBD

## How to install (development)

**Prerequisites:**

The lifecycle controllers includes a Mutating Webhook which requires TLS certificates to be mounted as a volume in its pod. The certificate creation
is handled automatically by [cert-manager](https://cert-manager.io). To install **cert-manager**, follow their [installation instructions](https://cert-manager.io/docs/installation/).

When cert-manager is installed, use the following commands to deploy the operator:

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
