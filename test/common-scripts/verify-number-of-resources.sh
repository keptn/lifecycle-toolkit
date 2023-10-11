#!/bin/bash

 RETRY_COUNT=3
 SLEEP_TIME=5
 RESOURCE_TYPE=$1
 EXPECTED_NUMBER=$2
 NS=$3

 for i in $(seq 1 $RETRY_COUNT); do
     NR_RESOURCES=$(kubectl get ${RESOURCE_TYPE} -n ${NS} --no-headers | wc -l)

     if [[ ${NR_RESOURCES} -eq ${EXPECTED_NUMBER} ]]; then
       echo "Found expected number of ${RESOURCE_TYPE}: ${EXPECTED_NUMBER}"
       exit 0
     else
       echo "Number of ${RESOURCE_TYPE} is not equal to ${EXPECTED_NUMBER}"
     fi

     if [ "$i" -lt "$RETRY_COUNT" ]; then
         echo "Sleeping for ${SLEEP_TIME} seconds before retrying..."
         sleep ${SLEEP_TIME}
     fi
 done

 echo "Retried ${RETRY_COUNT} times, but expected number of ${RESOURCE_TYPE} could not be verified. Exiting..."
 exit 1

