---
title: Keptn for Non-Kubernetes Applications
description: Using Keptn with Non-Kubernetes Applications
weight: 95
---

Keptn Tasks can be triggered for workloads and applications
that are deployed outside of Kubernetes.
For example, Keptn could trigger load and performance tests
for an application that is deployed on a virtual machine.

To do this:

1. [Install and enable a Kubernetes cluster](#install-and-enable-a-kubernetes-cluster)
1. [Create a KeptnTaskDefinition](#create-a-keptntaskdefinition)
1. [Create a KeptnTask](#create-and-appy-a-keptntask)
1. [Re-run the KeptnTask](#re-run-the-keptntask)

## Install and enable a Kubernetes cluster

You must still
[install](../install/install/#use-helm-chart)
and
[enable](../install/install/#enable-keptn-for-your-cluster)
Keptn on a Kubernetes cluster,
but this can be a very lightweight, single-node KinD cluster; see
[Create local Kubernetes cluster](../install/k8s/#create-local-kubernetes-cluster).
Keptn only triggers on-demand `KeptnTask` resources
so resource utilization is minimal.

TODO: How is this cluster associated with the VM
where the deployment is running?

## Create a KeptnTaskDefinition

When you have Keptn installed, create a
[KeptnTaskDefinition](../yaml-crd-ref/taskdefinition/)
YAML file that defines what you want to execute.
See
[Deployment tasks](../implementing/tasks/)
for more information.

For example, you might create a `test-task-definition.yaml` file
with the following content:

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

## Create and apply a KeptnTask

You must manually create the
[KeptnTask](../yaml-crd-ref/task)
the resource.
In the standard operating mode, when Keptn is managing workloads,
the creation of the `KeptnTask` resource is automatic.
Here though, we must create it ourselves.

Moreover, you must create a new (and uniquely named)
`KeptnTask` resource
each time you want to rerun this task
Each time you want to execute a `KeptnTask` resource,
you must manually create a
a new (and uniquely named)
[KeptnTask](../yaml-crd-ref/task)
YAML file to describe that resource.


The `KeptnTask` references the `KeptnTaskDefinition`
in the `spec.taskDefinition` field.
For example, you might create a `test-task.yaml` file
with the following content:

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

You can then apply this YAML file with the following command:

```yaml
kubectl --apply test-task.yaml
```

Applying this file causes Keptn to create a Job and a Pod
and run the executables defined
in the associated `KeptnTaskDefinition` resource.

Use the following commands to show the current status of the jobs:

```shell
kubectl get keptntasks 
kubectl get pods
```

## Re-run the KeptnTask

For subsequent KeptnTask runs,
the `KeptnTask` name and version fields must be unique,
so copy the `KeptnTask` file you have and update the following fields:

- `name`
- `spec.appVersion`
- `spec.workloadVersion`
- `spec.context.appVersion`
- `spec.context.workloadVersion`

A standard practice is to just increment the values of these fields.
For example, you could create a `test-task-2.yaml` file
with the following content:

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

You can then apply this file with the following command:

```yaml
kubectl --apply test-task-2.yaml
```


