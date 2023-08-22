---
title: Keptn Lifecycle Scheduler
linktitle: Scheduler
description: Basic understanding of Keptn's Lifecycle Scheduler
weight: 80
cascade:
---

The **Keptn Scheduler** is an integral component of the Keptn Lifecycle Toolkit that orchestrates
the deployment process.
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
Additionally it registers itself as a [Permit plugin](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/#permit).

## How does the Keptn Scheduler works

Firstly the Mutating Webhook checks for annotations on Pods to see if it is annotated with
[Keptn specific annotations](https://main.lifecycle.keptn.sh/docs/implementing/integrate/#basic-annotations).
If the annotations are present, the Webhook assigns the **Keptn Scheduler** to the Pod.
This ensures that the Keptn Scheduler only gets Pods that have been annotated for it.

If the Pod is annotated with Keptn specific annotations, the Keptn Scheduler retrieves
the WorkloadInstance CRD that is associated with the Pod.
The **WorkloadInstance CRD** contains information about the `pre-deployment` checks that
need to be performed before the Pod can be scheduled.

The Keptn Scheduler then checks the status of the WorkloadInstance CRD to see
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

The Keptn Scheduler processes the following information from the WorkloadInstance CRD:

- The name of the pre-deployment checks that need to be performed.
- The status of the pre-deployment checks.
- The deadline for the pre-deployment checks to be completed.
- The Keptn Scheduler checks the status of the `pre-deployment` checks every 10 seconds.
If the checks have not finished successfully within 5 minutes, the Keptn Scheduler will not allow the Pod to be scheduled.

If all of the `pre-deployment` checks have finished successfully and the deadline has not been reached,
the Keptn Scheduler allows the Pod to be scheduled.
If any of the `pre-deployment` checks have not finished successfully or the deadline has
been reached, the Keptn Scheduler tells Kubernetes to check again later.

Also the Keptn Scheduler will not schedule Pods to nodes that have failed `pre-deployment`
checks in the past.
This helps to prevent Pods from being scheduled to nodes that are not ready for them.
