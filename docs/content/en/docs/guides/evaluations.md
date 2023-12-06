---
title: Evaluations
description: Understand Keptn evaluations and how to use them
weight: 700
---

An evaluation is run either pre- or post-deployment
to do a simple comparison of a metric value to a target value.
For example, an evaluation can check whether more than four CPUs are available,
whether a remote database is accessible,
or check for other problems in your infrastructure.
The deployment is kept in a pending state
until the infrastructure is capable of accepting deployments again.
You can also define
[deployment tasks](tasks.md)
that are run pre- and post-deployment

Use the
[Analysis](slo)
feature if you want to do more complex investigations
that may include multiple criteria with weights and scoring applied.

To implement evaluations, you must:

1. [Install and enable Keptn](../installation)
   in your cluster.
   You will need both the
   Keptn Metrics Operator (`metrics-operator`)
   and the Keptn Lifecycle Operator (`lifecycle-operator`)
   [components](../components).
1. Define a
   [KeptnMetricsProvider](../reference/crd-reference/metricsprovider.md)
   resource for each data source you want to use for evaluations.
   You can define multiple instances of multiple types of data providers.
1. Annotate all
   [workloads](https://kubernetes.io/docs/concepts/workloads/)
   ([Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
   [StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
   [DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
   and
   [ReplicaSets](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/)
   with
   [basic annotations](integrate/#basic-annotations).
1. Generate the required
   [KeptnApp](../reference/crd-reference/app.md)
   resources following the instructions in
   [Auto app discovery](auto-app-discovery.md).
1. Define a
   [KeptnMetric](../reference/crd-reference/metric.md)
   resource for each piece of data
   (defined as a query from one of your data sources)
   that you need for your evaluations.
1. [Create KetnEvaluationDefinition](#create-keptnevaluationdefinition-resources)
   resources for each evaluation you want to perform.
1. Annotate the appropriate `KeptnApp`
   resources for each `KeptnEvaluationDefinition` resource
   you want to run pre- and/or post-deployment

## Create KeptnEvaluationDefinition resources

A
[KeptnEvaluationDefinition](../reference/crd-reference/evaluationdefinition.md)
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
[KeptnMetric](../reference/crd-reference/metric.md) resource
that is named  `available-cpus`.
This is defined in the example
[metric.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/base/metric.yaml)
file:
{{< embed path="/examples/sample-app/base/metric.yaml" >}}

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
  [KeptnEvaluation](../reference/api-reference/lifecycle/v1alpha3/#keptnevaluation)
  resource.

## Annotate the KeptnApp resource

To define the pre- and post-deployment evaluations to run,
you must manually edit the
[KeptnApp](../reference/crd-reference/app.md)
YAML file to provide an annotation
for each `KeptnEvaluationDefinition` resource to be run
pre- and post-deployment.
The annotations are:

```yaml
keptn.sh/pre-deployment-evaluations: <evaluation-name>
keptn.sh/post-deployment-evaluations: <evaluation-name>
```

   > **Caveat:** Be very careful when implementing pre-deployment evaluations
     since, if one fails, Keptn prevents the deployment from running.
   >

The value of these annotations corresponds
to the values of the `name` field of each
[KeptnTaskDefinition](../reference/crd-reference/taskdefinition.md)
resource.
These resources contain re-usable "functions"
that can execute before and after the deployment.

If everything is fine, the deployment continues and afterward,
a Slack notification is sent with the result of the deployment

## Example of pre- and post-deployment actions

A comprehensive example of pre- and post-deployment
evaluations and tasks can be found in our
[examples folder](https://github.com/keptn/lifecycle-toolkit/tree/main/examples/sample-app),
where we use [Podtato-Head](https://github.com/podtato-head/podtato-head)
to run some simple pre-deployment checks.

To run the example, download the example
then issue the following commands:

```shell
cd ./examples/podtatohead-deployment/
kubectl apply -f .
```

Afterward, use the following command
to monitor the status of the deployment:

```shell
kubectl get keptnworkloadversion -n podtato-kubectl -w
```

The deployment for a workload stays in a `Pending` state
until the respective pre-deployment check is successfully completed.
Afterwards, the deployment starts and when the workload is deployed,
the post-deployment checks start.
