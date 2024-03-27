---
comments: true
---

# Deployment tasks

A
[KeptnTaskDefinition](../reference/crd-reference/taskdefinition.md)
resource defines one or more "executables"
(functions, programs, scripts, etc.)
that Keptn runs
as part of the pre- and post-deployment phases of a
[KeptnApp](../reference/crd-reference/app.md) or
[KeptnWorkload](../reference/api-reference/lifecycle/v1/index.md#keptnworkload).

- pre-deployment (before the pod is scheduled)
- post-deployment (after the pod is scheduled)

These `KeptnTask` resources and the
`KeptnEvaluation` resources (discussed in
[Evaluations](./evaluations.md))
are part of the Keptn Release Lifecycle Management.

A
[KeptnTask](../reference/api-reference/lifecycle/v1/index.md#keptntask)
executes as a runner in an application
[container](https://kubernetes.io/docs/concepts/containers/),
which runs as part of a Kubernetes
[job](https://kubernetes.io/docs/concepts/workloads/controllers/job/).
A `KeptnTaskDefinition` includes calls to executables to be run.

To implement a `KeptnTask`:

- Define a
  [KeptnTaskDefinition](../reference/crd-reference/taskdefinition.md)
  resource that defines the runner to use for the container
  and the executables to be run
  pre- and post-deployment
- Apply [basic-annotations](./integrate.md#basic-annotations)
  to your workloads.
- Generate the required
  [KeptnApp](../reference/crd-reference/app.md)
  resources following the instructions in
  [Auto app discovery](auto-app-discovery.md).
- Annotate your workload YAML files
  to associate your `KeptnTaskDefinition`
  with the pre-/post-deployment tasks that should be run.
- Create the appropriate `KeptnAppContext` resource
  to associate `KeptnTaskDefinition` resources to the generated `KeptnApp`.

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

Keptn provides a general Kubernetes container runtime
that you can configure to do almost anything you want:

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
[KeptnTaskDefinition](../reference/crd-reference/taskdefinition.md)
reference page for the synopsis and examples for each runner.

## Run a task associated with your workload deployment

To define pre-/post-deployment tasks,
you must manually edit the YAML files
to add annotations for your tasks to the appropriate
workload YAML file.

Specify one of the following annotations/labels
for each task you want to execute:

```yaml
keptn.sh/pre-deployment-tasks: <task-name>
keptn.sh/post-deployment-tasks: <task-name>
```

The value of each annotation corresponds
to the value of the `name` field of the
[KeptnTaskDefinition](../reference/crd-reference/taskdefinition.md)
resource.

## Run a task associated with your entire KeptnApp

To execute pre-/post-deployment tasks for a `KeptnApp`,
create a `KeptnAppContext` with the same name and in the same `namespace` as the `KeptnApp`.
The `KeptnAppContext` resource contains a list of
pre-/post-deployment tasks
that should be executed before and after the
workloads within the `KeptnApp` are deployed.

See the [Getting started guide](../getting-started/lifecycle-management.md#more-control-over-the-application)
for more information on how to configure a `KeptnAppContext` resource
to execute pre-/post-deployment checks.
`KeptnAppContext` is also used to collect user defined
metadata information to use during your Task execution.
As explained later in this guide.
(See [the context section](#context))

## Example of pre/post-deployment actions

A comprehensive example of pre-/post-deployment
evaluations and tasks can be found in our
[examples folder](https://github.com/keptn/lifecycle-toolkit/tree/main/examples/sample-app),
where we use [Podtato-Head](https://github.com/podtato-head/podtato-head)
to run some simple pre-deployment checks.
Check out the [readme](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/README.md)
to learn how to test this example on your machine.

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
  because it gives you the most flexibility.

- Use the `inline` syntax for one of the Keptn pre-defined runners
  (either `deno-runtime` or `python-runtime`)
  to code the actual calls inline in the `KeptnTaskDefinition` resource.
  See
  [Fields for pre-defined containers](../reference/crd-reference/taskdefinition.md#fields-for-predefined-containers)
  for more information.

- Create a script that calls the functions, programs, and scripts
  that need to execute sequentially
  and install this on a remote webserver that Keptn can access.
  Then use the `httpRef` syntax for one of the pre-defined runners
  to call this script from your `KeptnTaskDefinition`,
  which can set parameters for the script if appropriate.

For more details about implementing these options, see the
[KeptnTaskDefinition](../reference/crd-reference/taskdefinition.md)
page.

## Context

The Keptn task context includes details about the current deployment, application name, version, object type and other
user-defined metadata.
Keptn populates this metadata while running tasks before and after deployments, to provide the necessary context
associated with each task.

This contrasts with the Kubernetes context, which is a set of access parameters that defines the
specific cluster, user and namespace with which you interact.
For more information, see
[Configure Access to Multiple Clusters](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/).

For Tasks generated for applications running in Kubernetes,
Keptn populates a `KEPTN_CONTEXT` environment variable containing a set
of parameters that correlate a task
to a specific application/workload,
information about the phase in which the task is being executed,
as well as any metadata that has been attached to
the related application/workload

You can use this context information
in the `function` code in your
[KeptnTaskDefinition](../reference/crd-reference/taskdefinition.md)
resource.

`KEPTN_CONTEXT` is encoded as JSON and by default, contains:

- "appName"
- "appVersion"
- "workloadName"
- "workloadVersion"
- "taskType"
- "objectType"
- "traceparent"
- "metadata"

A Job created by a `KeptnTask` with `KEPTN_CONTEXT`, may look like the following

```yaml
{% include "./assets/tasks/job-context.yaml" %}
```

You can customize the metadata field to hold any key-value pair of interest to share among
your workloads and tasks in a `KeptnApp` (for instance a commit ID value).
To do so, the metadata needs to be specified for the workload or for the application.
Follow our guide on [Context and Metadata here](./metadata.md).

<!-- markdownlint-disable MD046 max-one-sentence-per-line-->

!!! note

    For an example of how to access the `KEPTN_CONTEXT`, follow our reference page examples
    [for deno](../reference/crd-reference/taskdefinition.md#accessing-keptn_context-in-a-deno-task)
    and [for python](../reference/crd-reference/taskdefinition.md#accessing-keptn_context-in-a-python-task).

<!-- markdownlint-enable MD046 max-one-sentence-per-line-->

## Parameterized functions

`KeptnTaskDefinition`s can use input parameters.
Simple parameters are passed as a single map of key values,
while the `secret` parameters refer to a single Kubernetes `secret`.

Consider the following example:

```yaml
{% include "./assets/tasks/slack.yaml" %}
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
{% include "./assets/tasks/dummy-task.yaml" %}
```

To pass multiple variables
you can create a Kubernetes secret using a JSON string:

```shell
kubectl create secret generic my-secret \
--from-literal=SECURE_DATA="{\"foo\": \"bar\", \"foo2\": \"bar2\"}"
```

```yaml
{% include "./assets/tasks/multi-secret.yaml" %}
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
{% include "./assets/tasks/secret-data.yaml" %}
```

Then, you can make use of that secret as follows:

```yaml
{% include "./assets/tasks/taskdef-secure-data.yaml" %}
```
