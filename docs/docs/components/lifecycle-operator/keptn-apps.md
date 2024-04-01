---
comments: true
---

# KeptnApp and KeptnWorkload resources

## Keptn Workloads

A
[KeptnWorkload](../../reference/api-reference/lifecycle/v1/index.md#keptnworkload)
resource augments a Kubernetes
[Workload](https://kubernetes.io/docs/concepts/workloads/)
with the ability to handle extra phases.
It can execute the pre-/post-deployment evaluations of a workload
and run pre-/post-deployment tasks.

In its state, it tracks the currently active workloads
(`DaemonSet`, `StatefulSet`, or `ReplicaSet` resources),
as well as the overall state of the Pre Deployment phase,
which Keptn can use to determine
whether the pods belonging to a workload
should be created and assigned to a node.
When it detects that the referenced object has reached its desired state
(e.g. all pods of a deployment are up and running),
it knows that a`PostDeploymentCheck` can be triggered.

The `KeptnWorkload` resources are created automatically
and without delay by the mutating webhook
as soon as the workload manifest is applied.

> **Note**
By default Keptn observes the state of the Kubernetes workloads
for 5 minutes.
After this timeout is exceeded, the deployment phase (from Keptn
viewpoint) is considered as `Failed` and Keptn does not proceed
with post-deployment phases (tasks, evaluations or promotion phase).
This timeout can be modified for the cluster by changing the value
of the `observabilityTimeout` field in the
[KeptnConfig](../../reference/crd-reference/config.md)
resource.

## Keptn Applications

A [KeptnApp](../../reference/crd-reference/app.md)
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
  before the scheduler binds the pods to the nodes.
  For information about how to disable the blocking
  functionality, please refer to
  [this section](./keptn-non-blocking.md#keptn-non-blocking-deployment-functionality).
* You can define post-deployment evaluations and tasks
  that run only after all the workloads have completed successfully.
* You can define promotion tasks that run only after all the post-deployment
  tasks and evaluations have completed successfully.

You control the content of a `KeptnApp` resource
with annotations or labels that are applied to each
[Workload](https://kubernetes.io/docs/concepts/workloads/)
([Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
and
[DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/))
plus specific tasks and evaluations that you define
for the `KeptnApp` resource itself:

* The annotations described in
  [Basic annotations](../../guides/integrate.md#basic-annotations)
  are used to automatically generate `KeptnApp` resources
  that contain the identifications required
  to run the Keptn observability features.
* You must create a `KeptnAppContext` resource
  that has the same name and namespace of your `KeptnApp` to define
  the evaluations and tasks you want to run pre-/post-deployment.
  For more information check [how to create tasks](../../guides/tasks.md/#run-a-task-associated-with-your-entire-keptnapp)
  and [how to create evaluations](../../guides/evaluations.md/#create-keptnappcontext-for-app-level-evaluations).

The `KeptnApp` resources are generated automatically and
contain the identifications required to run the Keptn observability features.
The `spec.workloads.name` and a `spec.workloads.version` fields
are also generated automatically as long as all the workloads you want in your app
are correctly annotated. (See [basic annotations](#how-basic-annotations-are-implemented))

By default, the `KeptnApp` resources are updated every 30 seconds
when any of the Workloads have been modified;
The timeout is provided because it may take some time
to apply all `KeptnWorkload` resources to the cluster.
This interval can be modified for the cluster by changing the value
of the `keptnAppCreationRequestTimeoutSeconds` field in the
[KeptnConfig](../../reference/crd-reference/config.md)
resource.

## How basic annotations are implemented

The [Basic annotations](../../guides/integrate.md#basic-annotations)
page gives instructions for applying the annotations or labels
that identify the pods that Keptn should manage.

Three `keptn.sh` and three `app.kubernetes.io` keys are recognized.
They are equivalent; you can use either of them
and they can be implemented as either annotations or labels.
Annotations take precedence over labels,
and the `keptn.sh` keys take precedence over `app.kubernetes.io` keys.
In other words:

* The operator first checks if the `keptn.sh/*` key is present
  in the annotations, and then in the labels.
* If neither is the case, it looks for the `app.kubernetes.io/*` equivalent,
  again first in the annotations, then in the labels.

Keptn automatically generates appropriate
[KeptnApp](../../reference/crd-reference/app.md)
resources that are used for observability,
based on whether the `keptn.sh/app` or `app.kubernetes.io/part-of`
annotation/label is populated:

* If either of these labels/annotations are populated,
  Keptn automatically generates a `KeptnApp` resource
  that includes all workloads that have the same annotation/label,
  thus creating a `KeptnApp` resource for each defined grouping

* If only the `workload` and `version` annotations/labels are available
  (in other words, neither the `keptn.sh/app`
  or `app.kubernetes.io/part-of` annotation/label is populated),
  one `KeptnApp` resource is created automatically for each workload.
