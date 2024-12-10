---
comments: true
---

# Keptn integration with Scheduling

Keptn integrates with Kubernetes scheduling to block
the deployment of applications that do not satisfy Keptn defined pre-deployment checks.
The default scheduling paradigm is implemented using
[Kubernetes scheduling gates](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-scheduling-readiness/)
to gate Pods until the required deployment checks pass.

## Keptn Scheduling Gates

When a workload is applied to a Kubernetes cluster,
the Mutating Webhook checks each Pod for annotations.
If
[Keptn specific annotations](../guides/integrate.md#basic-annotations)
are present,
the Webhook adds a scheduling gate to the Pod called `keptn-prechecks-gate`.
This spec tells the Kubernetes scheduling framework
to wait for the Keptn checks before binding the pod to a node.

For example, a pod gated by Keptn looks like the following:

```yaml
{% include "./assets/gated.yaml" %}
```

If the `pre-deployment` checks and evaluations have finished successfully,
the WorkloadVersion Controller removes the gate from the Pod.
The default k8s scheduler then binds the Pod to a node.
If the `pre-deployment` checks have not yet finished successfully,
the gate stays and the Pod remains in the pending state.
When removing the gate,
the WorkloadVersion controller also adds the following annotation so that,
if the Pod is updated, it is not gated again:

```yaml
{% include "./assets/gate-removed.yaml" %}
```

## Integrating Keptn with your custom scheduler

Keptn scheduling logics are compatible with
the [Kubernetes Scheduler Framework](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/).
Keptn does not work with a custom scheduler unless it is implemented as
a [scheduler plugin](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/#plugin-configuration).
