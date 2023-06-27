---
title: Migration strategy
description: General guidelines for migrating your deployment to KLT
weight: 10
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

> **Note**
This section is under development.
Information that is published here has been reviewed for technical accuracy
but the format and content is still evolving.
We hope you will contribute your experiences
and questions that you have.

The Keptn Lifecycle Toolkit uses a different paradigm
than Keptn v1 and so migration is not a straight-forward process.
To help you wrap your mind around the process,
this lists Keptn v1 components
and identifies their disposition when you migrate to KLT.

TODO: Make the KLT terms links to relevant documentation
      after we figure out the mapping

## How many namespaces?

What set of these resources, etc need to be in the same namespace?
User has flexibility to decide.
See the Kubernetes
[Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
documentation for some basic information.
You can also search and find lots of "Best Practices for Namespaces"
documents published on the web.

Considerations:

* KLT primarily operates on Workloads (each of which has its own CRD)
  and Keptn Applications, (each of which has its own CRD
  and is defined by the value of the `part-of` annotation
  to each Workload.
* Each Keptn Metrics and evaluation resource
  identifies the namespace to which it applies.
* Each `KeptnApp` resource does identify the namespace to which it belongs
  so, if you configure multiple namespaces,
  you can have `KeptnApp` resources with the same name
  in multiple namespaces without having them conflict.
* You do not need separate namespaces for separate versions of your application.
  The `KeptnApp` resource includes fields to define
  the `version` as well as a `revision`
  (used if you have to rerun a deployment but want to retain the version number.

So, possible namespace designs run the gamut:

* Move all the Keptn v1 projects you had into a single namespace
* Move a group of Keptn v1 projects into a single namespace
* Define one namespace per Keptn v1 project
* Define one namespace per Keptn v1 stage

## Disposition of Keptn v1 components in KLT

* **project** -- May not have any meaning
* **stage** -- Depends on what the stage does,
  or maybe it has no meaning in KLT?
  * **deployment stage** becomes a `Deployment` workload in KLT
    or maybe it has no meaning and the focus is on the deployment task?
  * **testing stage** becomes a `KeptnTask` in KLT
    or maybe it has no meaning and the focus is on the test task?
  * The `Deployment` and associated `KeptnTask` resources
    form a `KeptnApp` resource
* **sequence** -- contains **tasks**
  * **deployment task** becomes a `Deployment` workload.
    You can code `KeptnTask Definition`
    and `KeptnEvaluationDefinition` resources
    that are configured to run either pre- or post-deployment tasks
  * **evaluation task** becomes a `KeptnEvaluation` resource
  * All other standard tasks
    (action, approval, get-action, rollback
    release, test)
    as well as custom task types
    you might define should be translated into
    `KeptnTaskDefinition` resources.
  * The `key:value` **properties** for each Keptn v1 sequence
    will be coded into the `KeptnTaskDefinition` resource
  * **TriggeredOn** -- TODO
* **SLI** -- should be translated into `KeptnMetric` resources
  with associated `KeptnMetricProvider` resources
  * Note that KLT allows you to support multiple data providers
    and multiple instances of each data provider for your SLI's.
* **SLO** -- KLT at this time does not support the full range
  of Quality Gates evaluations.
  Facilities such as weighting of SLI's and scoring of the evaluation
  do not currently exist.
  However, simple evaluations of an SLI can be defined
  as `KeptnEvaluationDefinition` resources.
* **Remediation** -- KLT does not currently support
  the same level of remediations as Keptn v1 does.
  KLT does provid some "Day 2" facilities:
  * Any query that is possible for your data provider post-deployment
    can be defined as a `KeptnMetricDefinition`;
    this is then reported as a Keptn Metric.
    Evaluation can be defined as a `KeptnEvaluationDefinition`.
  * `KeptnMetricDefinitions` can be retrieved and used
    to implement the Kubernetes HorizontalPodAutoscaler (HPA),
    which can detect the need for additional resources
    (more pods, memory, disk space, etc)
    and automatically add those resources to your configuration
    based on the `ReplicaSet` resources you have defined
* **Integrations and services in JES** -- Most can be defined
  as tasks using the `container-runtime` runner.
  In most cases, the code from the JES container
  can simply be moved into a `KeptnTaskDefinition` resource
  that uses the `container-runtime` runner.
  If the JES container code is written in JavaScript or TypeScript,
  you may be able to use the `deno-runtime` runner.
  If the JES container code is written in Python 3,
  you may be able to use the `python-runtime` runner.
  * No need for integrations for data providers;
    these are configured as `KeptnMetricsProvider` resources.
