---
title: Home
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
{{< blocks/cover title="Welcome to the Keptn Lifecycle Toolkit" image_anchor="top" height="half" color="primary" >}}
<div class="mx-auto">
 <a class="btn btn-lg -bg-green mr-3 mb-4" href="{{< relref "/docs" >}}">
  Docs <i class="fas fa-arrow-alt-circle-right ml-2"></i>
 </a>
    <a class="btn btn-lg btn-secondary mr-3 mb-4" href="https://github.com/keptn/lifecycle-toolkit/releases">
  Releases <i class="fab fa-github ml-2 "></i>
 </a>
</div>

[![Keptn Lifecycle Toolkit in a Nutshell](https://img.youtube.com/vi/K-cvnZ8EtGc/0.jpg)](https://www.youtube.com/watch?v=K-cvnZ8EtGc)
{{< /blocks/cover >}}
<!-- markdownlint-enable no-inline-html -->

{{% blocks/lead color="secondary" %}}

## Use Cases

We extend the K8s APIs with the following Use Cases:

{{% /blocks/lead %}}

{{< blocks/section color="light" >}}
{{% blocks/feature icon="home homeobservability" title="Deployment Observability" %}}
making ANY K8s Deployment OBSERVABLE

<!-- markdownlint-disable-next-line link-fragments -->
[read more](#deployment-observability)

{{% /blocks/feature %}}
{{% blocks/feature icon="home homedata" title="Deployment Data Access" %}}
standardizing access for all Observability Data for K8s

<!-- markdownlint-disable-next-line link-fragments -->
[read more](#data-access)
{{% /blocks/feature %}}

{{% blocks/feature icon="home homeorchestrate" title="Orchestrate Deployment Checks" %}}
orchestrating deployment checks as part of scheduler

<!-- markdownlint-disable-next-line link-fragments -->
[read more](#deployment-check-orchestration)

{{% /blocks/feature %}}

{{< /blocks/section >}}

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
<a class="btn -bg-green" href="./docs/getting-started/">
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
{{% readfile "partials/_index-data-access-left.md" %}}
</div>
<div class="whykeptn whykeptn-right ">
{{% readfile "partials/_index-data-access-right.md" %}}
</div>
</div>
<a class="btn -bg-green" href="./docs/getting-started/">
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
<a class="btn -bg-green" href="./docs/getting-started/">
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
