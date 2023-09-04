---
title: Codespaces
description: Use GitHub Codespaces to contribute documentation and software
weight: 200
---

Use GitHub codespaces as a pre-built and pre-configured development environment.
This is especially useful for Windows users
who may not have `make` installed.
It is also useful for Linux and MacOS users
who may not wish to download tools just to contribute to docs.

Review [this video](https://www.youtube.com/watch?v=sFNzOhZw7Eg) to see how this works.

[![Keptn + GitHub codespaces video](https://img.youtube.com/vi/sFNzOhZw7Eg/0.jpg)](https://www.youtube.com/watch?v=sFNzOhZw7Eg)

As shown in the video, the steps to set up a new Codespace are:

1. Create a fork of the repository.
   Keptn software and documentation are in the
   [link](https://github.com/keptn/lifecycle-toolkit)
   repository.
1. In your fork, click the green `Code` button
1. Switch to `Codespaces` tab and create a new codespace

You will be presented with a fully configured environment
with access to all the tools you require
to develop software or documentation for Keptn.

The interface is similar to that of
[Visual Studio Code](https://code.visualstudio.com/).

To develop or modify software or documentation, the steps are:

1. Make your modifications and test those modifications
1. Go back to Codespaces and click on the "Source Control" button on the left
1. Find the file(s) that you modified and click the "**+**" button
   to create a commit
   - Supply a commit message, adhering to the conventions for Keptn commits
   - Sign the commit by clicking the "**...**" button
     and selecting "Commit -> Commit Staged"
1. Click the "**...**" button again
   and select "Push" to push your changes to your fork
1. Go to the UI for your fork and create a PR in the normal way.
1. After your PR has been merged,
   go to your github page, select "Codespaces", and delete this codespace.
