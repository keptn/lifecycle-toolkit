# Keptn

Keptn provides a cloud-native approach for pre- and post-deployment,
and application health checks

<!-- markdownlint-disable MD012 -->
## Parameters

### Keptn Scheduler

| Name                                                         | Description                                                    | Value                     |
| ------------------------------------------------------------ | -------------------------------------------------------------- | ------------------------- |
| `scheduler.containerSecurityContext`                         | Sets security context                                          |                           |
| `scheduler.env.otelCollectorUrl`                             | sets url for open telemetry collector                          | `otel-collector:4317`     |
| `scheduler.image.repository`                                 | set image repository for scheduler                             | `ghcr.io/keptn/scheduler` |
| `scheduler.image.tag`                                        | set image tag for scheduler                                    | `v0.8.2`                  |
| `scheduler.imagePullPolicy`                                  | set image pull policy for scheduler                            | `Always`                  |
| `scheduler.livenessProbe`                                    | customizable liveness probe for the scheduler                  |                           |
| `scheduler.readinessProbe`                                   | customizable readiness probe for the scheduler                 |                           |
| `scheduler.resources`                                        | sets cpu and memory resurces/limits for scheduler              |                           |
| `schedulerConfig.leaderElection.leaderElect`                 | enables leader election for multiple replicas of the scheduler | `false`                   |
| `schedulerConfig.profiles[0].plugins.permit.enabled[0].name` | enables permit plugin                                          | `KLCPermit`               |
| `schedulerConfig.profiles[0].schedulerName`                  | changes scheduler name                                         | `keptn-scheduler`         |
| `scheduler.nodeSelector`                                     | adds node selectors for scheduler                              | `{}`                      |
| `scheduler.replicas`                                         | modifies replicas                                              | `1`                       |
| `scheduler.tolerations`                                      | adds tolerations for scheduler                                 | `[]`                      |
| `scheduler.topologySpreadConstraints`                        | add topology constraints for scheduler                         | `[]`                      |

### Keptn Lifecycle Operator common

| Name                                                   | Description                                                                                                                                                   | Value               |
| ------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- |
| `lifecycleManagerConfig.health.healthProbeBindAddress` | setup on what address to start the default health handler                                                                                                     | `:8081`             |
| `lifecycleManagerConfig.leaderElection.leaderElect`    | enable leader election for multiple replicas of the lifecycle operator                                                                                        | `true`              |
| `lifecycleManagerConfig.leaderElection.resourceName`   | define LeaderElectionID                                                                                                                                       | `6b866dd9.keptn.sh` |
| `lifecycleManagerConfig.metrics.bindAddress`           | MetricsBindAddress is the TCP address that the controller should bind to for serving prometheus metrics. It can be set to "0" to disable the metrics serving. | `127.0.0.1:8080`    |
| `lifecycleManagerConfig.webhook.port`                  | setup port for the lifecycle operator admission webhook                                                                                                       | `9443`              |
| `lifecycleOperator.replicas`                           | customize number of installed lifecycle operator replicas                                                                                                     | `1`                 |
| `lifecycleOperatorMetricsService`                      | Adjust settings here to change the k8s service for scraping Prometheus metrics                                                                                |                     |
| `lifecycleWebhookService`                              | Mutating Webhook Configurations for lifecycle Operator                                                                                                        |                     |
| `lifecycleWebhookService.ports[0].port`                |                                                                                                                                                               | `443`               |
| `lifecycleWebhookService.ports[0].protocol`            |                                                                                                                                                               | `TCP`               |
| `lifecycleWebhookService.ports[0].targetPort`          |                                                                                                                                                               | `9443`              |
| `lifecycleWebhookService.type`                         |                                                                                                                                                               | `ClusterIP`         |
| `lifecycleOperator.nodeSelector`                       | add custom nodes selector to lifecycle operator                                                                                                               | `{}`                |
| `lifecycleOperator.tolerations`                        | add custom tolerations to lifecycle operator                                                                                                                  | `[]`                |
| `lifecycleOperator.topologySpreadConstraints`          | add custom topology constraints to lifecycle operator                                                                                                         | `[]`                |

### Keptn Lifecycle Operator controller

| Name                                                                  | Description                                                 | Value                                 |
| --------------------------------------------------------------------- | ----------------------------------------------------------- | ------------------------------------- |
| `lifecycleOperator.containerSecurityContext`                          | Sets security context privileges                            |                                       |
| `lifecycleOperator.containerSecurityContext.allowPrivilegeEscalation` |                                                             | `false`                               |
| `lifecycleOperator.containerSecurityContext.capabilities.drop`        |                                                             | `["ALL"]`                             |
| `lifecycleOperator.containerSecurityContext.privileged`               |                                                             | `false`                               |
| `lifecycleOperator.containerSecurityContext.runAsGroup`               |                                                             | `65532`                               |
| `lifecycleOperator.containerSecurityContext.runAsNonRoot`             |                                                             | `true`                                |
| `lifecycleOperator.containerSecurityContext.runAsUser`                |                                                             | `65532`                               |
| `lifecycleOperator.containerSecurityContext.seccompProfile.type`      |                                                             | `RuntimeDefault`                      |
| `lifecycleOperator.env.keptnAppControllerLogLevel`                    | sets the log level of Keptn App Controller                  | `0`                                   |
| `lifecycleOperator.env.keptnAppCreationRequestControllerLogLevel`     | sets the log level of Keptn App Creation Request Controller | `0`                                   |
| `lifecycleOperator.env.keptnAppVersionControllerLogLevel`             | sets the log level of Keptn AppVersion Controller           | `0`                                   |
| `lifecycleOperator.env.keptnEvaluationControllerLogLevel`             | sets the log level of Keptn Evaluation Controller           | `0`                                   |
| `lifecycleOperator.env.keptnTaskControllerLogLevel`                   | sets the log level of Keptn Task Controller                 | `0`                                   |
| `lifecycleOperator.env.keptnTaskDefinitionControllerLogLevel`         | sets the log level of Keptn TaskDefinition Controller       | `0`                                   |
| `lifecycleOperator.env.keptnWorkloadControllerLogLevel`               | sets the log level of Keptn Workload Controller             | `0`                                   |
| `lifecycleOperator.env.keptnWorkloadInstanceControllerLogLevel`       | sets the log level of Keptn WorkloadInstance Controller     | `0`                                   |
| `lifecycleOperator.env.optionsControllerLogLevel`                     | sets the log level of Keptn Options Controller              | `0`                                   |
| `lifecycleOperator.env.otelCollectorUrl`                              | Sets the URL for the open telemetry collector               | `otel-collector:4317`                 |
| `lifecycleOperator.env.functionRunnerImage`                           | specify image for deno task runtime                         | `ghcr.io/keptn/deno-runtime:v1.0.1`   |
| `lifecycleOperator.env.pythonRunnerImage`                             | specify image for python task runtime                       | `ghcr.io/keptn/python-runtime:v1.0.0` |
| `lifecycleOperator.image.repository`                                  | specify registry for manager image                          | `ghcr.io/keptn/lifecycle-operator`    |
| `lifecycleOperator.image.tag`                                         | select tag for manager image                                | `v0.8.2`                              |
| `lifecycleOperator.imagePullPolicy`                                   | specify pull policy for manager image                       | `Always`                              |
| `lifecycleOperator.livenessProbe`                                     | custom livenessprobe for manager container                  |                                       |
| `lifecycleOperator.readinessProbe`                                    | custom readinessprobe for manager container                 |                                       |
| `lifecycleOperator.resources`                                         | specify limits and requests for manager container           |                                       |

### Global

| Name                      | Description                                                                                                                                     | Value           |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | --------------- |
| `kubernetesClusterDomain` | overrides domain.local                                                                                                                          | `cluster.local` |
| `imagePullSecrets`        | global value for image registry secret                                                                                                          | `[]`            |
| `schedulingGatesEnabled`  | enables the scheduling gates in lifecycle-operator. This feature is available in alpha version from K8s 1.27 or 1.26 enabling the alpha version | `false`         |
