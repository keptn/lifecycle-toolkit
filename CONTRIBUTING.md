# Contributing to the Keptn Lifecycle Toolkit

We are thrilled to have you join us as a contributor!
The Keptn Lifecycle Toolkit is a community-driven project and we
greatly value collaboration.
There are various ways to contribute to the Lifecycle Toolkit, and
all contributions are highly valued.
Please, explore the options below to learn more about how you can
contribute.

## Before you get started

### Code of Conduct

Please make sure to read and observe our
[Code of Conduct](https://github.com/keptn/.github/blob/main/CODE_OF_CONDUCT.md).

* **Create an issue**: If you have noticed a bug, want to contribute features,
or simply ask a question that for whatever reason you do not want to ask in the
[Lifecycle Toolkit Channels in the CNCF Slack workspace](https://cloud-native.slack.com/channels/keptn-lifecycle-toolkit-dev),
please [search the issue tracker](https://github.com/keptn/lifecycle-toolkit/issues?q=something)
to see if someone else in the community has already created a ticket.
If not, go ahead and [create an issue](https://github.com/keptn/lifecycle-toolkit/issues/new).

* **Start contributing**: We also have a list of
[good first issues](https://github.com/keptn/lifecycle-toolkit/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22).
If you want to work on it, just post a comment on the issue.

* **Add yourself**: Add yourself to the [list of contributors](CONTRIBUTORS.md)
along with your first pull request.

This document lays out how to get you started in contributing to
Keptn Lifecycle Toolkit, so please read on.

## How to Start?

If you are worried or don’t know where to start, check out our next section
explaining what kind of help we could use and where you can get involved.
You can reach out with the questions to
[Lifecycle Toolkit Channels](https://cloud-native.slack.com/channels/keptn-lifecycle-toolkit-dev)
on Slack and a mentor will surely guide you!

### Prerequisites

* [**Docker**](https://docs.docker.com/get-docker/): a tool for containerization,
which allows software applications to run in isolated environments
and makes it easier to deploy and manage them.
* A Kubernetes `cluster >= Kubernetes 1.24` .If you don’t have one,
we recommend Kubernetes-in-Docker(kind) to set up your local development environment.
* [**kubectl**](https://kubernetes.io/docs/tasks/tools/): a command-line interface tool used for deploying
and managing applications on Kubernetes clusters.
* [**kustomize**](https://kustomize.io/): a tool used for customizing Kubernetes resource configurations
and generating manifests.
* [**Helm**](https://helm.sh/): a package manager for Kubernetes that
simplifies the deployment and management of applications on a Kubernetes cluster.

## Related Technologies

Please check [Related Technologies](docs/content/en/contribute/general/technologies/_index.md).

## Linter Requirements

Please check [Linter Requirements](docs/content/en/contribute/docs/linter-requirements/_index.md).

## Submit a Pull Request 🚀

At this point, you should switch back to the `main` branch in your repository,
and make sure it is up to date with `main` branch of Keptn Lifecycle Toolkit:

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

> Note:
All PRs must include a commit message with a description of the changes made!

Make sure you **sign off your commits**.
To do this automatically check [this](https://github.com/keptn/lifecycle-toolkit/blob/main/CONTRIBUTING.md#auto-signoff-commit-messages).
Finally, go to GitHub and create a Pull Request.
There should be a PR template already prepared for you.
If not, you will find it at `.github/pull_request_template.md`.
Please describe what this PR is about and add a link to relevant GitHub issues.
If you changed something that is visible to the user, please add a screenshot.
Please follow the
[conventional commit guidelines](https://www.conventionalcommits.org/en/v1.0.0/) for your PR title.

If you only have one commit in your PR, please follow the guidelines for the message
of that single commit, otherwise the PR title is enough.
You can find a list of all possible feature types [here](#commit-types).

An example for a pull request title would be:

```bash
feat(api): New endpoint for feature X (#1234)
```

If you have **breaking changes** in your PR, it is important to note them in the PR
description but also in the merge commit for that PR.

When pressing "squash and merge", you have the option to fill out the commit message.
Please use that feature to add the breaking changes according to the
[conventional commit guidelines](https://www.conventionalcommits.org/en/v1.0.0/).
Also, please remove the PR number at the end and just add the issue number.

An example for a PR with breaking changes and the according merge commit:

```bash
feat(bridge): New button that breaks other things (#345)

BREAKING CHANGE: The new button added with #345 introduces new functionality that is not compatible with the previous type of sent events.
```

If your breaking change can be explained in a single line you can also use this form:

```bash
feat(bridge)!: New button that breaks other things (#345)
```

Following those guidelines helps us create automated releases where the commit
and PR messages are directly used in the changelog.

In addition, please always ask yourself the following questions:

**Based on the linked issue,**
**what changes within the PR would you expect as a reviewer?**

Your PR will usually be reviewed by the Keptn Lifecycle Toolkit team within a
couple of days, but feel free to let us know about your PR
[via Slack](https://cloud-native.slack.com/channels/keptn-lifecycle-toolkit-dev).

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

## Developer Certification of Origin (DCO)

All commits must be accompanied by a DCO sign-off.
See
[DCO](docs/content/en/contribute/general/dco)
for more information.
