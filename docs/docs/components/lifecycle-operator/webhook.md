---
comments: true
---

# Keptn Webhooks

Keptn uses
[Admission Webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#what-are-admission-webhooks)
to mutate resources.
To enable the webhook (and therefore Keptn Lifecycle Management)
for a certain namespace, the namespace must be annotated:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: podtato-kubectl
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled"  # this line tells the webhook to handle the namespace
```

The mutating webhook only modifies specifically annotated resources in the annotated namespace.
When the webhook receives a request for a new pod,
it looks for the workload annotations:

```yaml
keptn.sh/workload: "some-workload-name"
```

The mutation consists in changing the scheduler used for the deployment
with the Keptn Scheduler, or adding the
[Scheduling Gate](https://keptn.sh/stable/docs/components/scheduling/#keptn-scheduling-gates-for-k8s-127-and-above).
The webhook then creates a workload and app resource for each annotated resource.
You can also specify a custom app definition with the annotation:

```yaml
keptn.sh/app: "your-app-name"
```

In this case the webhook does not generate an app,
but it expects that the user will provide one.
Additionally, it computes a version string,
using a hash function that takes certain properties of the pod as parameters
(e.g. the images of its containers).
Next, it looks for an existing instance of a `Workload CRD`
for the specified workload name:

- If it finds the `Workload`,
  it updates its version according to the previously computed version string.
  In addition, it includes a reference to the ReplicaSet UID of the pod
  (i.e. the Pods owner),
  or the pod itself, if it does not have an owner.
- If it does not find a workload instance,
  it creates one containing the previously computed version string.
  In addition, it includes a reference to the ReplicaSet UID of the pod
  (i.e. the Pods owner), or the pod itself, if it does not have an owner.

It uses the following annotations for the specification
of the pre/post deployment checks that should be executed for the `Workload`:

- `keptn.sh/pre-deployment-tasks: task1,task2`
- `keptn.sh/post-deployment-tasks: task1,task2`

and for the Evaluations:

- `keptn.sh/pre-deployment-evaluations: my-evaluation-definition`
- `keptn.sh/post-deployment-evaluations: my-eval-definition`

After either one of those actions has been taken,
the pod is allowed to be scheduled.
