---
title: Standardize access to observability data
description: The Keptn Lifecycle Toolkit makes any Kubernetes deployment observable.
weight: 45
---

The Keptn Lifecycle Toolkit makes any Kubernetes deployment observable.
You can readily see why a deployment takes so long or why it fails,
even when using multiple deployment tools.
Keptn provides a Keptn application CRD
that aggregates multiple Workloads into a single resource
that can be monitored.

The Lifecycle Toolkit emits OpenTelemetry traces
for everything that happens in the Kubernetes pod scheduler
and can display this information on dashboard tools
such as Grafana and Jaeger.
DORA metrics are implemented out of the box
when the Lifecycle Toolkit is enabled
and you can define additional metrics you want to monitor.

This exercise builds on what you did in the
[Getting started with Keptn metrics](../metrics)
exercise:

* Install the Keptn Lifecycle Toolkit on your deployment cluster
* Enable the Lifecycle Toolkit for your cluster
  and integrate it into your cluster
* Define `KeptnMetrics` for the data you want to monitor.

This exercise shows how to standardize access
to the observability data for your cluster.
The steps are:

* Configure OpenTelemetry for the Lifecycle Toolkit
* Define the workloads to be included in your Keptn Application.
  You can specify this manually or use the application discovery feature
  to automatically generate a Keptn Application
  that includes all workloads on the cluster,
  regardless of the tools being used
  that includes all workloads running in the cluster,
* If you like, define `KeptnMetrics` for additional data you want to monitor.
* Increment the version number for either your Workload
  or your application to start aggragating data
* View the aggregated metrics and traces on Grafana
  or the dashboard of your choice

## Integrate KLT with your applications

Integrate the Lifecycle Toolkit with your applications
by annotating the Kubernetes `Deployment` and `Namespace` CRDs

### Annotate Deployment CRD

[Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/).
For our example, modify the simplenode-dev-deployment.yaml file
and add the following lines to the `template.metadat.annotations` section:

```yaml
...
app.kubernetes.io/name: simplenodeservice
```

If you have multiple deployments or stateful sets,
use the following syntax to reference the
[KeptnApp](../../yaml-crd-ref/app.md)
CRD:

```yaml
app.kubernetes.io/part-of: simpleapp
```

`simpleapp` is the name assigned to a
[KeptnApp](../../yaml-crd-ref/app.md)

Could instead use the Keptn annotations:

```yaml
keptn.sh/app: simpleapp
keptn.sh/workload: simplenode
keptn.sh/version: x.y.z
```

The version number is an arbitrary version for the Workload
and you can use any sort of numbering system you like.
Incrementing this value triggers a new execution of the Workload.

TODO: Should we discuss that KeptnApp has a version for the app
as well as a version for each workload that is in the app or
should we just link to the ref page for those details.

TODO: Is this version number automatically sync'ed with
the application version in the KeptnApp CRD?
If I have KeptnApp defined, should I change version number
in KeptnApp or in Workload or does it matter?

### Annotate Namespace CRD

TODO: This is the first time we've mentioned Namespace.
Should we mention something about setting the namespace and all that?
Maybe just link to
[Kubernetes Namespaces](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/).

Annotate the
[Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
CRD to tell the webhook to handle the namespace:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: simplenode-dev
  annotations:
    kept.sh/lifecycle-toolkit: "enabled"
```

This activates the Lifecycle toolkit on your
where you want Keptn to become active
and make deployments observable.activates the Lifecycle toolkit on your
where you want Keptn to become active
and make deployments observable.

## Start making deployments observable

To start feeding observability data for your deployments
onto a dashboard of your choice,
modify either your Deployment or KeptnApp CRD yaml file
to increment the version number
and commit that change to your repository.

## View the results

TODO: talk about the Grafana display.

If you also have Jaeger extension for Grafana installed on your cluster,
you can view full end-to-end trace for everything
that happens in your deployment.

## For more information

For more information, see
