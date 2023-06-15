---
title: Working with container runtime
description: How to run tasks in container runtime
weight: 95
---

Container runtime allows the Keptn Lifecycle Toolkit
to execute tasks in Kubernetes
[Containers](https://kubernetes.io/docs/concepts/containers/)
as part of a Kubernetes
[Job](https://kubernetes.io/docs/concepts/workloads/controllers/job/).
This is similar to the Keptn v1 Job Executor Service.

This is a new feature introduced in v0.8.0
that is useful for:

- Running load/performance testing tools
- Running tasks expressed in other languages such as Python

To implement this feature:

- Define a
  [KeptnTaskDefinition](../yaml-crd-ref/taskdefinition.md)
  resource that defines the container
- Populate a [KeptnApp](../yaml-crd-ref/app.md)
  resource that associates that KeptnTaskDefinition
  with the pre- and post-deployment tasks
  that should run in it

