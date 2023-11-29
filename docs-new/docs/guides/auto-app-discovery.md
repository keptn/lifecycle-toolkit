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
metadata:
  name: simpleapp
  namespace: simplenode-dev
```

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

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: podtato-kubectl
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled"

