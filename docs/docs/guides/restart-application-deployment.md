---
comments: true
---

# Redeploy/Restart an Application

A [KeptnApp](../reference/crd-reference/app.md) can fail
because of an unsuccessful pre-deployment evaluation
or pre-deployment task.
For example, this happens if the target value of a
[KeptnEvaluationDefinition](../reference/crd-reference/evaluationdefinition.md)
resource is misconfigured
or a pre-deployment evaluation checks the wrong URL.

After you fix the configuration
that caused the pre-deployment evaluation or task to fail,
you can increment the `spec.revision` value
and apply the updated `KeptnApp` manifest
to create a new revision of the `KeptnApp`
without modifying the `version`.

Afterwards, all related `KeptnWorkloadVersions`
automatically refer to the newly created revision of the `KeptnAppVersion`
to determine whether they are allowed
to enter their respective deployment phases.

To illustrate this, consider the following example:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: restartable-apps
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled"
---
apiVersion: lifecycle.keptn.sh/v1beta1
kind: KeptnAppContext
metadata:
  name: podtato-head
  namespace: restartable-apps
spec:
  preDeploymentTasks:
    - pre-deployment-check
---
apiVersion: lifecycle.keptn.sh/v1beta1
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

Applying these resources causes a `KeptnApp`
to be created.
The `KeptnApp` contains all workloads that are part of
the application, and a version number that is
derived from the workloads:

```yaml
apiVersion: lifecycle.keptn.sh/v1beta1
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: restartable-apps
  annotations:
    app.kubernetes.io/managed-by: "keptn"
spec:
  version: "0.1.2"
  revision: 1
  workloads:
    - name: podtato-head-entry
      version: "0.1.2"
```

Then, based of the `KeptnApp` and `KeptnAppContext`,
a `KeptnAppVersion` is created automatically and
the execution of the pre-deployment checks defined in
the `KeptnAppContext` starts.
In this example, the pre deployment checks
contain a task which clearly fails.
Therefore, the `KeptnAppVersion` is not able
to proceed with the deployment.

You can inspect the status of the created `KeptnAppVersion`:

```shell
$ kubectl get keptnappversions.lifecycle.keptn.sh -n restartable-apps
NAME                   APPNAME        VERSION   PHASE
podtato-head-0.1.2-ab1223js   podtato-head   0.1.1     AppPreDeployTasks
```

Notice that the `KeptnAppVersion` stays
in the `AppPreDeployTasks` phase for a while,
due to the pre-check trying to run
until the failure threshold is reached.
Eventually, the `KeptnAppVersion`'s `PredeploymentPhase`
is in a `Failed` state, with the remaining phases being `Deprecated`.

<!-- markdownlint-disable MD013 -->
```shell
$ kubectl get keptnappversions.lifecycle.keptn.sh -n restartable-apps -owide
NAME                          APPNAME        VERSION   PHASE               PREDEPLOYMENTSTATUS   PREDEPLOYMENTEVALUATIONSTATUS   WORKLOADOVERALLSTATUS   POSTDEPLOYMENTSTATUS   POSTDEPLOYMENTEVALUATIONSTATUS
podtato-head-0.1.2-ab1223js   podtato-head   0.1.2     AppPreDeployTasks   Failed                Deprecated                      Deprecated              Deprecated             Deprecated
```
<!-- markdownlint-enable MD013 -->

To fix the deployment of this application,
we first need to fix the task that has failed earlier.
To do so, edit the `pre-deployment-check` `KeptnTaskDefinition`:

```shell
kubectl -n restartable-apps edit keptntaskdefinitions.lifecycle.keptn.sh pre-deployment-check
```

Modify the manifest to look like this:

```yaml
apiVersion: lifecycle.keptn.sh/v1beta1
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

To restart the deployment of our `KeptnApplication`,
edit the manifest:

```shell
kubectl -n restartable-apps edit keptnapps.lifecycle.keptn.sh podtato-head
```

Increment the value of the `spec.revision` field by one:

```yaml
apiVersion: lifecycle.keptn.sh/v1beta1
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: restartable-apps
spec:
  version: "0.1.2"
  revision: 2 # Increased this value from 1 to 2
  workloads:
    - name: podtato-head-entry
      version: "0.1.2"
```

After those changes have been made,
you will notice a new revision of the `podtato-head` `KeptnAppVersion`:

```shell
$ kubectl get keptnappversions.lifecycle.keptn.sh -n restartable-apps       
NAME                          APPNAME        VERSION   PHASE
podtato-head-0.1.2-ab1223js   podtato-head   0.1.2     AppPreDeployTasks
podtato-head-0.1.2-xbhj9073   podtato-head   0.1.2     AppDeploy
```

See that the newly created revision `podtato-head-0.1.2-xbhj9073`
has made it beyond the pre-deployment check phase
and has reached its `AppDeployPhase`.

You can also verify the execution of the `pre-deployment-check`
by retrieving the list of `KeptnTasks` in
the `restartable-apps` namespace:

<!-- markdownlint-disable MD013 -->
```shell
$ kubectl get keptntasks.lifecycle.keptn.sh -n restartable-apps
NAME                             APPNAME        APPVERSION   WORKLOADNAME   WORKLOADVERSION   JOB NAME                              STATUS
pre-pre-deployment-check-49827   podtato-head   0.1.2                                         klc-pre-pre-deployment-check--77601   Failed
pre-pre-deployment-check-65056   podtato-head   0.1.2                                         klc-pre-pre-deployment-check--57313   Succeeded
```
<!-- markdownlint-enable MD013 -->

Notice that the previous failed instances are still available
for both `KeptnAppVersions` and `KeptnTasks`.
This may be useful historical data to keep track of
what went wrong during earlier deployment attempts.
