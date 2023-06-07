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
    createResourceReport "$logsDir/$namespace" "$namespace" "ConfigMaps" false
    createResourceReport "$logsDir/$namespace" "$namespace" "KeptnApp" false
    createResourceReport "$logsDir/$namespace" "$namespace" "KeptnAppVersion" false
    createResourceReport "$logsDir/$namespace" "$namespace" "KeptnEvaluationDefinition" false
    createResourceReport "$logsDir/$namespace" "$namespace" "KeptnEvaluationProvider" false
    createResourceReport "$logsDir/$namespace" "$namespace" "KeptnEvaluation" false
    createResourceReport "$logsDir/$namespace" "$namespace" "KeptnTaskDefinition" false
    createResourceReport "$logsDir/$namespace" "$namespace" "KeptnTask" false
    createResourceReport "$logsDir/$namespace" "$namespace" "KeptnWorkload" false
    createResourceReport "$logsDir/$namespace" "$namespace" "KeptnWorkloadInstance" false
    createResourceReport "$logsDir/$namespace" "$namespace" "KeptnMetric" false
    createResourceReport "$logsDir/$namespace" "$namespace" "KeptnMetricsProvider" false
    
done
