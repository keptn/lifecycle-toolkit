---
comments: true
---

# Integrate Keptn with your Applications

Keptn works
on top of the default scheduler for the cluster.
This allows it to:

- Track all activities of all workloads on the cluster,
  no matter what tool is used for the deployment
- Inject pre/post-deployment checks into all workloads.

Keptn monitors resources
that have been applied into the Kubernetes cluster
and reacts if it finds a workload with special annotations/labels.
Keptn uses metadata that is added to the Kubernetes workloads
to identify the workloads of interest.

To integrate Keptn with your workloads:

- You must first
  [install and enable](../installation/index.md#basic-installation)
  Keptn.
- Annotate or label your
  [workloads](https://kubernetes.io/docs/concepts/workloads/)
  ([Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
  [StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
  [DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
  and
  [ReplicaSets](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/)
  with either Keptn or Kubernetes keys.

    - [Basic annotations](#basic-annotations) or labels are required for all Keptn features except Keptn metrics.
    - [Pre- and post-deployment checks](#basic-annotations) are required only for the Release lifecycle management
      feature.

Keptn uses these annotations to the Kubernetes workloads to create the
[KeptnWorkload](../reference/api-reference/lifecycle/v1/index.md#keptnworkload)
and
[KeptnApp](../reference/crd-reference/app.md)
resources that it uses to provide observability
and release lifecycle management.

> Note: Annotations are not required if you are only using the
  `metrics-operator` component of Keptn
  to observe Keptn metrics.

## Basic annotations

Keptn automatically discovers `KeptnApp` resources,
based on the annotations or labels.
This enables the Keptn observability features
(based on OpenTelemetry) for existing applications,
without additional Keptn configuration.

Keptn monitors your
[Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
and
[ReplicaSets](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/),
and
[DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/)
resources in the namespaces where Keptn is enabled.
If Keptn finds any of these resources and the resource has either
the `keptn.sh` or the `kubernetes` annotations/labels,
it creates appropriate
[KeptnWorkload](../reference/api-reference/lifecycle/v1/index.md#keptnworkload)
and
[KeptnApp](../reference/crd-reference/app.md)
resources for the version it detects.

The basic keptn.sh keys that can be used for annotations or labels are:

```yaml
keptn.sh/workload: myAwesomeWorkload
keptn.sh/version: myAwesomeWorkloadVersion
keptn.sh/app: myAwesomeAppName
keptn.sh/container: myAwesomeContainer
```

Alternatively, you can use Kubernetes keys for annotations or labels.
These are part of the Kubernetes
[Recommended Labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/):

```yaml
app.kubernetes.io/name: myAwesomeWorkload
app.kubernetes.io/version: myAwesomeWorkloadVersion
app.kubernetes.io/part-of: myAwesomeAppName
```

These keys are defined as:

- `keptn.sh/workload` or `app.kubernetes.io/name`: Determines the name
  of the generated
  [KeptnWorkload](../reference/api-reference/lifecycle/v1/index.md#keptnworkload)
  resource.
- `keptn.sh/version` or `app.kubernetes.io/version`:
  Determines the version of the `KeptnWorkload`
  that represents the Workload.
  If the Workload has no `version` annotation/labels
  and the pod has only one container,
  Keptn takes the image tag as version
  (unless it is "latest").
- `keptn.sh/app` or `app.kubernetes.io/part-of`: Determines the name
  of the generated `KeptnApp` representing your Application.
  All workloads that share the same value for this label
  are consolidated into the same `KeptnApp` resource
  that you can generate following the instructions in
  [Auto app discovery](auto-app-discovery.md).
- `keptn.sh/container`: Determines the name of the container in the workload,
  from which Keptn extracts the version.
  This applies to single- and multi-container
  workloads.
  If the given container name does not match any container in the workload
  no version can be determined.
  Note that there is no equivalent `app.kubernetes.io/` annotation/label for this label.

Keptn automatically generates appropriate
[KeptnApp](../reference/crd-reference/app.md)
resources that are used for observability,
based on whether the `keptn.sh/app` or `app.kubernetes.io/part-of`
annotation/label is populated:

- If either of these annotations/labels are populated,
  Keptn automatically generates a `KeptnApp` resource
  that includes all workloads that have the same annotation/label,
  thus creating a `KeptnApp` resource for each defined grouping

- If only the `workload` and `version` annotations/labels are available
  (in other words, neither the `keptn.sh/app`
  or `app.kubernetes.io/part-of` annotation/label is populated),
  Keptn creates a `KeptnApp` resource for each `KeptnWorkload`
  and your observability output traces the individual `Keptnworkload` resources
  but not the combined workloads that constitute your deployed application.

See
[Keptn Applications and Keptn Workloads](../components/lifecycle-operator/keptn-apps.md)
for architectural information about how `KeptnApp` and `KeptnWorkloads`
are implemented.

## Annotations vs. labels

The same keys can be used as
[annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/)
or
[labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/).
Both annotations and labels are can be attached to Kubernetes objects.
Some key differences between the two:

- Annotations
    - Are not used to identify and select objects
    - Can contain up to 262144 chars
    - Metadata in an annotation can be small or large,
      structured or unstructured,
      and can include characters not permitted by labels

- Labels
    - Can be used to select objects
      and to find collections of objects that satisfy certain conditions
    - Can contain up to 63 chars
    - Are appropriate for identifying attributes of objects
      that are meaningful and relevant to users
      but do not directly imply semantics to the core system

Annotations take precedence over labels,
and the `keptn.sh` keys take precedence over `app.kubernetes.io` keys.
In other words:

- The operator first checks if the `keptn.sh` key is present
  in the annotations, and then in the labels.
- If neither is the case, it looks for the `app.kubernetes.io` equivalent,
  again first in the annotations, then in the labels.

In general, annotations are more appropriate than labels
for integrating Keptn with your applications
because they store references, names, and version information
so the 63 char limitation is quite restrictive.
However, labels can be used if you specifically need them
and can accommodate the size restriction.
