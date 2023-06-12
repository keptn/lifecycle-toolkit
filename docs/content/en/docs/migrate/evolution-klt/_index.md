---
title: Evolution of KLT
description: Understand the paradigm of KLT, how it evolved from Keptn v1, and whether you should migrate to KLT
weight: 5
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

The Keptn products provide tool-agnostic ways
to manage and observe software deployment.
Keptn connects to the tools you choose
through an open event standard
to automate delivery and remediation sequences.

Two different products are available:

* **First Generation:** Keptn v1.y.z LTS, referred to as "Keptn v1"
* **Second Generation:** Keptn Lifecycle Toolkit, referred to as KLT

In this page, we give a high-level overview
of each of these products
to help you understand the similarities and differences
and which is most appropriate for your environment.

> **Note**
This section is under development.
Information that is published here has been reviewed for technical accuracy
but the format and content is still evolving.
We hope you will contribute your experiences
and questions that you have.

## Keptn v1 (Keptn v1.y.z LTS)

Development of Keptn v1 began in 2019.
In 2022, it reached incubation status
and the LTS (Long-Term Support) version was released.
Keptn v1 is designed for the classic monolithic development world,
where many developers work against a single Git repository
on different components that have dependencies within that monolith.
Software applications are developed
using tools like Jenkins to run Continuous Integration.
By building, testing, and validating the entire application
each time a change is made,
this creates a well-defined, monolithic application
that operations could then deploy, observe, and operate.
The complexity of the app was all in the monolithic repository.

Keptn v1 is a general purpose orchestration engine
for cloud and non-cloud native tooling on Kubernetes.
It connects all different types of tools you use in your toolchain
from build to test to deployment,
including notification tools and monitoring tools,
All these tools are connected
through the CloudEvents open standard.

This enables more automated delivery and remediation
based on observability data.
Observability is in the center of Keptn.
Keptn v1 connects different testing and evaluation tools
to execute sequences
that ensure that deployments get into production safely.

The quality gates feature enables you to define
sophisticated analyses of the health of your deployment
at each stage.
They use Service-Level Indicators (SLIs),
which is information provided as queries to your monitoring software
(such as Prometheus, Dynatrace, or DataDog),
and allow you to define Service Level Objectives (SLOs)
which can be based on sophisticated calculations
of multiple SLIs over multiple build/test cycles.
SLIs can be weighted and the results presented
as a general score of the health of your deployment.

Keptn v1 also supports a remediation use case
that can react to a problem such as a Prometheus alert,
and execute sequences of tasks
that bring the system back to a healthy state.
Remediation calls your tools and looks at observability data
from whatever data source you choose --
Prometheus, Dynatrace, DataDog --
to see if the system got back to a healthy state.

## Keptn Lifecycle Toolkit (KLT)

Development of the Keptn Lifecycle Toolkit began in 2022
in response to industry trends and feedback from Keptn users.

Specifically, Cloud Native GitOps has shifted complexity
from devopment to operations.
GitOps promotes a microservice-based paradigm
where each development team runs,configures,
and builds their own Git repository
then uses tools such as Argo and Flux to push
their microservice out to the target environment.
This makes it easier for developers to work on their individual services:
they become more innovative and deploy faster.
but more complex for people who must operate these apps.

Operations is now responsible for determining
exactly what the application is --
is it these five services in a certain version?
Is it another flavor of the application
with a different version of one service?
Who is resposible for the individual tests
as you start to deploy these different services together?

GitOps deploys a service into the environment
but testing, security checks, and validation are not part of GitOps.
Different teams using different tools
and different flavors of Kubernetes
make it more complex for operations
to validate the health of an application,
ensuring that everything is secure, everything is orchestratable,
and everything is observable.

Keptn v1 can define metrics, response time, failure rates,
and use a variety of tools
to validate whether your overall deployment is actually healthy.
But it has flaws in the GitOps world:

- Single-pod ready does not mean application ready.
  If you deploy different services in two different pods,
  and they are ready,
  it does not mean that your overall deployment is actually healthy

- Very difficult to integrate Keptn into all the different tools,
  especially with many development teams using different tools
  because they had to trigger Keptn from each pipeline.
  We saw that a serial integration approach would work better.

- Keptn v1 does not use native Kubernetes CRDs
  to configure SLIs/SLOs, tasks, and other thing
  that are part of the environment where the application deployment lives.
  Instead, it uses its own Git repo and its own configuration files.

- As OpenTelementry and Prometheus have become more prevalent,
  it makes sense to extract the data
  and provide the results in OpenTelemetry and Prometheus
  rather than storing the data in a special Keptn database
  that can be accessed using Keptn APIs.

- Keptn V1 can not define an application
  as a sum of separate microservices,
  each of which may be deployed at different times
  using different tools.
  This means that it cannot validate the full application.

KLT solves these problems by providing an operator
that can observe and orchestrate application-aware workload life cycles.
This operator leverages Kubernetes webhooks
and extends the Kubernetes scheduler
to support pre- and post-deployment hooks.
This means that when any deployment tool
deploys a workload change into Kubernetes,
KLT uses Kubernetes capabilities
to do pre- and post-deployment checks.
This means that:

- Information is observable and is made available
  using OpenTelemetry
- We can provide traces that capture information
  about how long each step takes for the entire deployment cycle
- KLT knows when each workload is being deployed
  and which version of the workload/service
  belongs to this version of the application
- KLT can also provide traces.

- When the operator detects a new version of a service,
  it can execute pre- and post-deployment evaluations and tasks.

## Why to choose KLT

- You are running (or moving towards) GitOps
- Your application is being assembled
  from a set of microservices,
  each of which is developed in its own repo
- You need cloud native observability and tracing
  as well as the ability to add
  pre- and post-deployment evaluations and tasks
  for the deployments you are running

## Why to choose Keptn v1

- You are not using Kubernetes to deploy your application
- You are developing your software as an application
  in a monolithic repository
  rather than assembling it from multiple repos,
  each of which develops a microservice that is part of the application.
- You need the sophisticated SLI/SLO evaluation provided by quality gates,
  including the ability to weight different SLIs
  and receive a composite score
- You want the sophisticated remediation facilities
  that Keptn v1 provides
