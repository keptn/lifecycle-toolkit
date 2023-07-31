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

You control the content of a `KeptnApp` resource
with annotations or labels that are applied to each
[Workload](https://kubernetes.io/docs/concepts/workloads/)
([Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
[DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
and
[ReplicaSets](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/):

* The annotations described in
  [Basic annotations](../../../implementing/integrate/#basic-annotations)
  are used to automatically generate `KeptnApp` resources
  that contain the identifications required
  to run the KLT observability features.
* You must manually add the annotations described in
  [Pre- and post-deployment checks](../../../implementing/integrate/#pre--and-post-deployment-checks)
  to the basic `KeptnApp` manifest to define
  the evaluations and tasks you want to run pre- and post-deployment.

By default, the `KeptnApp` resources are updated every 30 seconds
when any of the Workloads have been modified;
The timeout is provided because it may take some time
to apply all `KeptnWorkload` resources to the cluster.
This interval can be modified for the cluster by changing the value
of the `keptnAppCreationRequestTimeoutSeconds` field in the
[KeptnConfig](../../../yaml-crd-ref/config.md)
resource.
