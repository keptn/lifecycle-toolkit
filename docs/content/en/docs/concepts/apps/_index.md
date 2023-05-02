---
title: Apps
description: Learn what Keptn Apps are and how to use them
icon: concepts
layout: quickstart
weight: 10
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

An App contains information about all workloads and checks associated with an application.
It will use the following structure for the specification of the pre/post deployment and pre/post evaluations checks
that should be executed at app level:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: podtato-kubectl
spec:
  version: "1.3"
  workloads:
  - name: podtato-head-left-arm
    version: 0.1.0
  - name: podtato-head-left-leg
    version: 1.2.3
  postDeploymentTasks:
  - post-deployment-hello
  preDeploymentEvaluations:    
  - my-prometheus-definition
```

While changes in the workload version will affect only workload checks, a change in the app version will also cause a
new execution of app level checks.

## Automatic App Discovery

The Keptn Lifecycle Toolkit also provides the option to automatically discover `KeptnApp`s, based on the
recommended Kubernetes labels `app.kubernetes.io/part-of`, `app.kubernetes.io/name` `app.kubernetes.io/version`.
This allows users to enable the observability features provided by the Lifecycle Toolkit for
their existing applications, without the need for creating any Keptn-related custom resources.

To enable the automatic discovery of `KeptnApp`s for your existing applications, the following steps will
be required:

1. Make sure the namespace of your application is enabled to be managed by the Keptn Lifecycle Toolkit,
by adding the annotation `keptn.sh/lifecycle-toolkit: "enabled"` to your namespace.
2. Make sure the following labels and/or annotations are present in the pod template
specs of your Workloads (i.e. `Deployments`/`StatefulSets`/`DaemonSets`/`ReplicaSets`) within your application:
    - `app.kubernetes.io/name`: Determines the name of the generated `KeptnWorkload` representing the
    Workload.
    - `app.kubernetes.io/version`: Determines the version of the `KeptnWorkload` representing the Workload.
    - `app.kubernetes.io/part-of`: Determines the name of the generated `KeptnApp` representing your
    Application.
    All Workloads that share the same value for this label will be consolidated into the same `KeptnApp`.

As an example, consider the following application, consisting of several deployments, which is going to be
deployed into a KLT-enabled namespace:

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

Applying these resources will then result in the creation of the following `KeptnApp`:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: podtato-kubectl
  annotations:
    app.kubernetes.io/managed-by: "klt"
spec:
  version: "<version string based on a hash of all containing workloads>"
  workloads:
  - name: podtato-head-frontend
    version: 0.1.0
  - name: podtato-head-hat
    version: 1.1.1
```

Due to the creation of this resource, you will now get observability of your application's deployments due to
the OpenTelemetry tracing features provided by the Keptn Lifecycle Toolkit:

![Application deployment trace](assets/trace.png)
