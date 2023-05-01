---
title: KeptnApp
description: Define all workloads and checks associated with an application
weight: 10
---

`KeptnApp` specifies a concept of a running application with a list of workloads and the possibility to define pre/post deployment and pre/post evaluations checks
that should be executed at the Keptn application level.
It contains information about all workloads and checks
that are associated with a Keptn application.

## Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnApp
metadata:
  name: <app-name>
  namespace: <app-namespace>
spec:
  version: "x.y"
  workloads:
  - name: <workload-name>
    version: x.y.z
  - name: podtato-head-left-leg
    version: x.y.z
  preDeploymentTasks:
  - <list of tasks>
  postDeploymentTasks:
  - <list of tasks>
  preDeploymentEvaluations:
  - <list of evaluations>
  postDeploymentEvaluations:
  - <list of evaluations>
```

## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.
   Must be set to `KeptnApp`

* **metadata**
  * **name** -- Unique name of this application.
    Names must comply with the
    [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
    specification.

* **spec**
  * **version** -- version of the Keptn application.
    Changing this version number causes a new execution
    of all application-level checks
  * **workloads**
    * **name** - name of this Kubernetes
      [workload](https://kubernetes.io/docs/concepts/workloads/).
      Use the same naming rules listed above for the application name.
      Provide one entry for each workload
      associated with this Keptn application.
    * **version** -- version number for this workload.
      Changing this number causes a new execution
      of checks for this workload only.
  * **postDeploymentTasks** -- list each task to be run
    as part of the post-deployment stage.
    Task names must match the value of the `name` field
    for the associated [KeptnTaskDefinition](taskdefinition.md) CRD.
  * **preDeploymentEvaluations** -- list each evaluation to be run
    as part of the pre-deployment stage.
    Evaluation names must match the value of the `name` field
    for the associated [KeptnEvaluationDefinition](evaluationdefinition.md) CRD.

## Usage

Kubernetes defines
[workloads](https://kubernetes.io/docs/concepts/workloads/)
but does not define applications.
The Keptn Lifecycle Toolkit adds the concept of applications
defined as a set of workloads that can be executed.
A Keptn application is a file that is inserted
into the repository of the deployment engine
(ArgoCD, Flux, etc)
and is then deployed by that deployment engine.

## Example

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: podtato-kubectl
spec:
  version: "1.3"
  workloads:
  - name: podtato-head-left-arm
    version: 0.1.0
  - name: podtato-head-left-leg
    version: 1.2.3
  postDeploymentTasks:
  - post-deployment-hello
  preDeploymentEvaluations:
  - my-prometheus-definition
```

## Files

## Differences between versions

The `KeptnApp` CRD is the same for
all lifecycle API versions.

## See also

* Link to reference pages for any related CRDs
