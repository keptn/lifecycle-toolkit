---
title: Keptn for Non-Kubernetes Applications
description: Using Keptn with Non-Kubernetes Applications
weight: 20
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

It is possible to trigger Keptn Tasks for workloads and applications that are deployed outside of Kubernetes.

For example, to trigger a load test for an application deployed on a virtual machine.

You will still need to deploy Keptn on a Kubernetes cluster, but this can be a very lightweight, single-node kind cluster.
Keptn's only job is to trigger on-demand tasks, so resource utilization will be minimal.

## Step 1: Create a KeptnTaskDefinition

When you have Keptn installed, [create a KeptnTaskDefinition](../implementing/tasks/).

A `KeptnTaskDefinition` defines **what** you want to execute.

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

## Step 2: Create a KeptnTask for each run

Each time you want to execute a `KeptnTaskDefinition`, a new (and uniquely named) `KeptnTask` must be created.

In the standard operating mode, when Keptn is managing workloads, the creation of the `KeptnTask` CR is automatic.

Here though, we must create it ourselves.

The `KeptnTask` references the `KeptnTaskDefinition` in the `spec.taskDefinition` field:

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

Applying this file will cause Keptn to create a Job and a Pod and run the `KeptnTaskDefinition`.

`kubectl get keptntasks` and `kubectl get pods` will show the current status of the jobs.

## Running More KeptnTasks

For subsequent KeptnTask runs, the `KeptnTask` name needs to be unique, update the follow fields:

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
