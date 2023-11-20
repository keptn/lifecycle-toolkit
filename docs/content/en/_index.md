---
title: Keptn Kubernetes Orchestration - Supercharge Your Deployments
cascade:
  - _target:
      path: "/docs/**"
    sitemap:
      priority: 1.0

  - _target:
      path: "/docs-*/**"
    sitemap:
      priority: 0.1

  - _target:
      path: "/community/**"
    type: docs
  - _target:
      path: "/community/readme.md"
    draft: true
  - _target:
      path: "/community/_index.md"
    title: Community
    menu:
      main:
        weight: 20
---

<!-- markdownlint-disable no-inline-html -->
<!-- markdownlint-disable-next-line line-length -->
{{< blocks/cover title="" image_anchor="top" height="full" color="primary" >}}
<div class="mx-auto">
<div class="mb-4 d-none d-md-block " >
<picture >
    <img alt="keptn" src="/images/home/hero/keptn-logo-white.svg">
</picture>
</div>
  <h2 class="styled display-3 mt-0 mb-4">Cloud native application lifecycle orchestration </h2>
<div>
 <a class="btn btn-lg btn-primary mr-3 mb-4" href="{{< relref "/docs" >}}">
  Docs <i class="fas fa-arrow-alt-circle-right ml-2"></i>
 </a>
    <a class="btn btn-lg btn-secondary mr-3 mb-4" href="{{< relref "/docs/getting-started/" >}}">
  Get Started <i class="fa-solid fa-rocket ml-2"></i>
 </a>
</div>
</div>
<div class="usecasebox">

## Use Cases

We extend the K8s APIs with the following Use Cases:

<div class="row usecases">
{{% blocks/feature icon="home homeobservability" title="Deployment Observability" %}}
Make ANY Kubernetes Deployment observable

<!-- markdownlint-disable-next-line link-fragments -->
 <a class="btn btn-lg -bg-light mr-3 mb-4" href="#deployment-observability">
  read more <i class="fas fa-arrow-alt-circle-down ml-2"></i>
 </a>
{{% /blocks/feature %}}
{{% blocks/feature icon="home homeorchestrate" title="Gather metrics from anywhere" %}}
Standardize access for all Observability Data for K8s

<!-- markdownlint-disable-next-line link-fragments -->
 <a class="btn btn-lg -bg-light mr-3 mb-4" href="#gather-metrics-from-anywhere">
  read more <i class="fas fa-arrow-alt-circle-down ml-2"></i>
 </a>
{{% /blocks/feature %}}

{{% blocks/feature icon="home homedata" title="Orchestrate Deployment Checks" %}}
Gain confidence in your work with pre-/post-deployment checks

<!-- markdownlint-disable-next-line link-fragments -->
 <a class="btn btn-lg -bg-light mr-3 mb-4" href="#orchestrate-deployment-checks">
  read more <i class="fas fa-arrow-alt-circle-down ml-2"></i>
 </a>
{{% /blocks/feature %}}

</div>
</div>

{{< /blocks/cover >}}
<!-- markdownlint-enable no-inline-html -->

{{% blocks/lead color="light" %}}

{{< youtube K-cvnZ8EtGc >}}

{{% /blocks/lead %}}

<!-- markdownlint-disable no-inline-html -->
{{% blocks/lead color="white"%}}
<div class="mx-auto">
<div class="d-flex flex-row flex-wrap" >
<div class="whykeptn whykeptn-left">
{{% readfile "partials/_index-observability-left.md" %}}
</div>
<div class="whykeptn whykeptn-right w-25">
{{% readfile "partials/_index-observability-right.md" %}}
</div>
</div>
<a class="btn -bg-green" href="./docs/intro/#observability">
    Get Started!
</a>
</div>

{{% /blocks/lead %}}
<!-- markdownlint-enable no-inline-html -->

<!-- markdownlint-disable no-inline-html -->
{{% blocks/lead color="light" %}}
<div class="mx-auto">
<div class="d-flex flex-row flex-wrap" >
<div class="whykeptn whykeptn-left w-25">
{{% readfile "partials/_index-gather-metrics-left.md" %}}
</div>
<div class="whykeptn whykeptn-right ">
{{% readfile "partials/_index-gather-metrics-right.md" %}}
</div>
</div>
<a class="btn -bg-green" href="./docs/intro/#metrics">
    Get Started!
</a>
</div>
{{% /blocks/lead %}}
<!-- markdownlint-enable no-inline-html -->

<!-- markdownlint-disable no-inline-html -->
{{% blocks/lead color="white" %}}
<div class="mx-auto">
<div class="d-flex flex-row flex-wrap" >
<div class="whykeptn whykeptn-left">
{{% readfile "partials/_index-deployment-checks-left.md" %}}
</div>
<div class="whykeptn whykeptn-right w-25 text-center">
{{% readfile "partials/_index-deployment-checks-right.md" %}}
</div>
</div>
<a class="btn -bg-green" href="./docs/intro/#release-lifecycle-management">
    Get Started!
</a>
</div>

{{% /blocks/lead %}}
<!-- markdownlint-enable no-inline-html -->

{{< blocks/section color="dark" >}}
{{% blocks/feature icon="fa-lightbulb" title="Keptn Recordings" %}}
See Keptn in Action

<!-- markdownlint-disable-next-line no-inline-html -->
<a class="btn -bg-white rounded-lg" href="https://youtube.com/playlist?list=PL6i801Rjt9DbikPPILz38U1TLMrEjppzZ">
  Watch now!
 </a>
{{% /blocks/feature %}}

{{% blocks/feature icon="fab fa-github" title="Contributions welcome!" %}}
We do a [Pull Request](https://github.com/keptn/lifecycle-toolkit/pulls) contributions workflow on **GitHub**.
New users are always welcome!

<!-- markdownlint-disable-next-line no-inline-html -->
<a class="btn -bg-white rounded-lg" href="https://github.com/keptn/lifecycle-toolkit">
  Contribute on GitHub
 </a>
{{% /blocks/feature %}}

{{% blocks/feature icon="fab fa-twitter" title="Follow us on Twitter!" %}}
For announcement of latest features etc.

<!-- markdownlint-disable-next-line no-inline-html -->
<a class="btn -bg-white rounded-lg" href="https://twitter.com/keptnProject">
  Follow us!
 </a>
{{% /blocks/feature %}}

{{< /blocks/section >}}
