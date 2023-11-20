---
title: Technologies and concepts you should know
description: Technologies to get familiar before working with Keptn
weight: 100
---

You should understand some related technologies
to effectively use and contribute to Keptn.
This section provides links to some materials that can help your learning.
The information has been gathered from the community and is subject to alteration.
If you have suggestions about additional content that should be included in this list,
please submit an issue.

## Kubernetes

Keptn runs on Kubernetes and primarily works with deployments that run on Kubernetes
so a good understanding of Kubernetes is essential
for working with and contributing to Keptn.

* **Understand the basics of Kubernetes**
  * [ ] [Kubernetes official documentation](https://kubernetes.io/docs/concepts/overview/)
  * [ ] [Kubernetes For Beginner](https://youtu.be/X48VuDVv0do)
* **Kubernetes Architecture**
  * [ ] [Philosophy](https://youtu.be/ZuIQurh_kDk)
  * [ ] [Kubernetes Deconstructed: Understanding Kubernetes by Breaking It Down](https://www.youtube.com/watch?v=90kZRyPcRZw)
* **CRD**
  * [ ] [Custom Resource Definition (CRD)](https://www.youtube.com/watch?v=xGafiZEX0YA)
  * [ ] [Kubernetes Operator simply explained in 10 mins](https://www.youtube.com/watch?v=ha3LjlD6g7g)
  * [ ] [Writing Kubernetes Controllers for CRDs](https://www.youtube.com/watch?v=7wdUa4Ulwxg)
* **Kube-builder Tutorial**
  * [ ] [book.kubebuilder.io](https://book.kubebuilder.io/introduction.html)
* **Isitobservable**
  * [ ] Keptn has tight integrations with Observability tools and therefore knowing how to _Observe a System_ is important.
  * [ ] [Isitobservable website](https://isitobservable.io/)
  * [ ] [Is it Observable?
    with Henrik Rexed](https://www.youtube.com/watch?v=aMwk2qo0v40)

## Development tools

* **Go language**
  Most of the Keptn software and many of the test tools
  are written in the Go language.
  * [ ] [Go web page](https://go.dev/)
  has tutorials and documentation.
  * [ ] [Ginkgo library](https://github.com/onsi/ginkgo/blob/master/README.md)
    is used with the
    [Gomega matcher](https://onsi.github.io/gomega/)
    to implement component tests and end-to-end tests.
* **KUTTL (KUbernetes Test TooL)**
  Some test tools are written in KUTTL
  * [ ] [KUTTL web page](https://kuttl.dev/)
  has information to get you started
* **Markdown**
  Keptn documentation is authored in Markdown
  and processed with Hugo using the `docsy` theme.
  * [ ] [Markdown Guide](https://www.markdownguide.org/)

## Understanding SLO, SLA, SLIs

* **Overview**
  * [ ] [Overview](https://www.youtube.com/watch?v=tEylFyxbDLE)
  * [ ] [The Art of SLOs (Service Level Objectives)](https://www.youtube.com/watch?v=E3ReKuJ8ewA)

### Operator SDK

* **Go-based Operators**
  * [ ] [Go operator tutorial from RedHat](https://docs.okd.io/latest/operators/operator_sdk/golang/osdk-golang-tutorial.html)
