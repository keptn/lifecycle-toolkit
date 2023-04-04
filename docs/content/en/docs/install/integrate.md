---
title: Integrate KLT with your applications
description: How to integrate the Keptn Lifecycle Toolkit into your Kubernetes cluster
icon: concepts
layout: quickstart
weight: 45
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

Use Kubernetes annotations and labels
to integrate the Keptn Lifecycle Toolkit into your Kubernetes cluster.

The Keptn Lifecycle Toolkit monitors manifests
that have been applied against the Kubernetes API
and reacts if it finds a workload with special annotations/labels.
This is a four-step process:

* Annotate your workload(s)
* Create a `KeptnApp` custom resource that references those workloads
* Create the `KeptnTaskDefinition`s you need
* Enable the target namespace by annotating it

## Annotate workload(s)

For this, you should annotate your
[Workload](https://kubernetes.io/docs/concepts/workloads/)
with (at least) the following annotations:

```yaml
keptn.sh/app: myAwesomeAppName
keptn.sh/workload: myAwesomeWorkload
keptn.sh/version: myAwesomeWorkloadVersion
```

Alternatively, you can use Kubernetes
[Recommended Labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/)
to annotate your workload:

```yaml
app.kubernetes.io/part-of: myAwesomeAppName
app.kubernetes.io/name: myAwesomeWorkload
app.kubernetes.io/version: myAwesomeWorkloadVersion
```

Note the following:

* The Keptn Annotations/Labels take precedence
  over the Kubernetes recommended labels.
* If the Workload has no version annotation/labels
  and the pod has only one container,
  the Lifecycle Toolkit takes the image tag as version
  (if it is not "latest").

This process is demonstrated in the
[Keptn Lifecycle Toolkit: Installation and KeptnTask Creation in Mintes](https://www.youtube.com/watch?v=Hh01bBwZ_qM)
video.

## Pre- and post-deployment checks

Further annotations are necessary
to run pre- and post-deployment checks:

```yaml
keptn.sh/pre-deployment-tasks: verify-infrastructure-problems
keptn.sh/post-deployment-tasks: slack-notification,performance-test
```

The value of these annotations are
Keptn [CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
called `KeptnTaskDefinition`s.
These CRDs contain re-usable "functions"
that can execute before and after the deployment.
In this example, before the deployment starts,
a check for open problems in your infrastructure is performed.
If everything is fine, the deployment continues and afterward,
a slack notification is sent with the result of the deployment
and a pipeline to run performance tests is invoked.
Otherwise, the deployment is kept in a pending state
until the infrastructure is capable of accepting deployments again.

A more comprehensive example can be found in our
[examples folder](https://github.com/keptn/lifecycle-toolkit/tree/main/examples/sample-app),
where we use [Podtato-Head](https://github.com/podtato-head/podtato-head)
to run some simple pre-deployment checks.

To run the example, use the following commands:

```shell
cd ./examples/podtatohead-deployment/
kubectl apply -f .
```

Afterward, you can monitor the status of the deployment using

```shell
kubectl get keptnworkloadinstance -n podtato-kubectl -w
```

The deployment for a Workload stays in a `Pending`
state until the respective pre-deployment check is completed.
Afterwards, the deployment starts and when it is marked  `Succeeded`,
the post-deployment checks start.
