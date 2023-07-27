---
title: Install, enable, and integrate KLT into your cluster
description: How to add KLT to your cluster
weight: 30
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

> **Note**
This section is under development.
Information that is published here has been reviewed for technical accuracy
but the format and content is still evolving.
We hope you will contribute your experiences
and questions that you have.

After you have confirmed that your deployments are running as they should,
it is time to install and enable KLT on your cluster
and to integrate it into your workloads.
You can then define your first Keptn Application
that you will implement to familiarize yourself with KLT.

Once you have integrated KLT into your workloads,
KLT begins to gather DORA metrics
that you can view to get some useful information
about how your deployments are behaving.

## Install KLT

Follow the instructions in
[Install and enable KLT](../../install/install.md)
to install KLT in your cluster.
We recommend installing from the Helm chart
rather than the manifest.

In most cases,
it is not necessary to modify the Helm chart before installation
but instructions are provided if you need them.

The most common modification required when beginning the migration process
is to replace the default KLT `cert-manager` with another cert-manager.
See
[Use your own cert-manager](../../install/cert-manager.md) for details.

## Enable KLT

To enable KLT, annotate the
[Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
resource for each namespace where you want to run KLT.
See
[Enable KLT for your cluster](../../install/install.md/#enable-klt-for-your-cluster)
for details.

## Integrate KLT with your Workloads

Follow the instructions in
[Basic annotations](../../implementing/integrate/#basic-annotations)
to annotate your
[Workloads](https://kubernetes.io/docs/concepts/workloads/)
([Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
[DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
and
[ReplicaSets](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/)
in the namespaces where KLT is enabled.
KLT creates
[KeptnWorkload](../../crd-ref/lifecycle/v1alpha3/#keptnworkload)
resources based on the workload name and version you specify
in these annotations.

Recommended but optional is that you also populate the
`keptn.sh/app` or `app.kubernetes.io/part-of` annotation
to identify the workloads that should be included in each
[KeptnApp](../../yaml-crd-ref/app.md)
resource
so that KLT can automatically generate the `KeptnApp` resources you need.
See
[Use Keptn automatic app discovery](../../implementing/integrate/#use-keptn-automatic-app-discovery)
for more details.

In some cases, the Keptn v1 `projects` you have defined
may map to KLT `KeptnApp` resources
but in most cases you need to architect
what applications you need.

As you begin your migration,
you may want to define a single `KeptnApp` resource
that uses one of your `Deployment` resources
and work through the migration process for that one `Deployment`.
This gives you a chance to familiarize yourself with KLT
before introducing the complexity of multiple Keptn Applications.

## View DORA metrics

KLT begins collecting DORA metrics for your cluster
as soon as you annotate the Workload resources you are using.
See
[Dora metrics](../../implementing/dora)
for more information about DORA metrics
and how to view them.
