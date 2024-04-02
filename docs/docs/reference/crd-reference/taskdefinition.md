---
comments: true
---

# KeptnTaskDefinition

A `KeptnTaskDefinition` defines tasks
that Keptn runs as part of the pre- and post-deployment phases of a
[KeptnApp](./app.md) or
[KeptnWorkload](../api-reference/lifecycle/v1/index.md#keptnworkload).

A Keptn task executes as a
[runner](https://docs.gitlab.com/runner/executors/kubernetes.html#how-the-runner-creates-kubernetes-pods)
in an application
[container](https://kubernetes.io/docs/concepts/containers/),
which runs as part of a Kubernetes
[job](https://kubernetes.io/docs/concepts/workloads/controllers/job/).

Each `KeptnTaskDefinition` can use exactly one container with one runner.
which is one of the following,
differentiated by the `spec` section:

- The `container-runtime` runner provides
  functionality to run
  a standard Kubernetes container inside
  a Kubernetes job.
  You define the container image, and
  the arbitrary application inside it.
  This gives you the flexibility
  to define tasks using the language and facilities of your choice,
  although it is more complicated that using one of the pre-defined runtimes.
  See
  [Synopsis for container-runtime](#synopsis-for-container-runtime)
  and
  [Example for a container-runtime runner](#example-for-a-container-runtime-runner).

- Pre-defined containers

    - Use the pre-defined `deno-runtime` runner
      to define tasks using
      [Deno](https://deno.com/)
      scripts,
      which use JavaScript with a few limitations.
      You can use this to specify simple actions
      without having to define a full container.
      See [runtime examples](#examples-for-deno-runtime-and-python-runtime-runners)
    - Use the pre-defined `python-runtime` runner
      to define your task using
      [Python 3](https://www.python.org/).
      See [runtime examples](#examples-for-deno-runtime-and-python-runtime-runners)
      for practical usage of the pre-defined containers.

## Synopsis for all runners

The `KeptnTaskDefinition` Yaml files for all runners
include the same lines at the top.
These are described here.

```yaml
{% include "../../assets/crd/examples/synopsis-for-all-runners.yaml"  %}
```

### Fields used for all containers

- **apiVersion** -- API version being used.

- **kind** -- Resource type.
  Must be set to `KeptnTaskDefinition`

- **metadata**
    - **name** -- Unique name of this task or container.
      This is the name used to insert this task or container
      into the `preDeployment` or `postDeployment` list.
      Names must comply with the
      [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
      specification.
- **spec**
    - **deno | python | container** (required) -- Define the container type
      to use for this task.
      Each task can use one type of runner,
      identified by this field:

        - **deno** -- Use a `deno-runtime` runner
          and code the functionality in Deno script,
          which is similar to JavaScript and Typescript.
          See
          [Synopsis for deno](./#deno-runtime-synopsis).
        - **python** -- Use a `python-runtime` function
          and code the functionality in Python 3.
          See
          [Synopsis for python](./#python-runtime-synopsis).
        - **container** -- Use the runner defined
          for the `container-runtime` container.
          This is a standard Kubernetes container
          for which you define the image, runner, runtime parameters, etc.
          and code the functionality to match the container you define.
          See
          [Synopsis for container-runtime container](#synopsis-for-container-runtime).

    - **retries** -- specifies the number of times
      a job executing the `KeptnTaskDefinition`
      should be restarted if an attempt is unsuccessful.
    - **timeout** -- specifies the maximum time
      to wait for the task to be completed successfully.
      The value supplied should specify the unit of measurement;
      for example, `5s` indicates 5 seconds and `5m` indicates 5 minutes.
      If the task does not complete successfully within this time frame,
      it is considered to be failed.

## Synopsis for container-runtime

Use the `container-runtime` to specify your own
[Kubernetes container](https://kubernetes.io/docs/concepts/containers/)
and define the task you want to execute.

Task sequences that are not part of the lifecycle workflow and
should be handled by the pipeline engine tools being used
such as Jenkins, Argo Workflows, Flux, and Tekton.

If you are migrating from Keptn v1,
you can use a `container-runtime` to execute
almost anything that you implemented with JES for Keptn v1.

```yaml
{% include "../../assets/crd/examples/synopsis-for-container-runtime.yaml"  %}
```

### Fields used only for container-runtime

- **spec**
    - **container** -- Container definition.
        - **name** -- Name of the container that will run,
          which is not the same as the `metadata.name` field
          that is used in the `KeptnTaskDefinition` resource.
        - **image** -- name of the image you defined according to
          [image reference](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#image)
          and
          [image concepts](https://kubernetes.io/docs/concepts/containers/images/)
          and pushed to a registry
        - **other fields** -- The full list of valid fields is available at
          [ContainerSpec](../api-reference/lifecycle/v1/index.md#containerspec),
          with additional information in the Kubernetes
          [Container](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#Container)
          spec documentation.

## Synopsis for predefined containers

The predefined containers allow you to easily define a task
using either Deno or Python syntax.
You do not need to specify the image, volumes, and so forth.
Instead, just provide either a Deno or Python script
and Keptn sets up the container and runs the script as part of the task.

<!-- markdownlint-disable MD046 -->

=== "Deno-runtime synopsis"

    When using the `deno-runtime` runner to define a task,
    the executables are coded in
    [Deno-script](https://deno.com/manual),
    (which is mostly the same as JavaScript and TypeScript)
    and executed in the
    `deno-runtime` runner,
    which is a lightweight runtime environment
    that executes in your namespace.
    Note that Deno has tighter restrictions
    for permissions and importing data
    so a script that works properly elsewhere
    may not function out of the box when run in the `deno-runtime` runner.
    In this case you may want to use a custom container instead. 

    ```yaml
    {% include "../../assets/crd/examples/synopsis-for-deno-runtime-container.yaml" %}
    ```

=== "Python-runtime synopsis"

    When using the `python-runtime` runner to define a task,
    the executables are coded in python3.
    The runner enables the following packages: requests, json, git, yaml.
    
    ??? example 
        ```yaml
        {% include "../../assets/crd/python-libs.yaml" %}
        ```

    Note that other libraries may not function out of the box 
    in the `python-runtime` runner. 
    In this case you may want to use a custom container instead.

    ```yaml
    {% include "../../assets/crd/examples/synopsis-for-python-runtime-runner.yaml" %}
    ```

<!-- markdownlint-enable MD046 -->

### Fields for predefined containers

- **spec** -- choose either `deno` or `python`
    - **deno | python**
        - **deno** -- Specify that the task uses the `deno-runtime`
          and is expressed as a [Deno](https://deno.com/) script.
          Refer to [deno runtime](https://github.com/keptn/lifecycle-toolkit/tree/main/runtimes/deno-runtime)
          for more information about this runner.
        - **python** -- Identifies this as a Python runner.

            - **inline | httpRef | functionRef | ConfigMapRef** -- choose the syntax
              used to call the executables.
              Only one of these can be specified per `KeptnTaskDefinition` resource:

                - **inline** -- Include the actual executable code to execute.
                        You can code a sequence of executables here
                        that need to be run in order
                        as long as they are executables that are part of the lifecycle workflow.
                        Task sequences that are not part of the lifecycle workflow
                        should be handled by the pipeline engine tools being used
                        such as Jenkins, Argo Workflows, Flux, and Tekton.
                        See examples of usage for [deno](./#inline-script-for-deno)
                        and for [python](./#inline-script-for-python).

                - **httpRef** -- Specify a script to be executed at runtime
                        from the remote webserver that is specified.
                        This syntax allows you to call a general function
                        that is used in multiple places,
                        possibly with different parameters
                        that are provided in the calling `KeptnTaskDefinition` resource.
                        Another `KeptnTaskDefinition` resource could call this same script
                        but with different parameters.
                        Only one script can be executed.
                        Any other scripts listed here are silently ignored.
                        See examples of usage for [deno](./#httpref-for-deno)
                        and for [python](./#httpref-for-python).

                - **functionRef** -- Execute another `KeptnTaskDefinition` resources.
                    Populate this field with the value(s) of the `metadata.name` field
                    for each `KeptnDefinitionTask` to be called.
                    Like the `httpRef` syntax,this is commonly used
                    to call a general function that is used in multiple places,
                    possibly with different parameters
                    that are set in the calling `KeptnTaskDefinition` resource.
                    To be able to run the pre-/post-deployment task, you must create
                    the `KeptnAppContext` resource and link the `KeptnTaskDefinition`
                    in the pre-/post-deployment section of `KeptnAppContext`.
                    The `KeptnTaskDefinition` called with `functionref`
                    is the `parent task` whose runner is used for the execution
                    even if it is not the same runner defined in the
                    calling `KeptnTaskDefinition`.
                    Only one `KeptnTaskDefinition` resources can be listed
                    with the `functionRef` syntax
                    although that `KeptnTaskDefinition` can call multiple
                    executables (programs, functions, and scripts).
                    Any calls to additional `KeptnTaskDefinition` resources
                    are silently ignored.
                    See examples of usage for [deno](./#functionref-for-deno)
                    and [python](./#functionref-for-python).
                - **ConfigMapRef** -- Specify the name of a
                  [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/)
                  resource that contains the function to be executed.
                  See examples of usage for [deno](./#configmapref-for-deno)
                  and for [python](./#configmapref-for-python).

            - **parameters** -- An optional field
                to supply input parameters to a function.
                Keptn passes the values defined inside the `map` field
                as a JSON object.
                See [Parameterized functions](../../guides/tasks.md#parameterized-functions)
                for more information.
                Also see examples for [deno](./#env-var-in-deno)
                and
                [python](./#env-var-in-python).

            - **secureParameters** -- An optional field
                used to pass a Kubernetes secret.
                The `secret` value is the Kubernetes secret name
                that is mounted into the runtime and made available to functions
                using the `SECURE_DATA` environment variable.

                Note that, currently, only one secret can be passed
                per `KeptnTaskDefinition` resource.

                See [Create secret text](../../guides/tasks.md#create-secret-text)
                for details.
                Also see examples on secret usage in tasks runner
                for [deno](./#env-var-in-deno) and [python](./#env-var-in-python).

## Usage

A Task executes the TaskDefinition of a
[KeptnApp](app.md) or a
[KeptnWorkload](../api-reference/lifecycle/v1/index.md#keptnworkload).
The execution is done by spawning a Kubernetes
[Job](https://kubernetes.io/docs/concepts/workloads/controllers/job/)
to handle a single Task.
In its state, it tracks the current status of this Kubernetes Job.

When using a container runtime that includes a volume,
an `EmptyDir` volume is created
with the same name as is specified the container `volumeMount`.
Note that, if more `volumeMount`s are specified,
only one volume is created with the name of the first `volumeMount`.
By default, the size of this volume is 1GB.
If the memory limit for the container is set,
the size of the volume is 50% of the memory allocated for the node.

A task can be executed either pre-deployment or post-deployment
as specified in the pod template specs of your Workloads
([Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
[DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/),
and
[ReplicaSets](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/))
and in the
[KeptnApp](app.md) resource.
See
[Annotations to KeptnApp](../../guides/tasks.md/#run-a-task-associated-with-your-entire-keptnapp)
for details.
Note that the annotation identifies the task by `name`.
This means that you can modify the `function` code in the resource definition
and the revised code is picked up without additional changes.

All `KeptnTaskDefinition` resources specified to the `KeptnAppContext` resource
at the same stage (either pre- or post-deployment) run in parallel.
You can run multiple executables sequentially
either by using the `inline` syntax for a predefined container image
or by creating your own image
and running it in the Keptn `container-runtime` runner.
See
[Executing sequential tasks](../../guides/tasks.md#executing-sequential-tasks)
for more information.

## Example for a container-runtime runner

For an example of a `KeptnTaskDefinition` that defines a custom container.
see
[container-task.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/base/container-task.yaml).
This is a trivial example that just runs `busybox`,
then spawns a shell and runs the `sleep 30` command:

```yaml
{% include "../../assets/crd/task-definition.yaml" %}
```

This task is then referenced in the
[appcontext.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/sample-app/version-2/appcontext.yaml)
file.

## Examples for deno-runtime and python-runtime runners

<!-- markdownlint-disable MD046 max-one-sentence-per-line -->

??? abstract "Inline scripts"

    === "Inline script for deno"

        This example defines a full-fledged Deno script
        within the `KeptnTaskDefinition` YAML file:
    
        ```yaml
        {% include "../../assets/crd/examples/inline-script-for-deno-script.yaml" %} 
        ```

    === "Inline script for python"
    
        You can embed python code directly in the task definition.
        This example prints data stored in the parameters map:
    
        ```yaml
        {% include "../../assets/crd/python-inline.yaml" %}
        ```

??? abstract "HttpRef"

    === "httpRef for deno"
    
        This example fetches the Deno script from a remote webserver at runtime:
    
        ```yaml
        {% include "../../assets/crd/examples/httpref-script-for-deno-script.yaml" %}
        ```
          
        For other examples, see the [sample-app](https://github.com/keptn-sandbox/lifecycle-toolkit-examples/blob/main/sample-app/version-1/app-pre-deploy.yaml).
        and [sample-app/version-1](https://github.com/keptn-sandbox/lifecycle-toolkit-examples/blob/main/sample-app/version-1/app-pre-deploy.yaml)
        PodtatoHead example for a more complete example.
    
    === "httpRef for python"
    
        ```yaml
        {% include "https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/lifecycle-operator/config/samples/python_execution/taskdefinition_pyfunction_upstream_hellopy.yaml" %}
        ```

??? abstract "FunctionRef"

    === "functionRef for deno"
    
        This example calls another defined task,
        illustrating how one `KeptnTaskDefinition` can build
        on top of other `KeptnTaskDefinition`s.
        In this case, it calls `slack-notification-dev`,
        passing `parameters` and `secureParameters` to that other task:
    
        ```yaml
        {% include "../../assets/crd/examples/functionref-for-deno-script.yaml" %} 
        ```

    === "functionRef for python"
    
        You can refer to an existing `KeptnTaskDefinition`.
        this example calls the inline example
        but overrides the data printed with what is specified in the task:
    
        ```yaml
        {% include "../../assets/crd/python-recursive.yaml" %}
        ```

??? abstract "ConfigMap and ConfigMapRef"

    === "ConfigMapRef for deno"
    
        This example references a `ConfigMap` by the name of `dev-configmap`
        that contains the code for the function to be executed.
    
        ```yaml
        {% include "../../assets/crd/examples/configmap-for-deno-script.yaml" %} 
        ```
    
    === "ConfigMapRef for python"
    
        In this example the python runner refers to an existing configMap 
        called `python-test-cm`
    
        ```yaml
        {% include "../../assets/crd/python-configmap.yaml" %}
        ```

??? abstract "Accessing KEPTN_CONTEXT environment variable"

    For Tasks triggered as pre- and post- deployment of applications
    on Kubernetes, Keptn populates an environment variable called `KEPTN_CONTEXT`.
    As all environment variables, this can be accessed using language specific methods.

    === "Accessing KEPTN_CONTEXT in a Deno task"
    
        ```javascript
        let context = Deno.env.get("KEPTN_CONTEXT");
        console.log(context);
        ```
    
    === "Accessing KEPTN_CONTEXT in a Python task"
    
        ```python
        import os
        import yaml
        data = os.getenv('KEPTN_CONTEXT')
        dct = yaml.safe_load(data)
        meta= dct['metadata']
        print(meta)
        ```

??? abstract "Passing secrets, environment variables and modifying the runner command"

    === "Env var in deno"
    
    The following example shows how to pass data inside the parameter map,
    and how to load a secret in your code.
    The deno command does not takes modifiers so filling the `cmdParameters`
    will do nothing.

    ```yaml
    {% include "../../assets/crd/deno-context.yaml" %}
    ```

    === "Env var in python"
    
    The following example shows how to pass data inside the parameter map,
    how to load a secret in your code,
    and how to modify the python command.
    In this case the container runs with the `-h` option,
    which prints the help message for the python3 interpreter:

    ```yaml
    {% include "../../assets/crd/python-context.yaml" %}
    ```

<!-- markdownlint-enable MD046 max-one-sentence-per-line-->

## More examples

See
the [lifecycle-operator/config/samples](https://github.com/keptn/lifecycle-toolkit/tree/main/lifecycle-operator/config/samples)
directory for more example `KeptnTaskDefinition` YAML files.

## Files

API Reference:

- [KeptnTaskDefinition](../api-reference/lifecycle/v1/index.md#keptntaskdefinition)
- [KeptnTaskDefinitionList](../api-reference/lifecycle/v1/index.md#keptntaskdefinitionlist)
- [KeptnTaskDefinitionSpec](../api-reference/lifecycle/v1/index.md#keptntaskdefinitionspec)
- [FunctionReference](../api-reference/lifecycle/v1/index.md#functionreference)
- [FunctionSpec](../api-reference/lifecycle/v1/index.md#runtimespec)
- [FunctionStatus](../api-reference/lifecycle/v1/index.md#functionstatus)
- [HttpReference](../api-reference/lifecycle/v1/index.md#httpreference)
- [Inline](../api-reference/lifecycle/v1/index.md#inline)

## Differences between versions

The `KeptnTaskDefinition` support for
the `container-runtime` and `python-runtime` is introduced in v0.8.0.
This modifies the synopsis in the following ways:

- Add the `spec.container` field.
- Add the `python` descriptor for the `python-runtime` runner.
- Add the `container` descriptor for the `container-runtime` runner.
- Add the `deno` descriptor to replace `function`
  for the `deno-runtime` runner.
  The `function` identifier for the `deno-runtime` runner
  is deprecated;
  it still works for v 0.8.0 but will be dropped from future releases.
- The `spec.function` field is changed to be a pointer receiver.
  This aligns it with the `spec.container` field,
  which must be a pointer,
  and enables `KeptnTask` to omit it when it is empty,
  which it must be when `spec.container` is populated.

## Limitations

- Only one
  [runtime](https://kubernetes.io/docs/setup/production-environment/container-runtimes/)
  is allowed per `KeptnTaskDefinition` resource.

- Only one secret can be passed
  per `KeptnTaskDefinition` resource.

## See also

- [KeptnApp](app.md)
- [Working with tasks](../../guides/tasks.md)
- [KeptnApp and KeptnWorkload resources](../../components/lifecycle-operator/keptn-apps.md)
- [Release Lifecycle Management](../../getting-started/lifecycle-management.md)
- [Executing sequential tasks](../../guides/tasks.md#executing-sequential-tasks)
