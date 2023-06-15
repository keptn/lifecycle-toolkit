---
title: Working with container runtime
description: How to run tasks in container runtime
weight: 95
---

Container runtime allows the Keptn Lifecycle Toolkit
to execute tasks in Kubernetes containers as part of a Kubernetes Job.
This is similar to the Keptn v1 Job Executor Service.

This is a new feature introduced in v0.8.0
that is useful for:

- Running load/performance testing tools
- Running tasks expressed in other languages such as Python
- [TODO] What else?
- Running SLO validations

To implement this feature:

- Define a
  [KeptnTaskDefinition](../yaml-crd-ref/taskdefinition.md)
  resource that defines the container
- Populate a [KeptnApp](../yaml-crd-ref/app.md)
  resource that associates that KeptnTaskDefinition
  with the pre- and post-deployment tasks
  that should run in it

## Using containers for SLO validation

One special use case for container runtimes
is to implement SLO (Service-Level Objectives) validation coming from
OpenTelementry workload traces.
Basically, you define OpenTelemetry traces
as SLIs (Service-Level Indicators)
that can be queried from the OpenTelemetry backend using tools such as
[Grafana Tempo](https://grafana.com/oss/tempo/).
