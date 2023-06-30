---
title: KLT End-to-end exercise
description: Implement full deployment orchestration
weight: 85
---

This page gives instructions for installing the Keptn Lifecycle Toolkit
and running a simple Keptn application to familiarize yourself
with how the Keptn Lifecycle Toolkit works and
implements full deployment orchestration.

You will learn how to do the following:

* Use the Keptn Lifecycle Toolkit to control the deployment of your application
* Connect the lifecycle-toolkit to Prometheus
* Use pre-deployment tasks to check if a dependency is met before deploying a workload
* Use post-deployment tasks on an application level to send a notification

## Prerequisites

To complete this exercise, you need:

* A Kubernetes cluster running Kubernetes 1.24 or later.

* [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
installed on your system.

You can use an existing cluster
or you can create a local cluster to use; see
[Create local Kubernetes cluster](../../install/k8s.md/#create-local-kubernetes-cluster)
for details.

Run the following and verify that both client and server versions
are running Kubernetes versions greater than or equal to v1.24.
In this example, both client and server are at v1.24.0
so are appropriate for the Keptn Lifecycle Toolkit.

```shell
kubectl version --short
```

The output should look like the following

```shell
Client Version: v1.24.0
Kustomize Version: v4.5.4
Server Version: v1.24.0
```

## Clone the repository for this exercise

We provide a repository that contains the sample application
used in this exercise
as well as some helper scripts
that make it easier for you to set up your environment.
Use the following command to clone this repository:

```shell
git clone https://github.com/keptn-sandbox/lifecycle-toolkit-examples.git
cd lifecycle-toolkit-examples
```

## Install the required observability features

The Keptn Lifecycle Toolkit emits OpenTelemetry data as standard
but the toolkit does not come pre-bundled with observability backend tooling.
This is deliberate as it provides flexibility
for you to bring your own observability backend
that consumes this emitted data.

In order to use the observability features of the lifecycle toolkit,
we need a monitoring and tracing backend.

In this guide, we use:

* [Prometheus](https://prometheus.io/) for Metrics
* [Jaeger](https://jaegertracing.io) for Traces
* [Grafana](https://github.com/grafana/) for Dashboarding

Install these together with the KLT with the following command:

```shell
make install
```

> **Note**
To export traces to the OpenTelemetry Collector, you need a
 [KeptnConfig](../../yaml-crd-ref/config.md)
 CRD with `spec.OTelCollectorUrl` specified
 in the namespace where KLT is installed.

## The Demo Application

For this demonstration, we use a slightly modified version of
[the PodTatoHead](https://github.com/podtato-head/podtato-head) application.

<!-- markdown-link-check-disable-next-line -->
![img.png](assets/podtatohead.png)

Over time, we will evolve this application
from a simple manifest to a Keptn-managed application:

* We install it with **kubectl**
  then add pre- and post-deployment tasks.
  * For this, we check if the entry service is available
    before the other services are scheduled.
* We then add evaluations to ensure
    that our infrastructure is in good shape before we deploy the application.
* Finally, we evolve to a GitOps driven deployment
  and notify an external webhook service when the deployment has finished.

## Install the Demo Application (Version 1)

In the first version of the Demo application,
the Keptn Lifecycle Toolkit evaluates metrics provided by Prometheus
and checks if the specified amount of CPUs are available
before deploying the application

To install it, simply apply the manifest:

```shell
make deploy-version-1
```

You can watch the progress of the deployment as follows:

### Watch workload state

When the Lifecycle Toolkit detects workload labels
("app.kubernetes.io/name" or "keptn.sh/workload") on a resource,
a KeptnWorkloadInstance (kwi) resource is created.
Using this resource you can watch the progress of the deployment.

```shell
kubectl get keptnworkloadinstances -n podtato-kubectl
```

This shows the current status of the Workloads
and in which phase they are at the moment.
You can get more detailed information about the workloads
by describing one of the resources:

```shell
kubectl describe keptnworkloadinstances podtato-head-podtato-head-entry -n podtato-kubectl
```

> **Note**
The event stream of the object contains more detailed information

<details>
<summary>Watch application state</summary>
Although you didn't specify an application in your manifest,
the Lifecycle Toolkit assumes that this is a single-service application
and creates an ApplicationVersion (kav) resource for you.

Using `kubectl get keptnappversions -n podtato-kubectl`,
you can see the state of these resources.
</details>

<details>
<summary>Watch pods</summary>
Obviously, you should see that the pods are starting normally.
You can watch the state of the pods using:

```shell
kubectl get pods -n podtato-kubectl
```

</details>

Furthermore, you can port-forward the podtato-head service
to your local machine and access the application via your browser:

```shell
make port-forward-grafana
```

In your browser (<http://localhost:3000>),
log in with the user `admin` and the password `admin`).
You can open the Dashboard `Keptn Applications`
and see the current state of the application,
which should be similar to the following:

<!-- markdown-link-check-disable-next-line -->
![grafana.png](assets/grafana.png)

In this screen you get the following information:

* Successful/Failed Deployments
* Time between Deployments
* Deployment Time per Version
* The link to the Trace of the deployment

After some time (~60 seconds),
you should see one more failed deployment in your dashboard.
You can click on the link to the trace and see the reason for the failure:

<!-- markdown-link-check-disable-next-line -->
![trace-failed.png](assets/trace-failed.png)

In this case, we see the name of the failed pre-deployment evaluation
and the reason for the failure.
In this case, the minimum amount of CPUs is not met.
This is a problem we can solve by changing the threshold in the evaluation file.

## Install the Demo Application (Version 2)

To achieve this, we changed the operator in the evaluation file
(sample-app/version-2/app-pre-deploy-eval) from `<` to `>`
and applied the new manifest:

```shell
kubectl apply -f sample-app/version-2
```

After this, you can inspect the new state of the application
using the same commands as before.
You should see that the deployment is now successful
and that the trace is also updated.
You should also see in the Grafana Dashboards
that the deployment was successful.

Congratulations!
You have successfully deployed an application
using the Keptn Lifecycle Toolkit!
