# ⚠️This repository is still in experimental phase - we discourage any use in production. The code is provided with no expectation of support or maintenance.

# Keptn Lifecycle Controller

This repository hosts a prototype of the Keptn Lifecycle Controller.
The goal of this prototype is to introduce a more “cloud-native” approach for pre- and post-deployment, as well as the concept of application health checks.

The Keptn Lifecycle Controller is composed of the following components:

- Keptn Lifecycle Operator
- Keptn Scheduler

The Keptn Lifecycle Operator contains several controllers for Keptn CRDs and a Mutating Webhook.
The Keptn Scheduler ensures that Pods are started only after the pre-deployment checks have finished.

## Architecture

![](./assets/architecture.jpg)

A Kubernetes Manifest, which is annotated with Keptn specific annotations, gets applied to the Kubernetes Cluster.
Afterward, the Keptn Scheduler gets injected (via Mutating Webhook), and Kubernetes Events for Pre-Deployment are sent to the event stream.
The Event Controller watches for events and triggers a Kubernetes Job to fullfil the Pre-Deployment.
After the Pre-Deployment has finished, the Keptn Scheduler schedules the Pod to be deployed.
The Application and Service Controllers watchfor the workload resources to finish and then generate a Post-Deployment Event.
After the Post-Deployment checks, SLOs can be validated using an interface for retrieving SLI data from a provider, e.g, [Prometheus](https://prometheus.io/).
Finally, Keptn Lifecycle Controller exposes Metrics and Traces of the whole Deployment cycle with [OpenTelemetry](https://opentelemetry.io/).

## How to use

TBD

## How to install

TBD

## License

Please find more information in the [LICENSE](LICENSE) file.