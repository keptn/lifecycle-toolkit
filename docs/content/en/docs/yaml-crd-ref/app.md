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
  revision: x
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
  * **revision** -- revision of a `version`.
    The value is an integer that can be modified
    to trigger another deployment of a `KeptnApp` of the same version.
    For example, increment this number to restart a `KeptnApp` version
    that failed to deploy, perhaps because a
    `preDeploymentEvaluation` or `preDeploymentTask` failed.
  * **workloads**
    * **name** - name of this Kubernetes
      [workload](https://kubernetes.io/docs/concepts/workloads/).
      Use the same naming rules listed above for the application name.
      Provide one entry for each workload
      associated with this Keptn application.
    * **version** -- version number for this workload.
      Changing this number causes a new execution
      of checks for this workload only.
  * **preDeploymentTasks** -- list each task to be run
    as part of the pre-deployment stage.
    Task names must match the value of the `name` field
    for the associated [KeptnTaskDefinition](taskdefinition.md) resource.
  * **postDeploymentTasks** -- list each task to be run
    as part of the post-deployment stage.
    Task names must match the value of the `name` field
    for the associated [KeptnTaskDefinition](taskdefinition.md) resource.
  * **preDeploymentEvaluations** -- list each evaluation to be run
    as part of the pre-deployment stage.
    Evaluation names must match the value of the `name` field
    for the associated [KeptnEvaluationDefinition](evaluationdefinition.md)
    resource.
  * **postDeploymentEvaluations** -- list each evaluation to be run
    as part of the post-deployment stage.
    Evaluation names must match the value of the `name` field
    for the associated [KeptnEvaluationDefinition](evaluationdefinition.md)
    resource.

## Usage

Kubernetes defines
[workloads](https://kubernetes.io/docs/concepts/workloads/)
but does not define applications.
The Keptn Lifecycle Toolkit adds the concept of applications
defined as a set of workloads that can be executed.
A `KeptnApp` resource is inserted
into the repository of the deployment engine
(ArgoCD, Flux, etc)
and is then deployed by that deployment engine.

You can create a `KeptnApp` resource by populating the YAML file
or you can use the automatic application discovery feature
to automatically generate a `KeptnApp` resource
based on Kubernetes labels.
This allows you to use the KLT observability features for existing resources
without manually populating any Keptn related resources.

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

The `spec.Revision` field is introduced in v1alpha2.

## See also

* Link to information about app discovery in implementing/integrate
