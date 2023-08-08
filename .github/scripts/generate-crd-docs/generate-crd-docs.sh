#!/bin/bash

# CRD docs auto generation script
#
# This script goes through all API definitions in the lifecycle-operator/apis folder
# and generates docs from code for each API group and version
#
# Inputs: None

# renovate: datasource=github-releases depName=elastic/crd-ref-docs
GENERATOR_VERSION=v0.0.9
API_DOMAIN="keptn.sh"
OPERATOR_API_ROOT='lifecycle-operator/apis/'
METRICS_API_ROOT='metrics-operator/api/'
TEMPLATE_DIR='.github/scripts/generate-crd-docs/templates'
RENDERER='markdown'
RENDERER_CONFIG_FILE='.github/scripts/generate-crd-docs/crd-docs-generator-config.yaml'
PATH=$PATH:$(go env GOPATH)/bin

echo "Checking if code generator tool is installed..."
if ! test -s crd-ref-docs; then
  echo "Docs generator is not installed, installing now..."
  go install github.com/elastic/crd-ref-docs@${GENERATOR_VERSION}
else
  echo "Docs generator is installed, continuing..."
fi

echo "Running CRD docs auto-generator..."

for api_group in "$OPERATOR_API_ROOT"*; do

  sanitized_api_group="${api_group#$OPERATOR_API_ROOT}"
  INDEX_PATH="./docs/content/en/docs/crd-ref/$sanitized_api_group/_index.md"

  if [ ! -f "$INDEX_PATH" ]; then
    echo "API group index file doesn't exist for group $sanitized_api_group. Creating it now..."
    # Use sanitized_api_group and make first char uppercase
    API_GROUP="$(tr '[:lower:]' '[:upper:]' <<< "${sanitized_api_group:0:1}")${sanitized_api_group:1}"
    export API_GROUP
    envsubst < './.github/scripts/generate-crd-docs/templates/index-template.md' > "$INDEX_PATH"
    unset API_GROUP
  fi
  for api_version in "$api_group"/*; do
    sanitized_api_version="${api_version#$OPERATOR_API_ROOT$sanitized_api_group/}"

    OUTPUT_PATH="./docs/content/en/docs/crd-ref/$sanitized_api_group/$sanitized_api_version"

    echo "Arguments:"
    echo "TEMPLATE_DIR: $TEMPLATE_DIR"
    echo "OPERATOR_API_ROOT: $OPERATOR_API_ROOT"
    echo "API_GROUP: $sanitized_api_group"
    echo "API_VERSION: $sanitized_api_version"
    echo "RENDERER: $RENDERER"
    echo "RENDERER_CONFIG_FILE: $RENDERER_CONFIG_FILE"
    echo "OUTPUT_PATH: $OUTPUT_PATH/_index.md"

    echo "Creating docs folder $OUTPUT_PATH..."
    mkdir -p "$OUTPUT_PATH"

    echo "Generating CRD docs for $sanitized_api_group.$API_DOMAIN/$sanitized_api_version..."
    crd-ref-docs \
      --templates-dir "$TEMPLATE_DIR" \
      --source-path="./$api_version" \
      --renderer="$RENDERER" \
      --config "$RENDERER_CONFIG_FILE" \
      --output-path "$OUTPUT_PATH/_index.md"
    echo "---------------------"
  done
done

# Metrics API


sanitized_api_group="metrics"
INDEX_PATH="./docs/content/en/docs/crd-ref/$sanitized_api_group/_index.md"

if [ ! -f "$INDEX_PATH" ]; then
  echo "API group index file doesn't exist for group $sanitized_api_group. Creating it now..."
  # Use sanitized_api_group and make first char uppercase
  API_GROUP="$(tr '[:lower:]' '[:upper:]' <<< "${sanitized_api_group:0:1}")${sanitized_api_group:1}"
  export API_GROUP
  envsubst < './.github/scripts/generate-crd-docs/templates/index-template.md' > "$INDEX_PATH"
  unset API_GROUP
fi

for api_version in "$METRICS_API_ROOT"*; do
  sanitized_api_version="${api_version#$METRICS_API_ROOT}"
  OUTPUT_PATH="./docs/content/en/docs/crd-ref/$sanitized_api_group/$sanitized_api_version"

  echo "Arguments:"
  echo "TEMPLATE_DIR: $TEMPLATE_DIR"
  echo "METRICS_API_ROOT: $METRICS_API_ROOT"
  echo "API_GROUP: $sanitized_api_group"
  echo "API_VERSION: $sanitized_api_version"
  echo "RENDERER: $RENDERER"
  echo "RENDERER_CONFIG_FILE: $RENDERER_CONFIG_FILE"
  echo "OUTPUT_PATH: $OUTPUT_PATH/_index.md"

  echo "Creating docs folder $OUTPUT_PATH..."
  mkdir -p "$OUTPUT_PATH"

  echo "Generating CRD docs for $sanitized_api_group.$API_DOMAIN/$sanitized_api_version..."
  crd-ref-docs \
    --templates-dir "$TEMPLATE_DIR" \
    --source-path="./$api_version" \
    --renderer="$RENDERER" \
    --config "$RENDERER_CONFIG_FILE" \
    --output-path "$OUTPUT_PATH/_index.md"
  echo "---------------------"
done
