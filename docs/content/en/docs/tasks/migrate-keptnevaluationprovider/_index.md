---
title: Migrate KeptnEvaluationProvider to KeptnMetricsProvider
description: Migrate deprecated KeptnEvaluationProvider to KeptnMetricsProvider.
icon: concepts
layout: quickstart
weight: 20
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

## Migrate KeptnEvaluationProvider to KeptnMetricsProvider

Due to the recent changes by splitting the `klt-operator` into two separate operators: `lifecycle-operator` and
`metrics-operator` (introduced in version 0.7.0), some of the functionality and behavior have been moved,
changed, or renamed. The `KeptnEvaluationProvider` CRD was deprecated and replaced by the `KeptnMetricsProvider`
CRD as a part of the `metrics.keptn.sh/v1alpha2` API group. During the upgrade from version 0.6.0 (or sooner)
to 0.7.0, there is a need to migrate manually `KeptnEvaluationProvider` CRs by using an external bash script:

```sh
curl -sL https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/.github/scripts/keptnevaluationprovider_migrator.sh | bash
```

This script will fetch all existing `KeptnEvaluationProvider` CRs and migrate them to the `KeptnMetricsProvider` CRs.
Additionally, the script stores the migrated manifests in your current working directory.

The script will also ask the user to apply the created manifests and delete the deprecated
`KeptnEvaluationProvider` CRs.

> **Note:**
Please be aware that only `KeptnEvaluationProvider` CRs from the cluster your kubecontext is pointing
to will be fetched and migrated.
