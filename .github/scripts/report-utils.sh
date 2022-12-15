#!/bin/bash

ignore="--ignore-not-found"

createResourceReport () {
    path=$1
    namespace=$2
    resource=$3
    withLogs=$4

    mkdir -p "$path/$resource"

    kubectl get "$resource" -n "$namespace" "$ignore" > "$path/$resource/list-$resource.txt"

    for r in $(kubectl get "$resource" -n "$namespace" "$ignore" -o jsonpath='{.items[*].metadata.name}'); do
        kubectl describe "$resource/$r" -n "$namespace" > "$path/$resource/$r-describe.txt"

        if $withLogs ; then
        kubectl logs "$resource/$r" --all-containers=true -n "$namespace" > "$path/$resource/$r-logs.txt"
        fi
    done
}
