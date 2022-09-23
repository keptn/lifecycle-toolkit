# Hello, Keptn!

## Goal
This example shows how to define a function inline and pass over parameters to this function.

## Usage
* Edit task.yaml and add your name to `spec.parameters.map.name`
* Apply the manifests: `kubectl apply -f *.yaml`

## Outcome
* A KeptnTaskDefinition `hello-keptn` should be created
* A KeptnTask `hello-developer` should be created
* You can track the state of the job with `kubectl get KeptnTask hello-developer`
```                                                                                                          
NAME                            APPLICATION   WORKLOAD      VERSION   JOB NAME                           STATUS
hello-developer                 my-app        my-workload   1.0       klc-my-app-my-workload-1.0-57692   Succeeded
```