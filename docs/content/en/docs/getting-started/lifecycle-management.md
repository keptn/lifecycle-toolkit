---
title: Release Lifecycle Management
description: Add KeptnTasks to deployments
weight: 50
---

The Release Lifecycle Management tools run
pre- and post-deployment tasks and checks
for your existing cloud-native deployments
to make them more robust.
This tutorial introduces these tools.

> This tutorial assumes you have already completed the
[Getting started with Keptn Observability](../getting-started/)
exercise.
> Please ensure you've finished that before attempting this guide.

## Keptn Pre and Post Deployment Tasks

When Keptn is successfully monitoring your deployments, it can also run arbitrary tasks and SLO evaluations:

- pre-deployment (before the pod is scheduled) and
- post-deployment (after the post is scheduled)

> Pre and post deployments can also run on a KeptnApp level.
> See [annotations to KeptnApp](../implementing/integrate.md#annotations-to-keptnapp)

## Prerequisites: Deploy webhook sink

During this exercise, you will configure Keptn to trigger a webhook before and after a deployment has successfully completed.

For demo purposes, a place is required to send those request.
Install the [open source webhook.site tool](https://github.com/webhooksite/webhook.site/tree/master/kubernetes) now.

This will provide a place, on cluster, to send (and view) web requests.

> Note: If you have your own endpoint, you can skip this step.

```shell
kubectl apply -f https://raw.githubusercontent.com/webhooksite/webhook.site/master/kubernetes/namespace.yml
kubectl apply -f https://raw.githubusercontent.com/webhooksite/webhook.site/master/kubernetes/redis.deployment.yml
kubectl apply -f https://raw.githubusercontent.com/webhooksite/webhook.site/master/kubernetes/laravel-echo-server.deployment.yml
kubectl apply -f https://raw.githubusercontent.com/webhooksite/webhook.site/master/kubernetes/webhook.deployment.yml
kubectl apply -f https://raw.githubusercontent.com/webhooksite/webhook.site/master/kubernetes/service.yml
```

Wait until all pods are running in the `webhook` namespace then port-forward and view the webhook sink page:

```shell
kubectl -n webhook wait --for=condition=Ready pods --all
kubectl -n webhook port-forward svc/webhook 8084
```

Open a browser and go to `http://localhost:8084`

You should see a page like this with a unique URL (your ID will be different).

![webhook.site page](../assets/webhook.site.1.png)

Make a note of that unique URL.

### Verify Webhook Sink

Open a new browser table and go to your unique URL.
The page should remain blank, but when toggling back to `http://localhost:8084`, you should see a new entry.

Every request sent to that unique URL will be logged here.

![webhook.site entry](../assets/webhook.site.2.png)

## Add a Post Deployment Task

Add a task which will trigger after a deployment.

Change `UUID` to whatever value you have.
Apply this manifest:

```yaml
---
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition
metadata:
  name: send-event
  namespace: keptndemo
spec:
  retries: 0
  timeout: 5s
  container:
    name: curlcontainer
    image: curlimages/curl:latest
    args: [
        '-X',
        'POST',
        'http://webhook.webhook.svc.cluster.local:8084/YOUR-UUID-HERE',
        '-H',
        'Content-Type: application/json',
        '-d',
        '{ "from": "keptn send-event" }'
    ] 
```

### Verify it works

Verify that the KeptnTaskDefinition above actually works.

Trigger an on-demand task execution to verify that the job and pod function correctly.

In the following steps we will have Keptn orchestrate this for us automatically.

Apply this manifest:

```yaml
---
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTask
metadata:
  name: runsendevent1
  namespace: keptndemo
spec:
  taskDefinition: send-event
  context:
    appName: "my-test-app"
    appVersion: "1.0.0"
    objectType: ""
    taskType: ""
    workloadName: "my-test-workload"
    workloadVersion: "1.0.0"
```

If it works, `kubectl -n keptndemo get jobs` should show:

```shell
NAME                  COMPLETIONS   DURATION   AGE
runsendevent1-*****   1/1           6s         2m
```

`kubectl -n keptndemo get pods` will show the successfully executed pod.

The webhook sync should show this:

![webhook sync](../assets/webhook.site.3.png)

Incidentally, this is exactly how you can use Keptn with [applications deployed outside of Kubernetes](../implementing/tasks-non-k8s-apps.md).

> Note: If you want to trigger this multiple times, you must change the KeptnTask name.
>
> For example, by changing `runsendevent1` to `runsendevent2`

## Ask Keptn to trigger task after Deployment

Annotate the demo application `Deployment` manifest to have Keptn automatically trigger the task after every deployment.

Recall the `Deployment` from the [Observability](../getting-started/observability.md#step-3-deploy-demo-application)
Getting started guide.

Add a new label so the `labels` section looks like this:

```yaml
...
labels:
    app.kubernetes.io/part-of: keptndemoapp
    app.kubernetes.io/name: nginx
    app.kubernetes.io/version: 0.0.2
    keptn.sh/post-deployment-tasks: "send-event"
...
```

Increase the version number to `0.0.2` and re-apply the manifest.

Here is a full version of the new YAML:

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: keptndemo
  labels:
    app.kubernetes.io/name: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: nginx
  template:
    metadata:
      labels:
        app.kubernetes.io/part-of: keptndemoapp
        app.kubernetes.io/name: nginx
        app.kubernetes.io/version: 0.0.2
        keptn.sh/post-deployment-tasks: "send-event"
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

> Best Practice: Start with post deployment tasks.
> Pre-deployment tasks can potentially block deployments (see below).

### What Happens Next?

1. The deployment will be applied
1. When the pods are running, Keptn will automatically create a `KeptnTask` resource for version `0.0.2` of this KeptnApp
1. The `KeptnTask` will create a Kubernetes Job
1. The Kubernetes Job will create a Pod
1. The pod will run curl and send a new event to the event sink

### Pre-deployment Tasks

Keptn Tasks can also be executed pre-deployment (before the pods are scheduled).
Do this by using the `keptn.sh/pre-deployment-tasks` label.

> Note: If a pre-deployment task fails, the pod will remain in a Pending state.

## Further Information

There is a lot more you can do with KeptnTasks.
See [pre and post deployment checks page](../implementing/integrate.md#pre--and-post-deployment-checks) to find out more.

## What's next?

Keptn can also run pre and post deployment SLO evaluations.

Continue the Keptn learning journey by adding evaluations.
