---
title: Installation and upgrade
description: Learn how to install, configure, and upgrade the Keptn Lifecycle Toolkit
weight: 15
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

This section provides details about how to install and configure
the components of the Keptn Lifecycle Toolkit
either as a local cluster you use for study, testing, and demonstrations
or as part of an existing production cluster.
The steps are:

1. Understand the [Software versions and resources](reqs.md)
   that are required.
1. Be sure that your cluster includes the components discussed in
   [Prepare your cluster for KLT](k8s.md/#prepare-your-cluster-for-klt).
1. [Bring or create your Kubernetes cluster](k8s.md).
1. [Replace the default cert-manager](cert-manager.md) (optional).
   This step is only required if you want to replace
   the default KLT cert-manager with another cert-manager.
1. [Install the Keptn Lifecycle Toolkit](install.md).
1. [Enable KLT for your cluster](install.md/#enable-klt-for-your-cluster)
1. [Enable Keptn Lifecycle Toolkit](install.md/#enable-klt-for-your-cluster).
   This step is not required if you only want to run Keptn Metrics
   but is required for all other KLT features.

1. Run the following command to ensure that your Kuberetes cluster
   is ready to implement the Lifecycle Toolkit:

   ```shell
   kubectl get pods -n keptn-lifecycle-toolkit-system
   ```

   You should see pods for the following components:
   - certificate-operator (or another cert manager)
   - lifecycle-operator
   - scheduler
   - metrics-operator

Unless you are only using the customized Keptn metrics feature,
you now need to:

- Follow the instructions in
  [Annotate workload](../../implementing/integrate/#basic-annotations)
  to integrate the Lifecycle Toolkit into your Kubernetes cluster
  by applying basic annotations to your `Deployment` resource.
- Follow the instructions in
  [Define a Keptn application](../../implementing/integrate/#define-a-keptn-application)
  to create a Keptn application that aggragates
  all the `workloads` for your deployment into a single
  [KeptnApp](../../yaml-crd-ref/app) resource.

This section also includes:

1. How to [Upgrade](upgrade.md)
   to a new version of the Keptn Lifecycle Toolkit
