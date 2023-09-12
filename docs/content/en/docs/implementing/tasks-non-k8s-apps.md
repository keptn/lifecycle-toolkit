---
title: Keptn for Non-Kubernetes Applications
description: Using Keptn with Non-Kubernetes Applications
weight: 95
---

It is possible to trigger Keptn Tasks for workloads and applications
that are deployed outside of Kubernetes.
For example, Keptn could trigger load and performance tests
for an application that is deployed on a virtual machine.

You must still
[install](../install/install/#use-helm-chart)
and
[enable](../install/install/#enable-keptn-for-your-cluster)
Keptn on a Kubernetes cluster,
but this can be a very lightweight, single-node KinD cluster; see
[Create local Kubernetes cluster](../install/k8s/#create-local-kubernetes-cluster).
Keptn only triggers on-demand tasks, so resource utilization is minimal.

## Step 1: Create a KeptnTaskDefinition

When you have Keptn installed, create a
[KeptnTaskDefinition](../yaml-crd-ref/taskdefinition/)
resource that defines what you want to execute.
See
[Deployment tasks](../implementing/tasks/)
for more information.

For example:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
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

This example uses the `container-runtime` runner
because it allows the most flexibility
but you can also use the `deno-runtime` or `python-runtime` runners.

## Step 2: Create a KeptnTask for each run

Each time you want to execute a `KeptnTaskDefinition`,
a new (and uniquely named) `KeptnTask` must be created.

In the standard operating mode, when Keptn is managing workloads,
the creation of the `KeptnTask` CR is automatic.
Here though, we must create it ourselves.

The `KeptnTask` references the `KeptnTaskDefinition`
in the `spec.taskDefinition` field:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTask
metadata:
  name: runhelloworld1
spec:
  workload: "my-workload"
  workloadVersion: "1.0.0"
  appVersion: "1.0.0"
  app: "my-app"
  taskDefinition: helloworldtask
  context:
    appName: "my-app"
    appVersion: "1.0.0"
    objectType: ""
    taskType: ""
    workloadName: "my-workload"
    workloadVersion: "1.0.0"
```

TODO: This file does not match what I see in the API Reference.
See specific comments in the `KeptnTask` reference page.
When we resolve those issues, I will modify this file appropriately.

Applying this file causes Keptn to create a Job and a Pod
and run the `KeptnTaskDefinition`.

Use the following commands to show the current status of the jobs:

```shell
kubectl get keptntasks 
kubectl get pods
```

## Running More KeptnTasks

For subsequent KeptnTask runs, the `KeptnTask` name needs to be unique,
so update the following fields:

- `name`
- `spec.appVersion`
- `spec.workloadVersion`
- `spec.context.appVersion`
- `spec.context.workloadVersion`

For example:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTask
metadata:
  name: runhelloworld2
spec:
  workload: "my-workload"
  workloadVersion: "1.0.1"
  appVersion: "1.0.1"
  app: "my-app"
  taskDefinition: helloworldtask
  context:
    appName: "my-app"
    appVersion: "1.0.1"
    objectType: ""
    taskType: ""
    workloadName: "my-workload"
    workloadVersion: "1.0.1"
```
