---
comments: true
---

# Deploying Applications with Flux

Flux is a tool for keeping Kubernetes clusters in sync with sources
of configuration (like Git repositories), and automating updates to
configuration when there is new code to deploy.

This section shows an already existing use case of running
[pre and post-deployment jobs with Flux](https://fluxcd.io/flux/use-cases/running-jobs/)
and how Keptn makes it simpler and and more straight-forward.

## High-level structure of the git repository

Since Flux uses a GitOps approach to continuous delivery, the git
repository structure needs to look like the following:

```markdown
├── apps
│   └── dev
│       ├── podtato-head.yaml
│       └── kustomize.yaml
└── clusters
    └── dev
        ├── flux-system
        │   ├── gotk-components.yaml
        │   ├── gotk-sync.yaml
        │   └── kustomize.yaml
        ├── podtato-head-source.yaml
        └── podtato-head-kustomization.yaml
```

The `apps` directory contains application manifests, that will be deployed.
The `clusters` directory contains Flux configuration manifests and custom
resources, that are needed for the delivery.
`apps` and `clusters` directories can live in two separate repositories,
but for simplicity of this excercise, we will keep them in a single one.

## Set up environment

Before starting, you need to install Flux CLI and Keptn.
You can find the installation instructions of Keptn [here](./../installation/index.md)
and for Flux [here](https://fluxcd.io/flux/installation/).

After successfully installing Flux CLI and Keptn, you need to
retrieve your git repository credentials.
You can use any available git providers, but be sure to store your `token`
for later usage.
For simplicity, we will use `GitHub`.

In the end, you need to install Flux to your cluster.
This step will require
[bootstrapping the git repository](https://fluxcd.io/flux/installation/bootstrap/)
in order to set up all Flux structures.
For that, you can use the following command:

```bash
GITHUB_USER=<user-name> GITHUB_TOKEN=<token> \
flux bootstrap github \
  --owner=$GITHUB_USER \
  --repository=podtato-head \
  --branch=main \
  --path=./clusters/dev \
  --personal
```

The bootstrap command above does the following:

* Creates a git repository `podtato-head` on your `GitHub` account.
* Adds Flux component manifests to the repository -
creates `./clusters/dev/flux-system/*` structure.
* Deploys Flux Components to your Kubernetes Cluster.
* Configures Flux components to track the path `./clusters/dev/` in the repository.

## Creating application

Now it's time to add the application together with pre- and
post-deployments checks to the repository.

Firstly, clone the `podtato-head` repository to your local machine.

Add the following manifests into `podtato-head.yaml` representing the application
and store it into `./apps/dev/` directory of your repository:

```yaml
{% include "./assets/flux/app.yaml" %}
```

Additionally, create a `kustomize.yaml` file right next to it:

```yaml
{% include "./assets/flux/kustomize.yaml" %}
```

You can commit and push these manifests to your git repository.

> **Note**
Notice, that the application has pre- and post-deployment tasks defined
in the manifests.
This enhances the
[Flux pre and post-deployment jobs](https://fluxcd.io/flux/use-cases/running-jobs/)
setup with added observability out of the box.

## Set up continuous delivery for the application

Firstly, we need to create a
[GitRepository](https://fluxcd.io/flux/components/source/gitrepositories/)
Flux custom resource
pointing to a repository where the application manifests are present.
In this case, it will be our `podtato-head` repository.
To create it, we will use Flux CLI:

```bash
flux create source git podtato-head \
  --url=<git-repo-url> \
  --branch=main \
  --interval=1m \
  --export > ./clusters/dev/podtato-head-source.yaml
```

which will result output similar to:

```yaml
{% include "./assets/flux/gitrepository.yaml" %}
```

In the last step, create a
[Kustomization](https://fluxcd.io/flux/components/kustomize/kustomizations/)
Flux custom resource to deploy the `podtato-head` application to the cluster.
You can create it using the Flux CLI:

```bash
flux create kustomization podtato-head \
  --target-namespace=podtato-kubectl \
  --source=podtato-head \
  --path="./apps/dev" \
  --prune=true \
  --wait=true \
  --interval=30m \
  --retry-interval=2m \
  --health-check-timeout=3m \
  --export > ./clusters/dev/podtato-head-kustomization.yaml
```

which will result output similar to:

```yaml
{% include "./assets/flux/flux-kustomization.yaml" %}
```

Now commit and push the resources you created in the recent steps.
After pushing them, Flux should pick up the configuration and
deploy your application into the cluster.

## Watch Flux sync of the application

You can watch the synchronization of the application
with Flux CLI

```bash
$ flux get kustomizations --watch

NAME          REVISION             SUSPENDED  READY   MESSAGE
flux-system   main@sha1:4e9c917f   False      True    Applied revision: main@sha1:4e9c917f
podtato-head  main@sha1:44154333   False      True    Applied revision: main@sha1:44154333

```

or using `kubectl`:

```shell
$ kubectl get keptnappversion -n podtato-head

NAME                           APPNAME        VERSION   PHASE
podtato-head-v0.1.0-6b86b273   podtato-head   v0.1.0    Completed
```

Every time you update the application, the changes will be
synced to the cluster.

## Possible follow-ups

You can set up a multi-stage delivery with Flux, same
as it was done with `ArgoCD`.
You can follow similar steps to the
[ArgoCD multi-stage delivery with Keptn](../guides/multi-stage-application-delivery.md)
user guide.
