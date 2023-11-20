---
title: KeptnTaskDefinition
description: Define tasks that can be run pre- or post-deployment
weight: 89
---

A `KeptnTaskDefinition` defines tasks
that Keptn runs as part of the pre- and post-deployment phases of a
[KeptnApp](./app.md) or
[KeptnWorkload](../crd-ref/lifecycle/v1alpha3/#keptnworkload).

A Keptn task executes as a
[runner](https://docs.gitlab.com/runner/executors/kubernetes.html#how-the-runner-creates-kubernetes-pods)
in an application
[container](https://kubernetes.io/docs/concepts/containers/),
which runs as part of a Kubernetes
[job](https://kubernetes.io/docs/concepts/workloads/controllers/job/).

Each `KeptnTaskDefinition` can use exactly one container with one runner.
which is one of the following,
differentiated by the `spec` section:

* The `custom-runtime` runner provides
  a standard Kubernetes application container
  that is run as part of a Kubernetes job.
  You define the runner, an application,
  and its runtime dependencies.
  This gives you the flexibility
  to define tasks using the language and facilities of your choice,
  although it is more complicated that using one of the pre-defined runtimes.
  See
  [Synopsis for container-runtime](#synopsis-for-container-runtime)
  and
  [Examples for a container-runtime runner](#examples-for-a-container-runtime-runner).
* Pre-defined containers

  * Use the pre-defined `deno-runtime` runner
    to define tasks using
    [Deno](https://deno.com/)
    scripts,
    which use a syntax similar to JavaScript and Typescript,
    with a few limitations.
    You can use this to specify simple actions
    without having to define a full container.
    See
    [Synopsis for Deno-runtime container](#deno-runtime)
    and
    [Deno-runtime examples](#examples-for-deno-runtime-runner).
  * Use the pre-defined `python-runtime` runner
    to define your task using
    [Python 3](https://www.python.org/).
    See
    [Synopsis for python-runtime runner](#python-runtime)
    and
    [Examples for a python-runtime runner](#examples-for-a-python-runtime-runner).

## Synopsis for all runners

The `KeptnTaskDefinition` Yaml files for all runners
include the same lines at the top.
These are described here.

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnTaskDefinition
metadata:
  name: <task-name>
spec:
  deno | python | container
  ...
  retries: <integer>
  timeout: <duration>
```

### Fields used for all containers

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.
   Must be set to `KeptnTaskDefinition`

* **metadata**
  * **name** -- Unique name of this task or container.
    This is the name used to insert this task or container
    into the `preDeployment` or `postDeployment` list.
    Names must comply with the
    [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
    specification.
* **spec**
  * **deno | python | container** (required) -- Define the container type
  to use for this task.
  Each task can use one type of runner,
  identified by this field:
    * **deno** -- Use a `deno-runtime` runner
    and code the functionality in Deno script,
    which is similar to JavaScript and Typescript.
    See
    [Synopsis for deno-runtime container](#deno-runtime)
    * **python** -- Use a `python-runtime` function
    and code the functionality in Python 3.
    See
    [Synopsis for python-runtime runner](#python-runtime)
    * **container** -- Use the runner defined
      for the `container-runtime` container.
      This is a standard Kubernetes container
      for which you define the image, runner, runtime parameters, etc.
      and code the functionality to match the container you define.
      See
      [Synopsis for container-runtime container](#synopsis-for-container-runtime).
  * **retries** -- specifies the number of times
    a job executing the `KeptnTaskDefinition`
    should be restarted if an attempt is unsuccessful.
  * **timeout** -- specifies the maximum time
    to wait for the task to be completed successfully.
    The value supplied should specify the unit of measurement;
    for example, `5s` indicates 5 seconds and `5m` indicates 5 minutes.
    If the task does not complete successfully within this time frame,
    it is considered to be failed.

## Synopsis for container-runtime

Use the `container-runtime` to specify your own
[Kubernetes container](https://kubernetes.io/docs/concepts/containers/)
and define the task you want to execute.

Task sequences that are not part of the lifecycle workflow and
should be handled by the pipeline engine tools being used
such as Jenkins, Argo Workflows, Flux, and Tekton.

If you are migrating from Keptn v1,
you can use a `container-runtime` to execute
almost anything that you implemented with JES for Keptn v1.

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnTaskDefinition
metadata:
  name: <task-name>
spec:
  container:
    name: <container-name>
    image: <image-name>
    <other fields>
```

### Spec used only for container-runtime

* **spec**
  * **container** -- Container definition.
    * **name** -- Name of the container that will run,
      which is not the same as the `metadata.name` field
      that is used in the `KeptnApp` resource.
    * **image** -- name of the image you defined according to
      [image reference](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#image)
      and
      [image concepts](https://kubernetes.io/docs/concepts/containers/images/)
      and pushed to a registry
    * **other fields** -- The full list of valid fields is available at
      [ContainerSpec](../crd-ref/lifecycle/v1alpha3/#containerspec),
      with additional information in the Kubernetes
      [Container](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#Container)
      spec documentation.

## Synopsis for pre-defined containers

The pre-defined containers allow you to easily define a task
using either Deno or Python syntax.
You do not need to specify the image, volumes, and so forth.
Instead, just provide either a Deno or Python script
and Keptn sets up the container and runs the script as part of the task.

### deno-runtime

When using the `deno-runtime` runner to define a task,
the executables are coded in
[Deno-script](https://deno.com/manual),
(which is mostly the same as JavaScript and TypeScript)
and executed in the
`deno-runtime` runner,
which is a lightweight runtime environment
that executes in your namespace.
Note that Deno has tighter restrictions
for permissions and importing data
so a script that works properly elsewhere
may not function out of the box when run in the `deno-runtime` runner.

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnTaskDefinition
metadata:
  name: <task-name>
spec:
  deno:
    inline | httpRef | functionRef | ConfigMapRef
    parameters:
      map:
        textMessage: "This is my configuration"
    secureParameters:
      secret: <secret-name>
```

### python-runtime

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnTaskDefinition
metadata:
  name: <task-name>
spec:
    python:
      inline | httpRef | functionRef | ConfigMapRef
      parameters:
        map:
          textMessage: "This is my configuration"
      secureParameters:
        secret: <secret-name>
```

### Fields for pre-defined containers

* **spec** -- choose either `deno` or `python`
  * **deno** -- Specify that the task uses the `deno-runtime`
    and is expressed as a [Deno](https://deno.com/) script.
    Refer to [deno runtime](https://github.com/keptn/lifecycle-toolkit/tree/main/runtimes/deno-runtime)
    for more information about this runner.
  * **python** -- Identifies this as a Python runner.

* **inline | httpRef | functionRef | ConfigMapRef** -- choose the syntax
  used to call the executables.
  Only one of these can be specified per `KeptnTaskDefinition` resource:

  * **inline** - Include the actual executable code to execute.
    You can code a sequence of executables here
    that need to be run in order
    as long as they are executables that are part of the lifecycle workflow.
    Task sequences that are not part of the lifecycle workflow
    should be handled by the pipeline engine tools being used
    such as Jenkins, Argo Workflows, Flux, and Tekton.

    * **deno example:**
        [Example 1: inline script for a Deno script](#example-1-inline-script-for-a-deno-script)

    * **python example:**
        [Example 1: inline code for a python-runtime runner](#example-1-inline-code-for-a-python-runtime-runner)

  * **httpRef** - Specify a script to be executed at runtime
      from the remote webserver that is specified.

      This syntax allows you to call a general function
      that is used in multiple places,
      possibly with different parameters
      that are provided in the calling `KeptnTaskDefinition` resource.
      Another `KeptnTaskDefinition` resource could call this same script
      but with different parameters.

      Only one script can be executed.
      Any other scripts listed here are silently ignored.

    * **deno example:**
        [Example 2: httpRef script for a Deno script](#example-2-httpref-script-for-a-deno-script)
    * **python example:**
        [Example 2: httpRef for a python-runtime runner](#example-2-httpref-for-a-python-runtime-runner)

  * **functionRef** -- Execute another `KeptnTaskDefinition` resources.
      Populate this field with the value(s) of the `metadata.name` field
      for each `KeptnDefinitionTask` to be called.

      Like the `httpRef` syntax,this is commonly used
      to call a general function that is used in multiple places,
      possibly with different parameters
      that are set in the calling `KeptnTaskDefinition` resource.

      You must annotate the `KeptnApp` resource to run the
      calling `KeptnTaskDefinition` resource.

      The `KeptnTaskDefinition` called with `functionref`
      is the `parent task` whose runner is used for the execution
      even if it is not the same runner defined in the
      calling `KeptnTaskDefinition`.

      Only one `KeptnTaskDefinition` resources can be listed
      with the `functionRef` syntax
      although that `KeptnTaskDefinition` can call multipe
      executables (programs, functions, and scripts)..
      Any calls to additional `KeptnTaskDefinition` resources
      are silently ignored.

    * **deno example:**
        [Example 3: functionRef for a Deno script](#example-3-functionref-for-a-deno-script)
    * **python example:**
        [Example 3: functionRef for a python-runtime runner](#example-3-functionref-for-a-python-runtime-runner)

  * **ConfigMapRef** - Specify the name of a
      [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/)
      resource that contains the function to be executed.

    * **deno example:**
        [Example 5: ConfigMap for a Deno script](#example-5-configmap-for-a-deno-script)
    * **python example:**
        [Example 4: ConfigMapRef for a python-runtime runner](#example-4-configmapref-for-a-python-runtime-runner)

  * **parameters** - An optional field to supply input parameters to a function.
    Keptn passes the values defined inside the `map` field
    as a JSON object.
    See
    [Passing secrets, environment variables, and modifying the python command](#passing-secrets-environment-variables-and-modifying-the-python-command)
    and
    [Parameterized functions](../implementing/tasks.md#parameterized-functions)
    for more information.

    * **deno example:**
      [Example 3: functionRef for a Deno script](#example-3-functionref-for-a-deno-script)
    * **python example:**
      [Example 3: functionRef for a python-runner runner](#example-3-functionref-for-a-python-runtime-runner)

  * **secureParameters** -- An optional field used to pass a Kubernetes secret.
    The `secret` value is the Kubernetes secret name
    that is mounted into the runtime and made available to functions
    using the `SECURE_DATA` environment variable.

    Note that, currently, only one secret can be passed
    per `KeptnTaskDefinition` resource.

    See [Create secret text](../implementing/tasks.md#create-secret-text)
    for details.

    * **deno example:**
      [Example 3: functionRef for a Deno script](#example-3-functionref-for-a-deno-script)
    * **python example:**
      [Example 3: functionRef for a python-runner runner](#example-3-functionref-for-a-python-runtime-runner)

## Usage

A Task executes the TaskDefinition of a
[KeptnApp](app.md) or a
[KeptnWorkload](../crd-ref/lifecycle/v1alpha3/#keptnworkload).
The execution is done by spawning a Kubernetes
[Job](https://kubernetes.io/docs/concepts/workloads/controllers/job/)
to handle a single Task.
In its state, it tracks the current status of this Kubernetes Job.

When using a container runtime that includes a volume,
an `EmptyDir` volume is created
with the same name as is specified the container `volumeMount`.
Note that, if more `volumeMount`s are specified,
only one volume is created with the name of the first `volumeMount`.
By default, the size of this volume is 1GB.
If the memory limit for the container is set,
the size of the volume is 50% of the memory allocated for the node.

A task can be executed either pre-deployment or post-deployment
as specified in the pod template specs of your Workloads
([Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
[DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
and
[ReplicaSets](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/))
and in the
[KeptnApp](app.md) resource.
See
[Pre- and post-deployment tasks](../implementing/integrate.md#pre--and-post-deployment-checks)
for details.
Note that the annotation identifies the task by `name`.
This means that you can modify the `function` code in the resource definition
and the revised code is picked up without additional changes.

All `KeptnTaskDefinition` resources specified to the `KeptnApp` resource
at the same stage (either pre- or post-deployment) run in parallel.
You can run multiple executables sequentially
either by using the `inline` syntax for a pre-defined container image
or by creating your own image
and running it in the Keptn `container-runtime` runner.
See
[Executing sequential tasks](../implementing/tasks.md#executing-sequential-tasks)
for more information.

## Examples for a container-runtime runner

For an example of a `KeptnTaskDefinition` that defines a custom container.
 see
[container-task.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/base/container-task.yaml).
This is a trivial example that just runs `busybox`,
then spawns a shell and runs the `sleep 30` command:

{{< embed path="/examples/sample-app/base/container-task.yaml" >}}

This task is then referenced in the
[app.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/version-3/app.yaml)
file.

## Examples for deno-runtime runner

### Example 1: inline script for a Deno script

This example defines a full-fledged Deno script
within the `KeptnTaskDefinition` YAML file:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition
metadata:
  name: hello-keptn-inline
spec:
  deno:
    inline:
      code: |
        let text = Deno.env.get("DATA");
        let data;
        let name;
        data = JSON.parse(text);

        name = data.name
        console.log("Hello, " + name + " new");
```

### Example 2: httpRef script for a Deno script

This example fetches the Deno script from a remote webserver at runtime:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition
metadata:
  name: hello-keptn-http
spec:
  deno:
    httpRef:
      url: "https://www.example.com/yourscript.js"
```

For another example, see the
[sample-app](https://github.com/keptn-sandbox/lifecycle-toolkit-examples/blob/main/sample-app/version-1/app-pre-deploy.yaml).

See the
[sample-app/version-1](https://github.com/keptn-sandbox/lifecycle-toolkit-examples/blob/main/sample-app/version-1/app-pre-deploy.yaml)
PodtatoHead example for a more complete example.

### Example 3: functionRef for a Deno script

This example calls another defined task,
illustrating how one `KeptnTaskDefinition` can build
on top of other `KeptnTaskDefinition`s.
In this case, it calls `slack-notification-dev`,
passing `parameters` and `secureParameters` to that other task:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition
metadata:
  name: slack-notification-dev
spec:
  deno:
    functionRef:
      name: slack-notification
    parameters:
      map:
        textMessage: "This is my configuration"
    secureParameters:
      secret: slack-token
```

### Example 4: ConfigMapRef for a Deno script

This example references a `ConfigMap` by the name of `dev-configmap`
that contains the code for the function to be executed.

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition
metadata:
  name: keptntaskdefinition-sample
spec:
  deno:
    configMapRef:
      name: dev-configmap
```

### Example 5: ConfigMap for a Deno script

This example illustrates the use of both a `ConfigMapRef` and a `ConfigMap`:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnTaskDefinition
metadata:
  name: scheduled-deployment
spec:
  function:
    configMapRef:
      name: scheduled-deployment-cm-1
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: scheduled-deployment-1
data:
  code: |
    let text = Deno.env.get("DATA");
    let data;
    if (text != "") {
        data = JSON.parse(text);
    }
    let targetDate = new Date(data.targetDate)
    let dateTime = new Date();
    if(targetDate < dateTime) {
        console.log("Date has passed - ok");
        Deno.exit(0);
    } else {
        console.log("It's too early - failing");
        Deno.exit(1);
    }
    console.log(targetDate);
```

## Examples for a python-runtime runner

### Example 1: inline code for a python-runtime runner

You can embed python code directly in the task definition.
This example prints data stored in the parameters map:
{{< embed path="/lifecycle-operator/config/samples/python_execution/taskdefinition_pyfunction_inline.yaml" >}}

### Example 2: httpRef for a python-runtime runner

You can refer to code stored online.
For example, we have a few examples available in the
[python-runtime samples](https://github.com/keptn/lifecycle-toolkit/tree/main/runtimes/python-runtime/samples)
tree.

Consider the following:
{{< embed path="/lifecycle-operator/config/samples/python_execution/taskdefinition_pyfunction_configmap.yaml" >}}

### Example 3: functionRef for a python-runtime runner

You can refer to an existing `KeptnTaskDefinition`.
This example calls the inline example
but overrides the data printed with what is specified in the task:
{{< embed path="/lifecycle-operator/config/samples/python_execution/taskdefinition_pyfunction_recursive.yaml" >}}

### Example 4: ConfigMapRef for a python-runtime runner

{{< embed path="/lifecycle-operator/config/samples/python_execution/taskdefinition_pyfunction_configmap.yaml" >}}

### Allowed libraries for the python-runtime runner

The following example shows how to use some of the allowed packages, namely:
requests, json, git, and yaml:

{{< embed path="/lifecycle-operator/config/samples/python_execution/taskdefinition_pyfunction_inline_printargs_py.yaml">}}

### Passing secrets, environment variables and modifying the python command

The following examples show how to pass data inside the parameter map,
how to load a secret in your code,
and how to modify the python command.
In this case the container runs with the `-h` option,
which prints the help message for the python3 interpreter:

{{< embed path="/lifecycle-operator/config/samples/python_execution/taskdefinition_pyfunction_use_envvars.yaml" >}}

## More examples

See the [lifecycle-operator/config/samples](https://github.com/keptn/lifecycle-toolkit/tree/main/lifecycle-operator/config/samples/function_execution)
directory for more example `KeptnTaskDefinition` YAML files.

## Files

API Reference:

* [KeptnTaskDefinition](../crd-ref/lifecycle/v1alpha3/_index.md#keptntaskdefinition)
* [KeptnTaskDefinitionList](../crd-ref/lifecycle/v1alpha3/_index.md#keptntaskdefinitionlist)
* [KeptnTaskDefinitionSpec](../crd-ref/lifecycle/v1alpha3/_index.md#keptntaskdefinitionspec)
* [FunctionReference](../crd-ref/lifecycle/v1alpha3/_index.md#functionreference)
* [FunctionSpec](../crd-ref/lifecycle/v1alpha3/_index.md#runtimespec)
* [FunctionStatus](../crd-ref/lifecycle/v1alpha3/_index.md#functionstatus)
* [HttpReference](../crd-ref/lifecycle/v1alpha3/_index.md#httpreference)
* [Inline](../crd-ref/lifecycle/v1alpha3/_index.md#inline)

## Differences between versions

The `KeptnTaskDefinition` support for
the `container-runtime` and `python-runtime` is introduced in v0.8.0.
This modifies the synopsis in the following ways:

* Add the `spec.container` field.
* Add the `python` descriptor for the `python-runtime` runner.
* Add the `container` descriptor for the `container-runtime` runner.
* Add the `deno` descriptor to replace `function`
  for the `deno-runtime` runner.
  The `function` identifier for the `deno-runtime` runner
  is deprecated;
  it still works for v 0.8.0 but will be dropped from future releases.
* The `spec.function` field is changed to be a pointer receiver.
  This aligns it with the `spec.container` field,
  which must be a pointer,
  and enables `KeptnTask` to omit it when it is empty,
  which it must be when `spec.container` is populated.

## Limitations

* Only one
  [runtime](https://kubernetes.io/docs/setup/production-environment/container-runtimes/)
  is allowed per `KeptnTaskDefinition` resource.

* Only one secret can be passed
  per `KeptnTaskDefinition` resource.

## See also

* [KeptnApp](app.md)
* [Working with tasks](../implementing/tasks.md)
* [Pre- and post-deployment tasks](../implementing/integrate.md#pre--and-post-deployment-checks)
* [KeptnApp and KeptnWorkload resources](../architecture/keptn-apps.md).
* [Orchestrate deployment checks](../intro/usecase-orchestrate.md)
* [Executing sequential tasks](../implementing/tasks.md#executing-sequential-tasks)
