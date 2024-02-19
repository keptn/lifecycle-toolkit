---
comments: true
---

# Auto app discovery

The automatically generated `KeptnApp` file
aggregates the workloads to include in the application,
based on annotations made to the workloads themselves.
This enables you to run Keptn observability features on your cluster.

Afterward, you can monitor the status of the deployment using
a command like the following:

```shell
kubectl get keptnworkloadversion -n podtato-kubectl -w
```

The generated `KeptnApp` file includes `metadata`
that names this `KeptnApp` and identifies the Namespace where it resides.

```yaml
{% include "./assets/auto-app-discovery/keptnapp-metadata.yaml" %}
```

It also includes a `spec.workloads` list
that defines the workloads to be included.

As an example, consider the following application,
consisting of multiple deployments,
which is going to be deployed into a Keptn-enabled namespace.
Note that:

1. Keptn is enabled for the namespace where your application runs.
1. The `Deployment` workloads are annotated appropriately.
   This example does not use other workloads.

```yaml
{% include "./assets/auto-app-discovery/deployments.yaml" %}
```

Applying these resources results in the creation
of the following `KeptnApp` resource:

```yaml
{% include "./assets/auto-app-discovery/keptnapp.yaml" %}
```

With the `KeptnApp` resource created,
you get observability of your application's deployments
by using the OpenTelemetry tracing features
that are provided by Keptn:

![Application deployment trace](./assets/trace.png)

To execute pre-/post-deployment checks for a `KeptnApp`,
create a `KeptnAppContext` with the same name and in the same `namespace` as the `KeptnApp`.
The `KeptnAppContext` contains a list of
pre-/post-deployment tasks, evaluations, and promotion tasks
that should be executed before and after the
workloads within the `KeptnApp` are deployed.

See the [Getting started guide](../getting-started/lifecycle-management.md#more-control-over-the-application)
for more information on how to configure a `KeptnAppContext`
to execute pre-/post-deployment checks or promotion tasks.
