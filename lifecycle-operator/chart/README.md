# Keptn

Keptn introduces a more cloud-native approach for pre- and post-deployment,
as well as the concept of application health checks

<!-- markdownlint-disable MD012 -->
## Parameters

### Keptn Scheduler

| Name                                                                             | Description                                                    | Value                     |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------- | ------------------------- |
| `scheduler.scheduler.containerSecurityContext`                                   | Sets security context                                          |                           |
| `scheduler.scheduler.env.otelCollectorUrl`                                       | sets url for open telemetry collector                          | `otel-collector:4317`     |
| `scheduler.scheduler.image.repository`                                           | set image repository for scheduler                             | `ghcr.io/keptn/scheduler` |
| `scheduler.scheduler.image.tag`                                                  | set image tag for scheduler                                    | `v0.8.2`                  |
| `scheduler.scheduler.imagePullPolicy`                                            | set image pull policy for scheduler                            | `Always`                  |
| `scheduler.scheduler.livenessProbe`                                              | customizable liveness probe for the scheduler                  |                           |
| `scheduler.scheduler.readinessProbe`                                             | customizable readiness probe for the scheduler                 |                           |
| `scheduler.scheduler.resources`                                                  | sets cpu and memory resurces/limits for scheduler              |                           |
| `schedulerConfig.schedulerConfigYaml.leaderElection.leaderElect`                 | enables leader election for multiple replicas of the scheduler | `false`                   |
| `schedulerConfig.schedulerConfigYaml.profiles[0].plugins.permit.enabled[0].name` | enables permit plugin                                          | `KLCPermit`               |
| `schedulerConfig.schedulerConfigYaml.profiles[0].schedulerName`                  | changes scheduler name                                         | `keptn-scheduler`         |
| `scheduler.nodeSelector`                                                         | adds node selectors for scheduler                              | `{}`                      |
| `scheduler.replicas`                                                             | modifies replicas                                              | `1`                       |
| `scheduler.tolerations`                                                          | adds tolerations for scheduler                                 | `[]`                      |
| `scheduler.topologySpreadConstraints`                                            | add topology constraints for scheduler                         | `[]`                      |

### Keptn Lifecycle Operator common

| Name                                                                               | Description                                                                                                                                                   | Value               |
| ---------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- |
| `lifecycleManagerConfig.controllerManagerConfigYaml.health.healthProbeBindAddress` | setup on what address to start the default health handler                                                                                                     | `:8081`             |
| `lifecycleManagerConfig.controllerManagerConfigYaml.leaderElection.leaderElect`    | enable leader election for multiple replicas of the lifecycle operator                                                                                        | `true`              |
| `lifecycleManagerConfig.controllerManagerConfigYaml.leaderElection.resourceName`   | define LeaderElectionID                                                                                                                                       | `6b866dd9.keptn.sh` |
| `lifecycleManagerConfig.controllerManagerConfigYaml.metrics.bindAddress`           | MetricsBindAddress is the TCP address that the controller should bind to for serving prometheus metrics. It can be set to "0" to disable the metrics serving. | `127.0.0.1:8080`    |
| `lifecycleManagerConfig.controllerManagerConfigYaml.webhook.port`                  | setup port for the lifecycle operator admission webhook                                                                                                       | `9443`              |
| `lifecycleOperator.replicas`                                                       | customize number of installed lifecycle operator replicas                                                                                                     | `1`                 |
| `lifecycleOperatorMetricsService`                                                  | Adjust settings here to change the k8s service for scraping Prometheus metrics                                                                                |                     |
| `lifecycleWebhookService`                                                          | Mutating Webhook Configurations for lifecycle Operator                                                                                                        |                     |
| `lifecycleWebhookService.ports[0].port`                                            |                                                                                                                                                               | `443`               |
| `lifecycleWebhookService.ports[0].protocol`                                        |                                                                                                                                                               | `TCP`               |
| `lifecycleWebhookService.ports[0].targetPort`                                      |                                                                                                                                                               | `9443`              |
| `lifecycleWebhookService.type`                                                     |                                                                                                                                                               | `ClusterIP`         |
| `lifecycleOperator.nodeSelector`                                                   | add custom nodes selector to lifecycle operator                                                                                                               | `{}`                |
| `lifecycleOperator.tolerations`                                                    | add custom tolerations to lifecycle operator                                                                                                                  | `[]`                |
| `lifecycleOperator.topologySpreadConstraints`                                      | add custom topology constraints to lifecycle operator                                                                                                         | `[]`                |

### Keptn Lifecycle Operator controller

| Name                                                                          | Description                                                 | Value                                 |
| ----------------------------------------------------------------------------- | ----------------------------------------------------------- | ------------------------------------- |
| `lifecycleOperator.manager.containerSecurityContext`                          | Sets security context privileges                            |                                       |
| `lifecycleOperator.manager.containerSecurityContext.allowPrivilegeEscalation` |                                                             | `false`                               |
| `lifecycleOperator.manager.containerSecurityContext.capabilities.drop`        |                                                             | `["ALL"]`                             |
| `lifecycleOperator.manager.containerSecurityContext.privileged`               |                                                             | `false`                               |
| `lifecycleOperator.manager.containerSecurityContext.runAsGroup`               |                                                             | `65532`                               |
| `lifecycleOperator.manager.containerSecurityContext.runAsNonRoot`             |                                                             | `true`                                |
| `lifecycleOperator.manager.containerSecurityContext.runAsUser`                |                                                             | `65532`                               |
| `lifecycleOperator.manager.containerSecurityContext.seccompProfile.type`      |                                                             | `RuntimeDefault`                      |
| `lifecycleOperator.manager.env.keptnAppControllerLogLevel`                    | sets the log level of Keptn App Controller                  | `0`                                   |
| `lifecycleOperator.manager.env.keptnAppCreationRequestControllerLogLevel`     | sets the log level of Keptn App Creation Request Controller | `0`                                   |
| `lifecycleOperator.manager.env.keptnAppVersionControllerLogLevel`             | sets the log level of Keptn AppVersion Controller           | `0`                                   |
| `lifecycleOperator.manager.env.keptnEvaluationControllerLogLevel`             | sets the log level of Keptn Evaluation Controller           | `0`                                   |
| `lifecycleOperator.manager.env.keptnTaskControllerLogLevel`                   | sets the log level of Keptn Task Controller                 | `0`                                   |
| `lifecycleOperator.manager.env.keptnTaskDefinitionControllerLogLevel`         | sets the log level of Keptn TaskDefinition Controller       | `0`                                   |
| `lifecycleOperator.manager.env.keptnWorkloadControllerLogLevel`               | sets the log level of Keptn Workload Controller             | `0`                                   |
| `lifecycleOperator.manager.env.keptnWorkloadInstanceControllerLogLevel`       | sets the log level of Keptn WorkloadInstance Controller     | `0`                                   |
| `lifecycleOperator.manager.env.optionsControllerLogLevel`                     | sets the log level of Keptn Options Controller              | `0`                                   |
| `lifecycleOperator.manager.env.otelCollectorUrl`                              | Sets the URL for the open telemetry collector               | `otel-collector:4317`                 |
| `lifecycleOperator.manager.env.functionRunnerImage`                           | specify image for deno task runtime                         | `ghcr.io/keptn/deno-runtime:v1.0.1`   |
| `lifecycleOperator.manager.env.pythonRunnerImage`                             | specify image for python task runtime                       | `ghcr.io/keptn/python-runtime:v1.0.0` |
| `lifecycleOperator.manager.image.repository`                                  | specify registry for manager image                          | `ghcr.io/keptn/lifecycle-operator`    |
| `lifecycleOperator.manager.image.tag`                                         | select tag for manager image                                | `v0.8.2`                              |
| `lifecycleOperator.manager.imagePullPolicy`                                   | specify pull policy for manager image                       | `Always`                              |
| `lifecycleOperator.manager.livenessProbe`                                     | custom livenessprobe for manager container                  |                                       |
| `lifecycleOperator.manager.readinessProbe`                                    | custom readinessprobe for manager container                 |                                       |
| `lifecycleOperator.manager.resources`                                         | specify limits and requests for manager container           |                                       |

### Global

| Name                      | Description                                                                                                                                     | Value           |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | --------------- |
| `kubernetesClusterDomain` | overrides domain.local                                                                                                                          | `cluster.local` |
| `imagePullSecrets`        | global value for image registry secret                                                                                                          | `[]`            |
| `schedulingGatesEnabled`  | enables the scheduling gates in lifecycle-operator. This feature is available in alpha version from K8s 1.27 or 1.26 enabling the alpha version | `false`         |
