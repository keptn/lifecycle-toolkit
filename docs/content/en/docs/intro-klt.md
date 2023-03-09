---
title: Introduction to the Keptn Lifecycle Toolkit
linktitle: Introduction to the Keptn Lifecycle Toolkit
description: Understand the Keptn Lifecycle Toolkit
weight: 05
cascade:
  github_subdir: "docs/content/en/docs"
  path_base_for_github_subdir: "/content/en/docs-dev"
---

The Keptn Lifecycle Toolkit (KLT) implements observability
for deployments that are implemented with standard GitOps tools
such as ArgoCD, Flux, and Gitlab
and brings application awareness to your Kubernetes cluster.

These standard GitOps deployment tools
do an excellent job at deploying applications
but do not handle all issues
that are required to ensure that your deployment is usable.
The Keptn Lifecycle Toolkit "wraps" a standard Kubernetes GitOps deployment
with the capability to automatically handle issues
before and after the actual deployment.

Pre-deployment issues:

* Send appropriate notifications that this deployment is about to happen
* Check that downstream services meet their SLOs
* Verify that your infrastructure is ready
* Ensure that your infrastructure
  has the resources necessary for a successful deployment

Post-deployment issues:

* Integrate with tooling beyond the standard Kubernetes probes
* Automatically test the deployment
* Ensure that the deployment is meeting its SLOs
* Identify any downstream issues that may be caused by this deployment
* Send appropriate notifications about whether the deployment was successful or unsuccessful

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

All of these things can be executed for a workload or for a KeptnApp,
which is a collection of multiple workloads.

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
can accomodate complex scoring algorithms for SLO evaluations,
and can implement remediations (self-healing) for problems discovered
on the production site.

In a December 2022 Keptn Community meeting,
we discussed the differences and similarities
between Keptn and the Keptn Lifecycle Toolkit
to help you decide which best fits your needs.
View the recording:
[Compare Keptn V1 and the Keptn Lifecycle Toolkit](https://www.youtube.com/watch?v=0nCbrG_RFos)

## Overviews of Keptn Lifecycle Toolkit

A number of presentations are available to give an overview
of the Keptn Lifecycle Toolkit:

* [What is keptn, how it works, and how to get started!](https://www.dynatrace.com/news/blog/what-is-keptn-how-it-works-and-how-to-get-started/)

* [Observability and Orchestration of your Deployment](https://www.youtube.com/watch?v=0nCbrG_RFos)

* [Keptn Lifecycle Toolkit Demo Tutorial on k3s, with ArgoCD for GitOps, OTel, Prometheus and Grafana](https://www.youtube.com/watch?v=6J_RzpmXoCc)

* The "IsItObservable Tutorial for KLT" will be released to YouTube soon.
* You can explore the [GitHub repository](https://github.com/isItObservable/keptn-lifecycle-Toolkit)
* that accompanies this video in the meantime.
