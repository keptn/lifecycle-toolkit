# Keptn Lifecycle Toolkit

KLT introduces a more cloud-native approach for pre- and post-deployment, as well as the concept of application health
checks

<!-- markdownlint-disable MD012 -->
## Parameters

### Controller Log level

| Name                                                | Description                                              | Value |
| --------------------------------------------------- | -------------------------------------------------------- | ----- |
| `operator.keptnappController.logLevel`              | Sets the log level of Keptn App Controller               | `0`   |
| `operator.keptnappversionController.logLevel`       | Sets the log level of Keptn Version Controller           | `0`   |
| `operator.keptnevaluationController.logLevel`       | Sets the log level of Keptn Evaluation Controller        | `0`   |
| `operator.keptntaskController.logLevel`             | Sets the log level of Keptn Task Controller              | `0`   |
| `operator.keptntaskdefinitionController.logLevel`   | Sets the log level of Keptn Task Defintion Controller    | `0`   |
| `operator.keptnworkloadController.logLevel`         | Sets the log level of Keptn Workload Controller          | `0`   |
| `operator.keptnworkloadinstanceController.logLevel` | Sets the log level of Keptn Workload Instance Controller | `0`   |
| `operator.metricsController.logLevel`               | Sets the log level of Keptn Metrics Controller           | `0`   |
| `operator.optionsController.logLevel`               | Sets the log level of Keptn Options Controller           | `0`   |

### OpenTelemetry

| Name                | Description                                   | Value                 |
| ------------------- | --------------------------------------------- | --------------------- |
| `otelCollector.url` | Sets the URL for the open telemetry collector | `otel-collector:4317` |

### General

| Name                         | Description                                          | Value    |
| ---------------------------- | ---------------------------------------------------- | -------- |
| `deployment.imagePullPolicy` | Sets the image pull policy for kubernetes deployment | `Always` |


## Other info

<!-- markdownlint-enable MD012 -->
