---
comments: true
---

# Keptn for non-Kubernetes deployments

Keptn can interact with deployments that are not running on Kubernetes
even though Keptn mainly targets Kubernetes deployments.
The following functionality is available for deployments not on Kubernetes:

- [Run Keptn tasks](#run-keptntask-for-a-deployment-not-on-kubernetes)
- [Run Keptn analysis](#run-keptn-analysis-for-a-deployment-not-on-kubernetes)

To use either of these features,
you must set up a Kubernetes cluster and
[install](../installation/index.md#basic-installation)
Keptn on it,
but this can be a very lightweight, single-node KinD cluster; see
[Create local Kubernetes cluster](../installation/k8s.md#create-local-kubernetes-cluster).
Keptn only runs  on-demand `KeptnTask` and `KeptnAnalysis` resources
so resource utilization is minimal.

## Run KeptnTask for a deployment not on Kubernetes

Keptn tasks running on a Kubernetes cluster can be triggered for
[workloads](https://kubernetes.io/docs/concepts/workloads/)
and applications that are deployed outside of Kubernetes.
For example, Keptn can run (or "trigger")
load and performance tests
for an application that is deployed on a virtual machine,
or any other environment, as long as it can be reached
by the container running the KeptnTask.
It does this by specifying a container image that should be executed.
You specify the container in a `KeptnTaskDefinition` resource; see
[Deployment tasks](../guides/tasks.md) for more information.
The `KeptnTask` runs as a Kubernetes
[job](https://kubernetes.io/docs/concepts/workloads/controllers/job/)
on the cluster where Keptn is installed.

To implement this, install Keptn on a Kubernetes cluster
as described above, then::

- [Create a KeptnTaskDefinition](#create-a-keptntaskdefinition)
- [Create and apply a KeptnTask](#create-and-apply-a-keptntask)

### Create a KeptnTaskDefinition

When you have Keptn installed, create a
YAML file that defines what you want to execute
as a `KeptnTaskDefinition` resource.
See
[Deployment tasks](../guides/tasks.md)
and the
[KeptnTaskDefinition](../reference/crd-reference/taskdefinition.md)
reference page for more information.

For example, you might create a `test-task-definition.yaml` file
with the following content:

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnTaskDefinition
metadata:
  name: helloworldtask
spec:
  retries: 0
  timeout: 30s
  container:
    name: cowsay
    image: rancher/cowsay:latest
    args:
      - 'hello world'
```

This example uses the `container-runtime` runner,
but you can instead use the `deno-runtime` or `python-runtime` runner.
See
[Runners and containers](../guides/tasks.md#runners-and-containers)
for more information.

### Create and apply a KeptnTask

You must manually create the
[KeptnTask](../reference/crd-reference/task.md) resource.
In the standard operating mode,
when Keptn is managing
[workloads](https://kubernetes.io/docs/concepts/workloads/)
for deployments running on Kubernetes,
the creation of the `KeptnTask` resource is automatic.

Moreover, each time you want to execute a `KeptnTask`,
you must manually create a new (and uniquely named) `KeptnTask` resource.

The `KeptnTask` resource references the `KeptnTaskDefinition`
that you created above
in the `spec.taskDefinition` field.
For example, you might create a `test-task.yaml` file
with the following content:

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnTask
metadata:
  name: runhelloworld1
spec:
  taskDefinition: helloworldtask
  context:
    appName: "my-app"
    appVersion: "1.0.0"
    objectType: ""
    taskType: ""
    workloadName: "my-workload"
    workloadVersion: "1.0.0"
```

You can then apply this YAML file with the following command:

```shell
kubectl apply -f test-task.yaml -n my-keptn-annotated-namespace
```

Applying this file causes Keptn to create a Kubernetes
[job](https://kubernetes.io/docs/concepts/workloads/controllers/job/)
and run the executables defined
in the associated `KeptnTaskDefinition` resource.

Use the following commands to show the current status of the jobs:

```shell
kubectl get keptntasks -n my-keptn-annotated-namespace
kubectl get pods -n my-keptn-annotated-namespace
```

For subsequent KeptnTask runs,
the values of the `KeptnTask` `name` and `version` fields must be unique,
so copy the `KeptnTask` yaml file you have and update the
`metadata.name` field.

A standard practice is to just increment the value of the suffix field.
For example, you could create a `test-task-2.yaml` file
with the `metadata.name` field set to `runhelloworld2`:

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnTask
metadata:
  name: runhelloworld2
spec:
  taskDefinition: helloworldtask
  context:
    appName: "my-app"
    appVersion: "1.0.1"
    objectType: ""
    taskType: ""
    workloadName: "my-workload"
    workloadVersion: "1.0.1"
```

Use the following command to apply this resource:

```shell
kubectl apply -f test-task-2.yaml -n my-keptn-annotated-namespace
```

See the
[Deployment tasks](../guides/tasks.md)
guide and associated reference pages
for more information about running Keptn tasks
with deployments that do not run on Kubernetes.

## Run Keptn analysis for a deployment not on Kubernetes

The Keptn analyses feature
analyzes Service Level Objectives (SLOs)
based on Service Level Indicators (SLIs).
It can apply weights to reach a composite score
about the health of the deployment,
similar to what the metrics evaluations of the
Keptn v1 quality gates feature provided.
The data used can come from multiple instances
of multiple data providers
(such as Prometheus, Thanos, Dynatrace, and DataDog).

A Keptn analysis can be run for any application running anywhere
as long Keptn can access a monitoring provider endpoint
that serves metrics for the application.
You can point to multiple instances of the supported monitoring providers
(Prometheus, Thanos, Dynatrace, Datadog, and dql)
so the application itself can run anywhere.

To implement a Keptn analysis for your deployment:

- Create a `KeptnMetricProvider` resource
  for each data source to be used for your analysis.
  This specifies the URL for the data source,
  assigns a `name` that Keptn uses to reference that provider,
  and can define a secret for the data provider if necessary.

- Create `AnalysisValueTemplate` resources for each SLI
  and an `AnalysisDefinition` resource that contains all SLOs
  to be used in your analysis

- Create and apply an `Analysis` resource
  to define each specific analysis you want to run.

See the
[Analysis](../guides/slo.md)
guide and the
[Analyzing Application Performance with Keptn](https://keptn.sh/stable/blog/2023/12/19/analyzing-application-performance-with-keptn/)
blog
for more details and examples for the Keptn analysis feature.
