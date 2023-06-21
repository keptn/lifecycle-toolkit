---
title: Tasks
description: Learn what Keptn Tasks are and how to use them
icon: concepts
layout: quickstart
weight: 10
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

### Keptn Task Definition

A `KeptnTaskDefinition` is a CRD used to define tasks that can be run by the Keptn Lifecycle Toolkit
as part of pre- and post-deployment phases of a deployment.
`KeptnTaskDefinition` resource can be created in the namespace where the application is running, or
in the default KLT namespace, which will be the fallback option for the system to search.
The task definition is a [Deno](https://deno.land/) script
Please, refer to the [function runtime](https://github.com/keptn/lifecycle-toolkit/tree/main/functions-runtime) for more
information about the runtime.
In the future, we also intend to support other runtimes, especially running a container image directly.

A task definition can be configured in three different ways:

- inline
- referring to an HTTP script
- referring to another `KeptnTaskDefinition`

An inline task definition looks like the following:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnTaskDefinition
metadata:
  name: deployment-hello
spec:
  function:
    inline:
      code: |
        console.log("Deployment Task has been executed");
```

In the code section, it is possible to define a full-fletched Deno script.

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
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

The runtime can also fetch the script on the fly from a remote webserver.
For this, the CRD should look like the
following:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnTaskDefinition
metadata:
  name: hello-keptn-http
spec:
  function:
    httpRef:
      url: <url>
```

An example is
available [here](https://github.com/keptn-sandbox/lifecycle-toolkit-examples/blob/main/sample-app/version-1/app-pre-deploy.yaml)
.

Finally, `KeptnTaskDefinition` can build on top of other `KeptnTaskDefinition`s.
This is a common use case where a general function can be re-used in multiple places with different parameters.

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
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

## Context

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

## Input Parameters and Secret Handling

As you might have noticed, Task Definitions also have the possibility to use input parameters.
The Lifecycle Toolkit passes the values defined inside the `map` field as a JSON object.
At the moment, multi-level maps are not supported.
The JSON object can be read through the environment variable `DATA` using `Deno.env.get("DATA");`.
K8s secrets can also be passed to the function using the `secureParameters` field.
Currently only one secret can be passed.
The secret must have a `key` called `SECURE_DATA`.
It can be accessed via the environment variable `Deno.env.get("SECURE_DATA")`.

For example:

```yaml
# kubectl create secret generic my-secret --from-literal=SECURE_DATA=foo

apiVersion: lifecycle.keptn.sh/v1alpha1
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

This methodology supports multiple variables by creating a K8s secret with a JSON string:

```yaml
# kubectl create secret generic my-secret \
# --from-literal=SECURE_DATA="{\"foo\": \"bar\", \"foo2\": \"bar2\"}"

apiVersion: lifecycle.keptn.sh/v1alpha1
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

### Keptn Task

A Task is responsible for executing the TaskDefinition of a workload.
The execution is done spawning a K8s Job to handle a single Task.
In its state, it keeps track of the current status of the K8s Job created.
