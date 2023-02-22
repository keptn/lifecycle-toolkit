# Keptn Lifecycle Toolkit

KLT introduces a more cloud-native approach for pre- and post-deployment, as well as the concept of application health
checks

<!-- markdownlint-disable MD012 -->
## Parameters

### Keptn Scheduler

| Name                                                                             | Description                                                    | Value                     |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------- | ------------------------- |
| `keptnScheduler.keptnScheduler.containerSecurityContext`                         | Sets security context                                          |                           |
| `keptnScheduler.keptnScheduler.env.otelCollectorUrl`                             | sets url for open telemetry collector                          | `otel-collector:4317`     |
| `keptnScheduler.keptnScheduler.image.repository`                                 | set image repository for scheduler                             | `ghcr.io/keptn/scheduler` |
| `keptnScheduler.keptnScheduler.image.tag`                                        | set image tag for scheduler                                    | `202302231677157645`      |
| `keptnScheduler.keptnScheduler.imagePullPolicy`                                  | set image pull policy for scheduler                            | `Always`                  |
| `keptnScheduler.keptnScheduler.livenessProbe`                                    | customizable liveness probe for the scheduler                  |                           |
| `keptnScheduler.keptnScheduler.readinessProbe`                                   | customizable readiness probe for the scheduler                 |                           |
| `keptnScheduler.keptnScheduler.resources`                                        | sets cpu and memory resurces/limits for scheduler              |                           |
| `schedulerConfig.schedulerConfigYaml.leaderElection.leaderElect`                 | enables leader election for multiple replicas of the scheduler | `false`                   |
| `schedulerConfig.schedulerConfigYaml.profiles[0].plugins.permit.enabled[0].name` | enables permit plugin                                          | `KLCPermit`               |
| `schedulerConfig.schedulerConfigYaml.profiles[0].schedulerName`                  | changes scheduler name                                         | `keptn-scheduler`         |
| `keptnScheduler.nodeSelector`                                                    | adds node selectors for scheduler                              | `{}`                      |
| `keptnScheduler.replicas`                                                        | modifies replicas                                              | `1`                       |
| `keptnScheduler.tolerations`                                                     | adds tolerations for scheduler                                 | `[]`                      |
| `keptnScheduler.topologySpreadConstraints`                                       | add topology constraints for scheduler                         | `[]`                      |

### Keptn Cert Manager common

| Name                                                                         | Description                                    | Value               |
| ---------------------------------------------------------------------------- | ---------------------------------------------- | ------------------- |
| `kltCertManager.replicas`                                                    | customize number of replicas                   | `1`                 |
| `kltCertManagerMetricsService.ports[0].name`                                 | TODO  TODO  TODO                               | `https`             |
| `kltCertManagerMetricsService.ports[0].port`                                 | TODO  TODO  TODO                               | `8443`              |
| `kltCertManagerMetricsService.ports[0].protocol`                             | TODO  TODO  TODO                               | `TCP`               |
| `kltCertManagerMetricsService.ports[0].targetPort`                           | TODO  TODO  TODO                               | `https`             |
| `kltCertManagerMetricsService.type`                                          | TODO  TODO TODO                                | `ClusterIP`         |
| `kltCertManager.nodeSelector`                                                | specify custom node selectors for cert manager | `{}`                |
| `kltCertManager.tolerations`                                                 | customize tolerations for cert manager         | `[]`                |
| `kltCertManager.topologySpreadConstraints`                                   | add topology constraints for cert manager      | `[]`                |
| `klcManagerConfig.controllerManagerConfigYaml.health.healthProbeBindAddress` | TODO  TODO  TODO                               | `:8081`             |
| `klcManagerConfig.controllerManagerConfigYaml.leaderElection.leaderElect`    | TODO  TODO  TODO                               | `true`              |
| `klcManagerConfig.controllerManagerConfigYaml.leaderElection.resourceName`   | TODO  TODO  TODO                               | `6b866dd9.keptn.sh` |
| `klcManagerConfig.controllerManagerConfigYaml.metrics.bindAddress`           | TODO  TODO  TODO                               | `127.0.0.1:8080`    |
| `klcManagerConfig.controllerManagerConfigYaml.webhook.port`                  | TODO  TODO  TODO                               | `9443`              |
| `klcWebhookService.ports[0].port`                                            | TODO  TODO  TODO                               | `443`               |
| `klcWebhookService.ports[0].protocol`                                        | TODO  TODO  TODO                               | `TCP`               |
| `klcWebhookService.ports[0].targetPort`                                      | TODO  TODO  TODO                               | `9443`              |
| `klcWebhookService.type`                                                     | TODO  TODO  TODO                               | `ClusterIP`         |

### Keptn Cert Manager RBAC proxy

| Name                                                    | Description                                                | Value                                |
| ------------------------------------------------------- | ---------------------------------------------------------- | ------------------------------------ |
| `kltCertManager.kubeRbacProxy.containerSecurityContext` | Sets security context for RBAC proxy                       |                                      |
| `kltCertManager.kubeRbacProxy.image.repository`         | setup proxy image repository                               | `gcr.io/kubebuilder/kube-rbac-proxy` |
| `kltCertManager.kubeRbacProxy.image.tag`                | specify proxy image tag                                    | `v0.13.0`                            |
| `kltCertManager.kubeRbacProxy.imagePullPolicy`          | specify proxy image pull policy                            | `IfNotPresent`                       |
| `kltCertManager.kubeRbacProxy.livenessProbe`            | custom RBAC proxy liveness probe                           |                                      |
| `kltCertManager.kubeRbacProxy.readinessProbe`           | custom RBAC proxy readiness probe                          |                                      |
| `kltCertManager.kubeRbacProxy.resources`                | custom RBAC proxy's limits and requests for cpu and memory |                                      |

### Keptn Cert Manager controller

| Name                                              | Description                                      | Value                            |
| ------------------------------------------------- | ------------------------------------------------ | -------------------------------- |
| `kltCertManager.manager.containerSecurityContext` | Sets security context for the cert manager       |                                  |
| `kltCertManager.manager.image.repository`         | specify repo for manager image                   | `ghcr.io/keptn/klt-cert-manager` |
| `kltCertManager.manager.image.tag`                | select tag for manager container                 | `202302231677157645`             |
| `kltCertManager.manager.imagePullPolicy`          | select image pull policy for manager container   | `Always`                         |
| `kltCertManager.manager.livenessProbe`            | custom RBAC proxy liveness probe                 |                                  |
| `kltCertManager.manager.readinessProbe`           | custom manager readiness probe                   |                                  |
| `kltCertManager.manager.resources`                | custom limits and requests for manager container |                                  |

### Keptn Lifecycle Operator common

| Name                                                     | Description                                               | Value       |
| -------------------------------------------------------- | --------------------------------------------------------- | ----------- |
| `klcControllerManager.replicas`                          | customize number of installed lifecycle operator replicas | `1`         |
| `klcControllerManagerMetricsService.ports[0].name`       | TODO  TODO  TODO                                          | `https`     |
| `klcControllerManagerMetricsService.ports[0].port`       | TODO  TODO  TODO                                          | `8443`      |
| `klcControllerManagerMetricsService.ports[0].protocol`   | TODO  TODO  TODO                                          | `TCP`       |
| `klcControllerManagerMetricsService.ports[0].targetPort` | TODO  TODO  TODO                                          | `https`     |
| `klcControllerManagerMetricsService.ports[1].name`       | TODO  TODO  TODO                                          | `metrics`   |
| `klcControllerManagerMetricsService.ports[1].port`       | TODO  TODO  TODO                                          | `2222`      |
| `klcControllerManagerMetricsService.ports[1].protocol`   | TODO  TODO  TODO                                          | `TCP`       |
| `klcControllerManagerMetricsService.ports[1].targetPort` | TODO  TODO  TODO                                          | `metrics`   |
| `klcControllerManagerMetricsService.type`                | TODO TODO TODO                                            | `ClusterIP` |
| `klcControllerManager.nodeSelector`                      | add custom nodes selector to lifecycle operator           | `{}`        |
| `klcControllerManager.tolerations`                       | add custom tolerations to lifecycle operator              | `[]`        |
| `klcControllerManager.topologySpreadConstraints`         | add custom topology constraints to lifecycle operator     | `[]`        |

### Keptn Lifecycle Operator RBAC proxy

| Name                                                          | Description                                          | Value                                |
| ------------------------------------------------------------- | ---------------------------------------------------- | ------------------------------------ |
| `klcControllerManager.kubeRbacProxy.containerSecurityContext` | Sets security context privileges for RBAC proxy      |                                      |
| `klcControllerManager.kubeRbacProxy.image.repository`         | specify registry for RBAC proxy image                | `gcr.io/kubebuilder/kube-rbac-proxy` |
| `klcControllerManager.kubeRbacProxy.image.tag`                | select tag for RBAC proxy image                      | `v0.13.0`                            |
| `klcControllerManager.kubeRbacProxy.imagePullPolicy`          | custom pull policy for RBAC proxy image              | `IfNotPresent`                       |
| `klcControllerManager.kubeRbacProxy.livenessProbe`            | custom livenessprobe for proxy RBAC container        |                                      |
| `klcControllerManager.kubeRbacProxy.readinessProbe`           | custom readinessprobe for proxy RBAC container       |                                      |
| `klcControllerManager.kubeRbacProxy.resources`                | specify limits and requests for proxy RBAC container |                                      |

### Keptn Lifecycle Operator controller

| Name                                                                             | Description                                       | Value                                          |
| -------------------------------------------------------------------------------- | ------------------------------------------------- | ---------------------------------------------- |
| `klcControllerManager.manager.containerSecurityContext`                          | Sets security context privileges                  |                                                |
| `klcControllerManager.manager.containerSecurityContext.allowPrivilegeEscalation` |                                                   | `false`                                        |
| `klcControllerManager.manager.containerSecurityContext.capabilities.drop`        |                                                   | `["ALL"]`                                      |
| `klcControllerManager.manager.containerSecurityContext.privileged`               |                                                   | `false`                                        |
| `klcControllerManager.manager.containerSecurityContext.runAsGroup`               |                                                   | `65532`                                        |
| `klcControllerManager.manager.containerSecurityContext.runAsNonRoot`             |                                                   | `true`                                         |
| `klcControllerManager.manager.containerSecurityContext.runAsUser`                |                                                   | `65532`                                        |
| `klcControllerManager.manager.containerSecurityContext.seccompProfile.type`      |                                                   | `RuntimeDefault`                               |
| `klcControllerManager.manager.env.otelCollectorUrl`                              | Sets the URL for the open telemetry collector     | `otel-collector:4317`                          |
| `klcControllerManager.manager.env.exposeKeptnMetrics`                            | enable metrics exporter                           | `true`                                         |
| `klcControllerManager.manager.env.functionRunnerImage`                           | specify image for task runtime                    | `ghcr.keptn.sh/keptn/functions-runtime:v0.6.0` |
| `klcControllerManager.manager.image.repository`                                  | specify registry for manager image                | `ghcr.io/keptn/keptn-lifecycle-operator`       |
| `klcControllerManager.manager.image.tag`                                         | select tag for manager image                      | `202302231677157645`                           |
| `klcControllerManager.manager.imagePullPolicy`                                   | specify pull policy for manager image             | `Always`                                       |
| `klcControllerManager.manager.livenessProbe`                                     | custom livenessprobe for manager container        |                                                |
| `klcControllerManager.manager.readinessProbe`                                    | custom readinessprobe for manager container       |                                                |
| `klcControllerManager.manager.resources`                                         | specify limits and requests for manager container |                                                |

### Keptn Metrics Operator common

| Name                                                                                     | Description                                             | Value               |
| ---------------------------------------------------------------------------------------- | ------------------------------------------------------- | ------------------- |
| `metricsOperatorController.replicas`                                                     | customize number of installed metrics operator replicas | `1`                 |
| `metricsOperatorControllerMetricsService.ports[0].name`                                  |                                                         | `https`             |
| `metricsOperatorControllerMetricsService.ports[0].port`                                  |                                                         | `8443`              |
| `metricsOperatorControllerMetricsService.ports[0].protocol`                              |                                                         | `TCP`               |
| `metricsOperatorControllerMetricsService.ports[0].targetPort`                            |                                                         | `https`             |
| `metricsOperatorControllerMetricsService.ports[1].name`                                  |                                                         | `custom-metrics`    |
| `metricsOperatorControllerMetricsService.ports[1].port`                                  |                                                         | `443`               |
| `metricsOperatorControllerMetricsService.ports[1].targetPort`                            |                                                         | `custom-metrics`    |
| `metricsOperatorControllerMetricsService.ports[2].name`                                  |                                                         | `metrics`           |
| `metricsOperatorControllerMetricsService.ports[2].port`                                  |                                                         | `2222`              |
| `metricsOperatorControllerMetricsService.ports[2].protocol`                              |                                                         | `TCP`               |
| `metricsOperatorControllerMetricsService.ports[2].targetPort`                            |                                                         | `metrics`           |
| `metricsOperatorControllerMetricsService.type`                                           |                                                         | `ClusterIP`         |
| `metricsOperatorManagerConfig.controllerManagerConfigYaml.health.healthProbeBindAddress` | TODO  TODO  TODO                                        | `:8081`             |
| `metricsOperatorManagerConfig.controllerManagerConfigYaml.leaderElection.leaderElect`    | TODO  TODO  TODO                                        | `true`              |
| `metricsOperatorManagerConfig.controllerManagerConfigYaml.leaderElection.resourceName`   | TODO  TODO  TODO                                        | `3f8532ca.keptn.sh` |
| `metricsOperatorManagerConfig.controllerManagerConfigYaml.metrics.bindAddress`           | TODO  TODO  TODO                                        | `127.0.0.1:8080`    |
| `metricsOperatorManagerConfig.controllerManagerConfigYaml.webhook.port`                  | TODO  TODO  TODO                                        | `9443`              |
| `metricsOperatorController.nodeSelector`                                                 | add custom nodes selector to metrics operator           | `{}`                |
| `metricsOperatorController.tolerations`                                                  | add custom tolerations to metrics operator              | `[]`                |
| `metricsOperatorController.topologySpreadConstraints`                                    | add custom topology constraints to metrics operator     | `[]`                |

### Keptn Metrics Operator controller

| Name                                                                                  | Description                                       | Value                            |
| ------------------------------------------------------------------------------------- | ------------------------------------------------- | -------------------------------- |
| `metricsOperatorController.manager.containerSecurityContext`                          | Sets security context privileges                  |                                  |
| `metricsOperatorController.manager.containerSecurityContext.allowPrivilegeEscalation` |                                                   | `false`                          |
| `metricsOperatorController.manager.containerSecurityContext.capabilities.drop`        |                                                   | `["ALL"]`                        |
| `metricsOperatorController.manager.image.repository`                                  | specify registry for manager image                | `ghcr.io/keptn/metrics-operator` |
| `metricsOperatorController.manager.image.tag`                                         | select tag for manager image                      | `nil`                            |
| `metricsOperatorController.manager.livenessProbe`                                     | custom livenessprobe for manager container        |                                  |
| `metricsOperatorController.manager.readinessProbe`                                    | custom readinessprobe for manager container       |                                  |
| `metricsOperatorController.manager.resources`                                         | specify limits and requests for manager container |                                  |

### Global

| Name                      | Description                            | Value           |
| ------------------------- | -------------------------------------- | --------------- |
| `kubernetesClusterDomain` | overrides domain.local                 | `cluster.local` |
| `imagePullSecrets`        | global value for image registry secret | `[]`            |
