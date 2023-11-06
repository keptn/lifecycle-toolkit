---
title: Introduction to Keptn
description: An introduction to Keptn and the usecases.
weight: 10
---

Keptn integrates seamlessly with cloud-native deployment tools
such as ArgoCD, Flux, and Gitlab
to bring application awareness to your Kubernetes cluster.
Keptn supplements the standard deployment tools
with features to help you ensure that your deployments are in
a healthy state.

For information about the history of the Keptn project, see the
[Keptn Lifecycle Toolkit is now Keptn!](https://medium.com/keptn/keptn-lifecycle-toolkit-is-now-keptn-e0812217bf46)
blog.

Keptn includes multiple features
that can be implemented independently or together.
It targets three main use cases:
Metrics, Observability, and Release lifecycle management.

## Metrics

The Keptn metrics feature extends the functionality of
[Kubernetes metrics](https://kubernetes.io/docs/concepts/cluster-administration/system-metrics/):

* Allows you to define metrics
  from multiple data sources in your Kubernetes cluster.

* Supports deployment tools like Argo, Flux, KEDA, HPA, or
  Keptn for automated decision-making based on observability data.

* Handles observability data from multiple instances
  of multiple observability solutions
  – Prometheus, Dynatrace, Datadog and others –
  as well as data that comes directly from your cloud provider
  such as AWS, Google, or Azure.

* Enhances the Kubernetes
  [Horizontal Pod Autoscaling](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
  facility.

The Keptn metrics server unifies and standardizes
access to data from various sources,
simplifying configuration and integration into a single set of metrics.

To learn more, see:

* [Getting started with Keptn metrics](../getting-started/metrics.md)
* [Keptn Metrics](../implementing/evaluatemetrics.md) User Guide

## Observability

Keptn ensures observability for Kubernetes deployments
by creating a comprehensive trace
of all Kubernetes activities within a deployment.
Keptn observability makes it easy to understand
deployment durations and failures across multiple deployment strategies.

* Provides observability data for standard Kubernetes [workload](https://kubernetes.io/docs/concepts/workloads/) resources
  as well as
  [KeptnApp](https://lifecycle.keptn.sh/docs/yaml-crd-ref/app/)
  resources (which connect logically related [workloads](https://kubernetes.io/docs/concepts/workloads/))
  using different deployment strategies.

* Captures
  [DORA metrics](../implementing/dora.md)
  and exposes them as OpenTelemetry metrics out of the box.

* Reports traces and custom Keptn metrics from configured data providers
   using OpenTelemetry.

* Enables monitoring of new logs from log monitoring solutions.

* Information can be displayed on standard dashboard tools
  like Grafana.

Keptn is tool- and vendor neutral
and does not depend on particular tooling.
Keptn emits signals at every stage
([Kubernetes events](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/event-v1/),
[CloudEvents](https://cloudevents.io/), and
OpenTelemetry metrics and traces)
to ensure that your deployments are observable.

To learn more, see:

* [Getting started with Keptn Observability](../getting-started/observability.md)
* [Standardize observability](usecase-observability.md/)
* [DORA metrics](../implementing/dora.md) User Guide
* [OpenTelemetry observability](../implementing/otel.md) User Guide

## Release lifecycle management

The Release lifecycle management tools run in conjunction
with the standard Kubernetes deployment tools
to make deployments more robust.
Keptn "wraps" a standard Kubernetes deployment
with the capability to automatically handle issues
before and after the actual deployment.

These tools run checks and tasks before or after deployment initiation.

* Pre-deployment tasks such as checking for dependent services,
  image scanning, and setting the cluster to be ready for the deployment.

* Pre-deployment evaluations such as checking whether the cluster
  has enough resources for the deployment.

* Post-deployment tasks such as triggering tests,
  triggering a deployment to another cluster,
  or sending notifications that the deployment succeeded or failed.

* Post-deployment evaluations to evaluate the deployment,
  evaluate test results,
  or confirm software health against SLOs
  like performance and user experience.

All `KeptnTask` resources that are defined by `KeptnTaskDefinition` resources
at the same level (either pre-deployment or post-deployment) run in parallel.
Task sequences that are not part of the lifecycle workflow
should be handled by the pipeline engine tools rather than Keptn.
A `KeptnTask` resource can be defined to run multiple executables
(functions, programs, and scripts)
that are part of the lifecycle workflow.
The executables within a `KeptnTask` resource
run in sequential order.

Keptn tasks and evaluations can be run
for either a Kubernetes [workload](https://kubernetes.io/docs/concepts/workloads/) (single service) resource
or a
[KeptnApp](../yaml-crd-ref/app.md) resource,
which is a single, cohesive unit that groups multiple [workloads](https://kubernetes.io/docs/concepts/workloads/).
For more information, see:

* [Getting started with release lifecycle management](../getting-started/lifecycle-management.md)
* [Deployment tasks](../implementing/tasks.md) User Guide
* [Evaluations](../implementing/evaluations.md) User Guide
* [Manage release lifecycle](usecase-orchestrate.md)
* [KeptnApp and KeptnWorkload resources](../architecture/keptn-apps.md)
