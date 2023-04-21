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
`
* **kind** -- Resource type.
   Must be set to `KeptnConfig`.`

* **metadata**
  * **name** -- Unique name of this set of configurations.
    * Must be an alphanumeric string and, by convention, is all lowercase.
    * Can include the special characters `_`, `-`, (others?)
    * Should not include spaces.

* **spec**
  * **oTelCollectorUrl** -- The URL and port of the OpenTelemtry collector
  * **keptnAppCreationRequestTimeoutSeconds** --
     interval in which automatic app discovery searches for workloads
     to put into the same auto-generated [KeptnApp](app.md).

    Both values can be defined in one CRD
    or you can create separate `KeptnApp` CRDs for each one.

## Usage

## Example

### oTel example

This example specifies the URL of the OpenTelemetry collector:

```yaml
apiVersion: options.keptn.sh/v1alpha2
kind: KeptnConfig
metadata:
  name: otel-url
spec:
  OTelCollectorUrl: 'otel-collector:4317'
```

### App collection example

This example specifies the interval for
running the automatic app discovery:

```yaml
apiVersion: options.keptn.sh/v1alpha2
kind: KeptnConfig
metadata:
  name: app-collection-timeout
spec:
  OTelCollectorUrl: 'otel-collector:4317'
  keptnAppCreationRequestTimeoutSeconds: 40
```

## Files

API Reference:

* [KeptnTaskDefinition](../crd-ref/lifecycle/v1alpha3/_index.md#keptntaskdefinition)
* [KeptApp](app.md)

## Differences between versions

The `keptnAppCreationRequestTimeoutSeconds` field
is new in the `v1alpha2` App version.

## See also
