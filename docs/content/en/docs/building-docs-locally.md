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