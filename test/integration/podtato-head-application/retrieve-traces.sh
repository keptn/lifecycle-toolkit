#!/bin/bash

kubectl port-forward -n keptn-lifecycle-toolkit-system svc/jaeger-query 16686 &

RETRY_COUNT=10
SLEEP_TIME=5

for i in $(seq 1 $RETRY_COUNT); do
    # Retrieve the custom metric value
    TRACE_RESPONSE=$(curl -s "http://localhost:16686/api/traces?service=lifecycle-operator&limit=20&lookback=1h&operation=podtato-head-1.3-6b86b273")

    echo "$TRACE_RESPONSE"

    LENGTH_ITEMS=$(echo $TRACE_RESPONSE | jq '.data | length')

    if [[ $LENGTH_ITEMS -gt 0 ]]; then
        echo "Found exported traces"
        exit 0
    else
        echo "No exported traces found"
    fi

    if [ "$i" -lt "$RETRY_COUNT" ]; then
        echo "Sleeping for ${SLEEP_TIME} seconds before retrying..."
        sleep ${SLEEP_TIME}
    fi
done

echo "Retried ${RETRY_COUNT} times, but custom metric value did not meet the condition. Exiting..."
exit 1