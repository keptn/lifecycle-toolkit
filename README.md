# Keptn

![build](https://img.shields.io/github/actions/workflow/status/keptn/lifecycle-toolkit/CI.yaml?branch=main)
![Codecov](https://img.shields.io/codecov/c/github/keptn/lifecycle-toolkit?token=KPGfrBb2sA)
![goversion](https://img.shields.io/github/go-mod/go-version/keptn/lifecycle-toolkit?filename=lifecycle-operator/go.mod)
![version](https://img.shields.io/github/v/release/keptn/lifecycle-toolkit)
[![GitHub Discussions](https://img.shields.io/github/discussions/keptn/lifecycle-toolkit)](https://github.com/keptn/lifecycle-toolkit/discussions)
[![Artifacthub Badge](https://img.shields.io/badge/Keptn-blue?style=flat&logo=artifacthub&label=Artifacthub&link=https%3%2F%2Fartifacthub.io%2Fpackages%2Fhelm%2Flifecycle-toolkit%2Fkeptn)](https://artifacthub.io/packages/helm/lifecycle-toolkit/keptn)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/3588/badge)](https://www.bestpractices.dev/projects/3588)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/keptn/lifecycle-toolkit/badge)](https://securityscorecards.dev/viewer/?uri=github.com/keptn/lifecycle-toolkit)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkeptn%2Flifecycle-toolkit.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkeptn%2Flifecycle-toolkit?ref=badge_shield&issueType=license)
[![CLOMonitor](https://img.shields.io/endpoint?url=https://clomonitor.io/api/projects/cncf/keptn/badge)](https://clomonitor.io/projects/cncf/keptn)

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
- ![status](https://img.shields.io/badge/status-beta-yellow)
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

## Community

Find details on regular hosted community events in the [keptn/community repo](https://github.com/keptn/community)
and our Slack channel(s) in the [CNCF Slack workspace.](https://cloud-native.slack.com/messages/keptn/)

## Roadmap

You can find our roadmap [here](https://github.com/orgs/keptn/projects/10).

## Governance

- [Community Membership](https://github.com/keptn/community/blob/main/community-membership.md):
  Guidelines for community engagement, contribution expectations,
  and the process for becoming a community member at different levels.

- [Members and Charter](https://github.com/keptn/community/blob/main/governance/members-and-charter.md):
  Describes the formation and responsibilities of the Keptn Governance Committee,
  including its scope, members, and core responsibilities.

## Installation

Keptn can be installed on any Kubernetes cluster
running Kubernetes >=1.24.

For users running [vCluster](https://www.vcluster.com/),
please note that you may need to modify
your configuration before installing Keptn; see
[Running Keptn with vCluster](https://main.lifecycle.keptn.sh/docs/install/install//#running-keptn-with-vcluster)
for more information.

Use the following command sequence
to install the latest release of Keptn:

```shell
helm repo add keptn https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn keptn/keptn -n keptn-system --create-namespace --wait
```

### Keptn and namespaces

Keptn must be installed in its own namespace
that does not run other major components or deployments.

By default, the Keptn lifecycle orchestration
monitors all namespaces in the cluster
except for a few namespaces that are reserved
for specific Kubernetes and other components.
You can modify the Helm chart to specify the namespaces
where the Keptn lifecycle orchestration is allowed.
For more information, see the "Namespaces and Keptn" page in the
[Configuration](https://keptn.sh/stable/docs/installation/configuration/index.md)
section of the documentation.

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
  and implement its various features.
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

## Contributing

For more information about contributing to Keptn, please
refer to the [Contribution guide](https://keptn.sh/stable/docs/contribute/)
section of the documentation.

To set up your local Keptn development environment, please follow
[these steps](https://keptn.sh/stable/docs/contribute/software/dev-environ/#first-steps)
for new contributors.

## License

Please find more information in the [LICENSE](LICENSE) file.

## Thanks to all the people who have contributed 💜

<!-- markdownlint-disable-next-line MD033 -->
<a href="https://github.com/keptn/lifecycle-toolkit/graphs/contributors">
<!-- markdownlint-disable-next-line MD033 MD045 -->
  <img src="https://contrib.rocks/image?repo=keptn/lifecycle-toolkit" />
</a>

Made with [contrib.rocks](https://contrib.rocks).
