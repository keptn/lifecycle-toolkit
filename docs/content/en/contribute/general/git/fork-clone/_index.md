---
title: Fork and clone the repository
description: How to get a local version of the Keptn repository
weight: 20
---

Perform the following steps to create a copy
of the Keptn repository on your local machine:

1. Fork the Keptn repository:

     * Log into GitHub (or create a GitHub account and then log into it).
     * Go to the [Keptn lifecycle-toolkit repository](https://github.com/keptn/lifecycle-toolkit).
     * Click the **Fork** button at the top of the screen.
     * Choose the user for the fork from the options you are given,
       usually your GitHub ID.

   A copy of this repository is now available in your GitHub account.

2. Get the string to use when cloning your fork:

     * Click the green "Code" button on the UI page.
     * Select the protocol to use for this clone (either HTTPS or SSH).
     * A box is displayed that gives the URL for the selected protocol.
       Click the icon at the right end of that box to copy that URL.

3. Run the **git clone** command from the shell of a local directory
    to clone the forked repository to a directory on your local machine,
    pasting in the URl you saved in the previous step.

    For example, if you are using HTTPS:

    ```console
    git clone https://github.com/<UserName>/lifecycle-toolkit
    ```

    Or if you are using SSH:

    ```console
    git clone git@github.com:<UserName>/lifecycle-toolkit.git
    ```

    Where <*UserName*> is your GitHub username.
    The lifecycle-toolkit directory is now available in the local directory.

4. Associate your clone with `upstream`.

   * In a shell, go to the root folder of the project
     and run *git status* to confirm that you are on the `main` branch.
   * Type the following to associate `upstream` with your clone,
     pasting in the string for the main repo that you copied above.:

     ```console
     git remote add upstream https://github.com/keptn/lifecycle-toolkit.git
     ```

You are now ready to
[create a local branch](../branch-create)
and begin to create the software or documentation modifications.
