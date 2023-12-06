<!-- markdownlint-disable MD041 -->
```yaml
apiVersion: metrics.keptn.sh/v1beta1
kind: KeptnMetricsProvider
metadata:
  name: datadog-provider
  namespace: podtato-kubectl
spec:
  type: datadog
  targetServer: "<datadog-url>"
  secretKeyRef:
    name: datadog-secret
---
apiVersion: v1
kind: Secret
metadata:
  name: datadog-secret
data:
  DD_CLIENT_API_KEY: api-key
  DD_CLIENT_API_KEY: app-key
type: Opaque
```
<!-- markdownlint-enable MD041 -->
