---
comments: true
---

# Technologies and concepts you should know

You should understand some related technologies
to effectively use and contribute to Keptn.
This section provides links to some materials that can help your learning.
The information has been gathered from the community and is subject to alteration.
If you have suggestions about additional content that should be included in this list,
please submit an issue.

## Containers

Containers are a basic prerequisite in the path of learning kubernetes and a good understanding
of containers will help a lot for working with and contributing to Keptn.

* **Understand the Basics of Containers**
    * [Docker's Official Getting Started Guide](https://docs.docker.com/get-started/)
    * [Docker for beginner](https://youtu.be/3c-iBn73dDE?si=tilXJsYAxZGEMGg1)
* **Container Architecture**
    * [Understanding Docker Architecture](https://docs.docker.com/get-started/overview/#docker-architecture)
* **Container Concepts**
    * [Docker Images and Containers Explained](https://circleci.com/blog/docker-image-vs-container/)
    * [Container Security Explained](https://www.youtube.com/watch?v=b_euX_M82uI)
    * [Container Orchestration Explained](https://youtu.be/kBF6Bvth0zw?si=bUj2bdMbk9xmnF_G)
* **Useful Tools**
    * [Docker Compose Overview](https://docs.docker.com/compose/)
    * [Docker Compose Tutorial](https://youtu.be/SXwC9fSwct8?si=dXaxQVxx0QhW7sku)
    * [Docker Swarm](https://docs.docker.com/engine/swarm/)

## Kubernetes

Keptn runs on Kubernetes and primarily works with deployments that run on Kubernetes
so a good understanding of Kubernetes is essential
for working with and contributing to Keptn.

* **Understand the basics of Kubernetes**
    * [Kubernetes official documentation](https://kubernetes.io/docs/concepts/overview/)
    * [Kubernetes For Beginner](https://youtu.be/X48VuDVv0do)
* **Kubernetes Architecture**
    * [Philosophy](https://youtu.be/ZuIQurh_kDk)
    * [Kubernetes Deconstructed: Understanding Kubernetes by Breaking It Down](https://www.youtube.com/watch?v=90kZRyPcRZw)
* **CRD**
    * [Custom Resource Definition (CRD)](https://www.youtube.com/watch?v=xGafiZEX0YA)
    * [Kubernetes Operator simply explained in 10 mins](https://www.youtube.com/watch?v=ha3LjlD6g7g)
    * [Writing Kubernetes Controllers for CRDs](https://www.youtube.com/watch?v=7wdUa4Ulwxg)
* **Kube-builder Tutorial**
    * [book.kubebuilder.io](https://book.kubebuilder.io/introduction.html)
* **Isitobservable**
    * Keptn has tight integrations with Observability tools and therefore knowing how to _Observe a System_ is important.
    * [Isitobservable website](https://isitobservable.io/)
    * [Is it Observable?
    with Henrik Rexed](https://www.youtube.com/watch?v=aMwk2qo0v40)

## Development tools

* **Go language**
  Most of the Keptn software and many of the test tools
  are written in the Go language.
    * [Go web page](https://go.dev/)
  has tutorials and documentation.
    * [Ginkgo library](https://github.com/onsi/ginkgo/blob/master/README.md)
    is used with the
    [Gomega matcher](https://onsi.github.io/gomega/)
    to implement component tests and end-to-end tests.
* **Chainsaw**
  Some test tools are written with chainsaw
    * [Chainsaw web page](https://kyverno.github.io/chainsaw/)
  has information to get you started
* **Markdown**
  Keptn documentation is authored in Markdown
  and processed with Hugo using the `docsy` theme.
    * [Markdown Guide](https://www.markdownguide.org/)

## Understanding SLOs and SLIs

* **Overview**
    * [Overview](https://www.youtube.com/watch?v=tEylFyxbDLE)
    * [The Art of SLOs (Service Level Objectives)](https://www.youtube.com/watch?v=E3ReKuJ8ewA)

### Operator SDK

* **Go-based Operators**
    * [Go operator tutorial from RedHat](https://docs.okd.io/latest/operators/operator_sdk/golang/osdk-golang-tutorial.html)
