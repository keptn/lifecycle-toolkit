---
title: Keptn integration with Scheduling
linktitle: Scheduler and Scheduling Gates
description: Basic understanding of how Keptn integrates with Kubernetes Pod Scheduling
weight: 80
---

Keptn integrates with Kubernetes scheduling to block
the deployment of applications that do not satisfy Keptn defined pre-deployment checks.
The default scheduling paradigm is different
depending on the version of Kubernetes you are running:

* On Kubernetes versions 1.26 and older,
  Keptn uses the **Keptn Scheduler** to block application deployment when appropriate
  and orchestrate the deployment process.

* On Kubernetes version 1.27 and greater,
  scheduling is implemented using
  [Kubernetes scheduling gates](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-scheduling-readiness/).

These two implementations are discussed below.

## Keptn Scheduling Gates for K8s 1.27 and above

When Keptn is running on Kubernetes version 1.27 and greater
and the Keptn Helm value `lifecycleOperator.schedulingGatesEnabled` is set to `true`,
Keptn uses the
[Pod Scheduling Readiness K8s API](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-scheduling-readiness)
to gate Pods until the required deployment checks pass.

When a workload is applied to a Kubernetes cluster,
the Mutating Webhook checks each Pod for annotations.
If
[Keptn specific annotations](../../implementing/integrate.md#basic-annotations)
are present,
the Webhook adds a scheduling gate to the Pod called `keptn-prechecks-gate`.
This spec tells the Kubernetes scheduling framework
to wait for the Keptn checks before binding the pod to a node.

For example, a pod gated by Keptn looks like the following:

{{< embed path="/docs/assets/scheduler-gates/gated.yaml" >}}

If the `pre-deployment` checks have finished successfully,
the WorkloadVersion Controller removes the gate from the Pod.
The default k8s scheduler can then allow the Pod to be bound to a node.
If the `pre-deployment` checks have not yet finished successfully,
the gate stays and the Pod remains in the pending state.
When removing the gate,
the WorkloadVersion controller also adds the following annotation so that,
if the spec is updated, the Pod is not gated again:

{{< embed path="/docs/assets/scheduler-gates/gate-removed.yaml" >}}

## Keptn Scheduler for K8s 1.26 and earlier

The **Keptn Scheduler** works by registering itself as a Permit plugin within the Kubernetes
scheduling cycle that ensures that Pods are scheduled to a node until and unless the
pre-deployment checks have finished successfully.
This helps to prevent Pods from being scheduled to nodes that are not yet ready for them,
which can lead to errors and downtime.
Furthermore, it also allows users to control the deployment of an application based on
customized rules that can take into consideration more parameters than what the default
scheduler has (typically CPU and memory values).

The Keptn Scheduler uses the Kubernetes
[Scheduler Framework](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/) and is based on the
[Scheduler Plugins Repository](https://github.com/kubernetes-sigs/scheduler-plugins/tree/master).
Additionally, it registers itself as
a [Permit plugin](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/#permit).

### How does the Keptn Scheduler works

Firstly the Mutating Webhook checks for annotations on Pods to see if it is annotated with
[Keptn specific annotations](../../implementing/integrate.md#basic-annotations).
If the annotations are present, the Webhook assigns the **Keptn Scheduler** to the Pod.
This ensures that the Keptn Scheduler only gets Pods that have been annotated for it.
A Pod `test-pod` modified by the Mutating Webhook looks as follows:

{{< embed path="/docs/assets/scheduler-gates/scheduler.yaml" >}}

If the Pod is annotated with Keptn specific annotations, the Keptn Scheduler retrieves
the WorkloadVersion CRD that is associated with the Pod.
The **WorkloadVersion CRD** contains information about the `pre-deployment` checks that
need to be performed before the Pod can be scheduled.

The Keptn Scheduler then checks the status of the WorkloadVersion CRD to see
if the `pre-deployment` checks have finished successfully.
If the pre-deployment checks have finished successfully, the **Keptn Scheduler** allows
the Pod to be scheduled to a node.
If the `pre-deployment` checks have not yet finished, the Keptn Scheduler tells Kubernetes to check again later.

It is important to note that the Keptn Scheduler is a plugin to the default Kubernetes scheduler.
This means that all of the checks that are performed by the default Kubernetes scheduler
will also be performed by the **Keptn Scheduler**.
For example, if there is not enough capacity on any node to schedule the Pod,
the Keptn Scheduler will not be able to schedule it, even if the `pre-deployment`
checks have finished successfully.

The Keptn Scheduler processes the following information from the WorkloadVersion CRD:

* The name of the pre-deployment checks that need to be performed.
* The status of the pre-deployment checks.
* The deadline for the pre-deployment checks to be completed.
* The Keptn Scheduler checks the status of the `pre-deployment` checks every 10 seconds.
If the checks have not finished successfully within 5 minutes,
the Keptn Scheduler does not allow the Pod to be scheduled.

If all of the `pre-deployment` checks have finished successfully and the deadline has not been reached,
the Keptn Scheduler allows the Pod to be scheduled.
If any of the `pre-deployment` checks have not finished successfully or the deadline has
been reached, the Keptn Scheduler tells Kubernetes to check again later.

Also the Keptn Scheduler will not schedule Pods to nodes that have failed `pre-deployment`
checks in the past.
This helps to prevent Pods from being scheduled to nodes that are not ready for them.

## Integrating Keptn with your custom scheduler

Keptn scheduling logics are compatible with
the [Scheduler Framework](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/).
Keptn does not work with a custom scheduler unless it is implemented as
a [scheduler plugin](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/#plugin-configuration).
