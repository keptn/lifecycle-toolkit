---
title: Word list
description: Proper use of terms for the Keptn project documentation
weight: 800
---

This document summarizes information
about the proper use of terminology for the Keptn project.

The Keptn project does not maintain a formal style guide
but should adhere to recommendations in:

* [CNCF Style Guide](https://github.com/cncf/foundation/blob/main/style-guide.md).

* The [Kubernetes documentation](https://kubernetes.io/docs/home/)
  is a good reference for Kubernetes terms.
   In particular:

  * [Concepts](https://kubernetes.io/docs/concepts/)
  * [Reference](https://kubernetes.io/docs/reference/)

* The [Google developer documentation style guide](https://developers.google.com/style)
  is a comprehensive reference for technical writers.
   In particular:

  * [Word list](https://developers.google.com/style/word-list)
    includes good information about words and phrases
    that are commonly used in technical documentation

## Keptn project

This is the proper name of the project that was developed
under the code name of "Keptn Lifecycle Toolkit (KLT)".
The earlier project is called "Keptn v1".

* Keptn is capitalized when used in prose as the name of the project,
  although the logo uses a lowercase "k".
  Use `keptn` if it is part of a command name, pathname,
  an argument to a command or function, etc.

* As a project name that is trademarked,
  you should not use an apostrophe-s to make it a possessive ("Keptn's")
  or hyphenate it (as in "Keptn-specific").

The Keptn project is a "toolkit" with three use cases, named:

* Metrics (or Deployment data access)

* Observability (or Deployment observability)

* Release lifecycle management (or Orchestrate deployment checks)

## CRD, resource, etc

Keptn makes extensive use of Kubernetes
[Custom resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).
It is important to use the related terminology correctly:

* A "Resource Definition (RD)" is the definition (or syntax)
  of a resource that is part of the official Kubernetes API

* A "Custom Resource Definition (CRD)" is the definition
  (or syntax) of a resource that Keptn (or some other product)
  adds to Kubernetes

* An instance of a CRD or RD that a user creates is a custom resource
  or just a resource but not a CRD or RD.
  Most of the time, we recommend just using the term "resource".

* The first occurence of a CRD name in a section should be a link to the
  [CRD YAML Reference](../../../docs/yaml-crd-ref)
  page if there is one.
  Otherwise, it should be a link to the appropriate spot in the
  [API Reference](../../../docs/crd-ref)
  section.

* Occurrences of a resource name that are not links to a reference page
  should be enclosed in tics so they render as code-case.
