---
title: Workloads
description: Learn what Keptn Workloads are and how to use them
icon: concepts
layout: quickstart
weight: 10
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

A `KeptnWorkload`resource contains information about
which tasks should be performed during the `preDeployment`
or `postDeployment` phase of a deployment.
In its state,
it keeps track of the currently active `Workload Instances`,
which are responsible for doing those checks
for a particular instance of a Deployment/StatefulSet/ReplicaSet
(e.g. a Deployment of a certain version).

## KeptnWorkload

A `KeptnWorkload` resource augments a Kubernetes
[Workload](https://kubernetes.io/docs/concepts/workloads/)
with the ability to handle extra phases.
KLT generates the `KeptnWorkload` resource
from metadata information;
it is not necessary to manually create a YAML file that defines it.

A `KeptnWorkload` instance is responsible for executing
the pre- and post deployment checks of a workload.
In its state, it keeps track of the current status of all checks,
as well as the overall state of the Pre Deployment phase,
which the scheduler can use to determine whether the deployment should proceed and whether the pod should be bound to a node or stay pending.
that a node meets the criteria required
to allow a pod to be bind to a node.
`KeptnWorkload` instances have a reference
to the respective Pod/DeamonSet/StatefulSet/ReplicaSet,
to check whether it has reached the desired state.
If it detects that the referenced object has reached its desired state
(e.g. all pods of a deployment are up and running),
it knows that a `PostDeploymentCheck` can be triggered.
