---
comments: true
---

# Context metadata

This guide walks you through the usage of context metadata in Keptn.

After reading this guide you will be able to:

- add metadata to applications, workloads
  and tasks in two ways
- add metadata information about applications to their traces

## Introduction

Context metadata in Keptn is defined as all additional information that can
be added to applications, workloads and tasks.

Context metadata can be added to your application
in either of the following ways:

- Use the `keptn.sh/metadata` annotation in any of
your workloads you plan to deploy with Keptn
- Apply a `KeptnAppContext` custom resource

In the following section we will explore both.

Common uses of this feature include:

- adding a commit ID or references related to your GitOps and CI-CD tooling
- referencing different stages and actors
  such as: who committed the change, in what repo, for what ticket...

### Before you start

1. [Install Keptn](../installation/index.md)
2. Deploy an application, for instance, you can follow
   [a demo app installation here](../getting-started/observability.md#step-3-deploy-demo-application)

To collect traces you will require Jaeger.
To visualise and inspect the traces, you can either use
the Jaeger UI or Grafana.

1. Install
   [Grafana](https://grafana.com/grafana/)
   following the instructions in [Grafana Setup](https://grafana.com/docs/grafana/latest/setup-grafana/)
   or the visualization tool of your choice.
2. Install
   [Jaeger](https://www.jaegertracing.io/)
   or a similar tool for traces following the instructions in
   [Jaeger Setup](https://www.jaegertracing.io/docs/1.50/getting-started/).

## Include metadata in workload traces

To enrich workload traces with custom metadata, use the
`keptn.sh/metadata` annotation in your
[Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment)
resource.
The comma-separated list of key-value pairs specified in the annotation
is added as key-value attributes to the workload trace.

Modify your workload (for example your YAML file containing a Deployment)
adding an annotation with any metadata you prefer.
If you want to add multiple key-value pairs, separate them with commas.
For instance, to add information about a stage and owning team, you could add:
`keptn.sh/metadata: "stage=dev,team=my-team"`.

To see the changes Keptn must redeploy: increment the `app.kubernetes.io/version` value
(ex. if you are following our getting started guide, change the version
from `0.0.2` to `0.0.3`) or change the `keptn.sh/version` value
if you used the Keptn specific labels in your deployment YAML file.

This way, after the re-deployment, the workload trace will contain the `stage=dev` attribute.

## Include metadata in application traces

Similar to the previous step, custom metadata can also be added to application traces via the
[KeptnAppContext](../reference/api-reference/lifecycle/v1/index.md#keptnappcontext) CRD.

`KeptnAppContext` is a custom resource definition in Keptn that allows you to add metadata
and links to traces for a specific application.
This enables you to enrich your Keptn resources and your traces with additional
information, making it easier to understand and analyze
the performance of your applications.

In the `.spec.metadata` field you can define multiple key-value pairs, which are propagated
to the application trace as attributes in the same manner as for workloads.

> **Note** The key-value pairs that are added to the application trace are also added
to each workload trace that is part of the application.
If the same key is specified for both
application and workload metadata attributes,
values specified for the workload take precedence.

A `KeptnAppContext` custom resource looks like the following:

```yaml
{% include "./assets/metadata/keptn-app-context.yaml" %}
```

After applying the `KeptnAppContext` to your cluster, you need to increment the version of your
application by modifying your deployment file and changing the
value of the`app.kubernetes.io/version` field (or `keptn.sh/version` if you used the Keptn specific labels earlier).
(Alternatively you could apply the context resource before or together with the workloads.)

After deploying the `KeptnAppContext` resource and re-deploying the application,
Keptn triggers another deployment of your application with the new context metadata,
and all traces will contain the new metadata as defined in the above code example.
In other words, you should be able to see the application trace as well as the workload trace
contain the defined metadata as key-value attributes.

## What's next?

Congratulations!
You've learned how to use the `KeptnAppContext` CRD to add
metadata to applications and their traces.
This can be valuable for understanding the context of your traces and
establishing connections between
different versions and stages of your application.

Explore more about [traces on Keptn](./otel.md).
The paragraph on
[linking traces between different application](./otel.md#advanced-tracing-configurations-in-keptn-linking-traces)
also uses KeptnAppContext.
