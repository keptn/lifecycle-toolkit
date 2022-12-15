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