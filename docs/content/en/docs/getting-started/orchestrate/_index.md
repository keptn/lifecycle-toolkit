---
title: Orchestrating pre- and post-deployment tasks and evaluations
description: Learn how the Keptn Lifecycle Toolkit can orchestrate deployment checks.
weight: 55
---

In this exercise, we will configure the Keptn Lifecyle Toolkit
to run deployment checks as part of your deployment.
Whether you are deploying your software with
Argo, Flux, another deployment engine, or even `kubectl apply`,
the Lifecycle Toolkit can do the following:

* Pre-deploy: Validate external dependencies,
  confirm that images are scanned, and so forth

* Post-deply: Execute tests, notify stakeholders,
  promote to the next stage

* Automatically validate against your SLO (Service Level Objectives)


KLT sits in the job scheduler and can trace the deployment
from start to end.
KLT is also application aware,
so we can extend the deployment
with tasks and evaluations that
are run either before or after the actual deployment.
You can also validate any metric,
either pre- or post-deployment,
using the metrics from the Keptn Metrics Server introduced in
[Getting started with Keptn metrics](../metrics).
This means that you can be sure that the environment is healthy
and has adequate resources before you begin the deployment.
After the deployment succeeds,
use Keptn metrics to confirm that your deployed software is really healthy --
not just that the pods are running but validate against SLOs
such as SLOs to measure performance and user experience.
You can also check for new logs that came in from a log monitoring solution.

## Using this exercise

This exercise is based on the
[simplenode-dev](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd)
example.
You can clone that repo to access it locally
or just look at it for examples
as you implement the functionality "from scratch"
on your local Kubernetes deployment cluster.

The steps to implement pre- and post-deployment orchestration are:

1. Bring or create a Kubernetes deployment cluster
1. Install the Keptn Lifecycle Toolkit on your cluster
1. Enable KLT for your cluster
1. Integrate KLT with your cluster by annotating
   the Kubernetes `Deployment` and `Namespace` CRDs
1. Define evaluations to be performed pre- and post-deployment
1. Define tasks to be performed pre- and post-deployment

## Bring or create a Kubernetes deployment cluster

## Install KLT on your cluster

## Enable KLT for your cluster

To enable KLT for your cluster, annotate the
[Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
CRD.
In this example, this is defined in the
[simplenode-dev-ns.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/simplenode-dev-ns.yaml)
file, which looks like this:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: simplenode-dev
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled"
```

You see the annotation line that enables `lifecycle-toolkit`.
This line tells the webhook to handle the namespace

## Integrate KLT with your cluster

To integrate KLT with your cluster, annotate the Kubernetes
[Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
CRD.
In this example, this is defined in the
[simplenode-dev-deployment.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/simplenode-dev-deployment.yaml)
file, which includes the following lines:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simplenode
  namespace: simplenode-dev
...
template:
    metadata:
      labels:
        app: simplenode
        app.kubernetes.io/name: simplenodeservice
      annotations:
        # keptn.sh/app: simpleapp
        keptn.sh/workload: simplenode
        keptn.sh/version: 1.0.2
        keptn.sh/pre-deployment-evaluations: evaluate-dependencies
        keptn.sh/pre-deployment-tasks: notify
        keptn.sh/post-deployment-evaluations: evaluate-deployment
        keptn.sh/post-deployment-tasks: notify
...
```

For more information about using annotations and labels
to integrate KLT into your deployment cluster, see
[Integrate KLT with your applications[../../implementing/integrate).

## Define evaluations to be performed pre- and post-deployment

An `evaluation` is a KeptnMetric that has a defined target value.
Evaluations are CRDs that are defined in a
[KeptinEvaluationDefinition](../../yaml-crd-ref/evaluationdefinition.md)
yaml file.
For our example, evaluations are defined in the
[keptn-evaluations.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/keptn-evaluations.yaml) file.
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

Tasks are CRDs that are defined in a
[KeptnTaskDefinition](../../yaml-crd-ref/taskdefinition.md)
file.
For our example, the tasks are defined in the
[keptn-tasks.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/keptn-tasks.yaml)
file
As an example,
we have a `notify` task that composes some markdown text
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
[Implement Slack notifications](../../implementing/slack.md).
The code to be executed is expressed as a
[Deno](https://deno.land/)
script, which uses JavaScript syntax.
It can be embedded in the definition file
or pulled in from a remote webserver that is specified.
For this example, the code to be executed is embedded in this file
although, in practice,
this script would probably be located on a remote webserver.

You can view the actual JavaScript code for the task in the repository.
You see that "context" is important in this code.
This refers to the context in which this code executes --
for which application, for which version, for which Workload

Because the slack server that is required to execute this task
is protected by a secret, the task definition also specifies that secret.
