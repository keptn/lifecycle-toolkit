---
title: KeptnConfig
description: Define configuration values
weight: 20
---

`KeptnConfig` defines configuration values for the Keptn Lifecycle Toolkit.

## Yaml Synopsis

```yaml
apiVersion: options.keptn.sh/v?alpha?
kind: KeptnConfig
metadata:
  name: <configuration-name>
spec:
  OTelCollectorUrl: '<otelurl:port>'
  keptnAppCreationRequestTimeoutSeconds: <#-seconds>
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
  * **OTelCollectorUrl** -- The URL and port of the OpenTelemetry collector
    This field must be populated
    in order to export traces to the OpenTelemetry Collector.
  * **keptnAppCreationRequestTimeoutSeconds** --
    interval in which automatic app discovery searches for workloads
    to put into the same auto-generated [KeptnApp](app.md).
    The default value is 30 (seconds).

## Usage

Each cluster should have a single `KeptnConfig` CRD
that describes all configurations for that cluster.

## Example

### OTel example

This example specifies the URL of the OpenTelemetry collector
and that the automatic app discovery should be run every 40 seconds:

```yaml
apiVersion: options.keptn.sh/v1alpha2
kind: KeptnConfig
metadata:
  name: klt-config
spec:
  OTelCollectorUrl: 'otel-collector:4317'
  keptnAppCreationRequestTimeoutSeconds: 40
```

## Files

API Reference:

* [KeptnTaskDefinition](../crd-ref/lifecycle/v1alpha3/_index.md#keptntaskdefinition)

## Differences between versions

## See also

* [KeptnApp](../yaml-crd-ref/app.md)
* [OpenTelemetry observability](../implementing/otel.md)
* [Keptn automatic app discovery](../implementing/integrate/_index.md/#use-keptn-automatic-app-discovery)
