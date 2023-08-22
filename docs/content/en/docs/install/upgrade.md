---
title: Upgrade
description: How to upgrade to the latest version of the Lifecycle Toolkit
layout: quickstart
weight: 45
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

If you installed the previous version of the Lifecycle Toolkit using `helm`,
you can upgrade to the latest version
by running the same command sequence used to install KLT:

```shell
helm repo add klt https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn klt/klt \
   -n keptn-lifecycle-toolkit-system --create-namespace --wait
```

Use the `--set` flag or download and edit the `values.yaml` file
to modify the configuration as discussed on the
[Install the Lifecycle Toolkit](../install/) page.

> **Warning**
If you installed your Lifecycle Toolkit instance from the Manifest,
additional steps are required to use the Helm Chart to upgrade.
Contact us on Slack for assistance.

## Migrate from v0.6.0 to v0.7.0

Keptn Lifecycle Toolkit Version v0.7.0
introduces the `metrics-operator`,
which is now separate from the `lifecycle-operator`.
Some functionality and behavior has been moved, changed, or renamed.

Specifically, the `KeptnMetricsProvider` CRD replaces
the now-deprecated `KeptnEvaluationProvider` CRD.
Consequently, you must manually migrate the existing functionality
to the `KeptnMetricsProvider` CRD.
Execute the following external bash script to do this:

```shell
curl -sL https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/.github/scripts/keptnevaluationprovider_migrator.sh | bash
```

This fetches and migrates all `KeptnEvaluationProvider` CRs
for the cluster at which your kubernetes  context is pointing.
If you have multiple clusters,
you must run this command for each one.

This script does the following:

* Fetch all existing `KeptnEvaluationProvider` CRs
* Migrate each to `KeptnMetricsProvider` CRs
* Store the migrated manifests in the current working directory
* Ask the user to apply the created manifests
  and delete the deprecated `KeptnEvaluationProvider` CRs.
