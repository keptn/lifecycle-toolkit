#!/bin/bash

# Readme generator for Keptn Helm Chart
#
# This script will install the readme generator if it's not installed already
# and then it will generate the readme from the local helm values
#
# Dependencies:
# Node >=16

# renovate: datasource=github-releases depName=bitnami-labs/readme-generator-for-helm
GENERATOR_VERSION="2.6.0"

echo "Checking if readme generator is installed already..."
if [[ $(npm list -g | grep -c "readme-generator-for-helm@${GENERATOR_VERSION}") -eq 0 ]]; then
  echo "Readme Generator v${GENERATOR_VERSION} not installed, installing now..."
  git clone https://github.com/bitnami-labs/readme-generator-for-helm.git
  cd ./readme-generator-for-helm || exit
  git checkout ${GENERATOR_VERSION}
  npm ci
  cd ..
  npm install -g ./readme-generator-for-helm
else
  echo "Readme Generator is already installed, continuing..."
fi

echo "Generating Keptn readme now..."
readme-generator --values=./chart/values.yaml --readme=./chart/README.md

echo "Generating lifecycle operator readme now..."
cat ./lifecycle-operator/chart/values.yaml ./lifecycle-operator/chart/doc.yaml > ./lifecycle-operator/chart/rendered.yaml
readme-generator --values=./lifecycle-operator/chart/rendered.yaml --readme=./lifecycle-operator/chart/README.md
rm ./lifecycle-operator/chart/rendered.yaml

echo "Generating keptn cert manager readme now..."
cat ./klt-cert-manager/chart/values.yaml ./klt-cert-manager/chart/doc.yaml > ./klt-cert-manager/chart/rendered.yaml
readme-generator --values=./klt-cert-manager/chart/rendered.yaml --readme=./klt-cert-manager/chart/README.md
rm ./klt-cert-manager/chart/rendered.yaml

echo "Generating keptn metrics operator readme now..."
cat ./metrics-operator/chart/values.yaml ./metrics-operator/chart/doc.yaml > ./metrics-operator/chart/rendered.yaml
readme-generator --values=./metrics-operator/chart/rendered.yaml --readme=./metrics-operator/chart/README.md
rm ./metrics-operator/chart/rendered.yaml

# Please be aware, the readme file needs to exist and needs to have a Parameters section, as only this section will be re-generated
