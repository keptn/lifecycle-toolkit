#!/bin/bash

logsDir="logs"

# shellcheck source=report_utils.sh
source report-utils.sh

# Go through each namespace in the cluster
for namespace in $(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}'); do

    mkdir -p "$logsDir/$namespace"
    createResourceReport "$logsDir/$namespace" "$namespace" "Pods" true
    createResourceReport "$logsDir/$namespace" "$namespace" "Deployments" false
    createResourceReport "$logsDir/$namespace" "$namespace" "Daemonsets" false
    createResourceReport "$logsDir/$namespace" "$namespace" "Statefulsets" false
    createResourceReport "$logsDir/$namespace" "$namespace" "Jobs" false
    
done
