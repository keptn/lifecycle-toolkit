---
title: KeptnConfig
description: Define configuration values
weight: 20
---

`KeptnConfig` defines configuration values for the Keptn Lifecycle Toolkit.
Currently, it is used to configure the URL of the OpenTelemetry collector.

## Yaml Synopsis

```yaml
apiVersion: options.keptn.sh/v?alpha?
kind: KeptnConfig
metadata:
  name: <collector-name>
spec:
  OTelCollectorUrl: '<otelurl:port>'
```

## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.
   Must be set to `KeptnConfig`.`

* **metadata**
  * **name** -- Unique name of this collector.
    This is the name that identifies this collector
    in the `KeptnConfigList` CR.
    * Must be an alphanumeric string and, by convention, is all lowercase.
    * Can include the special characters `_`, `-`, (others?)
    * Should not include spaces.

## Usage

The `KeptnConfigList` CR is generated automatically
and contains a list of all valid `KeptnConfig` CR's.

## Example

```yaml
apiVersion: options.keptn.sh/v1alpha1
kind: KeptnConfig
metadata:
  name: keptnconfig-sample
spec:
  OTelCollectorUrl: 'otel-collector:4317'
```

## Files

API Reference:

* [KeptnTaskDefinition](../../crd-ref/lifecycle/v1alpha3/_index.md#keptntaskdefinition)

## See also
