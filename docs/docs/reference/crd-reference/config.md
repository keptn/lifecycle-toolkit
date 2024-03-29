---
comments: true
---

# KeptnConfig

`KeptnConfig` defines Keptn configuration values.

## Yaml Synopsis

```yaml
apiVersion: options.keptn.sh/v1alpha1
kind: KeptnConfig
metadata:
  name: <configuration-name>
spec:
  OTelCollectorUrl: '<otelurl:port>'
  keptnAppCreationRequestTimeoutSeconds: <#-seconds>
  cloudEventsEndpoint: <endpoint>
  blockDeployment: true | false
  observabilityTimeout: <duration>
```

## Fields

* **apiVersion** -- API version being used.
* **kind** -- Resource type.
  Must be set to `KeptnConfig`.

* **metadata**
    * **name** -- Unique name of this set of configurations.
      Names must comply with the
      [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
      specification.

* **spec**
    * **OTelCollectorUrl** -- The URL and port of the OpenTelemetry collector.
      This field must be populated in order to export traces to the OpenTelemetry Collector.
    * **keptnAppCreationRequestTimeoutSeconds** --
      Interval in which automatic app discovery searches for [workloads](https://kubernetes.io/docs/concepts/workloads/)
      to put into the same auto-generated [KeptnApp](app.md).
      The default value is 30 (seconds).
    * **cloudEventsEndpoint** -- Endpoint where the lifecycle operator posts Cloud Events.
    * **blockDeployment** -- If set to `true` (default), application deployment is blocked until the
      pre-deployment tasks and evaluations succeed.
      You can set this field to `false` when building up
      your pre-deployment tasks and evaluations
      so that your application is deployed
      even if the pre-deployment tasks and/or evaluations fail.
      For more information see the
      [non-blocking deployment section](../../components/lifecycle-operator/keptn-non-blocking.md).
    * **observabilityTimeout** -- specifies the maximum time
      to observe the deployment phase of
      [KeptnWorkload](../api-reference/lifecycle/v1/index.md).
      The value supplied should specify the unit of measurement;
      for example, `5m` indicates 5 minutes and `1h` indicates 1 hour.
      If the workload is not deployed successfully within this time frame,
      it is considered to be failed.

## Usage

Each cluster should have a single `KeptnConfig` CRD that describes all configurations for that cluster.

## Example

This example specifies:

* the URL of the OpenTelemetry collector
* automatic app discovery that should be run every 40 seconds
* CloudEvents endpoint URL
* blocking functionality of the deployment of the application is disabled in case
  of the pre-deployment task or evaluation failure

```yaml
apiVersion: options.keptn.sh/v1alpha1
kind: KeptnConfig
metadata:
  name: keptn-config
spec:
  OTelCollectorUrl: 'otel-collector:4317'
  keptnAppCreationRequestTimeoutSeconds: 40
  cloudEventsEndpoint: 'http://endpoint.com'
  blockDeployment: false
  observabilityTimeout: 10m
```

## Files

API Reference:

* [KeptnTaskDefinition](../api-reference/lifecycle/v1/index.md#keptntaskdefinition)

## Differences between versions

## See also

* [KeptnApp](./app.md)
* [OpenTelemetry observability](../../guides/otel.md)
* [Keptn automatic app discovery](../../guides/auto-app-discovery.md)
* [Keptn non-blocking deployment](../../components/lifecycle-operator/keptn-non-blocking.md)
