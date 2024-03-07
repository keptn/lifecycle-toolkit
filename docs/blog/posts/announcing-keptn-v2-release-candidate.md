---
date: 2024-03-07
authors: [agardnerit]
description: >
  Announcing the Official Release Candidate for Keptn v2
categories:
  - Installation
  - Upgrade
comments: true
---

# Announcing the Official Release Candidate for Keptn v2

The Keptn project is proud to announce a release candidate for what will become Keptn v2.

This release candidate will provide the project with enough end-user validation and bug-finding so that we can release Keptn v2.

We anticipate (as long as no release blockers are identified) that Keptn v2 will be released approximately 1 month after this announcement.

<!-- more -->

## What Problem Does Keptn Solve?

Deploying software using Kubernetes is, on the surface, easy. Just use `kubectl` or a GitOps solution like ArgoCD to deploy a YAML file and that’s it, right?

Well, no. Not normally. There is a lot more to ensuring a healthy deployment. The pod(s) may be running, but that doesn’t automatically mean that the application is healthy.

This is the problem domain that Keptn acts upon. By using Keptn, in combination with your standard deployment tooling or practices, you can move from "I **guess** it’s OK" to "I know it’s OK".

Keptn allows you to wrap governance and automated checks around the deployment process to ensure that the end-to-end process of deploying is healthy, and your application is meeting the SLOs you’ve defined.

## What’s New in Keptn v2?

User feedback to the Keptn project has been clear, and we have listened. Keptn v2 will bring 3 key new features:

1.	Non-Blocking Tasks 
2.	A new "promotion" stage
3.	Ability to pass OpenTelemetry contextual information into, and out of, Keptn

### Non-Blocking Tasks

Keptn offers the ability to perform arbitrary tasks and SLO evaluations both before a deployment and after a deployment. By design, any pre-deployment task or SLO evaluation which fails will block the deployment. Often, this is the behaviour you want – if a downstream dependency is unavailable or unhealthy, you probably don’t want to complete the deployment.
However, for new Keptn users, this behaviour can appear drastic and cause deployments to be "pending" without an obvious cause.

In Keptn v2 this blocking behaviour for pre tasks and evaluationis now configurable. When creating the [KeptnConfig](https://keptn.sh/stable/docs/reference/crd-reference/config/) resource, set `spec.blockDeployment: [true|false]`. The default behaviour is for Keptn to block deployments (ie. `spec.blockDeployment: true`)

### The Promotion Stage

Keptn v2 introduces a new "promotion" stage. Keptn is commonly used alongside [GitOps practices](https://opengitops.dev/) and thus, users want to have a dedicated way to promote an application to the next stage in my environment.
The new stage is enabled by default and can be controlled via the Helm flag: `lifecycleOperator.promotionTasksEnabled: [true|false]`

Further information can be found in the [official documentation](https://keptn.sh/stable/docs/guides/multi-stage-application-delivery/).

### Pass OpenTelemetry Contexts through Keptn

Usually, Keptn will not act in isolation during a deployment. In the logical timeline of events, there will be tools acting prior to Keptn (such as GitOps tools like ArgoCD) and tools that act after Keptn (such as security scanning tools or post-deployment checks).

If the tools prior to Keptn begin generating OpenTelemetry data (eg. Spans and traces), it would be beneficial to see "Keptn’s portion" of the work in the correct context (ie. As part of the same distributed trace). Similarly, anything that happens after Keptn (but still in the same logical "deployment" operation) should be included in the end-to-end trace / timing view.
To achieve this, Keptn now accepts tools to pass w3c trace IDs into Keptn and Keptn will pass W3C trace IDs out of Keptn (thanks to [@geoffrey1330](https://github.com/geoffrey1330) for providing this functionality!)
Now, you will see a true end-to-end picture of everything in the logical operation, potentially from "PR merged" all the way to "deployment complete".

## [Download Keptn v2 RC1 Now!](https://artifacthub.io/packages/helm/lifecycle-toolkit/keptn)

[Keptn v2 Release Candidate 1](https://github.com/keptn/lifecycle-toolkit/releases) is available now on GitHub. Please provide any feedback via the #keptn Slack channel in the [CNCF Slack workspace](https://communityinviter.com/apps/cloud-native/cncf).