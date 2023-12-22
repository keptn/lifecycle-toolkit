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

{% include "./assets/auto-app-discovery.md_1.yaml" %}

It also includes a `spec.workloads` list
that defines the workloads to be included.

Pre-/post-deployment tasks and evaluations for the `KeptnApp`
can also be added to this resource manually
but this is not required for observability.

As an example, consider the following application,
consisting of multiple deployments,
which is going to be deployed into a Keptn-enabled namespace.
Note that:

1. Keptn is enabled for the namespace where your application runs.
1. The `Deployment` workloads are annotated appropriately.
   This example does not use other workloads.

{% include "./assets/auto-app-discovery.md_2.yaml" %}

Applying these resources results in the creation
of the following `KeptnApp` resource:

{% include "./assets/auto-app-discovery.md_3.yaml" %}

With the `KeptnApp` resource created,
you get observability of your application's deployments
by using the OpenTelemetry tracing features
that are provided by Keptn:

![Application deployment trace](./assets/trace.png)
