### Create Secret for Slack here

```
kubectl create secret generic slack-notification --from-literal=SECURE_DATA='{"slack_hook":"<WebHook>","text":"Deployed PodTatoHead Application"}' -n podtato-kubectl -oyaml --dry-run > slack-secret.yaml
```

<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=858843d8-8da2-4ce5-a325-e5321c770a78" />