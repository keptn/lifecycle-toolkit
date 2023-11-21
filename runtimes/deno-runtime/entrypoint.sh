#!/bin/sh

set -eu

deno run --allow-net --allow-env=DATA,SECURE_DATA,KEPTN_CONTEXT "$SCRIPT"
