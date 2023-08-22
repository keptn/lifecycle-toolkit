---
title: Manage release lifecycle
description: How KLT orchestrates pre- and post-deployment evaluations and tasks
weight: 20
---

In this exercise, we will configure the Keptn Lifecycle Toolkit
to run deployment checks as part of your deployment.
Whether you are deploying your software with
Argo, Flux, another deployment engine, or even `kubectl apply`,
the Lifecycle Toolkit can do the following:

* Pre-deploy: Validate external dependencies,
  confirm that images are scanned, and so forth

* Post-deploy: Execute tests, notify stakeholders,
  promote to the next stage

* Automatically validate against your SLO (Service Level Objectives)

KLT sits on top of the Kubernetes scheduler
and can do the following:

* Trace the deployment from start to end
* KLT is application aware,
  so can extend the deployment with tasks and evaluations that
  are run either before or after your whole application starts the deployment
  or at the individual workload level.
* Validate any Keptn metric,
  either pre- or post-deployment,
  using the metrics from the Keptn Metrics Server introduced in
  [Getting started with Keptn metrics](usecase_metrics.md).

This means that you can be sure that the environment is healthy
and has adequate resources before you begin the deployment.
After the deployment succeeds,
use Keptn metrics to confirm that your deployed software is really healthy --
not just that the pods are running but validate against SLOs
such as performance and user experience.
You can also check for new logs that came in from a log monitoring solution.

## Using this exercise

This exercise shows how to implement
pre- and post-deployment evaluations and tasks
for your application.
It is based on the
[simplenode-dev](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd)
example.

The steps to implement pre- and post-deployment orchestration are:

1. [Define evaluations to be performed pre- and post-deployment](#define-evaluations-to-be-performed-pre--and-post-deployment)
1. [Define tasks to be performed pre- and post-deployment](#define-tasks-to-be-performed-pre--and-post-deployment)
1. [Integrate evaluations and tasks into the cluster](#integrate-evaluations-and-tasks-into-the-cluster)

This is the third of three exercises in the
[Introducing the Keptn Lifecycle Toolkit](_index.md)
series.
You may want to complete the other exercises before doing this exercise
although that is not required:

* In the
  [Getting started with Keptn metrics](usecase_metrics.md)
  exercise, you learn how to define and use Keptn metrics.
* In [Standardize observability](usecase-observability.md),
  you learn how to standardize access
  to the observability data for your cluster.

If you are installing the Keptn Lifecycle Toolkit on an existing cluster
or in a local cluster you are creating for this exercise
and did not previously set up your cluster for the
[Standardize observability](usecase-observability.md) exercise,
you need to do the following:

1. Follow the instructions in
   [Install and update](../install/_index.md)
   to install and enable KLT on your cluster.
1. Follow the instructions in
   [Annotate workload](../implementing/integrate/#basic-annotations)
   to integrate the Lifecycle Toolkit into your Kubernetes cluster
   by applying basic annotations to your `Deployment` resource.
   This also creates appropriate
   [KeptnApp](../yaml-crd-ref/app.md) resources
   which aggregate workloads that are combined into the released product,
   regardless of the tools being used.

## Define evaluations to be performed pre- and post-deployment

An `evaluation` is a KeptnMetric that has a defined target value.
Evaluations are resources that are defined in a
[KeptnEvaluationDefinition](../yaml-crd-ref/evaluationdefinition.md)
yaml file.
In our example, evaluations are defined in the
[keptn-evaluations.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/keptn-evaluations.yaml)
file.
For example, the definition of the `evaluate-dependencies` evaluation
looks like this:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnEvaluationDefinition
metadata:
  name: evaluate-dependencies
  namespace: simplenode-dev
spec:
  objectives:
    - keptnMetricRef:
        name: available-cpus
        namespace: simplenode-dev
      evaluationTarget: ">4"
```

You see that the `available-cpus` metric is defined in the
[keptn-metric.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/keptn-metric.yaml)
file.
The `evaluationTarget` is set to be `>4`,
so this evaluation makes sure that more than 4 CPUs are available.
You could include objectives and additional metrics in this evaluation.

## Define tasks to be performed pre- and post-deployment

Tasks are resources that are defined in a
[KeptnTaskDefinition](../yaml-crd-ref/taskdefinition.md)
file.
In our example, the tasks are defined in the
[keptn-tasks.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/keptn-tasks.yaml)
file.
As an example,
we have a `notify` task that composes some Markdown text
to be sent as Slack notifications
The `KeptnTaskDefinition` looks like this:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition
metadata:
  name: notify
spec:
  function:
    inline:
      code: | 
            <javascript code>
    secureParameters:
      secret: slack-notification
```

For more information about sending Slack notifications with KLT, see
[Implement Slack notifications](../implementing/slack.md).
The code to be executed is expressed as a
[Deno](https://deno.land/)
script, which uses JavaScript syntax.
It can be embedded in the definition file
or pulled in from a remote webserver that is specified.
For this example, the code to be executed is embedded in this file
although, in practice,
this script would probably be located on a remote webserver.

Note that, beginning with KLT 0.8.0,
you can also use Python 3 to define your task,
or you can define a standard Kubernetes container
that uses the image, runner, and runtime dependencies that you choose.
For more information, see
[Working with Keptn tasks](../implementing/tasks).

You can view the actual JavaScript code for the task in the repository.
You see that "context" is important in this code.
This refers to the context in which this code executes --
for which application, for which version, for which Workload.

Because the slack server that is required to execute this task
is protected by a secret, the task definition also specifies that secret.

## Integrate evaluations and tasks into the cluster

Follow the instructions in
[Annotate workload](../implementing/integrate/#pre--and-post-deployment-checks)
to integrate the evaluations and tasks you defined
into the cluster
by applying annotations to the `Deployment` resource.
See the
[simplenode-dev-deployment.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/simplenode-dev-deployment.yaml)
file for an example.
