---
title: Software development environment
description: How to set up and use the development environment to develop and test Keptn software
weight: 30
---

This page gives instructions and hints for setting up a development environment
and then develop, test, and deploy your software changes in that environment.
This material was presented at the
11 September 2023 New Contributors meeting.
You can view the video
[here](https://www.youtube.com/watch?v=UcmULstMYXQ).

To prepare to contribute to the Keptn project, we recommend that you:

* Study the [Keptn documentation](https://lifecycle.keptn.sh/docs/)
  to understand what Keptn does and how it works.
* Familiarize yourself with the
  [lifecycle-toolkit](https://github.com/keptn/lifecycle-toolkit)
  repository, which is the primary repository for
  Keptn software and documentation.
  In particular, study the sections for the four main Keptn components:
  
  * [lifecycle-operator](https://github.com/keptn/lifecycle-toolkit/tree/main/lifecycle-operator)
  * [metrics-operator](https://github.com/keptn/lifecycle-toolkit/tree/main/metrics-operator)
  * [scheduler](https://github.com/keptn/lifecycle-toolkit/tree/main/scheduler)
  * [klt-cert-manager](https://github.com/keptn/lifecycle-toolkit/tree/main/klt-cert-manager)

  Each of these is described in the
  [Architecture](../../../docs/architecture/)
  section of the documentation
  and most include a *README* file with more information.
* Study the material in
  [Technologies and concepts you should know](../../general/technologies).
* Create an account for yourself on
  [GitHub](https://github.com)
  if you do not already have an account.
* Set up a fork of the [lifecycle-toolkit](https://github.com/keptn/lifecycle-toolkit) repository to use with your development.

## View repository

When you view the
[lifecycle-toolkit](https://github.com/keptn/lifecycle-toolkit)
repository, you see that Keptn is composed of multiple components,
each of which is discussed in the Architecture
[Architecture](../../../docs/architecture/)
documentation:

* Three Kubernetes operators
  * `metrics-operator`
  * `lifecycle-operatory`
  * `cert-manager`
* Keptn `scheduler`

At the top level of the repository,
you also see the `runtimes` directory.
This defines the runners that you can use when defining
tasks to be run either pre- or post-deployment.
These are discussed in
[Runners and containers](../../../docs/implementing/tasks.md#runners-and-containers).

## Install software

To test and develop software for the Keptn project,
you need to install the following on your system:

* [**Docker**](https://docs.docker.com/get-docker/): a tool for containerization,
  which allows software applications to run in isolated environments
  and makes it easier to deploy and manage them.
* A Kubernetes cluster running an appropriate version of Kubernetes.
  See [Supported Kubernetes versions](../../../docs/install/reqs.md/#supported-kubernetes-versions)
  for details.
  Most contributors create a local
  Kubernetes-in-Docker(KinD) cluster.
  This is adequate for developing software for Keptn.
  See
  [Kubernetes cluster](../../../docs/install/k8s.md/#create-local-kubernetes-cluster)
  for instructions.
* [**kubectl**](https://kubernetes.io/docs/tasks/tools/):
  a command-line interface tool used for deploying
  and managing applications on Kubernetes clusters.
* [**kustomize**](https://kustomize.io/): a tool used
  for customizing Kubernetes resource configurations
  and generating manifests.
* [**Helm**](https://helm.sh/): a package manager for Kubernetes
  that simplifies the deployment and management of applications
  on a Kubernetes cluster.
* [**Go-lang**](https://go.dev/): the language used to code the Keptn software.

## First steps

1. Follow the instructions in
   [Fork and clone the repository](../../general/git/fork-clone/)
   to get a local copy of the software.

1. Keptn provides a tool that deploys the development version of the software
   on your Kubernetes cluster and pushes the built image to your private repository.
   You identify your private repository with the `RELEASE_REGISTRY=` argument
   and can add any `TAG` arguments you like.
   For example, the following command builds the environment
   and pushes the image to the `docker.io/exampleuser` github repository:

   ```shell
   make build-deploy-dev-environment RELEASE_REGISTRY=docker.io/exampleuser TAG=main
   ```

   The build commands are defined in the
   [Makefile](https://github.com/keptn/lifecycle-toolkit/blob/main/Makefile)
   located in the root directory of your clone.
   This file includes a number of environment variables
   that can be specified as required.

1. After this runs, verify that pods are running on your Kubernetes cluster
   for the four components of the product.

## Code your changes

You are now ready to make your changes to the source code.

1. Follow the instructions in
   [Create local branch](../../general/git/branch-create/)
   to create a branch for your changes.

1. Make your changes to the appropriate component.

1. Deploy the component you modified and push the image to your private Github repository.
   Note that you do not need to rebuild all components,
   only the one you modified.
   For example, if your modifications are to the `metrics-operator`, run:

   ```shell
   make build-deploy-metrics-operator RELEASE_REGISTRY=docker.io/exampleuser TAG=my-feature
   ```

## Testing

Keptn includes a set of tests that are run on each PR that is submitted.
We require that all changes pass
unit tests, component tests, end-to-end tests, and integration tests
before you create a PR with your changes.

If your change introduces a new feature,
you may need to update the test suites to cover your changes.
These tests use basic go-library, Ginkgo or KUTTL tests.
You can ask the maintainers to tell you where to put your additional test data.

Tests are run on your local machine.
Study the detailed log that is produced to identify why the test failed.
Study these errors, modify your code, and rerun the test until it passes.

1. Use your IDE to run unit tests on your code.

1. Run the integration tests from the root directory of your clone:

   ```shell
   make integration-test-local
   ```

   `integration-test-local` cleans up after the test.

1. From the `lifecycle-operator` directory, run the component test:

   ```shell
   make component-test
   ```

1. From the `lifecycle-operator` directory, run the end-to-end tests:

   ```shell
   make e2e-test
   ```

## Create and manage the PR

When all the tests have passed,
you can follow the instructions in
[Create PR](../../general/git/pr-create/)
to create your PR.
Be sure to monitor your PR as discussed in
[PR review process](../../general/git/review/)
until it is merged.
