### Create Secret for Slack here

```
kubectl create secret generic slack-notification --from-literal=SECURE_DATA='{"slack_hook":"<WebHook>","text":"Deployed PodTatoHead Application"}' -n podtato-kubectl -oyaml --dry-run > slack-secret.yaml
```
