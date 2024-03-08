---
date: 2024-03-10
authors: [agardnerIT]
description: >
  This blog post announces the release candidate for Keptn v2.
comments: true
---

# Announcing the Official Release Candidate for Keptn v2

The Keptn project is proud to announce a release candidate for what will become Keptn v2.
We are sharing this release candidate
so that the community and end users can participate in the testing
to ensure that the actual release is as robust as possible.
We anticipate (as long as no release blockers are identified)
that Keptn v2 will be released approximately 1 month after this announcement.

## What Problem Does Keptn Solve?

Deploying software using Kubernetes is, on the surface, easy.
Just use kubectl or a GitOps solution like ArgoCD to deploy a YAML file and that’s it, right?
Well, no.
Not normally.
There is a lot more to ensuring a healthy deployment.
The pod(s) may be running, but that doesn’t automatically mean that the application is healthy.
This is the problem domain that Keptn acts upon.
By using Keptn, in combination with your standard
deployment tooling or practices, you can move from “I guess it’s OK” to “I know it’s OK”.
Keptn allows you to wrap governance and automated checks around the deployment process to ensure that
the end-to-end process of deploying is healthy and your application is meeting the SLOs you’ve defined.

## What’s New in Keptn v2-rc?

User feedback to the Keptn project has been clear, and we’ve listened.
Keptn v2-rc will bring the following new features:

- Non-Blocking Tasks
- A new “promotion” stage

### Non-blocking tasks and evaluations

Keptn offers the ability to perform arbitrary tasks and SLO evaluations both before a deployment and after a deployment.
By design, any pre-deployment task or SLO evaluation that fails will block the deployment.
Often, this is the behaviour you want – if a downstream dependency is unavailable or unhealthy,
you probably don’t want to complete the deployment.
However, for new Keptn users, this behaviour can appear drastic and cause deployments to be
“pending” without an obvious cause.
In Keptn v2 this blocking behaviour for pre tasks and evaluations can be temporarily disabled for the cluster
until you are sure that your tasks and evaluations
are performing appropriately.
When creating the KeptnConfig resource, set `spec.blockDeployment: [true|false]`.
The default behaviour is for Keptn to block deployments (i.e. `spec.blockDeployment: true`)

### The Promotion Stage

Keptn v2 introduces a new “promotion” stage
to support multi-stage application delivery.
Keptn is commonly used alongside [GitOps practices](https://opengitops.dev/) and thus,
users want to have a dedicated way to promote an application to the next stage in my environment.
The new stage is disabled by default and can be controlled via the Helm flag:
`lifecycleOperator.promotionTasksEnabled: [true|false]`.
The upcoming stable release will be shipped with this feature enabled out of the box.

Further information can be found in
The
[Multi-stage application delvery](https://keptn.sh/stable/docs/guides/multi-stage-application-delivery/)
guide.

## Try it out

Now, you can see a true end-to-end picture of everything in the logical operation,
potentially from “PR merged” all the way to “deployment complete”.
[Download Keptn v2 RC1 Now!](https://artifacthub.io/packages/helm/lifecycle-toolkit/keptn)

[Keptn v2 Release Candidate 1](https://github.com/keptn/lifecycle-toolkit/releases) is available now on GitHub.
Please provide any feedback via the #keptn Slack channel in the
[CNCF Slack workspace](https://communityinviter.com/apps/cloud-native/cncf) or raising issues in our
[GitHub repository](https://github.com/keptn/lifecycle-toolkit/issues).
