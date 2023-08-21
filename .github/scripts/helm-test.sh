#!/bin/bash

# Keptn Lifecycle Toolkit Helm Testing
#
# This script supports the comparison of standard values and expected templated results to helm chart
# it is used to make sure changes to the chart are intentional and produce expected outcomes

echo "running Helm tests"
  tests=$(find ./.github/scripts/.helm-tests -maxdepth 1 -mindepth 1 -type d )

  errors=0
  successful=0
  failures=""

  for test in $tests
  do
    echo "Testing $test"
    helm template --namespace helmtests -f $test/values.yaml ./helm/chart > $test/helm_tests_output.yaml
    if [ $? -ne 0 ]
    then
      echo "Error: helm template failed for test in $test"
      errors=$((errors + 1))
      failures+="\n test $test failed"
    else
      diff -u "$test/helm_tests_output.yaml" "$test/result.yaml"
      if [ $? -ne 0 ]
      then
        echo "Error: test in $test not successful"
        errors=$((errors + 1))
        failures+="\n test $test failed"
      else
        echo "Info: test in $test successful"
        successful=$((successful + 1))
      fi
    fi
  done

  echo "run $((errors + successful)) tests: successful $successful, errors $errors"
  echo -e $failures
  if [ $errors -gt 0 ]
  then
    exit 1
  fi
