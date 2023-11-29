# Keptn Observability

Keptn provides sophisticated observability features
that enhance your existing cloud-native deployment environment.
These features are useful whether or not you use a GitOps strategy.

The following is an imperative walkthrough.

## Prerequisites

- [Docker](https://docs.docker.com/get-started/overview/)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)
- [Helm](https://helm.sh/docs/intro/install/)
- A Kubernetes cluster >= 1.24 (we recommend [Kubernetes kind](https://kind.sigs.k8s.io/docs/user/quick-start/))
  (`kind create cluster`)

## Objectives

- Install Keptn on your cluster
- Annotate a namespace and deployment to enable Keptn
- Install Grafana and Observability tooling to view DORA metrics and OpenTelemetry traces

## System Overview

By the end of this page, here is what will be built.
The system will be built in stages.

![system overview](../assets/install01.png)

## The Basics: A Deployment, Keptn and DORA Metrics

To begin our exploration of the Keptn observability features, we will:

- Deploy a simple application called `keptndemo`.

Keptn will monitor the deployment and generate:

- An OpenTelemetry trace per deployment
- DORA metrics

![the basics](../assets/install02.png)

Notice though that the metrics and traces have nowhere to go.
That will be fixed in a subsequent step.

## Step 1: Install Keptn

Install Keptn using Helm:

```shell
helm repo add keptn https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn keptn/keptn -n keptn-system --create-namespace --wait
```

Keptn will need to know where to send OpenTelemetry traces.
Of course, Jaeger is not yet installed so traces have nowhere to go (yet),
but creating this configuration now means the system is preconfigured.

Save this file as `keptnconfig.yaml`.
It doesn't matter where this file is located on your local machine:

```yaml
