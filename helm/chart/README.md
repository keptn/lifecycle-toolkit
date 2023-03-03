# Keptn Lifecycle Toolkit

KLT introduces a more cloud-native approach for pre- and post-deployment, as well as the concept of application health
checks

<!-- markdownlint-disable MD012 -->
## Parameters

### Keptn Scheduler

| Name                                                                             | Description                                                    | Value                     |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------- | ------------------------- |
| `scheduler.scheduler.containerSecurityContext`                                   | Sets security context                                          |                           |
| `scheduler.scheduler.env.otelCollectorUrl`                                       | sets url for open telemetry collector                          | `otel-collector:4317`     |
| `scheduler.scheduler.image.repository`                                           | set image repository for scheduler                             | `ghcr.io/keptn/scheduler` |
| `scheduler.scheduler.image.tag`                                                  | set image tag for scheduler                                    | `202303031677839700`      |
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

### Keptn Cert Manager common

| Name                                                                               | Description                                    | Value               |
| ---------------------------------------------------------------------------------- | ---------------------------------------------- | ------------------- |
| `certificateOperator.replicas`                                                     | customize number of replicas                   | `1`                 |
| `certificateOperator.nodeSelector`                                                 | specify custom node selectors for cert manager | `{}`                |
| `certificateOperator.tolerations`                                                  | customize tolerations for cert manager         | `[]`                |
| `certificateOperator.topologySpreadConstraints`                                    | add topology constraints for cert manager      | `[]`                |
| `lifecycleManagerConfig.controllerManagerConfigYaml.health.healthProbeBindAddress` | TODO  TODO  TODO                               | `:8081`             |
| `lifecycleManagerConfig.controllerManagerConfigYaml.leaderElection.leaderElect`    | TODO  TODO  TODO                               | `true`              |
| `lifecycleManagerConfig.controllerManagerConfigYaml.leaderElection.resourceName`   | TODO  TODO  TODO                               | `6b866dd9.keptn.sh` |
| `lifecycleManagerConfig.controllerManagerConfigYaml.metrics.bindAddress`           | TODO  TODO  TODO                               | `127.0.0.1:8080`    |
| `lifecycleManagerConfig.controllerManagerConfigYaml.webhook.port`                  | TODO  TODO  TODO                               | `9443`              |

### Keptn Cert Manager controller

| Name                                                   | Description                                      | Value                                |
| ------------------------------------------------------ | ------------------------------------------------ | ------------------------------------ |
| `certificateOperator.manager.containerSecurityContext` | Sets security context for the cert manager       |                                      |
| `certificateOperator.manager.image.repository`         | specify repo for manager image                   | `ghcr.io/keptn/certificate-operator` |
| `certificateOperator.manager.image.tag`                | select tag for manager container                 | `202303031677839700`                 |
| `certificateOperator.manager.imagePullPolicy`          | select image pull policy for manager container   | `Always`                             |
| `certificateOperator.manager.livenessProbe`            | custom RBAC proxy liveness probe                 |                                      |
| `certificateOperator.manager.readinessProbe`           | custom manager readiness probe                   |                                      |
| `certificateOperator.manager.resources`                | custom limits and requests for manager container |                                      |

### Keptn Lifecycle Operator common

| Name                                          | Description                                               | Value       |
| --------------------------------------------- | --------------------------------------------------------- | ----------- |
| `lifecycleOperator.replicas`                  | customize number of installed lifecycle operator replicas | `1`         |
| `lifecycleWebhookService.ports[0].port`       | TODO  TODO  TODO                                          | `443`       |
| `lifecycleWebhookService.ports[0].protocol`   | TODO  TODO  TODO                                          | `TCP`       |
| `lifecycleWebhookService.ports[0].targetPort` | TODO  TODO  TODO                                          | `9443`      |
| `lifecycleWebhookService.type`                | TODO  TODO  TODO                                          | `ClusterIP` |
| `lifecycleOperator.nodeSelector`              | add custom nodes selector to lifecycle operator           | `{}`        |
| `lifecycleOperator.tolerations`               | add custom tolerations to lifecycle operator              | `[]`        |
| `lifecycleOperator.topologySpreadConstraints` | add custom topology constraints to lifecycle operator     | `[]`        |

### Keptn Lifecycle Operator controller

| Name                                                                          | Description                                             | Value                                          |
| ----------------------------------------------------------------------------- | ------------------------------------------------------- | ---------------------------------------------- |
| `lifecycleOperator.manager.containerSecurityContext`                          | Sets security context privileges                        |                                                |
| `lifecycleOperator.manager.containerSecurityContext.allowPrivilegeEscalation` |                                                         | `false`                                        |
| `lifecycleOperator.manager.containerSecurityContext.capabilities.drop`        |                                                         | `["ALL"]`                                      |
| `lifecycleOperator.manager.containerSecurityContext.privileged`               |                                                         | `false`                                        |
| `lifecycleOperator.manager.containerSecurityContext.runAsGroup`               |                                                         | `65532`                                        |
| `lifecycleOperator.manager.containerSecurityContext.runAsNonRoot`             |                                                         | `true`                                         |
| `lifecycleOperator.manager.containerSecurityContext.runAsUser`                |                                                         | `65532`                                        |
| `lifecycleOperator.manager.containerSecurityContext.seccompProfile.type`      |                                                         | `RuntimeDefault`                               |
| `lifecycleOperator.manager.env.keptnAppControllerLogLevel`                    | sets the log level of Keptn App Controller              | `0`                                            |
| `lifecycleOperator.manager.env.keptnAppVersionControllerLogLevel`             | sets the log level of Keptn AppVersion Controller       | `0`                                            |
| `lifecycleOperator.manager.env.keptnEvaluationControllerLogLevel`             | sets the log level of Keptn Evaluation Controller       | `0`                                            |
| `lifecycleOperator.manager.env.keptnTaskControllerLogLevel`                   | sets the log level of Keptn Task Controller             | `0`                                            |
| `lifecycleOperator.manager.env.keptnTaskDefinitionControllerLogLevel`         | sets the log level of Keptn TaskDefinition Controller   | `0`                                            |
| `lifecycleOperator.manager.env.keptnWorkloadControllerLogLevel`               | sets the log level of Keptn Workload Controller         | `0`                                            |
| `lifecycleOperator.manager.env.keptnWorkloadInstanceControllerLogLevel`       | sets the log level of Keptn WorkloadInstance Controller | `0`                                            |
| `lifecycleOperator.manager.env.optionsControllerLogLevel`                     | sets the log level of Keptn Options Controller          | `0`                                            |
| `lifecycleOperator.manager.env.otelCollectorUrl`                              | Sets the URL for the open telemetry collector           | `0`                                            |
| `lifecycleOperator.manager.env.functionRunnerImage`                           | specify image for task runtime                          | `ghcr.keptn.sh/keptn/functions-runtime:v0.6.0` |
| `lifecycleOperator.manager.image.repository`                                  | specify registry for manager image                      | `ghcr.io/keptn/lifecycle-operator`             |
| `lifecycleOperator.manager.image.tag`                                         | select tag for manager image                            | `202303031677839700`                           |
| `lifecycleOperator.manager.imagePullPolicy`                                   | specify pull policy for manager image                   | `Always`                                       |
| `lifecycleOperator.manager.livenessProbe`                                     | custom livenessprobe for manager container              |                                                |
| `lifecycleOperator.manager.readinessProbe`                                    | custom readinessprobe for manager container             |                                                |
| `lifecycleOperator.manager.resources`                                         | specify limits and requests for manager container       |                                                |

### Keptn Metrics Operator common

| Name                                                                             | Description                                             | Value               |
| -------------------------------------------------------------------------------- | ------------------------------------------------------- | ------------------- |
| `metricsOperator.replicas`                                                       | customize number of installed metrics operator replicas | `1`                 |
| `metricsOperatorService.ports[0].name`                                           |                                                         | `https`             |
| `metricsOperatorService.ports[0].port`                                           |                                                         | `8443`              |
| `metricsOperatorService.ports[0].protocol`                                       |                                                         | `TCP`               |
| `metricsOperatorService.ports[0].targetPort`                                     |                                                         | `https`             |
| `metricsOperatorService.ports[1].name`                                           |                                                         | `custom-metrics`    |
| `metricsOperatorService.ports[1].port`                                           |                                                         | `443`               |
| `metricsOperatorService.ports[1].targetPort`                                     |                                                         | `custom-metrics`    |
| `metricsOperatorService.ports[2].name`                                           |                                                         | `metrics`           |
| `metricsOperatorService.ports[2].port`                                           |                                                         | `2222`              |
| `metricsOperatorService.ports[2].protocol`                                       |                                                         | `TCP`               |
| `metricsOperatorService.ports[2].targetPort`                                     |                                                         | `metrics`           |
| `metricsOperatorService.type`                                                    |                                                         | `ClusterIP`         |
| `metricsManagerConfig.controllerManagerConfigYaml.health.healthProbeBindAddress` | TODO  TODO  TODO                                        | `:8081`             |
| `metricsManagerConfig.controllerManagerConfigYaml.leaderElection.leaderElect`    | TODO  TODO  TODO                                        | `true`              |
| `metricsManagerConfig.controllerManagerConfigYaml.leaderElection.resourceName`   | TODO  TODO  TODO                                        | `3f8532ca.keptn.sh` |
| `metricsManagerConfig.controllerManagerConfigYaml.metrics.bindAddress`           | TODO  TODO  TODO                                        | `127.0.0.1:8080`    |
| `metricsManagerConfig.controllerManagerConfigYaml.webhook.port`                  | TODO  TODO  TODO                                        | `9443`              |
| `metricsWebhookService.ports[0].port`                                            | TODO  TODO  TODO                                        | `443`               |
| `metricsWebhookService.ports[0].protocol`                                        | TODO  TODO  TODO                                        | `TCP`               |
| `metricsWebhookService.ports[0].targetPort`                                      | TODO  TODO  TODO                                        | `9443`              |
| `metricsWebhookService.type`                                                     | TODO  TODO  TODO                                        | `ClusterIP`         |
| `metricsOperator.nodeSelector`                                                   | add custom nodes selector to metrics operator           | `{}`                |
| `metricsOperator.tolerations`                                                    | add custom tolerations to metrics operator              | `[]`                |
| `metricsOperator.topologySpreadConstraints`                                      | add custom topology constraints to metrics operator     | `[]`                |

### Keptn Metrics Operator controller

| Name                                                                        | Description                                       | Value                            |
| --------------------------------------------------------------------------- | ------------------------------------------------- | -------------------------------- |
| `metricsOperator.manager.containerSecurityContext`                          | Sets security context privileges                  |                                  |
| `metricsOperator.manager.containerSecurityContext.allowPrivilegeEscalation` |                                                   | `false`                          |
| `metricsOperator.manager.containerSecurityContext.capabilities.drop`        |                                                   | `["ALL"]`                        |
| `metricsOperator.manager.image.repository`                                  | specify registry for manager image                | `ghcr.io/keptn/metrics-operator` |
| `metricsOperator.manager.image.tag`                                         | select tag for manager image                      | `202303031677839700`             |
| `metricsOperator.manager.env.exposeKeptnMetrics`                            | enable metrics exporter                           | `true`                           |
| `metricsOperator.manager.env.metricsControllerLogLevel`                     | sets the log level of Metrics Controller          | `0`                              |
| `metricsOperator.manager.livenessProbe`                                     | custom livenessprobe for manager container        |                                  |
| `metricsOperator.manager.readinessProbe`                                    | custom readinessprobe for manager container       |                                  |
| `metricsOperator.manager.resources`                                         | specify limits and requests for manager container |                                  |

### Global

| Name                      | Description                            | Value           |
| ------------------------- | -------------------------------------- | --------------- |
| `kubernetesClusterDomain` | overrides domain.local                 | `cluster.local` |
| `imagePullSecrets`        | global value for image registry secret | `[]`            |
