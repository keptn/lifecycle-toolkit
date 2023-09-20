#!/bin/bash

NAMESPACE="keptn-lifecycle-toolkit-system"
RETRY_COUNT=3
SLEEP_TIME=5

for i in $(seq 1 $RETRY_COUNT); do
    VAR=$(kubectl logs -n keptn-lifecycle-toolkit-system deployments/lifecycle-operator | grep -c "Error while parsing response")
    # shellcheck disable=SC1072
    if [ "$VAR" -ge 1 ]; then
      echo "Controller could access secret"
      exit 0
    fi
    if [ "$i" -lt "$RETRY_COUNT" ]; then
            echo "Sleeping for ${SLEEP_TIME} seconds before retrying..."
            sleep ${SLEEP_TIME}
    fi
done
echo "Retried ${RETRY_COUNT} times, but custom metric value did not meet the condition. Exiting..."exit 1
