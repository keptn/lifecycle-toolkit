---
comments: true
---

# Published Doc Structure

New writing goes to the `main` branch and can be viewed on the Keptn website under the `latest` version.
We have staging and production levels for our documentation which are as follows:
<!-- markdownlint-disable MD007 -->
* **Stable:** official documentation of the current Keptn release
    * [link](https://keptn.sh/latest/): This is the build of the latest `keptn-*` Git tag.

* **Latest:** documentation being staged for the next Keptn release
    * [link](https://keptn.sh/latest/): This is the latest build of the `main` branch.
<!-- markdownlint-enable MD007 -->
Let's take a look what happens when your changes are merged in `main`.

## Stable - Official documentation (Production)

This set of documentation pertains to the most recent stable Keptn release and resides within the
latest `keptn-*` tagged commit on GitHub.
When a new version of Keptn is released, the contents of the `main` branch are published as
the new `stable` release documentation.

* build: on each release of Keptn (`keptn-*` Git tag)

A new version is generated when we push a new `keptn-*` tag to the GitHub repository through the
[release pipeline](https://github.com/keptn/lifecycle-toolkit/tree/main/.github/workflows/release.yml).
This means, that the content of the old Keptn tag will be replaced by the newly released version
that was tagged on the `main` branch.

## Latest - Development documentation (Staging)

This page contains the documentation being staged for the next Keptn release.
It contains information about new and changed features and functionality
as well as general documentation improvements.
It is built regularly and can be easily accessed from the version dropdown menu on the documentation site.

* build: on each push to `main`

This version represents the pre-release iteration of the documentation for the upcoming Keptn release.

## Previous Versions

Keptn documentation is versioned but only the documentation for the current stable release of Keptn is published.
By default, the version for the current Keptn release
is displayed on the documentation page.
