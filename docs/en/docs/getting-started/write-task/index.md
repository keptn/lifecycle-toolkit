## Add Pre-Deployment Task

---

**TL;DR**

You can also used the prepared manifest and apply it directly using: `kubectl apply -k demo-application/with-pre-deployment/` and proceed [here](#watch-workload-behavior-with-pre-deployment-task).

---
In this step, we will use a pre-defined Keptn Function to check if a service is available before we deploy it.

Let's assume that the other services should not start before the entry service is available.

We could achieve this in two different ways:

<details>
<summary>Use a hosted Keptn Function</summary>
In this case, published this function in our repository and you can simply reference to it in your KeptnTaskDefinition Manifest as this:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha1
kind: KeptnTaskDefinition
metadata:
  name: pre-deployment-check-entry
  namespace: podtato-kubectl
spec:
  function:
    httpRef:
      url: https://raw.githubusercontent.com/keptn/lifecycle-toolbox/main/functions-runtime/samples/ts/http.ts
    parameters:
      map:
        url: http://podtato-head-entry.podtato-kubectl.svc.cluster.local:9000
```

Note, that we referred to the URL function in `.spec.function.httpRef.url`

</details>

<details>
<summary>Create the Task Inline</summary>
Alternatively, you could also create the function directly in the KeptnTaskDefinition manifest. This would look like this:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha1
kind: KeptnTaskDefinition
metadata:
  name: pre-deployment-check-entry
  namespace: podtato-kubectl
spec:
  function:
    inline:
      code: |
        let text = Deno.env.get("DATA");
        let data;
        data = JSON.parse(text);
  
        try {
          let resp = await fetch(data.url);
        }
        catch (error){
          console.error("Could not fetch url");
          Deno.exit(1);
        }
    parameters:
      map:
        url: http://podtato-head-entry.podtato-kubectl.svc.cluster.local:9000
```

In this case, we added the typescript code directly in the manifest. This is valuable if you want to deploy the function code nearby your application and don't need to share it or don't want to rely on an external service.
</details>

### Change manifests and apply them
Now that we specified our pre-deployment task, we can add this to our workloads. Therefore, you can add the `keptn.sh/pre-deployment-tasks: pre-deployment-check-entry` annotation to all workloads except the entry service in your manifest.yaml.

Before we apply the new manifests, we need to drop the old entry-service to see the effect of the pre-deployment tasks by executing `kubectl delete deploy -n keptn podtato-head-entry`.

Whatever path you have chosen before, create this manifest (`workload-pre-deploy.yaml`) nearby your application manifest, raise the version of your application in `app.yaml`, and the workload versions in the `manifest.yaml` and apply it.

### Watch workload behavior with pre-deployment task
Now you should see that the entry service is deployed first and the other services are not deployed until the entry service is available.

Using `kubectl get pods` you see that the pre-deployment tasks fail as long as the entry workload is not available. Once the entry workload is available, the pre-deployment tasks succeed and the other workloads will be scheduled.

As before, you can also watch the behavior of the application using `kubectl get keptnappversions -n podtato-kubectl` and the workload state using `kubectl get kwi -n podtato-kubectl`.

## I want to ...
* Deploy the Demo Application
    * [Using kubectl](./deploy-kubectl)
    * [Using ArgoCD](./deploy-argocd)
* [Observe my deployment using Prometheus and Grafana](./observability.md)
* [Write a Keptn Task](./writing-a-task.md)

* [Modify the Demo Application](./modify-app.md)
