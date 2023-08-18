---
title: Create local branch
description: How to create a local branch and make modifications in it
weight: 30
---

After you
[fork and clone](../fork-clone)
the Keptn repository and set `upstream` in your local machine,
you need to create a local branch where you will make your changes.

## Create a new branch and make your changes

Be sure that your branch is based on and sync'ed with `main`,
unless you intend to create a derivative PR.
The following sequence of commands does that:

```console
git checkout main
git pull upstream main
git push origin main
git checkout -b <my-new-branch>
```

Execute the following and check the output
to ensure that your branch is set up correctly:

```console
git status
```

Now you can make your changes, build and test them locally,
then create a PR to add these changes to the documentation set.

* For documentation changes:
  * Do the writing you want to do in your local branch
  * Check the formatted version in your IDE
    or at `localhost:1314/docs-dev`
    to ensure that it is rendering correctly
    and that all links are valid..
    See [Build Documentation Locally](../../../docs/local-building)
    for more information.
  * Run `make markdownlint-fix` to check and fix the markdown code.

* For software changes:
  * Create the new code in your local branch.
  * Create and run unit tests for your new code.
  * Run other appropriate test to ensure that your code works correctly.

When you have completed the checking and testing of your work,
it is time to push your changes and create a PR that can be reviewed.
See [Create PR](../pr-create) for details.
