---
title: KeptnEvaluationProvider
description: Define the data provider for evaluations
weight: 33
---

A `KeptnEvaluationProvider` identifies an evaluation provider
that provides data for evaluations done
during the pre- and post-analysis phases of a workload or application.

## Yaml Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnEvaluationProvider
metadata: <provider-name>
source:
  name: prometheus | dynatrace | datadog
spec:
  targetServer: "http://prometheus-k8s.monitoring.svc.cluster.local:9090"
  secretName: prometheusLoginCredentials
```

## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.  Must be set to `KeptnTaskDefinition`

* **metadata**
  * **name** -- Unique name of this data provider.
    * Must be an alphanumeric string and, by convention, is all lowercase.
    * Can include the special characters `_`, `-`, (others?)
    * Should not include spaces.

* **spec**
  * **source** -- Type of data provider being used
    Note that you can configure one each of the different providers
    for your KLT cluster
    but you cannot currently configure more than one instance
    of any of the providers.
  * **targetServer** -- Location of the data provider
  * **secretName** -- Secret used to access the data provider

## Usage

## Files

API Reference:

* [KeptnEvaluationProvider](../../crd-ref/lifecycle/v1alpha3/#keptnevaluationprovider)

## Differences between versions

The `KeptnEvaluationProvider` CRD is the same for
all `v1alpha?` library versions.

## See also

* [KeptnEvaluationDefinition](evaluationdefinition.md)
