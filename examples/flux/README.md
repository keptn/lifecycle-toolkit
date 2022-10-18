# Use lifecycle-controller together with flux

<<<<<<< HEAD
This tutorial shows how to use the lifecycle-controller together with [flux](https://fluxcd.io/).

## TL;DR
* Set up a Personal Access Token according to: https://fluxcd.io/flux/installation/#github-and-github-enterprise
* Set up flux and the Repository: `make install GITHUB_REPO=https://github.com/<YOUR_GITHUB_HANDLE>/flux-demo GITHUB_USER=<YOUR_GITHUB_HANDLE> GITHUB_TOKEN=<YOUR_GITHUB_TOKEN>`
=======
This tutorial shows how to use the lifecycle-controller together with [Flux](https://fluxcd.io/).

## TL;DR
* Set up a Personal Access Token according to: https://fluxcd.io/flux/installation/#github-and-github-enterprise
* Set up Flux and a GitHub repository: `make install GITHUB_REPO=https://github.com/<YOUR_GITHUB_HANDLE>/flux-demo GITHUB_USER=<YOUR_GITHUB_HANDLE> GITHUB_TOKEN=<YOUR_GITHUB_TOKEN>`
>>>>>>> refs/remotes/origin/feat/flux_example
* Apply manifests to the Repository: `make manifests`
* Watch the progress of the deployment using: `kubectl get keptnapplicationversions -n podtato-kubectl`
  * This might take a while

## Prerequisites
<<<<<<< HEAD
The flux cli should be installed. See [here](https://fluxcd.io/docs/installation/) for more information.
=======
The Flux CLI should be installed. See [here](https://fluxcd.io/docs/installation/) for more information.
>>>>>>> refs/remotes/origin/feat/flux_example

MacOS: `brew install fluxcd/tap/flux`
bash: `curl -s https://fluxcd.io/install.sh | sudo bash`
chocolatey: `choco install flux`

## Bootstrap your repository and install flux
Follow the instructions in the quickstart guide: https://fluxcd.io/docs/get-started/

## Installing the Demo Application
To install the demo application, you can check in the configuration provided in the config-repository to the repository you created in the previous step.

You can watch the progress of the deployment using:
> `kubectl get pods -n podtato-kubectl`
* See that the pods are pending until the pre-deployment tasks have passed
* Pre-Deployment Tasks are started
* Pods get scheduled

> `kubectl get keptnworkloadinstances -n podtato-kubectl`
* Get the current status of the workloads
* See in which phase your workload deployments are at the moment

> `kubectl get keptnapplicationversions -n podtato-kubectl`
* Get the current status of the application
* See in which phase your application deployment is at the moment

After some time all resources should be in a succeeded state. Taking a look on the kustomization resource, you can see that the deployment has been updated to the latest version.
<<<<<<< HEAD
> `k describe kustomizations.kustomize.toolkit.fluxcd.io podtatohead -n default`
=======
> `kubectl describe kustomizations.kustomize.toolkit.fluxcd.io podtatohead -n default`
>>>>>>> refs/remotes/origin/feat/flux_example
