# KubeCon 2022 NA Demo

This demonstration is based on the [Observability Example](../observability) and should do the following:

![img.png](assets/big-picture.png)

## Prepare Secret for Slack Notification
As a first step, create an incoming webhook according to the instructions
https://api.slack.com/messaging/webhooks

Afterwards create a secret with the created webhook
> kubectl create secret generic slack-notification --from-literal=SECURE_DATA='{"slack_hook":<YOUR_HOOK>,"text":"Deployed PodTatoHead Application"}' -n podtato-kubectl -oyaml --dry-run=client > base/slack-secret.yaml

## Deploy the Observability Part and Keptn-lifecycle-toolkit
> make install

## Port-Forward Grafana
> make port-forward-grafana

If you want to port-forward to a different port, please execute:
> make port-forward-grafana GRAFANA_PORT_FORWARD=<port>

## Deploy Version 1 of the PodTatoHead
> make deploy-version-1

Now watch the progress on the cluster
> kubectl get keptnworkloadinstances
> kubectl get keptnappversions

You could also open up a browser and watch the progress in Jaeger. You can find the Context ID in the "TraceId" Field of the KeptnAppVersion

The deployment should fail because of too few cpu resources

## Deploy Version 2 of the PodTatoHead
> make deploy-version-2

* Watch the progress of the deployments
* After some time, you should see that everything is successful

## Deploy Version 3
> make deploy-version-3

* This should only change one service, you can see that only this changed in the trace
