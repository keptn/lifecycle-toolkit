# Troubleshooting Common Problems with KLT

Welcome to the troubleshooting guide for KLT (Keptn Lifecycle Toolkit). 

This guide will help you address common issues that you might encounter while using KLT for managing your workloads. Below are some common problems and their solutions:

## KLT is installed but it isn't aware of my workloads

If you are facing an issue where KLT is installed but doesn't seem to be aware of your workloads, follow these steps:

1. Ensure that the namespace you wish to target is [annotated correctly](https://lifecycle.keptn.sh/docs/install/install/#enable-klt-for-your-cluster).
2. Make sure your workloads (e.g., Deployment manifests) have the [three required annotations](https://lifecycle.keptn.sh/docs/implementing/integrate/#annotate-workloads).

## KLT is causing my pods to be pending

If your pods are stuck in a pending state and KLT seems to be the cause, it might be due to a pre-evaluation task. Follow these steps:

The most probable reason is that a pre-evaluation task in your workload is either failing or hasn't completed yet.

Failing pre-evaluation tasks will prevent a pod from being deployed.

Check the logs of the pre-evaluation task Kubernetes Jobs for insights. For instance, if your application is in the `prod` namespace:

```bash
kubectl -n prod get pods
kubectl -n prod logs job/...
```

## I have pending Pods after KLT is uninstalled

> **_NOTE:_**  This section particularly affects clusters managed by ArgoCD.

If you have uninstalled Keptn Lifecycle Toolkit and are now facing issues scheduling or deleting pods, follow these steps:

ArgoCD does not delete various CRDs and webhooks, when uninstalling applications, causing lingering resources.

For cleanup instructions, refer to this [issue](https://github.com/keptn/lifecycle-toolkit/issues/1828).


## I cannot see DORA metrics or OpenTelemetry traces

KLT will automatically generate DORA metrics and Otel traces for every deployment, but by default it doesn't know where to send them. You need an OpenTelemetry collector installed and configured on the cluster.

[The OpenTelemetry observability page](https://lifecycle.keptn.sh/docs/implementing/otel/) contains more information on how to configure this.