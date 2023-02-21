---
title: KeptnTaskDefinition
description: Define tasks that can be run pre- or post-deployment
weight: 89
---

A `KeptnTaskDefinition` is a CRD used to define tasks
that can be run by the Keptn Lifecycle Toolkit
as part of pre- and post-deployment phases of a deployment.
The task definition is a [Deno](https://deno.land/) script.
Please, refer to the [function runtime](https://github.com/keptn/lifecycle-toolkit/tree/main/functions-runtime)
for more information about the runtime.
In the future, we also intend to support other runtimes,
especially running a container image directly.

## Synopsis

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

* **apiVersion** - API version being used.
  The default for KLT Release 0.5.0 is `lifecycle.keptn.sh/v1alpha2`.
  * This must match <whatever>
  * Other information
`
* **kind** CRD type.  This is `KeptnTaskDefinition`

* **name** Unique name of this task.
  This must be an alphanumeric string and, by convention, is all lowercase.
  It can use the special characters `_`, `-` ... <what others>.
  It should not inclue spaces.

* **function** - Code to be executed.
  This can be expressed as one of the following:

  * **inline** - Include the actual executable code to execute.
    This can be written as a full-fledged Deno script.
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
  that is used in multiple place with different parameters.
  An example is:
   ```yaml
   spec:
     function:
       functionRef:
         name: slack-notification
   ```

  This can also be used to group a set of tasks into a single `KeptnTaskDefinitions`,
  such as defining a `KeptnTaskDefinition` for testing
  and have it call a `KeptnTaskDefinition` for each type of test to be run.

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
  The `secret` value is the Kubernetes secrete name
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
<!-- Instructions and guidelines for when and how to customize a CRD -->

## Examples

## Files

## Differences between versions

## See also
