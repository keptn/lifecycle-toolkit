---
comments: true
---

# Word list

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

## Keptn terminology

### Keptn project

This is the proper name of the project that was developed
under the code name of "Keptn Lifecycle Toolkit (KLT)".
The earlier project is called "Keptn v1".

* Keptn is capitalized when used in prose as the name of the project.
  Use `keptn` if it is part of a command name, pathname,
  an argument to a command or function, etc.

* As a project name that is trademarked,
  you should not use an apostrophe-s to make it a possessive ("Keptn's")
  or hyphenate it (as in "Keptn-specific").

The Keptn project is a "toolkit" with three use cases, named:

* Metrics
* Observability
* Release lifecycle management

## Kubernetes terminology

The Keptn documentation frequently uses Kubernetes terminology.
Here are some guidelines for using these terms in Keptn documentation.

* Do not duplicate information from the Kubernetes documentation
  into the Keptn documentation.
  We should not be maintaining documentation
  for software in other projects and products.

* Kubernetes concepts and objects (such as workload or resource)
  should be lowercase unless they are the proper name of an object.

    This means that "workload" is not capitalized, but "Pod" and "Deployment" are.

* The first instance of one of these terms in a section
  should be a link to the relevant Kubernetes documentation.

* Avoid using generic references to Kubernetes documentation.
  Instead, link to the particular section
  that contains the relevant information.

* The dictionary of Kubernetes terms that is used by the
  Spell checker
  is in the `cspell`
  [k8s.txt](https://github.com/check-spelling/cspell-dicts/blob/main/dictionaries/k8s/dict/k8s.txt)
  file.

### CRD, resource, etc

Keptn makes extensive use of Kubernetes
[Custom resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).
It is important to use the related terminology correctly:

* A "Resource Definition (RD)" is the definition (or syntax)
  of a resource that is part of the official Kubernetes API

* A "Custom Resource Definition (CRD)" is the definition
  (or syntax) of a resource that Keptn (or some other software)
  adds to Kubernetes

* An instance of a CRD or RD that a user creates is a custom resource, CR
  or just a resource but not a CRD or RD.
  Most of the time, we recommend just using the term "resource".

* The first occurrence of a CRD name in a section should be a link to the
  [CRD Reference](../../reference/crd-reference/index.md) page if there is one.
  Otherwise, it should be a link to the appropriate spot in the
  [API Reference](../../reference/api-reference/index.md) section.

* Occurrences of a resource name that are not links to a reference page
  should be enclosed in back-ticks, so that they render as inline `code-case`.
