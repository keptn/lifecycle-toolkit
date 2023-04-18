# Add Application Awareness

In the previous step, we installed the demo application without any application awareness.
This means that the Lifecycle
Toolkit assumed that every workload is a single-service application at the moment and created the Application resources
for you.

To get the overall state of an application, we need a grouping of workloads, called KeptnApp in the Lifecycle Toolkit.
To get this working, we need to modify our application manifest with two things:

* Add an "app.kubernetes.io/part-of" or "keptn.sh/app" label to the deployment
* Create an application resource

## Preparing the Manifest and create an App resource

---

### TL;DR

You can also used the prepared manifest and apply it directly using: `kubectl apply -k sample-app/version-2/` and
proceed [here](#watch-application-behavior).

---

### Otherwise

Create a temporary directory and copy the base manifest there:

```shell
mkdir ./my-deployment
cp demo-application/base/manifest.yml ./my-deployment
```

Now, open the manifest in your favorite editor and add the following label to the deployments, e.g.:

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podtato-head-right-leg
  namespace: podtato-kubectl
  labels:
    app: podtato-head
spec:
  selector:
    matchLabels:
      component: podtato-head-right-leg
  template:
    metadata:
      labels:
        component: podtato-head-right-leg
      annotations:
        keptn.sh/workload: "right-leg"
        keptn.sh/version: "0.1.0"
        keptn.sh/app: "podtato-head"
    spec:
      terminationGracePeriodSeconds: 5
      containers:
        - name: server
          image: ghcr.io/podtato-head/right-leg:0.2.7
          imagePullPolicy: Always
          ports:
            - containerPort: 9000
          env:
            - name: PODTATO_PORT
              value: "9000"
```

Now, update the version of the workloads in the manifest to `0.2.0`.

Finally, create an application resource (app.yaml) and save it in the directory as well:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: podtato-kubectl
spec:
  version: "0.1.0"
  workloads:
    - name: left-arm
      version: "0.1.1"
    - name: left-leg
      version: "0.1.1"
    - name: entry
      version: "0.1.1"
    - name: right-arm
      version: "0.1.1"
    - name: left-arm
      version: "0.1.1"
    - name: hat
      version: "0.1.1"
```

Now, apply the manifests:

```shell
kubectl apply -f ./my-deployment/.
```

## Watch Application behavior

Now, your application gets deployed in an application aware way.
This means that pre-deployment tasks and evaluations
would be executed if you would have any.
The same would happen for post-deployment tasks and evaluations after the last
workload has been deployed successfully.

Now that you defined your application, you could watch the state of the whole application using:

```shell
kubectl get keptnappversions -n podtato-kubectl`
```

You should see that the application is in a progressing state as long as the workloads (`kubectl get kwi`) are
progressing.
After the last application has been deployed, and post-deployment tasks and evaluations are finished (there
are none at this point), the state should switch to completed.

Now, we have deployed an application and are able to get the total state of the application state.
Metrics and traces
get exported and now we're ready to dive deeper in the world of Pre- and Post-Deployment Tasks.
