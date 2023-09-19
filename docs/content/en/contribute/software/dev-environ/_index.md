---
title: Set up the development environment
description: How to set up an environment to develop and test Keptn software
weight: 30
---

This page gives instructions and hints for setting up an environment
you can use to develop, test, and deploy your software changes.
This material was presented at the
11 September 2023 New Contributors meeting.
You can view the video
[here](https://www.youtube.com/watch?v=UcmULstMYXQ).

To prepare to contribute to the Keptn project, we recommend:

* Familiarize yourself with the
  [lifecycle-toolkit](https://github.com/keptn/lifecycle-toolkit)
  repository, which is the primary repository for
  Keptn software and documentation.
* Study the [Keptn documentation](https://lifecycle.keptn.sh/docs/)
  to understand what Keptn does and how it works.
* Study the material about
  [Related Technologies](../../general/technologies)
  that you need to understand to develop Kubernetes software.

## View repository

When you view the
[lifecycle-toolkit](https://github.com/keptn/lifecycle-toolkit)
repository, you see that Keptn is composed of multiple components,
each of which is discussed in the Architecture documentation:

* Two Kubernetes operators, `metrics-operator` and `lifecycle-operatory`
* Keptn `scheduler`
* Keptn `cert-manager`

At the top level of the repository,
you also see the `runtimes` directory.
This defines the runners that you can use when defining
tasks to be run either pre- or post-deployment.

## Install software

To test and develop software for the Keptn project,
you need to install the following on your system:

* [**Docker**](https://docs.docker.com/get-docker/): a tool for containerization,
which allows software applications to run in isolated environments
and makes it easier to deploy and manage them.
* A Kubernetes cluster running an appropriate version of Kubernetes.
  See [Supported Kubernetes versions](../../../docs/install/reqs.md/#supported-kubernetes-versions)
  for details.
  If you need to set up a local Kubernetes cluster
  we recommend Kubernetes-in-Docker(KinD).
  This is adequate for developing software for Keptn.
* [**kubectl**](https://kubernetes.io/docs/tasks/tools/):
  a command-line interface tool used for deploying
  and managing applications on Kubernetes clusters.
* [**kustomize**](https://kustomize.io/): a tool used
  for customizing Kubernetes resource configurations
  and generating manifests.
* [**Helm**](https://helm.sh/): a package manager for Kubernetes
  that simplifies the deployment and management of applications
  on a Kubernetes cluster.
* [**Go-lang](https://go.dev/): the language used to code the Keptn software.
