---
title: Kubernetes cluster
description: Bring or install a Kubernetes cluster 
icon: concepts
layout: quickstart
weight: 25
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

The Keptn Lifecycle Toolkit is meant to be installed
into an existing Kubernetes cluster
that runs your deployment software.
See [Requirements](reqs.md) for information about supported releases
and advice about resources required.

You can also create a local cluster using packages
such as KinD, Minikube, K3s, and K3d
that can be used for testing, study, and demonstration purposes.

## Create local Kubernetes cluster

You can use tools such as
[KinD](https://kind.sigs.k8s.io/),
[k3d](https://k3d.io/),
[k3s](https://k3s.io/),
and [Minikube](https://minikube.sigs.k8s.io/docs/)
to set up a local, lightweight Kubernetes cluster
where you can install the Keptn Lifecycle Toolkit
for personal study, demonstrations, and testing.

The [Keptn Lifecycle Toolkit: Installation and KeptnTask Creation in Minutes](https://www.youtube.com/watch?v=Hh01bBwZ_qM)
video  demonstrates how to create a KinD cluster.
on which you can install the Lifecycle Toolkit.
The basic steps are:

1. Download, install, and run [Docker](https://docs.docker.com/get-docker/)
1. Download [KinD](https://kind.sigs.k8s.io/)
1. Create the local KinD cluster with the following command:

   ```shell
   kind create cluster
   ```

1. When the cluster has been created,
   run the following to verify that the cluster is working
   and that it is running a supported version of Kubernetes
   with the following command:

   ```shell
   kubectl version --short
   ```
