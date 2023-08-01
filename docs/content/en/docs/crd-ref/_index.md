---
title: API Reference
description: Reference information about the KLT CRDs
weight: 100
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

This section provides comprehensive reference information about all
[Custom Resource Definitions (CRDs)](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
that are defined for the Keptn Lifecycle Toolkit.
This section is auto-generated from source code.

Each CRD is an object of an API library.
Keptn APIs follow the Kubernetes API versioning scheme.
and are themselves composed of objects and sub-objects.
The Keptn CRDs extend the base Kubernetes API
with new objects and functionality.
Keptn APIs follow API versioning conventions recommended by Kubernetes.

Keptn generates most of the resources it needs
without requiring manual input.
[Manifest CRD Reference](../yaml-crd-ref)
contains reference pages for the manifests
that must be populated manually.

Use `kubectl` to inspect the current contents of any Keptn resource:

1. List all resources of the specified type within a certain namespace.
   For example, to list all the `KeptnApp` resources
   in the `namespace1` namespace, the command is:

   ```shell
   kubectl get keptnapps -n namespace1
   ```

1. Get the current manifest for the specified resource.
   For example, to view the manifest for the `my-keptn-app` resource
   in the `namespace1` namespace, the command is:

   ```shell
   kubectl get keptnapp -n <namespace> my-keptn-app -oyaml
   ```

For more information about the APIs and Custom Resources,
see the Kubernetes documentation:

* [API Overview](https://kubernetes.io/docs/reference/using-api/)

* [Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

* [API versioning](https://kubernetes.io/docs/reference/using-api/#api-versioning)

* [Understanding Kubernetes Objects](https://kubernetes.io/docs/concepts/overview/working-with-objects/kubernetes-objects/)
