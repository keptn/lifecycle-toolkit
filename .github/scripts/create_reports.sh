#!/bin/bash

ignore="--ignore-not-found"
logsDir="logs"

createResourceReport () {
  namespace=$1
  resource=$2
  withLogs=$3

  mkdir -p "$logsDir/$namespace/$resource"

  kubectl get "$resource" -n "$namespace" "$ignore" > "$logsDir/$namespace/$resource/list-$resource.txt"

  for r in $(kubectl get "$resource" -n "$namespace" "$ignore" -o jsonpath='{.items[*].metadata.name}'); do
    kubectl describe "$resource/$r" -n "$namespace" > "$logsDir/$namespace/$resource/$r-describe.txt"

    if $withLogs ; then
      kubectl logs "$resource/$r" --all-containers=true -n "$namespace" > "$logsDir/$namespace/$resource/$r-logs.txt"
    fi
  done
}

# Create a folder to store the logs
mkdir -p "$logsDir"

# Go through each namespace in the cluster
for namespace in $(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}'); do

    mkdir -p "$logsDir/$namespace"
    createResourceReport "$namespace" "Pods" true
    createResourceReport "$namespace" "Deployments" false
    createResourceReport "$namespace" "Daemonsets" false
    createResourceReport "$namespace" "Statefulsets" false
    createResourceReport "$namespace" "Jobs" false
    createResourceReport "$namespace" "KeptnApp" false
    createResourceReport "$namespace" "KeptnAppVersion" false
    createResourceReport "$namespace" "KeptnEvaluationDefinition" false
    createResourceReport "$namespace" "KeptnEvaluationProvider" false
    createResourceReport "$namespace" "KeptnEvaluation" false
    createResourceReport "$namespace" "KeptnTaskDefinition" false
    createResourceReport "$namespace" "KeptnTask" false
    createResourceReport "$namespace" "KeptnWorkload" false
    createResourceReport "$namespace" "KeptnWorkloadInstance" false
    
done

