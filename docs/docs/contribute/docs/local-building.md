---
comments: true
---

# Build Documentation Locally

You can run MkDocs locally so that you can view the formatted version
of what you are writing before you push it to GitHub.
We provide a MkDocs run environment in a Docker container,
which simplifies the setup
and makes it easier to upgrade your local build environment
as the software is updated.

To set up a local MkDocs build:

1. Install Docker Desktop:

    * [Install on macOS](https://docs.docker.com/desktop/install/mac-install/)
    * [Install on Linux](https://docs.docker.com/desktop/install/linux-install/)
    * [Install on Windows](https://docs.docker.com/desktop/install/windows-install/)

1. Execute the following command from the root of your clone:

     ```shell
     make docs-serve
     ```

     It will continue running in its own shell.

    > **Note**
    To utilize the `Makefile`, you must have GNU **make**
    available on your local machine.
    Versions are available for all the usual operating systems.

1. Start contributing!
Note that MkDocs updates the rendered documentation each time you write the file.

1. Enter the following in a browser to view the website:

    `http://localhost:8000`

1. Use Ctrl+C to stop the local MkDocs server when you are done.

1. To restart the continuous build, run `make docs-serve` again.
