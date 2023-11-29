---
title: Create PR
description: Create and submit a PR with your changes
weight: 40
---

When you have completed the checking and testing of your work
on your local branch as described in
[Create local branch](../branch-create),
it is time to push your changes and create a PR that can be reviewed.
This is a two-step process:

1. [Push new content from your local branch](#push-new-content-from-your-local-branch)
1. [Create the PR on the github web site](#create-pr-on-the-github-web-site)

## Push new content from your local branch

The steps to push your new content from your local branch
to the repository are:

1. When you have completed the writing you want to do,
   close all files in your branch and run `git status` to confirm
   that it correctly reflects the files you have modified, added, and deleted
   and does not include any files that you do not want to push.

1. Switch back to the `main` branch in your repository,
   and ensure that it is up to date
   with the `main` Keptn branch:

   ```bash
   git remote add upstream https://github.com/keptn/lifecycle-toolkit.git
   git checkout main
   git pull upstream main
   ```

   Then update your feature branch from your local copy of `main` and push it:

   ```bash
   git checkout feature/123/foo
   git rebase main
   git push --set-upstream origin feature/123/foo
   ```

1. Add and commit your changes.
   The `git commit -s` command commits the files
   and signs that you are contributing this intellectual property
   to the Keptn project.
   See the DCO docs for more information.
   Here, we commit all modified files but you can specify individual files
   to the `git add` command.

   ```console
   git add .
   git commit -s
   ```

   Use vi commands to add a description of the PR
   (should be 80 characters or less) to the commit.
   The title text should be prefixed with an appropriate
   [commit type](#commit-types)
   to conform to our semantic commit scheme.
   This title is displayed as the title of the PR in listings.

   You can add multiple lines explaining the PR here but, in general,
   it is better to only supply the PR title here;
   you can add more information and edit the PR title
   when you create the PR on the GitHub UI page.

1. Push your branch to github.
   If you cloned your fork to use SSH, the command is:

      ```console
      git push --set-upstream origin <branch-name>
      ```

      > **Note**
      You can just issue the `git push` command.
      Git responds with an error message
      that gives you the full command line to use;
      you can copy this command and paste it into your shell to push the content.

## Create PR on the github web site

To create the actual PR that can be reviewed
and eventually merged, go to the
<https://github.com/keptn/lifecycle-toolkit> page.
You should see a yellow box that identifies your recent pushes.
Click the `Compare & pull request` button in that box
to open a PR template that you can populate.

> **Note**
  The PR template can also be found at `.github/pull_request_template.md`.

You need to provide the following information:

* Title for the PR.
   Follow the
  [conventional commit guidelines](https://www.conventionalcommits.org/en/v1.0.0/)
  for your PR title.
  * Title should begin with an appropriate
    [commit type](#commit-types).feature type.
  * The first word after the feature type should be lowercase.

    An example for a pull request title is:

    ```bash
    feat(api): new endpoint for feature X
    ```

* Full description of what the PR is about.

  * Link to relevant GitHub issue(s).
     Use the phrase `Closes <issue>` for this link;
       is ensures that the issue is closed when this PR is merged.
        this PR does not completely satisfy the issue,
       e some other phrasing for the link to the issue.
  * Describe what this PR does,
    including related work that will be in other PRs.
  * If you changed something that is visible to the user,
    add a screenshot.
  * Describe tests that are included or were run to test this PR.
  * Anything else that will help reviewers understand
    the scope and purpose of this PR.

* If you have **breaking changes** in your PR,
  it is important to note them in both the PR description
  and in the merge commit for that PR.

   When pressing "squash and merge",
   you have the option to fill out the commit message.
   Please use that feature to add the breaking changes according to the
   [conventional commit guidelines](https://www.conventionalcommits.org/en/v1.0.0/).
   Also, please remove the PR number at the end and just add the issue number.

   An example for a PR with breaking changes and the according merge commit:

   ```bash
   feat(bridge): New button that breaks other things (#345) 

   BREAKING CHANGE: The new button added with #345 introduces
   new functionality that is not compatible with the previous
   type of sent events.
   ```

   If your breaking change can be explained in a single line,
   you can also use this form:

   ```bash
   feat(bridge)!: New button that breaks other things (#345)
   ```

   Following these guidelines helps us create automated releases
   where the commit and PR messages are directly used in the changelog.

When you have filled in the PR template,
you should also quickly scroll down to see the changes
that are being made with this commit
to ensure that this PR contains what you want reviewed.

When you are satisfied that your PR is ready for review,
click the `Create pull request` button.

### Commit Types

**Type** can be:

* `feat`: a new feature
* `fix`: a bug fix
* `build`: changes that affect the build system or external dependencies
* `chore`: other changes that don't modify source or test files
* `ci`: changes to our CI configuration files and scripts
* `docs`: documentation only changes
* `perf`: a code change that improves performance
* `refactor`: a code change that neither fixes a bug nor adds a feature
* `revert`: reverts a previous commit
* `style`: changes that do not affect the meaning of the code
* `test`: adding missing tests or correcting existing tests

## Check PR build

As soon as you create the PR,
a number of tests and checks are run.
Be sure to check the results immediately
and fix any problems that are found.
Click the `Details` link on the line for each each failed test
to get details about the errors found.

The most common errors for documentation PRs are:

* Markdown errors found by `markdownlint`.
  Most of these can be fixed
  by running `make markdownlint-fix` on your local branch
  then pushes the changes.
* Cross-reference errors.
  To quickly find the errors in the report,
  search for the `dead` string on the `Details` page.

When you have resolved all build errors
you move into the
[PR review process](../review).
