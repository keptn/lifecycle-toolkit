---
title: Auto app discovery
description: Use Keptn automatic app discovery
weight: 10
---

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

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podtato-head-frontend
  namespace: podtato-kubectl
spec:
  template:
    metadata:
      labels:
        app.kubernetes.io/name: podtato-head-frontend
        app.kubernetes.io/part-of: podtato-head
        app.kubernetes.io/version: 0.1.0
    spec:
      containers:
        - name: podtato-head-frontend
          image: podtato-head-frontend
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podtato-head-hat
  namespace: podtato-kubectl
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app.kubernetes.io/name: podtato-head-hat
        app.kubernetes.io/part-of: podtato-head
        app.kubernetes.io/version: 0.1.1
    spec:
      containers:
        - name: podtato-head-hat
          image: podtato-head-hat
```

Applying these resources results in the creation
of the following `KeptnApp` resource:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: podtato-kubectl
  annotations:
    app.kubernetes.io/managed-by: "keptn"
spec:
  version: "<version string based on a hash of all containing workloads>"
  workloads:
  - name: podtato-head-frontend
    version: 0.1.0
  - name: podtato-head-hat
    version: 1.1.1
```

With the `KeptnApp` resource created,
you get observability of your application's deployments
by using the OpenTelemetry tracing features
that are provided by Keptn:

![Application deployment trace](../assets/trace.png)