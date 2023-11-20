---
title: Evaluations
description: Understand Keptn evaluations and how to use them
weight: 150
---

A
[KeptnEvaluationDefinition](../yaml-crd-ref/evaluationdefinition.md)
resource contains a list of `objectives`,
each of which checks whether a defined `KeptnMetric` resource
meets a defined target value.
The example
[app-pre-deploy-eval.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/version-3/app-pre-deploy-eval.yaml)
file specifies the `app-pre-deploy-eval-2` evaluation as follows:
{{< embed path="/examples/sample-app/version-3/app-pre-deploy-eval.yaml" >}}

The `evaluationTarget` is set to be `>1`,
so this evaluation ensures that more than 1 CPU is available
before the [workload](https://kubernetes.io/docs/concepts/workloads/) or application is deployed.

This evaluation references the
[KeptnMetric](../yaml-crd-ref/metric.md) resource
that is named  `available-cpus`.
This is defined in the example
[metric.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/base/metric.yaml)
file:
{{< embed path="/examples/sample-app/base/metric.yaml" >}}

To run an evaluation on one of your
[Workloads](https://kubernetes.io/docs/concepts/workloads/)
([Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
[DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
or
[ReplicaSets](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/),
you must:

* Annotate your [workloads](https://kubernetes.io/docs/concepts/workloads/)
  to identify the `KeptnEvaluationDefinition` resource you want to run
  pre- and post-deployment for the specific [workloads](https://kubernetes.io/docs/concepts/workloads/).
* Manually edit all
  [KeptnApp](../yaml-crd-ref/app.md) resources
  to specify the `KeptnEvaluationDefinition` to be run
  pre- and post-deployment evaluations for the `KeptnApp` itself.

See
[Pre- and post-deployment checks](./integrate.md#pre--and-post-deployment-checks)
for details.

Note the following:

* One `KeptnEvaluationDefinition` resource can include
  multiple `objective` fields that reference additional metrics.
  In this example, you might want to also query
  available memory, disk space, and other resources
  that are required for the deployment.
* The `KeptnMetric` resources that are referenced
  in a `KeptnEvaluationDefinition` resource
  * can be defined on different namespaces in the cluster
  * can query different instances of different types of metric providers
* All objectives within a `KeptnEvaluationDefinition` resource
  are evaluated in order.
  If the evaluation of any objective fails,
  the `KeptnEvaluation` itself fails.
* You can define multiple evaluations
  for each stage (pre- and post-deployment).
  These evaluations run in parallel so the failure of one evaluation
  has no effect on whether other evaluations are completed.
* The results of each evaluation
  is written to a
  [KeptnEvaluation](../crd-ref/lifecycle/v1alpha3/#keptnevaluation)
  resource.
