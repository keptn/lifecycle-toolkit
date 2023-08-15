---
title: Set up the development environment
description: How to set up an environment to develop and test Keptn software
weight: 30
---

To test and develop software for the Keptn project,
you need to install the following on your system:

* [**Docker**](https://docs.docker.com/get-docker/): a tool for containerization,
which allows software applications to run in isolated environments
and makes it easier to deploy and manage them.
* A Kubernetes cluster running an appropriate version of Kubernetes.
  See [Supported Kubernetes versions](docs/install/reqs/#supported-kubernetes-versions)
  for details.
  If you need to set up a local Kubernetes cluster
  we recommend Kubernetes-in-Docker(kind).
* [**kubectl**](https://kubernetes.io/docs/tasks/tools/):
  a command-line interface tool used for deploying
  and managing applications on Kubernetes clusters.
* [**kustomize**](https://kustomize.io/): a tool used
  for customizing Kubernetes resource configurations
  and generating manifests.
* [**Helm**](https://helm.sh/): a package manager for Kubernetes
  that simplifies the deployment and management of applications
  on a Kubernetes cluster.
