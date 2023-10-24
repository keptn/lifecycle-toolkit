# Keptn Helm Chart

Keptn provides a “cloud-native” approach for managing the application release lifecycle
metrics, observability, health checks, with pre- and post-deployment evaluations and tasks.

<!-- markdownlint-disable MD012 -->

## Parameters

### Keptn

| Name                        | Description                                            | Value  |
| --------------------------- | ------------------------------------------------------ | ------ |
| `certManager.enabled`       | Enable this value to install Keptn Certificate Manager | `true` |
| `lifecycleOperator.enabled` | Enable this value to install Keptn Lifecycle Operator  | `true` |
| `metricsOperator.enabled`   | Enable this value to install Keptn Metrics Operator    | `true` |

### Global parameters

| Name                       | Description                                                               | Value |
| -------------------------- | ------------------------------------------------------------------------- | ----- |
| `global.imageRegistry`     | Global Docker image registry                                              | `""`  |
| `global.imagePullSecrets`  | Global Docker registry secret names as an array                           | `[]`  |
| `global.commonLabels`      | Common annotations to add to all Keptn resources. Evaluated as a template | `{}`  |
| `global.commonAnnotations` | Common annotations to add to all Mongo resources. Evaluated as a template | `{}`  |
