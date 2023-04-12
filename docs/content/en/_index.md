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

<!-- markdownlint-disable MD033 -->
<!-- markdownlint-disable-next-line MD013 -->
{{< blocks/cover title="Welcome to the Keptn Lifecycle Toolkit Documentation" image_anchor="top" height="half" color="primary" >}}
<div class="mx-auto">
 <a class="btn btn-lg btn-primary mr-3 mb-4" href="{{< relref "/docs" >}}">
  Docs <i class="fas fa-arrow-alt-circle-right ml-2"></i>
 </a>
    <a class="btn btn-lg btn-primary mr-3 mb-4" href="https://github.com/keptn/lifecycle-toolkit/releases">
  Releases <i class="fab fa-github ml-2 "></i>
 </a>
</div>
{{< /blocks/cover >}}
<!-- markdownlint-enable MD033 -->

{{% blocks/lead color="white" %}}
[![Keptn Lifecycle Toolkit in a Nutshell](https://img.youtube.com/vi/K-cvnZ8EtGc/0.jpg)](https://www.youtube.com/watch?v=K-cvnZ8EtGc)
{{% /blocks/lead %}}
git add
{{< blocks/section color="dark" >}}
{{% blocks/feature icon="fa-lightbulb" title="Keptn Recordings" %}}
See Keptn [in Action](https://youtube.com/playlist?list=PL6i801Rjt9DbikPPILz38U1TLMrEjppzZ)
{{% /blocks/feature %}}

{{% blocks/feature icon="fab fa-github" title="Contributions welcome!" url="https://github.com/keptn/lifecycle-toolkit" %}}
We do a [Pull Request](https://github.com/keptn/lifecycle-toolkit/pulls) contributions workflow on **GitHub**.
New users are always welcome!
{{% /blocks/feature %}}

{{% blocks/feature icon="fab fa-twitter" title="Follow us on Twitter!" url="https://twitter.com/keptnProject" %}}
For announcement of latest features etc.
{{% /blocks/feature %}}

{{< /blocks/section >}}
