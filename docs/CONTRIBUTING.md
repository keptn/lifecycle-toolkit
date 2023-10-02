# Contribute to the Keptn documentation

This document provides information about contributing to
the [Keptn Lifecycle Toolkit documentation](https://lifecycle.keptn.sh/docs/),
which is part of the [Keptn](https://keptn.sh) website.

The Keptn Lifecycle Toolkit documentation is authored with
[markdown](https://www.markdownguide.org/basic-syntax/)
and rendered using the Hugo
[Docsy](https://www.docsy.dev/) theme.

We welcome and encourage contributions of all levels.
You can make modifications using the GitHub editor;
this works well for small modifications but,
if you are making significant changes,
you may find it better to fork and clone the repository
and make changes using the text editor or IDE of your choice.
You can also run the Docsy based website locally
to check the rendered documentation
and then push your changes to the repository as a pull request.

If you need help getting started,
feel free to ask for help on the `#help-contributing` or `#keptn-docs` channels on the [Keptn Slack](https://keptn.sh/community/#slack).
We were all new to this once and are happy to help you!

## Guidelines for Contributing

Please check [Contribution Guidelines](content/en/contribute/docs/contrib-guidelines-docs/_index.md).

## Building the Documentation Locally

Please check [Building the Documentation Locally](content/en/contribute/docs/local-building/_index.md).

## Interacting with github

The documentation source is stored on github.com
and you use the standard github facilities to modify it.
Please check [Working with Git](content/en/contribute/general/git/_index.md).

## Choosing the correct branch when contributing

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

This set of documentation pertains to the latest KLT release and resides within an
isolated branch known as `page`.
When a new version of KLT is launched, the contents of the `development`
branch are rolled into this branch.
Additionally, it's important to recognize that any
document changes made using the "Edit this page" feature are seamlessly integrated into this branch.

This uses the `latest` label so that links to a doc page
remain valid across software and documentation updates.

* build: on each push to `page` with documentation changes
* build-environment: production
* config folder: [production](./config/production/)

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
* config folder: [main](./config/staging/)

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

### Developer Certification of Origin (DCO)

All commits must be accompanied by a DCO sign-off.
 See
[DCO](content/en/contribute/general/dco)
for more information.

## Source File Structure

Please check [Source File Structure](content/en/contribute/docs/source-file-structure/_index.md)..

## Guidelines for working on documentation in development versus already released documentation

[This material will be provided when we define the versioning scheme to use]

### Documentation for new features

Most documentation changes should be made to the docs-dev branch,
which means creating a PR in the `lifecycle-toolkit` repository
under the `docs/content/en/docs` directory.
You can view the local build as described above.
We are releasing new versions of the software frequently
so this makes new content available reasonably quickly.

### Documentation for published docs

If a critical problem needs to be solved immediately,
you can modify the documentation source in the sandbox.
In this case, modify the files in the
`keptn-sandbox/lifecycle-toolkit-docs` repository directly.
You can view these changes locally on the `localhost:1314` port.

Note that changes made to the docs in the sandbox
will be overwritten so the same changes should be applied
to the corresponding doc source in the `lifecycle-toolkit` documentation.
