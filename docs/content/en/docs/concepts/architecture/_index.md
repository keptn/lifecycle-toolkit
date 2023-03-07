---
title: Architecture
linktitle: Architecture
description: Understand the details of how Keptn Lifecycle Toolkit works
weight: 80
cascade:
---


### Architecture

The Keptn Lifecycle Toolkit consists of two main components: 

* Keptn Lifecycle Operator
* Keptn Scheduler

The Keptn Lifecycle Operator contains several controllers for **Keptn CRDs** 
and a **Mutating Webhook**.
The Keptn Scheduler guarantees that Pods are initiated only after 
the Pre-Deployment checks are completed.

A **Kubernetes Manifest**, which is annotated with Keptn specific annotations, 
gets applied to the Kubernetes Cluster.
The Keptn Scheduler is then added through a Mutating Webhook, 
and events related to Pre-Deployment are sent to the event stream in Kubernetes. 
The Event Controller monitors these events and starts a Kubernetes Job to complete 
the Pre-Deployment process. Once the Pre-Deployment is done, the Keptn Scheduler 
schedules the pod to be deployed.

The **KeptnApp** and **KeptnWorkload** Controllers watch for the workload resources to finish 
and then generate a Post-Deployment Event.
After the Post-Deployment checks, SLOs can be validated by using an interface that retrieves 
SLI data from a provider, such as Prometheus [Prometheus](https://prometheus.io/).
Lastly, the Keptn Lifecycle Toolkit provides Metrics and Traces of the entire Deployment process 
with [OpenTelemetry](https://opentelemetry.io/).

<!-- ![KLT Architecture](./assets/architecture.png) -->

## How it works

<!-- The following sections will provide insights on each component of the Keptn Lifecycle Toolkit in terms of their purpose,
responsibility, and communication with other components.
Furthermore, there will be a description on what CRD they monitor and a general overview of their fields.

### Webhook

Annotating a namespace subjects it to the effects of the mutating webhook:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: podtato-kubectl
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled"  # this lines tells the webhook to handle the namespace
```

However, the mutating webhook will modify only resources in the annotated namespace that have Keptn annotations.
When the webhook receives a request for a new pod, it will look for the workload annotations:

```yaml
keptn.sh/workload: "some-workload-name"
```

The mutation consists in changing the scheduler used for the deployment with the Keptn Scheduler. Webhook then creates a
workload and app resource per annotated resource.
You can also specify a custom app definition with the annotation:

```yaml
keptn.sh/app: "your-app-name"
```

In this case the webhook will not generate an app, but it will expect that the user will provide one.
The webhook should be as fast as possible and should not create/change any resource.
Additionally, it will compute a version string, using a hash function that takes certain properties of the pod as
parameters
(e.g. the images of its containers).
Next, it will look for an existing instance of a `Workload CRD` for the given workload name:

- If it finds the `Workload`, it will update its version according to the previously computed version string.
  In addition, it will include a reference to the ReplicaSet UID of the pod (i.e. the Pods owner),
  or the pod itself, if it does not have an owner.
- If it does not find a workload instance, it will create one containing the previously computed version string.
  In addition, it will include a reference to the ReplicaSet UID of the pod (i.e. the Pods owner), or the pod itself, if
  it does not have an owner.

It will use the following annotations for
the specification of the pre/post deployment checks that should be executed for the `Workload`:

- `keptn.sh/pre-deployment-tasks: task1,task2`
- `keptn.sh/post-deployment-tasks: task1,task2`

and for the Evaluations:

- `keptn.sh/pre-deployment-evaluations: my-evaluation-definition`
- `keptn.sh/post-deployment-evaluations: my-eval-definition`

After either one of those actions has been taken, the webhook will set the scheduler of the pod and allow the pod to be
scheduled.

### Scheduler

After the Webhook mutation, the Keptn-Scheduler will handle the annotated resources. The scheduling flow follows the
default scheduler behavior,
since it implements a scheduler plugin based on
the [scheduling framework]( https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/).
For each pod, at the very end of the scheduling cycle, the plugin verifies whether the pre deployment checks have
terminated, by retrieving the current status of the WorkloadInstance. Only if that is successful, the pod is bound to a
node.
 -->
