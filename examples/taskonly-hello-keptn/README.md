# Hello, Keptn!

## Goal
This example shows how to define a function inline and pass over parameters to this function.

## Variants
* `inline` - shows how to specify a function in the KeptnTaskDefinition
* `http` - fetches the Script from the Web
* `upstream` - shows how functions could be reused
## Usage
* Edit task.yaml and add your name to `spec.parameters.map.name`
* Choose the corresponding folder
* Apply the manifests: `kubectl apply -f .`

## Outcome
* A KeptnTaskDefinition `hello-keptn-<variant>` should be created
* A KeptnTask `hello-developer` should be created
* You can track the state of the job with `kubectl get KeptnTask hello-developer`
```                                                                                                          
NAME                            APPLICATION   WORKLOAD      VERSION   JOB NAME                           STATUS
hello-developer                 my-app        my-workload   1.0       klc-my-app-my-workload-1.0-57692   Succeeded
```