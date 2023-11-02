---
title: Deployment tasks
description: Learn how to work with Keptn tasks
weight: 90
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

A
[KeptnTaskDefinition](../yaml-crd-ref/taskdefinition.md/)
resource defines one or more "executables"
(functions, programs, scripts, etc)
that Keptn runs
as part of the pre- and post-deployment phases of a
[KeptnApp](../yaml-crd-ref/app.md) or
[KeptnWorkload](../crd-ref/lifecycle/v1alpha3/#keptnworkload).

- pre-deployment (before the pod is scheduled)
- post-deployment (after the pod is scheduled)

These `KeptnTask` resources and the
`KeptnEvaluation` resources (discussed in
[Evaluations](./evaluations.md))
are part of the Keptn Release Lifecycle Management.

A
[KeptnTask](../crd-ref/lifecycle/v1alpha3/#keptntask)
executes as a runner in an application
[container](https://kubernetes.io/docs/concepts/containers/),
which runs as part of a Kubernetes
[job](https://kubernetes.io/docs/concepts/workloads/controllers/job/).
A `KeptnTaskDefinition` includes calls to executables to be run.

To implement a `KeptnTask`:

- Define a
  [KeptnTaskDefinition](../yaml-crd-ref/taskdefinition.md)
  resource that defines the runner to use for the container
  and the executables to be run
pre- and post-deployment
- Apply [basic-annotations](./integrate.md#basic-annotations)
  to your workloads to integrate your task with Kubernetes and,
  if desired, creates a `KeptnApp` resource
  that consolidates multiple workloads into a single application
- Annotate the appropriate
  [KeptnApp](../yaml-crd-ref/app.md)
  resource to associate your `KeptnTaskDefinition`
  with the pre- and post-deployment tasks that should be run;
  see
  [Pre- and post-deployment tasks and checks](./integrate.md#pre--and-post-deployment-checks)
  for more information

This page provides information to help you create your tasks:

- Code your task in an appropriate [runner](#runners-and-containers)
- How to control the
  [execution order](#executing-sequential-tasks)
  of functions, programs, and scripts
  since all `KeptnTask` resources at the same level run in parallel
- Understand how to use [Context](#context)
  that contains a Kubernetes cluster, a user, a namespace,
  the application name, workload name, and version.
- Use [parameterized functions](#parameterized-functions)
  if your task requires input parameters
- [Create secret text](#create-secret-text)
  and [pass secrets to a function](#pass-secrets-to-a-function)
  if necessary.

## Runners and containers

Each `KeptnTaskDefinition` can use exactly one container with one runner.
The runner you use determines the language you can use
to define the task.
The `spec` section of the `KeptnTaskDefinition`
defines the runner to use for the container:

Keptn provides a general Kubernetes that you can configure
to do almost anything you want:

- The `container-runtime` runner provides
  a pure custom Kubernetes application container
  that you define to includes a runtime, an application
  and its runtime dependencies.
  This gives you the greatest flexibility
  to define tasks using the language and facilities of your choice

Keptn also includes two "pre-defined" runners:

- Use the `deno-runtime` runner to define tasks using Deno scripts,
  which use JavaScript/Typescript syntax with a few limitations.
  You can use this to specify simple actions
  without having to define a container.
- Use the `python-runtime` runner
  to define your task using Python 3.

For the pre-defined runners (`deno-runtime` and `python-runtime`),
the actual code to be executed
can be configured in one of four different ways:

- inline
- referring to an HTTP script
- referring to another `KeptnTaskDefinition`
- referring to a
  [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/)
  resource that is populated with the function to execute

See the
[KeptnTaskDefinition](../yaml-crd-ref/taskdefinition.md)
reference page for the synopsis and examples for each runner.

## Executing sequential tasks

All `KeptnTask` resources that are defined by
`KeptnTaskDefinition` resources at the same level
(either pre-deployment or post-deployment) execute in parallel.
This is by design, because Keptn is not a pipeline engine.
**Task sequences that are not part of the lifecycle workflow
should not be handled by Keptn**
but should instead be handled by the pipeline engine tools being used
such as Jenkins, Argo Workflows, Flux, and Tekton.

If your lifecycle workflow includes
a sequence of executables that need to be run in order,
you can put them all in one `KeptnTaskDefinition` resource,
which can execute a virtually unlimited number
of programs, scripts, and functions,
as long as they are all using the same runner.
You have the following options:

- Encode all your steps in the language of your choice
  and build an image
  that Keptn executes in a `container-runtime` runner.
  This is often the best solution if you need to execute complex sequences
  because it gives you the most flexibility..

- Use the `inline` syntax for one of the Keptn pre-defined runners
  (either `deno-runtime` or `python-runtime`)
  to code the actual calls inline in the `KeptnTaskDefinition` resource.
  See
  [Fields for pre-defined containers](../yaml-crd-ref/taskdefinition.md/#fields-for-pre-defined-containers)
  for more information.

- Create a script that calls the functions, programs, and scripts
  that need to execute sequentially
  and install this on a remote webserver that Keptn can access.
  Then use the `httpRef` syntax for one of the pre-defined runners
  to call this script from your `KeptnTaskDefinition`,
  which can set parameters for the script if appropriate.

For more details about implementing these options, see the
[KeptnTaskDefinition](../yaml-crd-ref/taskdefinition.md)
page.

## Context

A Kubernetes context is a set of access parameters
that contains a Kubernetes cluster, a user, a namespace,
the application name, workload name, and version.
For more information, see
[Configure Access to Multiple Clusters](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/).

You may need to include context information in the `function` code
included in the YAML file that defines a
[KeptnTaskDefinition](../yaml-crd-ref/taskdefinition.md)
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
Simple parameters are passed as a single map of key values,
while the `secret` parameters refer to a single Kubernetes `secret`.

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

- Keptn passes the values
  defined inside the `map` field as a JSON object.
- Multi-level maps are not currently supported.
- The JSON object can be read through the environment variable `DATA`
  using `Deno.env.get("DATA");`.
- Currently only one secret can be passed.
  The secret must have a `key` called `SECURE_DATA`.
  It can be accessed via the environment variable `Deno.env.get("SECURE_DATA")`.

## Working with secrets

A special case of parameterized functions
is to pass secrets that may be required
to access data that your task requires.

### Create secret text

To create a secret to use in a `KeptnTaskDefinition`,
execute this command:

```shell
kubectl create secret generic my-secret --from-literal=SECURE_DATA=foo
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

To pass multiple variables
you can create a Kubernetes secret using a JSON string:

```shell
kubectl create secret generic my-secret \
--from-literal=SECURE_DATA="{\"foo\": \"bar\", \"foo2\": \"bar2\"}"
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

### Pass secrets to a function

[Kubernetes secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
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
