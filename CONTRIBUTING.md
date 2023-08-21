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

## Working with git

See [Working with Git](docs/content/en/contribute/general/git)

Your PR will usually be reviewed by the Keptn Lifecycle Toolkit team within a
couple of days, but feel free to let us know about your PR
[via Slack](https://cloud-native.slack.com/channels/keptn-lifecycle-toolkit-dev).

## Auto signoff commit messages

We have a DCO check that runs on every PR to verify that the commit has been signed off.

To sign off the commits use `-s` flag, you can can use

```bash
git commit -s -m "my awesome contribution"
```

To sign off the last commit you made, you can use

```bash
git commit --amend --signoff
```

or the command below to sign off the last 2 commits you made

```bash
git rebase HEAD~2 --signoff
```

This process is sometimes inconvenient but you can automate it
by creating a pre-commit git hook as follows:

1. Create the hook:

    ``` bash
    touch .git/hooks/prepare-commit-msg
    ```

2. Add the following to the `prepare-commit-msg` file:

    ```bash
    SOB=$(git var GIT_AUTHOR_IDENT | sed -n 's/^\(.*>\).*$/Signed-off-by: \1/p')
    grep -qs "^$SOB" "$1" || echo "$SOB" >> "$1"
    ```

3. Give it execution permissions by calling:

    ```bash
    chmod +x ./.git/hooks/prepare-commit-msg
    ```
