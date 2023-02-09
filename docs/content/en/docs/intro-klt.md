---
title: Introduction to the Keptn Lifecycle Toolkit
linktitle: Introduction to the Keptn Lifecycle Toolkit
description: Understand the Keptn Lifecycle Toolkit
weight: 05
cascade:
  github_subdir: "docs/content/en/docs"
  path_base_for_github_subdir: "/content/en/docs-dev"
---

The Keptn Lifecycle Toolkit (KLT) makes your deployments observable,
brings application-awareness to your Kubernetes cluster,
and helps you reliably deliver your application with:

* Pre-Deployment Tasks: e.g. checking for dependant services,
  setting the cluster to be ready for the deployment, etc.
* Pre-Deployment Evaluations: e.g. evaluate metrics
  before your application gets deployed (e.g. layout of the cluster)
* Post-Deployment Tasks: e.g. trigger a test,
  trigger a deployment to another cluster, etc.
* Post-Deployment Evaluations: e.g. evaluate the deployment,
  evaluate the test results, etc.

All of these things can be executed on a workload or on an application level,
whereby a Keptn application is a collection of multiple workloads.

## Compare Keptn Lifecycle Toolkit and Keptn LTS

The Keptn Lifecycle Controller (KLT) is a Keptn subproject
whose design reflects lessons we learned while developing Keptn LTS.
KLT recognizes that tools such as Argo and Flux
are very good at deploying applications
so adds pre-deployment and post-deployment evaluations and actions.
For many installations, this provides the functionality they need
with much less complexity than the Keptn LTS project.

Keptn LTS is a fully-incubated, LTS release
that can deploy applications on platforms other than Kubernetes
and currently offers some functionality that KLT does not
such as the abiity to accomodate complex scoring algorithms
for SLO evaluations
and the ability to implement remediations (self-healing) for problems discovered
on the production site.
Keptn LTS can also be more difficult to set up and maintain
thank the Keptn Lifecycle Toolkit.

In a December 2022 Keptn Community meeting, 
we discussed the differences and similarities
between Keptn and the Keptn Lifecycle Toolkit
to help you decide which best fits your needs.
View the recording:
[Compare Keptn V1 and the Keptn Lifecycle Toolkit](https://www.youtube.com/watch?v=0nCbrG_RFos)

## Overviews of Keptn Lifecycle Toolkit

A number of presentations are available to give an overview
of the Keptn Lifecycle Toolkit:

* [Observability and Orchestration of your Deployment](https://www.youtube.com/watch?v=0nCbrG_RFos)

* [Keptn Lifecycle Toolkit Demo Tutorial on k3s, with ArgoCD for GitOps, OTel, Prometheus and Grafana](https://www.youtube.com/watch?v=6J_RzpmXoCc)

* The "IsItObservable Tutorial for KLT" will be released to YouTube soon.
  You can explore the GitHub repository that accompanies this video
  at https://github.com/isItObservable/keptn-lifecycle-Toolkit in the meantime.
