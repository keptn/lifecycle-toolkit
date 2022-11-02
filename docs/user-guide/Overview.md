# How to use the Lifecycle Controller?

## Annotating workloads

The Keptn Lifecycle Controller monitors manifests that have been applied against the Kubernetes API and reacts if it finds a workload with special annotations/labels.
For this, you should annotate your [Workload](https://kubernetes.io/docs/concepts/workloads/) with (at least) the following two annotations:

```yaml
keptn.sh/app: myAwesomeAppName
keptn.sh/workload: myAwesomeWorkload
keptn.sh/version: myAwesomeWorkloadVersion
```

Alternatively, you can use Kubernetes [Recommended Labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/) to annotate your workload:

```yaml
app.kubernetes.io/part-of: myAwesomeAppName
app.kubernetes.io/name: myAwesomeWorkload
app.kubernetes.io/version: myAwesomeWorkloadVersion
```

In general, the Keptn Annotations/Labels take precedence over the Kubernetes recommended labels. If there is no version annotation/label and there is only one container in the pod, the Lifecycle Controller will take the image tag as version (if it is not "latest").

In case you want to run pre- and post-deployment checks, further annotations are necessary:

```yaml
keptn.sh/pre-deployment-tasks: verify-infrastructure-problems
keptn.sh/post-deployment-tasks: slack-notification,performance-test
```

The value of these annotations are Keptn [CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
called [KeptnTaskDefinition](#keptn-task-definition)s. These CRDs contains re-usable "functions" that can
executed before and after the deployment. In this example, before the deployment starts, a check for open problems in your infrastructure
is performed. If everything is fine, the deployment continues and afterward, a slack notification is sent with the result of
the deployment and a pipeline to run performance tests is invoked. Otherwise, the deployment is kept in a pending state until
the infrastructure is capable to accept deployments again.

## Examples

A more comprehensive example can be found in our [examples folder](/examples/podtatohead-deployment/) where we
use [Podtato-Head](https://github.com/podtato-head/podtato-head) to run some simple pre-deployment checks.

To run the example, use the following commands:

```bash
cd ./examples/podtatohead-deployment/
kubectl apply -f .
```

Afterward, you can monitor the status of the deployment using

```bash
kubectl get keptnworkloadinstance -n podtato-kubectl -w
```

The deployment for a Workload will stay in a `Pending` state until the respective pre-deployment check is completed. Afterward, the deployment will start and when it is `Succeeded`, the post-deployment checks will start.