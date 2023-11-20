---
title: Upgrade
description: How to upgrade to the latest version of Keptn
layout: quickstart
weight: 45
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

If you installed the previous version of Keptn using `helm`,
you can upgrade to the latest version
by running the same command sequence used to install Keptn:

```shell
helm repo add klt https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn klt/klt \
   -n keptn-lifecycle-toolkit-system --create-namespace --wait
```

Use the `--set` flag or download and edit the `values.yaml` file
to modify the configuration as discussed on the
[Install Keptn](../install/) page.

> **Warning**
If you installed your Keptn instance from the Manifest,
additional steps are required to use the Helm Chart to upgrade.
Contact us on Slack for assistance.

## Upgrade to Helm from a manifest installation

> **Warning**
Upgrade to Helm from a manifest installation can cause loss
of Keptn custom resources.

Keptn v.0.7.0 and later can be installed with Helm charts;
Keptn v.0.8.3 can only be installed with Helm.
If you previously installed Keptn from manifests,
you can not directly upgrade with Helm but must back up your manifests,
then reinstall Keptn from a Helm chart and re-apply your manifests.

To start the upgrade process, follow the steps below:

1. To not loose all of your data, we encourage you to do a backup of the manifests,
which you applied to the cluster (`Pods`, `Deployments`,
`StatefulSets`, `DaemonSets`, `KeptnApps`,... ).
After the re-installation of Keptn with Helm, you can re-apply
these manifests and restart the Keptn deployment process.

1. Completely remove your Keptn installation with the following command sequence:

```shell
your-keptn-version=<your-keptn-version>
kubectl delete -f \
     https://github.com/keptn/lifecycle-toolkit/releases/download/$your-keptn-version/manifest.yaml
```

1. Use Helm to install a clean version of Keptn:

```shell
helm repo add klt https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn klt/klt -n keptn-lifecycle-toolkit-system --create-namespace --wait
```

For information about  advanced installation options, refer to
[Modify Helm configuration options](install.md).

1. After the installation finishes, restore the manifests from you backup

> **Warning** Please be aware that Keptn applications will start the deployment process from the
beginning and the system is not guaranteed to return
to the exact state it was in before re-installation, even if you created the backup correctly.

## Migrate from v0.6.0 to v0.7.0

Keptn Version v0.7.0
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
