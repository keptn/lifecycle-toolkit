---
title: Requirements
description: Supported software versions and information about resources required
weight: 15
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

## Supported Kubernetes versions

Keptn requires Kubernetes v1.24.0 or later.

Run the following to ensure that both client and server versions
are running Kubernetes versions greater than or equal to v1.24.
In this example, both client and server are at v1.24.0
so Keptn will work.

```shell
kubectl version --short
```

```shell
Client Version: v1.24.0
Kustomize Version: v4.5.4
Server Version: v1.24.0
```

Keptn uses a different scheduling method
when using Kubernetes versions 1.26 and older
or if the `schedulingGatesEnabled` Helm chart value is set to `false`
for Kubernetes versions greater than 1.26.  See
[Keptn integration with Scheduling](../architecture/components/scheduler/)
for details.

Note that you must modify your vCluster configuration
before installing Keptn when using Kubernetes versions 1.26 and older
or if the `schedulingGatesEnabled` Helm chart value is set to `false`
for Kubernetes versions greater than 1.26.  See
[Running Keptn with vCluster](install.md/#running-keptn-with-vcluster)
for more information.

## Resource requirements

## cert-manager

Keptn includes a lightweight cert-manager
that is used for installation and Webhooks.
You can configure a different cert-manager
before you install Keptn.
See [Implement your own cert-manager](../operate/cert-manager.md)
for instructions.
