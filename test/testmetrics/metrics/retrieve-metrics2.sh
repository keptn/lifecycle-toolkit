#!/bin/bash

# Usage: ./retrieve-metrics.sh <component_name>

component_name=$1

# Replace with the actual URL where the metrics are exposed
metrics_url="http://localhost:2222/metrics"

# Fetch metrics for the specified component
metrics=$(curl -s $metrics_url | grep "${component_name}_active")

echo "Metrics for $component_name:"
echo "$metrics"
