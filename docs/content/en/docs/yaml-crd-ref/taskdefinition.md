---
title: KeptnTaskDefinition
description: Define tasks that can be run pre- or post-deployment
weight: 89
---

A `KeptnTaskDefinition` defines tasks
that are run by the Keptn Lifecycle Toolkit
as part of the pre- and post-deployment phases of a
[KeptnApp](./app.md) or
[KeptnWorkload](../concepts/workloads/).

A Keptn task runs as an application
[container](https://kubernetes.io/docs/concepts/containers/),
which runs as part of a Kubernetes
[job](https://kubernetes.io/docs/concepts/workloads/controllers/job/).

Each `KeptnTaskDefinition` can use exactly one container,
which is one of the following,
differentiated by the `spec` section:

* KLT includes a Deno-runtime container
  that you can use to define tasks using Deno scripts,
  which is basically JavaScript/Typescript with a few limitations.
  You can use this to specify simple actions
  without having to define a container.
  See
  [Deno-runtime synopsis](#yaml-synopsis-for-deno-runtime-container)
  and
  [Deno-runtime examples](#examples-for-deno-runtime).
* The `custom-runtime` provides custom Kubernetes application containers,
  that you define to includes a runtime,  an application
  and its runtime dependencies.
  This gives you the flexibility,
  to define tasks using the lanugage and facilities of your choice,
  although it is more complicated that using one of the pre-defined runtimes.
  See
  [Yaml synopsis for container-runtime](#yaml-synopsis-for-container-runtime)
  and
  [Custom container examples](#examples-for-a-custom-container).
* KLT also includes a `python-runtime` container
  This can be used to define your task using Python 3.
  See
  [Yaml synopsis for python-container](#yaml-synopsis-for-python-runtime-container)
  and
  [Python container examples](#examples-for-a-python-container).

## Yaml Synopsis for all containers

The `KeptnTaskDefinition` Yaml files for all container runtimes
include the same lines at the top.
These are described here.

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnTaskDefinition
metadata:
  name: <task-name>
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

## Yaml Synopsis for Deno-runtime container

When using the Deno-runtime container to define a task,
the `function` is coded in JavaScript
and executed in
[Deno](https://deno.com/runtime),
which is a lightweight runtime environment
that executes in your namespace.
Note that Deno has tighter restrictions
for permissions and importing data
so a script that works properly elsewhere
may not function out of the box when run in Deno.

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnTaskDefinition
metadata:
  name: <task-name>
spec:
  function:
    inline | httpRef | functionRef | ConfigMapRef
    parameters:
      map:
        textMessage: "This is my configuration"
    secureParameters:
      secret: slack-token
```

### Spec fields for Deno-runtime definitions

* **spec**
  * **function** -- Code to be executed,
    expressed as a [Deno](https://deno.land/) script.
    Refer to [function runtime](https://github.com/keptn/lifecycle-toolkit/tree/main/functions-runtime)
    for more information about the runtime.

    The `function` can be defined as one of the following:

    * **inline** - Include the actual executable code to execute.
      This can be written as a full-fledged Deno script
      that is included in this file.
      For example:

      ```yaml
      function:
        inline:
          code: |
            console.log("Deployment Task has been executed");
      ```

    * **httpRef** - Specify a Deno script to be executed at runtime
      from the remote webserver that is specified.
      For example:

      ```yaml
      name: hello-keptn-http
        spec:
            function:
              httpRef:
                url: "https://www.example.com/yourscript.js"
      ```

    * **functionRef** -- Execute one or more `KeptnTaskDefinition` resources
      that have been defined.
      Populate this field with the value(s) of the `name` field
      for the `KeptnTaskDefinition`(s) to be called.
      This is commonly used to call a general function
      that is used in multiple places, possibly with different parameters.
      An example is:

       ```yaml
       spec:
         function:
           functionRef:
             name: slack-notification
       ```

      This can also be used to group a set of tasks
      into a single `KeptnTaskDefinition`,
      such as defining a `KeptnTaskDefinition` for testing.
      In this case, it calls other, existing `KeptnTaskDefinition`s
      for each type of test to be run,
      specifying each by the value of the `name` field.
    * **ConfigMapRef** - Specify the name of a
      [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/)
      resource that contains the function to be executed.

  * **parameters** - An optional field to supply input parameters to a function.
    The Lifecycle Toolkit passes the values defined inside the `map` field
    as a JSON object.
    For example:

     ```yaml
       spec:
         parameters:
           map:
             textMessage: "This is my configuration"
     ```

     See
     [Parameterized functions](../implementing/tasks/#parameterized-functions)
     for more information.

  * **secureParameters** -- An optional field used to pass a Kubernetes secret.
    The `secret` value is the Kubernetes secret name
    that is mounted into the runtime and made available to functions
    using the `SECURE_DATA` environment variable.
    For example:

    ```yaml
    secureParameters:
      secret: slack-token
    ```

    Note that, currently, only one secret can be passed.

    See [Create secret text](../implementing/tasks/#create-secret-text)
    for details.

## Yaml Synopsis for container-runtime

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnTaskDefinition
metadata:
  name: <task-name>
spec:
  container
    name: <container-name>
    image: <image-name>
    <other fields>
```

### Spec used only for container-runtime

The `container-runtime` can be used to specify
your own container image and define almost task you want to do.
If you are migrating from Keptn v1,
you can use a `container-runtime` to execute
almost anything that you implemented with JES for Keptn v1.

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

## Yaml Synopsis for Python-runtime container

The `python-runtime` is built on the `container-runtime`
provides a way easily define a task using Python 3.
You do not need to specify the image, volumes, and so forth.
Instead, just provide a Python script
and KLT sets up the container and runs the script as part of the task.

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnTaskDefinition
metadata:
  name: <task-name>
spec:
    python
      inline | httpRef | functionRef | ConfigMapRef
      parameters:
        map:
          textMessage: "This is my configuration"
      secureParameters:
        secret: slack-token
```

TODO: Fix synopsis above

### Spec used only for python-runtime definitions

The `python-container` can be used to define tasks using  Python 3 code.

* **spec**
  * **python** -- Identifies this as a Python container
    * **inline** -- Include the actual Python 3.1 code to execute.
      For example, the following example
      prints data stored in the parameters map:

      ```yaml
      function:
        inline:
          code: |
            console.log("Deployment Task has been executed");
      ```

TODO: Need python code above

    * **httpRef** - Specify a Deno script to be executed at runtime
      from the remote webserver that is specified.
      For example:

      ```yaml
      name: hello-keptn-http
        spec:
            function:
              httpRef:
                url: "https://www.example.com/yourscript.py"
      ```

    * **functionRef** -- Execute one or more `KeptnTaskDefinition` resources
      that have been defined.
      Populate this field with the value(s) of the `metadata.name` field
      for each `KeptnDefinitionTask` to be called.
      This is commonly used to call a general function
      that is used in multiple places,
      possibly with different parameters.
      An example is:

      ```yaml
       spec:
         function:
           functionRef:
             name: slack-notification
       ```

      This can also be used to group a set of tasks
      into a single `KeptnTaskDefinition`,
      such as defining a `KeptnTaskDefinition` for testing.
      In this case, it calls other, existing `KeptnTaskDefinition`s
      for each type of test to be run,
      specifying each by the value of the `name` field.

    * **ConfigMapRef** -- Specify the name of a
      [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/)
      resource that contains the function to be executed.

  * **parameters** - An optional field to supply input parameters to a function.
    The Lifecycle Toolkit passes the values defined inside the `map` field
    as a JSON object.
    For example:

     ```yaml
       spec:
         parameters:
           map:
             textMessage: "This is my configuration"
     ```

     See
     [Parameterized functions](../implementing/tasks/#parameterized-functions)
     for more information.

  * **secureParameters** -- An optional field used to pass a Kubernetes secret.
    The `secret` value is the Kubernetes secret name
    that is mounted into the runtime and made available to functions
    using the `SECURE_DATA` environment variable.
    For example:

    ```yaml
    secureParameters:
      secret: slack-token
    ```

    Note that, currently, only one secret can be passed.

    See [Create secret text](../implementing/tasks/#create-secret-text)
    for details.

## Usage

A Task executes the TaskDefinition of a
[KeptnApp](app.md) or a
[KeptnWorkload](../concepts/workloads).
The execution is done by spawning a Kubernetes
[Job](https://kubernetes.io/docs/concepts/workloads/controllers/job/)
to handle a single Task.
In its state, it tracks the current status of this Kubernetes Job.

When using a container runtime that includes a volume,
an `EmptyDir` volume is created with the same name as is specified the container `volumeMount`. Please be aware, that, if more `volumeMount`s are specified, only one volume is created with the name of the first `volumeMount`.
By default, the size of this volume is 1GB.
If the memory limit for the container is set,
the size of the volume is 50% of the memory allocated for the node.

A task can be executed either pre-deployment or post-deployment
as specified in the pod template specs of your Workloads
([Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
[DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
and
[ReplicaSets](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/).
See
[Pre- and post-deployment tasks](../implementing/integrate/#pre--and-post-deployment-checks)
for details.
Note that the annotation identifies the task by `name`.
This means that you can modify the `function` code in the resource definition
and the revised code is picked up without additional changes.

## Examples for Deno-runtime

### Example 1: inline script for a Deno script

This example defines a full-fledged Deno script
within the `KeptnTaskDefinition` YAML file:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition
metadata:
  name: hello-keptn-inline
spec:
  function:
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
  function:
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
  function:
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
  function:
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

## Examples for a custom container

### Example 1: Custom container

For an example of a `KeptnTaskDefinition` that defines a custom container.
 see
[container-task.yaml](<https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/base/container-task.yaml>.
The `spec` includes:

```yaml
spec:
  container:
    name: testy-test
    image: busybox:1.36.0
    command:
      - 'sh'
      - '-c'
      - 'sleep 30'
```

This task is then referenced in

[app.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/version-3/app.yaml).

This is a a trivial example that just runs `busybox`,
then spawns a shell and runs the `sleep 30` command.

## Examples for a python container

### Example 1: inline code for a Python container

You can embed python code directly in the task definition.
This example prints data stored in the parameters map:
{{< readfile file="/yaml_py/taskdefinition_pyfunction_inline.yaml" code="true" lang="yaml" >}}

### Example 2: httpRef for a Python container

You can refer to code stored online.
For example, we have a few examples available
[here](https://github.com/keptn/lifecycle-toolkit/tree/main/python-runtime/samples).
Consider the following:
{{< readfile file="/yaml_py/taskdefinition_pyfunction_configmap.yaml" code="true" lang="yaml" >}}

### Example 3: functionRef for a Python container

You can refer to an existing `KeptnTaskDefinition`.
This example calls the inline example
but overrides the data printed with what is specified in the task:
{{< readfile file="/yaml_py/taskdefinition_pyfunction_recursive.yaml" code="true" lang="yaml" >}}

### Example 4: ConfigMapRef for a Python container

{{< readfile file="/yaml_py/taskdefinition_pyfunction_configmap.yaml" code="true" lang="yaml" >}}

### Allowed libraries for the Python runtime

The following example shows how to use few of the allowed packages, namely:
requests, json, git, and yaml:

{{< readfile file="/yaml_py/taskdefinition_pyfunction_inline_printargs_py.yaml" code="true" lang="yaml" >}}

### Passing secrets, environment variables and modifying the python command

The following examples show how to pass data inside the parameter map,
how to load a secret in your code,
and how to modify the python command.
In this case the container run with the `-h` option
which prints the help message for the python3 interpreter:

{{< readfile file="/yaml_py/taskdefinition_pyfunction_use_envvars.yaml" code="true" lang="yaml" >}}

## More examples

See the [operator/config/samples](https://github.com/keptn/lifecycle-toolkit/tree/main/operator/config/samples/function_execution)
directory for more example `KeptnTaskDefinition` YAML files.

## Files

API Reference:

* [KeptnTaskDefinition](../crd-ref/lifecycle/v1alpha3/_index.md#keptntaskdefinition)
* [KeptnTaskDefinitionList](../crd-ref/lifecycle/v1alpha3/_index.md#keptntaskdefinitionlist)
* [KeptnTaskDefinitionSpec](../crd-ref/lifecycle/v1alpha3/_index.md#keptntaskdefinitionspec)
* [FunctionReference](../crd-ref/lifecycle/v1alpha3/_index.md#functionreference)
* [FunctionSpec](../crd-ref/lifecycle/v1alpha3/_index.md#functionspec)
* [FunctionStatus](../crd-ref/lifecycle/v1alpha3/_index.md#functionstatus)
* [HttpReference](../crd-ref/lifecycle/v1alpha3/_index.md#httpreference)
* [Inline](../crd-ref/lifecycle/v1alpha3/_index.md#inline)

## Differences between versions

The `KeptnTaskDefinition` support for containers is introduced in v0.8.0
to support the container runtime feature.
This modifies the synopsis in two ways:

* Add the `spec.container` field.
* The `spec.function` field is changed to pointer receiver.
  This aligns it with the `spec.container` field,
  which must be a pointer,
  and enables `KeptnTask` to omit it when it is empty,
  which it must be when `spec.container` is populated.

## Limitations

Only one
[runtime](https://kubernetes.io/docs/setup/production-environment/container-runtimes/)
is allowed per `KeptnTaskDefinition`.

## See also

* [KeptnApp](app.md)
* [Working with tasks](../implementing/tasks)
* [Working with container runtimes](../implementing/container.md)
* [Pre- and post-deployment tasks](../implementing/integrate/#pre--and-post-deployment-checks)
* [Orchestrate deployment checks](../getting-started/orchestrate)
