---
title: Introduction to Keptn Lifecycle Toolkit
description: An introduction to KLT and the usecases.
weight: 20
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

This section contains tutorials on how to use the Keptn Lifecycle Toolkit.

When you understand the toolkit, begin with the [getting started guide](../getting-started/).

----

The Keptn Lifecycle Toolkit (KLT) implements observability
for deployments and seamlessly integrates with GitOps tools
such as ArgoCD, Flux, and Gitlab
and brings application awareness to your Kubernetes cluster.

These standard GitOps deployment tools
do an excellent job at deploying applications
but do not handle all issues
that are required to ensure that your deployment is usable.
The Keptn Lifecycle Toolkit "wraps" a standard Kubernetes GitOps deployment
with the capability to automatically handle issues
before and after the actual deployment.

Pre-deployment issues that Keptn Lifecycle Toolkit can handle:

* Send appropriate notifications that this deployment is about to happen
* Check that downstream services meet their SLOs
* Verify that your infrastructure is ready
* Ensure that your infrastructure
  has the resources necessary for a successful deployment

Post-deployment issues that Keptn Lifecycle Toolkit can handle:

* Integrate with tooling beyond the standard Kubernetes probes
* Automatically test the deployment
* Ensure that the deployment is meeting its SLOs
* Identify any downstream issues that may be caused by this deployment
* Send appropriate notifications
  about whether the deployment was successful or unsuccessful

KLT can evaluate both workload (single service) tests
and SLO evaluations before and after the actual deployment.
Multiple workloads can also be logically grouped and evaluated
as a single cohesive unit called a `KeptnApp`.
In other words, a `KeptnApp` is a collection of multiple workloads.

KLT is tool- and vendor neutral and does not depend on particular GitOps tooling.
KLT emits signals at every stage
(Kubernetes events, OpenTelemetry metrics and traces)
to ensure that your deployments are observable.
It supports the following steps:

* Pre-Deployment Tasks: e.g. checking for dependant services,
  setting the cluster to be ready for the deployment, etc.
* Pre-Deployment Evaluations: e.g. evaluate metrics
  before your application gets deployed (e.g. layout of the cluster)
* Post-Deployment Tasks: e.g. trigger a test,
  trigger a deployment to another cluster, etc.
* Post-Deployment Evaluations: e.g. evaluate the deployment,
  evaluate the test results, etc.

All of these things can be executed for a workload or for a [KeptnApp](https://lifecycle.keptn.sh/docs/yaml-crd-ref/app/).

## Main features of : Metrics, Observability and Release lifecycle

* **Custom Metrics:** The Custom Keptn Metrics feature in the
Keptn Lifecycle Toolkit allows you to define metrics from
multiple data sources in your Kubernetes cluster.
It supports deployment tools like Argo, Flux, KEDA, HPA, or
Keptn for automated decision-making based on observability data.
Your observability data may come from multiple observability solutions
– Prometheus, Dynatrace, Datadog and others – or may be data that comes
directly from your cloud provider such as AWS, Google, or Azure.
The Keptn Metrics Server unifies and standardizes access to data from
various sources, simplifying configuration and integration into a single
set of metrics.

* **Observability:** The Keptn Lifecycle Toolkit (KLT) ensures observability
for Kubernetes deployments by creating a comprehensive trace of all Kubernetes
activities within a deployment.
It introduces the concept of applications, which connect logically related
workloads using different deployment strategies.
With KLT, you can easily understand deployment durations and failures across
multiple strategies.
It captures DORA metrics and exposes them as OpenTelemetry metrics.
The observability data includes out-of-the-box DORA metrics, traces from
OpenTelemetry, and custom Keptn metrics from configured data providers.
Visualizing this information is possible using dashboard tools like Grafana.

* **Release Lifecycle:** The Lifecycle Toolkit offers versatile functionalities
for deployment scenarios, including pre-deployment validation, image scanning,
and post-deployment tasks like test execution and stakeholder notification.
It automatically validates against Service Level Objectives (SLOs) and provides
end-to-end deployment traceability.
The toolkit extends deployments with application-aware tasks and evaluations,
allowing checks before or after deployment initiation.
It validates Keptn metrics using the Keptn Metrics Server, ensuring a healthy
environment and confirming software health against SLOs like performance and
user experience.
Additionally, it enables monitoring of new logs from log monitoring solutions.

To get started with Keptn Lifecycle Toolkit, refer to the
[Getting Started Exercises](https://lifecycle.keptn.sh/docs/getting-started/)
for detailed instructions and examples.
This guide will walk you through the installation process and help you set up
your environment for using KLT effectively.

## Compare Keptn Lifecycle Toolkit and Keptn LTS

The Keptn Lifecycle Controller (KLT) is a Keptn subproject
whose design reflects lessons we learned while developing Keptn LTS.
KLT recognizes that tools such as Argo and Flux
are very good at deploying applications.
However, these deployment tools do not provide
pre-deployment and post-deployment evaluations and actions;
this is what KLT adds.

Keptn LTS is a long-term support release
that can deploy applications on platforms other than Kubernetes,
can accommodate complex scoring algorithms for SLO evaluations,
and can implement remediations (self-healing) for problems discovered
on the production site.

Keptn Lifecycle Toolkit includes multiple features
that can be implemented independently or together.
Different features are at different levels of stability.
See the [Keptn Lifecycle Toolkit README file](https://github.com/keptn/lifecycle-toolkit/blob/main/README.md)
for a list of the features that have been implemented to date
and their level of stability.

In a December 2022 Keptn Community meeting,
we discussed the differences and similarities
between Keptn and the Keptn Lifecycle Toolkit
to help you decide which best fits your needs.
View the recording:
[Compare Keptn V1 and the Keptn Lifecycle Toolkit](https://www.youtube.com/watch?v=-cKyUKFjtwE&t=170s)

## Overviews of Keptn Lifecycle Toolkit

A number of presentations are available to help you understand
the Keptn Lifecycle Toolkit:

* [Orchestrating and Observing GitOps Deployments with Keptn](https://www.youtube.com/watch?v=-cKyUKFjtwE&t=11s)
  discusses the evolution of Keptn
  and the concepts that drive the Keptn Lifecycle Toolkit,
  then gives a simple demonstration of a Keptn Lifecycle Controller implementation.

* [Introducing Keptn Lifecycle Toolkit](https://youtu.be/449HAFYkUlY)
  gives an overview of what KLT does and how to implement it.

* [Keptn Lifecycle Toolkit Demo Tutorial on k3s, with ArgoCD for GitOps, OTel, Prometheus and Grafana](https://www.youtube.com/watch?v=6J_RzpmXoCc)
  is a short video demonstration of how the Keptn Lifecycle Toolkit works.
  You can download the exercise and run it for yourself;
  notes below the video give a link to the github repo.
  The README file in that repo gives instructions for installing the software
  either automatically or manually.

* [What is the Keptn Lifecycle Toolkit?](https://isitobservable.io/observability/kubernetes/what-is-the-keptn-lifecycle-toolkit)
  blog discusses KLT as part of the "Is It Observable?" series.
  This links to:

  * [What is Keptn Lifecycle Toolkit?](https://www.youtube.com/watch?v=Uvg4uG8AbFg)
    is a video that steps through the process of integrating KLT
    with your existing cloud native cluster.

* [Keptn Lifecycle Toolkit: Installation and KeptnTask Creation in Minutes](https://www.youtube.com/watch?v=Hh01bBwZ_qM)
  demonstrates how to install KLT and create your first KeptnTask in less than ten minutes.
  
* You can explore the [GitHub repository](https://github.com/isItObservable/keptn-lifecycle-Toolkit)
  for more information.
