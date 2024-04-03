# Release Checklist

This page should be used to document the current process to create a full release for Keptn.

## Checklist

1. (in-between every release: resolve conflicts in release-please-manifest file)
1. run security scans
1. release python runtime
1. update helm tests and helm chart docs with new python runtime version in renovate PR
1. merge python runtime renovate PRs
1. release deno runtime
1. update helm tests and helm chart docs with new deno runtime version in renovate PR
1. merge deno runtime renovate PRs
1. re-generate cert-manager Helm chart docs
1. release cert-manager
1. bump cert-manager chart version in charts repo and merge chart release PR
1. bump cert-manager chart version in main chart
1. manually update the cert manager library inside metrics and lifecycle operator with current commit hash from master
1. re-generate metrics-operator Helm chart docs
1. release metrics operator
1. bump metrics operator chart version in charts repo and merge chart release PR
1. bump metrics operator chart version in main chart
1. release scheduler
1. fix scheduler pr conflicts re-generate lifecycle-operator Helm chart docs
1. release lifecycle operator
1. bump lifecycle operator chart version in charts repo and merge chart release PR
1. bump lifecycle-operator chart version in main chart
1. helm dep update to fix Chart.lock files
1. klt release bump umbrella chart version in main chart
1. release the Keptn chart
1. update `stable` tag to the new keptn version and rebuild stable tag on readthedocs
1. verify docs
1. manually test examples
