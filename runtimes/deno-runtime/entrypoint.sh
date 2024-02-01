#!/bin/sh

set -eu

deno run --allow-net --allow-write --allow-read --allow-env=DATA,SECURE_DATA,KEPTN_CONTEXT "$SCRIPT"
