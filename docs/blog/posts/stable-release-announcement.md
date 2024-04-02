---
date: 2024-03-15
authors: [agardnerIT]
description: >
  This blog post announces the release candidate for Keptn v2.
categories:
  - Announcement
comments: true
slug: announcing-the-official-release-candidate-for-keptn-v2
---

# Announcing the Release Candidate for Keptn v2

The Keptn project is proud to announce a release candidate for what will become Keptn v2.
User feedback to the Keptn project has been clear, and we’ve listened.
We thank the users who have given us feedback
and all the community members who have contributed to this project.

We are sharing this release candidate
so that the community and end users can participate in the testing
to ensure that the actual release is as robust as possible.
We anticipate (as long as no release blockers are identified)
that Keptn v2 will be released approximately 1 month after this announcement.

<!-- more -->

## What problem does Keptn solve?

Deploying software using Kubernetes is, on the surface, easy.
Just use `kubectl` or a GitOps solution like
[ArgoCD](https://argoproj.github.io/cd/)
to deploy a YAML file and that’s it, right?
Well, no.
Not normally.
There is a lot more to ensuring a healthy deployment.
The pod(s) may be running, but that doesn’t automatically mean that the application is healthy.
This is the problem domain that Keptn acts upon.
By using Keptn, in combination with your standard
deployment tooling or practices, you can move from “I guess it’s OK” to “I know it’s OK”.
Keptn allows you to wrap governance and automated checks around the deployment process to ensure that
the end-to-end process of deploying is healthy and your application is meeting the SLOs you’ve defined.

## What’s New in the Keptn v2 release candidate?

Keptn v2-rc introduces the following new features to the previous Keptn release:

- Non-Blocking Tasks
- A new “promotion” stage

### Non-blocking tasks and evaluations

Keptn offers the ability to define tasks and SLO evaluations that run either before or after a deployment.
By design, any pre-deployment task or SLO evaluation that fails will block the deployment.
Often, this is the behaviour you want –- if a downstream dependency is unavailable or unhealthy,
you probably don’t want to complete the deployment.
However, when first testing and implementing Keptn in your development environment,
this may cause deployments to be “pending” without an obvious cause.

In Keptn v2 this blocking behaviour for pre tasks and evaluations can be temporarily disabled for the cluster
until you are sure that your tasks and evaluations are performing appropriately.
To implement this feature, set `spec.blockDeployment: [true|false]` in the
[KeptnConfig](../../docs/reference/crd-reference/config.md) resource.
The default behaviour is for Keptn to block deployments (i.e. `spec.blockDeployment: true`).

### The Promotion Stage

Keptn v2 introduces a new “promotion” phase
to support multi-stage application delivery.
Keptn is commonly used as part of [GitOps practices](https://opengitops.dev/) and thus,
users want a dedicated way to promote an application to the next stage in my environment.
In this release candidate,
the new stage is disabled by default and can be controlled via the Helm flag:
`lifecycleOperator.promotionTasksEnabled: [true|false]`.
We expect that the upcoming stable release will be shipped with this feature enabled out of the box.

Further information can be found in
The
[Multi-stage application delivery](../../docs/guides/multi-stage-application-delivery.md)
guide.

## Try it out

With this release of Keptn, you can see a true end-to-end picture of everything in the logical order,
potentially from “PR merged” all the way to “deployment complete”.

To try it out, download Keptn v2 RC1 now:

- From
[ArtifactHub](https://artifacthub.io/packages/helm/lifecycle-toolkit/keptn)
- From
[GitHub](https://github.com/keptn/lifecycle-toolkit/releases)

Please provide any feedback -- good and bad -- by posting in the
[#Keptn CNCF Slack channel](https://cloud-native.slack.com/archives/C017GAX90GM)
or by raising issues in our
[GitHub repository](https://github.com/keptn/lifecycle-toolkit/issues).
