# Keptn Certificate Manager

Keptn Certificate Manager handles certificates for Keptn but can also handle certs for any other Kubernetes
resource.

<!-- markdownlint-disable MD012 -->
## Parameters

### Keptn Certificate Operator common

| Name                        | Description                                    | Value |
| --------------------------- | ---------------------------------------------- | ----- |
| `replicas`                  | customize number of replicas                   | `1`   |
| `nodeSelector`              | specify custom node selectors for cert manager | `{}`  |
| `tolerations`               | customize tolerations for cert manager         | `[]`  |
| `topologySpreadConstraints` | add topology constraints for cert manager      | `[]`  |

### Keptn Certificate Operator controller

| Name                       | Description                                                               | Value                                |
| -------------------------- | ------------------------------------------------------------------------- | ------------------------------------ |
| `containerSecurityContext` | Sets security context for the cert manager                                |                                      |
| `image.repository`         | specify repo for manager image                                            | `ghcr.io/keptn/certificate-operator` |
| `image.tag`                | select tag for manager container                                          | `v1.1.0`                             |
| `imagePullPolicy`          | select image pull policy for manager container                            | `Always`                             |
| `env.labelSelectorKey`     | specify the label selector to find resources to generate certificates for | `keptn.sh/inject-cert`               |
| `env.labelSelectorValue`   | specify the value for the label selector                                  | `true`                               |
| `livenessProbe`            | custom RBAC proxy liveness probe                                          |                                      |
| `readinessProbe`           | custom manager readiness probe                                            |                                      |
| `resources`                | custom limits and requests for manager container                          |                                      |

### Global

| Name                      | Description                            | Value           |
| ------------------------- | -------------------------------------- | --------------- |
| `kubernetesClusterDomain` | overrides domain.local                 | `cluster.local` |
| `imagePullSecrets`        | global value for image registry secret | `[]`            |
