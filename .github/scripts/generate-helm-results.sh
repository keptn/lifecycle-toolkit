#!/bin/bash

echo "re-generating Helm test results"
tests=$(find ./.github/scripts/.helm-tests -maxdepth 1 -mindepth 1 -type d )

# shellcheck source=.github/scripts/helm-utils.sh
source .github/scripts/helm-utils.sh

performHelmRepoUpdate

for test in $tests
do
  echo "Re-generating $test"
  helm template keptn-test --namespace helmtests -f $test/values.yaml ./chart > $test/result.yaml
done
