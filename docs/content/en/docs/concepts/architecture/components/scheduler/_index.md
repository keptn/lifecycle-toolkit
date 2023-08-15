---
title: Keptn Lifecycle Scheduler
linktitle: Scheduler
description: Basic understanding of Keptn's Lifecycle Scheduler
weight: 80
cascade:
---

The **Keptn Scheduler**, an integral component of the Keptn Lifecycle Toolkit which orchestrates
the deployment process.
The **Keptn Scheduler** works by registering itself as a Permit plugin with Kubernetes that
is responsible for ensuring that Pods are scheduled to a node until and unless the
pre-deployment checks have finished successfully.
This is important because it helps to prevent Pods from being scheduled to nodes that are
not yet ready for them, which can lead to errors and downtime

## How does the Keptn Scheduler works

The Keptn Scheduler first checks the annotations on the Pod to see if it is annotated with
[Keptn specific annotations](https://main.lifecycle.keptn.sh/docs/implementing/integrate/#basic-annotations).
If the Pod is annotated with Keptn specific annotations, the Keptn Scheduler retrieves
the WorkloadInstance CRD that is associated with the Pod.
The **WorkloadInstance CRD** contains information about the `pre-deployment` checks that
need to be performed before the Pod can be scheduled.

The Keptn Scheduler then checks the status of the WorkloadInstance CRD to see if the
`pre-deployment` checks have finished successfully.
If the pre-deployment checks have finished successfully, the **Keptn Scheduler** allows
the Pod to be scheduled to a node.
If the `pre-deployment` checks have not yet finished, the Keptn Scheduler tells Kubernetes to check again later.

The Keptn Scheduler processes the following information from the WorkloadInstance CRD:

- The name of the pre-deployment checks that need to be performed.
- The status of the pre-deployment checks.
- The deadline for the pre-deployment checks to be completed.

The Keptn Scheduler uses this information to determine whether or not the Pod can be scheduled.
If all of the pre-deployment checks have finished successfully and the deadline has not been reached,
the Keptn Scheduler allows the Pod to be scheduled.
If any of the `pre-deployment` checks have not finished successfully or the deadline has
been reached, the Keptn Scheduler tells Kubernetes to check again later.

The Keptn Scheduler also has the ability to blacklist nodes that have failed pre-deployment checks.
This means that the Keptn Scheduler will not schedule Pods to nodes that have failed pre-deployment
checks in the past.
This helps to prevent Pods from being scheduled to nodes that are not ready for them.

The **Keptn Scheduler** is a powerful tool that can be used to ensure that Pods are
only scheduled to nodes that are ready for them.
This helps to prevent errors and downtime, and it also helps to improve the overall reliability of your Kubernetes cluster.

This project uses the
Kubernetes [Scheduler Framework](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/)
and is based on the [Scheduler Plugins Repository](https://github.com/kubernetes-sigs/scheduler-plugins/tree/master)
