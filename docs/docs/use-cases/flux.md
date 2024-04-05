---
comments: true
---

# Deploying Applications with Flux and Keptn

[Flux](https://fluxcd.io/)
is a tool that keeps Kubernetes clusters synchronized with sources
of configuration (such as Git repositories) and automates updates to
the configuration when new code is deployed.

This section shows how to add Keptn
[pre- and post-deployment tasks](../guides/tasks.md)
to an existing Flux use case that runs
[pre and post-deployment jobs with Flux](https://fluxcd.io/flux/use-cases/running-jobs/).
Adding Keptn makes it simpler and more straight-forward
to run the Flux pre and post-deployment jobs
and provides added observability out of the box.

## High-level structure of the Git repository

Flux uses a GitOps approach for continuous delivery.
The Git
repository structure for our use case looks like the following:

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
The `apps` and `clusters` directories can live in two separate repositories,
but for simplicity of this excercise, we will keep them in a single one.

You see that the Keptn runs pre/post-deployment tasks
rather than using the Flux `pre-deploy` and `post-deploy` directories in the Git repository.
The Keptn process is easier to implement
and contains more information than the Flux jobs do.

## Set up your environment

Before we start, you need to install Flux CLI on your local machine and Keptn on your cluster.
You can find the installation instructions for Keptn [here](../installation/index.md)
and for Flux [here](https://fluxcd.io/flux/installation/).

After successfully installing Flux CLI and Keptn, you need to
retrieve your Git repository credentials.
You can use any available Git providers, but be sure to store your `token`
for later usage.
For simplicity, we will use GitHub.

You then need to
[bootstrap the Git repository](https://fluxcd.io/flux/installation/bootstrap/),
supplying the `token` you saved,
to install Flux in your cluster.
This sets up all necessary Flux structures in the repository.
Use the following command to do this:

```shell
GITHUB_USER=<user-name> GITHUB_TOKEN=<token> \
flux bootstrap github \
  --owner=$GITHUB_USER \
  --repository=podtato-head \
  --branch=main \
  --path=./clusters/dev \
  --personal
```

The bootstrap command above does the following:

* Creates a Git repository `podtato-head` on your GitHub account.
* Adds Flux component manifests to the repository and
creates `./clusters/dev/flux-system/*` structure.
* Deploys Flux components to your Kubernetes Cluster.
* Configures Flux components to track the path `./clusters/dev/` in the repository.

## Creating the application

Now it's time to add the Keptn application that defines the Keptn pre- and
post-deployment checks to the repository.

First, clone the `podtato-head` repository to your local machine.

Add the following manifests that represent the application into the  `podtato-head.yaml` file
and store it in the `./apps/dev/` directory of your repository:

```yaml
{% include "./assets/flux/app.yaml" %}
```

See the Keptn
[CRD reference](../reference/crd-reference/index.md)
documentation for more information about the Keptn resources that are used.
See
[Basic annotations](../guides/integrate.md#basic-annotations)
for a description of the labels and annotations that Keptn uses.
Additionally, create a `kustomize.yaml` file right next to it:

```yaml
{% include "./assets/flux/kustomize.yaml" %}
```

You can commit and push these manifests to the Git repository.

The application has Keptn pre- and post-deployment tasks defined
in the manifests.
This enhances the
[Flux pre and post-deployment jobs](https://fluxcd.io/flux/use-cases/running-jobs/)
setup with added observability out of the box.

## Set up continuous delivery for the application

First, we need to create a Flux
[GitRepository](https://fluxcd.io/flux/components/source/gitrepositories/)
resource
pointing to a repository where the application manifests are present.
In this case, it will be our `podtato-head` repository.
To create this resource, we use Flux CLI:

```shell
flux create source git podtato-head \
  --url=<git-repo-url> \
  --branch=main \
  --interval=1m \
  --export > ./clusters/dev/podtato-head-source.yaml
```

This results in output similar to:

```yaml
{% include "./assets/flux/gitrepository.yaml" %}
```

In the last step, create a Flux
[Kustomization](https://fluxcd.io/flux/components/kustomize/kustomizations/)
resource to deploy the `podtato-head` application to the cluster.
You can create it using the Flux CLI:

```shell
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

which results in output similar to:

```yaml
{% include "./assets/flux/flux-kustomization.yaml" %}
```

Now, commit and push the resources you created in the previous steps.
After pushing them, Flux picks up the configuration and
deploys your application into the cluster.

## Watch Flux sync of the application

You can watch the synchronization of the application
using the Flux CLI:

```shell
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

Each time you update the application, the changes are
synced to the cluster.

## Possible follow-ups

You can also set up multi-stage delivery with Flux,
following the steps described for `ArgoCD` in the
[ArgoCD multi-stage delivery with Keptn](../guides/multi-stage-application-delivery.md)
user guide.
