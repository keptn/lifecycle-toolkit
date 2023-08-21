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

Please check [Contribution Guidelines](content/en/contribute/docs/contribution-guidelines/_index.md).

## Building the Documentation Locally

Please check [Building the Documentation Locally](content/en/contribute/docs/local-building/index.md).

## Interacting with github

The documentation source is stored on github.com
and you use the standard github facilities to modify it.
The following notes are provided for people who need a little help.

### Fork and clone the repository

Perform the following steps to create a copy of this repository on your local machine:

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
    pasting in the URl you saved in the previous step:

    ```console
    git clone https://github.com/<UserName>/lifecycle-toolkit
    ```

    or

    ```console
    git clone git@github.com:<UserName>/lifecycle-toolkit.git
    ```

    Where <*UserName*> is your GitHub username.
    The lifecycle-toolkit directory is now available in the local directory.

4. Associate your clone with `upstream`.
   To do this, use the same string you used to clone your fork.

   * Be sure that you are in the root folder of the project
     and run *git status* to confirm that you are on the `main` branch.
   * Type the following to associate `upstream` with your clone,
     pasting in the string for the main repo that you copied above.:

     ```console
     git remote add upstream https://github.com/keptn/lifecycle-toolkit.git 
     ```

### Create a new branch and make your changes

Now you have your upstream setup in your local machine.
You need to create a local branch where you will make your changes.
Be sure that your branch is based on and sync'ed with `main`,
unless you intend to create a derivative PR:

The following sequence of commands does that:

```console
git checkout main
git pull upstream main
git push origin main
git checkout -b <my-new-branch>
```

Now you can make your changes, build them locally to check formatting,
then create a PR to add these changes to the documentation set

### Submitting new content through a pull request

If you have forked and cloned the repository,
you can modify the documentation or add new material
by editing the markdown file(s) in the local clone of your fork
and then submitting a *pull request (PR)*.

Note that you can also modify the source using the GitHub editor.
This is very useful when you want to fix a typo or make a small editorial change
but, if you are doing significant writing,
it is generally better to work on files in your local clone.

1. Execute the following and check the output to ensure that your branch is set up correctly:

   ```console
   git status
   ```

1. Do the writing you want to do in your local branch, checking the formatted version at `localhost:1314/docs-dev`.

1. When you have completed the writing you want to do, close all files in your branch and run `git status` to confirm
   that it correctly reflects the files you have modified, added, and deleted.

1. Add and commit your changes.
   Here, we commit all modified files but you can specify individual files to the
   `git add` command.
The `git commit -s` command commits the files and signs that you are contributing this intellectual property to the
Keptn project.

   ```console
   git add .
   git commit -s
   ```

   Use vi commands to add a description of the PR
   (should be 80 characters or less) to the commit.
   The title text should be prefixed with `docs:`
   to conform to our semantic commit scheme.
   This title is displayed as the title of the PR in listings.
   You can add multiple lines explaining the PR here but, in general,
   it is better to only supply the PR title here;
   you can add more information and edit the PR title
   when you create the PR on the GitHub UI page.

1. Push your branch to github:
   * If you cloned your fork to use SSH, the command is:

      ```console
      git push --set-upstream origin <branch-name>
      ```

      > **Note**
      You can just issue the `git push` command.
      Git responds with an error message that gives you the full command line to use; you can copy this command and
      paste it into your shell to push the content.

   * If you cloned your fork to use HTTP, the command is:

      ```console
      git push <need options/arguments>
      ```

1. Create the PR by going to the [lifecycle-toolkit](https://github.com/keptn/lifecycle-toolkit) GitHub repository.
   * You should see a yellow shaded area that says something like:

     ```console
     <GitID>:<branch> had recent pushes less than a minute ago
     ```

   * Click on the button in that shaded area marked:

     ```console
     Compare & pull request
     ```

   * Check that the title of the PR is correct; click the "Edit" button to modify it.
Add "WIP" (Work in Progress) or "Draft" to the title if the PR is not yet ready for general review.
   * Add a full description of the work in the PR, including any notes for reviewers, a reference to the relevant GitHub
      issue (if any), and tags for specific people (if any) who may be interested in this PR.
   * Carefully review the changes GitHub displays for this PR to ensure that they are what you want.
   * Click the green "Create pull request" button to create the PR.
      You may want to record the PR number somewhere for future reference although you can always find the PR in the
      GitHub lists of open and closed PRs.
   * GitHub automatically populates the "Reviewers" block.
   * If this PR is not ready for review, click the "Still in progress?
  Convert to draft" string under the list of
      reviewers.
      People can still review the content but can not merge the PR until you remove the "Draft" status.
   * The block of the PR that reports on checks will include the following item:

     ```console
     This pull request is still a work in progress
     Draft pull requests cannot be merged.
     ```

   * When the PR is ready to be reviewed, approved, and merged, click the "Ready to review" button to remove the "Draft"
      status.
  Then, if you added "WIP" or "Draft" to the PR title, remove it now.

1. Your PR should be reviewed within a few days.
   Watch for any comments that may be added by reviewers and implement or
   respond to the recommended changes as soon as possible.

   * If a reviewer makes a GitHub suggestion and you agree with the change, just click "Accept this change" to create a
      commit for that modification.
      You can also group several suggestions into a single commit using the GitHub tools.
   * You can make other changes using the GitHub editor or you can work in your local branch to make modifications.

      * If changes have been made using the GitHub editor, you will need to do a `git pull` request to pull those
         commits back to your local branch before you push the new changes.
      * After modifying the local source, issue the `git add .`, `git commit`, and `git push` commands to push your
         changes to github.

1. When your PR has the appropriate approvals, it will be merged and the revised content should be published on the
   website within a few minutes.

1. When your PR has been approved and merged,
   you can delete your local branch with the following command:

   ```console
   git branch -d <branch-name>
   ```
   
## Publishing

We are using Netlify to publish our pages.
There are 3 different types of publication:

1. pull request previews
2. development documentation aka staging (build of `main` branch) - [link](https://main.lifecycle.keptn.sh)
3. official documentation aka production (build of `page` branch) - [link](https://lifecycle.keptn.sh)

Within the navigation bar, we do have version links pointing to the different publications - if it makes sense.
For example, we are not linking from development and production to pull request previews.

### Pull request preview

- build: on each pull request with documentation changes
- build-environment: development
- config folder: [_default](./config/_default/)

The pull request preview will be generated if documentation files have been touched - this is configured in the [netlify.toml](../netlify.toml).

This preview should help contributors to inspect their changes within our usual page release.
Furthermore, it allows reviewers to inspect the rendered documentation without building it on their own.

### Development page

- build: on each push to `main` with documentation changes
- build-environment: main
- config folder: [main](./config/staging/)

This page reflects the current development status of the documentation.
It will be built regularly and can be easily accessed.

It should allow bleeding-edge users and contributors to see the current state and help with debugging etc.

### Official documentation

- build: on each push to `page` with documentation changes
- build-environment: production
- config folder: [production](./config/production/)

This documentation set contains all released versions of KLT and is stored in an orphaned branch called `page`.

Each version has its own `docs` folder named `docs-<version>`.
Except for the latest version which will be within the `docs` folder.

Each version-specific documentation contains a `version` file containing the version string.
This is important so we do know which version it contains - specifically important for `docs` of the latest version.


### Developer Certification of Origin (DCO)

Licensing is very important to open source projects.
It helps ensure the software continues to be available under the
terms that the author desired.

Keptn uses [Apache License 2.0](https://github.com/keptn/lifecycle-toolkit/blob/main/LICENSE) to strike a balance
between open contributions and allowing you to use the software however you would like to.

The license tells you what rights you have that are provided by the copyright holder.
It is important that the contributor fully understands what rights they are licensing and agrees to them.
Sometimes the copyright holder isn't the contributor, such as when the contributor is doing work on behalf of a company.

To make a good faith effort to ensure these criteria are met,
Keptn requires the Developer Certificate of Origin (DCO) process to be followed.

The DCO is an attestation attached to every contribution made by every developer.
In the commit message of the contribution, the developer simply adds a Signed-off-by statement and
thereby agrees to the DCO, which you can find below or at <http://developercertificate.org/>.

```text
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.


Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

#### DCO Sign-Off Methods

The DCO requires a sign-off message in the following format to appear on each commit in the pull request:

```text
Signed-off-by: Humpty Dumpty <humpty.dumpty@example.com>
```

The DCO text can either be manually added to your commit body,
or you can add either **-s** or **--signoff** to your usual git commit commands.
If you are using the GitHub UI to make a change you can add the sign-off message directly to the commit message
when creating the pull request.
If you forget to add the sign-off you can also amend a previous commit
with the sign-off by running **git commit --amend -s**.
If you've pushed your changes to GitHub already you'll need to force push your branch after this with **git push -f**.

**ATTRIBUTION**:

* <https://probot.github.io/apps/dco/>
* <https://github.com/opensearch-project/common-utils/blob/main/CONTRIBUTING.md>
* <https://code.asam.net/simulation/wiki/-/wikis/docs/project_guidelines/ASAM-DCO?version_id=c510bffb1195dc04deb9db9451112669073f0ba5>
* <https://thesofproject.github.io/latest/contribute/contribute_guidelines.html>

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
