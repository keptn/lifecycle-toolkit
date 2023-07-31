---
title: KeptnApp and KeptnWorkload resources
linktitle: Keptn Applications and Keptn Workloads
description: How Keptn applications work
weight: 50
cascade:
---

## Keptn Workloads

A
[KeptnWorkload](../../../crd-ref/lifecycle/v1alpha3/#keptnworkload)
resource augments a Kubernetes
[Workload](https://kubernetes.io/docs/concepts/workloads/)
with the ability to handle extra phases.
It can execute the pre- and post-deployment evaluations of a Workload
and run pre- and post-deployment tasks.

In its state, it keeps track of the currently active `Workload Instances`,
(`Pod`, `DaemonSet`, `StatefulSet`, and `ReplicaSet` resources),
as well as the overall state of the Pre Deployment phase,
which the scheduler can use to determine whether the pods belonging to a workload should proceed.
When it detects that the referenced object has reached its desired state
(e.g. all pods of a deployment are up and running),
it knows that a`PostDeploymentCheck` can be triggered.

The KeptnWorkload resources are created automatically
by the mutating webhook as soon as a pod for the workload
(i.e. `Deployment`, `StatefulSet`, `DaemonSet`, `ReplicaSet`)
is about to be started.
The KeptnWorkloads are created automatically and without delay by the webhook.

## Keptn Applications

A [KeptnApp](../../../yaml-crd-ref/app.md)
resource combines multiple Kubernetes
[workloads](https://kubernetes.io/docs/concepts/workloads/)
into a single entity
that represent the application that is published.
Note that the Kubernetes documentation
often refers to workloads as applications,
but each workload actually corresponds to one version
of one deployable microservice,
not the amalgamation of multiple microservices
that typically comprise the released software.

Implementing Keptn applications provides the following benefits:

* Observability tools report on the deployment
  of all workloads together rather than individually.
* You can define pre-deployment evaluations and tasks
  that must all complete successfully
  before the scheduler creates the pods for any of the workloads.
* You can define post-deployment evaluations and tasks
  that run only after all the workloads have completed successfully.

The `KeptnApp` resources that are generated automatically
contain the identifications required to run the KLT observability features.
The `spec.workloads.name` and a `spec.workloads.version` fields
that define evaluations and tasks to be run
pre- and post-deployment are not generated automatically
but must be input manually.

By default, the `KeptnApp` resources are updated every 30 seconds
when any of the Workloads have been modified;
The timeout is provided because it may take some time
to apply all `KeptnWorkload` resources to the cluster.
This interval can be modified for the cluster by changing the value
of the `keptnAppCreationRequestTimeoutSeconds` field in the
[KeptnConfig](../../../yaml-crd-ref/config.md)
resource.
