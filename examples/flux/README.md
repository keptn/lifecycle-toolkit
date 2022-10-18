# Use lifecycle-controller together with flux

This tutorial shows how to use the lifecycle-controller together with [flux](https://fluxcd.io/).

## TL;DR
* Set up a Personal Access Token according to: https://fluxcd.io/flux/installation/#github-and-github-enterprise
* Set up flux and the Repository: `make install GITHUB_REPO=https://github.com/<YOUR_GITHUB_HANDLE>/flux-demo GITHUB_USER=<YOUR_GITHUB_HANDLE>`

## Prerequisites
The flux cli should be installed. See [here](https://fluxcd.io/docs/installation/) for more information.

MacOS: `brew install fluxcd/tap/flux`
bash: `curl -s https://fluxcd.io/install.sh | sudo bash`
chocolatey: `choco install flux`

## Bootstrap your repository

