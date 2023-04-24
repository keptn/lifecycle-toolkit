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
<!-- markdownlint-enable no-inline-html -->

{{% blocks/lead color="primary" %}}

## Use Cases

We extend the K8s Pod Scheduler with following Use Cases:

{{% /blocks/lead %}}

{{< blocks/section color="gray" >}}
{{% blocks/feature icon="fa-lightbulb" title="Deployment Observability" %}}
making ANY K8s Deployment OBSERVABLE

[read more](#deployment-observability)

{{% /blocks/feature %}}
{{% blocks/feature icon="fa-lightbulb" title="Deployment Data Access" %}}
standardizing access for all Observability Data for K8s

[read more](#data-access)
{{% /blocks/feature %}}

{{% blocks/feature icon="fa-lightbulb" title="Orchestrate Deployment Checks" %}}
orchestrating deployment checks as part of scheduler

[read more](#deployment-check-orchestration)

{{% /blocks/feature %}}

{{< /blocks/section >}}

{{% blocks/lead color="blue" %}}
[![Keptn Lifecycle Toolkit in a Nutshell](https://img.youtube.com/vi/K-cvnZ8EtGc/0.jpg)](https://www.youtube.com/watch?v=K-cvnZ8EtGc)
{{% /blocks/lead %}}

<!-- markdownlint-disable no-inline-html -->
{{% blocks/lead color="primary" title="" %}}

## Deployment Observability

### making ANY K8s Deployment OBSERVABLE

<div style="text-align: left">

If you deploy with ArcoCD, Flux, GitLab, kubectl ...
we provide you:

* Automated **App-Aware DORA** metrics (OTel Metrics)
* Troubleshoort failed deployments  (OTel Traces)
* Trace deployments from git to cloud  (traces accross stages)

</div>
<a class="btn -bg-green" href="./docs/quickstart/">
    Get Started!
</a>
{{% /blocks/lead %}}
<!-- markdownlint-enable no-inline-html -->

<!-- markdownlint-disable no-inline-html -->
{{% blocks/lead color="white" %}}

## Data Access

### standardizing access for all Observability Data for K8s

<div style="text-align: left">

To drive decisions based on metrics use Keptn Metrics Server to:

* Define Keptn Metrics once for Dynatrace, DataDog, AWS, Azure, GCP, ...
* Access all those metrics via Prometheus or K8s Metric API
* Eliminate the need of multiple plugins for Argo Roolouts, KEDA, HPA, ...

</div>
<a class="btn -bg-green" href="./docs/quickstart/">
    Get Started!
</a>
{{% /blocks/lead %}}
<!-- markdownlint-enable no-inline-html -->

<!-- markdownlint-disable no-inline-html -->
{{% blocks/lead color="primary" %}}

## Deployment Check Orchestration

### orchestrating deployment checks as part of scheduler

<div style="text-align: left">

To reduce complexity of custom checks use Keptn to:

* Pre-Deploy:
  * validate external dependenicies
  * confiorm images are scanned
  * ...
* Post-Deploy:
  * Execute tests
  * Notify Stakeholders
  * Promote to next stage
  * ...
* Automatically validate against your SLO (Service Level Objectives)

</div>
<a class="btn -bg-green" href="./docs/quickstart/">
    Get Started!
</a>
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
