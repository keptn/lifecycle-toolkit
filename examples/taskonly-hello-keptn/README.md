# Hello, Keptn!

## Goal
This example shows how to define an inline function and pass over parameters to its execution.

## Variants
* `inline` - shows how to specify a function directly in the `KeptnTaskDefinition`
* `http` - fetches the Script from the Web
* `upstream` - shows how functions could be reused

## Usage

* Edit task.yaml and add your name to `spec.parameters.map.name`
* Navigate to the corresponding folder
* Apply the manifests: `kubectl apply -f .`

## Outcome

* A KeptnTaskDefinition `hello-keptn-<variant>` should be created
* A KeptnTask `hello-developer` should be created
* You can track the state of the job with `kubectl get KeptnTask hello-developer`
```                                                                                                          
NAME                            APPLICATION   WORKLOAD      VERSION   JOB NAME                           STATUS
hello-developer                 my-app        my-workload   1.0       klc-my-app-my-workload-1.0-57692   Succeeded
```

<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=858843d8-8da2-4ce5-a325-e5321c770a78" />