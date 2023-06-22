---
title: Source File Structure
description: Structure source files with Metadata
weight: 400
---

The source files for the [Keptn Lifecycle Toolkit](https://lifecycle.keptn.sh/docs) are stored under
the *docs/content/en/docs* directory in the repository.
The build strategy is to build everything except for files that are explicitly ignored
and files that include the `hidden: true` string in the file's metadata section

The order in which the files are displayed is determined by the value of the `weight` field
in the metadata section of *_index.md* and *index.md* files located throughout the directory tree.

The metadata section of these files contains at least three fields.
As an example, the metadata section for the *Concepts* section of the documentation includes the following fields:

```yaml
title: Concepts
description: Learn about underlying concepts of the keptn lifecycle toolkit.
icon: concepts
layout: quickstart
weight: 50
```

The meaning of these fields is:

* **title** -- title displayed for the section or file
* **description** -- subtext displayed for the section or subsection
* **icon** -- An optional field specifying an icon associated with the section
* **layout** -- The layout template to be used for rendering the section
* **weight** -- order in which section or subsection is displayed relative to other sections and
  subsections at the same level.

In this case, the weight of 50 means that this section is displayed
after sections that have weights of 49 and lower
and before sections that have weights of 51 and higher.
If two files have the same weight,
their order is determined alphabetically,
but this is a bad practice.

Some other fields are sometimes used in the metadata.

### Top level structure

The current tools do not support versioning.
To work around this limitation, the docs are arranged with some general topics that generally apply to all releases and
then subsections for each release that is currently supported.

The system for assigning weights for the docs landing page is:

* General introductory material uses weight values under 100.
* Sections for individual releases use weight values of 9**.
* Sections for general but advanced info use weight value of 1***.

### Subdirectory structure

Each subdirectory contains topical subdirectories for each chapter in that section.
Each topical subdirectory contains:

* An *index.md* file that has the metadata discussed above plus the text for the section
* An *assets* subdirectory where graphical files for that topic are stored.
* No *assets* subdirectory is present if the topic has no graphics.
