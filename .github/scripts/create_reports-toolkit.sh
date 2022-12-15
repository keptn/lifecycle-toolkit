#!/bin/bash

logsDir="logs"

# shellcheck source=report_utils.sh
source report_utils.sh

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
    
done

