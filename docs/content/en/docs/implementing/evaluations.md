---
title: Evaluations
description: Understand Keptn evaluations and how to use them
weight: 150
---
A `KeptnEvaluation` checks whether a `KeptnMetric` meets a defined target value.
A `KeptnEvaluation` resource is defined in a
[KeptnEvaluationDefinition](../yaml-crd-ref/evaluationdefinition.md)
yaml file.
The example
[app-pre-deploy-eval.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/version-3/app-pre-deploy-eval.yaml)
file specifies the `app-pre-deploy-eval-2` evaluation as follows:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnEvaluationDefinition
metadata:
  name: app-pre-deploy-eval-2
  namespace: podtato-kubectl
spec:
  objectives:
    - keptnMetricRef:
        name: available-cpus
        namespace: podtato-kubectl
      evaluationTarget: ">4"`
```

The `evaluationTarget` is set to be `>4`,
so this evaluation ensures that more than 4 CPUs are available
before the workload or application is deployed.

The example
[metric.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/base/metric.yaml)
file defines the
[KeptnMetric](../yaml-crd-ref/metric.md) resource
that is named  `available-cpus`:

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: KeptnMetric
metadata:
  name: available-cpus
  namespace: podtato-kubectl
spec:
  provider:
    name: my-provider
  query: "sum(kube_node_status_capacity{resource='cpu'})"
  fetchIntervalSeconds: 10
```

To integrate your `KeptnEvaluation` resource, you must:

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
