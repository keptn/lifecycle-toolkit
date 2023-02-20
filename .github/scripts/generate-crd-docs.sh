#!/bin/bash

API_ROOT='operator/apis'
TEMPLATE_DIR='./template'
RENDERER='markdown'
RENDERER_CONFIG_FILE='crd-render-config.yaml'

for api_group in "$API_ROOT"/*; do
  for api_version in "$API_ROOT/$api_group"/*; do
    crd-ref-docs \
      --templates-dir "$TEMPLATE_DIR" \
      --source-path="./$API_ROOT/$api_group/$api_version" \
      --renderer="$RENDERER" \
      --config "$RENDERER_CONFIG_FILE" \
      --output-path "./docs/content/en/docs/crd-ref/$api_group/$api_version/_index.md"
    done
done
