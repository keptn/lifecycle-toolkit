---
title: Workloads
description: Learn what Keptn Workloads are and how to use them
icon: concepts
layout: quickstart
weight: 10
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

A Workload contains information about which tasks should be performed during the `preDeployment` as well as
the `postDeployment`
phase of a deployment. In its state it keeps track of the currently active `Workload Instances`, which are responsible
for doing those checks for
a particular instance of a Deployment/StatefulSet/ReplicaSet (e.g. a Deployment of a certain version).

### Keptn Workload Instance

A Workload Instance is responsible for executing the pre- and post deployment checks of a workload. In its state, it
keeps track of the current status of all checks, as well as the overall state of
the Pre Deployment phase, which can be used by the scheduler to tell that a pod can be allowed to be placed on a node.
Workload Instances have a reference to the respective Deployment/StatefulSet/ReplicaSet, to check if it has reached the
desired state. If it detects that the referenced object has reached
its desired state (e.g. all pods of a deployment are up and running), it will be able to tell that
a `PostDeploymentCheck` can be triggered.
