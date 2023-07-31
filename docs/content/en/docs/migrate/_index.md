---
title: Migrating to the Keptn Lifecycle Toolkit
description: Notes to help you migrate from Keptn v1 to KLT
weight: 900
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

The Keptn Lifecycle Toolkit uses a different paradigm
than that used for Keptn v1
and so migration from Keptn v1 will not be a straight-forward process.
In this section, we will assemble information to help people
who want to move from Keptn v1 as it becomes available.

> **Note**
This section is under development.
Information that is published here has been reviewed for technical accuracy
but the format and content is still evolving.
We hope you will contribute your experiences
and questions that you have.

These instructions mostly assume that you want to utilize
the full Keptn Lifecycle Toolkit.
Note, however, that you can install and use some functionality
such as Keptn Metrics and Observability
without implementing all KLT features.

This section currently includes the following topics:

* [Evolution of KLT](evolution-klt)
  Understand the paradigm of KLT and how it evolved from Keptn v1.
  Also see whether migrating to KLT is appropriate for your deployments

* [Migration strategy](strategy) --
  Overview of the recommended migration strategy

* [Set up Kubernetes cluster with a deployment engine](setup) --
  Begin with a Kubernetes cluster with a deployment engine
  that deploys your software.
  This can be a new installation or an existing installation
  that meets the requirements for KLT.

* [Install and integrate KLT into your cluster](install-integrate]

* [Set up metrics and observability](metrics-observe) --
  When you identify the data-sources being used
  and provide KLT with information about your OpenTelemetry collector,
  KLT begins to accumulate information that you can view
  for your deployment.

* [Migrate CI/CD functionality](cicd)

* [Migrate Quality Gates to Keptn Evaluations](evaluations)

* [Migrate JES services to Keptn Tasks](jes)

* [Migrate remediation services to Day 2 monitoring](day2)
