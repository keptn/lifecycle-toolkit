---
date: 2024-03-07
authors: [agardnerit,bacherfl]
description: >
  Announcing the Official Release Candidate for Keptn v2
categories:
  - Installation
  - Upgrade
comments: true
---

# Announcing the Official Release Candidate for Keptn v2

The Keptn project is proud to announce a release candidate for what will become Keptn v2.
This release candidate provides the project with enough end-user validation
and bug-finding so that we can release Keptn v2.
We anticipate (as long as no release blockers are identified) that Keptn
v2 will be released approximately 1 month after this announcement.

<!-- more -->

## What Problem Does Keptn Solve?

Deploying software using Kubernetes is, on the surface, easy.
Just use `kubectl` or a GitOps solution like ArgoCD to deploy a YAML file and that's it, right?
Well, no.
Not normally.
There is a lot more to ensuring a healthy deployment,
and there are a lot of potential problems that can
prevent the successful deployment of a Kubernetes application, such as:

- Not enough resources (i.e. CPU, Memory) being available in the cluster.
- An external dependency (e.g. a Database hosted outside the cluster) not being reachable
- Unhealthy state of your infrastructure (e.g. due to vulnerabilities)
- The cluster or namespace not being monitored correctly - You would not want to deploy into the dark in this case.
- Unexpected performance degradation after rolling out a new version of your application

This is the problem domain on which Keptn acts.
By using Keptn, in combination with your standard deployment tooling or practices,
you can move from "I guess it's OK" to "I know it's OK".
Keptn allows you to wrap governance and automated checks around the deployment
process to ensure that the end-to-end process of deploying is
healthy and your application is meeting the SLOs you've defined.
This is done by configuring Keptn to execute KeptnTasks and
KeptnEvaluations before and after the deployment of a workload.
With these KeptnTasks, you can implement checks that ensure your
infrastructure is in a healthy state before a new version of a deployment is rolled out.
This way, you can avoid the problems mentioned above by:

- Verifying that your cluster has enough resources to run a workload
- Checking the reachability of external services that might be required by your workloads
- Ensure that your cluster is monitored by your monitoring provider
- Verifying that your monitoring provider did not detect any open problems within your cluster

After a successful deployment of a workload, Keptn also provides the means to do post-deployment checks.
These can be crucial in ensuring that the performance of your newly deployed workload meets your expectations,
before promoting it into the next stage.

## What’s New in Keptn v2?

User feedback to the Keptn project has been clear, and we have listened.
Keptn v2 brings 3 key new features:

- Non-Blocking Tasks 
- A new "promotion" stage
- Ability to pass OpenTelemetry contextual information into, and out of, Keptn

### Non-Blocking Tasks

Keptn offers the ability to perform arbitrary tasks and SLO
evaluations both before a deployment and after a deployment.
By design, any pre-deployment task or SLO evaluation
which fails will block the deployment.
Often, this is the behaviour you want – if a downstream dependency is unavailable or unhealthy,
you probably don't want to complete the deployment.
However, for new Keptn users, this behaviour can appear drastic
and cause deployments to be "pending" without an obvious cause.
In Keptn v2 this blocking behaviour for pre tasks and evaluations is now configurable.
When creating the [KeptnConfig](https://keptn.sh/stable/docs/reference/crd-reference/config/) resource,
set `spec.blockDeployment: [true|false]`.
The default behaviour is for Keptn to block deployments (i.e. `spec.blockDeployment: true`)
until the related pre-deployment tasks and
evaluations defined in the Keptn configuration have passed successfully.

### The Promotion Stage

Keptn v2 introduces a new "promotion" phase.
Keptn is commonly used alongside [GitOps practices](https://opengitops.dev/) and thus,
users want to have a dedicated way to promote an application to the next stage in my environment.
Being a recently added feature this is currently hidden behind a feature flag in the release candidate.
However, it can easily be enabled during the installation
of Keptn via the Helm flag: `lifecycleOperator.promotionTasksEnabled: [true|false]`.

Further information can be found in the
[Multi-stage application delivery](https://keptn.sh/stable/docs/guides/multi-stage-application-delivery/)
guide.

### Pass OpenTelemetry Contexts through Keptn

Usually, Keptn does not act in isolation during a deployment.
In the logical timeline of events, there will be tools acting prior to
Keptn (such as GitOps tools like ArgoCD) and tools that act after
Keptn (such as security scanning tools or post-deployment checks).

If the tools prior to Keptn begin generating OpenTelemetry data (e.g. Spans and traces),
it would be beneficial to see "Keptn’s portion" of the work in the correct context
(i.e. as part of the same distributed trace).
Similarly, anything that happens after Keptn (but still in the same logical "deployment" operation)
should be included in the end-to-end trace / timing view.
To achieve this, Keptn now accepts tools to pass w3c trace IDs into Keptn
and Keptn will pass W3C trace IDs out of Keptn
(thanks to all contributors working on this for providing this functionality!)
Now, you can see a true end-to-end picture of everything in the logical operation,
potentially from "PR merged" all the way to "deployment complete".

## How to get started

To get started with Keptn, head over to the [Keptn installation instructions](https://keptn.sh/stable/docs/installation/).
After installing Keptn, you can get started by going through the several guides listed below that
will introduce you to the features provided by Keptn step by step.

- [Integrate Keptn with your Applications](https://keptn.sh/stable/docs/guides/integrate/):
This shows you how to configure your applications to be deployed with Keptn.
- [Deployment tasks](https://keptn.sh/stable/docs/guides/tasks/):
This builds up on the previous guide and will explain in detail how you can enhance your application
deployments by adding pre- and post-deployment tasks to ensure you deploy into a healthy environment.
- [Evaluations](https://keptn.sh/stable/docs/guides/evaluations/#annotate-the-workload-resource-for-workload-level-evaluations):
This shows you how to use Keptn to evaluate metrics relevant for an application as part of
the pre- and post-deployment checks.
- [OpenTelemetry Observability](https://keptn.sh/stable/docs/guides/otel/):
This shows how to configure Keptn for making use of its observability features so you
have a holistic overview of your application deployments and any issues that arise during deployment.
- [Multi-stage Application Delivery](https://keptn.sh/stable/docs/guides/multi-stage-application-delivery/):
This guide serves as an example of how you can integrate Keptn with ArgoCD and GitHub to
take the previous concepts and apply them across multiple application environments,
giving you full observability of an application deployment from development into production.

Going through these guides will give you a good initial overview of Keptn,
but be sure to also check out the [additional guides](https://keptn.sh/stable/docs/guides/).
Also, the list of guides is being updated actively by the maintainers,
so feel free to keep an eye on it to see if there is something new added that you could
potentially also apply to your workflow.

[Keptn v2 Release Candidate 1](https://github.com/keptn/lifecycle-toolkit/releases) is available now on GitHub. Please provide any feedback via the #keptn Slack channel in the [CNCF Slack workspace](https://communityinviter.com/apps/cloud-native/cncf).