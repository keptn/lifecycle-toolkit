---
title: Source File Structure
description: Structure source files with Metadata
weight: 400
---

The source files for the Keptn documentation
are stored in the same github repository as the source code.
This page explains how the documentation source files are organized.

> **Note** The structure of the documentation
  and the source code for the documentation is evolving.
  You may find small discrepencies between
  what is documented here and what is currently implemented.

## Primary documentation set

The source for the
[Keptn Lifecycle Toolkit](https://lifecycle.keptn.sh/docs)
documentation is stored under
the *docs/content/en/docs* directory in the repository.

The subdirectories correspond to the contents listed in the right frame.
In the order they appear in the rendered docs, the subdirectories are:

* **intro** (Introduction to the Keptn Lifecycle Toolkit):
  A brief overview of Keptn, its features and use cases, and its history
* **getting-started** (Getting started):
  A hands-on exercise that demonstrates the capabilities of Keptn
* **tutorials** (Tutorials):
  Additional hands-on exercises to introduce new users to Keptn
* **install** (Installation and Upgrade):
  Requirements and instructions for installing and enabling Keptn
* **operate** (Operate Keptn):
  Guides about running and managing the Keptn cluster
* **implementing** (User Guides):
  This is currently a catch-all section
  for guides and how-to material about implementing Keptn features.
  It needs to be restructured and subdivided
* **architecture** (Architecture):
  Information about how Keptn actually works

  * **components** (Keptn Components)

    * **lifecycle-operator** (Keptn Lifecycle Operator)
    * **metrics-operator** (Keptn Metrics Operator)
    * **scheduler** (Keptn integration with Scheduling)

  * **deployment-flow** (Flow of deployment)
  * **keptn-apps** (KeptnApp and KeptnWorkload)
  * **cert-manager** (Keptn Certificate Manager)

* **crd-ref** (API Reference):
  Comprehensive information about all the APIs that define the Keptn CRDs.
  This section is auto-generated from source code
  and should never be modified in the *docs* directory.
  The source for the authored text can be modified
  in the source code files for the APIs
* **yaml-crd-ref** (CRD Reference):
  Reference pages for the CRs that users must populate.
  This is a subset of the CRDs documented in the *API Reference* section
* **migrate** (Migrating to the Keptn Lifecycle Toolkit):
  Information to help users who are migrating to Keptn
  from Keptn v1

## Contributing guide

The source for the
[Contributing to Keptn](https://lifecycle.keptn.sh/contribute/)
guides
(which are accessed from the **Contributing** tab on the documentation page)
is stored under the *docs/content/en/contribute* directory.

The subdirectories correspond to the contents listed in the right frame.
In the order they appear, the subdirectories are:

* **general** (General information about contributing):
  Information that is applicable to all contributors,
  whether contributing software or documentation

* **software** (Software contributions):
  Information that is specific to software contributions

* **docs** (Documentation contributions):
  Information that is specific to documentation contributions.

We also have *CONTRIBUTING.md* files located in the
home directory of the *lifecycle-toolkit* repository
(general and software development information)
and in the *lifecycle-toolkit/docs* directory
(documentation development information only).
These are the traditional locations for contributing information
but the amount of content we have was making them unwieldy
so we are in the process of moving most content from these files
into the *Contributing guide*,
leaving links to the various sections in the *CONTRIBUTING.md* guides.

## Build strategy

This section discusses how the individual files and directories
are assembled into the documentation set.
See
[Published Doc Structure](../publish)
for information about the branches where the documentation is published.

All files in the directories are built
except for files that are explicitly `ignored`
and files that include the `hidden: true` string in the file's metadata section.

The order in which the files are displayed
is determined by the value of the `weight` field
in the metadata section of *_index.md*, *index.md*,
and *topic.md* files that are located throughout the directory tree.

The metadata section of these files contains at least three fields.
As an example, the metadata section for the *Installation and upgrade* section
of the documentation includes the following fields:

```yaml
---
title: Installation and Upgrade
description: Learn how to install and upgrade the Keptn Lifecycle Toolkit
weight: 30
---
```

The meaning of these fields is:

* **title** -- title displayed for the section or file
* **description** -- subtext displayed for the section or subsection
* **weight** -- order in which section or subsection is displayed
  relative to other sections and subsections at the same level.

In this case, the weight of 30 means that this section is displayed
after sections that have weights of 29 and lower
and before sections that have weights of 31 and higher.
If two files have the same weight,
their order is determined alphabetically,
but this is a bad practice.
When you create a new section or a new page,
be sure to check the weight of the files
that immediately precede and follow this file
to be sure that you are not assigning the same weight to your new file.

The system for assigning weights for the docs landing page
allows for maximum flexibility as we create new sections:

* General introductory material uses weight values under 100.
* Guide material about using specific Keptn features
  use weight value of 2**.
* Reference material uses weight values of 5**.
* Other documents use weight values of 7**.

Some other fields are sometimes used in the metadata, including:

* **icon** -- optional field specifying an icon associated with the section
* **layout** -- layout template to be used for rendering the section
* **hidden** -- if set to `true`, this page is not included in the
  documentation set that is built
* **hidechildren** -- if set to `true`,
  the listing in the right margin of subsections of this page is omitted.
  In most cases, that listing is a convenient navigational aid for the reader
  but it can be omitted in special cases.

You can use these fields if you need them
but check the rendering carefully
to ensure that they are playing out as they should.

## Subdirectory structure

Each subdirectory contains topical subdirectories for each chapter in that section.
Each topical subdirectory may contain:

* An *_index.md* or *index.md* file that has the metadata discussed above
  plus the text for the section.
  If this is a subdirectory that contains subdirectories for other pages,
  the *_index.md* or *index.md* file
  contains introductory content for the section.
* An *assets* subdirectory where graphical files for that topic are stored.
  No *assets* subdirectory is present if the topic has no graphics.
