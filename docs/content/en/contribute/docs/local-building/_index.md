---
title: Build Documentation Locally
description: This guide explains how to create a local version of the documentation
weight: 400
---

## Building the documentation locally

You can run Docsy locally so that you can view the formatted version
of what you are writing before you push it to github.
We provide a Docsy run environment in a Docker container,
which simplifies the setup
and makes it easier to upgrade your local build environment
as the software is updated.

To set up a local Docsy build:

1. Install Docker Desktop:

   * [Install on macOS](https://docs.docker.com/desktop/install/mac-install/)
   * [Install on Linux](https://docs.docker.com/desktop/install/linux-install/)
   * [Install on Windows](https://docs.docker.com/desktop/install/windows-install/)

1. Execute the following command from the `docs` folder of your clone:

   ```shell
   make server
   ```

   It will continue running in its own shell.

   > **Note**
   To utilize the `Makefile`, you must have GNU **make**
   available on your local machine.
   Versions are available for all the usual Operating Systems.

1. Start contributing!
Note that Hugo updates the rendered documentation each time you write the file.

1. Enter the following in a browser to view the website:

    `http://localhost:1314`

   > **Note**
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

     ```shell
     grip <file>.md
     ```

     Point your browser at `localhost:6419` to view the formatted file.
     The document updates automatically
     each time you write your changes to disk.
