---
title: Working with Keptn tasks
description: Learn how to work with Keptn tasks
weight: 90
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

Keptn tasks are defined in a
[KeptnTaskDefinition](../../yaml-crd-ref/taskdefinition.md/)
resource.
A task definition includes a function
that defines the action taken by that task.
It can be configured in one of three different ways:

- inline
- referring to an HTTP script
- referring to another `KeptnTaskDefinition`
- referring to a
  [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/)
  resource that is populated with the function to execute

### Context

A Kubernetes context is a set of access parameters
that contains a Kubernetes cluster, a user, a namespace,
the application name, workload name, and version.
For more information, see
[Configure Access to Multiple Clusters](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/).

You may need to include context information in the `function` code
included in the YAML file that defines a
[KeptnTaskDefinition](../../yaml-crd-ref/taskdefinition.md)
resource.
For an example of how to do this, see the
[keptn-tasks.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/keptn-tasks.yaml)
file.

A context environment variable is available via `Deno.env.get("CONTEXT")`.
It can be used like this:
  
```javascript
let context = Deno.env.get("CONTEXT");
    
if (context.objectType == "Application") {
    let application_name = contextdata.appName;
    let application_version = contextdata.appVersion;
}       
        
if (context.objectType == "Workload") {
    let application_name = contextdata.appName;
    let workload_name = contextdata.workloadName;
    let workload_version = contextdata.workloadVersion;
}
```

## Parameterized functions

`KeptnTaskDefinition`s can use input parameters.
Consider the following example:

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

Note the following about using parameters with functions:

- The Lifecycle Toolkit passes the values
  defined inside the `map` field as a JSON object.
- Multi-level maps are not currently supported.
- The JSON object can be read through the environment variable `DATA`
  using `Deno.env.get("DATA");`.
- Currently only one secret can be passed.
  The secret must have a `key` called `SECURE_DATA`.
  It can be accessed via the environment variable `Deno.env.get("SECURE_DATA")`.

## Create secret text

To create a secret to use in a `KeptnTaskDefinition`,
execute this command:

```shell
# kubectl create secret generic my-secret --from-literal=SECURE_DATA=foo
```
```yaml
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

```shell
# kubectl create secret generic my-secret \
# --from-literal=SECURE_DATA="{\"foo\": \"bar\", \"foo2\": \"bar2\"}"
```
```yaml
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

## Pass secrets to a function

In the previous example, you see that
Kubernetes
[secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
can be passed to the function
using the `secureParameters` field.

Here, the `secret` value is the name of the Kubernetes secret,
which contains a field with the key `SECURE_DATA`.  
The value of that field is then available to the function's runtime
via an environment variable called `SECURE_DATA`.

For example, if you have a task function that should make use of secret data,
you must first ensure that the secret containing the `SECURE_DATA` key exists
For example:

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
