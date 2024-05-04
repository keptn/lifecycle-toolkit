#!/bin/bash

performHelmRepoUpdate () {
    helm repo add keptn "https://charts.lifecycle.keptn.sh"
    helm repo update

    for chart_dir in ./lifecycle-operator/chart \
            ./metrics-operator/chart \
            ./keptn-cert-manager/chart \
            ./chart; do
        # shellcheck disable=SC2164
        cd "$chart_dir"
        echo "updating charts for" $chart_dir
        helm dependency update
        helm dependency build
        # shellcheck disable=SC2164
        cd -  # Return to the previous directory
    done
}
