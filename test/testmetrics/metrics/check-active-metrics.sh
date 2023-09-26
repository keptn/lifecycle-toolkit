#!/bin/bash

# Usage: ./retrieve-metrics.sh <component_name>

component_name=$1

metrics_url="http://lifecycle-operator-metrics-service.keptn-lifecycle-toolkit-system.svc.cluster.local:2222/metrics"

# Fetch keptn_lifecycle_active metrics
metrics=$(curl -s $metrics_url | grep "${component_name}")

echo "$metrics"
