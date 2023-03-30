---
title: KeptnTaskDefinition
description: Define tasks that can be run pre- or post-deployment
weight: 89
---

A `KeptnTaskDefinition` defines tasks
that are run by the Keptn Lifecycle Toolkit
as part of the pre- and post-deployment phases of a `KeptnApp`.

## Yaml Synopsis

```yaml
apiVersion: <lifecycle.keptn.sh/v?alpha?>
kind: KeptnTaskDefinition
metadata:
  name: <task-name>
spec:
  function:
    inline: | httpRef | functionRef
    [parameters:
      map:
        textMessage: "This is my configuration"]
    [secureParameters:
      secret: slack-token]
```

## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.  Must be set to `KeptnTaskDefinition`

* **name** -- Unique name of this task.
  * Must be an alphanumeric string and, by convention, is all lowercase.
  * Can include the special characters `_`, `-`, <what others>.
  * Should not include spaces.

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

  * **functionRef** -- Execute another `KeptnTaskDefinition` that has been defined.
    Populate this field with the value of the `name` field
    for the `KeptnTaskDefinition` to be called.
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
    In this case, it other, existing `KeptnTaskDefinition`s
    for each type of test to be run,
    specifying each by the value of the `name` field.

* **parameters** - An optional field to supply input parameters to a function.
  The Lifecycle Toolkit passes the values defined inside the `map` field
  as a JSON object.
  For example:

   ```yaml
   spec:
     parameters:
       map:
         textMessage: "This is my configuration"

   The JSON object can be read
   through the `DATA` environment variable using `Deno.env.get("DATA");`.

   Multi-level maps are not supported at this time.

   Currently only one secret can be passed.
   The secret must have a `key` called `SECURE_DATA`.
   It can be accessed via the environment variable `Deno.env.get("SECURE_DATA")`.
   See [Context](#context) for details.

* **secureParameters** -- An optional field used to pass a Kubernetes secret.
  The `secret` value is the Kubernetes secret name
  that is mounted into the runtime and made available to functions
  using the `SECURE_DATA` environment variable.
  For example:

  ```yaml
  secureParameters:
    secret: slack-token

   See [Create secret text](#create-secret-text) for details.

## Usage

A Task is responsible for executing the TaskDefinition of a workload.
The execution is done by spawning a Kubernetes Job to handle a single Task.
In its state, it tracks the current status of this Kubernetes Job.

### Context

A context environment variable is available via `Deno.env.get("CONTEXT")`.
It can be used like this:

```javascript
let context = Deno.env.get("CONTEXT");

if (contextdata.objectType == "Application") {
    let application_name = contextdata.appName;
    let application_version = contextdata.appVersion;
}

if (contextdata.objectType == "Workload") {
    let application_name = contextdata.appName;
    let workload_name = contextdata.workloadName;
    let workload_version = contextdata.workloadVersion;
}
```

### Create secret text

```yaml
# kubectl create secret generic my-secret --from-literal=SECURE_DATA=foo

apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition
metadata:
  name: dummy-task
  namespace: "default"
spec:
  function:
    secureParameters:
      secret: my-secret
    inline:
      code: |
        let secret_text = Deno.env.get("SECURE_DATA");
        // secret_text = "foo"
```

This methodology supports multiple variables
by creating a Kubernetes secret with a JSON string:

```yaml
# kubectl create secret generic my-secret \
# --from-literal=SECURE_DATA="{\"foo\": \"bar\", \"foo2\": \"bar2\"}"

apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition
metadata:
  name: dummy-task
  namespace: "default"
spec:
  function:
    secureParameters:
      secret: my-secret
    inline:
      code: |
        let secret_text = Deno.env.get("SECURE_DATA");
        let secret_text_obj = JSON.parse(secret_text);
        // secret_text_obj["foo"] = "bar"
        // secret_text_obj["foo2"] = "bar2"
```

## Examples

**Example 1** defines a full-fledged Deno script
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

**Example 2** fetches the Deno script from a remote webserver at runtime:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition
metadata:
  name: hello-keptn-http
spec:
  function:
    httpRef:
      url: <url>
```

See the
[sample-app/version-1](https://github.com/keptn-sandbox/lifecycle-toolkit-examples/blob/main/sample-app/version-1/app-pre-deploy.yaml)
PodtatoHead example for a more complete example.

**Example 3** calls another defined task,
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

<<<<<<< HEAD
=======
**More examples**

* See the [operator/config/samples](https://github.com/keptn/lifecycle-toolkit/tree/main/operator/config/samples)
directory for more example `KeptnTaskDefinition` YAML files.
Separate examples are provided for each API version.
For example, the `lifecycle_v1alpha3_keptntaskdefinition` file
contains examples for the `v1alpha3` library.

>>>>>>> 0333070 (library -> API)
## Files

API Reference:

* [KeptnTaskDefinition](../../api-ref/lifecycle/v1alpha3/#keptntaskdefinition)
* [KeptnTaskDefinitionList](../../api-ref/lifecycle/v1alpha3/#keptntaskdefinitionlist)
* [KeptnTaskDefinitionSpec](../../api-ref/lifecycle/v1alpha3/#keptntaskdefinitionspec)
* [FunctionReference](../../api-ref/lifecycle/v1alpha3/#functionreference)
* [FunctionSpec](../../api-ref/lifecycle/v1alpha3/#functionspec)
* [FunctionStatus](../../api-ref/lifecycle/v1alpha3/#functionstatus)
* [HttpReference](../../api-ref/lifecycle/v1alpha3/#httpreference)
* [Inline](../../api-ref/lifecycle/v1alpha3/#inline)

## Differences between versions

The `KeptnTaskDefinition` is the same for
all `v1alpha?` library versions.

## See also

* Link to "use-case" guide pages that do something interesting with this CRD
* Link to reference pages for any related CRDs
