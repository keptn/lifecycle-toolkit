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

## Yaml Synopsis

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

## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.
   Must be set to `KeptnTaskDefinition`

* **metadata**
  * **name** -- Unique name of this task.
    Names must comply with the
    [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
    specification.

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

## Usage

A Task executes the TaskDefinition of a
[KeptnApp](app.md) or [KeptnWorkload].
The execution is done by spawning a Kubernetes
[Job](https://kubernetes.io/docs/concepts/workloads/controllers/job/)
to handle a single Task.
In its state, it tracks the current status of this Kubernetes Job.

The `function` is coded in JavaScript
and executed in
[Deno](https://deno.com/runtime),
which is a lightweight runtime environment
that executes in your namespace.
Note that Deno has tighter restrictions
for permissions and importing data
so a script that works properly elsewhere
may not function out of the box when run in Deno.

A task can be executed either pre-deployment or post-deployment
as specified in the `Deployment` resource;
see
[Pre- and post-deployment tasks](../implementing/integrate/#pre--and-post-deployment-checks)
for details.
Note that the annotation identifies the task by `name`.
This means that you can modify the `function` code in the resource definition
and the revised code is picked up without additional changes.

## Examples

### Example 1: inline script

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

### Example 2: httpRef script

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

### Example 3: functionRef

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

### Example 4: ConfigMapRef

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

### Example 5: ConfigMap

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

### More examples

See the [operator/config/samples](https://github.com/keptn/lifecycle-toolkit/tree/main/operator/config/samples)
directory for more example `KeptnTaskDefinition` YAML files.
Separate examples are provided for each API version.
For example, the `lifecycle_v1alpha3_keptntaskdefinition` file
contains examples for the `v1alpha3` version of the lifecycle API group.

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

The `KeptnTaskDefinition` is the same for
all `v1alpha?` library versions.

## See also

* [Working with tasks](../implementing/tasks)
* [Pre- and post-deployment tasks](../implementing/integrate/#pre--and-post-deployment-checks)
* [Orchestrate deployment checks](../getting-started/orchestrate)
