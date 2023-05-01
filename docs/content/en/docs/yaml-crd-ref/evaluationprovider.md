---
title: KeptnEvaluationProvider (deprecated)
description: Define the evaluation provider
weight: 13
---

In earlier releases of the Lifecycle Toolkit,
`KeptnEvaluationProvider` defined the data provider
used by [KeptnEvaluationDefinition](

## Yaml Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnTaskDefinition
metadata:
  name: <task-name>
```

## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.
   Must be set to `KeptnTaskDefinition`

* **name** -- Unique name of this task.
  * Must be an alphanumeric string and, by convention, is all lowercase.
  * Can include the special characters `_`, `-`, (others?)
  * Should not include spaces.

## Usage

### Create secret text

## Examples

## Files

API Reference:

## Differences between versions

The `KeptnEvaluationProvider` is deprecated in the v1alpha3 API version.
`KeptnEvaluationDefinition` now gets provider information from the
[KeptnMetricsProvider](metricsprovider.md) CR.

## See also

* [KeptnEvaluationDefinition](evaluationdefinition.md)
