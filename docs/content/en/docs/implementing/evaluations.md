---
title: Evaluations
description: Understand Keptn evaluations and how to use them
weight: 150
---
A `KeptnEvaluation` checks whether a value in `KeptnMetric` meets a defined target value defined in `KeptnEvaluationDefinition`.
A `KeptnEvaluation` resource is defined in a
[KeptnEvaluationDefinition](../yaml-crd-ref/evaluationdefinition.md)
yaml file.
The example
[app-pre-deploy-eval.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/version-3/app-pre-deploy-eval.yaml)
file specifies the `app-pre-deploy-eval-2` evaluation as follows:
{{< embed path="/examples/sample-app/version-3/app-pre-deploy-eval.yaml" >}}

The `evaluationTarget` is set to be `>1`,
so this evaluation ensures that more than 1 CPU is available
before the workload or application is deployed.

The example
[metric.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/base/metric.yaml)
file defines the
[KeptnMetric](../yaml-crd-ref/metric.md) resource
that is named  `available-cpus`:
{{< embed path="/examples/sample-app/base/metric.yaml" >}}

To run an evaluation on one of your
[Workloads](https://kubernetes.io/docs/concepts/workloads/)
([Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
[DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
or
[ReplicaSets](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/),
you must:

* Annotate your [Workloads](https://kubernetes.io/docs/concepts/workloads/)
  [Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
  [StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
  and
  [DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/)
  to include each evaluation and task you want to run
  pre- and post-deployment for the specific workloads.
* Manually edit all
  [KeptnApp](../yaml-crd-ref/app.md) resources
  to specify the evaluations to be run
  pre- and post-deployment for the `KeptnApp` itself.

See [Pre- and post-deployment checks](../implementing/integrate/#pre--and-post-deployment-checks)
for details.

Note the following:

* One `KeptnEvaluationDefinition` resource can include
  multiple `objective` fields that reference additional metrics.
  In this example, you might want to also query
  available memory, disk space, and other resources
  that are required for the deployment.
* The `KeptnMetric` resources that are referenced
  in a `KeptnMetricDefinition` resource
  can be defined on different namespaces in the cluster.
* All `KeptnEvaluation` resources
  that are defined by `KeptnEvaluationDefinition` resources at the same level
  (either pre-deployment or post-deployment)
  execute in parallel.
* All objectives within a `KeptnEvaluationDefinition` resource
  are evaluated in order.
  If the evaluation of an objective fails,
  the `KeptnEvaluation` itself fails
  and the objectives listed after the failed objection
  are not evaluated.
