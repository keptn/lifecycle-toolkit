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

## [Guidelines for contributing](https://github.com/keptn/lifecycle-toolkit/tree/main/docs/content/en/docs/guidelines-for-contributing.md)

## Building the documentation locally

You can run Docsy locally so that you can view the formatted version
of what you are writing before you push it to github.
We provide a Docsy run environment in a Docker container,
which simplifies the set up
and makes it easier to upgrade your local build environment
as the software is updated.

To set up a local Docsy build:

1. Install Docker Desktop:

   * [Install on macOS](https://docs.docker.com/desktop/install/mac-install/)
   * [Install on Linux](https://docs.docker.com/desktop/install/linux-install/)
   * [Install on Windows](https://docs.docker.com/desktop/install/windows-install/)

1. Build the Keptn Docsy repo:

   ```console
   make build
   ```

   > **Note:**
   To utilize the `Makefile`, you must have GNU **make**
   available on your local machine.
   Versions are available for all the usual Operating Systems.

1. Execute the following command from the `docs` folder of your clone:

   ```console
   make server
   ```

   It will continue running in its own shell.

1. Start contributing!
Note that Hugo updates the rendered documentation each time you write the file.

1. Enter the following in a browser to view the website:

    `http://localhost:1314/docs-dev/`

   > **Note:**
   By default, Hugo serves the local docs on port 1313.
   We have modified that port for the lifecycle-toolkit docs
   to avoid conflicts with the keptn.github.io docs, which use
   port 1313 for local builds.

1. Use Ctrl+C to stop the local Hugo server when you are done.

1. To restart the continuous build:

   * Restart Docker-Desktop, if necessary
   * If changes have been made to the build tools:
     * make clone
     * make build
   * Run `make server`

### Building markdown files without Hugo

The Hugo generator described above only renders
the markdown files under the */content/docs* directory.
If you need to render another markdown file
(such as this *CONTRIBUTING.md* file)
to check your formatting, you have the following options:

* If you are using an IDE to author the markdown text,
     use the markdown preview browser for the IDE.
* You can push your changes to GitHub
     and use the GitHub previewer (*View Page*).
* You can install and use the
     [grip](https://github.com/joeyespo/grip/blob/master/README.md) previewer
     to view the rendered content locally.
     When *grip* is installed,
     you can format the specified file locally
     by running the following in its own shell:

     ```console
     grip <file>.md
     ```

     Point your browser at `localhost:6419` to view the formatted file.
     The document updates automatically
     each time you write your changes to disk.

## [Interacting with github](https://github.com/keptn/lifecycle-toolkit/tree/main/docs/content/en/docs/interacting-with-github.md)

## Source file structure

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
* **weight** -- order in which section or subsection is desplayed relative to other sections and
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
No *assets* subdirectory is present if the topic has no graphics.

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
