---
comments: true
---

# Source File Structure

The source files for the Keptn documentation
are stored in the same GitHub repository as the source code for the software.
This page explains how the documentation source files are organized.

## Specifying the doc structure

The documentation builds are controlled by the
[mkdocs.yml](https://github.com/keptn/lifecycle-toolkit/blob/main/mkdocs.yml)
file in the root directory of the Keptn repository.
The documentation structure is defined under the `nav` section
that is roughly the second half of the file.
The following snippet illustrates how this is structured:

```markdown
...
nav:
  - Home:
      - index.md
  - Documentation:
      - docs/index.md
    ...
      - Use Cases:
          - docs/use-cases/index.md
          - Day 2 Operations: docs/use-cases/day-2-operations.md
          - Keptn + HorizontalPodAutoscaler: docs/use-cases/hpa.md
          - Keptn for non-Kubernetes deployments: docs/use-cases/non-k8s.md
   ...
```

* The `Home` item specifies the file for the `keptn.sh` landing page.
* The `Documentataion` item is the beginning of the specification
  for what files are are included and the order they will be listed.
* Each subitem to `Documentation` is a section of the docs
  as displayed in the left frame.
* Under each section are the individual pages,
  listed in the order they are displayed in the left frame.
  Each page line shows the title of the page
  that will be displayed in the left frame
  and the path to the source file.

  Note that the page title displayed in the main canvas
  is defined by the value of the H1 header in the page source file.
  When creating a new page or modifying the title,
  it is important to ensure that the title in the page source
  and the title in the `mkdocs.yml` file match.

## Coding the documentation

Keptn documentation is written using the Markdown language,
with each page written in a separate file.
Markdown supports many variants and the build tools we use
impose a few special requirements that are discussed here.

### Top of source requirements

The top of each documentation source file must include:

- A comments block
- A level 1 header with the title of the page.
  This must be preceded and followed by a single blank line.

  Remember that this is the title
  displayed in the main canvas of the doc page.
  The title displayed in the left frame
  is determined by the page line in the `mkdocs.yml` file.
  When creating or renaming a page,
  be sure that the title in the page source
  and the title in the `mkdocs.yml` file match.

For example:

```markdown
---
comments: true
---

# Uninstall

```

### Nested list format

Traditional markdown coding has sub-items in a list
lined up with the text of the item above.
MkDocs, however, requires that:

- Subitems be indented two additional spaces
- Subsequent paragraphs for this sub-item
  must be indented an additional two spaces
  in order to line up with the text of the first paragraph.
- Code blocks within the list must also be indented two spaces.

Alas, `make markdownlint-fix` does not currently understand
the formatting required for MkDocs
and so will "correct" the indentation
and you will end up with all list items
rendered as left-justified
and the numbers for numbered lists do not increment
so all numbered items begin with "1.".
Consequently, you must disable that particular markdownlint check
at the beginning of the list
and then re-enable it at the end of the list.

- Disable markdownlint check for list indentation:

    ```markdown
    <!-- markdownlint-disable MD007 -->`
    ```

- Re-enable markdownlint check for list indentation:

    ```markdown
    <!-- markdownlint-disable MD007 -->`
    ```

The following snippet illustrates this formatting.
This is for a numbered list with unordered sub-lists.
The same coding is used for an totally unordered list.

```markdown
<!-- markdownlint-disable MD007 -->

1. First item of numbered list.
   And more description

     * **subitem-1** -- Some description about this item,
       which may go on and on.
       Note that the mark is indented
       two spaces from the text of the preceding item.

         It could also go into a second paragraph
         which must be indented an additional two spaces.

          * **sub-subitem-1** -- Some description about this item,
            which may go on and on.
          * **sub-subitem-2** -- Some description about this item
            which may go on and on.

1. Second item of numbered list
   with more description.

     ```markdown
     Code blocks must also be indented two spaces
     from the text that introduces them.

1. Third item of numbered list
     ```
<!-- markdownlint-enable MD007 -->
```

## Primary documentation set

The source for the
[Keptn](https://lifecycle.keptn.sh/docs)
documentation is stored under
the `docs/docs` directory in the repository.

The subdirectories with content are:

- `assets`: This folder is used to save assets such as code examples that are used throughout the documentation.
  Many subfolders also contain an `assets` folder,
  usually containing graphics files (.png, .jpg, etc)
  to keep such files closer to the content where they are referenced.
- `components`: Information about how the different subcomponents of Keptn work
- `contribute`: Contains information on how to contribute software, tests, and documentation to Keptn
- `core-concepts`: A brief overview of Keptn, its features and use cases, and its history
- `getting-started`: Hands-on exercises that demonstrate the capabilities of Keptn
- `guides`: Guides and how-to material about using Keptn features
- `installation`: Requirements and instructions for installing and enabling Keptn
- `migrate`: Information to help users who are migrating to Keptn from Keptn v1
- `reference`: Reference pages for the CRDs and APIs that Keptn provides
- `use-cases`: Examples and exercises of using Keptn in specific scenarios

### Working with reference pages

The Keptn documentation includes two reference sections
that document the Keptn APIs and CRDs.
For background information, see:

- [Kubernetes API Concepts](https://kubernetes.io/docs/reference/using-api/api-concepts/)
- [Kubernetes API Reference](https://kubernetes.io/docs/reference/kubernetes-api/)

#### API Reference

The
[API Reference](../../reference/api-reference/index.md)
pages are autogenerated from source code.
This is a comprehensive list of all APIs and resources Keptn uses.

Descriptive text for the APIs is authored in the source code itself.
Each operator has its own API with different versions.
The source locations are:

- [Lifecycle API](https://github.com/keptn/lifecycle-toolkit/tree/main/lifecycle-operator/apis/lifecycle)
- [Metrics API](https://github.com/keptn/lifecycle-toolkit/tree/main/metrics-operator/api)
- [Options API](https://github.com/keptn/lifecycle-toolkit/tree/main/lifecycle-operator/apis/options)

The text is coded in a limited form of markdown.

To regenerate the autogenerated API reference docs,
execute the following script
in the root directory of your `lifecycle-toolkit` clone:

```shell
./.github/scripts/generate-crd-docs/generate-crd-docs.sh
```

#### CRD Reference

The [CRD Reference](../../reference/crd-reference/index.md) pages
describe the YAML manifest files used to populate resources
for the small set of CRDs that users must populate themselves.
These are currently authored manually
and should provide as much information as possible about the resource.
These are intended to be reference pages that are used regularly
by users after they are familiar with how to use Keptn.
For new users, the
[Guides](https://lifecycle.keptn.sh/docs/implementing/)
provide introductory material about how to use various features of Keptn.

A template to create a new CRD Reference page
is available [here](assets/yaml-crd-ref-template.md).

## Contribution guide

The source for the
[Contributing to Keptn](https://lifecycle.keptn.sh/contribute/)
guides
(which are accessed from the **Contributing** tab on the documentation page)
is stored under the `docs/contribute` directory.

The subdirectories of the contribution guide are:

- **general** (General information):
  Information that is applicable to all contributors,
  whether contributing software or documentation
- **software** (Software contributions):
  Information that is specific to software contributions
- **docs** (Documentation contributions):
  Information that is specific to documentation contributions

We also have *CONTRIBUTING.md* files located in the
home directory of the *lifecycle-toolkit* repository
(general and software development information)
and in the *lifecycle-toolkit/docs* directory
(documentation development information only).
These are the traditional locations for contributing information
but the amount of content we have was making them unwieldy,
so we are in the process of moving most content from these files
into the *Contributing guide*,
leaving links to the various sections in the *CONTRIBUTING.md* guides.

## Build strategy

This section discusses how the individual files and directories
are assembled into the documentation set.
See
[Published Doc Structure](./publish.md)
for information about the branches where the documentation is published.

All files in the directories are built.

The order in which the files are displayed
is determined by their order in the `nav` field
of the `mkdocs.yml` file.

## Subdirectory structure

Each subdirectory contains topical subdirectories for each chapter in that section.
Each topical subdirectory may contain:

- An *index.md* file that has the text for the section.
  If this is a subdirectory that contains subdirectories for other pages,
  the *index.md* file
  contains introductory content for the section.
- An *assets* subdirectory where graphical files for that topic are stored.
  No *assets* subdirectory is present if the topic has no graphics.
