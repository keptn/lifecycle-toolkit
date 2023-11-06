---
title: Keptn for Non-Kubernetes Applications
description: Using Keptn with Non-Kubernetes Applications
weight: 95
---

Keptn Tasks running on a Kubernetes cluster
can be triggered for [workloads](https://kubernetes.io/docs/concepts/workloads/) and applications
that are deployed outside of Kubernetes.
For example, Keptn could trigger load and performance tests
for an application that is deployed on a virtual machine.

To do this:

- [Install Keptn on a Kubernetes cluster](#install-keptn-on-a-kubernetes-cluster)
- [Create a KeptnTaskDefinition](#create-a-keptntaskdefinition)
- [Create and apply a KeptnTask](#create-and-apply-a-keptntask)
- [Re-run the KeptnTask](#re-run-the-keptntask)

## Install Keptn on a Kubernetes cluster

You must set up a Kubernetes cluster and
[install](../install/install.md/#basic-installation)
Keptn on it,
but this can be a very lightweight, single-node KinD cluster; see
[Create local Kubernetes cluster](../install/k8s.md/#create-local-kubernetes-cluster).
Keptn only triggers on-demand `KeptnTask` resources
so resource utilization is minimal.

## Create a KeptnTaskDefinition

When you have Keptn installed, create a
YAML file that defines what you want to execute
as a `KeptnTaskDefinition` resource.
See
[Deployment tasks](./tasks.md)
and the
[KeptnTaskDefinition](../yaml-crd-ref/taskdefinition.md/)
reference page for more information.

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

This example uses the `container-runtime` runner,
but you can instead use the `deno-runtime` or `python-runtime` runner.
See
[Runners and containers](./tasks.md#runners-and-containers)
for more information.

## Create and apply a KeptnTask

You must manually create the
[KeptnTask](../yaml-crd-ref/task.md) resource.
In the standard operating mode, when Keptn is managing [workloads](https://kubernetes.io/docs/concepts/workloads/),
the creation of the `KeptnTask` resource is automatic.

Moreover, each time you want to execute a `KeptnTask`,
you must manually create a new (and uniquely named) `KeptnTask` resource.

The `KeptnTask` resource references the `KeptnTaskDefinition`
that you created above
in the `spec.taskDefinition` field.
For example, you might create a `test-task.yaml` file
with the following content:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
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

Applying this file causes Keptn to create a Job and a Pod
and run the executables defined
in the associated `KeptnTaskDefinition` resource.

Use the following commands to show the current status of the jobs:

```shell
kubectl get keptntasks -n my-keptn-annotated-namespace
kubectl get pods -n my-keptn-annotated-namespace
```

## Re-run the KeptnTask

For subsequent KeptnTask runs,
the `KeptnTask` name and version fields must be unique,
so copy the `KeptnTask` yaml file you have and update the
`metadata.name` field.

A standard practice is to just increment the value of the suffix field.
For example, you could create a `test-task-2.yaml` file
with the `metadata.name` field set to `runhelloworld2`:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
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

You can then apply this file with the following command:

```shell
kubectl apply -f test-task-2.yaml -n my-keptn-annotated-namespace
```
