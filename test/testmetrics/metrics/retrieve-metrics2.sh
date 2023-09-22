#!/bin/bash

# Usage: ./retrieve-metrics.sh <component_name>

component_name=$1

# Use kubectl port-forward to forward local port 2222 to the metrics service
kubectl -n keptn-lifecycle-toolkit-system port-forward service/lifecycle-operator-metrics-service 2222 &

# Sleep for a moment to allow port-forwarding to be established
sleep 2

# Metrics URL is now localhost:2222
metrics_url="http://localhost:2222/metrics"

# Fetch metrics for the specified component
metrics=$(curl -s $metrics_url | grep "${component_name}_active")

echo "Metrics for $component_name:"
echo "$metrics"

# Kill the port-forwarding process
pkill -f "kubectl -n keptn-lifecycle-toolkit-system port-forward"
