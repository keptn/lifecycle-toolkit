---
title: Set up Kubernetes deployment environment
description: Set up Kubernetes cluster with a deployment engine for your software
weight: 20
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

> **Note**
This section is under development.
Information that is published here has been reviewed for technical accuracy
but the format and content is still evolving.
We hope you will contribute your experiences
and questions that you have.

The Keptn Lifecycle Toolkit (KLT) supplements Kubernetes deployments
done with tools such as Argo CD, Flux, or even `kubectl apply`.
This is a major paradigm switch from Keptn v1,
which works best when it orchestrates deployments.

Consequently, the first step in your migration path
is to set up a new Kubernetes cluster
and implement the deployment of all your software in that cluster.
You can use the deployment tool(s) of your choice,
and you can use different tools to deploy different components.
KLT interacts with the Kubernetes API
and so interacts with all deployments in the same way,
regardless of the deployment engine used.

So the first steps of your migration are:

* Check the [Requirements](../install/reqs.md) page
  to ensure that the Kubernetes cluster you create
  is appropriate for KLT.
* Create your [Kubernetes cluster](../install/k8s.md).
* Install the deployment engine(s) in that cluster,
  following the documentation for each tool.
* Set up deployments for all components
  that should be built as part of your software.
* Test the deployments to ensure that they are working correctly.

[Prepare your cluster for KLT](../install/k8s/#prepare-your-cluster-for-klt)
lists other software that you will need to install in your cluster
to use Keptn Metrics and the standardized observability feature.
You can install these now or later.

When you have verified that your deployments are working well,
you are ready to
[Install and integrate KLT into your cluster](../install).
