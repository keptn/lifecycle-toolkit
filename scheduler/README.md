# scheduler

// TODO(user): Add simple overview of use/purpose

## Description

// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

You’ll need a Kubernetes cluster v0.24.0 or higher to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a
local cluster for testing, or run against a remote cluster.
> **Note:**
Your controller will automatically use the current context in your kubeconfig file (i.e. whatever
cluster `kubectl cluster-info` shows).

### Running on the cluster

1. Build and push your image to the location specified by `RELEASE_REGISTRY`:

```sh
make build-and-push-local RELEASE_REGISTRY=<some-registry>
```

> **Note:**
Run `make --help` for more information on all potential `make` targets

1. Generate your release manifest

```sh
make release-manifests RELEASE_REGISTRY=<some-registry>
```

1. Deploy the scheduler using kubectl:

```sh
kubectl apply -f ./config/rendered/release.yaml # install the scheduler
```

### Uninstall

To delete the scheduler:

```sh
kubectl delete -f ./config/rendered/release.yaml # uninstall the scheduler
```

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works

This project uses the
Kubernetes [Scheduler Framework](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/)
and is based on the [Scheduler Plugins Repository](https://github.com/kubernetes-sigs/scheduler-plugins/tree/master).
