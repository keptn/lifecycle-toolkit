---
title: Introduction to Keptn
description: An introduction to Keptn and the usecases.
weight: 10
---

Keptn implements observability
for deployments and seamlessly integrates with deployment tools
such as ArgoCD, Flux, and Gitlab
and brings application awareness to your Kubernetes cluster.

These standard deployment tools
do an excellent job at deploying applications
but do not handle all issues
that are required to ensure that your deployment is usable.
Keptn "wraps" a standard Kubernetes deployment
with the capability to automatically handle issues
before and after the actual deployment.

Keptn includes multiple features
that can be implemented independently or together.
It targets three main use cases:
Custom metrics, Observability, and Release lifecycle management.

## Custom metrics

The Custom Keptn metrics feature extends the functionality of
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

The Keptn metrics server unifies and standardizes
access to data from various sources,
simplifying configuration and integration into a single set of metrics.

## Observability

Keptn ensures observability for Kubernetes deployments
by creating a comprehensive trace
of all Kubernetes activities within a deployment.
Keptn observability makes it easy to understand
deployment durations and failures across multiple deployment strategies.

* Provides observability data for standard Kubernetes workload resources
  as well as
  [KeptnApp](https://lifecycle.keptn.sh/docs/yaml-crd-ref/app/)
  resources (which connect logically related workloads)
  using different deployment strategies.

* Captures
  [DORA metrics](../implementing/dora/)
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

## Release lifecycle management

The Release lifecycle management tools run in conjunction
with the standard Kubernetes deployment tools
to make deployments more robust.

These tools run checks and tasks before or after deployment initiation.

* Pre-deployment tasks such as checking for dependent services,
  image scanning, and setting the cluster to be ready for the deployment

* Pre-deployment evaluations such as checking the layout of the cluster

* Post-deployment tasks such as triggering tests,
  triggering a deployment to another cluster,
  or sending notifications that the deployment succeeded or failed

* Post-deployment evaluations to evaluate the deployment,
  evaluate test results,
  or confirm software health against SLOs
  like performance and user experience

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
for either a Kubernetes workload (single service) resource
or a
[KeptnApp](https://lifecycle.keptn.sh/docs/yaml-crd-ref/app/) resource,
which is a single, cohesive unit that groups multiple workloads.

To familiarize yourself with how Keptn works, refer to the
[Getting started with Keptn](../getting-started/)
and the
[Getting Started Exercises](https://lifecycle.keptn.sh/docs/getting-started/).

For information about the history of the Keptn project,
see the blog.
