---
comments: true
---

# vCluster installation

Keptn running on Kubernetes versions 1.26 and older
uses a custom
[scheduler](../../components/scheduling.md),
so it does not work with
[Virtual Kubernetes Clusters](https://www.vcluster.com/)
("vClusters") out of the box.
This is also an issue
if the `lifecycleOperator.schedulingGatesEnabled` Helm value is set to `false`
for Kubernetes version 1.27 and later.
See
[Keptn integration with Scheduling](../../components/scheduling.md)
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
