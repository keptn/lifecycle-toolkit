---
title: KeptnEvaluationProvider (deprecated)
description: Define the evaluation provider
weight: 13
---

In earlier releases of the Keptn Lifecycle Toolkit,
`KeptnEvaluationProvider` defined the data provider
used by [KeptnEvaluationDefinition](

## Yaml Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnTaskDefinition
metadata:
  name: <provider-name>
```

## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.
   Must be set to `KeptnTaskDefinition`

* **name** -- Unique name of this task.
  Names must comply with the
  [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
  specification.

## Usage

### Create secret text

## Examples

## Files

API Reference:

* [KeptnTaskDefinition](../crd-ref/lifecycle/v1alpha3/_index.md#keptntaskdefinition)

## Differences between versions

The `KeptnEvaluationProvider` is deprecated in the v1alpha3 API version.
`KeptnEvaluationDefinition` now gets provider information from the
[KeptnMetricsProvider](metricsprovider.md) CR.

## See also

* [KeptnEvaluationDefinition](evaluationdefinition.md)
