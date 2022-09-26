# scheduler
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started
You’ll need a Kubernetes cluster v0.24.0 or higher to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Build and push your image to the location specified by `RELEASE_REGISTRY`:
	
```sh
make push-release-images RELEASE_REGISTRY=<some-registry>
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

2. Deploy the scheduler using helm, make sure it is in the same namespace of the operator:

```sh
cd manifest/install/charts
helm upgrade --install keptn-scheduler keptn-scheduler/ --set scheduler.image=<your-generated-image>  -n operator-system
```

### Uninstall 
To delete the scheduler:

```sh
helm uninstall keptn-scheduler -n operator-system
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project uses the Kubernetes [Scheduler Framework](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/)
and is based on the [Scheduler Plugins Repository](https://github.com/kubernetes-sigs/scheduler-plugins/tree/master).

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

