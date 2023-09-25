#!/bin/bash

# Usage: ./retrieve-metrics.sh <component_name>

component_name=$1

# Metrics URL is now localhost:2222
metrics_url="http://lifecycle-operator-metrics-service.keptn-lifecycle-toolkit-system.svc.cluster.local:2222/metrics"

# Fetch metrics for the specified component
metrics=$(curl -s $metrics_url | grep "${component_name}_active")

echo "Metrics for $component_name:"
echo "$metrics"
