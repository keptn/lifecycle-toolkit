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

Keptn makes use of a custom scheduler
when running on Kubernetes v1.26 and earlier.
For Kubernetes v1.27 and later, scheduling is
implemented using
[Kubernetes scheduling gates](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-scheduling-readiness/),
unless the `schedulingGatesEnabled` Helm value is set to `false`.
See
[Keptn integration with Scheduling](../architecture/components/scheduler.md)
for details.

If Keptn is installed on a [vCluster](https://www.vcluster.com/) with
Kubernetes v1.26 or earlier, some extra configuration
needs to be added for full compatibility.
See
[Running Keptn with vCluster](#running-keptn-with-vcluster)
for more information.

## Running Keptn with vCluster

Keptn running on Kubernetes versions 1.26 and older
uses a custom
[scheduler](../architecture/components/scheduler.md),
so it does not work with
[Virtual Kubernetes Clusters](https://www.vcluster.com/)
("vClusters") out of the box.
This is also an issue
if the `lifecycleOperator.schedulingGatesEnabled` Helm value is set to `false`
for Kubernetes version 1.27 and later.
See
[Keptn integration with Scheduling](../architecture/components/scheduler.md)
for details.

To solve this problem:

1. Follow the instructions in
   [Separate vCluster Scheduler](https://www.vcluster.com/docs/architecture/scheduling#separate-vcluster-scheduler)
   to modify the vCluster `values.yaml` file
   to use a virtual scheduler.

1. Create or upgrade the vCluster,
   following the instructions in that same document.

1. Follow the instructions in the section below
   to install Keptn in that vCluster.
