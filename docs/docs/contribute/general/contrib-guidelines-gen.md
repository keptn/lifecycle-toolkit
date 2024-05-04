---
comments: true
---

# Contribution Guidelines

Before using Keptn
as a contributor to the Keptn `lifecycle-toolkit` repository,
it is expected that you comply with the guidelines while
making contributions towards the repository.

These guidelines are appropriate for both software and documentation.
For additional guidelines that are relevant only to documentation, see
[Contribution guidelines for documentation](../docs/contrib-guidelines-docs.md).

## Guidelines for contributing

* Always fork the repository then clone that fork to your local system
  rather than clone `main` directly.
  Keptn, like most open source projects,
  severely restricts who can push changes directly to the `main` branch
  to protect the integrity of the repository.
* Smaller PR's are easier to review and so generally get processed more quickly;
  when possible, chunk your work into smallish PR's.
  However, we recognize that documentation work sometimes requires larger PRs,
  such as when writing a whole new section or reorganizing existing files.
* You may want to squash your commits before creating the final PR,
  to avoid conflicting commits.
  This is **not mandatory**; the maintainers will squash your commits
  during the merge when necessary.
* Be sure that the description of the pull request itself
  is meaningful and clear.
  This helps reviewers understand each commit
  and provides a good history after the PR is merged.
* If your PR is not reviewed in a timely fashion,
  feel free to post a gentle reminder to the `#keptn` Slack channel.
* Resolve review comments and suggestions promptly.

If you see a problem and are unable to fix it yourself
or have an idea for an enhancement,
please create an issue on the GitHub repository.

* Provide specific and detailed information about the problem
  and possible solutions to help others understand the issue.
* When reporting a bug, provide a detailed list of steps to reproduce the bug.
  If possible, also attach screenshots that illustrate the bug.
* If you want to do the work on an issue,
  include that information in your description of the issue
  or in a comment to the issue.

## Proposing new work

* When proposing new work, start by creating an issue or ticket in the project's
  [issue tracker](https://github.com/keptn/lifecycle-toolkit/issues).
  If you would like to work on this issue, include that in a comment to the issue.
* Actively participate in the
  [refinement sessions](refinement-guide.md)
  that are part of the regularly scheduled
  [community meetings](https://community.cncf.io/keptn-community/).
    * In these sessions, everyone discusses the proposed work, whether it is a good idea,
      what exactly should be done and how it aligns with the project goals.
    * After the discussions, maintainers engage in a process known as **Scrum Poker**.
      This involves a voting mechanism where maintainers collectively assess the size
      and complexity of the proposed work, helping to decide whether it should proceed.
    * When the issue is refined, it can be assigned and work can begin.
* A PR should be created within a week of when an issue is assigned to you
  or the issue may be reassigned to someone else.
  When necessary, you can request an extension from the maintainers
  but the general expectation is that issues will be resolved
  soon after they are assigned.
* In general, you must complete one issue before you are assigned another one.
  Exceptions may be made when one issue is nearly completed
  but we discourage contributors claiming multiple issues at one time.
