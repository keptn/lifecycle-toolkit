---
comments: true
---

# Codespaces

Use GitHub codespaces as a pre-built and pre-configured development environment.
This is especially useful for Windows users
who may not have `make` installed.
It is also useful for Linux and MacOS users
who may not wish to download tools just to contribute to docs.

Review [this video](https://www.youtube.com/watch?v=HdiXPgvfgQw) to see how this works.

[![Keptn + GitHub codespaces video](https://i.ytimg.com/vi/HdiXPgvfgQw/hqdefault.jpg)](https://www.youtube.com/watch?v=HdiXPgvfgQw)

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

The only additional tool required for software development is [Kustomize](https://kustomize.io/),
which is a Kubernetes configuration transformation tool.
To install it, simply change your directory to the workspace
as shown in the video and copy and paste the below-mentioned two commands.
The first command downloads the precompiled binary of Kustomize;
the second command grants the necessary permissions to the tool and sets the path.
You can always download the latest binaries from the official [Kustomize GitHub repository](https://github.com/kubernetes-sigs/kustomize/releases).

```bash

curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash

sudo install -o root -g root -m 0755 kustomize /usr/local/bin/kustomize
```

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
   go to your GitHub page, select "Codespaces", and delete this codespace.
