# Keptn

Keptn provides a cloud-native approach for pre- and post-deployment,
and application health checks

<!-- markdownlint-disable MD012 -->
## Parameters

### Keptn Lifecycle Operator common


### Global parameters

| Name                                                    | Description                                                                                                                                                   | Value               |
| ------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- |
| `global.certManagerEnabled`                             | Enable this value to install Keptn Certificate Manager                                                                                                        | `true`              |
| `global.imageRegistry`                                  | Global container image registry                                                                                                                               | `ghcr.io`           |
| `global.imagePullSecrets`                               | Global Docker registry secret names as an array                                                                                                               | `[]`                |
| `global.imagePullPolicy`                                | select global image pull policy                                                                                                                               | `""`                |
| `global.commonLabels`                                   | Common labels to add to all Keptn resources. Evaluated as a template                                                                                          | `{}`                |
| `global.commonAnnotations`                              | Common annotations to add to all Keptn resources. Evaluated as a template                                                                                     | `{}`                |
| `global.caInjectionAnnotations`                         | CA injection annotations for cert-manager.io configuration                                                                                                    | `{}`                |
| `lifecycleOperatorConfig.health.healthProbeBindAddress` | setup on what address to start the default health handler                                                                                                     | `:8081`             |
| `lifecycleOperatorConfig.leaderElection.leaderElect`    | enable leader election for multiple replicas of the lifecycle operator                                                                                        | `true`              |
| `lifecycleOperatorConfig.leaderElection.resourceName`   | define LeaderElectionID                                                                                                                                       | `6b866dd9.keptn.sh` |
| `lifecycleOperatorConfig.metrics.bindAddress`           | MetricsBindAddress is the TCP address that the controller should bind to for serving prometheus metrics. It can be set to "0" to disable the metrics serving. | `127.0.0.1:8080`    |
| `lifecycleOperatorConfig.webhook.port`                  | setup port for the lifecycle operator admission webhook                                                                                                       | `9443`              |
| `lifecycleWebhookService`                               | Mutating Webhook Configurations for lifecycle Operator                                                                                                        |                     |
| `lifecycleWebhookService.ports[0].port`                 |                                                                                                                                                               | `443`               |
| `lifecycleWebhookService.ports[0].protocol`             |                                                                                                                                                               | `TCP`               |
| `lifecycleWebhookService.ports[0].targetPort`           |                                                                                                                                                               | `9443`              |
| `lifecycleWebhookService.type`                          |                                                                                                                                                               | `ClusterIP`         |

### Keptn Lifecycle Operator controller

| Name                                                                  | Description                                                                    | Value                                 |
| --------------------------------------------------------------------- | ------------------------------------------------------------------------------ | ------------------------------------- |
| `lifecycleOperator.containerSecurityContext`                          | Sets security context privileges                                               |                                       |
| `lifecycleOperator.containerSecurityContext.allowPrivilegeEscalation` |                                                                                | `false`                               |
| `lifecycleOperator.containerSecurityContext.capabilities.drop`        |                                                                                | `["ALL"]`                             |
| `lifecycleOperator.containerSecurityContext.privileged`               |                                                                                | `false`                               |
| `lifecycleOperator.containerSecurityContext.runAsGroup`               |                                                                                | `65532`                               |
| `lifecycleOperator.containerSecurityContext.runAsNonRoot`             |                                                                                | `true`                                |
| `lifecycleOperator.containerSecurityContext.runAsUser`                |                                                                                | `65532`                               |
| `lifecycleOperator.containerSecurityContext.seccompProfile.type`      |                                                                                | `RuntimeDefault`                      |
| `lifecycleOperator.env.functionRunnerImage`                           | specify image for deno task runtime                                            | `ghcr.io/keptn/deno-runtime:v2.0.3`   |
| `lifecycleOperator.env.keptnAppControllerLogLevel`                    | sets the log level of Keptn App Controller                                     | `0`                                   |
| `lifecycleOperator.env.keptnAppCreationRequestControllerLogLevel`     | sets the log level of Keptn App Creation Request Controller                    | `0`                                   |
| `lifecycleOperator.env.keptnAppVersionControllerLogLevel`             | sets the log level of Keptn AppVersion Controller                              | `0`                                   |
| `lifecycleOperator.env.keptnEvaluationControllerLogLevel`             | sets the log level of Keptn Evaluation Controller                              | `0`                                   |
| `lifecycleOperator.env.keptnTaskControllerLogLevel`                   | sets the log level of Keptn Task Controller                                    | `0`                                   |
| `lifecycleOperator.env.keptnTaskDefinitionControllerLogLevel`         | sets the log level of Keptn TaskDefinition Controller                          | `0`                                   |
| `lifecycleOperator.env.keptnWorkloadControllerLogLevel`               | sets the log level of Keptn Workload Controller                                | `0`                                   |
| `lifecycleOperator.env.keptnWorkloadVersionControllerLogLevel`        | sets the log level of Keptn WorkloadVersion Controller                         | `0`                                   |
| `lifecycleOperator.env.keptnDoraMetricsPort`                          | sets the port for accessing lifecycle metrics in prometheus format             | `2222`                                |
| `lifecycleOperator.env.optionsControllerLogLevel`                     | sets the log level of Keptn Options Controller                                 | `0`                                   |
| `lifecycleOperator.env.pythonRunnerImage`                             | specify image for python task runtime                                          | `ghcr.io/keptn/python-runtime:v1.0.4` |
| `lifecycleOperator.image.registry`                                    | specify the container registry for the lifecycle-operator image                | `""`                                  |
| `lifecycleOperator.image.repository`                                  | specify registry for manager image                                             | `keptn/lifecycle-operator`            |
| `lifecycleOperator.image.tag`                                         | select tag for manager image                                                   | `v0.9.2`                              |
| `lifecycleOperator.image.imagePullPolicy`                             | specify pull policy for the manager image. This overrides global values        | `""`                                  |
| `lifecycleOperator.livenessProbe`                                     | custom liveness probe for manager container                                    |                                       |
| `lifecycleOperator.readinessProbe`                                    | custom readinessprobe for manager container                                    |                                       |
| `lifecycleOperator.resources`                                         | specify limits and requests for manager container                              |                                       |
| `lifecycleOperator.nodeSelector`                                      | add custom nodes selector to lifecycle operator                                | `{}`                                  |
| `lifecycleOperator.replicas`                                          | customize number of installed lifecycle operator replicas                      | `1`                                   |
| `lifecycleOperator.tolerations`                                       | add custom tolerations to lifecycle operator                                   | `[]`                                  |
| `lifecycleOperator.topologySpreadConstraints`                         | add custom topology constraints to lifecycle operator                          | `[]`                                  |
| `lifecycleOperatorMetricsService`                                     | Adjust settings here to change the k8s service for scraping Prometheus metrics |                                       |

### Global

| Name                      | Description                                                                                                                                     | Value                                                          |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------- |
| `kubernetesClusterDomain` | overrides cluster.local                                                                                                                         | `cluster.local`                                                |
| `annotations`             | add deployment level annotations                                                                                                                | `{}`                                                           |
| `podAnnotations`          | adds pod level annotations                                                                                                                      | `{}`                                                           |
| `schedulingGatesEnabled`  | enables the scheduling gates in lifecycle-operator. This feature is available in alpha version from K8s 1.27 or 1.26 enabling the alpha version | `false`                                                        |
| `promotionTasksEnabled`   | enables the promotion task feature in the lifecycle-operator.                                                                                   | `false`                                                        |
| `allowedNamespaces`       | specifies the allowed namespaces for the lifecycle orchestration functionality                                                                  | `[]`                                                           |
| `deniedNamespaces`        | specifies a list of namespaces where the lifecycle orchestration functionality is disabled, ignored if `allowedNamespaces` is set               | `["cert-manager","keptn-system","observability","monitoring"]` |

### Keptn Scheduler

| Name                                                         | Description                                                             | Value                 |
| ------------------------------------------------------------ | ----------------------------------------------------------------------- | --------------------- |
| `scheduler.nodeSelector`                                     | adds node selectors for scheduler                                       | `{}`                  |
| `scheduler.replicas`                                         | modifies replicas                                                       | `1`                   |
| `scheduler.containerSecurityContext`                         | Sets security context                                                   |                       |
| `scheduler.env.otelCollectorUrl`                             | sets url for open telemetry collector                                   | `otel-collector:4317` |
| `scheduler.image.registry`                                   | specify the container registry for the scheduler image                  | `""`                  |
| `scheduler.image.repository`                                 | set image repository for scheduler                                      | `keptn/scheduler`     |
| `scheduler.image.tag`                                        | set image tag for scheduler                                             | `v0.9.2`              |
| `scheduler.image.imagePullPolicy`                            | specify pull policy for the manager image. This overrides global values | `""`                  |
| `scheduler.livenessProbe`                                    | customizable liveness probe for the scheduler                           |                       |
| `scheduler.readinessProbe`                                   | customizable readiness probe for the scheduler                          |                       |
| `scheduler.resources`                                        | sets cpu and memory resources/limits for scheduler                      |                       |
| `scheduler.topologySpreadConstraints`                        | add topology constraints for scheduler                                  | `[]`                  |
| `schedulerConfig.profiles[0].schedulerName`                  | changes scheduler name                                                  | `keptn-scheduler`     |
| `schedulerConfig.leaderElection.leaderElect`                 | enables leader election for multiple replicas of the scheduler          | `false`               |
| `schedulerConfig.profiles[0].plugins.permit.enabled[0].name` | enables permit plugin                                                   | `KLCPermit`           |
| `scheduler.tolerations`                                      | adds tolerations for scheduler                                          | `[]`                  |
