---
title: Restart an Application Deployment
description: Learn how to restart an unsuccessful Keptn Application Deployment.
icon: concepts
layout: quickstart
weight: 20
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

## Restart an Application Deployment

During the deployment of a `KeptnApp`, it might be that the deployment fails due to an unsuccessful pre-deployment
evaluation or pre-deployment task.
This could happen because of, e.g., a misconfigured target value of a `KeptnEvaluationDefinition`, or a wrong URL being
checked in a pre deployment check.

To retry a `KeptnApp` deployment without incrementing the version of the `KeptnApp`, we introduced the concept of **
revisions** for a `KeptnAppVersion`.
This means that
whenever the spec of a `KeptnApp` changes, even though the version stays the same, the KLT Operator will create a new
revision of the `KeptnAppVersion` referring to the `KeptnApp`.

This way, when a `KeptnApp` failed due to a misconfigured pre-deployment check, you can first fix the configuration of
the `KeptnTaskDefinition`/`KeptnEvaluationDefinition`, then
increase the value of `spec.revision` of the `KeptnApp` and finally apply the updated `KeptnApp` manifest.
This will result in a restart of the `KeptnApp`.
Afterwards, all related `KeptnWorkloadInstances` will automatically refer to the newly
created revision of the `KeptnAppVersion` to determine whether they are allowed to enter their respective deployment
phase.

To illustrate this, let's have a look at the following example:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: restartable-apps
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled"
---
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: restartable-apps
spec:
  version: "0.1.1"
  revision: 1
  workloads:
    - name: podtato-head-entry
      version: "0.1.2"
  preDeploymentTasks:
    - pre-deployment-check
---
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnTaskDefinition
metadata:
  name: pre-deployment-check
  namespace: restartable-apps
spec:
  function:
    inline:
      code: |
        console.error("I failed")
        process.exit(1)
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podtato-head-entry
  namespace: restartable-apps
  labels:
    app: podtato-head
spec:
  selector:
    matchLabels:
      component: podtato-head-entry
  template:
    metadata:
      labels:
        component: podtato-head-entry
        keptn.sh/workload: podtato-head-entry
        keptn.sh/app: podtato-head
        keptn.sh/version: "0.1.2"
    spec:
      terminationGracePeriodSeconds: 5
      containers:
        - name: server
          image: ghcr.io/podtato-head/entry:0.2.7
          imagePullPolicy: Always
          ports:
            - containerPort: 9000
          env:
            - name: PODTATO_PORT
              value: "9000"
```

In this example, the `KeptnApp` executes a pre-deployment check which clearly fails due to the `pre-deployment-check`
task, and will therefore not be able to proceed with the deployment.

After applying this manifest, you can inspect the status of the created `KeptnAppVersion`:

```shell
$ kubectl get keptnappversions.lifecycle.keptn.sh -n restartable-apps
NAME                   APPNAME        VERSION   PHASE
podtato-head-0.1.1-1   podtato-head   0.1.1     AppPreDeployTasks
```

You will notice that the `KeptnAppVersion` will stay in the `AppPreDeployTasks` phase for a while, due to the pre-check
trying to run until a certain failure threshold is reached.
Eventually, you will find the `KeptnAppVersion`'s `PredeploymentPhase` to be in a `Failed` state, with the remaining
phases being `Deprecated`.

<!-- markdownlint-disable MD013 -->
```shell
$ kubectl get keptnappversions.lifecycle.keptn.sh -n restartable-apps -owide
NAME                   APPNAME        VERSION   PHASE               PREDEPLOYMENTSTATUS   PREDEPLOYMENTEVALUATIONSTATUS   WORKLOADOVERALLSTATUS   POSTDEPLOYMENTSTATUS   POSTDEPLOYMENTEVALUATIONSTATUS
podtato-head-0.1.1-1   podtato-head   0.1.1     AppPreDeployTasks   Failed                Deprecated                      Deprecated              Deprecated             Deprecated
```
<!-- markdownlint-enable MD013 -->

Now, to fix the deployment of this application, we first need to fix the task that has failed earlier.
To do so, edit the `pre-deployment-check` `KeptnTaskDefinition` to the
following (`kubectl -n restartable-apps edit keptntaskdefinitions.lifecycle.keptn.sh pre-deployment-check`):

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnTaskDefinition
metadata:
  name: pre-deployment-check
  namespace: restartable-apps
spec:
  function:
    inline:
      code: |
        console.error("Success")
```

After we have done that, we can restart the deployment of our `KeptnApplication` by incrementing the `spec.revision`
field by one
(`kubectl -n restartable-apps edit keptnapps.lifecycle.keptn.sh podtato-head`):

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: restartable-apps
spec:
  version: "0.1.1"
  revision: 2 # Increased this value from 1 to 2
  workloads:
    - name: podtato-head-entry
      version: "0.1.2"
  preDeploymentTasks:
    - pre-deployment-check
```

After those changes have been made, you will notice a new revision of the `podtato-head` `KeptnAppVersion`:

```shell
$ kubectl get keptnappversions.lifecycle.keptn.sh -n restartable-apps       
NAME                   APPNAME        VERSION   PHASE
podtato-head-0.1.1-1   podtato-head   0.1.1     AppPreDeployTasks
podtato-head-0.1.1-2   podtato-head   0.1.1     AppDeploy
```

As you will see, the newly created revision `podtato-head-0.1.1-2` has made it beyond the pre-deployment check phase and
has reached its `AppDeployPhase`.

You can also verify the execution of the `pre-deployment-check` by retrieving the list of `KeptnTasks` in
the `restartable-apps` namespace:

<!-- markdownlint-disable MD013 -->
```shell
$ kubectl get keptntasks.lifecycle.keptn.sh -n restartable-apps
NAME                             APPNAME        APPVERSION   WORKLOADNAME   WORKLOADVERSION   JOB NAME                              STATUS
pre-pre-deployment-check-49827   podtato-head   0.1.1                                         klc-pre-pre-deployment-check--77601   Failed
pre-pre-deployment-check-65056   podtato-head   0.1.1                                         klc-pre-pre-deployment-check--57313   Succeeded
```
<!-- markdownlint-enable MD013 -->

You will notice that for both the `KeptnAppVersions` and `KeptnTasks` the previous failed instances are still available,
as this might be useful historical data to keep track of
what went wrong during earlier deployment attempts.
