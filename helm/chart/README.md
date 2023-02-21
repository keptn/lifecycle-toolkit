# Keptn Lifecycle Toolkit

KLT introduces a more cloud-native approach for pre- and post-deployment, as well as the concept of application health
checks

<!-- markdownlint-disable MD012 -->
## Parameters

### OpenTelemetry

| Name                | Description                                   | Value                 |
| ------------------- | --------------------------------------------- | --------------------- |
| `otelCollector.url` | Sets the URL for the open telemetry collector | `otel-collector:4317` |

### General

| Name                         | Description                                          | Value    |
| ---------------------------- | ---------------------------------------------------- | -------- |
| `deployment.imagePullPolicy` | Sets the image pull policy for kubernetes deployment | `Always` |

