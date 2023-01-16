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
The task definition is a [Deno](https://deno.land/) script
Please, refer to the [function runtime](https://github.com/keptn/lifecycle-toolkit/functions-runtime/) for more information about the runtime.
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

In the code section, it is possible to define a full-fletched Deno script. An example for that would be:

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

The runtime can also fetch the script on the fly from a remote webserver. For this, the CRD should look like the following:

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

As you might have noticed, Task Definitions also have the possibility to use input parameters.
The Lifecycle Toolkit passes the values defined inside the `map` field as a JSON object.
At the moment, multi-level maps are not supported.
The JSON object can be read through the environment variable `DATA` using `Deno.env.get("DATA");`.
K8s secrets can also be passed to the function using the `secureParameters` field.
Here, the `secret` value is the K8s secret name that will be mounted into the runtime and made available to the function via the environment variable `SECURE_DATA`.


### Keptn Task

A Task is responsible for executing the TaskDefinition of a workload.
The execution is done spawning a K8s Job to handle a single Task.
In its state, it keeps track of the current status of the K8s Job created.
