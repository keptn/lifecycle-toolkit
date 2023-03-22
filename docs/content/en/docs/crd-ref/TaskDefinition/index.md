---
title: KeptnTaskDefinition
description: Define tasks that can be run pre- or post-deployment
weight: 89
---

A `KeptnTaskDefinition` defines tasks
that can be run by the Keptn Lifecycle Toolkit
as part of pre- and post-deployment phases of a deployment.

## Yaml Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
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
  The default for KLT Release 0.5.0 is `lifecycle.keptn.sh/v1alpha2`.
  * This must match <whatever>
  * Other information
`
* **kind** -- CRD type.  Must be set to `KeptnTaskDefinition`

* **name** -- Unique name of this task.
  * Must be an alphanumeric string and, by convention, is all lowercase.
  * Can include the special characters `_`, `-`, <what others>.
  * Should not include spaces.

* **function** -- Code to be executed,
  expressed as a [Deno](https://deno.land/) script.
  Refer to [function runtime](https://github.com/keptn/lifecycle-toolkit/tree/main/functions-runtime)
  for more information about the runtime.
  In the future, we intend to support additional runtimes,
  especially running a container image directly.

  The `function` can be defined as one of the following:

  * **inline** - Include the actual executable code to execute.
    This can be written as a full-fledged Deno script
    that is included in this file.
    For example:

    ```function:
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
            url: <url>
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
    into a single `KeptnTaskDefinitions`,
    such as defining a `KeptnTaskDefinition` for testing.
    In this case, it calls a `KeptnTaskDefinition`
    for each type of test to be run.

    <Explain what happens if one task fails.
    Will the subsequent tasks execute or does the pipeline stop
    or pass control to whatever would execute next.
    Can I control that behavior?

* **parameters** - An optional field to supply input parameters to a function.
  The Lifecycle Toolkit passes the values defined inside the `map` field
  as a JSON object.
  At the moment, multi-level maps are not supported.
  For example:
   ```spec:
       parameters:
         map:
           textMessage: "This is my configuration"
   ```

   The JSON object can be read
   through the `DATA` environment variable using `Deno.env.get("DATA");`.

* **secureParameters** -- An optional field used to pass a Kubernetes secret.
  The `secret` value is the Kubernetes secret name
  that is mounted into the runtime
  and made available to the function
  using the `SECURE_DATA` environment variable.

  ```yaml
      secureParameters:
        secret: slack-token
   ```

## Usage

A Task is responsible for executing the TaskDefinition of a workload.
The execution is done by spawning a K8s Job to handle a single Task.
In its state, it keeps track of the current status of the K8s Job created.

<!-- How this CRD is "activated".  For example, which event uses this CRD -->
<!-- Can I execute tasks in parallel? -->
<!-- Instructions and guidelines for when and how to customize a CRD -->

## Examples

This section can do any of the following:

* Include annotated examples
* Link to formal `examples` that we provide;
  include an annotation about what they illustrate.
  In this case, it would be nice to point to examples
  of the `inline`, `httpRef`, and `functionRef` code
  for the `function` field.

## Files

* Link to source code file where this is defined.
* Should the links to the API Ref go here instead of "See also"?

## Differences between versions

## See also

* Link to "use-case" guide pages that do something interesting with this CRD
* Link to reference pages for any related CRDs

API Reference:

* [KeptnTaskDefinition](../../api-ref/lifecycle/v1alpha3/#keptntaskdefinition)
* [KeptnTaskDefinitionList](../../api-ref/lifecycle/v1alpha3/#keptntaskdefinitionlist)
* [KeptnTaskDefinitionSpec](../../api-ref/lifecycle/v1alpha3/#keptntaskdefinitionspec)
* [FunctionReference](../../api-ref/lifecycle/v1alpha3/#functionreference)
* [FunctionSpec](../../api-ref/lifecycle/v1alpha3/#functionspec)
* [FunctionStatus](../../api-ref/lifecycle/v1alpha3/#functionstatus)
* [HttpReference](../../api-ref/lifecycle/v1alpha3/#httpreference)
* [Inline](../../api-ref/lifecycle/v1alpha3/#inline)
