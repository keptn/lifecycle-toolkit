#!/bin/bash

METRIC_NAME=$1
RETRY_COUNT=3
SLEEP_TIME=5

for i in $(seq 1 $RETRY_COUNT); do
    # Retrieve the custom metric value
    METRIC_VALUE=$(kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/${NAMESPACE}/keptnmetrics.metrics.sh/${METRIC_NAME}/${METRIC_NAME}")

    LENGTH_ITEMS=$(echo $METRIC_VALUE | jq '.items | length')

    if [[ $LENGTH_ITEMS == 1 ]]; then
        echo "Found the expected metric $METRIC_NAME"
        exit 0
    else
        echo "The length of the property .items of $METRIC_NAME is not 1, it is: $LENGTH_ITEMS"
    fi

    if [ "$i" -lt "$RETRY_COUNT" ]; then
        echo "Sleeping for ${SLEEP_TIME} seconds before retrying..."
        sleep ${SLEEP_TIME}
    fi
done

echo "Retried ${RETRY_COUNT} times, but custom metric value did not meet the condition. Exiting..."
exit 1
