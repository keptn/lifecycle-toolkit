#!/bin/bash
kubectl -n $NAMESPACE wait --for condition=established --timeout=30s crd/analyses.metrics.keptn.sh
crd_status=$(kubectl get crd/analyses.metrics.keptn.sh -n $NAMESPACE -o jsonpath='{.status.pass}')
expected_status="true"
if [ "$crd_status" == "$expected_status" ]; then
  echo "CRD status assertion passed: Status is as expected"
else
  echo "CRD status assertion failed: Actual: $crd_status, Expected: $expected_status"
  exit 1
fi
