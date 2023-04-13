# Sample Application

This example should demonstrate the capabilities of the lifecycle toolkit as illustrated here:

![img.png](assets/big-picture.png)

## PostDeployment Slack Notification

This section describes how to **prepare and enable** post-deployment tasks to send notifications to slack using webhooks.

### Create Slack Webhook

In the first step, create an incoming slack webhook.
Necessary information is available in the [slack api page](https://api.slack.com/messaging/webhooks).
Once you create the webhook, you will get a URL similar to below example.

`https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX`

`T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX` is the secret part of the webhook which we would need in the next step.

### Create slack-secret

Create a `slack-secret.yaml` definition using the following command.
This will create a kubernetes secret named `slack-secret.yaml` in the [base](./base) directory.

```shell
kubectl create secret generic slack-secret --from-literal=SECURE_DATA='{"slack_hook":<YOUR_HOOK_SECRET>,"text":"Deployed PodTatoHead Application"}' -n podtato-kubectl -oyaml --dry-run=client > base/slack-secret.yaml
```

### Enable post deployment task

To enable Slack notification add `post-deployment-notification` in as a postDeploymentTasks in the
[app.yaml](base/app.yaml) file as shown below.

```yaml
  postDeploymentTasks:
    - post-deployment-notification
```

## Deploy the Observability Part and Keptn-lifecycle-toolkit

```make install```

## Port-Forward Grafana

```make port-forward-grafana```

If you want to port-forward to a different port, please execute:
```make port-forward-grafana GRAFANA_PORT_FORWARD=<port>```

## Deploy Version 1 of the PodTatoHead

```make deploy-version-1```

Now watch the progress on the cluster
```kubectl get keptnworkloadinstances```
```kubectl get keptnappversions```

You could also open up a browser and watch the progress in Jaeger.
You can find the Context ID in the "TraceId" Field of
the KeptnAppVersion

The deployment should fail because of too few cpu resources

## Deploy Version 2 of the PodTatoHead

```make deploy-version-2```

* Watch the progress of the deployments
* After some time, you should see that everything is successful

## Deploy Version 3

```make deploy-version-3```

* This should only change one service, you can see that only this changed in the trace

<!-- markdownlint-disable-next-line MD033 MD013 -->
<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=858843d8-8da2-4ce5-a325-e5321c770a78" />
