---
title: Keptn Lifecycle Operator
linktitle: Lifecycle Operator
description: Basic understanding of the Keptn Lifecycle Operator
weight: 80
---


**Keptn's Lifecycle Operator** is
a Kubernetes [operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
that automates the deployment and management
of the Keptn components in a Kubernetes cluster.
The Keptn Lifecycle Operator contains several controllers for **Keptn CRDs**
and a **Mutating Webhook**.

Here's a brief overview:

**Keptn CRDs:** Keptn Lifecycle Operator contains
several controllers that manage and reconcile different types of Keptn CRDs
such as the Project Controller, Service Controller, and Stage Controller.

**Mutating Webhook:** automatically injects Keptn labels
and annotations into Kubernetes resources,
such as deployments and services.
These labels and annotations are used to enable Keptn's automation
and monitoring capabilities.
