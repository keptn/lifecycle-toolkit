---
title: Migrate KeptnEvaluationProvider to KeptnMetricsProvider
description: Migrate deprecated KeptnEvaluationProvider to KeptnMetricsProvider.
icon: concepts
layout: quickstart
weight: 20
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

## Migrate KeptnEvaluationProvider to KeptmMetricsProvider

Due to the recent changes by splitting the `klt-operator` into two separate operators: `klt-operator` and
`metric-operator` (introduced in version 0.7.0), some of the functionality and behavior have been moved, changed, or renamed. The `KeptnEvaluationProvider` CRD was deprecated and replaced by the `KeptnMetricsOperator` CRD as a part of the `metrics.keptn.sh/v1alpha2` API group. During the upgrade from version 0.6.0 (or sooner) to 0.7.0, there is a need to migrate manually `KeptnEvaluationProvider` CRDs by using an external bash script:

```sh
curl -sL https://raw.githubusercontent.com/keptn/lifecycle-toolkit/epic/split-metrics-operator/.github/scripts/keptnevaluationprovider_migrator.sh | bash
```

This script will fetch all existing `KeptnEvaluationProvider` CRs, migrate them to `KeptnMetricsProvider`
and store the migrated manifests to a manifests file of your current working directory. Additionally, it
will directly apply the created manifests and also delete the deprecated`KeptnEvaluationProvider` CRs,
if the user wishes to do so.

**Note:** Please be aware, that only `KeptnEvaluationProvider` CRs from the cluster your kubecontext is
pointing to will be fetched and migrated.