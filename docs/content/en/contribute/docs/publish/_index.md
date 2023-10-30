---
title: Published Doc Structure
description: Structure of the published documentation
weight: 450
---

New writing goes to the `main` branch and can be viewed on the Releases -> development dropdown menu.
We have staging and production levels for our documentation which are as follows:

* **Latest:** official documentation of the current Keptn release
  * [link](https://lifecycle.keptn.sh):
      This is the build of the `page` branch.

* **Development:** documentation being staged for the next Keptn release
  * [link](https://main.lifecycle.keptn.sh):
   This is the latest build of the `main` branch.

* **Previous versions:** documentation for earlier releases.
   These are listed at [link](https://github.com/keptn/lifecycle-toolkit/tree/page/docs/content/en).

* **Contribute:** current version of the "Contribute" guide
   that is available from a tab on the documentation site.

Let's take a look what happens when your changes are merged in `main` and `page` branch respectively.

## Latest -- Official documentation (Production)

This set of documentation pertains to the latest Keptn release and resides within an
isolated branch known as `page`.
When a new version of Keptn is launched, the contents of the `development`
branch are rolled into this branch.
Additionally, it's important to recognize that any
document changes made using the "Edit this page" feature are seamlessly integrated into this branch.

This uses the `latest` label so that links to a doc page
remain valid across software and documentation updates.

* build: on each push to `page` with documentation changes
* build-environment: production
* config folder: [production](https://github.com/keptn/lifecycle-toolkit/tree/main/docs/config/production)

A new version is generated when we push the `main` branch to production to release a new version of the docs page.
This means, that the content of the old version on the `page` branch will be copied over
to a `docs-<version>` folder and the new version will be pushed into the `docs` folder.
This way, no changes or older versions get overwritten.

## Development documentation (Staging)

This page contains the documentation being staged for the next Keptn release.
It contains information about new and changed features and functionality
as well as general documentation improvements.
It is built regularly and can be easily accessed from the `Releases` tab on the documentation site.

* build: on each push to `main` with documentation changes
   from a user's local branch, from the github editor, or from codespaces
* build-environment: main
* config folder: [staging](https://github.com/keptn/lifecycle-toolkit/tree/main/docs/config/staging)

This version represents the pre-release iteration of the documentation for the upcoming Keptn release.
Pull requests originating from a user's local branch, the GitHub editor, or codespaces are merged into this branch.

When a new Keptn version is officially launched, this branch is elevated to the status of `latest`.
In exceptional cases, a pull request that includes vital documentation enhancements may be discreetly
advanced to `latest` without the need for a software release.

## Previous Versions

Keptn documentation is versioned.
By default, the version for the current Keptn release
is displayed on the documentation page but users can select other versions from the Releases tab.
The previous versions of the Keptn Documentation are available [here](https://github.com/keptn/lifecycle-toolkit/tree/page/docs/content/en).
