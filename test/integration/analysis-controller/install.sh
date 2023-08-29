#!/bin/bash
current_time=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
initial_timestamp=$(date -u -d "$current_time" +%s)
seconds_to_add=20
new_timestamp=$((initial_timestamp + seconds_to_add))
time_20_seconds_later=$(date -u -d "@$new_timestamp" +"%Y-%m-%dT%H:%M:%SZ")
export CURRENT_TIME="$current_time"
export LATER="$time_20_seconds_later"
echo "templating time" $CURRENT_TIME $LATER
cp "01-install.yaml" tmp
envsubst <"01-install.yaml" >out.yaml
rm "01-install.yaml"
mv out.yaml "01-install.yaml"

