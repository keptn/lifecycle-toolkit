---
date: 2024-02-01
authors: [odubajDT]
description: >
  In this blog post you will learn how to configure Keptn via the new Helm umbrella chart.
categories:
  - SRE
  - Helm
  - Installation
  - Upgrade
comments: true
---

# Configure Keptn via the new Helm umbrella chart

In the dynamic world of cloud-native and microservices, installation and continuous upgrades of distributed systems
is a priority.
Embarking on a journey into the world of containerized applications often leads us into a labyrinth of configuration
and deployment intricacies.
Amidst this complexity, Helm charts emerge as the guiding light, streamlining the orchestration of Kubernetes-based
applications.

In this article, we're delving into the topic of Helm charts, in particular Umbrella charts - an approach that
simplifies the deployment and maintenance of microservice application architectures.
Let's uncover how these consolidated charts pave the way for smoother deployments, enhance scalability, and bring
harmony to the orchestration of complex Kubernetes environments.
<!-- more -->

## What is Helm?

[Helm](https://helm.sh/) is a powerful package manager for Kubernetes.
It simplifies and streamlines the deployment and management of applications on Kubernetes clusters.
It acts as a tool for defining, installing, and upgrading complex Kubernetes applications using pre-configured
templates called charts.
These charts encapsulate all the necessary Kubernetes resources and configuration files needed to deploy an
application, making it easier for developers and operators to manage applications consistently across different
environments.
Helm's ability to package, version, and share applications as charts facilitates collaboration and standardization
within the Kubernetes ecosystem, significantly improving the efficiency of deploying and maintaining containerized
applications.
The basic configuration mechanism for deployments is a `values.yaml` file which contains all the configuration options
with the default values for the given software.

## What is Keptn?

[Keptn](https://lifecycle.keptn.sh/) is a CNCF project for continuous delivery and automated operations.
It helps developers and platform engineering teams automate deployment, monitoring, and management of applications
running in cloud environments.
Keptn works with standard deployment software like ArgoCD or Flux, consolidating metrics, observability, and analysis
for all the microservices that comprise your deployed software as well as providing checks and executables that can
run before and/or after the deployment.

It also comprises multiple components with different feature sets as well as different installation and upgrade
requirements so proper configuration can be cumbersome for people unfamiliar with Helm charts.

Keptn uses the Helm umbrella charts to manage multiple sub-charts (one sub-chart per component) as a single entity.
This concept provides us the opportunity to install and manage each Keptn component as a fully independent microservice,
but additionally also as a complete toolkit.

In Keptn, just use a single `values.yaml` file to configure both the global umbrella chart parameters that apply to
all components, and each individual sub-chart's parameters.
Each individual Keptn component has all the possible helm configuration parameters documented in the
[Keptn repository](https://github.com/keptn/lifecycle-toolkit).

In the individual sub-charts, we can set parameters such as image `tag`, `imagePullPolicy` or `annotations`.
Additionally, global parameters which affect all the sub-charts can be set, like the global image `registry` or
global `annotations`.
These options are documented [here](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/README.md).

Note that setting the values for individual components overrides the global ones.
For example, if the same `annotation` option is configured for the global chart and the component, the component value
is favored.

## Installation of Keptn with default values

Let's look at a real-life example of a clean installation of Keptn using umbrella charts.
For basic installation of all components with the use of default values for configuration parameters you need to
execute the following commands:

```shell
helm repo add keptn https://charts.lifecycle.keptn.sh 
helm repo update 
helm upgrade --install keptn keptn/keptn -n keptn-system --create-namespace 
```

You can see that it's not necessary to create a custom "values.yaml" file when using the default values.
Keptn will be installed with all the components (`metrics-operator`, `lifecycle-operator`, `cert-manager`), with
default Helm values set for each of them.

## Installation of Keptn with custom values

Let's now look at a use-case, where we want to use only Keptn metrics and therefore do not need to install the
`lifecycle-operator`, which provides functionality for Observability and/or Release Lifecycle Management.
Additionally to `metrics-operator`, we need to install `keptn-cert-manager,` which will provide certificates
webhooks present in the `metrics-operator`.
For this all, we need to create a `values.yaml` file with the following content:

```yaml
{% include "./keptn-helm-umbrella-charts/values.yaml" %}
```

This `values.yaml` file disables the `lifecycle-operator` and it sets the used `image registry` and `tag` of
the `cert-manager`. `tag`, `annotations` and `replicas` are set for the `metrics-operator`.
These parameters are configuring each installed component via their sub-charts.
Also, we configure the global annotations for all installed components (`metrics-operator` and `cert-manager`
only, since `lifecycle-operator` is disabled).
Note that the `myMetricsKey` annotation is set to different values in global parameters (`globalValue2`) and
in `metrics-operator` sub-chart parameters (`metricsValue1`).
Here the annotation defined in the `metrics-operator` component configuration will be used.

When the `values.yaml` file is ready, we use it to install Keptn with the following commands:

```shell
helm repo add keptn https://charts.lifecycle.keptn.sh 
helm repo update 
helm upgrade --install keptn keptn/keptn -n keptn-system --create-namespace --values=values.yaml 
```

More information about the installation and upgrade options can be found in the official
[Keptn documentation](https://lifecycle.keptn.sh/docs/install/) and in the linked Helm values documentation.

## Summary

Let's sum up what we have seen.
We started with a brief description of Keptn, Helm and Umbrella charts.
Next, we described how the installation of Keptn components can be customized with Helm values that are under
the aegis of the Helm umbrella chart.
In the end, we went through an example of the configurable parameters of individual components or global values
and how the custom `values.yaml` file can be used for installation.

## Useful links

- <https://lifecycle.keptn.sh/docs/install/>
- <https://github.com/keptn/lifecycle-toolkit/>
- <https://github.com/keptn/lifecycle-toolkit-charts/>
- <https://helm.sh/>
