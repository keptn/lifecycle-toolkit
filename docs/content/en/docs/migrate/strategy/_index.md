---
title: Migration strategy
description: Considerations for architecting the migration
weight: 10
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

> **Note**
This section is under development.
Information that is published here has been reviewed for technical accuracy
but the format and content is still evolving.
We hope you will contribute your experiences
and questions that you have.

Before you begin the migration project,
we suggest that you run through the exercises in
[Getting started](../../getting-started/)
to familiarize yourself with how the Keptn Lifecycle Toolkit works.
When you are ready to begin the migration,
follow the instructions in
[Installation and upgrade](../../install)
to set up your Kubernetes cluster
and install KLT on it.

On this page, we discuss some of the major paradigm shifts
in KLT relative to Keptn v1
and then discuss how elements of your Keptn v1 can be implemented
for the Keptn Lifecycle Toolkit.

## The Keptn Lifecycle Toolkit paradigm

The Keptn Lifecycle Toolkit uses a different paradigm
than Keptn v1 and so migration requires technical adjustments.
Much of the logic and functionality of your Keptn v1 projects
can be migrated to KLT
but must be rewritten to utilize KLT components.

Some key points:

* KLT uses native Kubernetes CRDs
  to configure SLIs/SLOs, tasks, and other elements
  that are part of the environment where the application deployment lives
  rather than using its own Git repo and its
  [shipyard.yaml](https://keptn.sh/docs/1.0.x/reference/files/shipyard/)
  file as Keptn V1 does.
  * See the
    [CRD Reference](../../yaml-crd-ref)
    section for pages that describe the Keptn manifests
    that you populate manually for KLT.
  * See the
    [API Reference](../../crd-ref)
    for a comprehensive reference to all resources
    defined for Keptn.

* KLT is not a delivery tool
  because it does not provide promotion.
  Instead, it works with standard deployment tools
  such as ArgoCD, Flux, even `kubectl -- apply`.
  KLT then prevents the scheduling and deployment
  of workflows if the environment does not meet
  the user-defined requirements.

* KLT operates on a
  [KeptnApp](../../yaml-crd-ref/app.md)
  resource
  that is an amalgamation of multiple Kubernetes workloads,
  which together comprise the released application.
  Each workload deploys a separate microservice,
  which can be deployed at different times using different tools.

* KLT integrates with continuous delivery tools
  to insures that a `KeptnApp` is not deployed
  if it does not meet the user-defined requirements
  for all the associated workloads.
  It also exposes metrics to evaluate the success/status of a deployment.

* KLT provides an operator
  that can observe and orchestrate application-aware workload life cycles.
  This operator leverages Kubernetes webhooks
  and extends the Kubernetes scheduler
  to support pre- and post-deployment hooks.
  When the operator detects a new version of a service
  (implemented as a Kubernetes
  [Workload](https://kubernetes.io/docs/concepts/workloads/)),
  it can execute pre- and post-deployment evaluations and tasks
  using Kubernetes capabilities.

* KLT provides extensive observability data
  using OpenTelemetry and Prometheus
  rather than storing the data in a special Keptn database.
  This data can be displayed using Grafana and Jaeger
  or the dashboard of your choice.

For in-depth information about KLT components
and how they work, see the
[Architecture](../../concepts/architecture)
section.

## Disposition of Keptn v1 components in KLT

To help you wrap your mind around the migration process,
this lists Keptn v1 components
and identifies their possible disposition when you migrate to KLT.

### project

A Keptn v1 project represents an arbitrary, higher-level grouping of services.
A project is defined by a
[shipyard.yaml](https://keptn.sh/docs/1.0.x/reference/files/shipyard/)
file.

KLT does not recognize projects.
Many Keptn v1 projects may translate into `KeptnApp` resources
but others will not.
For example, if your Keptn v1 project has a large number of services,
you may want to aggregate them into separate `KeptnApp` resources.
A project has a 1:1 mapping to a Git repo,
but, since few applications are stored in a mono-repo,
it is difficult to implement a project-as-application paradigm.

In general, the Keptn v1 project is a useful as a reference list
when migrating to KLT
but may not directly translate to a KLT resource.

### service

A Keptn v1 service models a smaller chunk of a project.
Most projects include many services.
In a micro-services world,
a service may represent a micro-service
but it could instead be a wrapper for something else,
like "the entire public website"

KLT does not have the concept of a service.
When migrating to KLT,
you need to analyze what each service is doing
and translate that into an appropriate resource.
The closest analogy is a Kubernetes
[workload](https://kubernetes.io/docs/concepts/workloads/)
but some services may be translated into
[KeptnTaskDefinition](../../yaml-crd-ref/app.md)
or other resources.
 See
[Working with Keptn tasks](../../implementing/tasks/)
for more information.

For example:

* A Keptn v1 service that runs chaos or load tests
  can probably be translated into
  a `KeptnTask` using the `container-runner`.
* A Keptn v1 service that runs a database
  can probably be translated
  into a Kubernetes `StateFulSet` workload; see
  [Workload Resources](https://kubernetes.io/docs/concepts/workloads/controllers/)
  for more information.
* A Keptn v1 service that runs a webserver
  can probably be translated into
  a Kubernetes `Deployment` workload.

### stage

A stage is a subsection of a project.
Because KLT is not a delivery tool,
it has no concept of a `stage`
but rather depends on a deployment engine.
However, the logic of the stages can be useful
when architecting the migration:

* A **deployment stage** -- may define sequences of tasks
    that should be translated into
    [KeptnTaskDefinition](../../yaml-crd-ref/taskdefinition.md)
    resources that are executed pre- and post-deployment
* A **testing stage** may define sequences of tasks
    that should be translated into `KeptnTaskDefinition` resources
    that are executed pre- and post-deployment.

Stage functionality could be implemented in many different ways.
Some functionality might be implemented in different namespaces
or even in different KLT-enabled clusters,
allowing a tool such as ArgoCD to handle promotion.

### sequence

A sequence is an ordered and user-defined sequence of tasks
that are defined in a
[shipyard.yaml](https://keptn.sh/docs/1.0.x/reference/files/shipyard/)
file.
The shipyard controller micro-service reads the shipyard file
and (when the time is right),
emits a `taskName.triggered` cloud event onto the event bus.
The shipyard controller waits to receive a `taskName.started` event
and a correspondingly equal number of `taskName.finished` events
before the shipyard controller reads the next task
and emits a `taskName.finished` event for that task.

In this way, you can define arbitrary sequences of any tasks
at any length and also link (or chain) sequences together
to form (primitive) workflows.
When migrating, these sequences of tasks can often be translated into
[KeptnTaskDefinition](../../yaml-crd-ref/taskdefinition.md)
resources that are defined to run either pre- or post-deployment
of the pod-scheduling phase.

The `shipyard` file is a general purpose workflow engine
that is backed by cloud events.
It is not opinionated to a tool, platform, technology,
or a particular "slice" of the lifecycle.
The `TriggeredOn` property allows
a Keptn v1 sequence to be triggered at any time
by a user or another system.
This capability can be used, for example,
to trigger a data encryption and backup operation,
or a file movement over a network, or other arbitrary activities
that may or may not have anything to do with an application's lifecycle.

When migrating to KLT,
sequences that are not part of the lifecycle workflow
should not be handled by KLT
but should instead be handled by the pipeline engine tools being used
such as Jenkins, Argo Workflows, Flux, and Tekton.

### task

Keptn v1 defines some specific types of tasks,
each of which is translated to a KLT resource
that is appropriate for the activity:

* A Keptn v1 **deployment task** becomes a
  [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
  [StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
  or [DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
  workload.
  You can code
  [KeptnTaskDefinition](../../yaml-crd-ref/taskdefinition.md)
  and
  [KeptnEvaluationDefinition](../../yaml-crd-ref/evaluationdefinition.md)
  resources that are configured
  to run either pre- or post-deployment tasks
* An **evaluation task** becomes a
  [KeptnEvaluationDefinition](../../yaml-crd-ref/evaluationdefinition.md)
  resource.
* All other standard tasks
  (**action**, **approval**, **get-action**, **rollback**,
  **release**, **test**)
  as well as custom task types
  that might be defined should be translated into
  `KeptnTaskDefinition` resources.
* The `key:value` **properties** for each Keptn v1 sequence
  should be coded into the `KeptnTaskDefinition` resource
  as appropriate.

Note that all `KeptnTask` resources at the same level
(either pre-deployment or post-deployment)
execute in parallel
whereas Keptn v1 sequences and tasks can not execute in parallel.

* If you have actions that need to execute sequentially,
  create a single `KeptnTaskDefinition` that calls each action in order.
* If you have tasks that can execute in parallel,
  migrating to KLT may improve the performance of the deployment.

### SLIs

Keptn v1
[SLIs](https://keptn.sh/docs/1.0.x/reference/files/sli/)
(Service Level Indicators)
represent queries from the data provider
such as Prometheus, Dynatrace, or Datadog,
which is configured as a Keptn integration.

When migrating to KLT, you need to define a
[KeptnMetricsProvider](../../yaml-crd-ref/metricsprovider.md)
resource for the data provider(s) you are using.
Note that KLT allows you to support multiple data providers
and multiple instances of each data provider for your SLIs
whereas Keptn v1 only allows you to use one SLI per project.

The queries defined for the Keptn v1 SLIs
should be translated into
[KeptnMetric](../../yaml-crd-ref/metric.md)
resources.

### SLOs

KLT at this time does not support the full range
of Quality Gates evaluations that are represented by
[SLOs](https://keptn.sh/docs/1.0.x/reference/files/slo/).
Facilities such as weighting of SLIs and scoring of the evaluation
do not currently exist.
This functionality is under development; see
[Epic 1646](https://github.com/keptn/lifecycle-toolkit/issues/1646).

However, simple evaluations of an SLI can be defined as
[KeptnEvaluationDefinition](../../yaml-crd-ref/evaluationdefinition.md)
resources.

### Remediation

KLT does not currently support the same level of
[remediations](https://keptn.sh/docs/1.0.x/reference/files/remediation/)
as Keptn v1 does,
but it does provide limited "Day 2" facilities:

* Any query that is possible for your data provider post-deployment
  can be defined as a `KeptnMetricDefinition`;
  this is then reported as a Keptn Metric.
  Evaluation can be defined as a `KeptnEvaluationDefinition`.
* `KeptnMetricsDefinition` resources can be retrieved and used
  to implement the Kubernetes HorizontalPodAutoscaler (HPA),
  which can detect the need for additional resources
  (more pods, memory, disk space, etc.)
  and automatically add those resources to your configuration
  based on the `ReplicaSet` resources you have defined.
  See
  [Using the HorizontalPodAutoscaler](../../implementing/evaluatemetrics.md/#using-the-horizontalpodautoscaler)
  for more information.

### Integrations and services in JES

Most functionality coded using the Keptn v1
[JES](https://github.com/keptn-contrib/job-executor-service)
(Job Executor Service) facility
can simply be moved into a `KeptnTaskDefinition` resource
that uses the
[container-runtime runner](../../yaml-crd-ref/taskdefinition.md/#yaml-synopsis-for-container-runtime).
If the JES container code is written in JavaScript or TypeScript,
you may be able to use the `deno-runtime` runner.
If the JES container code is written in Python 3,
you may be able to use the `python-runtime` runner.

Note that there is no need for integrations for data providers in KLT;
these are configured as `KeptnMetricsProvider` resources.
