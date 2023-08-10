---
title: Requirements
description: Supported software versions and information about resources required
weight: 15
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

## Supported Kubernetes versions

The Keptn Lifecycle Controller requires Kubernetes v1.24.0 or later.

Run the following to ensure that both client and server versions
are running Kubernetes versions greater than or equal to v1.24.
In this example, both client and server are at v1.24.0
so the Keptn Lifecycle Toolkit will work.

```shell
kubectl version --short
```

```shell
Client Version: v1.24.0
Kustomize Version: v4.5.4
Server Version: v1.24.0
```

KLT is not currently compatible with
[vcluster](<https://github.com/loft-sh/vcluster>).

## Resource requirements

## cert-manager

KLT includes a lightweight cert-manager
that is used for installation and Webhooks.
You can configure a different cert-manager
before you install KLT.
See [Implement your own cert-manager](../operate/cert-manager.md)
for instructions.
