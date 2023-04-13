---
title: KeptnEvaluationDefinition
description: Define all evaluations associated with an application
weight: 30
---

A `KeptnEvaluationDefinition` defines evaluation tasks
that can be run by the Keptn Lifecycle Toolkit
as part of pre- and post-analysis phases of a workload or application.

## Yaml Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnEvaluationDefinition
metadata:
  name: <evaluation-name>
spec:
  source: prometheus | dynatrace | datadog
  objectives:
    - name: query-1
      query: "xxxx"
      evaluationTarget: <20
    - name: query-2
      query: "yyyy"
      evaluationTarget: >4
```

## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.  Must be set to `KeptnEvaluationDefinition`

* **metadata**
  * **name** -- Unique name of this evaluation
    such as `pre-deploy-eval` or `post-deploy-resource-eval`.
    * Must be an alphanumeric string and, by convention, is all lowercase.
    * Can include the special characters `_`, `-`, (others?)
    * Should not include spaces.

* **spec**
  * **source** -- Name of the data provider being used for this evaluation.
    The value of the `source` field must match
    the string used for the `name` field
    in the corresponding [KeptnEvaluationProvider](evaluationprovider.md) CRD.

    Each `KeptnEvaluationDefinition` CRD can use only one data provider;
    if you are using multiple data provider, you must create
    `KeptnEvaluationProvider` and `KeptnEvaluationDefinition` CRDs for each.

    Currently, you can only access one occurrance of each type of data provider
    in your KLT cluster.

  * **objectives** -- define the evaluations to be performed.
     Each objective is expressed as a `query` and an `evaluationTarget` value.

    * **query** -- Any query that is supported by the data provider.
    * **value** -- Desired value of the query,
       expressed as an arithmatic formula,
       usually less than (`<`) or greater than (`>`)

## Usage

## Examples

## Files

API Reference:

## Differences between versions

The `KeptnTaskDefinition` is the same for
all `v1alpha?` library versions.

## See also

* [KeptnEvaluationProvider](evaluationprovider.md)
