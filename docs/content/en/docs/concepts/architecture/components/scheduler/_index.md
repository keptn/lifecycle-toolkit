---
title: Keptn Lifecycle Scheduler
linktitle: Scheduler
description: Basic understanding of Keptn's Lifecycle Scheduler
weight: 80
cascade:
---

**Keptn's Lifecycle Scheduler** replaces the
[Kubernetes scheduler](https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/)
to allow users to schedule events and tasks to occur
at specific times during the application lifecycle.
The Lifecycle Scheduler can trigger events such as
deployment, testing, and remediation at specific times or intervals.
The Keptn Scheduler guarantees that Pods are initiated only after
the Pre-Deployment checks are completed.
