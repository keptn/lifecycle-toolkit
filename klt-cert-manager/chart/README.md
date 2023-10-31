# Keptn Certificate Manager

Keptn Certificate Manager handles certificates for Keptn but can also handle certs for any other Kubernetes
resource.

<!-- markdownlint-disable MD012 -->
## Parameters

### Global parameters

| Name                       | Description                                                               | Value |
| -------------------------- | ------------------------------------------------------------------------- | ----- |
| `global.imageRegistry`     | Global container image registry                                           | `""`  |
| `global.imagePullSecrets`  | Global Docker registry secret names as an array                           | `[]`  |
| `global.commonLabels`      | Common annotations to add to all Keptn resources. Evaluated as a template | `{}`  |
| `global.commonAnnotations` | Common annotations to add to all Keptn resources. Evaluated as a template | `{}`  |

### Keptn Certificate Operator common

| Name                        | Description                                    | Value           |
| --------------------------- | ---------------------------------------------- | --------------- |
| `nodeSelector`              | specify custom node selectors for cert manager | `{}`            |
| `replicas`                  | customize number of replicas                   | `1`             |
| `tolerations`               | customize tolerations for cert manager         | `[]`            |
| `topologySpreadConstraints` | add topology constraints for cert manager      | `[]`            |
| `kubernetesClusterDomain`   | overrides cluster.local                        | `cluster.local` |
| `annotations`               | add deployment level annotations               | `{}`            |
| `podAnnotations`            | adds pod level annotations                     | `{}`            |

### Keptn Certificate Operator controller

| Name                       | Description                                                               | Value                        |
| -------------------------- | ------------------------------------------------------------------------- | ---------------------------- |
| `containerSecurityContext` | Sets security context for the cert manager                                |                              |
| `env.labelSelectorKey`     | specify the label selector to find resources to generate certificates for | `keptn.sh/inject-cert`       |
| `env.labelSelectorValue`   | specify the value for the label selector                                  | `true`                       |
| `image.registry`           | specify the container registry for the certificate-operator image         | `ghcr.io`                    |
| `image.repository`         | specify repo for manager image                                            | `keptn/certificate-operator` |
| `image.tag`                | select tag for manager container                                          | `v1.2.0`                     |
| `imagePullPolicy`          | select image pull policy for manager container                            | `Always`                     |
| `livenessProbe`            | custom RBAC proxy liveness probe                                          |                              |
| `readinessProbe`           | custom manager readiness probe                                            |                              |
| `resources`                | custom limits and requests for manager container                          |                              |
