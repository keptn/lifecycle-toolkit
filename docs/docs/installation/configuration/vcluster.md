---
comments: true
---

# vCluster installation

Keptn running on Kubernetes versions 1.26 and older
does not work with
[Virtual Kubernetes Clusters](https://www.vcluster.com/)
("vClusters") out of the box.
See
[Keptn integration with Scheduling](../../components/scheduling.md)
for details.

To solve this problem:

1. Follow the instructions in
   [Separate vCluster Scheduler](https://www.vcluster.com/docs/vcluster/configure/vcluster-yaml/control-plane/other/advanced/virtual-scheduler)
   to modify the vCluster `values.yaml` file
   to use a virtual scheduler.

1. Create or upgrade the vCluster,
   following the instructions in that same document.

1. Follow the instructions in the section below
   to install Keptn in that vCluster.
