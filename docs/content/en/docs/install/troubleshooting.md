---
title: Troubleshooting Guide
linktitle: Troubleshooting Guide
description: Troubleshoot common problems with Keptn
weight: 100
---

Welcome to the Keptn troubleshooting guide.

This guide will help you address common issues that you might encounter while using Keptn
for managing your [workloads](https://kubernetes.io/docs/concepts/workloads/).
Below are some common problems and their solutions:

## Keptn is installed but it is not aware of my workloads

If you are facing an issue where Keptn is installed but does not seem to be aware of your
[workloads](https://kubernetes.io/docs/concepts/workloads/), follow these steps:

1. Ensure that the namespace you wish to target is [annotated correctly](install.md/#enable-keptn-lifecycle-operator-observability).
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

## I have pending Pods after Keptn is uninstalled

> **Note**_ This section particularly affects clusters managed by ArgoCD.

If you have uninstalled Keptn and are now facing issues scheduling or deleting pods, follow these steps:

ArgoCD does not delete various CRDs and webhooks, when uninstalling applications, causing lingering resources.

### For cleanup instructions

1. Install Keptn & CRDs w/ Argo
2. Uninstall Keptn via Argo
3. Cluster is frozen, Pods cannot be deleted or scheduled

Trying to schedule a new pod results an `mpod.keptn.sh` error.

Reason: Argo doesn't remove the CRDs and global objects (mutating webhooks) and the
[default failurePolicy is to fail](https://github.com/keptn/lifecycle-toolkit/blob/650ecba95624ed3dc2bd61bf1f86578f450223a5/operator/config/webhook/manifests.yaml#L17).

to get further information refer to this [issue](https://github.com/keptn/lifecycle-toolkit/issues/1828).

## I cannot see DORA metrics or OpenTelemetry traces

Keptn will automatically generate DORA metrics and OTel traces for every deployment, but
by default it does not know where to send them.

You need an OpenTelemetry collector
installed and configured on the cluster.

[The OpenTelemetry observability page](https://lifecycle.keptn.sh/docs/implementing/otel/)
contains more information on how to configure this.
