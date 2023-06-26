---
title: Migration strategy
description: General guidelines for migrating your deployment to KLT
weight: 10
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

> **Note**
This section is under development.
Information that is published here has been reviewed for technical accuracy
but the format and content is still evolving.
We hope you will contribute your experiences
and questions that you have.

No two migrations are alike but these are some general guidelines
for how to approach the project:

1. Create a new Kubernetes cluster that is not running anything else
   and use that to build out your deployment environment.
1. Install and configure the deployment tool(s) you want to use
   to deploy the components of your software.
   You can use different deployment tools for different components.
1. Install KLT.

TODO: For migration, is it best to try to do all the applications at once
or should they do one application then move onto the next?

TODO: Similarly, would they start out doing one pre-deployment evaluation,
then one post-deployment evaluation, then one pre-deployment task,
then one post-deployment task?
Or would you just dig in and do them all?

1. Integrate KLT with your applications:
   - If you are only using Custom Metrics and Observability,
     you only need to do basic annotations.
   - If you want to do pre- and post-deployment evaluations and tasks,
     include those annotations as well.
   - In all cases, define a KeptnApp resource
     for each application you are deploying.
     You can do this manually but the easiest approach
     is to use the Keptn automatic app discovery feature
     to create the basic KeptnApp resource.
     You can then manually modify that resource as needed
     but easily create the basic structure.

1. Set up Keptn Metrics and Observability.
   - Configure your data sources as `KeptnMetricProviders`.
   - Install and configure OpenTelemetry.
   - Install and configure the data providers you want to use.
     Prometheus, Dynatrace, and Datadog are currently supported.
     KLT can access multiple data providers
     and multiple copies of each data provider.
   - Populate a KeptnMetricsProvider for each copy
     of each data provider you are using.
   - Run some deployments and use the metrics and observability features
     to monitor the process
1. Set up CI/CD tasks
1. Convert "simple" quality gates into KeptnMetric
   and KeptnEvaluation resources.
   TODO: Yeah, this will be a big section!
   - TODO: Is there a way to modify queries on the data provider
     or maybe implement a task that compares some values and
     does some math to migrate more complex quality gates?
1. Convert Keptn v1 remediation sequences into "Day 2"
   KeptnEvaluations and KeptnTaskDefinition resources.
