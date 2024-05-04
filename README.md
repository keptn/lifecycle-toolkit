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
- Checking the Application Health in a declarative (cloud-native) way
- Standardized way to run pre- and post-deployment tasks
- Provide out-of-the-box Observability
- Deployment lifecycle management

![Operator Maturity Model with third level circled in](./assets/operator-maturity.jpg)

Keptn can be seen as a general purpose and declarative
[Level 3 operator](https://operatorframework.io/operator-capabilities/)
for your Application.
For this reason, Keptn is agnostic to deployment tools
that are used and works with any GitOps solution.

For more information about the core concepts of Keptn, see
our core concepts
[documentation section](https://keptn.sh/stable/docs/core-concepts/).

## Status

Status of the different features:

- ![status](https://img.shields.io/badge/status-stable-brightgreen)
  Observability: expose [OTel](https://opentelemetry.io/) metrics and traces of your deployment.
- ![status](https://img.shields.io/badge/status-stable-brightgreen)
  K8s Custom Metrics: expose your Observability platform via the [Custom Metric API](https://github.com/kubernetes/design-proposals-archive/blob/main/instrumentation/custom-metrics-api.md).
- ![status](https://img.shields.io/badge/status-stable-brightgreen)
  Release lifecycle: handle pre- and post-checks of your Application deployment.
- ![status](https://img.shields.io/badge/status-stable-brightgreen)
  Certificate Manager: automatically configure TLS certificates with the
  [Keptn Certificate Manager](https://keptn.sh/stable/docs/components/certificate-operator/).
  You can instead
  [configure your own certificate manager](https://keptn.sh/stable/docs/installation/configuration/cert-manager/) to provide
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
[Running Keptn with vCluster](https://keptn.sh/stable/docs/installation/configuration/vcluster/)
for more information.

Use the following command sequence
to install the latest release of Keptn:

```shell
helm repo add keptn https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn keptn/keptn -n keptn-system --create-namespace --wait
```

### Monitored namespaces

Keptn must be installed in its own namespace
that does not run other major components or deployments.

By default, the Keptn lifecycle orchestration
monitors all namespaces in the cluster
except for a few namespaces that are reserved
for specific Kubernetes and other components.
You can modify the Helm chart to specify the namespaces
where the Keptn lifecycle orchestration is allowed.
For more information, see the "Namespaces and Keptn" page in the
[Configuration](https://keptn.sh/stable/docs/installation/configuration/)
section of the documentation.

## More information

For more info about Keptn, please see our
[documentation](https://keptn.sh/stable/docs/).

You can also find a number of video presentations and demos
about Keptn on the
[YouTube Keptn channel](https://www.youtube.com/@keptn).
Videos that refer to the "Keptn Lifecycle Controller"
are relevant for the Keptn project.

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
