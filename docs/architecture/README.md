# Architecture - Keptn Lifecycle Controller

Documents the Keptn Lifecycle Controller architecture that could be used for deeper understanding of how it works under the hood.

> **_NOTE:_**  This section is under development

## Overview

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

![](/assets/architecture.png)


## How it works

The following sections will provide insights on each component of the Keptn Lifecycle controller in terms of their purpose, responsibility, and communication with other components.
Furthermore, there will be a description on what CRD they monitor and a general overview of their fields.

### Webhook

Annotating a namespace subjects it to the effects of the mutating webhook:

```
apiVersion: v1
kind: Namespace
metadata:
  name: podtato-kubectl
  annotations:
    keptn.sh/lifecycle-controller: "enabled"  # this lines tells the webhook to handle the namespace
```
However, the mutating webhook will modify only resources in the annotated namespace that have Keptn annotations.
When the webhook receives a request for a new pod, it will look for the workload annotations:

```
keptn.sh/workload
```
The mutation consists in changing the scheduler used for the deployment with the Keptn Scheduler. Webhook then creates a workload and app resource per annotated resource. 
You can also specify a custom app definition with the annotation:

```
keptn.sh/app
```
In this case the webhook will not generate an app, but it will expect that the user will provide one.
The webhook should be as fast as possible and should not create/change any resource.
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

and for the Evaluations:

  - `keptn.sh/pre-deployment-evaluations: my-evaluation-definition`
  - `keptn.sh/post-deployment-evaluations: my-eval-definition`

After either one of those actions has been taken, the webhook will set the scheduler of the pod and allow the pod to be scheduled.


### Scheduler

After the Webhook mutation, the Keptn-Scheduler will handle the annotated resources. The scheduling flow follows the default scheduler behavior,
since it implements a scheduler plugin based on the [scheduling framework]( https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/).
For each pod, at the very end of the scheduling cycle, the plugin verifies whether the pre deployment checks have terminated, by retrieving the current status of the WorkloadInstance. Only if that is successful, the pod is bound to a node.


### Keptn App

An App contains information about all workloads and checks associated with an application.
It will use the following structure for the specification of the pre/post deployment and pre/post evaluations checks that should be executed at app level:

```
apiVersion: lifecycle.keptn.sh/v1alpha1
kind: KeptnApp
metadata:
name: podtato-head
namespace: podtato-kubectl
spec:
version: "1.3"
workloads:
- name: podtato-head-left-arm
version: 0.1.0
- name: podtato-head-left-leg
postDeploymentTasks:
- post-deployment-hello
preDeploymentEvaluations:    
- my-prometheus-definition
```
While changes in the workload version will affect only workload checks,  a change in the app version will also cause a new execution of app level checks.

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

### Keptn Evaluation Definition
A `KeptnEvaluationDefinition` is a CRD used to define evaluation tasks that can be run by the Keptn Lifecycle Controller
as part of pre- and post-analysis phases of a workload or application.

A Keptn evaluation definition looks like the following:

```yaml
apiVersion: keptn.sh/v1
kind: KeptnEvaluationDefinition
metadata:
  name: my-prometheus-evaluation
spec:
  source: prometheus
  objectives:
    - name: query-1
      query: "xxxx"
      evaluationTarget: <20
    - name: query-2
      query: "yyyy"
      evaluationTarget: >4
```


### Keptn Evaluation Provider
A `KeptnEvaluationProvider` is a CRD used to define evaluation provider, which will provide data for the 
pre- and post-analysis phases of a workload or application.

A Keptn evaluation provider looks like the following:

```yaml
apiVersion: keptn.sh/v1
kind: KeptnEvaluationProvider
metadata:
  name: prometheus
spec:
  targetServer: "http://prometheus-k8s.monitoring.svc.cluster.local:9090"
  secretName: prometheusLoginCredentials
```


