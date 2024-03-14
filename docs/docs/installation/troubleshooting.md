---
comments: true
---

# Troubleshooting Guide

Welcome to the Keptn troubleshooting guide.

This guide will help you address common issues that you might encounter while using Keptn
for managing your [workloads](https://kubernetes.io/docs/concepts/workloads/).
Below are some common problems and their solutions:

## Keptn is installed but it is not aware of my workloads

If you are facing an issue where Keptn is installed but does not seem to be aware of your
[workloads](https://kubernetes.io/docs/concepts/workloads/), follow these steps:

1. Ensure that the namespace you wish to target is [annotated correctly](index.md#basic-installation).
2. Make sure your [workloads](https://kubernetes.io/docs/concepts/workloads/)
   (e.g., Deployment manifests) have the [three required annotations](https://lifecycle.keptn.sh/docs/implementing/integrate/#annotate-workloads).

## Keptn is causing my pods to be pending

If your pods are stuck in a pending state and Keptn seems to be the cause, it might be due
to a pre-deployment task or evaluation.
Follow these steps:

The most probable reason is that a pre-deployment task in your
[workload](https://kubernetes.io/docs/concepts/workloads/) is either failing or has not completed yet.

Failing pre-deployment evaluation tasks will prevent a pod from being scheduled.

Check the logs of the pre-deployment task Kubernetes Jobs for insights.

For instance, if
your application is in the `prod` namespace:

```shell
kubectl -n prod get pods
kubectl -n prod logs job/...
```

> **Note**
The blocking behavior can be changed by configuring non-blocking deployment
functionality.
More information can be found in the
[Keptn non-blocking deployment section](../components/lifecycle-operator/keptn-non-blocking.md).

## I have pending Pods after Keptn is uninstalled

> **Note** This section particularly affects clusters where
Keptn was installed via ArgoCD.

If you have uninstalled Keptn, originally installed via ArgoCD
and are now facing issues scheduling or deleting pods for other applications,
you probably didn't enable
[cascading deletion](https://kubernetes.io/docs/concepts/architecture/garbage-collection/#cascading-deletion)
of the application during installation, which is disabled by default in ArgoCD.

To fix this problem, you need to install Keptn via ArgoCD again, with the use
of `finalizers` in your Argo Application.
For more information see the
[Deploy Keptn via ArgoCD](./configuration/argocd.md) section for more information.

## I cannot see DORA metrics or OpenTelemetry traces

Keptn will automatically generate DORA metrics and OTel traces for every deployment, but
by default it does not know where to send them.

You need an OpenTelemetry collector
installed and configured on the cluster.

[The OpenTelemetry observability page](https://lifecycle.keptn.sh/docs/implementing/otel/)
contains more information on how to configure this.
