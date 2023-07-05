#!/bin/bash

RETRY_COUNT=3
SLEEP_TIME=5

for i in $(seq 1 $RETRY_COUNT); do
    NR_JOBS=$(kubectl get jobs -n ${NAMESPACE} --no-headers | wc -l)

    if [[ ${NR_JOBS} -eq 1 ]]; then
      echo "Found exactly 1 job"
      exit 0
    else
      echo "Number of jobs is not equal to 1"
    fi

    if [ "$i" -lt "$RETRY_COUNT" ]; then
        echo "Sleeping for ${SLEEP_TIME} seconds before retrying..."
        sleep ${SLEEP_TIME}
    fi
done

echo "Retried ${RETRY_COUNT} times, but expected number of resources could not be verified. Exiting..."
exit 1
