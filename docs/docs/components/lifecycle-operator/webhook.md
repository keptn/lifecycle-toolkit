---
comments: true
---

# Mutating Webhook

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
When the webhook receives a request for a new pod, it first either replaces the scheduler
with the Keptn Scheduler, or adds the
[Scheduling Gate](https://keptn.sh/stable/docs/components/scheduling/#keptn-scheduling-gates-for-k8s-127-and-above).

In the next step it looks for the workload annotations:

```yaml
keptn.sh/workload: "some-workload-name"
keptn.sh/version: "some-workload-version"
```

If the `keptn.sh/version` annotation is missing, the webhook computes a version string,
using a hash function that takes certain properties of the pod as parameters
(e.g. the images of its containers).
Next, it looks for an existing instance of a `KeptnWorkload`
for the specified workload name:

- If it finds the `KeptnWorkload`,
  it updates its version according to the previously computed version string.
  In addition, it includes a reference to the ReplicaSet UID of the pod
  (i.e. the Pods owner).
- If it does not find a workload instance,
  it creates one containing the previously computed version string.

Afterwards the webhook looks for the application annotation:

```yaml
keptn.sh/app: "your-app-name"
```

The webhook searches for the `KeptnAppCreationRequest` resource with the name stored in the `keptn.sh/app`
annotations.
If it doesn't find it, it creates it and the automatic creation of the `KeptnApp` is afterwards
handled by the `KeptnAppCreationRequest Controller`.

The `keptn.sh/app` annotation is not mandatory for single-service applications.
If you have a multi-service application, you must add it to all workloads
to define which workloads belong to the application.

The Pod can also contain information about the definition of pre or
post-deployment tasks or evaluations for each workload.
These are specified via these annotations:

- `keptn.sh/pre-deployment-tasks: task1,task2`
- `keptn.sh/post-deployment-tasks: task1,task2`

and for the Evaluations:

- `keptn.sh/pre-deployment-evaluations: my-evaluation-definition`
- `keptn.sh/post-deployment-evaluations: my-eval-definition`

The lists of tasks or evaluations are parsed and stored in the `KeptnWorkload`
resource created in the previous steps.
