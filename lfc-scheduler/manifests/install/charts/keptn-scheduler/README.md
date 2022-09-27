# Scheduler-plugins as a second scheduler in cluster

## Installation

Quick start instructions for the setup and configuration of as-a-second-scheduler using Helm.

### Prerequisites

- [Helm](https://helm.sh/docs/intro/quickstart/#install-helm)

### Installing the chart

#### Install chart using Helm v3.0+

```bash
$ git clone git@github.com::keptn-sandbox/lifecycle-controller.git
$ cd lifecycle-controller/lfc-scheduler/manifests/install/charts
$ helm install keptn-scheduler keptn-scheduler/
```

#### Verify that scheduler and plugin-controller pod are running properly.

```bash
$ kubectl get deploy -n scheduler-plugins
NAME                           READY   UP-TO-DATE   AVAILABLE   AGE
keptn-scheduler                 1/1     1            1           7s
```

### Configuration

The following table lists the configurable parameters of the scheduler chart and their default values.

| Parameter                               | Description                   | Default                                                                                         |
| --------------------------------------- |-------------------------------|-------------------------------------------------------------------------------------------------|
| `scheduler.name`                        | Scheduler name                | `scheduler-plugins-scheduler`                                                                   |
| `scheduler.image`                       | Scheduler image               | `k8s.gcr.io/scheduler-plugins/kube-scheduler:v0.23.10`                                          |
| `scheduler.namespace`                   | Scheduler namespace           | `scheduler-plugins`                                                                             |
| `scheduler.replicaCount`                | Scheduler replicaCount        | `1`                                                                                             |
| `plugins.enabled`                       | Plugins enabled by default    | `["Coscheduling","CapacityScheduling","NodeResourceTopologyMatch", "NodeResourcesAllocatable"]` |
| `plugins.enabled`                       | Plugins disabled by default   | `["PrioritySort"]`                                                                              |

