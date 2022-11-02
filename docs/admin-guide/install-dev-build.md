## Installing a developer build

The [GitHub CLI](https://cli.github.com/) can be used to download the manifests of the latest CI build.

```bash
gh run list --repo keptn/lifecycle-controller # find the id of a run
gh run download 3152895000 --repo keptn/lifecycle-controller # download the artifacts
kubectl apply -f ./keptn-lifecycle-operator-manifest/release.yaml # install the operator
kubectl apply -f ./scheduler-manifest/release.yaml # install the scheduler
```

Instead, if you want to build and deploy the operator into your cluster directly from the code, you can type:

```bash
RELEASE_REGISTRY=<YOUR_DOCKER_REGISTRY>
# (optional)ARCH=<amd64(default)|arm64v8>
# (optional)TAG=<YOUR_PREFERRED_TAG (defaulting to current time)>

# Build and deploy the dev images to the current kubernetes cluster
make build-deploy-dev-environment

```