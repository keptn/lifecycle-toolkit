---
title: Write Keptn Tasks
description: Learn how to use the Keptn Lifecycle Toolkit and explore basic features.
icon: concepts
layout: quickstart
weight: 20
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

## Keptn Task Definition

A `KeptnTaskDefinition` is a CRD used to define tasks that can be run by the Keptn Lifecycle Toolkit
as part of pre- and post-deployment phases of a deployment.
The task definition is a [Deno](https://deno.land/) script.
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
Kubernetes secrets can also be passed to the function using the `secureParameters` field.

Here, the `secret` value is the name of the K8s secret containing a field with the key `SECURE_DATA`.  
The value of that field will then be available to the functions runtime via an environment variable called `SECURE_DATA`.

For example, if you have a task function that should make use of secret data, you must first ensure that the secret
containing the `SECURE_DATA` key exists, as e.g.:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: deno-demo-secret
  namespace: default
type: Opaque
data:
  SECURE_DATA: YmFyCg== # base64 encoded string, e.g. 'bar'
```

Then, you can make use of that secret as follows:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition
metadata:
  name: deployment-hello
  namespace: "default"
spec:
  function:
    secureParameters:
      secret: deno-demo-secret
    inline:
      code: |
        console.log("Deployment Hello Task has been executed");

        let foo = Deno.env.get('SECURE_DATA');
        console.log(foo);
        Deno.exit(0);
```
