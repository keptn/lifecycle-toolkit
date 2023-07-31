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

The Keptn Lifecycle Toolkit uses a different paradigm
than Keptn v1 and so migration is not a straight-forward process.
Much of the logic and functionality of your Keptn v1 projects
can be migrated to KLT
but must be rewritten to utilize KLT components.

## How many namespaces?

You have significan flexibility to decide how many namespaces to use
and how to set them up.
See the Kubernetes
[Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
documentation for some basic information.
You can also search and find lots of "Best Practices for Namespaces"
documents published on the web.

Some considerations for KLT:

* KLT primarily operates on Kubernetes
  [Workload](https://kubernetes.io/docs/concepts/workloads/)
  resources and
  [KeptnApp](../../yaml-crd-ref/app.md)
   resources
  that are activated and defined by annotations to each Workload.
* [KeptnMetricsProvider](../../yaml-crd-ref/metricsprovider.md)
  resources need to be located
  in the same namespace as the associated
  [KeptnMetric](../../yaml-crd-ref/metric.md)
  resources.
  But
  [KeptnEvaluationDefinition](../../yaml-crd-ref/evaluationdefinition.md)
  resources that are used for pre- and post-deployment
  can reference metrics from any namespace.
  So you can create `KeptnMetrics` in a centralized namespace
  (such as `keptn-lifecycle-toolkit`)
  and access those metrics in evaluations on all namespaces in the cluster.
* Each `KeptnApp` resource identifies the namespace to which it belongs.
  If you configure multiple namespaces,
  you can have `KeptnApp` resources with the same name
  in multiple namespaces without having them conflict.
* You do not need separate namespaces for separate versions of your application.
  The `KeptnApp` resource includes fields to define
  the `version` as well as a `revision`
  (used if you have to rerun a deployment
  but want to retain the version number).

So, possible namespace designs run the gamut:

* Move all the Keptn v1 projects you had into a single namespace
* Move a group of Keptn v1 projects into a single namespace
* Define one namespace per Keptn v1 project
* Define one namespace per Keptn v1 stage

## Disposition of Keptn v1 components in KLT

To help you wrap your mind around the process,
this lists Keptn v1 components
and identifies their possible disposition when you migrate to KLT.

### project

A Keptn v1 project represents an arbitrary, higher-level grouping of services.
A project is defined by a
[shipyard.yaml](https://keptn.sh/docs/1.0.x/reference/files/shipyard/)
file.

KLT does not recognize projects.
Some Keptn v1 projects may translate into `KeptnApp` resources
but many will not.
A project has a 1:1 mapping to a Git repo,
but, since few applications are stored in a mono-repo,
it is difficult to implement a project-as-application paradigm

In general, the Keptn v1 project is a useful as a reference list
when migrating to KLT
but does not directly translate to a KLT resource.

### service

A Keptn v1 service models a smaller chunk of a project.
Most projects include many services
In a micro-services world,
a service may represent a micro-service
but it could instead be a wrapper for something else,
like "the entire public website"

KLT does not have the concept of a service.
When migrating to KLT,
you need to analyze what each service is doing
and translate that into an appropriate resource.
The closest analogy is a Kubernetes
[orkload](https://kubernetes.io/docs/concepts/workloads/)
but some services may become
[Keptn tasks](../../implementing/tasks/)
or other resources.

### stage

A stage is a subsection of a project
and has no corresponding KLT component.
However, the logic of the stages can be useful
when architecting the migration:

* **deployment stage** -- KLT does not itself deploy the software
    but rather depends on a deployment engine.
    However, a `deployment` stage may define sequences of tasks
    that should be translated into
    [KeptnTaskDefinition](../../yaml-crd-ref/taskdefinition.md)
    resources that are executed pre- and post-deployment
* **testing stage** may define sequences of tasks
    that should be translated into `KeptnTaskDefinition` resources
    that are executed pre- and post-deployment.

    Note that all `KeptnTask` resources at the same level
    (either pre-deployment or post-deployment)
    execute in parallel
    whereas Keptn v1 sequences and tasks could execute in parallel.
    If you have actions that need to execute sequentially,
    create a single `KeptnTaskDefinition` that calls each action in order.
    If you have tasks that can execute in parallel,
    migrating to KLT may improve the performance of the deployment.

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
tasks that are not part of the lifecycle workflow
should not be handled by KLT
but should instead be handled by the pipeline engine tools being used
such as Jenkins, Argo Workflows, Flux, and Tekton.

### task

Keptn v1 defines some specific types of tasks,
each of which is translated to KLT resources
that are appropriate for the activity:

* A **deployment task** becomes a
  [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
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

### SLIs

Keptn v1
[SLIs](https://keptn.sh/docs/1.0.x/reference/files/sli/)
(Service Level Indicators)
represent queries from the data provider
such as Prometheus, Dynatrace, or Datadog,
which is configured as a Keptn integration.

When migrating to KLT, you need to define a
[KeptnMetricsProvider](../../yaml-crd-ref/metricsprovider.md)
resource for the data provider you are using.

The queries defined for the Keptn v1 SLIs
should be translated into
[KeptnMetric](../../yaml-crd-ref/metric.md)
resources.
Note that KLT allows you to support multiple data providers
and multiple instances of each data provider for your SLI's
whereas Keptn v1 only allows you to use one SLI per project.

### SLOs

KLT at this time does not support the full range
of Quality Gates evaluations that are represented by
[SLOs](https://keptn.sh/docs/1.0.x/reference/files/slo/).
Facilities such as weighting of SLI's and scoring of the evaluation
do not currently exist.
However, simple evaluations of an SLI can be defined as
[KeptnEvaluationDefinition](../../yaml-crd-ref/evaluationdefinition.md)
resources.

### Remediation

KLT does not currently support the same level of
[remediations](https://keptn.sh/docs/1.0.x/reference/files/remediation/)
as Keptn v1 does.
KLT does provide limited "Day 2" facilities:

* Any query that is possible for your data provider post-deployment
  can be defined as a `KeptnMetricDefinition`;
  this is then reported as a Keptn Metric.
  Evaluation can be defined as a `KeptnEvaluationDefinition`.
* `KeptnMetricsDefinition` resources can be retrieved and used
  to implement the Kubernetes HorizontalPodAutoscaler (HPA),
  which can detect the need for additional resources
  (more pods, memory, disk space, etc)
  and automatically add those resources to your configuration
  based on the `ReplicaSet` resources you have defined

### Integrations and services in JES

Most functionality coded using the Keptn v1
[JES](https://github.com/keptn-contrib/job-executor-service)
(Job Executor Service) facility
can be defined as tasks using the `container-runtime` runner.
In most cases, the code from the JES container
can simply be moved into a `KeptnTaskDefinition` resource
that uses the `container-runtime` runner.
If the JES container code is written in JavaScript or TypeScript,
you may be able to use the `deno-runtime` runner.
If the JES container code is written in Python 3,
you may be able to use the `python-runtime` runner.

Note that there is no need for integrations for data providers in KLT;
these are configured as `KeptnMetricsProvider` resources.
