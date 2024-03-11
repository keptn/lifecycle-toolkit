#!/bin/bash

KEPTN_NAMESPACE="keptn-system"
RETRY_COUNT=5
SLEEP_TIME=5
# span name which we expect to find in the logs of otel-collector
EXPECTED_SUBSTRINGS=(
    Name : keptnapp-0.0.1-6b86b273
    Name : AppPreDeployTasks
    Name : AppPreDeployEvaluations
    Name : AppDeploy
    Name : keptnapp-nginx/WorkloadDeploy
    Name : keptnapp-nginx/WorkloadPreDeployTasks
    Name : keptnapp-nginx/WorkloadPreDeployEvaluations
    Name : keptnapp-nginx/WorkloadDeploy
    Name : keptnapp-nginx/WorkloadPostDeployTasks
    Name : keptnapp-nginx/WorkloadPostDeployEvaluations
    Name : AppPostDeployTasks
    Name : AppPostDeployEvaluations
    keptn.deployment.app.namespace: $NAMESPACE
)

check_variable_contains_substrings() {
    # remove unneeded whitespaces
    local variable2=$(echo "$1" | sed 's/\s\+/ /g')
    local variable=$(echo "$variable2" | sed "s/Str($NAMESPACE)/$NAMESPACE/g")
    local -a substrings=("${@:2}")

    local contains_all=true

    # Iterate over the substrings array
    for substring in "${substrings[@]}"; do
        # Check if the variable contains the current substring
        if [[ ! $variable =~ $substring ]]; then
            # If the substring is not found, set the flag to false and break the loop
            contains_all=false
            break
        fi
    done

    # Return the flag
    echo "$contains_all"
}

for i in $(seq 1 $RETRY_COUNT); do
    VAR=$(kubectl logs -n "$KEPTN_NAMESPACE" deployment/otel-collector)
    result=$(check_variable_contains_substrings "$VAR" "${EXPECTED_SUBSTRINGS[@]}")
    # shellcheck disable=SC1072
    if [ "$result" = true ]; then
        echo "All traces found, test passed!"
        exit 0
    fi
    if [ "$i" -lt "$RETRY_COUNT" ]; then
            echo "Sleeping for ${SLEEP_TIME} seconds before retrying..."
            sleep ${SLEEP_TIME}
    fi
done

echo "Retried ${RETRY_COUNT} times, but correct traces were not found. Exiting..."
exit 1
