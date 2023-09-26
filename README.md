# Keptn

![build](https://img.shields.io/github/actions/workflow/status/keptn/lifecycle-toolkit/CI.yaml?branch=main)
![Codecov](https://img.shields.io/codecov/c/github/keptn/lifecycle-toolkit?token=KPGfrBb2sA)
![goversion](https://img.shields.io/github/go-mod/go-version/keptn/lifecycle-toolkit?filename=lifecycle-operator/go.mod)
![version](https://img.shields.io/github/v/release/keptn/lifecycle-toolkit)
[![GitHub Discussions](https://img.shields.io/github/discussions/keptn/lifecycle-toolkit)](https://github.com/keptn/lifecycle-toolkit/discussions)

This is the primary repository for
the Keptn software and documentation.
Keptn provides a “cloud-native” approach
for managing the application release lifecycle
metrics, observability, health checks,
with pre- and post-deployment evaluations and tasks.
It is an incubating project, under the umbrella of the
[Keptn Application Lifecycle working group](https://github.com/keptn/wg-app-lifecycle).

> **Note** Keptn was developed under the code name of
  "Keptn Lifecycle Toolkit" or "KLT" for short.
  The source code contains many vestiges of these names.

## Goals

Keptn provides Cloud Native teams with the following capabilities:

- Pre-requisite evaluation before deploying workloads and applications
- Finding out when an application (not just a workload) is ready and working
- Checking the Application Health in a declarative (cloud-native) way
- Standardized way to run pre- and post-deployment tasks
- Provide out-of-the-box Observability of the deployment cycle

![Operator Maturity Model with third level circled in](./assets/operator-maturity.jpg)

Keptn can be seen as a general purpose and declarative
[Level 3 operator](https://operatorframework.io/operator-capabilities/)
for your Application.
For this reason, Keptn is agnostic to deployment tools
that are used and works with any GitOps solution.

## Status

Status of the different features:

- ![status](https://img.shields.io/badge/status-stable-brightgreen)
  Observability: expose [OTel](https://opentelemetry.io/) metrics and traces of your deployment.
- ![status](https://img.shields.io/badge/status-alpha-orange)
  K8s Custom Metrics: expose your Observability platform via the [Custom Metric API](https://github.com/kubernetes/design-proposals-archive/blob/main/instrumentation/custom-metrics-api.md).
- ![status](https://img.shields.io/badge/status-beta-yellow)
  Release lifecycle: handle pre- and post-checks of your Application deployment.
- ![status](https://img.shields.io/badge/status-stable-brightgreen)
  Certificate Manager: automatically configure TLS certificates with the
  [Keptn Certificate Manager](https://lifecycle.keptn.sh/docs/concepts/architecture/cert-manager/).
  You can instead
  [configure your own certificate manager](https://lifecycle.keptn.sh/docs/install/cert-manager/) to provide
  [secure communication with the Kube API](https://kubernetes.io/docs/concepts/security/controlling-access/#transport-security).

<!---
alpha ![status](https://img.shields.io/badge/status-alpha-orange) )
beta ![status](https://img.shields.io/badge/status-beta-yellow) )
stable ![status](https://img.shields.io/badge/status-stable-brightgreen) )
-->
The status follows the
[Kubernetes API versioning schema](https://kubernetes.io/docs/reference/using-api/#api-versioning).

## Installation

Keptn can be installed on any Kubernetes cluster
running Kubernetes >=1.24.
Note that Keptn is not currently compatible with
[vcluster](https://github.com/loft-sh/vcluster).

Use the following command sequence
to install the latest release of Keptn:

```shell
helm repo add klt https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn klt/klt -n keptn-lifecycle-toolkit-system --create-namespace --wait
```

## More information

For more info about Keptn, please see our
[documentation](https://lifecycle.keptn.sh/docs/), specifically:

- [Introduction to Keptn](https://lifecycle.keptn.sh/docs/intro/)
  gives an overview of the Keptn facilities.
- [Getting started](https://lifecycle.keptn.sh/docs/getting-started/)
  includes some short exercises to introduce you to Keptn.
- [Installation and upgrade](https://lifecycle.keptn.sh/docs/install/)
  provides information about preparing your Kubernetes cluster
  then installing and enabling Keptn.
- [Implementing Keptn applications](https://lifecycle.keptn.sh/docs/implementing/)
  documents how to integrate Keptn to work with your existing deployment engine
  and implement its variouos features.
- [Architecture](https://lifecycle.keptn.sh/docs/concepts/architecture/) provides detailed technical information
  about how Keptn works.
- [CRD Reference](https://lifecycle.keptn.sh/docs/yaml-crd-ref/) and
  [API Reference](https://lifecycle.keptn.sh/docs/crd-ref/)
  provide detailed reference material for the custom resources
  used to configure Keptn.
- [Contributing to Keptn](https://lifecycle.keptn.sh/contribute/)
  provides information about how to contribute to the Keptn project.

You can also find a number of video presentations and demos
about Keptn on the
[YouTube Keptn channel](https://www.youtube.com/@keptn).
Videos that refer to the "Keptn Lifecycle Controller"
are relevant for the Keptn project.

## Architecture

Keptn is composed of the following components:

- Keptn Lifecycle Operator
- Keptn Scheduler

The Keptn Lifecycle Operator contains several controllers for Keptn CRDs
and a Mutating Webhook.
The Keptn Scheduler ensures that Pods are started
only after the pre-deployment checks have finished successfully.

A Kubernetes
[Manifest](https://monokle.io/learn/kubernetes-manifest-files-explained#:~:text=Kubernetes%20Manifest%20files!-,What%20is%20a%20Kubernetes%20Manifest%20File%3F,you%20want%20in%20your%20cluster).
which is annotated with Keptn specific annotations,
is applied to the Kubernetes Cluster.
Afterward, the Keptn Scheduler is injected (via Mutating Webhook),
and Kubernetes Events for Pre-Deployment are sent to the event stream.
The Event Controller watches for events
and triggers a Kubernetes Job to fullfil the Pre-Deployment.
After the Pre-Deployment has finished,
the Keptn Scheduler schedules the Pod to be deployed.
The KeptnApp and KeptnWorkload Controllers
watch for the workload resources to finish
and then generate a Post-Deployment Event.
After the Post-Deployment checks,
SLOs can be validated using an interface
for retrieving SLI data from a provider
e.g, [Prometheus](https://prometheus.io/).
Finally, Keptn exposes Metrics and Traces
of the entire Deployment cycle with
[OpenTelemetry](https://opentelemetry.io/).

![Keptn Architecture](./assets/architecture.png)

### Webhook

Annotating a namespace subjects it to the effects of the mutating webhook:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: podtato-kubectl
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled"  # this line tells the webhook to handle the namespace
```

The mutating webhook only modifies specifically annotated resources in the annotated namespace.
When the webhook receives a request for a new pod,
it looks for the workload annotations:

```yaml
keptn.sh/workload: "some-workload-name"
```

The mutation consists in changing the scheduler used for the deployment
with the Keptn Scheduler.
The webhook then creates a workload and app resource per annotated resource.
You can also specify a custom app definition with the annotation:

```yaml
keptn.sh/app: "your-app-name"
```

In this case the webhook does not generate an app,
but it expects that the user will provide one.
Additionally, it computes a version string,
using a hash function that takes certain properties of the pod as parameters
(e.g. the images of its containers).
Next, it looks for an existing instance of a `Workload CRD`
for the specified workload name:

- If it finds the `Workload`,
  it updates its version according to the previously computed version string.
  In addition, it includes a reference to the ReplicaSet UID of the pod
  (i.e. the Pods owner),
  or the pod itself, if it does not have an owner.
- If it does not find a workload instance,
  it creates one containing the previously computed version string.
  In addition, it includes a reference to the ReplicaSet UID of the pod
  (i.e. the Pods owner), or the pod itself, if it does not have an owner.

It uses the following annotations for the specification
of the pre/post deployment checks that should be executed for the `Workload`:

- `keptn.sh/pre-deployment-tasks: task1,task2`
- `keptn.sh/post-deployment-tasks: task1,task2`

and for the Evaluations:

- `keptn.sh/pre-deployment-evaluations: my-evaluation-definition`
- `keptn.sh/post-deployment-evaluations: my-eval-definition`

After either one of those actions has been taken,
the webhook sets the scheduler of the pod
and allows the pod to be scheduled.

### Scheduler

After the Webhook mutation, the Keptn-Scheduler handles the annotated resources.
The scheduling flow follows the default
[scheduler](https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/)
behavior,
since it implements a scheduler plugin based on the
[scheduling framework]( https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/).
For each pod, at the very end of the scheduling cycle,
the plugin verifies that the pre deployment checks have terminated
by retrieving the current status of the WorkloadInstance.
Only when that is successful is the pod bound to a node.

## Install a dev build

The [GitHub CLI](https://cli.github.com/) can be used to download the manifests of the latest CI build.

```bash
gh run list --repo keptn/lifecycle-toolkit # find the id of a run
gh run download 3152895000 --repo keptn/lifecycle-toolkit # download the artifacts
kubectl apply -f ./keptn-lifecycle-operator-manifest/release.yaml # install the lifecycle-operator
kubectl apply -f ./scheduler-manifest/release.yaml # install the scheduler
```

Instead, if you want to build and deploy the operator into your cluster
directly from the code, you can type:

```bash
RELEASE_REGISTRY=<YOUR_DOCKER_REGISTRY>
# (optional)ARCH=<amd64(default)|arm64v8>
# (optional)CHART_APPVERSION=<YOUR_PREFERRED_TAG (defaulting to current time)>

# Build and deploy the dev images to the current kubernetes cluster
make build-deploy-dev-environment
```

## License

Please find more information in the [LICENSE](LICENSE) file.

## Thanks to all the people who have contributed 💜

<!-- markdownlint-disable-next-line MD033 -->
<a href="https://github.com/keptn/lifecycle-toolkit/graphs/contributors">
<!-- markdownlint-disable-next-line MD033 -->
  <img src="https://contrib.rocks/image?repo=keptn/lifecycle-toolkit" />
</a>

Made with [contrib.rocks](https://contrib.rocks).

<!-- markdownlint-disable-next-line MD033 MD013 -->
<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=858843d8-8da2-4ce5-a325-e5321c770a78" />
