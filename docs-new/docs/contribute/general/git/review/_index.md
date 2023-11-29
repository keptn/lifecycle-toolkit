---
title: PR review process
description: How to navigate the review process
weight: 50
---

After you
[create your PR](../pr-create),
your PR must be approved and then merged
before it becomes part of the Keptn product.
This page discusses what you need to do during the review phase.

GitHub automatically assigns reviewers to your PR.
You can also tag other people in the description or comments.

Your PR will usually be reviewed within a couple of days,
but feel free to let us know about your PR
[via Slack](https://cloud-native.slack.com/channels/keptn-lifecycle-toolkit-dev).

You may want to record the PR number somewhere for future reference
although you can always find the PR in the GitHub lists of open and closed PRs.

## Draft (WIP) PRs

You may want to create a PR with work that is not ready for final review.
This happens when you want people to provide feedback on some of the work
before you finish the PR
or to protect the work you have done.

If this PR is not ready for review, click the "Still in progress?
s Convert to draft" string under the list of reviewers.
People can review the content but can not merge the PR
until you remove the "Draft" status.
The block of the PR that reports on checks will include the following item:

```console
This pull request is still a work in progress
Draft pull requests cannot be merged.
```

When the PR is ready to be reviewed, approved, and merged,
click the "Ready to review" button to remove the "Draft" status.
If you added "WIP" or "Draft" to the PR title, remove it now.

## Respond to review comments and suggestions

Watch for any comments that may be added by reviewers and implement or
respond to the recommended changes as soon as possible.
You should also check the build reports daily
to ensure that all tests are still passing.

* If a reviewer makes a GitHub suggestion and you agree with the change,
  click "Accept this change" to create a commit for that modification.
  Remember to include the DCO sign-off information in the commit message.

* You can make other changes using the GitHub editor.

* You can also work in your local branch to make modifications.
  However, if the PR has been modified on github,
  be sure to `pull` the changes back to your local branch
  before working in your local branch.

When your PR has the appropriate approvals,
it will be merged and the revised content should be published on the
website (as part of the `development` release)
within a few minutes.
You can now delete your local branch with the following command:

```console
git branch -d <branch-name>
```
