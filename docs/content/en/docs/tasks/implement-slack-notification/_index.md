---
title: Implement Slack Notification
description: Learn how to implement Slack notification as a post deployment task.
icon: concepts
layout: quickstart
weight: 24
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

## Create Slack Webhook

At first, create an incoming slack webhook.
Necessary information is available in the [slack api page](https://api.slack.com/messaging/webhooks).
Once you create the webhook, you will get a URL similar to below example.

`https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX`

`T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX` is the secret part of the webhook which we would need in the next step.

## Create slack-secret

Create a `slack-secret.yaml` definition using the following command.
This will create a kubernetes secret named `slack-secret.yaml` in the `examples/sample-app/base` directory.
Before running this command change your current directory into `examples/sample-app`.

```shell
kubectl create secret generic slack-secret --from-literal=SECURE_DATA='{"slack_hook":<YOUR_HOOK_SECRET>,"text":"Deployed PodTatoHead Application"}' -n podtato-kubectl -oyaml --dry-run=client > base/slack-secret.yaml
```

## Enable post deployment task

To enable Slack notification add `post-deployment-notification` in as a postDeploymentTasks in the
[examples/sample-app/base/app.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/base/app.yaml)
file as shown below.

```yaml
  postDeploymentTasks:
    - post-deployment-notification
```
